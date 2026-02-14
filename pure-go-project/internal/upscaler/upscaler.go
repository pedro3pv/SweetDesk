package upscaler

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"

	"pure-go-project/internal/providers"
	"pure-go-project/internal/runtime"
)

type UpscalerType string

const (
	RealCUGAN UpscalerType = "realcugan"
	LSDIR     UpscalerType = "lsdir"
)

// UpscaleOptions define opções de processamento
type UpscaleOptions struct {
	// ScaleFactor: fator de escala (ex: 2.0, 3.0, 4.0)
	// Se 0, usa o scale padrão do modelo
	ScaleFactor float64

	// TargetWidth/TargetHeight: resolução alvo específica
	// Se definidos, ScaleFactor é ignorado
	TargetWidth  int
	TargetHeight int

	// MaxResolution: limite máximo (largura ou altura)
	// Se a saída exceder, será redimensionada proporcionalmente
	MaxResolution int

	// KeepAspectRatio: mantém proporção ao usar TargetWidth/Height
	KeepAspectRatio bool
}

type Upscaler struct {
	modelType    UpscalerType
	modelPath    string
	session      *ort.Session[float32]
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	inputName    string
	outputName   string
	tileSize     int
	modelScale   int // Scale nativo do modelo
}

func NewUpscaler(modelType UpscalerType, modelDir string, onnxLibPath string) (*Upscaler, error) {
	// Usa o singleton para inicializar o runtime
	if err := runtime.GetInstance().Initialize(onnxLibPath); err != nil {
		return nil, fmt.Errorf("falha ao inicializar runtime: %w", err)
	}

	u := &Upscaler{
		modelType:  modelType,
		tileSize:   512,
		modelScale: 4,
	}

	switch modelType {
	case RealCUGAN:
		u.modelPath = filepath.Join(modelDir, "realcugan", "realcugan-pro.onnx")
	case LSDIR:
		u.modelPath = filepath.Join(modelDir, "lsdir", "4xLSDIR.onnx")
	default:
		return nil, fmt.Errorf("tipo de upscaler desconhecido: %s", modelType)
	}

	if _, err := os.Stat(u.modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("modelo não encontrado: %s", u.modelPath)
	}

	ort.SetSharedLibraryPath(onnxLibPath)

	log.Printf("Inspecionando modelo: %s", u.modelPath)
	inputs, outputs, err := ort.GetInputOutputInfo(u.modelPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter info do modelo: %w", err)
	}

	if len(inputs) == 0 || len(outputs) == 0 {
		return nil, fmt.Errorf("modelo sem inputs ou outputs válidos")
	}

	u.inputName = inputs[0].Name
	u.outputName = outputs[0].Name

	log.Printf("  Input: %s %v", u.inputName, inputs[0].Dimensions)
	log.Printf("  Output: %s %v", u.outputName, outputs[0].Dimensions)

	if len(inputs[0].Dimensions) >= 4 {
		if inputs[0].Dimensions[2] > 0 {
			u.tileSize = int(inputs[0].Dimensions[2])
			log.Printf("  Tile size: %d", u.tileSize)
		}
	}

	if len(outputs[0].Dimensions) >= 4 && len(inputs[0].Dimensions) >= 4 {
		if outputs[0].Dimensions[2] > 0 && inputs[0].Dimensions[2] > 0 {
			u.modelScale = int(outputs[0].Dimensions[2] / inputs[0].Dimensions[2])
			log.Printf("  Scale nativo: %dx", u.modelScale)
		}
	}

	inputShape := ort.NewShape(1, 3, int64(u.tileSize), int64(u.tileSize))
	inputData := make([]float32, 1*3*u.tileSize*u.tileSize)
	u.inputTensor, err = ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar input tensor: %w", err)
	}

	outputSize := u.tileSize * u.modelScale
	outputShape := ort.NewShape(1, 3, int64(outputSize), int64(outputSize))
	u.outputTensor, err = ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		u.inputTensor.Destroy()
		return nil, fmt.Errorf("falha ao criar output tensor: %w", err)
	}

	// Detectar hardware disponível (usa singleton de providers)
	hwConfig, err := providers.DetectHardware()
	if err != nil {
		log.Printf("⚠️  Erro ao detectar hardware: %v", err)
		hwConfig = &providers.HardwareConfig{Provider: providers.ProviderCPU}
	}

	PrintHardwareInfo(hwConfig)

	// Criar SessionOptions
	sessionOptions, err := ort.NewSessionOptions()
	if err != nil {
		u.inputTensor.Destroy()
		u.outputTensor.Destroy()
		return nil, fmt.Errorf("falha ao criar session options: %w", err)
	}
	defer sessionOptions.Destroy()

	// Obter prioridade de providers a partir da detecção
	providerPriority := providers.GetProviderPriority(hwConfig)

	// Tentar criar sessão com fallback entre providers
	u.session, err = u.createSessionWithFallback(
		providerPriority,
		hwConfig,
		sessionOptions,
	)

	if err != nil {
		u.inputTensor.Destroy()
		u.outputTensor.Destroy()
		return nil, fmt.Errorf("falha ao criar sessão: %w", err)
	}

	return u, nil
}

// Upscale com opções padrão (4x)
func (u *Upscaler) Upscale(inputPath, outputPath string) error {
	return u.UpscaleWithOptions(inputPath, outputPath, nil)
}

// UpscaleWithOptions permite controle total da resolução
func (u *Upscaler) UpscaleWithOptions(inputPath, outputPath string, opts *UpscaleOptions) error {
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("falha ao abrir imagem: %w", err)
	}
	defer imgFile.Close()

	imgData, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("falha ao decodificar imagem: %w", err)
	}

	bounds := imgData.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calcula resolução de saída
	targetWidth, targetHeight := u.calculateOutputSize(width, height, opts)

	log.Printf("📐 Entrada: %dx%d", width, height)
	log.Printf("🎯 Saída alvo: %dx%d", targetWidth, targetHeight)

	// Primeiro faz upscale com o modelo
	var upscaled image.Image
	if width <= u.tileSize && height <= u.tileSize {
		upscaled, err = u.processTile(imgData)
		if err != nil {
			return err
		}
	} else {
		upscaled, err = u.processTiled(imgData)
		if err != nil {
			return err
		}
	}

	// Ajusta para resolução final se necessário
	upscaledBounds := upscaled.Bounds()
	if upscaledBounds.Dx() != targetWidth || upscaledBounds.Dy() != targetHeight {
		log.Printf("🔧 Ajustando resolução final...")
		final := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
		draw.BiLinear.Scale(final, final.Bounds(), upscaled, upscaledBounds, draw.Over, nil)
		upscaled = final
	}

	return u.saveImage(upscaled, outputPath)
}

// calculateOutputSize determina a resolução de saída
func (u *Upscaler) calculateOutputSize(inputWidth, inputHeight int, opts *UpscaleOptions) (int, int) {
	// Usa opções padrão se não fornecidas
	if opts == nil {
		opts = &UpscaleOptions{ScaleFactor: float64(u.modelScale)}
	}

	var targetWidth, targetHeight int

	// Prioridade 1: Resolução específica
	if opts.TargetWidth > 0 || opts.TargetHeight > 0 {
		if opts.KeepAspectRatio {
			// Mantém proporção
			if opts.TargetWidth > 0 && opts.TargetHeight > 0 {
				// Ambos definidos: usa o que couber mantendo proporção
				scaleW := float64(opts.TargetWidth) / float64(inputWidth)
				scaleH := float64(opts.TargetHeight) / float64(inputHeight)
				scale := math.Min(scaleW, scaleH)
				targetWidth = int(float64(inputWidth) * scale)
				targetHeight = int(float64(inputHeight) * scale)
			} else if opts.TargetWidth > 0 {
				// Só largura: calcula altura
				scale := float64(opts.TargetWidth) / float64(inputWidth)
				targetWidth = opts.TargetWidth
				targetHeight = int(float64(inputHeight) * scale)
			} else {
				// Só altura: calcula largura
				scale := float64(opts.TargetHeight) / float64(inputHeight)
				targetWidth = int(float64(inputWidth) * scale)
				targetHeight = opts.TargetHeight
			}
		} else {
			// Força dimensões exatas (pode distorcer)
			targetWidth = opts.TargetWidth
			targetHeight = opts.TargetHeight
			if targetWidth == 0 {
				targetWidth = inputWidth * u.modelScale
			}
			if targetHeight == 0 {
				targetHeight = inputHeight * u.modelScale
			}
		}
	} else if opts.ScaleFactor > 0 {
		// Prioridade 2: Fator de escala
		targetWidth = int(float64(inputWidth) * opts.ScaleFactor)
		targetHeight = int(float64(inputHeight) * opts.ScaleFactor)
	} else {
		// Padrão: usa scale do modelo
		targetWidth = inputWidth * u.modelScale
		targetHeight = inputHeight * u.modelScale
	}

	// Aplica limite máximo se definido
	if opts.MaxResolution > 0 {
		if targetWidth > opts.MaxResolution || targetHeight > opts.MaxResolution {
			if targetWidth > targetHeight {
				scale := float64(opts.MaxResolution) / float64(targetWidth)
				targetWidth = opts.MaxResolution
				targetHeight = int(float64(targetHeight) * scale)
			} else {
				scale := float64(opts.MaxResolution) / float64(targetHeight)
				targetHeight = opts.MaxResolution
				targetWidth = int(float64(targetWidth) * scale)
			}
			log.Printf("⚠️  Limitando a %dx%d (MaxResolution=%d)", targetWidth, targetHeight, opts.MaxResolution)
		}
	}

	return targetWidth, targetHeight
}

func (u *Upscaler) processTile(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	resized := image.NewRGBA(image.Rect(0, 0, u.tileSize, u.tileSize))
	draw.BiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	data := u.inputTensor.GetData()
	for y := 0; y < u.tileSize; y++ {
		for x := 0; x < u.tileSize; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			data[0*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(r>>8) / 255.0
			data[1*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(g>>8) / 255.0
			data[2*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(b>>8) / 255.0
		}
	}

	if err := u.session.Run(); err != nil {
		return nil, fmt.Errorf("falha na inferência: %w", err)
	}

	outputData := u.outputTensor.GetData()
	outputSize := u.tileSize * u.modelScale
	result := image.NewRGBA(image.Rect(0, 0, outputSize, outputSize))

	for y := 0; y < outputSize; y++ {
		for x := 0; x < outputSize; x++ {
			r := uint8(clamp(outputData[0*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			g := uint8(clamp(outputData[1*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			b := uint8(clamp(outputData[2*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			result.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	if width < u.tileSize || height < u.tileSize {
		finalWidth := width * u.modelScale
		finalHeight := height * u.modelScale
		final := image.NewRGBA(image.Rect(0, 0, finalWidth, finalHeight))
		draw.BiLinear.Scale(final, final.Bounds(), result, result.Bounds(), draw.Over, nil)
		return final, nil
	}

	return result, nil
}

func (u *Upscaler) processTiled(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	outputWidth := width * u.modelScale
	outputHeight := height * u.modelScale
	result := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))

	totalTiles := ((width + u.tileSize - 1) / u.tileSize) * ((height + u.tileSize - 1) / u.tileSize)
	currentTile := 0

	for y := 0; y < height; y += u.tileSize {
		for x := 0; x < width; x += u.tileSize {
			currentTile++
			tileWidth := min(u.tileSize, width-x)
			tileHeight := min(u.tileSize, height-y)

			log.Printf("  Tile %d/%d", currentTile, totalTiles)

			tile := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
			draw.Draw(tile, tile.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

			processedTile, err := u.processTile(tile)
			if err != nil {
				return nil, fmt.Errorf("falha ao processar tile: %w", err)
			}

			draw.Draw(result,
				image.Rect(x*u.modelScale, y*u.modelScale, (x+tileWidth)*u.modelScale, (y+tileHeight)*u.modelScale),
				processedTile,
				image.Point{X: 0, Y: 0},
				draw.Src)
		}
	}

	return result, nil
}

func (u *Upscaler) saveImage(img image.Image, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório: %w", err)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de saída: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("falha ao codificar imagem: %w", err)
	}

	bounds := img.Bounds()
	log.Printf("✅ Salva: %s (%dx%d)", path, bounds.Dx(), bounds.Dy())
	return nil
}

func (u *Upscaler) Close() {
	if u.session != nil {
		u.session.Destroy()
	}
	if u.inputTensor != nil {
		u.inputTensor.Destroy()
	}
	if u.outputTensor != nil {
		u.outputTensor.Destroy()
	}
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HardwareConfig descreve capacidades relevantes para escolher providers
// ProviderTipo identifica o provider preferencial
type ProviderTipo string

const (
	ProviderCPU  ProviderTipo = "CPU"
	ProviderCUDA ProviderTipo = "CUDA"
)

// HardwareConfig descreve capacidades relevantes para escolher providers
type HardwareConfig struct {
	Provider    ProviderTipo
	UseCUDA     bool
	DeviceID    int
	MemoryLimit int64
}

// DetectHardware detecta (de forma simples) se CUDA está disponível.
// Atualmente é uma implementação conservadora que pode ser estendida.
func DetectHardware() (*HardwareConfig, error) {
	// Implementação simples: desabilita CUDA por padrão.
	// Futuramente podemos checar bibliotecas, variáveis de ambiente ou usar cuInit.
	return &HardwareConfig{Provider: ProviderCPU, UseCUDA: false, DeviceID: 0, MemoryLimit: 0}, nil
}

// PrintHardwareInfo mostra um resumo do hardware detectado
func PrintHardwareInfo(h *providers.HardwareConfig) {
	if h == nil {
		log.Println("Hardware: nenhum dado disponível")
		return
	}
	log.Printf("Hardware: Provider=%s DeviceID=%d MemLimit=%d", h.Provider, h.DeviceID, h.MemoryLimit)
}

func (u *Upscaler) createSessionWithFallback(
	provList []providers.ProviderType,
	hwConfig *providers.HardwareConfig,
	sessionOptions *ort.SessionOptions,
) (*ort.Session[float32], error) {
	var lastErr error

	for _, provider := range provList {
		log.Printf("🔧 Tentando provider: %s", provider)

		// Configurar provider específico
		if provider == providers.ProviderCoreML {
			// CoreML requer options
			err := sessionOptions.AppendExecutionProvider(string(provider), hwConfig.Options)
			if err != nil {
				log.Printf("⚠️  Falha ao configurar CoreML: %v", err)
				lastErr = err
				continue
			}
		} else if provider == providers.ProviderCUDA {
			// CUDA requer options
			err := sessionOptions.AppendExecutionProvider(string(provider), hwConfig.Options)
			if err != nil {
				log.Printf("⚠️  Falha ao configurar CUDA: %v", err)
				lastErr = err
				continue
			}
		}
		// CPU não precisa de configuração especial

		// Tentar criar sessão
		session, err := ort.NewSession[float32](
			u.modelPath,
			[]string{u.inputName},
			[]string{u.outputName},
			[]*ort.Tensor[float32]{u.inputTensor},
			[]*ort.Tensor[float32]{u.outputTensor},
		)

		if err == nil {
			log.Printf("✅ Usando %s com sucesso!", provider)
			return session, nil
		}

		log.Printf("❌ %s falhou: %v", provider, err)
		lastErr = err
	}

	return nil, fmt.Errorf("todos os providers falharam. Último erro: %w", lastErr)
}

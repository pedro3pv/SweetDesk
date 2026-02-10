package services

import (
	"bytes"
	"context"
	"embed"
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

	// Format: formato de saída ("png", "jpg")
	Format string
}

type Upscaler struct {
	ctx          context.Context
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

// NewUpscaler cria um novo upscaler com modelo embutido
func NewUpscaler(ctx context.Context, modelType UpscalerType, modelsFS embed.FS, onnxLibPath string) (*Upscaler, error) {
	u := &Upscaler{
		ctx:        ctx,
		modelType:  modelType,
		tileSize:   512,
		modelScale: 4,
	}

	var modelPath string
	switch modelType {
	case RealCUGAN:
		modelPath = "models/realcugan/realcugan-pro.onnx"
	case LSDIR:
		modelPath = "models/lsdir/4xLSDIR.onnx"
	default:
		return nil, fmt.Errorf("tipo de upscaler desconhecido: %s", modelType)
	}

	// Lê modelo da memória (embed.FS)
	modelData, err := modelsFS.ReadFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("modelo não encontrado no embed: %w", err)
	}

	ort.SetSharedLibraryPath(onnxLibPath)

	// Valores hardcoded (ajuste conforme seus modelos)
	u.inputName = "input"
	u.outputName = "output"

	// Cria tensores
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

	// Cria sessão ONNX a partir dos bytes
	u.session, err = ort.NewSessionWithONNXData[float32](
		modelData,
		[]string{u.inputName},
		[]string{u.outputName},
		[]*ort.Tensor[float32]{u.inputTensor},
		[]*ort.Tensor[float32]{u.outputTensor},
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

// UpscaleWithOptions permite controle total da resolução (arquivo -> arquivo)
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

	result, err := u.upscaleImage(imgData, opts)
	if err != nil {
		return err
	}

	return u.saveImage(result, outputPath)
}

// UpscaleBytes processa imagem a partir de bytes e retorna bytes
func (u *Upscaler) UpscaleBytes(imageData []byte, opts *UpscaleOptions) ([]byte, error) {
	// Decodifica imagem
	imgData, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar imagem: %w", err)
	}

	// Processa
	result, err := u.upscaleImage(imgData, opts)
	if err != nil {
		return nil, err
	}

	// Codifica resultado
	var buf bytes.Buffer
	format := "png"
	if opts != nil && opts.Format != "" {
		format = opts.Format
	}

	switch format {
	case "png":
		if err := png.Encode(&buf, result); err != nil {
			return nil, fmt.Errorf("falha ao codificar PNG: %w", err)
		}
	default:
		if err := png.Encode(&buf, result); err != nil {
			return nil, fmt.Errorf("falha ao codificar imagem: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// upscaleImage - lógica interna de processamento
func (u *Upscaler) upscaleImage(imgData image.Image, opts *UpscaleOptions) (image.Image, error) {
	// Verifica contexto
	select {
	case <-u.ctx.Done():
		return nil, u.ctx.Err()
	default:
	}

	bounds := imgData.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calcula resolução de saída
	targetWidth, targetHeight := u.calculateOutputSize(width, height, opts)
	log.Printf("📐 Entrada: %dx%d", width, height)
	log.Printf("🎯 Saída alvo: %dx%d", targetWidth, targetHeight)

	// Upscale com o modelo
	var upscaled image.Image
	var err error

	if width <= u.tileSize && height <= u.tileSize {
		upscaled, err = u.processTile(imgData)
		if err != nil {
			return nil, err
		}
	} else {
		upscaled, err = u.processTiled(imgData)
		if err != nil {
			return nil, err
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

	return upscaled, nil
}

// calculateOutputSize determina a resolução de saída
func (u *Upscaler) calculateOutputSize(inputWidth, inputHeight int, opts *UpscaleOptions) (int, int) {
	if opts == nil {
		opts = &UpscaleOptions{ScaleFactor: float64(u.modelScale)}
	}

	var targetWidth, targetHeight int

	// Prioridade 1: Resolução específica
	if opts.TargetWidth > 0 || opts.TargetHeight > 0 {
		if opts.KeepAspectRatio {
			if opts.TargetWidth > 0 && opts.TargetHeight > 0 {
				scaleW := float64(opts.TargetWidth) / float64(inputWidth)
				scaleH := float64(opts.TargetHeight) / float64(inputHeight)
				scale := math.Min(scaleW, scaleH)
				targetWidth = int(float64(inputWidth) * scale)
				targetHeight = int(float64(inputHeight) * scale)
			} else if opts.TargetWidth > 0 {
				scale := float64(opts.TargetWidth) / float64(inputWidth)
				targetWidth = opts.TargetWidth
				targetHeight = int(float64(inputHeight) * scale)
			} else {
				scale := float64(opts.TargetHeight) / float64(inputHeight)
				targetWidth = int(float64(inputWidth) * scale)
				targetHeight = opts.TargetHeight
			}
		} else {
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
			log.Printf("⚠️ Limitando a %dx%d (MaxResolution=%d)", targetWidth, targetHeight, opts.MaxResolution)
		}
	}

	return targetWidth, targetHeight
}

func (u *Upscaler) processTile(img image.Image) (image.Image, error) {
	// Verifica contexto
	select {
	case <-u.ctx.Done():
		return nil, u.ctx.Err()
	default:
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Redimensiona para tileSize
	resized := image.NewRGBA(image.Rect(0, 0, u.tileSize, u.tileSize))
	draw.BiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	// Converte para tensor
	data := u.inputTensor.GetData()
	for y := 0; y < u.tileSize; y++ {
		for x := 0; x < u.tileSize; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			data[0*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(r>>8) / 255.0
			data[1*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(g>>8) / 255.0
			data[2*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(b>>8) / 255.0
		}
	}

	// Inferência
	if err := u.session.Run(); err != nil {
		return nil, fmt.Errorf("falha na inferência: %w", err)
	}

	// Converte resultado
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

	// Ajusta se entrada era menor que tileSize
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
			// Verifica cancelamento
			select {
			case <-u.ctx.Done():
				return nil, u.ctx.Err()
			default:
			}

			currentTile++
			tileWidth := min(u.tileSize, width-x)
			tileHeight := min(u.tileSize, height-y)

			log.Printf("  Tile %d/%d", currentTile, totalTiles)

			// Extrai tile
			tile := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
			draw.Draw(tile, tile.Bounds(), img, image.Point{X: x, Y: y}, draw.Src)

			// Processa tile
			processedTile, err := u.processTile(tile)
			if err != nil {
				return nil, fmt.Errorf("falha ao processar tile: %w", err)
			}

			// Cola no resultado
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

// GetRecommendedModel retorna o modelo recomendado baseado no tipo de imagem
func GetRecommendedModel(imageType string) UpscalerType {
	switch imageType {
	case "anime":
		return RealCUGAN
	case "photo":
		return LSDIR
	default:
		return RealCUGAN
	}
}

// Close libera recursos
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

// Helper functions
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

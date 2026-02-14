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
	"sync"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"golang.org/x/image/draw"
	"gorgonia.org/tensor"
)

type UpscalerType string

const (
	RealCUGAN UpscalerType = "realcugan"
	LSDIR     UpscalerType = "lsdir"
)

type UpscaleOptions struct {
	ScaleFactor     float64
	TargetWidth     int
	TargetHeight    int
	MaxResolution   int
	KeepAspectRatio bool
	Format          string
}

type Upscaler struct {
	ctx        context.Context
	modelType  UpscalerType
	modelData  []byte // cached ONNX model bytes (immutable after init)
	tileSize   int
	modelScale int
	mu         sync.Mutex // serializes inference to limit memory usage
}

// NewUpscaler cria upscaler com onnx-go (pure Go, sem DLL)
func NewUpscaler(ctx context.Context, modelType UpscalerType, modelsFS embed.FS) (*Upscaler, error) {
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

	// Lê modelo do embed e valida uma vez
	modelData, err := modelsFS.ReadFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("modelo não encontrado: %w", err)
	}

	// Valida modelo ONNX na inicialização
	testBackend := gorgonnx.NewGraph()
	testModel := onnx.NewModel(testBackend)
	if err := testModel.UnmarshalBinary(modelData); err != nil {
		return nil, fmt.Errorf("falha ao carregar modelo: %w", err)
	}

	u.modelData = modelData
	log.Printf("✅ Modelo %s carregado (pure Go, sem DLL)", modelType)
	return u, nil
}

func (u *Upscaler) Upscale(inputPath, outputPath string) error {
	return u.UpscaleWithOptions(inputPath, outputPath, nil)
}

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

func (u *Upscaler) UpscaleBytes(imageData []byte, opts *UpscaleOptions) ([]byte, error) {
	imgData, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar imagem: %w", err)
	}

	result, err := u.upscaleImage(imgData, opts)
	if err != nil {
		return nil, err
	}

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

func (u *Upscaler) upscaleImage(imgData image.Image, opts *UpscaleOptions) (image.Image, error) {
	select {
	case <-u.ctx.Done():
		return nil, u.ctx.Err()
	default:
	}

	bounds := imgData.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	targetWidth, targetHeight := u.calculateOutputSize(width, height, opts)
	log.Printf("📐 Entrada: %dx%d", width, height)
	log.Printf("🎯 Saída alvo: %dx%d", targetWidth, targetHeight)

	var upscaled image.Image
	var err error

	if width <= u.tileSize && height <= u.tileSize {
		upscaled, err = u.processTile(imgData)
	} else {
		upscaled, err = u.processTiled(imgData)
	}

	if err != nil {
		return nil, err
	}

	// Ajusta resolução final se necessário
	upscaledBounds := upscaled.Bounds()
	if upscaledBounds.Dx() != targetWidth || upscaledBounds.Dy() != targetHeight {
		log.Printf("🔧 Ajustando resolução final...")
		final := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
		draw.BiLinear.Scale(final, final.Bounds(), upscaled, upscaledBounds, draw.Over, nil)
		upscaled = final
	}

	return upscaled, nil
}

func (u *Upscaler) processTile(img image.Image) (image.Image, error) {
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

	// Converte para tensor Gorgonia [1, 3, H, W]
	inputData := make([]float32, 1*3*u.tileSize*u.tileSize)
	for y := 0; y < u.tileSize; y++ {
		for x := 0; x < u.tileSize; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			inputData[0*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(r>>8) / 255.0
			inputData[1*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(g>>8) / 255.0
			inputData[2*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(b>>8) / 255.0
		}
	}

	inputTensor := tensor.New(
		tensor.WithShape(1, 3, u.tileSize, u.tileSize),
		tensor.WithBacking(inputData),
	)

	// Serialize inference calls to limit memory usage
	u.mu.Lock()
	defer u.mu.Unlock()

	// Cria backend+model frescos a cada inferência para evitar o bug do
	// TapeMachine.Reset() na onnx-go v0.5.0 que retorna typed-nil error
	// (Go nil-interface problem: interface não-nil com valor concreto nil).
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	if err := model.UnmarshalBinary(u.modelData); err != nil {
		return nil, fmt.Errorf("falha ao recriar modelo: %w", err)
	}

	// Set input no modelo
	if err := model.SetInput(0, inputTensor); err != nil {
		return nil, fmt.Errorf("falha ao setar input: %w", err)
	}

	// Executa inferência
	if err := backend.Run(); err != nil {
		return nil, fmt.Errorf("falha na inferência: %w", err)
	}

	// Obtém output
	outputs, err := model.GetOutputTensors()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter output: %w", err)
	}

	outputTensor := outputs[0]
	outputData := outputTensor.Data().([]float32)

	// Converte tensor para imagem
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

	// Ajusta se entrada era menor
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
			select {
			case <-u.ctx.Done():
				return nil, u.ctx.Err()
			default:
			}

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

func (u *Upscaler) calculateOutputSize(inputWidth, inputHeight int, opts *UpscaleOptions) (int, int) {
	if opts == nil {
		opts = &UpscaleOptions{ScaleFactor: float64(u.modelScale)}
	}

	var targetWidth, targetHeight int

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
		targetWidth = int(float64(inputWidth) * opts.ScaleFactor)
		targetHeight = int(float64(inputHeight) * opts.ScaleFactor)
	} else {
		targetWidth = inputWidth * u.modelScale
		targetHeight = inputHeight * u.modelScale
	}

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

func (u *Upscaler) saveImage(img image.Image, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório: %w", err)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("falha ao codificar: %w", err)
	}

	bounds := img.Bounds()
	log.Printf("✅ Salva: %s (%dx%d)", path, bounds.Dx(), bounds.Dy())
	return nil
}

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

func (u *Upscaler) Close() {
	// Gorgonia não precisa de cleanup explícito
	log.Println("✅ Upscaler fechado")
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

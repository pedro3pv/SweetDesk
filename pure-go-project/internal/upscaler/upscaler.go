package upscaler

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

// UpscalerType define o tipo de upscaler
type UpscalerType string

const (
	RealCUGAN UpscalerType = "realcugan"
	LSDIR     UpscalerType = "lsdir"
)

// Upscaler processa imagens com super-resolução
type Upscaler struct {
	modelType    UpscalerType
	modelPath    string
	session      *ort.Session[float32]
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	tileSize     int
	scale        int
}

// NewUpscaler cria um novo upscaler
func NewUpscaler(modelType UpscalerType, modelDir string, onnxLibPath string) (*Upscaler, error) {
	u := &Upscaler{
		modelType: modelType,
		tileSize:  512, // Tamanho do tile para processamento
		scale:     4,   // Fator de escala (4x)
	}

	// Define caminho do modelo baseado no tipo
	switch modelType {
	case RealCUGAN:
		u.modelPath = filepath.Join(modelDir, "realcugan", "realcugan-pro.onnx")
	case LSDIR:
		u.modelPath = filepath.Join(modelDir, "lsdir", "4xLSDIR.onnx")
	default:
		return nil, fmt.Errorf("tipo de upscaler desconhecido: %s", modelType)
	}

	// Verifica se modelo existe
	if _, err := os.Stat(u.modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("modelo não encontrado: %s", u.modelPath)
	}

	// Inicializa ONNX Runtime (se ainda não foi)
	ort.SetSharedLibraryPath(onnxLibPath)

	// Cria tensors (ajustar dimensões conforme o modelo)
	inputShape := ort.NewShape(int64(1), int64(3), int64(u.tileSize), int64(u.tileSize))
	inputData := make([]float32, 1*3*u.tileSize*u.tileSize)
	var err error
	u.inputTensor, err = ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar input tensor: %w", err)
	}

	outputSize := u.tileSize * u.scale
	outputShape := ort.NewShape(int64(1), int64(3), int64(outputSize), int64(outputSize))
	u.outputTensor, err = ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		u.inputTensor.Destroy()
		return nil, fmt.Errorf("falha ao criar output tensor: %w", err)
	}

	// Cria sessão
	u.session, err = ort.NewSession[float32](
		u.modelPath,
		[]string{"input"},
		[]string{"output"},
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

// Upscale processa uma imagem e retorna versão em alta resolução
func (u *Upscaler) Upscale(inputPath, outputPath string) error {
	// Carrega imagem
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

	// Se a imagem for menor que o tile size, processa inteira
	if width <= u.tileSize && height <= u.tileSize {
		result, err := u.processTile(imgData)
		if err != nil {
			return err
		}
		return u.saveImage(result, outputPath)
	}

	// Caso contrário, processa em tiles (para imagens grandes)
	return u.processTiled(imgData, outputPath)
}

// processTile processa uma única região da imagem
func (u *Upscaler) processTile(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Resize para tile size se necessário
	resized := image.NewRGBA(image.Rect(0, 0, u.tileSize, u.tileSize))
	draw.BiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	// Normaliza e converte para CHW
	data := u.inputTensor.GetData()
	for y := 0; y < u.tileSize; y++ {
		for x := 0; x < u.tileSize; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Normalização: [0, 255] -> [0, 1]
			data[0*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(r>>8) / 255.0
			data[1*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(g>>8) / 255.0
			data[2*u.tileSize*u.tileSize+y*u.tileSize+x] = float32(b>>8) / 255.0
		}
	}

	// Executa inferência
	if err := u.session.Run(); err != nil {
		return nil, fmt.Errorf("falha na inferência: %w", err)
	}

	// Converte output para imagem
	outputData := u.outputTensor.GetData()
	outputSize := u.tileSize * u.scale
	result := image.NewRGBA(image.Rect(0, 0, outputSize, outputSize))

	for y := 0; y < outputSize; y++ {
		for x := 0; x < outputSize; x++ {
			r := uint8(clamp(outputData[0*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			g := uint8(clamp(outputData[1*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			b := uint8(clamp(outputData[2*outputSize*outputSize+y*outputSize+x]*255.0, 0, 255))
			result.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	// Se a imagem original era menor, redimensiona o resultado
	if width < u.tileSize || height < u.tileSize {
		finalWidth := width * u.scale
		finalHeight := height * u.scale
		final := image.NewRGBA(image.Rect(0, 0, finalWidth, finalHeight))
		draw.BiLinear.Scale(final, final.Bounds(), result, result.Bounds(), draw.Over, nil)
		return final, nil
	}

	return result, nil
}

// processTiled processa imagem grande em tiles
func (u *Upscaler) processTiled(img image.Image, outputPath string) error {
	// Implementação simplificada - divide em tiles e processa cada um
	// Para produção, adicionar overlap e blending entre tiles
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	outputWidth := width * u.scale
	outputHeight := height * u.scale
	result := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))

	// Divide em tiles
	for y := 0; y < height; y += u.tileSize {
		for x := 0; x < width; x += u.tileSize {
			tileWidth := min(u.tileSize, width-x)
			tileHeight := min(u.tileSize, height-y)

			// Extrai tile
			tile := image.NewRGBA(image.Rect(0, 0, tileWidth, tileHeight))
			draw.Draw(tile, tile.Bounds(), img, image.Point{x, y}, draw.Src)

			// Processa tile
			processedTile, err := u.processTile(tile)
			if err != nil {
				return fmt.Errorf("falha ao processar tile (%d,%d): %w", x, y, err)
			}

			// Cola tile processado no resultado
			draw.Draw(result,
				image.Rect(x*u.scale, y*u.scale, (x+tileWidth)*u.scale, (y+tileHeight)*u.scale),
				processedTile,
				image.Point{0, 0},
				draw.Src)
		}
	}

	return u.saveImage(result, outputPath)
}

// saveImage salva a imagem processada
func (u *Upscaler) saveImage(img image.Image, path string) error {
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de saída: %w", err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return fmt.Errorf("falha ao codificar imagem: %w", err)
	}

	return nil
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

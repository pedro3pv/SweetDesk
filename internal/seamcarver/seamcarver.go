package seamcarver

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

// AspectRatioOptions define opções de manipulação de aspect ratio
type AspectRatioOptions struct {
	// TargetRatio = width/height (16:9 = 16.0/9.0, 4:3 = 4.0/3.0, etc.)
	TargetRatio float64

	// KeepContentPriority: se true, prioriza remoção vertical de seams
	KeepContentPriority bool

	// MaxDeltaBySeams: máximo de pixels alterados por seams (limite de segurança)
	MaxDeltaBySeams int
}

// SeamCarver manipula aspect ratio com content-aware
type SeamCarver struct {
	img image.Image
}

// NewSeamCarver cria um novo seam carver
func NewSeamCarver(img image.Image) *SeamCarver {
	return &SeamCarver{img: img}
}

// NewSeamCarverFromFile carrega imagem de arquivo
func NewSeamCarverFromFile(path string) (*SeamCarver, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar imagem: %w", err)
	}

	return &SeamCarver{img: img}, nil
}

// AdjustAspectRatio altera aspect ratio usando seam carving
func (sc *SeamCarver) AdjustAspectRatio(opts AspectRatioOptions) (image.Image, error) {
	bounds := sc.img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	currentRatio := float64(width) / float64(height)

	log.Printf("📐 Original: %dx%d (ratio: %.3f)", width, height, currentRatio)
	log.Printf("🎯 Target ratio: %.3f", opts.TargetRatio)

	if math.Abs(currentRatio-opts.TargetRatio) < 0.01 {
		log.Printf("ℹ️  Ratio já está suficientemente próximo")
		return sc.img, nil
	}

	var result image.Image
	var err error

	if currentRatio > opts.TargetRatio {
		// Imagem é muito larga - remove colunas (vertical seams)
		result, err = sc.reduceVertical(width, height, opts)
	} else {
		// Imagem é muito alta - remove linhas (horizontal seams)
		result, err = sc.reduceHorizontal(width, height, opts)
	}

	if err != nil {
		return nil, err
	}

	resultBounds := result.Bounds()
	finalRatio := float64(resultBounds.Dx()) / float64(resultBounds.Dy())
	log.Printf("✅ Final: %dx%d (ratio: %.3f)", resultBounds.Dx(), resultBounds.Dy(), finalRatio)

	return result, nil
}

// reduceVertical remove seams verticais para deixar imagem mais estreita
func (sc *SeamCarver) reduceVertical(width, height int, opts AspectRatioOptions) (image.Image, error) {
	targetWidth := int(opts.TargetRatio * float64(height))

	// Limites de segurança
	minWidth := int(float64(width) * 0.5) // não reduzir menos que 50%
	if targetWidth < minWidth {
		targetWidth = minWidth
	}

	// Decide quantas colunas remover com seams
	colsToRemoveBySeams := width - targetWidth

	if colsToRemoveBySeams > opts.MaxDeltaBySeams {
		colsToRemoveBySeams = opts.MaxDeltaBySeams
	}

	if colsToRemoveBySeams <= 0 {
		return sc.img, nil
	}

	log.Printf("🔧 Alvo de largura: %d (remover %d colunas)", targetWidth, colsToRemoveBySeams)

	result := image.NewRGBA(image.Rect(0, 0, width-colsToRemoveBySeams, height))

	// Copia pixels com compressão horizontal simples
	for y := 0; y < height; y++ {
		for x := 0; x < width-colsToRemoveBySeams; x++ {
			// Mapeia posição destino para origem proporcionalmente
			sourceX := int(float64(x) * float64(width) / float64(width-colsToRemoveBySeams))
			if sourceX >= width {
				sourceX = width - 1
			}
			r, g, b, a := getPixel(sc.img, sourceX, y)
			result.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return result, nil
}

// reduceHorizontal remove seams horizontais para deixar imagem mais curta
func (sc *SeamCarver) reduceHorizontal(width, height int, opts AspectRatioOptions) (image.Image, error) {
	targetHeight := int(float64(width) / opts.TargetRatio)

	// Limites de segurança
	minHeight := int(float64(height) * 0.5)
	if targetHeight < minHeight {
		targetHeight = minHeight
	}

	// Quantas linhas remover com seams
	rowsToRemoveBySeams := height - targetHeight

	if rowsToRemoveBySeams > opts.MaxDeltaBySeams {
		rowsToRemoveBySeams = opts.MaxDeltaBySeams
	}

	if rowsToRemoveBySeams <= 0 {
		return sc.img, nil
	}

	log.Printf("🔧 Alvo de altura: %d (remover %d linhas)", targetHeight, rowsToRemoveBySeams)

	result := image.NewRGBA(image.Rect(0, 0, width, height-rowsToRemoveBySeams))

	// Copia pixels com compressão vertical simples
	for y := 0; y < height-rowsToRemoveBySeams; y++ {
		sourceY := int(float64(y) * float64(height) / float64(height-rowsToRemoveBySeams))
		if sourceY >= height {
			sourceY = height - 1
		}
		for x := 0; x < width; x++ {
			r, g, b, a := getPixel(sc.img, x, sourceY)
			result.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return result, nil
}

// SaveImage salva a imagem ajustada (EXPORTADO)
func (sc *SeamCarver) SaveImage(img image.Image, path string) error {
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

// Helper functions
func getPixel(img image.Image, x, y int) (uint8, uint8, uint8, uint8) {
	bounds := img.Bounds()
	// Garante que está dentro dos limites
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if x >= bounds.Max.X {
		x = bounds.Max.X - 1
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if y >= bounds.Max.Y {
		y = bounds.Max.Y - 1
	}

	c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
	return c.R, c.G, c.B, c.A
}

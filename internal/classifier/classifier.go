package classifier

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"golang.org/x/image/draw"
	"gorgonia.org/tensor"
)

type Labels struct {
	Labels []string `json:"labels"`
}

type ClassificationResult struct {
	PredictedClass string
	Scores         map[string]float32
}

type Classifier struct {
	backend *gorgonnx.Graph
	model   *onnx.Model
	labels  Labels
}

// NewClassifier cria classifier com onnx-go (pure Go, sem DLL)
func NewClassifier(modelsFS embed.FS) (*Classifier, error) {
	c := &Classifier{}

	// Carrega labels
	metaData, err := modelsFS.ReadFile("models/classifier/meta.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read meta.json: %w", err)
	}

	if err := json.Unmarshal(metaData, &c.labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %w", err)
	}

	// Carrega modelo
	modelData, err := modelsFS.ReadFile("models/classifier/model.onnx")
	if err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
	}

	// Inicializa backend
	c.backend = gorgonnx.NewGraph()
	c.model = onnx.NewModel(c.backend)

	if err := c.model.UnmarshalBinary(modelData); err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	return c, nil
}

func (c *Classifier) Classify(imgPath string) (*ClassificationResult, error) {
	// Load image
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer imgFile.Close()

	imgData, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return c.ClassifyImage(imgData)
}

func (c *Classifier) ClassifyImage(imgData image.Image) (*ClassificationResult, error) {
	// Resize to 384x384
	resized := image.NewRGBA(image.Rect(0, 0, 384, 384))
	draw.BiLinear.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Over, nil)

	// Normalize and convert to CHW format [1,3,384,384]
	inputData := make([]float32, 1*3*384*384)
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Normalização: (pixel/255 - 0.5) / 0.5
			inputData[0*384*384+y*384+x] = (float32(r>>8)/255.0 - 0.5) / 0.5
			inputData[1*384*384+y*384+x] = (float32(g>>8)/255.0 - 0.5) / 0.5
			inputData[2*384*384+y*384+x] = (float32(b>>8)/255.0 - 0.5) / 0.5
		}
	}

	inputTensor := tensor.New(
		tensor.WithShape(1, 3, 384, 384),
		tensor.WithBacking(inputData),
	)

	// Set input
	if err := c.model.SetInput(0, inputTensor); err != nil {
		return nil, fmt.Errorf("failed to set input: %w", err)
	}

	// Run inference
	if err := c.backend.Run(); err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Get outputs
	outputs, err := c.model.GetOutputTensors()
	if err != nil {
		return nil, fmt.Errorf("failed to get outputs: %w", err)
	}

	outputData := outputs[0].Data().([]float32)

	// Process results
	scores := map[string]float32{
		c.labels.Labels[0]: outputData[0],
		c.labels.Labels[1]: outputData[1],
	}

	maxKey := ""
	maxVal := float32(-math.MaxFloat32)
	for k, v := range scores {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	return &ClassificationResult{
		PredictedClass: maxKey,
		Scores:         scores,
	}, nil
}

func (c *Classifier) Close() {
	// Gorgonia não precisa de cleanup explícito
}

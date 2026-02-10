package classifier

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"

	ort "github.com/yalue/onnxruntime_go"
	"golang.org/x/image/draw"
)

// Labels represents the model labels structure
type Labels struct {
	Labels []string `json:"labels"`
}

// ClassificationResult holds the classification output
type ClassificationResult struct {
	PredictedClass string
	Scores         map[string]float32
}

// Classifier handles anime/real image classification
type Classifier struct {
	modelPath    string
	metaPath     string
	session      *ort.Session[float32]
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	labels       Labels
	initialized  bool
}

// NewClassifier creates a new classifier instance
func NewClassifier(modelDir, onnxLibPath string) (*Classifier, error) {
	c := &Classifier{
		modelPath: filepath.Join(modelDir, "classifier", "model.onnx"),
		metaPath:  filepath.Join(modelDir, "classifier", "meta.json"),
	}

	// Initialize ONNX Runtime
	ort.SetSharedLibraryPath(onnxLibPath)
	if err := ort.InitializeEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX environment: %w", err)
	}

	// Load labels
	metaBytes, err := os.ReadFile(c.metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta.json: %w", err)
	}
	if err := json.Unmarshal(metaBytes, &c.labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %w", err)
	}

	// Create tensors
	inputShape := ort.NewShape(1, 3, 384, 384)
	inputData := make([]float32, 1*3*384*384)
	c.inputTensor, err = ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %w", err)
	}

	outputShape := ort.NewShape(1, 2)
	c.outputTensor, err = ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		c.inputTensor.Destroy()
		return nil, fmt.Errorf("failed to create output tensor: %w", err)
	}

	// Create session
	c.session, err = ort.NewSession[float32](
		c.modelPath,
		[]string{"input"},
		[]string{"output"},
		[]*ort.Tensor[float32]{c.inputTensor},
		[]*ort.Tensor[float32]{c.outputTensor},
	)
	if err != nil {
		c.inputTensor.Destroy()
		c.outputTensor.Destroy()
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	c.initialized = true
	return c, nil
}

// Classify processes an image and returns the classification result
func (c *Classifier) Classify(imgPath string) (*ClassificationResult, error) {
	if !c.initialized {
		return nil, fmt.Errorf("classifier not initialized")
	}

	// Load and decode image
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer imgFile.Close()

	imgData, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize to 384x384
	resized := image.NewRGBA(image.Rect(0, 0, 384, 384))
	draw.BiLinear.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Over, nil)

	// Normalize and convert to CHW format [1,3,384,384]
	data := c.inputTensor.GetData()
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			data[0*384*384+y*384+x] = (float32(r>>8)/255.0 - 0.5) / 0.5
			data[1*384*384+y*384+x] = (float32(g>>8)/255.0 - 0.5) / 0.5
			data[2*384*384+y*384+x] = (float32(b>>8)/255.0 - 0.5) / 0.5
		}
	}

	// Run inference
	if err := c.session.Run(); err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Process results
	outputData := c.outputTensor.GetData()
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

// Close releases all resources
func (c *Classifier) Close() {
	if c.session != nil {
		c.session.Destroy()
	}
	if c.inputTensor != nil {
		c.inputTensor.Destroy()
	}
	if c.outputTensor != nil {
		c.outputTensor.Destroy()
	}
	ort.DestroyEnvironment()
}

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

	"pure-go-project/internal/runtime" // ⭐ Importa o singleton
)

type Labels struct {
	Labels []string `json:"labels"`
}

type ClassificationResult struct {
	PredictedClass string
	Scores         map[string]float32
}

type Classifier struct {
	modelPath    string
	metaPath     string
	session      *ort.Session[float32]
	inputTensor  *ort.Tensor[float32]
	outputTensor *ort.Tensor[float32]
	labels       Labels
}

// NewClassifier cria uma nova instância do classificador
func NewClassifier(modelDir, onnxLibPath string) (*Classifier, error) {
	// ⭐ Usa o singleton para inicializar o runtime
	if err := runtime.GetInstance().Initialize(onnxLibPath); err != nil {
		return nil, fmt.Errorf("falha ao inicializar runtime: %w", err)
	}

	c := &Classifier{
		modelPath: filepath.Join(modelDir, "classifier", "model.onnx"),
		metaPath:  filepath.Join(modelDir, "classifier", "meta.json"),
	}

	// Carrega labels
	metaBytes, err := os.ReadFile(c.metaPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler meta.json: %w", err)
	}

	if err := json.Unmarshal(metaBytes, &c.labels); err != nil {
		return nil, fmt.Errorf("falha ao parsear labels: %w", err)
	}

	// Cria tensors
	inputShape := ort.NewShape(1, 3, 384, 384)
	inputData := make([]float32, 1*3*384*384)
	c.inputTensor, err = ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar input tensor: %w", err)
	}

	outputShape := ort.NewShape(1, 2)
	c.outputTensor, err = ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		c.inputTensor.Destroy()
		return nil, fmt.Errorf("falha ao criar output tensor: %w", err)
	}

	// Cria sessão
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
		return nil, fmt.Errorf("falha ao criar sessão: %w", err)
	}

	return c, nil
}

func (c *Classifier) Classify(imgPath string) (*ClassificationResult, error) {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir imagem: %w", err)
	}
	defer imgFile.Close()

	imgData, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar imagem: %w", err)
	}

	// Resize para 384x384
	resized := image.NewRGBA(image.Rect(0, 0, 384, 384))
	draw.BiLinear.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Over, nil)

	// Normalização
	data := c.inputTensor.GetData()
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			idx := y*384 + x
			data[0*384*384+idx] = (float32(r>>8)/255.0 - 0.5) / 0.5
			data[1*384*384+idx] = (float32(g>>8)/255.0 - 0.5) / 0.5
			data[2*384*384+idx] = (float32(b>>8)/255.0 - 0.5) / 0.5
		}
	}

	// Inferência
	if err := c.session.Run(); err != nil {
		return nil, fmt.Errorf("inferência falhou: %w", err)
	}

	// Resultados
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
}

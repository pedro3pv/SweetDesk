package main

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"

	ort "github.com/yalue/onnxruntime_go"
)

type Labels struct {
	Labels []string `json:"labels"`
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// ensureModelFiles checks if model files exist and downloads them if not
func ensureModelFiles(modelDir string) error {
	modelPath := filepath.Join(modelDir, "model.onnx")
	metaPath := filepath.Join(modelDir, "meta.json")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check and download model.onnx
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Println("Downloading model.onnx...")
		url := "https://huggingface.co/deepghs/anime_real_cls/resolve/main/caformer_s36_v1.3_fixed/model.onnx"
		if err := downloadFile(url, modelPath); err != nil {
			return fmt.Errorf("failed to download model.onnx: %w", err)
		}
		log.Println("model.onnx downloaded successfully")
	} else {
		log.Println("model.onnx already exists, skipping download")
	}

	// Check and download meta.json
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		log.Println("Downloading meta.json...")
		url := "https://huggingface.co/deepghs/anime_real_cls/resolve/main/caformer_s36_v1.3_fixed/meta.json"
		if err := downloadFile(url, metaPath); err != nil {
			return fmt.Errorf("failed to download meta.json: %w", err)
		}
		log.Println("meta.json downloaded successfully")
	} else {
		log.Println("meta.json already exists, skipping download")
	}

	return nil
}

func main() {
	modelDir := "./caformer_s36_v1.3_fixed"
	modelPath := filepath.Join(modelDir, "model.onnx")
	metaPath := filepath.Join(modelDir, "meta.json")

	// Ensure model files are downloaded
	if err := ensureModelFiles(modelDir); err != nil {
		log.Fatalf("Failed to ensure model files: %v", err)
	}

	// Inicializa biblioteca ONNX Runtime
	ort.SetSharedLibraryPath("/opt/homebrew/lib/libonnxruntime.dylib") // macOS
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	// Carrega labels
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		log.Fatal(err)
	}
	var labels Labels
	json.Unmarshal(metaBytes, &labels)

	// Cria tensors de input/output
	inputShape := ort.NewShape(1, 3, 384, 384)
	inputData := make([]float32, 1*3*384*384)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		log.Fatal(err)
	}
	defer inputTensor.Destroy()

	outputShape := ort.NewShape(1, 2)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		log.Fatal(err)
	}
	defer outputTensor.Destroy()

	// Cria sessão (com nomes de inputs/outputs)
	session, err := ort.NewSession[float32](
		modelPath,
		[]string{"input"},  // nome do input tensor
		[]string{"output"}, // nome do output tensor
		[]*ort.Tensor[float32]{inputTensor},
		[]*ort.Tensor[float32]{outputTensor},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Destroy()

	// Processa imagem
	imgPath := "test_img_1.jpg"
	imgFile, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	imgData, _, err := image.Decode(imgFile)
	imgFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Resize para 384x384
	resized := image.NewRGBA(image.Rect(0, 0, 384, 384))
	draw.BiLinear.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Over, nil)

	// Normalize e converte para CHW [1,3,384,384]
	data := inputTensor.GetData()
	for y := 0; y < 384; y++ {
		for x := 0; x < 384; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			data[0*384*384+y*384+x] = (float32(r>>8)/255.0 - 0.5) / 0.5
			data[1*384*384+y*384+x] = (float32(g>>8)/255.0 - 0.5) / 0.5
			data[2*384*384+y*384+x] = (float32(b>>8)/255.0 - 0.5) / 0.5
		}
	}

	// Executa inferência
	err = session.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Pega resultado
	outputData := outputTensor.GetData()
	values := map[string]float32{
		labels.Labels[0]: outputData[0],
		labels.Labels[1]: outputData[1],
	}

	maxKey := ""
	maxVal := float32(-math.MaxFloat32)
	for k, v := range values {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	fmt.Printf("Resultado: %s (real: %.4f, anime: %.4f)\n", maxKey, values["real"], values["anime"])
}

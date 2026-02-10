package main

import (
	"fmt"
	"log"
	"path/filepath"

	"pure-go-project/internal/classifier"
	"pure-go-project/internal/downloader"
	"pure-go-project/internal/upscaler"
)

func main() {
	modelDir := "./models"
	onnxLibPath := "/opt/homebrew/lib/libonnxruntime.dylib" // macOS
	inputImage := "test_img_1.jpg"

	// 1. Baixa todos os modelos necessários
	dl := downloader.NewDownloader(modelDir)
	if err := dl.EnsureAllModels(); err != nil {
		log.Fatalf("Erro ao baixar modelos: %v", err)
	}

	// 2. Classifica a imagem
	clf, err := classifier.NewClassifier(modelDir, onnxLibPath)
	if err != nil {
		log.Fatalf("Erro ao criar classificador: %v", err)
	}
	defer clf.Close()

	result, err := clf.Classify(inputImage)
	if err != nil {
		log.Fatalf("Erro na classificação: %v", err)
	}

	fmt.Printf("📊 Classificação: %s\n", result.PredictedClass)
	fmt.Printf("   Scores: real=%.4f, anime=%.4f\n\n",
		result.Scores["real"],
		result.Scores["anime"])

	// 3. Processa com upscaler apropriado
	var upscalerType upscaler.UpscalerType
	if result.PredictedClass == "anime" {
		upscalerType = upscaler.RealCUGAN
		fmt.Println("🎨 Processando com RealCUGAN-pro (anime)...")
	} else {
		upscalerType = upscaler.LSDIR
		fmt.Println("📸 Processando com 4xLSDIR (foto real)...")
	}

	ups, err := upscaler.NewUpscaler(upscalerType, modelDir, onnxLibPath)
	if err != nil {
		log.Fatalf("Erro ao criar upscaler: %v", err)
	}
	defer ups.Close()

	// Define nome do arquivo de saída
	outputImage := generateOutputPath(inputImage, string(upscalerType))

	if err := ups.Upscale(inputImage, outputImage); err != nil {
		log.Fatalf("Erro no upscaling: %v", err)
	}

	fmt.Printf("✅ Imagem processada salva em: %s\n", outputImage)
}

func generateOutputPath(inputPath, upscalerType string) string {
	ext := filepath.Ext(inputPath)
	base := inputPath[:len(inputPath)-len(ext)]
	return fmt.Sprintf("%s_%s_4x%s", base, upscalerType, ".png")
}

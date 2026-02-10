package main

import (
	"fmt"
	"log"

	"pure-go-project/internal/classifier"
	"pure-go-project/internal/downloader"
	"pure-go-project/internal/upscaler"
)

func main() {
	modelDir := "./models"
	onnxLibPath := "/opt/homebrew/lib/libonnxruntime.dylib"
	inputImage := "test_img_1.jpg"

	// Download e classificação
	dl := downloader.NewDownloader(modelDir)
	if err := dl.EnsureAllModels(); err != nil {
		log.Fatalf("Erro: %v", err)
	}

	clf, err := classifier.NewClassifier(modelDir, onnxLibPath)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	defer clf.Close()

	result, err := clf.Classify(inputImage)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	fmt.Printf("📊 Classificação: %s (%.2f%%)\n\n",
		result.PredictedClass,
		result.Scores[result.PredictedClass]*100)

	// Determina upscaler
	var upscalerType upscaler.UpscalerType
	if result.PredictedClass == "anime" {
		upscalerType = upscaler.RealCUGAN
	} else {
		upscalerType = upscaler.LSDIR
	}

	ups, err := upscaler.NewUpscaler(upscalerType, modelDir, onnxLibPath)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	defer ups.Close()

	// EXEMPLO 1: Upscale padrão 4x
	fmt.Println("\n=== 1. Upscale padrão (4x) ===")
	ups.Upscale(inputImage, "output_4x.png")

	// EXEMPLO 2: Fator de escala customizado (2x)
	fmt.Println("\n=== 2. Upscale 2x ===")
	ups.UpscaleWithOptions(inputImage, "output_2x.png", &upscaler.UpscaleOptions{
		ScaleFactor: 2.0,
	})

	// EXEMPLO 3: Resolução específica (Full HD)
	fmt.Println("\n=== 3. Resolução específica 1920x1080 ===")
	ups.UpscaleWithOptions(inputImage, "output_fullhd.png", &upscaler.UpscaleOptions{
		TargetWidth:     1920,
		TargetHeight:    1080,
		KeepAspectRatio: true, // Mantém proporção
	})

	// EXEMPLO 4: Apenas largura (altura proporcional)
	fmt.Println("\n=== 4. Largura 3840px (4K) ===")
	ups.UpscaleWithOptions(inputImage, "output_4k.png", &upscaler.UpscaleOptions{
		TargetWidth:     3840,
		KeepAspectRatio: true,
	})

	// EXEMPLO 5: Com limite máximo
	fmt.Println("\n=== 5. Upscale 4x limitado a 2560px ===")
	ups.UpscaleWithOptions(inputImage, "output_limited.png", &upscaler.UpscaleOptions{
		ScaleFactor:   4.0,
		MaxResolution: 2560, // Não ultrapassa 2560 em nenhuma dimensão
	})

	// EXEMPLO 6: Força dimensões exatas (pode distorcer)
	fmt.Println("\n=== 6. Força 1024x1024 (quadrado) ===")
	ups.UpscaleWithOptions(inputImage, "output_square.png", &upscaler.UpscaleOptions{
		TargetWidth:     1024,
		TargetHeight:    1024,
		KeepAspectRatio: false, // Força dimensões exatas
	})

	fmt.Println("\n✅ Processamento concluído!")
}

package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"pure-go-project/internal/classifier"
	"pure-go-project/internal/downloader"
	"pure-go-project/internal/seamcarver"
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

	// EXEMPLO 6: Workflow completo - seam carving + upscale para resolução específica
	fmt.Println("\n=== 6. Workflow completo: 16:9 → Upscale → 4K ===")

	// Carrega novamente
	imgFile2, err := os.Open(inputImage)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	img2, _, err := image.Decode(imgFile2)
	imgFile2.Close()
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	// Ajusta aspect ratio
	carver2 := seamcarver.NewSeamCarver(img2)
	adjusted, err := carver2.AdjustAspectRatio(seamcarver.AspectRatioOptions{
		TargetRatio:     16.0 / 9.0,
		MaxDeltaBySeams: 300,
	})
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	carver2.SaveImage(adjusted, "temp_adjusted.png")

	// Upscale para 4K (3840x2160)
	if err := ups.UpscaleWithOptions("temp_adjusted.png", "output_4k_16x9.png", &upscaler.UpscaleOptions{
		TargetWidth:     3840,
		TargetHeight:    2160,
		KeepAspectRatio: true,
	}); err != nil {
		log.Printf("Erro: %v", err)
	}

	// Limpa arquivo temporário
	os.Remove("temp_adjusted.png")

	fmt.Println("\n✅ Processamento concluído!")

	// EXEMPLO 7: Workflow completo - seam carving + upscale para 1024x1024
	fmt.Println("\n=== 7. Workflow completo: Ajusta para 1024x1024 ===")

	// Carrega novamente
	imgFile3, err := os.Open(inputImage)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	img3, _, err := image.Decode(imgFile3)
	imgFile3.Close()
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	// Ajusta aspect ratio para quadrado (1:1)
	carver3 := seamcarver.NewSeamCarver(img3)
	adjustedSquare, err := carver3.AdjustAspectRatio(seamcarver.AspectRatioOptions{
		TargetRatio:     1.0, // quadrado
		MaxDeltaBySeams: 300,
	})
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	carver3.SaveImage(adjustedSquare, "temp_adjusted_square.png")

	// Upscale para 1024x1024
	if err := ups.UpscaleWithOptions("temp_adjusted_square.png", "output_1024x1024.png", &upscaler.UpscaleOptions{
		TargetWidth:     1024,
		TargetHeight:    1024,
		KeepAspectRatio: false, // força quadrado
	}); err != nil {
		log.Printf("Erro: %v", err)
	}

	// Limpa arquivo temporário
	os.Remove("temp_adjusted_square.png")

	fmt.Println("\n✅ Processamento concluído!")
}

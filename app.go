package main

import (
	"SweetDesk/internal/services"
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx            context.Context
	imageProcessor *services.ImageProcessor
	upscaler       *services.Upscaler
	pixabayKey     string
	modelsFS       embed.FS
	tempLibPath    string
}

// NewApp creates a new App application struct
func NewApp(modelsFS embed.FS) *App {
	return &App{
		modelsFS: modelsFS,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize services
	a.imageProcessor = services.NewImageProcessor(ctx)

	// Extract and initialize ONNX library
	onnxLibPath, err := a.extractONNXLibrary()
	if err != nil {
		fmt.Printf("Failed to extract ONNX library: %v\n", err)
		return
	}
	a.tempLibPath = onnxLibPath

	// Initialize upscaler with RealCUGAN by default
	upscaler, err := services.NewUpscaler(ctx, services.RealCUGAN, a.modelsFS, onnxLibPath)
	if err != nil {
		fmt.Printf("Failed to initialize upscaler: %v\n", err)
	} else {
		a.upscaler = upscaler
	}

	// Get Pixabay API key from environment
	a.pixabayKey = os.Getenv("PIXABAY_API_KEY")
}

// extractONNXLibrary extracts the embedded ONNX library to temp directory
func (a *App) extractONNXLibrary() (string, error) {
	// Get platform-specific library data (defined in app_*.go files)
	libData, fileName := getONNXLibrary()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sweetdesk-onnx-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Write library to temp file
	libPath := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(libPath, libData, 0755); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to write library: %w", err)
	}

	fmt.Printf("✅ ONNX Runtime extracted to: %s\n", libPath)
	return libPath, nil
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	// Cleanup upscaler
	if a.upscaler != nil {
		a.upscaler.Close()
	}

	// Clean up temp library file
	if a.tempLibPath != "" {
		tempDir := filepath.Dir(a.tempLibPath)
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Printf("Failed to cleanup temp dir: %v\n", err)
		}
	}

	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SearchImages searches for images using Pixabay API
func (a *App) SearchImages(query string, page int, perPage int) ([]services.ImageResult, error) {
	if a.pixabayKey == "" {
		return nil, fmt.Errorf("Pixabay API key not configured")
	}

	provider := services.NewPixabayProvider(a.ctx, a.pixabayKey)
	options := services.SearchOptions{
		Page:        page,
		PerPage:     perPage,
		MinWidth:    1920,
		MinHeight:   1080,
		Orientation: "horizontal",
	}

	return provider.Search(query, options)
}

// DownloadImage downloads an image from a URL
func (a *App) DownloadImage(imageURL string) (string, error) {
	if a.pixabayKey == "" {
		return "", fmt.Errorf("Pixabay API key not configured")
	}

	provider := services.NewPixabayProvider(a.ctx, a.pixabayKey)
	data, err := provider.Download(imageURL)
	if err != nil {
		return "", err
	}

	return a.imageProcessor.ConvertToBase64(data), nil
}

// UpscaleImage upscales an image using AI
func (a *App) UpscaleImage(base64Data string, imageType string, scale int) (string, error) {
	if a.upscaler == nil {
		return "", fmt.Errorf("upscaler not initialized")
	}

	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}

	// Configure options based on scale
	options := &services.UpscaleOptions{
		ScaleFactor: float64(scale),
		Format:      "png",
	}

	upscaled, err := a.upscaler.UpscaleBytes(data, options)
	if err != nil {
		return "", err
	}

	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// ProcessImage is the main processing pipeline
func (a *App) ProcessImage(base64Data string, targetResolution string) (string, error) {
	if a.upscaler == nil {
		return "", fmt.Errorf("upscaler not initialized")
	}

	// 1. Decode image
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 2. Determine scale and options based on target resolution
	var options *services.UpscaleOptions
	switch targetResolution {
	case "4K", "3840x2160":
		options = &services.UpscaleOptions{
			TargetWidth:     3840,
			TargetHeight:    2160,
			KeepAspectRatio: true,
			Format:          "png",
		}
	case "1080p", "1920x1080":
		options = &services.UpscaleOptions{
			TargetWidth:     1920,
			TargetHeight:    1080,
			KeepAspectRatio: true,
			Format:          "png",
		}
	default:
		options = &services.UpscaleOptions{
			ScaleFactor: 4.0,
			Format:      "png",
		}
	}

	// 3. Upscale image
	upscaled, err := a.upscaler.UpscaleBytes(data, options)
	if err != nil {
		return "", fmt.Errorf("failed to upscale: %w", err)
	}

	// 4. Return result as base64
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

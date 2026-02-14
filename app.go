package main

import (
	"SweetDesk/internal/services"
	"context"
	"fmt"
	"os"

	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

// App struct
type App struct {
	ctx              context.Context
	coreProcessor    *processor.Processor
	imageProcessor   *services.ImageProcessor
	pixabayKey       string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	
	// Initialize SweetDesk-core processor
	var err error
	a.coreProcessor, err = processor.New(processor.Config{
		ModelDir: "./models",
		ONNXLibPath: os.Getenv("ONNX_LIB_PATH"),
	})
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to initialize SweetDesk-core: %v\n", err)
		a.coreProcessor = nil
	}
	
	// Initialize image processor for base64 conversions
	a.imageProcessor = services.NewImageProcessor(ctx)
	
	// Get Pixabay API key from environment
	a.pixabayKey = os.Getenv("PIXABAY_API_KEY")
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.coreProcessor != nil {
		a.coreProcessor.Close()
	}
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
	
	// Convert to base64 for frontend
	return a.imageProcessor.ConvertToBase64(data), nil
}

// ClassifyImage classifies an image as anime or photo using SweetDesk-core
func (a *App) ClassifyImage(base64Data string) (string, error) {
	if a.coreProcessor == nil {
		// Fallback to simple classification if core not available
		return "photo", nil
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	result, err := a.coreProcessor.Classify(data)
	if err != nil {
		// Fallback
		fmt.Printf("⚠️  Classification failed: %v\n", err)
		return "photo", nil
	}
	
	return result.Type, nil
}

// UpscaleImage upscales an image using SweetDesk-core AI engine
func (a *App) UpscaleImage(base64Data string, imageType string, scale int) (string, error) {
	if a.coreProcessor == nil {
		return "", fmt.Errorf("SweetDesk-core not initialized")
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	options := processor.ProcessOptions{
		ImageType:       imageType,
		TargetScale:     scale,
		OutputFormat:    "png",
	}
	
	upscaled, err := a.coreProcessor.Upscale(data, options)
	if err != nil {
		return "", err
	}
	
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// UpscaleToResolution upscales an image to a specific resolution
func (a *App) UpscaleToResolution(base64Data string, targetWidth int, targetHeight int) (string, error) {
	if a.coreProcessor == nil {
		return "", fmt.Errorf("SweetDesk-core not initialized")
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	options := processor.ProcessOptions{
		TargetWidth:     targetWidth,
		TargetHeight:    targetHeight,
		KeepAspectRatio: true,
		OutputFormat:    "png",
	}
	
	upscaled, err := a.coreProcessor.UpscaleToResolution(data, options)
	if err != nil {
		return "", err
	}
	
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// AdjustAspectRatio adjusts image aspect ratio using seam carving
func (a *App) AdjustAspectRatio(base64Data string, targetWidth int, targetHeight int) (string, error) {
	if a.coreProcessor == nil {
		return "", fmt.Errorf("SweetDesk-core not initialized")
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	options := processor.ProcessOptions{
		TargetWidth:     targetWidth,
		TargetHeight:    targetHeight,
		UseSeamCarving:  true,
		OutputFormat:    "png",
	}
	
	adjusted, err := a.coreProcessor.AdjustAspectRatio(data, options)
	if err != nil {
		return "", err
	}
	
	return a.imageProcessor.ConvertToBase64(adjusted), nil
}

// ProcessImage is the main processing pipeline - complete workflow
func (a *App) ProcessImage(base64Data string, targetResolution string, useSeamCarving bool, targetAspectRatio string) (string, error) {
	if a.coreProcessor == nil {
		return "", fmt.Errorf("SweetDesk-core not initialized")
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}
	
	// Parse target resolution
	var targetWidth, targetHeight int
	switch targetResolution {
	case "4K":
		targetWidth, targetHeight = 3840, 2160
	case "1440p":
		targetWidth, targetHeight = 2560, 1440
	case "1080p":
		targetWidth, targetHeight = 1920, 1080
	case "720p":
		targetWidth, targetHeight = 1280, 720
	default:
		targetWidth, targetHeight = 3840, 2160 // Default to 4K
	}
	
	// If seam carving is requested, adjust aspect ratio first
	if useSeamCarving && targetAspectRatio != "" {
		options := processor.ProcessOptions{
			TargetAspectRatio: targetAspectRatio,
			UseSeamCarving:    true,
			OutputFormat:      "png",
		}
		
		data, err = a.coreProcessor.AdjustAspectRatio(data, options)
		if err != nil {
			return "", fmt.Errorf("aspect ratio adjustment failed: %w", err)
		}
	}
	
	// Classify image
	classifyOptions := processor.ProcessOptions{OutputFormat: "png"}
	classifyResult, err := a.coreProcessor.ClassifyWithOptions(data, classifyOptions)
	imageType := "photo" // default fallback
	if err == nil {
		imageType = classifyResult.Type
	}
	
	// Upscale to target resolution
	upscaleOptions := processor.ProcessOptions{
		ImageType:       imageType,
		TargetWidth:     targetWidth,
		TargetHeight:    targetHeight,
		KeepAspectRatio: true,
		OutputFormat:    "png",
	}
	
	upscaled, err := a.coreProcessor.UpscaleToResolution(data, upscaleOptions)
	if err != nil {
		return "", fmt.Errorf("upscaling failed: %w", err)
	}
	
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

package main

import (
	"SweetDesk/internal/services"
	"context"
	"embed"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx            context.Context
	imageProcessor *services.ImageProcessor
	upscaler       *services.Upscaler
	pixabayKey     string
	modelsFS       embed.FS
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

	// Initialize upscaler (pure Go, no DLL needed!)
	upscaler, err := services.NewUpscaler(ctx, services.RealCUGAN, a.modelsFS)
	if err != nil {
		fmt.Printf("Failed to initialize upscaler: %v\n", err)
	} else {
		a.upscaler = upscaler
		fmt.Println("✅ Upscaler initialized (pure Go, no external dependencies)")
	}

	// Get Pixabay API key from environment
	a.pixabayKey = os.Getenv("PIXABAY_API_KEY")
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
func (a *App) ProcessImage(base64Data string, targetWidth int, targetHeight int, savePath string, fileName string) (string, error) {
	if a.upscaler == nil {
		return "", fmt.Errorf("upscaler not initialized")
	}

	// 1. Early validation of dimensions (before decoding Base64)
	if targetWidth <= 0 || targetHeight <= 0 {
		return "", fmt.Errorf("invalid target resolution: %dx%d (dimensions must be positive)", targetWidth, targetHeight)
	}

	// 2. Configurable maxResolution via UpscaleOptions
	const defaultMaxResolution = 16384
	maxResolution := defaultMaxResolution

	// Prepare options with possible MaxResolution override
	options := &services.UpscaleOptions{
		TargetWidth:     targetWidth,
		TargetHeight:    targetHeight,
		KeepAspectRatio: false,
		Format:          "png",
	}
	if options.MaxResolution > 0 && options.MaxResolution < defaultMaxResolution {
		maxResolution = options.MaxResolution
	}

	// 3. Per-dimension and total pixel checks (before decode)
	if targetWidth > maxResolution || targetHeight > maxResolution {
		return "", fmt.Errorf("target resolution %dx%d exceeds maximum allowed dimensions of %dx%d",
			targetWidth, targetHeight, maxResolution, maxResolution)
	}

	maxPixels := int64(maxResolution) * int64(maxResolution)
	totalPixels := int64(targetWidth) * int64(targetHeight)
	if totalPixels > maxPixels {
		estBytes := totalPixels * 4
		estMB := float64(estBytes) / (1024 * 1024)
		estGB := estMB / 1024
		return "", fmt.Errorf("target resolution %dx%d (%d megapixels) exceeds maximum pixel count (%.1f megapixels). Estimated RGBA buffer: %.1f MB (%.2f GB)",
			targetWidth, targetHeight, totalPixels/1_000_000, float64(maxPixels)/1_000_000, estMB, estGB)
	}

	// 4. Decode image (after validation)
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 5. Upscale image
	upscaled, err := a.upscaler.UpscaleBytes(data, options)
	if err != nil {
		return "", fmt.Errorf("failed to upscale: %w", err)
	}

	// 6. Save to file if savePath is provided
	if savePath != "" && fileName != "" {
		_, err := a.imageProcessor.SaveToFile(upscaled, savePath, fileName)
		if err != nil {
			return "", fmt.Errorf("failed to save image: %w", err)
		}
	}

	// 7. Return result as base64
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

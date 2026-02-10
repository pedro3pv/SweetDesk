package main

import (
	"SweetDesk/internal/seamcarver"
	"SweetDesk/internal/services"
	"context"
	"fmt"
	"math"
	"os"
)

// App struct
type App struct {
	ctx            context.Context
	imageProcessor *services.ImageProcessor
	upscaler       *services.Upscaler
	pythonBridge   *services.PythonBridge
	pixabayKey     string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	
	// Initialize services
	a.imageProcessor = services.NewImageProcessor(ctx)
	a.upscaler = services.NewUpscaler(ctx)
	
	// Initialize Python bridge (optional - may not be available)
	bridge, err := services.NewPythonBridge(ctx)
	if err == nil {
		a.pythonBridge = bridge
	}
	
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
		Page:    page,
		PerPage: perPage,
		MinWidth: 1920,
		MinHeight: 1080,
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

// ClassifyImage classifies an image as anime or photo
func (a *App) ClassifyImage(base64Data string) (string, error) {
	if a.pythonBridge == nil {
		// Fallback to simple classification if Python not available
		return "photo", nil
	}
	
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	result, err := a.pythonBridge.ClassifyImage(data)
	if err != nil {
		// Fallback
		return "photo", nil
	}
	
	return result.Type, nil
}

// UpscaleImage upscales an image using AI
func (a *App) UpscaleImage(base64Data string, imageType string, scale int) (string, error) {
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}
	
	model := a.upscaler.GetRecommendedModel(imageType)
	
	options := services.UpscaleOptions{
		Model:  model,
		Scale:  scale,
		Format: "png",
	}
	
	upscaled, err := a.upscaler.UpscaleImage(data, options)
	if err != nil {
		return "", err
	}
	
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// ProcessImage is the main processing pipeline
func (a *App) ProcessImage(base64Data string, targetResolution string, useSeamCarving bool) (string, error) {
	// 1. Decode image
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}
	
	// 2. Classify image
	imageType := "photo"
	if a.pythonBridge != nil {
		result, err := a.pythonBridge.ClassifyImage(data)
		if err == nil {
			imageType = result.Type
		}
	}
	
	// 3. Upscale image
	model := a.upscaler.GetRecommendedModel(imageType)
	scale := 4 // Default to 4x for 4K
	
	options := services.UpscaleOptions{
		Model:  model,
		Scale:  scale,
		Format: "png",
	}
	
	upscaled, err := a.upscaler.UpscaleImage(data, options)
	if err != nil {
		return "", fmt.Errorf("failed to upscale: %w", err)
	}
	
	// 4. Adjust aspect ratio if needed
	if useSeamCarving && a.pythonBridge != nil {
		// Get current dimensions
		img, _, err := a.imageProcessor.LoadImageFromBytes(upscaled)
		if err == nil {
			bounds := img.Bounds()
			currentWidth := bounds.Dx()
			currentHeight := bounds.Dy()
			
			// Calculate target 16:9 dimensions
			targetWidth := currentWidth
			targetHeight := (currentWidth * 9) / 16
			
			if targetHeight != currentHeight {
				seamOptions := services.SeamCarvingOptions{
					TargetWidth:  targetWidth,
					TargetHeight: targetHeight,
					EnergyMode:   "forward",
				}
				
				upscaled, err = a.pythonBridge.ApplySeamCarving(upscaled, seamOptions)
				if err != nil {
					// Fallback to original if seam carving fails
					fmt.Printf("Seam carving failed: %v\n", err)
				}
			}
		}
	}
	
	// 5. Return result as base64
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// ProcessImageWithCustomResolution processes an image with fully customizable resolution
func (a *App) ProcessImageWithCustomResolution(base64Data string, targetWidth int, targetHeight int, aspectRatioW int, aspectRatioH int, useSeamCarving bool) (string, error) {
	// 1. Decode image
	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}
	
	// 2. Classify image
	imageType := "photo"
	if a.pythonBridge != nil {
		result, err := a.pythonBridge.ClassifyImage(data)
		if err == nil {
			imageType = result.Type
		}
	}
	
	// 3. Determine target dimensions
	finalWidth := targetWidth
	finalHeight := targetHeight
	
	// If custom aspect ratio is specified, calculate dimensions
	if aspectRatioW > 0 && aspectRatioH > 0 && (targetWidth == 0 || targetHeight == 0) {
		// Load image to get current dimensions
		img, _, err := a.imageProcessor.LoadImageFromBytes(data)
		if err == nil {
			bounds := img.Bounds()
			currentWidth := bounds.Dx()
			
			aspectRatio := float64(aspectRatioW) / float64(aspectRatioH)
			
			// Calculate target dimensions maintaining aspect ratio
			if targetWidth > 0 {
				finalWidth = targetWidth
				finalHeight = int(float64(targetWidth) / aspectRatio)
			} else if targetHeight > 0 {
				finalHeight = targetHeight
				finalWidth = int(float64(targetHeight) * aspectRatio)
			} else {
				// Use current width and calculate height
				finalWidth = currentWidth
				finalHeight = int(float64(currentWidth) / aspectRatio)
			}
		}
	}
	
	// 4. Upscale image (if dimensions are larger than current)
	img, _, err := a.imageProcessor.LoadImageFromBytes(data)
	if err != nil {
		return "", fmt.Errorf("failed to load image: %w", err)
	}
	
	bounds := img.Bounds()
	currentWidth := bounds.Dx()
	currentHeight := bounds.Dy()
	
	upscaled := data
	
	// Only upscale if target is larger
	if finalWidth > currentWidth || finalHeight > currentHeight {
		model := a.upscaler.GetRecommendedModel(imageType)
		
		// Calculate required scale factor
		scaleW := float64(finalWidth) / float64(currentWidth)
		scaleH := float64(finalHeight) / float64(currentHeight)
		requiredScale := int(math.Max(scaleW, scaleH))
		
		if requiredScale < 2 {
			requiredScale = 2
		} else if requiredScale > 8 {
			requiredScale = 8
		}
		
		options := services.UpscaleOptions{
			Model:  model,
			Scale:  requiredScale,
			Format: "png",
		}
		
		upscaled, err = a.upscaler.UpscaleImage(data, options)
		if err != nil {
			return "", fmt.Errorf("failed to upscale: %w", err)
		}
	}
	
	// 5. Adjust to exact dimensions if needed
	if useSeamCarving {
		img, _, err := a.imageProcessor.LoadImageFromBytes(upscaled)
		if err == nil {
			bounds := img.Bounds()
			currentWidth := bounds.Dx()
			currentHeight := bounds.Dy()
			
			if currentWidth != finalWidth || currentHeight != finalHeight {
				// Use Go seamcarver for exact dimensions
				carver := seamcarver.NewSeamCarver(img)
				resizeOpts := seamcarver.ResizeOptions{
					TargetWidth:     finalWidth,
					TargetHeight:    finalHeight,
					MaxDeltaBySeams: 500,
				}
				
				resized, err := carver.ResizeToExactDimensions(resizeOpts)
				if err == nil {
					// Convert back to bytes
					upscaled, err = a.imageProcessor.EncodeImage(resized, "png", 100)
					if err != nil {
						fmt.Printf("Failed to encode resized image: %v\n", err)
					}
				} else {
					fmt.Printf("Resize to exact dimensions failed: %v\n", err)
				}
			}
		}
	} else {
		// Simple crop/scale to exact dimensions
		img, _, err := a.imageProcessor.LoadImageFromBytes(upscaled)
		if err == nil {
			bounds := img.Bounds()
			currentWidth := bounds.Dx()
			currentHeight := bounds.Dy()
			
			if currentWidth != finalWidth || currentHeight != finalHeight {
				// Simple resize using Go image libraries
				carver := seamcarver.NewSeamCarver(img)
				resizeOpts := seamcarver.ResizeOptions{
					TargetWidth:     finalWidth,
					TargetHeight:    finalHeight,
					MaxDeltaBySeams: 0, // No seam carving, just scale
				}
				
				resized, err := carver.ResizeToExactDimensions(resizeOpts)
				if err == nil {
					upscaled, err = a.imageProcessor.EncodeImage(resized, "png", 100)
					if err != nil {
						fmt.Printf("Failed to encode resized image: %v\n", err)
					}
				}
			}
		}
	}
	
	// 6. Return result as base64
	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

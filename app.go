package main

import (
	"SweetDesk/internal/services"
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// BatchItem represents a single image to process in a batch
type BatchItem struct {
	ID          string `json:"id"`
	Base64Data  string `json:"base64Data"`  // base64 image data or empty if downloadURL is used
	DownloadURL string `json:"downloadURL"` // URL to download image from
	Name        string `json:"name"`
	Dimension   string `json:"dimension"` // "WIDTHxHEIGHT"
}

// BatchItemStatus represents the processing status of a single item
type BatchItemStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"` // "pending", "processing", "done", "error"
	Error  string `json:"error,omitempty"`
}

// ProcessingStatus represents the overall batch processing state
type ProcessingStatus struct {
	IsProcessing bool              `json:"isProcessing"`
	Total        int               `json:"total"`
	Current      int               `json:"current"`
	Progress     int               `json:"progress"`
	Items        []BatchItemStatus `json:"items"`
	Done         bool              `json:"done"`
}

// App struct
type App struct {
	ctx            context.Context
	imageProcessor *services.ImageProcessor
	upscaler       *services.Upscaler
	pixabayKey     string
	modelsFS       embed.FS

	// Batch processing state
	procMu     sync.Mutex
	procStatus ProcessingStatus
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
func (a *App) domReady(ctx context.Context) {
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

// SelectDirectory opens the native OS directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	defaultDir := ""
	home, err := os.UserHomeDir()
	if err == nil {
		switch runtime.GOOS {
		case "windows":
			defaultDir = filepath.Join(home, "Pictures")
		case "darwin":
			defaultDir = filepath.Join(home, "Pictures")
		default:
			defaultDir = filepath.Join(home, "Pictures")
		}
	}

	result, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:            "Selecionar pasta para salvar",
		DefaultDirectory: defaultDir,
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return result, nil
}

// GetDefaultSavePath returns the default save path for the current OS
func (a *App) GetDefaultSavePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "Pictures", "SweetDesk")
	case "darwin":
		return filepath.Join(home, "Pictures", "SweetDesk")
	default:
		return filepath.Join(home, "Pictures", "SweetDesk")
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

// ProcessBatch processes a batch of images in the backend.
// Progress is emitted via Wails events so the frontend can re-render freely.
func (a *App) ProcessBatch(items []BatchItem, savePath string) {
	a.procMu.Lock()
	if a.procStatus.IsProcessing {
		a.procMu.Unlock()
		return // already processing
	}

	// Initialize status
	itemStatuses := make([]BatchItemStatus, len(items))
	for i, item := range items {
		itemStatuses[i] = BatchItemStatus{ID: item.ID, Status: "pending"}
	}
	a.procStatus = ProcessingStatus{
		IsProcessing: true,
		Total:        len(items),
		Current:      0,
		Progress:     0,
		Items:        itemStatuses,
		Done:         false,
	}
	a.procMu.Unlock()

	// Emit initial status
	a.emitProcessingStatus()

	// Process in a goroutine so the binding returns immediately
	go func() {
		defer func() {
			a.procMu.Lock()
			a.procStatus.IsProcessing = false
			a.procStatus.Done = true
			a.procStatus.Progress = 100
			a.procMu.Unlock()
			a.emitProcessingStatus()
		}()

		for i, item := range items {
			a.procMu.Lock()
			a.procStatus.Current = i
			a.procStatus.Items[i].Status = "processing"
			a.procMu.Unlock()
			a.emitProcessingStatus()

			err := a.processOneItem(item, savePath)

			a.procMu.Lock()
			if err != nil {
				log.Printf("❌ Failed to process item %s: %v", item.ID, err)
				a.procStatus.Items[i].Status = "error"
				a.procStatus.Items[i].Error = err.Error()
			} else {
				a.procStatus.Items[i].Status = "done"
			}
			a.procStatus.Progress = int(float64(i+1) / float64(len(items)) * 100)
			a.procMu.Unlock()
			a.emitProcessingStatus()
		}
	}()
}

// processOneItem handles downloading (if needed) and upscaling a single item
func (a *App) processOneItem(item BatchItem, savePath string) error {
	if a.upscaler == nil {
		return fmt.Errorf("upscaler not initialized")
	}

	// Get base64 data: either provided directly or download from URL
	base64Data := item.Base64Data
	if base64Data == "" && item.DownloadURL != "" {
		downloaded, err := a.DownloadImage(item.DownloadURL)
		if err != nil {
			return fmt.Errorf("failed to download image: %w", err)
		}
		base64Data = downloaded
	}
	if base64Data == "" {
		return fmt.Errorf("no image data available for item %s", item.ID)
	}

	// Parse dimensions from "WIDTHxHEIGHT"
	targetWidth, targetHeight := 3840, 2160
	if item.Dimension != "" {
		fmt.Sscanf(item.Dimension, "%dx%d", &targetWidth, &targetHeight)
	}

	// Sanitize filename
	fileName := item.Name
	if fileName == "" {
		fileName = fmt.Sprintf("wallpaper-%s.png", item.ID)
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".png"
	}

	// Process via existing ProcessImage pipeline
	_, err := a.ProcessImage(base64Data, targetWidth, targetHeight, savePath, fileName)
	return err
}

// GetProcessingStatus returns the current batch processing state.
// Called by the frontend on mount/re-mount to recover state.
func (a *App) GetProcessingStatus() ProcessingStatus {
	a.procMu.Lock()
	defer a.procMu.Unlock()
	// Return a copy
	status := a.procStatus
	itemsCopy := make([]BatchItemStatus, len(a.procStatus.Items))
	copy(itemsCopy, a.procStatus.Items)
	status.Items = itemsCopy
	return status
}

// emitProcessingStatus sends the current status to the frontend via Wails events
func (a *App) emitProcessingStatus() {
	a.procMu.Lock()
	status := a.procStatus
	itemsCopy := make([]BatchItemStatus, len(a.procStatus.Items))
	copy(itemsCopy, a.procStatus.Items)
	status.Items = itemsCopy
	a.procMu.Unlock()

	wailsRuntime.EventsEmit(a.ctx, "processing:status", status)
}

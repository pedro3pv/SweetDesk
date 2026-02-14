package main

import (
	"SweetDesk/internal/services"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pedro3pv/SweetDesk-core/pkg/types"
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
	coreBridge     *services.CoreBridge
	pixabayKey     string

	// Batch processing state
	procMu     sync.Mutex
	procStatus ProcessingStatus
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize services
	a.imageProcessor = services.NewImageProcessor(ctx)

	// Initialize SweetDesk-core bridge
	bridge, err := services.NewCoreBridge(ctx)
	if err != nil {
		// Do not terminate the entire application if the bridge fails to initialize.
		// Log the error and continue without the bridge (upscaling features disabled).
		log.Printf("SweetDesk-core bridge disabled: failed to initialize: %v", err)

		// Show a user-friendly warning dialog so the user understands that
		// upscaling features will be unavailable until the bridge is configured.
		_, _ = wailsRuntime.MessageDialog(ctx, wailsRuntime.MessageDialogOptions{
			Type:    wailsRuntime.WarningDialog,
			Title:   "Upscaling Unavailable",
			Message: "SweetDesk-core bridge could not be initialized. Upscaling features will be unavailable.\n\nDetails: " + err.Error(),
		})
	} else {
		a.coreBridge = bridge
		fmt.Println("✅ SweetDesk-core bridge initialized")
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
	// Cleanup SweetDesk-core bridge
	if a.coreBridge != nil {
		a.coreBridge.Close()
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
	if a.coreBridge == nil {
		return "", fmt.Errorf("coreBridge not initialized")
	}

	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", err
	}

	opts := &types.ProcessingOptions{
		TargetWidth:     0,
		TargetHeight:    0,
		ScaleFactor:     float64(scale),
		MaxResolution:   16384,
		KeepAspectRatio: true,
	}

	upscaled, err := a.coreBridge.UpscaleBytes(data, opts)
	if err != nil {
		return "", err
	}

	return a.imageProcessor.ConvertToBase64(upscaled), nil
}

// ProcessImage is the main processing pipeline
func (a *App) ProcessImage(base64Data string, targetWidth int, targetHeight int, savePath string, fileName string) (string, error) {
	if a.coreBridge == nil {
		return "", fmt.Errorf("coreBridge not initialized")
	}

	if targetWidth <= 0 || targetHeight <= 0 {
		return "", fmt.Errorf("invalid target resolution: %dx%d (dimensions must be positive)", targetWidth, targetHeight)
	}

	const defaultMaxResolution = 16384
	maxResolution := defaultMaxResolution

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

	data, err := a.imageProcessor.ConvertFromBase64(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	opts := &types.ProcessingOptions{
		TargetWidth:     targetWidth,
		TargetHeight:    targetHeight,
		ScaleFactor:     0,
		MaxResolution:   maxResolution,
		KeepAspectRatio: false,
	}

	upscaled, err := a.coreBridge.UpscaleBytes(data, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upscale: %w", err)
	}

	if savePath != "" && fileName != "" {
		_, err := a.imageProcessor.SaveToFile(upscaled, savePath, fileName)
		if err != nil {
			return "", fmt.Errorf("failed to save image: %w", err)
		}
	}

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

		// Prepare batch items for SweetDesk-core
		batchItems := make([]types.BatchItem, len(items))
		for i, item := range items {
			// Get base64 data: either provided directly or download from URL
			base64Data := item.Base64Data
			if base64Data == "" && item.DownloadURL != "" {
				downloaded, err := a.DownloadImage(item.DownloadURL)
				if err != nil {
					a.procMu.Lock()
					a.procStatus.Items[i].Status = "error"
					a.procStatus.Items[i].Error = err.Error()
					a.procMu.Unlock()
					a.emitProcessingStatus()
					continue
				}
				base64Data = downloaded
			}
			if base64Data == "" {
				a.procMu.Lock()
				a.procStatus.Items[i].Status = "error"
				a.procStatus.Items[i].Error = "no image data available"
				a.procMu.Unlock()
				a.emitProcessingStatus()
				continue
			}

			data, err := a.imageProcessor.ConvertFromBase64(base64Data)
			if err != nil {
				a.procMu.Lock()
				a.procStatus.Items[i].Status = "error"
				a.procStatus.Items[i].Error = err.Error()
				a.procMu.Unlock()
				a.emitProcessingStatus()
				continue
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

			// Save input to temp file
			if a.coreBridge == nil {
				a.procMu.Lock()
				a.procStatus.Items[i].Status = "error"
				a.procStatus.Items[i].Error = "core bridge not initialized"
				a.procMu.Unlock()
				a.emitProcessingStatus()
				continue
			}

			tmpDir := a.coreBridge.TmpDir
			if tmpDir == "" {
				tmpDir = os.TempDir()
			}

			tmpInput := filepath.Join(tmpDir, fmt.Sprintf("batch-%s-input.png", item.ID))
			if err := os.WriteFile(tmpInput, data, 0644); err != nil {
				a.procMu.Lock()
				a.procStatus.Items[i].Status = "error"
				a.procStatus.Items[i].Error = err.Error()
				a.procMu.Unlock()
				a.emitProcessingStatus()
				continue
			}
			// Ensure temporary input file is cleaned up when processing is done
			defer os.Remove(tmpInput)

			tmpOutput := filepath.Join(savePath, fileName)

			opts := &types.ProcessingOptions{
				TargetWidth:     targetWidth,
				TargetHeight:    targetHeight,
				ScaleFactor:     0,
				MaxResolution:   16384,
				KeepAspectRatio: false,
			}

			batchItems[i] = types.BatchItem{
				InputPath:  tmpInput,
				OutputPath: tmpOutput,
				Options:    opts,
			}
		}

		// Progress callback for UI
		progressCallback := func(current, total int, item types.BatchItem) {
			a.procMu.Lock()
			a.procStatus.Current = current - 1
			a.procStatus.Progress = int(float64(current) / float64(total) * 100)
			if current-1 < len(a.procStatus.Items) {
				a.procStatus.Items[current-1].Status = "processing"
			}
			a.procMu.Unlock()
			a.emitProcessingStatus()
		}

		// Call SweetDesk-core batch API
		if a.coreBridge != nil {
			_, err := a.coreBridge.ProcessBatch(batchItems, progressCallback)
			if err != nil {
				log.Printf("❌ Batch processing failed: %v", err)
			}
		}

		// Mark all items as done
		a.procMu.Lock()
		for i := range a.procStatus.Items {
			if a.procStatus.Items[i].Status == "processing" {
				a.procStatus.Items[i].Status = "done"
			}
		}
		a.procMu.Unlock()
		a.emitProcessingStatus()
	}()
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

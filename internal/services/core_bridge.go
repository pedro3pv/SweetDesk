package services

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
	"github.com/pedro3pv/SweetDesk-core/pkg/types"
)

// CoreBridge wraps SweetDesk-core processor for Wails integration.
// All upscaling and classification is delegated to the core library.
type CoreBridge struct {
	ctx       context.Context
	processor *processor.ImageProcessor
	tmpDir    string // temp directory for file-based operations
}

// NewCoreBridge creates a new CoreBridge instance.
// modelsFS should contain the embedded ONNX model files.
// If modelsFS is nil, models are expected at modelDir on the filesystem.
func NewCoreBridge(ctx context.Context, modelsFS fs.FS, modelDir string) (*CoreBridge, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create temp directory for intermediate files
	tmpDir, err := os.MkdirTemp("", "sweetdesk-bridge-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Build config for SweetDesk-core
	config := processor.Config{}

	if modelsFS != nil {
		config.ModelsFS = modelsFS
	} else if modelDir != "" {
		config.ModelDir = modelDir
	} else {
		os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("either modelsFS or modelDir must be provided")
	}

	// Optional ONNX Runtime path from environment
	if onnxPath := os.Getenv("ONNX_LIB_PATH"); onnxPath != "" {
		config.ONNXLibPath = onnxPath
	}

	// Create the core processor
	proc, err := processor.NewImageProcessor(config)
	if err != nil {
		os.RemoveAll(tmpDir)
		return nil, fmt.Errorf("failed to initialize SweetDesk-core processor: %w", err)
	}

	bridge := &CoreBridge{
		ctx:       ctx,
		processor: proc,
		tmpDir:    tmpDir,
	}

	log.Println("✅ SweetDesk-core bridge initialized")
	return bridge, nil
}

// ClassifyImage classifies an image (anime vs photo) from raw bytes.
// Returns the detected type and confidence score.
func (cb *CoreBridge) ClassifyImage(imageData []byte) (types.ImageType, float32, error) {
	if cb.processor == nil {
		return types.ImageTypeUnknown, 0, fmt.Errorf("processor not initialized")
	}

	// SweetDesk-core Classify expects a file path, write to temp
	tmpFile := filepath.Join(cb.tmpDir, "classify-input.png")
	if err := os.WriteFile(tmpFile, imageData, 0644); err != nil {
		return types.ImageTypeUnknown, 0, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	result, err := cb.processor.Classify(tmpFile)
	if err != nil {
		return types.ImageTypePhoto, 0, fmt.Errorf("classification failed: %w", err)
	}

	return result.Type, result.Confidence, nil
}

// UpscaleBytes upscales image bytes using SweetDesk-core with auto-classification.
// The core automatically selects the best model (RealCUGAN for anime, LSDIR for photos).
func (cb *CoreBridge) UpscaleBytes(imageData []byte, opts *types.ProcessingOptions) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	// Write input to temp file
	tmpInput := filepath.Join(cb.tmpDir, "upscale-input.png")
	if err := os.WriteFile(tmpInput, imageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input temp file: %w", err)
	}
	defer os.Remove(tmpInput)

	// Output temp file
	tmpOutput := filepath.Join(cb.tmpDir, "upscale-output.png")
	defer os.Remove(tmpOutput)

	// Process with options
	_, err := cb.processor.ProcessWithOptions(tmpInput, tmpOutput, opts)
	if err != nil {
		return nil, fmt.Errorf("upscaling failed: %w", err)
	}

	// Read result
	result, err := os.ReadFile(tmpOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}

	return result, nil
}

// ProcessFile processes a single image file (input path → output path).
// Combines classification + upscaling in one call.
func (cb *CoreBridge) ProcessFile(inputPath, outputPath string, opts *types.ProcessingOptions) (*types.ProcessingResult, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	return cb.processor.ProcessWithOptions(inputPath, outputPath, opts)
}

// ProcessBatch processes multiple images using SweetDesk-core batch API.
func (cb *CoreBridge) ProcessBatch(
	items []types.BatchItem,
	progressCallback types.ProgressCallback,
) (*types.BatchResult, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	return cb.processor.ProcessBatch(items, progressCallback)
}

// GetInfo returns processor information
func (cb *CoreBridge) GetInfo() map[string]interface{} {
	info := map[string]interface{}{
		"engine":  "SweetDesk-core",
		"status":  "ready",
		"version": "0.2.0",
	}
	return info
}

// Close releases all resources held by the bridge.
// Must be called when the application shuts down.
func (cb *CoreBridge) Close() error {
	log.Println("🗑️  Closing SweetDesk-core bridge...")

	var closeErr error
	if cb.processor != nil {
		closeErr = cb.processor.Close()
		cb.processor = nil
	}

	// Cleanup temp directory
	if cb.tmpDir != "" {
		os.RemoveAll(cb.tmpDir)
		cb.tmpDir = ""
	}

	log.Println("✅ SweetDesk-core bridge closed")
	return closeErr
}

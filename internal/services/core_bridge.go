package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
	"github.com/pedro3pv/SweetDesk-core/pkg/types"
)

// CoreBridge wraps SweetDesk-core processor for Wails integration.
// All upscaling and classification is delegated to the core library.
type CoreBridge struct {
	ctx       context.Context
	processor *processor.ImageProcessor
	TmpDir    string // temp directory for file-based operations
}

// NewCoreBridge creates a new CoreBridge instance.
// SweetDesk-core embeds models internally, so no modelsFS/modelDir is needed.
func NewCoreBridge(ctx context.Context) (*CoreBridge, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create temp directory for intermediate files
	tmpDir, err := os.MkdirTemp("", "sweetdesk-bridge-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Build config for SweetDesk-core
	// Core will auto-use its internally embedded models
	config := processor.Config{}

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
		TmpDir:    tmpDir,
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
	// Use unique temp file to avoid race conditions with concurrent operations
	tmpFile, err := os.CreateTemp(cb.TmpDir, "classify-*.png")
	if err != nil {
		return types.ImageTypeUnknown, 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := os.WriteFile(tmpPath, imageData, 0644); err != nil {
		return types.ImageTypeUnknown, 0, fmt.Errorf("failed to write temp file: %w", err)
	}

	result, err := cb.processor.Classify(tmpPath)
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

	// Create unique temp files to avoid race conditions with concurrent operations
	tmpInputFile, err := os.CreateTemp(cb.TmpDir, "upscale-input-*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create input temp file: %w", err)
	}
	tmpInput := tmpInputFile.Name()
	tmpInputFile.Close()
	defer os.Remove(tmpInput)

	if err := os.WriteFile(tmpInput, imageData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input temp file: %w", err)
	}

	// Create unique output temp file
	tmpOutputFile, err := os.CreateTemp(cb.TmpDir, "upscale-output-*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create output temp file: %w", err)
	}
	tmpOutput := tmpOutputFile.Name()
	tmpOutputFile.Close()
	defer os.Remove(tmpOutput)

	// Process with options
	_, err = cb.processor.ProcessWithOptions(tmpInput, tmpOutput, opts)
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
	if cb.TmpDir != "" {
		os.RemoveAll(cb.TmpDir)
		cb.TmpDir = ""
	}

	log.Println("✅ SweetDesk-core bridge closed")
	return closeErr
}

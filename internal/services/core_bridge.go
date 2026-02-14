package services

import (
	"context"
	"fmt"
	"os"

	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

// CoreBridge adapts SweetDesk-core processor to our existing services interface
type CoreBridge struct {
	ctx       context.Context
	processor *processor.Processor
}

// NewCoreBridge creates a new core bridge instance
func NewCoreBridge(ctx context.Context) (*CoreBridge, error) {
	modelDir := os.Getenv("SWEETDESK_MODEL_DIR")
	if modelDir == "" {
		modelDir = "./models"
	}

	onnxLibPath := os.Getenv("ONNX_LIB_PATH")

	p, err := processor.New(processor.Config{
		ModelDir:    modelDir,
		ONNXLibPath: onnxLibPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SweetDesk-core processor: %w", err)
	}

	return &CoreBridge{
		ctx:       ctx,
		processor: p,
	}, nil
}

// Close closes the core bridge
func (cb *CoreBridge) Close() error {
	if cb.processor != nil {
		cb.processor.Close()
	}
	return nil
}

// Classify classifies an image
func (cb *CoreBridge) Classify(imageData []byte) (string, error) {
	if cb.processor == nil {
		return "", fmt.Errorf("processor not initialized")
	}

	result, err := cb.processor.Classify(imageData)
	if err != nil {
		return "", err
	}

	return result.Type, nil
}

// Upscale upscales an image
func (cb *CoreBridge) Upscale(imageData []byte, imageType string, scale int) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	options := processor.ProcessOptions{
		ImageType:    imageType,
		TargetScale:  scale,
		OutputFormat: "png",
	}

	return cb.processor.Upscale(imageData, options)
}

// UpscaleToResolution upscales an image to specific dimensions
func (cb *CoreBridge) UpscaleToResolution(imageData []byte, width int, height int, keepAspectRatio bool) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	options := processor.ProcessOptions{
		TargetWidth:     width,
		TargetHeight:    height,
		KeepAspectRatio: keepAspectRatio,
		OutputFormat:    "png",
	}

	return cb.processor.UpscaleToResolution(imageData, options)
}

// AdjustAspectRatio adjusts the aspect ratio using seam carving
func (cb *CoreBridge) AdjustAspectRatio(imageData []byte, targetWidth int, targetHeight int) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	options := processor.ProcessOptions{
		TargetWidth:    targetWidth,
		TargetHeight:   targetHeight,
		UseSeamCarving: true,
		OutputFormat:   "png",
	}

	return cb.processor.AdjustAspectRatio(imageData, options)
}

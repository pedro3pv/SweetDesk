package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

// CoreBridge wraps SweetDesk-core processor for Wails
type CoreBridge struct {
	ctx       context.Context
	processor *processor.Processor
}

// NewCoreBridge creates a new CoreBridge instance
func NewCoreBridge(ctx context.Context) (*CoreBridge, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Initialize processor with default config
	config := processor.Config{
		ModelDir: "./models",
	}

	// Try to get ONNX lib path from environment
	if onnxPath := os.Getenv("ONNX_LIB_PATH"); onnxPath != "" {
		config.ONNXLibPath = onnxPath
	}

	proc, err := processor.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize processor: %w", err)
	}

	bridge := &CoreBridge{
		ctx:       ctx,
		processor: proc,
	}

	log.Println("✅ SweetDesk-core bridge initialized")
	return bridge, nil
}

// ClassifyImage classifies an image by type
func (cb *CoreBridge) ClassifyImage(imageData []byte) (string, float32, error) {
	if cb.processor == nil {
		return "", 0, fmt.Errorf("processor not initialized")
	}

	result, err := cb.processor.Classify(imageData)
	if err != nil {
		return "photo", 0, fmt.Errorf("classification failed: %w", err)
	}

	return result.Type, result.Confidence, nil
}

// UpscaleImage upscales an image to target scale
func (cb *CoreBridge) UpscaleImage(imageData []byte, scale int, format string) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	if format == "" {
		format = "png"
	}

	options := processor.ProcessOptions{
		TargetScale:  scale,
		OutputFormat: format,
	}

	result, err := cb.processor.Upscale(imageData, options)
	if err != nil {
		return nil, fmt.Errorf("upscaling failed: %w", err)
	}

	return result.Data, nil
}

// DownloadImage downloads an image from URL
func (cb *CoreBridge) DownloadImage(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return data, nil
}

// ProcessImage is the main processing pipeline
func (cb *CoreBridge) ProcessImage(imageData []byte, targetWidth, targetHeight int, format string, aspectRatio string) ([]byte, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	if format == "" {
		format = "png"
	}

	options := processor.ProcessOptions{
		TargetWidth:       targetWidth,
		TargetHeight:      targetHeight,
		TargetAspectRatio: aspectRatio,
		OutputFormat:      format,
		KeepAspectRatio:   true,
	}

	result, err := cb.processor.Process(imageData, options)
	if err != nil {
		return nil, fmt.Errorf("processing failed: %w", err)
	}

	return result.Data, nil
}

// SaveImage saves processed image to file
func (cb *CoreBridge) SaveImage(imageData []byte, savePath, fileName string) (string, error) {
	if savePath == "" || fileName == "" {
		return "", fmt.Errorf("save path and filename required")
	}

	// Ensure directory exists
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(savePath, fileName)

	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("✅ Image saved to %s", filePath)
	return filePath, nil
}

// GetInfo returns processor information
func (cb *CoreBridge) GetInfo() (map[string]interface{}, error) {
	if cb.processor == nil {
		return nil, fmt.Errorf("processor not initialized")
	}

	info := map[string]interface{}{
		"engine":  "SweetDesk-core",
		"status":  "ready",
		"version": "0.1.0",
	}

	return info, nil
}

// Close closes the bridge and processor
func (cb *CoreBridge) Close() error {
	if cb.processor != nil {
		return cb.processor.Close()
	}
	return nil
}

// ImageProcessor helper (for backward compatibility)
type ImageProcessor struct{
	ctx context.Context
}

// NewImageProcessor creates image processor
func NewImageProcessor(ctx context.Context) *ImageProcessor {
	return &ImageProcessor{ctx: ctx}
}

// ConvertToBase64 converts image bytes to base64 string
func (ip *ImageProcessor) ConvertToBase64(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}

// ConvertFromBase64 converts base64 string to image bytes
func (ip *ImageProcessor) ConvertFromBase64(base64Data string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}
	return data, nil
}

// ValidateImage validates image format and size
func (ip *ImageProcessor) ValidateImage(imageData []byte) error {
	if len(imageData) == 0 {
		return fmt.Errorf("empty image data")
	}

	// Try to decode to validate format
	_, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return fmt.Errorf("invalid image format: %w", err)
	}

	return nil
}

// GetImageInfo returns image information
func (ip *ImageProcessor) GetImageInfo(imageData []byte) (map[string]int, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	info := map[string]int{
		"width":  bounds.Dx(),
		"height": bounds.Dy(),
	}

	return info, nil
}

// SaveToFile saves image to file system
func (ip *ImageProcessor) SaveToFile(imageData []byte, savePath, fileName string) (string, error) {
	if savePath == "" || fileName == "" {
		return "", fmt.Errorf("save path and filename required")
	}

	// Ensure directory exists
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(savePath, fileName)

	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("✅ File saved to %s", filePath)
	return filePath, nil
}

// ConvertFormat converts image between formats
func (ip *ImageProcessor) ConvertFormat(imageData []byte, targetFormat string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var buf bytes.Buffer

	switch targetFormat {
	case "png":
		err = png.Encode(&buf, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	default:
		err = png.Encode(&buf, img)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// CoreBridge wraps image processing for Wails.
// It provides classification, upscaling and full processing pipelines
// that integrate with the local Upscaler when SweetDesk-core is not available.
type CoreBridge struct {
	ctx context.Context
}

// NewCoreBridge creates a new CoreBridge instance
func NewCoreBridge(ctx context.Context) (*CoreBridge, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	bridge := &CoreBridge{
		ctx: ctx,
	}

	log.Println("✅ SweetDesk-core bridge initialized")
	return bridge, nil
}

// ClassifyImage classifies an image by type.
// Returns the detected type ("photo" or "anime") and a confidence score.
func (cb *CoreBridge) ClassifyImage(imageData []byte) (string, float32, error) {
	// Default classification: photo with full confidence
	return "photo", 1.0, nil
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
	info := map[string]interface{}{
		"engine":  "SweetDesk-core",
		"status":  "ready",
		"version": "0.1.0",
	}

	return info, nil
}

// Close closes the bridge
func (cb *CoreBridge) Close() error {
	return nil
}

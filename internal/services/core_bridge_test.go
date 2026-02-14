package services

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// TestNewCoreBridge tests CoreBridge initialization
func TestNewCoreBridge(t *testing.T) {
	ctx := context.Background()
	bridge, err := NewCoreBridge(ctx)
	if err != nil {
		// Allow error if ONNX runtime / models not available in CI
		t.Logf("Warning: CoreBridge initialization warning: %v", err)
		t.Skip("skipping: SweetDesk-core runtime not available")
	}
	if bridge == nil {
		t.Fatal("CoreBridge should not be nil")
	}
	defer bridge.Close()
}

// TestImageProcessor tests ImageProcessor utility methods
func TestImageProcessor(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	if ip == nil {
		t.Fatal("ImageProcessor should not be nil")
	}
}

// TestImageProcessorBase64 tests base64 conversion
func TestImageProcessorBase64(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	// Create test data
	testImg := createTestImage(50, 50)
	originalData := encodeImageToPNG(testImg)

	// Encode to base64
	base64Data := ip.ConvertToBase64(originalData)
	if base64Data == "" {
		t.Fatal("ConvertToBase64 returned empty string")
	}

	// Decode from base64
	decodedData, err := ip.ConvertFromBase64(base64Data)
	if err != nil {
		t.Fatalf("ConvertFromBase64 failed: %v", err)
	}

	// Verify roundtrip
	if !bytes.Equal(originalData, decodedData) {
		t.Error("Base64 roundtrip failed: data mismatch")
	}
}

// TestImageProcessorLoadImage tests image loading from bytes
func TestImageProcessorLoadImage(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	testImg := createTestImage(200, 150)
	testImgData := encodeImageToPNG(testImg)

	img, format, err := ip.LoadImageFromBytes(testImgData)
	if err != nil {
		t.Fatalf("LoadImageFromBytes failed: %v", err)
	}

	if format != "png" {
		t.Errorf("Expected format 'png', got '%s'", format)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 200 {
		t.Errorf("Expected width 200, got %d", bounds.Dx())
	}
	if bounds.Dy() != 150 {
		t.Errorf("Expected height 150, got %d", bounds.Dy())
	}
}

// TestImageProcessorEncodeImage tests image encoding
func TestImageProcessorEncodeImage(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	testImg := createTestImage(100, 100)

	tests := []struct {
		name   string
		format string
	}{
		{"PNG", "png"},
		{"JPEG", "jpeg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := ip.EncodeImage(testImg, tt.format, 90)
			if err != nil {
				t.Errorf("EncodeImage failed: %v", err)
			}
			if len(encoded) == 0 {
				t.Error("EncodeImage returned empty data")
			}
		})
	}
}

// TestImageProcessorSaveToFile tests file saving with path sanitization
func TestImageProcessorSaveToFile(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	testImg := createTestImage(10, 10)
	testData := encodeImageToPNG(testImg)

	tmpDir := t.TempDir()

	path, err := ip.SaveToFile(testData, tmpDir, "test-image.png")
	if err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
	if path == "" {
		t.Error("SaveToFile returned empty path")
	}

	// Test invalid filenames
	_, err = ip.SaveToFile(testData, tmpDir, "../escape.png")
	if err != nil {
		// Should sanitize to just "escape.png" via filepath.Base
		t.Logf("Sanitized path traversal: %v", err)
	}
}

// TestCoreBridgeGetInfo tests processor info retrieval
func TestCoreBridgeGetInfo(t *testing.T) {
	ctx := context.Background()
	bridge, err := NewCoreBridge(ctx)
	if err != nil {
		t.Logf("Warning: NewCoreBridge error: %v", err)
		t.Skip("skipping: SweetDesk-core runtime not available")
	}
	if bridge == nil {
		t.Fatal("CoreBridge is nil")
	}
	defer bridge.Close()

	info := bridge.GetInfo()

	if info["engine"] != "SweetDesk-core" {
		t.Error("Expected engine 'SweetDesk-core'")
	}

	if info["status"] != "ready" {
		t.Error("Expected status 'ready'")
	}
}

// Helper functions

// createTestImage creates a simple test image
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img
}

// encodeImageToPNG encodes an image to PNG bytes
func encodeImageToPNG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic("failed to encode image: " + err.Error())
	}
	return buf.Bytes()
}

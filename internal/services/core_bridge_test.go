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
		// Allow error if processor not fully configured, but bridge should still be created
		t.Logf("Warning: CoreBridge initialization warning: %v", err)
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

	// Create test image
	testImg := createTestImage(100, 100)
	testImgData := encodeImageToPNG(testImg)

	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{
			name: "ValidateImage_ValidPNG",
			fn: func() error {
				return ip.ValidateImage(testImgData)
			},
			wantErr: false,
		},
		{
			name: "ValidateImage_Empty",
			fn: func() error {
				return ip.ValidateImage([]byte{})
			},
			wantErr: true,
		},
		{
			name: "ValidateImage_Invalid",
			fn: func() error {
				return ip.ValidateImage([]byte{0xFF, 0xD8, 0xFF}) // Incomplete JPEG header
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Errorf("Got error %v, want %v", err, tt.wantErr)
			}
		})
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

// TestImageProcessorGetImageInfo tests image info extraction
func TestImageProcessorGetImageInfo(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	testImg := createTestImage(200, 150)
	testImgData := encodeImageToPNG(testImg)

	info, err := ip.GetImageInfo(testImgData)
	if err != nil {
		t.Fatalf("GetImageInfo failed: %v", err)
	}

	if info["width"] != 200 {
		t.Errorf("Expected width 200, got %d", info["width"])
	}

	if info["height"] != 150 {
		t.Errorf("Expected height 150, got %d", info["height"])
	}
}

// TestImageProcessorConvertFormat tests format conversion
func TestImageProcessorConvertFormat(t *testing.T) {
	ctx := context.Background()
	ip := NewImageProcessor(ctx)

	testImg := createTestImage(100, 100)
	originalData := encodeImageToPNG(testImg)

	tests := []struct {
		name   string
		format string
	}{
		{"PNG", "png"},
		{"JPEG", "jpeg"},
		{"JPG", "jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converted, err := ip.ConvertFormat(originalData, tt.format)
			if err != nil {
				t.Errorf("ConvertFormat failed: %v", err)
			}
			if len(converted) == 0 {
				t.Error("ConvertFormat returned empty data")
			}
		})
	}
}

// TestCoreBridgeGetInfo tests processor info retrieval
func TestCoreBridgeGetInfo(t *testing.T) {
	ctx := context.Background()
	bridge, err := NewCoreBridge(ctx)
	if err != nil {
		t.Logf("Warning: NewCoreBridge error: %v", err)
	}
	if bridge == nil {
		t.Fatal("CoreBridge is nil")
	}
	defer bridge.Close()

	info, err := bridge.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}

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

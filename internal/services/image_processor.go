package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

// ImageProcessor handles all image processing operations
type ImageProcessor struct {
	ctx context.Context
}

// NewImageProcessor creates a new image processor instance
func NewImageProcessor(ctx context.Context) *ImageProcessor {
	return &ImageProcessor{
		ctx: ctx,
	}
}

// ProcessingOptions contains options for image processing
type ProcessingOptions struct {
	TargetResolution string // "4K", "5K", "8K"
	AspectRatio      string // "16:9", "21:9", "auto"
	UseSeamCarving   bool   // true for content-aware, false for crop
	Quality          int    // JPEG quality (1-100)
}

// ProcessingResult contains the result of image processing
type ProcessingResult struct {
	ImageData   []byte
	Width       int
	Height      int
	Format      string
	ModelUsed   string
	ProcessTime float64
}

// LoadImageFromBytes loads an image from byte array
func (ip *ImageProcessor) LoadImageFromBytes(data []byte) (image.Image, string, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}
	return img, format, nil
}

// EncodeImage encodes an image to bytes
func (ip *ImageProcessor) EncodeImage(img image.Image, format string, quality int) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("failed to encode jpeg: %w", err)
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, fmt.Errorf("failed to encode png: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Bytes(), nil
}

// ConvertToBase64 converts image bytes to base64 string
func (ip *ImageProcessor) ConvertToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ConvertFromBase64 converts base64 string to image bytes
func (ip *ImageProcessor) ConvertFromBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

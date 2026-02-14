package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
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

// SaveToFile saves image bytes to a file at the given path
func (ip *ImageProcessor) SaveToFile(data []byte, savePath string, fileName string) (string, error) {
	if savePath == "" {
		return "", fmt.Errorf("save path cannot be empty")
	}

	// Expand ~ or ~/ (and ~\ on Windows) to home directory
	if len(savePath) > 0 && savePath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		// Only ~ or ~/ or ~\
		if savePath == "~" {
			savePath = home
		} else if len(savePath) > 1 && (savePath[1] == '/' || savePath[1] == '\\') {
			// Trim the separator: ~/Pictures becomes "Pictures"
			rel := savePath[2:]
			savePath = filepath.Join(home, rel)
		}
		// If path is like ~username, do not expand (leave as is)
	}

	// Ensure directory exists
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", savePath, err)
	}

	// Sanitize fileName: extract only the base name
	safeName := filepath.Base(filepath.Clean(fileName))

	// Reject invalid filenames
	if safeName == "." || safeName == ".." || safeName == "" || safeName == string(filepath.Separator) {
		return "", fmt.Errorf("invalid filename: %s", fileName)
	}

	// Reject any remaining path separators (defense in depth)
	if containsSeparator(safeName) {
		return "", fmt.Errorf("filename contains path separators: %s", fileName)
	}

	// Reject Windows drive letters (e.g., C:image.png)
	if len(safeName) >= 2 && safeName[1] == ':' {
		return "", fmt.Errorf("filename contains drive letter: %s", fileName)
	}

	// Check for Windows reserved names
	if isWindowsReserved(safeName) {
		return "", fmt.Errorf("filename uses Windows reserved name: %s", fileName)
	}

	// Add default extension if missing
	ext := filepath.Ext(safeName)
	if ext == "" {
		safeName += ".png"
	}

	fullPath := filepath.Join(savePath, safeName)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	log.Printf("✅ Image saved: %s", fullPath)
	return fullPath, nil
}

// containsSeparator checks if a string contains any path separator
func containsSeparator(name string) bool {
	return strings.Contains(name, "/") || strings.Contains(name, "\\")
}

// isWindowsReserved checks if filename uses a Windows reserved name
func isWindowsReserved(name string) bool {
	// Remove extension to check base name
	baseName := strings.TrimSuffix(name, filepath.Ext(name))
	baseName = strings.ToUpper(baseName)

	// Windows reserved device names
	reserved := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	for _, r := range reserved {
		if baseName == r {
			return true
		}
	}
	return false
}

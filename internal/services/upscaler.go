package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Upscaler handles image upscaling operations
type Upscaler struct {
	ctx        context.Context
	binariesDir string
}

// NewUpscaler creates a new upscaler instance
func NewUpscaler(ctx context.Context) *Upscaler {
	binariesDir := filepath.Join(getCurrentDir(), "binaries", runtime.GOOS)
	
	return &Upscaler{
		ctx:        ctx,
		binariesDir: binariesDir,
	}
}

// UpscaleOptions contains options for upscaling
type UpscaleOptions struct {
	Model      string // "realesrgan-anime", "realesrgan-photo", "realcugan"
	Scale      int    // 2, 4, etc.
	TileSize   int    // GPU tile size (default 0 = auto)
	Format     string // "png", "jpg"
}

// UpscaleImage upscales an image using the specified model
func (u *Upscaler) UpscaleImage(imageData []byte, options UpscaleOptions) ([]byte, error) {
	// Determine which binary to use based on model
	var binaryName string
	var modelPath string
	
	switch options.Model {
	case "realesrgan-anime":
		binaryName = u.getBinaryName("realesrgan-ncnn-vulkan")
		modelPath = filepath.Join(u.binariesDir, "models", "realesrgan-x4plus-anime")
	case "realesrgan-photo":
		binaryName = u.getBinaryName("realesrgan-ncnn-vulkan")
		modelPath = filepath.Join(u.binariesDir, "models", "realesrgan-x4plus")
	case "realcugan":
		binaryName = u.getBinaryName("realcugan-ncnn-vulkan")
		modelPath = filepath.Join(u.binariesDir, "models", "models-pro")
	default:
		return nil, fmt.Errorf("unsupported model: %s", options.Model)
	}
	
	binaryPath := filepath.Join(u.binariesDir, binaryName)
	
	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("upscaler binary not found at %s", binaryPath)
	}
	
	// Create temporary files
	tmpInput, err := os.CreateTemp("", "upscale_input_*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create input temp file: %w", err)
	}
	defer os.Remove(tmpInput.Name())
	
	tmpOutputDir, err := os.MkdirTemp("", "upscale_output_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create output temp dir: %w", err)
	}
	defer os.RemoveAll(tmpOutputDir)
	
	if _, err := tmpInput.Write(imageData); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}
	tmpInput.Close()
	
	// Build command arguments
	args := []string{
		"-i", tmpInput.Name(),
		"-o", tmpOutputDir,
		"-n", filepath.Base(modelPath),
		"-s", fmt.Sprintf("%d", options.Scale),
		"-f", options.Format,
	}
	
	if options.TileSize > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", options.TileSize))
	}
	
	// Execute upscaler
	cmd := exec.CommandContext(u.ctx, binaryPath, args...)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("upscaling failed: %w, stderr: %s", err, stderr.String())
	}
	
	// Read output file
	outputName := filepath.Base(tmpInput.Name())
	outputName = outputName[:len(outputName)-len(filepath.Ext(outputName))] + "." + options.Format
	outputPath := filepath.Join(tmpOutputDir, outputName)
	
	result, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}
	
	return result, nil
}

// getBinaryName returns the platform-specific binary name
func (u *Upscaler) getBinaryName(baseName string) string {
	if runtime.GOOS == "windows" {
		return baseName + ".exe"
	}
	return baseName
}

// GetRecommendedModel returns the recommended upscaling model based on image type
func (u *Upscaler) GetRecommendedModel(imageType string) string {
	switch imageType {
	case "anime":
		return "realcugan"
	case "photo":
		return "realesrgan-photo"
	default:
		return "realesrgan-photo"
	}
}

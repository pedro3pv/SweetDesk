package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// PythonBridge handles communication with Python scripts
type PythonBridge struct {
	ctx        context.Context
	pythonPath string
	scriptsDir string
}

// NewPythonBridge creates a new Python bridge
func NewPythonBridge(ctx context.Context) (*PythonBridge, error) {
	pythonPath, err := findPython()
	if err != nil {
		return nil, fmt.Errorf("Python not found: %w", err)
	}
	
	scriptsDir := filepath.Join(getCurrentDir(), "python")
	
	return &PythonBridge{
		ctx:        ctx,
		pythonPath: pythonPath,
		scriptsDir: scriptsDir,
	}, nil
}

// ClassificationResult represents the result of image classification
type ClassificationResult struct {
	Type       string  `json:"type"`        // "anime", "photo", "art"
	Confidence float64 `json:"confidence"`
	Model      string  `json:"model"`
}

// ClassifyImage classifies an image as anime or photo
func (pb *PythonBridge) ClassifyImage(imageData []byte) (*ClassificationResult, error) {
	scriptPath := filepath.Join(pb.scriptsDir, "classify_image.py")
	
	// Create temporary file for image
	tmpFile, err := os.CreateTemp("", "img_*.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.Write(imageData); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()
	
	cmd := exec.CommandContext(pb.ctx, pb.pythonPath, scriptPath, tmpFile.Name())
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("classification failed: %w, stderr: %s", err, stderr.String())
	}
	
	var result ClassificationResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %w", err)
	}
	
	return &result, nil
}

// SeamCarvingOptions contains options for seam carving
type SeamCarvingOptions struct {
	TargetWidth  int
	TargetHeight int
	EnergyMode   string // "forward", "backward"
}

// ApplySeamCarving applies seam carving to resize an image
func (pb *PythonBridge) ApplySeamCarving(imageData []byte, options SeamCarvingOptions) ([]byte, error) {
	scriptPath := filepath.Join(pb.scriptsDir, "seam_carving.py")
	
	// Create temporary files
	tmpInput, err := os.CreateTemp("", "input_*.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create input temp file: %w", err)
	}
	defer os.Remove(tmpInput.Name())
	
	tmpOutput, err := os.CreateTemp("", "output_*.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create output temp file: %w", err)
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()
	
	if _, err := tmpInput.Write(imageData); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}
	tmpInput.Close()
	
	cmd := exec.CommandContext(pb.ctx, pb.pythonPath, scriptPath,
		tmpInput.Name(),
		tmpOutput.Name(),
		fmt.Sprintf("%d", options.TargetWidth),
		fmt.Sprintf("%d", options.TargetHeight),
		options.EnergyMode,
	)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("seam carving failed: %w, stderr: %s", err, stderr.String())
	}
	
	result, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}
	
	return result, nil
}

// findPython finds the Python executable
func findPython() (string, error) {
	candidates := []string{"python3", "python"}
	
	for _, candidate := range candidates {
		path, err := exec.LookPath(candidate)
		if err == nil {
			return path, nil
		}
	}
	
	return "", fmt.Errorf("Python executable not found")
}

// getCurrentDir returns the current working directory
func getCurrentDir() string {
	if runtime.GOOS == "darwin" {
		// On macOS, use the executable directory
		exe, err := os.Executable()
		if err == nil {
			return filepath.Dir(exe)
		}
	}
	
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

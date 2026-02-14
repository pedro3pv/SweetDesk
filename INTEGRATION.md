# SweetDesk-Core Integration Guide

## Overview

This guide covers integrating SweetDesk-Core into your Wails application for advanced image processing capabilities.

## Quick Start

### 1. Installation

```bash
go get -u github.com/pedro3pv/SweetDesk-core
go mod tidy
go generate ./...
```

### 2. Basic Usage

```go
package main

import (
    "log"
    "github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

func main() {
    // Initialize processor
    proc, err := processor.New(processor.Config{
        ModelDir: "./models",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer proc.Close()
    
    // Process image
    result, err := proc.Upscale(imageData, processor.ProcessOptions{
        TargetScale:  4,
        OutputFormat: "png",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Use result
    log.Printf("Upscaled to %dx%d", result.Width, result.Height)
}
```

---

## API Reference

### Processor Interface

#### Creating a Processor

```go
config := processor.Config{
    ModelDir:     "./models",      // Optional: model directory path
    ONNXLibPath:  "/path/to/lib",  // Optional: ONNX Runtime library path
    CacheSize:    1000,             // Optional: cache size in MB
    MaxWorkers:   4,                // Optional: number of worker threads
}

proc, err := processor.New(config)
if err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
defer proc.Close()
```

#### Available Methods

```go
type Processor interface {
    // Image classification
    Classify(imageData []byte) (*ClassifyResult, error)
    
    // Upscaling (4x fixed)
    Upscale(imageData []byte, opts ProcessOptions) (*UpscaleResult, error)
    
    // Content-aware resizing
    SeamCarve(imageData []byte, targetWidth, targetHeight int) (*SeamCarvingResult, error)
    
    // Full pipeline: classify -> upscale -> seam carve
    Process(imageData []byte, opts ProcessOptions) (*ProcessResult, error)
    
    // Resource cleanup
    Close() error
}
```

---

## Detailed Methods

### Classification

```go
result, err := proc.Classify(imageData)
if err != nil {
    return fmt.Errorf("classification failed: %w", err)
}

switch result.Type {
case processor.TypePhoto:
    log.Println("Detected: Photo")
case processor.TypeArtwork:
    log.Println("Detected: Artwork")
case processor.TypeScreenshot:
    log.Println("Detected: Screenshot")
default:
    log.Println("Unknown type")
}

log.Printf("Confidence: %.2f%%", result.Confidence*100)
```

**Result Structure:**

```go
type ClassifyResult struct {
    Type       string      // "photo", "artwork", "screenshot"
    Confidence float32     // 0.0-1.0
    Labels     []string    // Additional labels
    Metadata   map[string]interface{}
}
```

### Upscaling

```go
options := processor.ProcessOptions{
    TargetScale:   4,           // 2x, 4x, or 8x
    OutputFormat:  "png",       // "png" or "jpg"
    JPEGQuality:   95,          // 1-100 (if JPEG)
}

result, err := proc.Upscale(imageData, options)
if err != nil {
    return fmt.Errorf("upscaling failed: %w", err)
}

log.Printf("Upscaled %dx%d -> %dx%d",
    result.OriginalWidth, result.OriginalHeight,
    result.Width, result.Height)
```

**Result Structure:**

```go
type UpscaleResult struct {
    Data              []byte
    Width, Height     int
    OriginalWidth     int
    OriginalHeight    int
    Format            string
    ProcessingTime    time.Duration
    ModelUsed         string
}
```

### Content-Aware Resizing (Seam Carving)

```go
result, err := proc.SeamCarve(imageData, 1920, 1080)
if err != nil {
    return fmt.Errorf("seam carving failed: %w", err)
}

log.Printf("Resized to %dx%d",
    result.Width, result.Height)
log.Printf("Removed %d vertical, %d horizontal seams",
    result.VerticalSeamsRemoved,
    result.HorizontalSeamsRemoved)
```

**Result Structure:**

```go
type SeamCarvingResult struct {
    Data                  []byte
    Width, Height         int
    OriginalWidth         int
    OriginalHeight        int
    VerticalSeamsRemoved  int
    HorizontalSeamsRemoved int
    ProcessingTime        time.Duration
}
```

### Full Processing Pipeline

```go
options := processor.ProcessOptions{
    TargetWidth:        3840,           // Target resolution
    TargetHeight:       2160,           // 4K
    TargetAspectRatio:  "16:9",         // Aspect ratio constraint
    UseSeamCarving:     true,           // Enable content-aware resizing
    KeepAspectRatio:    true,           // Preserve original aspect ratio
    OutputFormat:       "png",          // Output format
}

result, err := proc.Process(imageData, options)
if err != nil {
    return fmt.Errorf("processing failed: %w", err)
}

log.Printf("Pipeline complete:")
log.Printf("  - Classification: %s (%.0f%% confidence)",
    result.Classification.Type,
    result.Classification.Confidence*100)
log.Printf("  - Upscaling: %dx%d -> %dx%d",
    result.OriginalWidth, result.OriginalHeight,
    result.UpscaledWidth, result.UpscaledHeight)
log.Printf("  - Total time: %v", result.TotalProcessingTime)
```

**Result Structure:**

```go
type ProcessResult struct {
    Data              []byte
    Width, Height     int
    OriginalWidth     int
    OriginalHeight    int
    UpscaledWidth     int
    UpscaledHeight    int
    Classification    ClassifyResult
    Format            string
    ModelsUsed        []string
    TotalProcessingTime time.Duration
}
```

---

## Integration Patterns

### Pattern 1: Single-Use Processing

For occasional image processing:

```go
func ProcessImageOnce(imageData []byte) ([]byte, error) {
    proc, err := processor.New(processor.Config{
        ModelDir: "./models",
    })
    if err != nil {
        return nil, err
    }
    defer proc.Close()
    
    result, err := proc.Upscale(imageData, processor.ProcessOptions{
        TargetScale:  4,
        OutputFormat: "png",
    })
    if err != nil {
        return nil, err
    }
    
    return result.Data, nil
}
```

### Pattern 2: Application-Level Singleton

For persistent processor instance (recommended):

```go
type App struct {
    processor *processor.Processor
}

func (a *App) startup(ctx context.Context) {
    var err error
    a.processor, err = processor.New(processor.Config{
        ModelDir: "./models",
    })
    if err != nil {
        fmt.Printf("Warning: processor init failed: %v\n", err)
        a.processor = nil
    }
}

func (a *App) shutdown(ctx context.Context) {
    if a.processor != nil {
        a.processor.Close()
    }
}

func (a *App) ProcessImage(base64Data string) (string, error) {
    if a.processor == nil {
        return "", fmt.Errorf("processor not available")
    }
    
    // Process...
    result, err := a.processor.Process(imageData, options)
    // Return base64 encoded result...
}
```

### Pattern 3: Worker Pool

For high-throughput processing:

```go
type ProcessorPool struct {
    queue    chan ProcessTask
    results  chan ProcessResult
    workers  int
    procs    []*processor.Processor
}

func NewProcessorPool(size int) (*ProcessorPool, error) {
    pool := &ProcessorPool{
        queue:   make(chan ProcessTask, 10),
        results: make(chan ProcessResult),
        workers: size,
        procs:   make([]*processor.Processor, size),
    }
    
    // Initialize worker processors
    for i := 0; i < size; i++ {
        proc, err := processor.New(processor.Config{
            ModelDir: "./models",
        })
        if err != nil {
            return nil, err
        }
        pool.procs[i] = proc
        
        // Start worker
        go pool.worker(proc)
    }
    
    return pool, nil
}

func (p *ProcessorPool) worker(proc *processor.Processor) {
    for task := range p.queue {
        result, err := proc.Process(task.Data, task.Options)
        p.results <- ProcessResult{
            Result: result,
            Error:  err,
        }
    }
}
```

### Pattern 4: Caching Layer

For repeated image processing:

```go
type CachedProcessor struct {
    proc  *processor.Processor
    cache map[string][]byte
}

func (cp *CachedProcessor) ProcessCached(imageHash string, imageData []byte, opts processor.ProcessOptions) ([]byte, error) {
    // Check cache
    if cached, ok := cp.cache[imageHash]; ok {
        return cached, nil
    }
    
    // Process
    result, err := cp.proc.Process(imageData, opts)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    cp.cache[imageHash] = result.Data
    
    return result.Data, nil
}
```

---

## Error Handling

### Graceful Degradation

```go
func (a *App) UpscaleImage(imageData []byte, scale int) ([]byte, error) {
    if a.processor == nil {
        // Fallback: return original image
        fmt.Println("Warning: processor unavailable, returning original")
        return imageData, nil
    }
    
    result, err := a.processor.Upscale(imageData, processor.ProcessOptions{
        TargetScale:  scale,
        OutputFormat: "png",
    })
    
    if err != nil {
        // Log error but don't crash
        fmt.Printf("Upscaling failed: %v\n", err)
        // Return original as fallback
        return imageData, nil
    }
    
    return result.Data, nil
}
```

### Error Classification

```go
import "github.com/pedro3pv/SweetDesk-core/pkg/errors"

result, err := proc.Process(imageData, options)
if err != nil {
    if errors.IsOutOfMemory(err) {
        // Handle OOM: reduce image size or scale
        log.Println("Out of memory, reducing image size")
    } else if errors.IsModelNotFound(err) {
        // Handle missing model: download or warn
        log.Println("Model not found")
    } else if errors.IsInvalidInput(err) {
        // Handle invalid input: validate format/size
        log.Println("Invalid image format")
    } else {
        // Generic error
        log.Printf("Processing error: %v", err)
    }
}
```

---

## Performance Optimization

### Configuration Tuning

```go
config := processor.Config{
    ModelDir:     "./models",
    MaxWorkers:   4,              // Match CPU cores
    CacheSize:    2000,           // Increase for batch processing
    PixelLimit:   200 * 1024 * 1024, // 200MP max
}

proc, err := processor.New(config)
```

### Batch Processing

```go
func ProcessBatch(images [][]byte) ([][]byte, error) {
    results := make([][]byte, len(images))
    
    for i, imageData := range images {
        result, err := proc.Process(imageData, options)
        if err != nil {
            // Skip failed image
            results[i] = imageData
            continue
        }
        results[i] = result.Data
    }
    
    return results, nil
}
```

### GPU Acceleration

```bash
# Enable NVIDIA GPU
export CUDA_VISIBLE_DEVICES=0

# Or AMD GPU (ROCm)
export HIP_VISIBLE_DEVICES=0

# Application auto-detects and uses GPU
```

---

## Debugging & Monitoring

### Enable Debug Logging

```bash
export LOG_LEVEL=debug
export DEBUG=true
```

### Performance Monitoring

```go
start := time.Now()
result, err := proc.Process(imageData, options)
duration := time.Since(start)

log.Printf("Processing took %v (%d bytes/sec)",
    duration,
    len(imageData)*1000/duration.Milliseconds())
```

### Model Information

```go
info, err := proc.GetModelInfo()
if err == nil {
    log.Printf("Classifier: %s (v%s)", info.Classifier.Name, info.Classifier.Version)
    log.Printf("Upscaler: %s (v%s)", info.Upscaler.Name, info.Upscaler.Version)
    log.Printf("Total size: %.1f MB", info.TotalSizeMB)
}
```

---

## Troubleshooting

### Common Issues

#### 1. "Models not found"

```bash
# Download models
go generate ./...

# Or specify custom path
export SWEETDESK_MODEL_DIR=/custom/path
```

#### 2. "ONNX Runtime error"

```bash
# Check library availability
ldd /usr/lib/libonnxruntime.so  # Linux
otool -L /usr/local/lib/libonnxruntime.dylib  # macOS

# Specify path if needed
export ONNX_LIB_PATH=/custom/path/libonnxruntime.so
```

#### 3. "Out of memory"

```bash
# Reduce concurrency
export MAX_CONCURRENT_TASKS=1

# Or process smaller images
# Split large images before processing
```

#### 4. "Slow processing"

```bash
# Enable GPU acceleration
export CUDA_VISIBLE_DEVICES=0

# Or increase workers
export MAX_WORKERS=8
```

---

## Examples

See `pkg/examples` in the repository for:

- Basic upscaling
- Classification workflow
- Batch processing
- Error handling
- GPU acceleration setup

---

## API Stability

- **Stable**: All methods marked as `Stable` won't change
- **Beta**: Methods may change in minor versions
- **Experimental**: May change anytime

Check documentation for stability status.

---

## Support

- [GitHub Issues](https://github.com/pedro3pv/SweetDesk-core/issues)
- [Documentation](https://github.com/pedro3pv/SweetDesk-core/wiki)
- [Discussions](https://github.com/pedro3pv/SweetDesk-core/discussions)

---

## License

SweetDesk-Core is MIT licensed. See LICENSE file for details.

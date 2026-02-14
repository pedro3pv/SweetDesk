# SweetDesk-Core Integration Guide

## Overview

SweetDesk has been refactored to use **SweetDesk-Core** as its backend image processing engine. SweetDesk-Core is a powerful ONNX Runtime-based processor that provides:

- ✅ **Intelligent Upscaling**: RealCUGAN (anime) and LSDIR (realistic) models
- ✅ **Automatic Classification**: Detects image type and chooses optimal model
- ✅ **Seam Carving**: Content-aware aspect ratio adjustment
- ✅ **Embedded Libraries**: ONNX Runtime bundled in executable
- ✅ **Hardware Acceleration**: CoreML (macOS), CUDA (Linux/Windows)
- ✅ **Tiling Support**: Handles large images efficiently

## Architecture Changes

### Before (Legacy)
```
Frontend (React/Vue)
    ↓
Wails App (app.go)
    ↓
Legacy Services
    ├── Python Bridge (classif + seam carving)
    ├── Local Upscaler (RealESRGAN/RealCUGAN binaries)
    └── Manual binary management
```

### After (SweetDesk-Core)
```
Frontend (React/Vue)
    ↓
Wails App (app.go)
    ↓
SweetDesk-Core Processor (unified)
    ├── ONNX Runtime (embedded)
    ├── Classifier (ONNX model)
    ├── Upscaler (ONNX models)
    └── Seam Carver (native Go)
```

## Key Changes

### 1. Dependencies (`go.mod`)

**Added:**
```go
require github.com/pedro3pv/SweetDesk-core v0.0.0-20260214000000-000000000000
```

**Removed (optional):**
- Python dependencies (no longer needed for classification/seam carving)
- Manual upscaler binary management

### 2. Application (`app.go`)

**Old approach:**
```go
type App struct {
    imageProcessor *services.ImageProcessor
    upscaler       *services.Upscaler        // External binaries
    pythonBridge   *services.PythonBridge    // Classification + seam carving
}
```

**New approach:**
```go
type App struct {
    coreProcessor    *processor.Processor    // All-in-one from SweetDesk-core
    imageProcessor   *services.ImageProcessor // Base64 conversions only
}
```

### 3. Core Bridge Adapter

Created `internal/services/core_bridge.go` for backwards compatibility:

```go
type CoreBridge struct {
    processor *processor.Processor
}

func (cb *CoreBridge) Classify(imageData []byte) (string, error)
func (cb *CoreBridge) Upscale(imageData []byte, imageType string, scale int) ([]byte, error)
func (cb *CoreBridge) UpscaleToResolution(imageData []byte, width, height int, keepAspectRatio bool) ([]byte, error)
func (cb *CoreBridge) AdjustAspectRatio(imageData []byte, targetWidth, targetHeight int) ([]byte, error)
```

## API Changes

### New Methods

#### `ClassifyImage(base64Data string) -> string`
Classifies image as "anime" or "photo"
```go
classification, err := app.ClassifyImage(base64ImageData)
// Returns: "anime" or "photo"
```

#### `UpscaleImage(base64Data, imageType string, scale int) -> string`
Upscales image by specified factor (2x, 4x, etc.)
```go
upscaled, err := app.UpscaleImage(base64ImageData, "photo", 4)
// Returns: base64 encoded upscaled image
```

#### `UpscaleToResolution(base64Data string, targetWidth, targetHeight int) -> string`
Upscales to specific dimensions while maintaining aspect ratio
```go
upscaled, err := app.UpscaleToResolution(base64ImageData, 3840, 2160) // 4K
// Returns: base64 encoded upscaled image
```

#### `AdjustAspectRatio(base64Data string, targetWidth, targetHeight int) -> string`
Adjusts aspect ratio using intelligent seam carving
```go
adjusted, err := app.AdjustAspectRatio(base64ImageData, 1920, 1080) // 16:9
// Returns: base64 encoded adjusted image
```

#### `ProcessImage(base64Data, targetResolution, useSeamCarving bool, targetAspectRatio string) -> string`
Complete workflow: classify → adjust aspect → upscale
```go
result, err := app.ProcessImage(
    base64ImageData, 
    "4K", 
    true, 
    "16:9"
)
// 1. Auto-classify image
// 2. Apply seam carving to 16:9 ratio
// 3. Upscale to 4K (3840x2160)
// Returns: base64 encoded final image
```

### Removed Methods

- `UpscaleImage()` with fixed scale factor (replaced by `UpscaleToResolution`)
- Direct access to Python bridge
- Manual model selection

## Setup Instructions

### 1. Install SweetDesk-Core

```bash
# Option A: Use remote (recommended)
go get github.com/pedro3pv/SweetDesk-core
go mod tidy

# Option B: Use local for development
# Uncomment in go.mod:
# replace github.com/pedro3pv/SweetDesk-core => ../SweetDesk-core
```

### 2. Download Models

Models are downloaded automatically on first run, or manually:

```bash
# Download all classification and upscaling models
DOWNLOAD_ALL_PLATFORMS=false go generate ./...
```

### 3. Environment Variables

```bash
# Optional: Specify custom model directory (default: ./models)
export SWEETDESK_MODEL_DIR="/path/to/models"

# Optional: Specify ONNX Runtime library path
# (uses embedded by default)
export ONNX_LIB_PATH="/path/to/libonnxruntime.so"

# Required for Pixabay integration
export PIXABAY_API_KEY="your_api_key"
```

### 4. Build and Run

```bash
# Development
wails dev

# Production
wails build
```

## Frontend Integration

### Example: Complete Image Processing Workflow

```javascript
// Import Wails runtime
import { Greet, ProcessImage } from './wails.js';

async function processUserImage(imageFile) {
    try {
        // 1. Read image file
        const fileContent = await imageFile.arrayBuffer();
        const base64 = arrayBufferToBase64(fileContent);

        // 2. Call backend with complete workflow
        const result = await ProcessImage(
            base64,
            "4K",        // Target: 4K resolution
            true,        // Enable seam carving
            "16:9"       // Target aspect ratio
        );

        // 3. Display result
        displayProcessedImage(result);
    } catch (error) {
        console.error('Processing failed:', error);
    }
}

function arrayBufferToBase64(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
}
```

## Performance Considerations

### Image Size Recommendations

| Input Size | Processing Time | Output Quality |
|-----------|-----------------|----------------|
| 512×512   | 1-2s            | Excellent      |
| 1024×1024 | 2-4s            | Excellent      |
| 2048×2048 | 4-8s            | Very Good      |
| 4096×4096 | 8-15s           | Good           |

### Tiling

SweetDesk-Core automatically tiles large images:
- Default tile size: 512×512
- Processed in parallel when available
- Results stitched seamlessly

## Troubleshooting

### Issue: "ONNX Runtime not found"

```bash
# Solution: Download embedded libraries
go generate ./...
go build
```

### Issue: Slow processing

```bash
# Check if hardware acceleration is available
export CUDA_VISIBLE_DEVICES=0  # Linux/Windows with NVIDIA GPU
# Or use CoreML on macOS (automatic)
```

### Issue: Out of memory with large images

```bash
# Enable system swap (Linux)
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

## Migration from Legacy

If you have existing code using the old APIs:

### Old API
```go
upscaledData, _ := app.upscaler.UpscaleImage(imageData, UpscaleOptions{
    Model: "realesrgan-photo",
    Scale: 4,
})
```

### New API
```go
upscaled, _ := app.ProcessImage(
    base64.StdEncoding.EncodeToString(imageData),
    "4K",
    false,
    "",
)
```

## File Structure

```
SweetDesk/
├── go.mod                              # Updated with SweetDesk-core
├── go.sum
├── app.go                              # Refactored with core processor
├── main.go                             # Unchanged
├── internal/
│   └── services/
│       ├── core_bridge.go              # NEW: Bridge adapter
│       ├── image_processor.go          # Updated: Minimal changes
│       ├── api_provider.go             # Unchanged: Pixabay
│       └── python_bridge.go            # DEPRECATED: Can be removed
├── frontend/                           # Unchanged
└── models/                             # Downloaded on first run
    ├── classifier/
    │   └── classifier.onnx
    ├── upscaler_anime/
    │   ├── realcugan.onnx
    │   └── ...
    └── upscaler_photo/
        ├── lsdir.onnx
        └── ...
```

## Future Improvements

- [ ] GPU acceleration metrics dashboard
- [ ] Batch processing API for multiple images
- [ ] Custom model support
- [ ] Real-time preview during processing
- [ ] Progress callbacks for long operations
- [ ] Model caching optimization

## Support

For issues or questions:
1. Check [SweetDesk-Core README](https://github.com/pedro3pv/SweetDesk-core)
2. Review example in `pkg/examples`
3. Open an issue on GitHub

## License

Both projects use MIT License - see LICENSE files for details.

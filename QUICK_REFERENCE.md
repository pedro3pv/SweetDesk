# 📓 SweetDesk - Quick Reference

## Initialize CoreBridge

```go
import "github.com/Molasses-Co/SweetDesk/internal/services"

ctx := context.Background()
bridge, err := services.NewCoreBridge(ctx)
if err != nil {
    log.Fatal(err)
}
defer bridge.Close()
```

## Classify Image

```go
// From URL
data, _ := bridge.DownloadImage("https://example.com/image.png")
imageType, confidence, _ := bridge.ClassifyImage(data)
fmt.Printf("%s (%.2f%%)\n", imageType, confidence*100)
```

## Upscale Image

```go
// 4x upscaling with PNG output
upscaled, _ := bridge.UpscaleImage(imageData, 4, "png")

// Save to disk
path, _ := bridge.SaveImage(upscaled, "./output", "upscaled.png")
```

## Process with Custom Options

```go
// Advanced processing
processed, _ := bridge.ProcessImage(
    imageData,
    2048,           // target width
    2048,           // target height
    "jpeg",         // output format
    "maintain",     // aspect ratio
)
```

## Image Processor Utilities

```go
processor := services.NewImageProcessor(ctx)

// Validate
err := processor.ValidateImage(imageData)

// Get info
info, _ := processor.GetImageInfo(imageData)
fmt.Printf("%dx%d\n", info["width"], info["height"])

// Convert format
jpeg, _ := processor.ConvertFormat(imageData, "jpeg")

// Base64
b64 := processor.ConvertToBase64(imageData)
original, _ := processor.ConvertFromBase64(b64)

// Save
path, _ := processor.SaveToFile(imageData, "./output", "result.png")
```

## Integration with Wails

```go
// main.go
app := wails.CreateApp(&options.App{
    Title: "SweetDesk",
    OnStartup: func(ctx context.Context) {
        bridge, _ := services.NewCoreBridge(ctx)
        // Bridge ready
    },
    Bind: []interface{}{
        bridge,
    },
})
```

## Frontend TypeScript

```typescript
import { CoreBridge } from './wails/go/services/CoreBridge';

// Classify
const [type, confidence] = await CoreBridge.ClassifyImage(bytes);

// Upscale
const upscaled = await CoreBridge.UpscaleImage(bytes, 4, 'png');

// Save
const path = await CoreBridge.SaveImage(upscaled, './out', 'result.png');

// Get processor info
const info = await CoreBridge.GetInfo();
console.log(info.engine); // "SweetDesk-core"
```

## Error Handling

```go
if err != nil {
    fmt.Printf("Original error: %v\n", errors.Unwrap(err))
    fmt.Printf("Wrapped error: %v\n", err)
}
```

## Build Commands

```bash
# Download dependencies
go mod download
go mod tidy

# Run tests
go test -v ./internal/services
go test -v -run TestNewCoreBridge ./internal/services

# Build desktop app
wails build
wails build -production

# Verify compilation
go build -v ./cmd/desktop
```

## Configuration

```bash
# Set custom model directory
export ONNX_LIB_PATH=/opt/onnxruntime/lib

# Enable debug logging
export SWEETDESK_DEBUG=true
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| "undefined reference" | Run `go mod tidy` |
| "models not found" | Verify `models/` directory structure |
| "processor not initialized" | Check `NewCoreBridge` error |
| "out of memory" | Reduce `tileSize` from 512 to 256 |
| "ONNX not found" | Install ONNX runtime library |

## File Locations

```
SweetDesk/
├── internal/services/
│   ├── core_bridge.go          # Main integration
│   ├── core_bridge_test.go     # Tests
│   ├── upscaler.go             # Upscaling engine
│   └── ...
├── models/
│   ├── realcugan/
│   │   └── realcugan-pro.onnx
│   └── lsdir/
│       └── 4xLSDIR.onnx
├── go.mod                      # Dependencies
├── FIXES_AND_INTEGRATION.md    # Detailed docs
├── INTEGRATION_SETUP.md        # Setup guide
└── QUICK_REFERENCE.md          # This file
```

## Useful Commands

```bash
# Check Go version
go version

# List dependencies
go list -m all | grep SweetDesk

# Verify module
go mod verify

# Format code
go fmt ./...

# Lint code
golangci-lint run ./...

# Generate docs
go doc github.com/Molasses-Co/SweetDesk/internal/services
```

## Performance Tips

- Use PNG for lossless processing
- Use JPEG for distribution (smaller file size)
- Process tiles for images > 2048x2048
- Use context cancellation for long operations
- Monitor memory with `top` or `htop` during batch processing

## Documentation

- **FIXES_AND_INTEGRATION.md** - All 10 fixes explained
- **INTEGRATION_SETUP.md** - Step-by-step setup
- **QUICK_REFERENCE.md** - This file
- **core_bridge_test.go** - Working examples

## Support

- GitHub Issues: [Molasses-Co/SweetDesk](https://github.com/Molasses-Co/SweetDesk/issues)
- Core Repository: [pedro3pv/SweetDesk-core](https://github.com/pedro3pv/SweetDesk-core)

---

**Last Updated:** 2026-02-14 | **Version:** 0.1.0


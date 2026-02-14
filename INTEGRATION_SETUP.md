# 🚀 SweetDesk-Core Integration Setup Guide

## Prerequisites

- Go 1.21+
- Git
- 4GB+ RAM
- ONNX Runtime (optional but recommended for performance)

---

## 🔜 Step 1: Clone and Setup

### Clone Repositories
```bash
# Clone main SweetDesk repository
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# For local development: clone SweetDesk-core locally
git clone https://github.com/pedro3pv/SweetDesk-core.git ../SweetDesk-core
```

### Update Dependencies
```bash
# Download all dependencies
go mod download

# Clean up unused dependencies
go mod tidy

# Verify dependencies
go mod verify
```

---

## 🔗 Step 2: Configure Integration

### Option A: Using GitHub Dependency (Production)
```bash
# Dependencies are automatically resolved from GitHub
go get -u github.com/pedro3pv/SweetDesk-core@latest
```

### Option B: Using Local Development Version

1. Uncomment the replace directive in `go.mod`:
```go
replace github.com/pedro3pv/SweetDesk-core => ../SweetDesk-core
```

2. Verify the path:
```bash
ls ../SweetDesk-core/go.mod  # Should exist
```

3. Update modules:
```bash
go mod tidy
```

---

## 📁 Step 3: Setup Models

### Create Models Directory Structure
```bash
mkdir -p models/realcugan
mkdir -p models/lsdir
```

### Download Models

**RealCUGAN Model** (for anime/artwork)
```bash
# Download from official source
cd models/realcugan
wget https://github.com/bilibili/Real-CUGAN/releases/download/v3/realcugan-pro.onnx
cd ../..
```

**LSDIR Model** (for photos)
```bash
# Download from official source
cd models/lsdir
wget https://github.com/chaofengc/Face-Restoration-Benchmark/releases/download/v0.0.1/4xLSDIR.onnx
cd ../..
```

### Verify Models
```bash
ls -lh models/realcugan/realcugan-pro.onnx
ls -lh models/lsdir/4xLSDIR.onnx
```

---

## 📋 Step 4: Build and Test

### Run Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./internal/services

# Run specific test
go test -v -run TestCoreBridge ./internal/services
```

### Build Application
```bash
# Install Wails CLI if not already installed
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build development version
wails build

# Build production version
wails build -production
```

---

## ✅ Step 5: Verification Checklist

### Compilation Check
```bash
# Verify compilation without errors
go build -v ./cmd/desktop

# Expected output: No errors, binary created
```

### Dependency Check
```bash
# List all dependencies
go list -m all | grep -E "(SweetDesk-core|wails|onnx|gorgonia)"

# Expected output:
# github.com/pedro3pv/SweetDesk-core v0.1.0
# github.com/wailsapp/wails/v2 v2.11.0
# github.com/owulveryck/onnx-go v0.4.0
# gorgonia.org/tensor v0.9.24
```

### CoreBridge Test
```bash
# Run bridge initialization test
go test -v -run TestNewCoreBridge ./internal/services

# Expected output:
# ✅ SweetDesk-core bridge initialized
```

### Image Processing Test
```bash
# Run image processor tests
go test -v -run TestImageProcessor ./internal/services

# Expected output:
# All tests passing
```

### Model Loading Test
```go
// Quick test file: test_models.go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

func main() {
	config := processor.Config{
		ModelDir: "./models",
	}

	proc, err := processor.New(config)
	if err != nil {
		log.Fatalf("Failed to init processor: %v", err)
	}
	defer proc.Close()

	info, _ := proc.GetInfo()
	fmt.Printf("✅ Processor ready: %v\n", info)
}
```

Run:
```bash
go run test_models.go
```

---

## 🖌 Integration Testing

### Test Case 1: Anime Image Upscaling
```bash
# Place test anime image
cp sample_anime.png test_input.png

# Create test program
cat > test_upscale.go << 'EOF'
package main

import (
	"fmt"
	"os"
	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

func main() {
	data, _ := os.ReadFile("test_input.png")
	
	config := processor.Config{ModelDir: "./models"}
	proc, _ := processor.New(config)
	defer proc.Close()
	
	opts := processor.ProcessOptions{TargetScale: 4, OutputFormat: "png"}
	result, err := proc.Upscale(data, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	os.WriteFile("test_output.png", result.Data, 0644)
	fmt.Println("✅ Upscaling successful")
}
EOF

go run test_upscale.go
ls -lh test_output.png
```

### Test Case 2: Format Conversion
```bash
# Test PNG to JPEG conversion
go test -v -run TestImageProcessorConvertFormat ./internal/services
```

### Test Case 3: Classification
```go
// Create test_classify.go
package main

import (
	"fmt"
	"os"
	"github.com/pedro3pv/SweetDesk-core/pkg/processor"
)

func main() {
	data, _ := os.ReadFile("test_input.png")
	
	config := processor.Config{ModelDir: "./models"}
	proc, _ := processor.New(config)
	defer proc.Close()
	
	result, _ := proc.Classify(data)
	fmt.Printf("Image Type: %s (%.2f%% confidence)\n", result.Type, result.Confidence*100)
}
EOF

go run test_classify.go
```

---

## 🗐 Environment Variables

### Optional Configuration
```bash
# Specify custom ONNX runtime library
export ONNX_LIB_PATH=/usr/lib/libonnxruntime.so

# Enable debug logging
export SWEETDESK_DEBUG=true

# Set model directory
export SWEETDESK_MODEL_DIR=/opt/models
```

---

## 🏗 Build Troubleshooting

### Issue: "undefined reference to processor"
**Solution:** Verify go.mod dependency:
```bash
grep "SweetDesk-core" go.mod
# Should output: github.com/pedro3pv/SweetDesk-core v0.1.0

go get -u github.com/pedro3pv/SweetDesk-core@latest
go mod tidy
```

### Issue: "models not found"
**Solution:** Verify models directory:
```bash
ls -R models/
# Should show:
# models/realcugan/realcugan-pro.onnx
# models/lsdir/4xLSDIR.onnx
```

### Issue: "ONNX runtime not found"
**Solution:** Install ONNX runtime:

**Ubuntu/Debian:**
```bash
apt-get install libonnxruntime1
```

**macOS:**
```bash
brew install onnxruntime
```

**Windows:**
Download from [ONNX Runtime Releases](https://github.com/microsoft/onnxruntime/releases)

### Issue: "Out of memory during upscaling"
**Solution:** Adjust tile size in `internal/services/upscaler.go`:
```go
u := &Upscaler{
    tileSize: 256,  // Reduced from 512
    modelScale: 4,
}
```

---

## 🐞 Performance Testing

### Benchmark Small Image
```bash
# Test with 512x512 image
time go run test_upscale.go

# Expected: 5-15 seconds depending on hardware
```

### Benchmark Large Image
```bash
# Test with 2048x2048 image
# Uses tiling internally
time go run test_upscale.go

# Expected: 30-60 seconds
```

### Memory Usage
```bash
# Monitor during processing
mem_before=$(free -h | awk 'NR==2{print $3}')
go run test_upscale.go
mem_after=$(free -h | awk 'NR==2{print $3}')
echo "Memory delta: $mem_before -> $mem_after"
```

---

## 🤐 Production Deployment

### Final Verification
```bash
# Build release binary
wails build -production -o sweetdesk

# Verify binary runs
./sweetdesk --help

# Test with sample image
./sweetdesk -input sample.png -output output.png
```

### Package Distribution
```bash
# Include models in distribution
mkdir -p distribution/models
cp -r models/* distribution/models/
cp sweetdesk distribution/

# Create archive
tar -czf sweetdesk-production.tar.gz distribution/
```

---

## 📚 Additional Resources

- [SweetDesk-Core Documentation](https://github.com/pedro3pv/SweetDesk-core/wiki)
- [ONNX Runtime Setup](https://onnxruntime.ai/docs/install/)
- [Wails Documentation](https://wails.io/docs/)
- [Go Modules Guide](https://go.dev/blog/using-go-modules)

---

## ✅ Verification Checklist (Final)

- [ ] Dependencies installed: `go mod download && go mod tidy`
- [ ] Models downloaded and verified
- [ ] Tests passing: `go test ./...`
- [ ] Build successful: `wails build`
- [ ] Runtime test working
- [ ] Image upscaling functional
- [ ] Format conversion working
- [ ] Classification model accessible
- [ ] Frontend can call backend methods
- [ ] Production build created

**Once all items checked, SweetDesk is ready for deployment! 🚀**


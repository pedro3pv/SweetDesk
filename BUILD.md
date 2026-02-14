# Build Instructions - SweetDesk with SweetDesk-Core

## Prerequisites

- Go 1.24.0 or later
- Node.js 18+ (for frontend)
- Wails v2 framework
- Git

### Platform-Specific

**macOS:**
```bash
# Install Xcode Command Line Tools
xcode-select --install

# Or install from App Store
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install build-essential libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install @development-tools pkg-config gtk3-devel webkit2gtk3-devel
```

**Windows:**
- Visual Studio Build Tools or Visual Studio Community
- MinGW-w64
- WebView2 Runtime

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk
```

### 2. Install Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 3. Install Dependencies

```bash
# Download Go dependencies
go mod download
go mod tidy

# Install frontend dependencies
cd frontend
npm install
cd ..
```

### 4. Download Models and Libraries

```bash
# This downloads ONNX Runtime and AI models
go generate ./...
```

## Development Build

### Local Development with Hot Reload

```bash
# Terminal 1: Start Wails dev server
wails dev

# This will:
# - Start the backend (go)
# - Start the frontend dev server (npm)
# - Open browser window
# - Hot reload on file changes
```

### Troubleshooting Dev Build

**Issue: Port 3000 already in use**
```bash
wails dev --port 3001
```

**Issue: Models not found**
```bash
# Re-download models
go generate ./...
wails dev
```

## Production Build

### Build for Current Platform

```bash
# Builds both frontend and backend
wails build

# Output:
# macOS: build/bin/SweetDesk.app
# Linux: build/bin/SweetDesk
# Windows: build/bin/SweetDesk.exe
```

### Build for Specific Platform

```bash
# macOS (Intel)
wails build -os darwin -arch amd64

# macOS (Apple Silicon)
wails build -os darwin -arch arm64

# Linux (x86-64)
wails build -os linux -arch amd64

# Linux (ARM64)
wails build -os linux -arch arm64

# Windows (x86-64)
wails build -os windows -arch amd64
```

### Build with All Platforms (SweetDesk-Core Models)

```bash
# Download ONNX Runtime for all platforms
DOWNLOAD_ALL_PLATFORMS=true go generate ./...

# Build for each platform
make release  # If Makefile exists
```

## Advanced Builds

### Debug Build with Symbols

```bash
wails build -debug
```

### Development Optimized

```bash
wails build -dev
```

### With Custom Frontend Build

```bash
# Build frontend first
cd frontend
npm run build
cd ..

# Then build app
wails build
```

## Environment Configuration

Create `.env` file in project root:

```bash
# Pixabay API
PIXABAY_API_KEY=your_api_key_here

# Model directory (default: ./models)
SWEETDESK_MODEL_DIR=./models

# ONNX Runtime (optional, uses embedded by default)
ONNX_LIB_PATH=./onnx-lib/libonnxruntime.so
```

## Local Development (SweetDesk-Core)

If developing SweetDesk-Core alongside SweetDesk:

### 1. Use Local Replace in go.mod

```bash
cd SweetDesk

# Edit go.mod
# Add or uncomment:
# replace github.com/pedro3pv/SweetDesk-core => ../SweetDesk-core

go mod tidy
```

### 2. Ensure SweetDesk-Core is built

```bash
cd ../SweetDesk-core
go generate ./...
```

### 3. Run development build

```bash
cd ../SweetDesk
wails dev
```

## Build Optimization

### Strip Binary (Linux/macOS)

```bash
wails build -skipbindings  # Skip Go bindings regeneration
```

### Compress Output

```bash
# After build
cd build/bin/

# Linux
strip SweetDesk
upx SweetDesk  # If upx installed

# macOS
strip SweetDesk.app/Contents/MacOS/SweetDesk
```

### Check Binary Size

```bash
# Before optimization
ls -lh build/bin/SweetDesk*

# After optimization
ls -lh build/bin/SweetDesk*
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: macos-latest
            os_name: darwin
            arch: amd64
          - os: ubuntu-latest
            os_name: linux
            arch: amd64
          - os: windows-latest
            os_name: windows
            arch: amd64

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          go mod download
          cd frontend && npm install
      - name: Download models
        run: go generate ./...
      - name: Build
        run: |
          go install github.com/wailsapp/wails/v2/cmd/wails@latest
          wails build -os ${{ matrix.os_name }} -arch ${{ matrix.arch }}
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: SweetDesk-${{ matrix.os_name }}-${{ matrix.arch }}
          path: build/bin/
```

## Verification

After build:

```bash
# Test binary exists
file build/bin/SweetDesk*

# macOS: Should show "Mach-O 64-bit executable"
# Linux: Should show "ELF 64-bit executable"
# Windows: Should show "PE32+ executable"
```

## Cleanup

```bash
# Remove build artifacts
rm -rf build/

# Clear Go cache
go clean -cache

# Remove downloaded models (careful!)
rm -rf models/
```

## Troubleshooting

### "command not found: wails"

```bash
# Ensure $GOPATH/bin is in PATH
echo $PATH | grep $GOPATH/bin

# If not, add to ~/.bashrc or ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
```

### "ONNX Runtime not found"

```bash
# Download models
go generate ./...
```

### "Frontend build failed"

```bash
cd frontend
npm cache clean --force
rm -rf node_modules
npm install
cd ..
wails build
```

### "Unable to find SDK"

**macOS:**
```bash
sudo xcode-select --reset
sudo xcode-select --install
```

**Windows:**
- Install Visual Studio Build Tools
- Ensure WebView2 Runtime is installed

## Performance Tips

### Speed up builds

```bash
# Use -skipbindings to skip regeneration
wails build -skipbindings

# Use -dev for faster builds (no optimizations)
wails dev
```

### Parallel compilation

```bash
# Increase parallelism
export GOMAXPROCS=8
wails build
```

## Next Steps

After successful build:

1. **Run the app**: Open the built binary
2. **Test locally**: Check all features work
3. **Set API keys**: Configure Pixabay API
4. **Download models**: First run will auto-download (or run `go generate ./...')
5. **Process images**: Test with sample images

For more information, see:
- [INTEGRATION.md](./INTEGRATION.md) - SweetDesk-Core integration details
- [README.md](./README.md) - Project overview
- [Development Guide](./DEVELOPMENT.md) - Code structure

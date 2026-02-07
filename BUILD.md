# SweetDesk Build Configuration

## Overview

This document describes the build system for SweetDesk, a cross-platform wallpaper processing application.

## Prerequisites

- Go 1.22+
- Node.js 18+
- Wails CLI v2.11+
- Python 3.8+ (for AI processing)

## Build Process

### 1. Download AI Binaries

```bash
./scripts/download-binaries.sh
```

This downloads:
- Real-ESRGAN NCNN Vulkan binaries (macOS, Linux, Windows)
- RealCUGAN NCNN Vulkan binaries (macOS, Linux, Windows)
- AI model files (.param and .bin files)

### 2. Install Dependencies

```bash
# Frontend
cd frontend && npm install

# Python
pip install -r python/requirements.txt
```

### 3. Build for Specific Platform

```bash
# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64

# Windows
wails build -platform windows/amd64
```

### 4. Build All Platforms

```bash
./scripts/build-all-platforms.sh
```

## Directory Structure

```
SweetDesk/
├── binaries/                    # AI upscaler binaries
│   ├── darwin/                  # macOS binaries
│   │   ├── realesrgan-ncnn-vulkan
│   │   ├── realcugan-ncnn-vulkan
│   │   └── models/              # AI models (.param, .bin)
│   ├── linux/                   # Linux binaries
│   └── windows/                 # Windows binaries
├── python/                      # Python scripts for AI
│   ├── classify_image.py        # Image classification
│   ├── seam_carving.py          # Content-aware resize
│   └── requirements.txt
├── internal/services/           # Go services
│   ├── api_provider.go          # Pixabay API integration
│   ├── image_processor.go       # Image processing
│   ├── python_bridge.go         # Python subprocess handler
│   └── upscaler.go              # AI upscaling
├── frontend/                    # Next.js frontend
│   └── src/
│       ├── app/
│       └── components/
└── scripts/                     # Build scripts
    ├── download-binaries.sh
    └── build-all-platforms.sh
```

## Cross-Platform Notes

### macOS
- Uses Vulkan or Metal for GPU acceleration
- Universal binary supports both Intel and Apple Silicon
- Binaries signed for distribution

### Linux
- Requires Vulkan support
- Distribution via AppImage or flatpak

### Windows
- Requires DirectX 12 or Vulkan
- Distributed as executable with installer

## In-Memory Processing

All image processing happens in memory using byte arrays:
- Images loaded from API or file → byte[]
- Processing pipeline operates on byte[] → byte[]
- Only writes to disk when user clicks "Save"

This ensures:
- Fast processing
- No temporary file cleanup needed
- Privacy (images never touch disk until saved)

## AI Model Sources

Production builds should download models from official sources:

- **Real-ESRGAN**: https://github.com/xinntao/Real-ESRGAN/releases
- **RealCUGAN**: https://github.com/nihui/realcugan-ncnn-vulkan/releases

Current implementation uses placeholders for development.

## Environment Variables

```bash
# Required for Pixabay integration
export PIXABAY_API_KEY="your_api_key_here"

# Optional: Enable debug mode
export SWEETDESK_DEBUG=1
```

## Testing Builds

```bash
# Run in development mode
wails dev

# Test specific platform build
./build/bin/SweetDesk  # or SweetDesk.exe on Windows
```

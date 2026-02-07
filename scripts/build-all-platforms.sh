#!/bin/bash
# Build script for SweetDesk - All platforms

set -e

echo "🍬 SweetDesk - Building for all platforms"
echo "=========================================="

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "❌ Wails CLI not found. Installing..."
    ./scripts/install-wails-cli.sh
fi

# Download binaries and models
echo ""
echo "📥 Downloading AI binaries..."
./scripts/download-binaries.sh

# Install frontend dependencies
echo ""
echo "📦 Installing frontend dependencies..."
cd frontend
npm install
cd ..

# Build for each platform
echo ""
echo "🔨 Building for all platforms..."
echo ""

# macOS
echo "→ Building for macOS..."
wails build -platform darwin/universal -o SweetDesk

# Linux
echo "→ Building for Linux..."
wails build -platform linux/amd64 -o SweetDesk

# Windows
echo "→ Building for Windows..."
wails build -platform windows/amd64 -o SweetDesk.exe

echo ""
echo "✅ Build complete!"
echo ""
echo "Outputs:"
echo "  - macOS: build/bin/SweetDesk.app"
echo "  - Linux: build/bin/SweetDesk"
echo "  - Windows: build/bin/SweetDesk.exe"
echo ""

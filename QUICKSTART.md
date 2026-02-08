# Quick Start Guide - SweetDesk

Get up and running with SweetDesk in 5 minutes!

## Step 1: Prerequisites (2 min)

Install the required tools:

```bash
# Check if you have the required tools
go version        # Should be 1.22+
node --version    # Should be 18+
python3 --version # Should be 3.8+
```

If missing, install them:
- **Go**: https://golang.org/dl/
- **Node.js**: https://nodejs.org/
- **Python**: https://python.org/downloads/

## Step 2: Clone and Setup (2 min)

```bash
# Clone the repository
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# Create environment file
cp .env.example .env

# Get a FREE Pixabay API key at https://pixabay.com/api/docs/
# Edit .env and add your key:
# PIXABAY_API_KEY=your_key_here
```

## Step 3: Install Wails (1 min)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

## Step 4: Run Development Mode (instant)

```bash
wails dev
```

That's it! The app should open and you can:

1. 🔍 **Search** for wallpapers using Pixabay
2. 📁 **Upload** your own images
3. 🚀 **Process** to 4K with AI upscaling
4. 💾 **Save** the result

## What You Can Do

### Search and Process

1. Click the "🔍 Search" tab
2. Type a query (e.g., "mountains", "ocean", "space")
3. Click an image to select it
4. Choose resolution (4K, 5K, or 8K)
5. Click "🚀 Process Image"
6. Wait for processing (30s - 2min)
7. Click "💾 Save Image"

### Upload and Process

1. Click the "📁 Upload" tab
2. Drag & drop an image or click to browse
3. Choose processing options
4. Click "🚀 Process Image"
5. Save the result

## Processing Options

### Resolution
- **4K**: 3840 × 2160 (most common, ~45s)
- **5K**: 5120 × 2880 (high quality, ~75s)
- **8K**: 7680 × 4320 (ultra quality, ~120s)

### Aspect Ratio Adjustment
- **Fast Crop**: Quick center crop (instant)
- **Content-Aware**: Smart seam carving (adds ~10s)

## Troubleshooting

### App doesn't start
```bash
# Reinstall dependencies
cd frontend && npm install && cd ..
go mod download
```

### "Python not found"
```bash
# macOS/Linux
brew install python3
# or
sudo apt install python3

# Windows: Download from python.org
```

### "API key not configured"
- Make sure you created `.env` from `.env.example`
- Add your Pixabay API key (get it free at https://pixabay.com/api/docs/)

### Processing fails
- Ensure Python dependencies are installed:
  ```bash
  pip install -r python/requirements.txt
  ```

## Next Steps

- ✅ Read [DEVELOPMENT.md](DEVELOPMENT.md) for development guide
- ✅ Check [API.md](API.md) for API documentation
- ✅ See [BUILD.md](BUILD.md) for production builds

## Performance Tips

### For Faster Processing

1. **Use smaller resolutions first** (4K instead of 8K)
2. **Disable content-aware resize** for speed
3. **Close other applications** to free RAM
4. **Use SSD** if available (faster I/O)

### Expected Processing Times

| Hardware | 4K Upscale | 8K Upscale |
|----------|------------|------------|
| M1/M2/M3 Mac | ~30-45s | ~90-120s |
| Intel i7+ | ~50-70s | ~120-180s |
| AMD Ryzen 7+ | ~45-60s | ~100-150s |

## Features Showcase

### 🎨 Anime Detection
The app automatically detects anime-style images and uses RealCUGAN for better line preservation.

### 📷 Photo Enhancement
Real-ESRGAN models enhance photo details while maintaining natural look.

### 🖼️ Smart Cropping
Content-aware seam carving preserves important image content when adjusting aspect ratio.

### 🌓 Dark Mode
Automatically follows your system dark/light mode preference.

## Building for Production

Want to share the app? Build it:

```bash
# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64

# Windows
wails build -platform windows/amd64
```

Output will be in `build/bin/`

## Get Help

- 🐛 **Bug Report**: https://github.com/Molasses-Co/SweetDesk/issues
- 💬 **Questions**: https://github.com/Molasses-Co/SweetDesk/discussions
- 📖 **Full Docs**: See README.md

---

**Happy wallpaper processing!** 🍬✨

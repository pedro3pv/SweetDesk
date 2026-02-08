# SweetDesk Development Guide

## Getting Started

### Prerequisites

1. **Go** (1.22+)
   ```bash
   # macOS
   brew install go
   
   # Linux
   sudo apt install golang-go
   
   # Windows
   # Download from https://golang.org/dl/
   ```

2. **Node.js** (18+)
   ```bash
   # macOS
   brew install node
   
   # Linux
   sudo apt install nodejs npm
   
   # Windows
   # Download from https://nodejs.org/
   ```

3. **Wails CLI** (v2.11+)
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **Python** (3.8+)
   ```bash
   # Usually pre-installed on macOS/Linux
   # Windows: Download from https://python.org/
   
   # Install dependencies
   pip install -r python/requirements.txt
   ```

### Initial Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Molasses-Co/SweetDesk.git
   cd SweetDesk
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env and add your PIXABAY_API_KEY
   ```

3. **Install dependencies**
   ```bash
   # Frontend
   cd frontend && npm install && cd ..
   
   # Go modules
   go mod download
   
   # Python
   pip install -r python/requirements.txt
   ```

4. **Download AI binaries** (optional for development)
   ```bash
   ./scripts/download-binaries.sh
   ```

### Development Mode

Run the app in development mode with hot reload:

```bash
wails dev
```

This will:
- Start the Go backend
- Start the Next.js dev server
- Open the app window
- Enable hot reload for both frontend and backend

### Project Structure

```
SweetDesk/
├── main.go                      # Wails entry point
├── app.go                       # Main app logic with Wails bindings
├── internal/
│   └── services/                # Business logic services
│       ├── api_provider.go      # API integrations (Pixabay, etc.)
│       ├── image_processor.go   # Image manipulation
│       ├── python_bridge.go     # Python subprocess management
│       └── upscaler.go          # AI upscaling logic
├── frontend/
│   ├── src/
│   │   ├── app/                 # Next.js app directory
│   │   └── components/          # React components
│   ├── wailsjs/                 # Generated Wails bindings
│   ├── package.json
│   └── next.config.ts
├── python/
│   ├── classify_image.py        # Image classification
│   ├── seam_carving.py          # Content-aware resize
│   └── requirements.txt
├── binaries/                    # AI upscaler binaries
│   ├── darwin/
│   ├── linux/
│   └── windows/
└── scripts/                     # Build and utility scripts
```

## Development Workflow

### Adding a New Feature

1. **Backend (Go)**
   - Add method to `app.go`
   - Methods will be auto-exposed to frontend via Wails
   - Example:
   ```go
   func (a *App) MyNewFeature(param string) (string, error) {
       // Your logic here
       return "result", nil
   }
   ```

2. **Frontend (React/TypeScript)**
   - Use Wails bindings:
   ```typescript
   const result = await window.go.main.App.MyNewFeature("param");
   ```
   - Create components in `frontend/src/components/`
   - Add pages in `frontend/src/app/`

3. **Testing**
   - Run in dev mode: `wails dev`
   - Test functionality
   - Check console for errors

### Adding a New API Provider

1. Create a new provider in `internal/services/api_provider.go`:
   ```go
   type UnsplashProvider struct {
       apiKey string
       client *http.Client
       ctx    context.Context
   }
   
   func (p *UnsplashProvider) Search(...) ([]ImageResult, error) {
       // Implementation
   }
   ```

2. Update `app.go` to support the new provider

3. Update frontend to show provider selection

### Working with Python Scripts

The Python scripts are called via subprocess. To debug:

1. Test scripts directly:
   ```bash
   python python/classify_image.py path/to/image.jpg
   ```

2. Check output format (should be valid JSON)

3. Ensure scripts are executable:
   ```bash
   chmod +x python/*.py
   ```

## Building for Production

### Single Platform

```bash
# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64

# Windows (cross-compile from Linux/macOS)
wails build -platform windows/amd64
```

### All Platforms

```bash
./scripts/build-all-platforms.sh
```

Output locations:
- macOS: `build/bin/SweetDesk.app`
- Linux: `build/bin/SweetDesk`
- Windows: `build/bin/SweetDesk.exe`

## Debugging

### Go Backend

Add debug prints:
```go
fmt.Printf("Debug: %v\n", variable)
```

Or use the logger:
```go
import "log"
log.Printf("Debug: %v", variable)
```

### Frontend

Use browser dev tools:
- In dev mode, press `Cmd+Option+I` (macOS) or `Ctrl+Shift+I` (Windows/Linux)
- Check Console tab for errors
- Use `console.log()` for debugging

### Python Scripts

Add debug output to stderr:
```python
import sys
print(f"Debug: {variable}", file=sys.stderr)
```

The Go backend will capture and log stderr output.

## Testing

### Manual Testing

1. Start dev mode: `wails dev`
2. Test each feature:
   - Search images
   - Upload image
   - Process image
   - Save result

### Integration Testing

Test the full pipeline:
1. Search for an image
2. Select and download it
3. Process with different options
4. Verify output quality

## Common Issues

### "Python not found"

Install Python 3:
```bash
# macOS
brew install python3

# Linux
sudo apt install python3

# Verify
python3 --version
```

### "Wails command not found"

Install Wails CLI:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Ensure $GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Frontend build errors

Clean and reinstall:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..
```

### "Model not found" errors

Download the AI binaries:
```bash
./scripts/download-binaries.sh
```

Or download manually from:
- https://github.com/xinntao/Real-ESRGAN/releases
- https://github.com/nihui/realcugan-ncnn-vulkan/releases

## Performance Optimization

### Backend

1. **Use goroutines for concurrent processing**
   ```go
   go processImage(image)
   ```

2. **Implement caching** for API results

3. **Pool connections** for external services

### Frontend

1. **Use React.memo** for expensive components
2. **Lazy load** images
3. **Debounce** search input

### Python

1. **Use NumPy** for fast array operations
2. **Consider Cython** for critical paths
3. **Add GPU support** for intensive processing

## Code Style

### Go

Follow standard Go conventions:
- Use `gofmt` for formatting
- Use meaningful variable names
- Add comments for exported functions

### TypeScript/React

Follow React best practices:
- Use functional components
- Use TypeScript for type safety
- Follow Next.js conventions

### Python

Follow PEP 8:
- Use 4 spaces for indentation
- Maximum line length: 88 (Black formatter)
- Use type hints where appropriate

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

### PR Checklist

- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] No console errors in dev mode

## Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [Next.js Documentation](https://nextjs.org/docs)
- [Go Documentation](https://golang.org/doc/)
- [Real-ESRGAN](https://github.com/xinntao/Real-ESRGAN)
- [Pixabay API](https://pixabay.com/api/docs/)

## Support

- GitHub Issues: https://github.com/Molasses-Co/SweetDesk/issues
- Discussions: https://github.com/Molasses-Co/SweetDesk/discussions

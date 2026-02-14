# Development Guidelines for SweetDesk

## Table of Contents

1. [Project Structure](#project-structure)
2. [Setup Environment](#setup-environment)
3. [Build & Run](#build--run)
4. [Code Guidelines](#code-guidelines)
5. [Testing Strategy](#testing-strategy)
6. [Debugging](#debugging)
7. [Common Tasks](#common-tasks)
8. [CI/CD Pipeline](#cicd-pipeline)

---

## Project Structure

```
SweetDesk/
├─ frontend/              # React/Vue UI (Wails frontend)
│  ├─ src/
│  │  ├─ components/   # React/Vue components
│  │  ├─ pages/       # Page components
│  │  ├─ services/    # API calls (Wails bindings)
│  │  ├─ hooks/       # Custom React hooks
│  │  ├─ styles/      # CSS/SCSS
│  │  ├─ types/       # TypeScript types
│  │  └─ App.tsx      # Main component
│  ├─ package.json
│  ├─ tsconfig.json
│  └─ vite.config.js
├─ cmd/app/              # Wails app entrypoint
│  └─ main.go            # Bootstrap & configuration
├─ internal/              # Private Go packages
│  ├─ app.go              # App struct & lifecycle
│  ├─ handlers/           # API handlers
│  │  ├─ image.go        # Image processing endpoints
│  │  ├─ upscale.go      # Upscaling endpoints
│  │  ├─ pixabay.go      # Pixabay API endpoints
│  │  └─ process.go      # Full pipeline endpoints
│  ├─ services/          # Business logic
│  │  ├─ image.go        # Image utilities
│  │  ├─ cache.go        # Caching layer
│  │  └─ pixabay.go      # Pixabay client
│  ├─ models/             # Data models
│  │  ├─ image.go
│  │  ├─ config.go
│  │  └─ result.go
│  └─ utils/              # Utilities
│     ├─ logger.go
│     ├─ errors.go
│     └─ constants.go
├─ pkg/                   # Public packages (if any)
│  └─ examples/            # Example implementations
├─ models/                # ML model files
│  ├─ classifier/
│  └─ upscaler/
├─ go.mod
├─ go.sum
├─ wails.json            # Wails configuration
├─ .env.example          # Environment variables template
├─ .github/
│  ├─ workflows/
│  │  ├─ test.yml       # Automated tests
│  │  ├─ build.yml      # Build binaries
│  │  └─ release.yml    # Release artifacts
│  └─ ISSUE_TEMPLATE/
├─ README.md
├─ MIGRATION_GUIDE.md     # From legacy to SweetDesk-core
└─ DEVELOPMENT.md         # This file
```

---

## Setup Environment

### Prerequisites

- **Go 1.21+** - Backend language
- **Node.js 18+** - Frontend build
- **npm or yarn** - Package management
- **Wails v2.11+** - Desktop framework
- **Git** - Version control

### macOS Setup

```bash
# Install Homebrew (if not installed)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install go node wails

# Verify installation
go version  # Go 1.21+
node -v    # Node 18+
wails version
```

### Linux Setup (Ubuntu/Debian)

```bash
# Install dependencies
sudo apt-get update
sudo apt-get install -y build-essential libgtk-3-dev libwebkit2gtk-4.1-dev

# Install Go
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Windows Setup

```powershell
# Using Chocolatey
choco install golang nodejs wails-cli

# Or manual installation
# 1. Download Go from https://go.dev/dl/
# 2. Download Node.js from https://nodejs.org/
# 3. Install Wails: go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone & Configure

```bash
# Clone repository
git clone https://github.com/Molasses-Co/SweetDesk.git
cd SweetDesk

# Copy environment configuration
cp .env.example .env

# Add your Pixabay API key
# Edit .env and set PIXABAY_API_KEY=your_key_here

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Download ML models (if not committed)
go generate ./...
```

---

## Build & Run

### Development Mode

```bash
# Start development server (hot reload)
wails dev

# The app will:
# - Open browser at http://localhost:34115
# - Hot reload on frontend changes
# - Restart on backend changes (Go files)
```

### Production Build

```bash
# Build optimized binary
wails build -o SweetDesk

# Output locations:
# - macOS: ./build/bin/SweetDesk.app/Contents/MacOS/SweetDesk
# - Linux: ./build/bin/SweetDesk
# - Windows: ./build/bin/SweetDesk.exe
```

### Frontend-only Development

```bash
# If you want to develop frontend without Wails
cd frontend
npm run dev

# Runs on http://localhost:5173
# (Note: Backend APIs won't work in this mode)
```

### Backend-only Testing

```bash
# Test Go packages
go test ./...

# Test with coverage
go test -cover ./...

# Test specific package
go test ./internal/handlers/...
```

---

## Code Guidelines

### Go Code Style

#### Naming Conventions

```go
// Package names: lowercase, no underscores
package handlers

// Exported (public): PascalCase
func (a *App) ProcessImage(data []byte) error { }

// Unexported (private): camelCase
func (a *App) validateInput(data []byte) bool { }

// Constants: ALL_CAPS for unexported, PascalCase for exported
const (
    maxImageSize = 100 * 1024 * 1024  // 100MB
    defaultScale = 4
)

const (
    MaxConcurrency = 4
)

// Variables: camelCase
var processingQueue = make(chan Task, 10)
```

#### Error Handling

```go
// Always check errors
data, err := ioutil.ReadFile(filename)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// Wrap with context
result, err := a.coreProcessor.Process(data, options)
if err != nil {
    return fmt.Errorf("processing failed: %w", err)
}

// Log and return early
if err := validate(input); err != nil {
    logger.Error("validation failed", err)
    return nil, err
}
```

#### Function Structure

```go
// 1. Receiver (if method)
func (a *App) Method() error {

// 2. Validate inputs early
    if a.coreProcessor == nil {
        return fmt.Errorf("core processor not initialized")
    }
    
// 3. Acquire resources
    defer func() {
        // Cleanup if needed
    }()
    
// 4. Main logic
    result := process(data)
    
// 5. Return results
    return nil
}
```

### Frontend Code Style

#### React Components

```typescript
// Use functional components with hooks
interface ImageCardProps {
  imageUrl: string;
  onUpscale: (scale: number) => void;
}

const ImageCard: React.FC<ImageCardProps> = ({ imageUrl, onUpscale }) => {
  const [isLoading, setIsLoading] = React.useState(false);
  
  const handleUpscale = async (scale: number) => {
    setIsLoading(true);
    try {
      await onUpscale(scale);
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <div className="image-card">
      <img src={imageUrl} alt="Preview" />
      <button onClick={() => handleUpscale(4)} disabled={isLoading}>
        Upscale 4x
      </button>
    </div>
  );
};

export default ImageCard;
```

#### Wails Service Integration

```typescript
// src/services/imageService.ts
import * as wails from '@wailsjs/runtime';
import { App } from '@wailsjs/go/main';

export const imageService = {
  async upscaleImage(base64: string, scale: number): Promise<string> {
    try {
      return await App.UpscaleImage(base64, "photo", scale);
    } catch (error) {
      console.error('Upscale failed:', error);
      throw new Error(`Failed to upscale: ${error}`);
    }
  },

  async processImage(base64: string, resolution: string): Promise<string> {
    try {
      return await App.ProcessImage(base64, resolution, true, "16:9");
    } catch (error) {
      console.error('Processing failed:', error);
      throw new Error(`Failed to process: ${error}`);
    }
  },
};
```

---

## Testing Strategy

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestProcessImage ./internal/handlers/

# Generate coverage report
go test -cover ./... > coverage.txt
go tool cover -html=coverage.txt
```

### Test Structure

```go
// internal/handlers/image_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUpscaleImage(t *testing.T) {
    // Arrange
    testData := []byte{/* PNG data */}
    app := setupTestApp(t)
    
    // Act
    result, err := app.UpscaleImage(testData, "photo", 4)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
    assert.Greater(t, len(result), len(testData))
}
```

### Integration Tests

```bash
# Run integration tests (slower, full pipeline)
go test -tags integration ./...
```

### Frontend Tests (if added)

```bash
cd frontend
npm test                    # Run tests
npm run test:coverage       # Coverage report
npm run test:watch          # Watch mode
```

---

## Debugging

### Backend Debugging

#### VS Code Setup

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug SweetDesk",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/app",
      "env": {
        "PIXABAY_API_KEY": "your_key"
      },
      "args": []
    }
  ]
}
```

Then press `F5` to start debugging.

#### Console Logging

```go
import "fmt"
import "log"

// Simple logging
fmt.Printf("Debug: processing image %d bytes\n", len(data))

// Structured logging
log.Printf("[INFO] Upscaling: scale=%d\n", scale)

// In production, use a proper logger
// (e.g., github.com/sirupsen/logrus, go.uber.org/zap)
```

### Frontend Debugging

```typescript
// Console logging
console.log('Processing image:', imageData);
console.error('Error occurred:', error);

// Browser DevTools
// Press F12 to open DevTools
// - Go to Console tab
// - Set breakpoints
// - Inspect network requests

// React DevTools
// Install React DevTools extension
// Inspect component hierarchy and state
```

---

## Common Tasks

### Adding a New Image Processing Feature

1. **Backend API Handler** (`internal/handlers/new_feature.go`):

```go
func (a *App) NewFeature(base64Data string, params interface{}) (string, error) {
    if a.coreProcessor == nil {
        return "", fmt.Errorf("core processor not available")
    }
    
    data, err := a.imageProcessor.ConvertFromBase64(base64Data)
    if err != nil {
        return "", err
    }
    
    // Process
    result, err := a.coreProcessor.NewFeature(data, params)
    if err != nil {
        return "", err
    }
    
    return a.imageProcessor.ConvertToBase64(result), nil
}
```

2. **Frontend Component** (`frontend/src/components/NewFeature.tsx`):

```typescript
const NewFeature: React.FC = () => {
  const [result, setResult] = React.useState<string | null>(null);
  
  const handleProcess = async (image: string) => {
    const res = await App.NewFeature(image, {});
    setResult(res);
  };
  
  return (
    <div>
      <button onClick={() => handleProcess(imageData)}>Process</button>
      {result && <img src={`data:image/png;base64,${result}`} />}
    </div>
  );
};
```

3. **Test** (`internal/handlers/new_feature_test.go`):

```go
func TestNewFeature(t *testing.T) {
    app := setupTestApp(t)
    testImage := getTestImage(t)
    
    result, err := app.NewFeature(testImage, nil)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

### Updating Dependencies

```bash
# Check for updates
go list -u -m all

# Update specific package
go get -u github.com/package/name

# Update SweetDesk-core
go get -u github.com/pedro3pv/SweetDesk-core

# Tidy up
go mod tidy
```

### Adding Environment Variables

1. Update `.env.example`:

```bash
echo "NEW_VAR=value" >> .env.example
```

2. Load in code:

```go
var newVar = os.Getenv("NEW_VAR")
if newVar == "" {
    log.Fatal("NEW_VAR not set")
}
```

---

## CI/CD Pipeline

### GitHub Actions Workflows

#### Test Workflow (`.github/workflows/test.yml`)

```yaml
name: Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: 1.21
      - run: go test -v ./...
```

#### Build Workflow (`.github/workflows/build.yml`)

```yaml
name: Build

on:
  push:
    tags: [v*]

jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - uses: actions/setup-node@v3
      - run: wails build
      - uses: actions/upload-artifact@v3
        with:
          name: SweetDesk-${{ matrix.os }}
          path: build/bin/*
```

### Local Pre-commit Checks

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Format code
gofmt -w ./...

# Run tests
go test ./... || exit 1

# Check linting (if golangci-lint installed)
golangci-lint run ./...
```

Install hook:

```bash
chmod +x .git/hooks/pre-commit
```

---

## Performance Optimization

### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Caching

```go
// Cache classification results
classificationCache := make(map[string]string)

func (a *App) ClassifyImageCached(hash string, data []byte) (string, error) {
    if result, ok := classificationCache[hash]; ok {
        return result, nil
    }
    
    result, err := a.coreProcessor.Classify(data)
    if err == nil {
        classificationCache[hash] = result.Type
    }
    
    return result.Type, err
}
```

---

## Resources

- [Wails Documentation](https://wails.io/docs/gettingstarted/firstproject)
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [SweetDesk-Core Repository](https://github.com/pedro3pv/SweetDesk-core)

---

## Getting Help

1. Check existing [GitHub Issues](https://github.com/Molasses-Co/SweetDesk/issues)
2. Review [README.md](./README.md) and [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)
3. Open a new issue with:
   - Description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, Node version)

Happy coding! 🚀

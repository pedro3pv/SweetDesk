# SweetDesk API Documentation

## Backend Go API (Wails Bindings)

The Go backend provides the following methods that can be called from the frontend via Wails bindings:

### Image Search

#### `SearchImages(query string, page int, perPage int) ([]ImageResult, error)`

Search for images using the Pixabay API.

**Parameters:**
- `query`: Search query string (e.g., "nature", "mountains")
- `page`: Page number for pagination (starting from 1)
- `perPage`: Number of results per page (max 100)

**Returns:**
- Array of `ImageResult` objects containing image metadata
- Error if the search fails

**Example:**
```javascript
const images = await window.go.main.App.SearchImages("mountains", 1, 12);
```

### Image Download

#### `DownloadImage(imageURL string) (string, error)`

Download an image from a URL and return it as base64.

**Parameters:**
- `imageURL`: Full URL to the image to download

**Returns:**
- Base64-encoded image data
- Error if download fails

**Example:**
```javascript
const base64Data = await window.go.main.App.DownloadImage(imageURL);
```

### Image Classification

#### `ClassifyImage(base64Data string) (string, error)`

Classify an image as "anime", "photo", or "art".

**Parameters:**
- `base64Data`: Base64-encoded image data

**Returns:**
- Classification type: "anime", "photo", or "art"
- Error if classification fails

**Note:** Falls back to "photo" if Python bridge is not available.

**Example:**
```javascript
const imageType = await window.go.main.App.ClassifyImage(base64Data);
```

### Image Upscaling

#### `UpscaleImage(base64Data string, imageType string, scale int) (string, error)`

Upscale an image using the appropriate AI model.

**Parameters:**
- `base64Data`: Base64-encoded image data
- `imageType`: Type of image ("anime" or "photo")
- `scale`: Upscaling factor (2, 4, or 8)

**Returns:**
- Base64-encoded upscaled image data
- Error if upscaling fails

**Models Used:**
- Anime: RealCUGAN-pro
- Photo: Real-ESRGAN (4xLSDIR)

**Example:**
```javascript
const upscaled = await window.go.main.App.UpscaleImage(base64Data, "photo", 4);
```

### Full Image Processing

#### `ProcessImage(base64Data string, targetResolution string, useSeamCarving bool) (string, error)`

Complete processing pipeline: classify, upscale, and adjust aspect ratio.

**Parameters:**
- `base64Data`: Base64-encoded input image data
- `targetResolution`: Target resolution ("4K", "5K", or "8K")
- `useSeamCarving`: Use content-aware resize (true) or simple crop (false)

**Returns:**
- Base64-encoded processed image data
- Error if processing fails

**Processing Steps:**
1. Classify image type (anime vs photo)
2. Upscale using appropriate model
3. Adjust aspect ratio to 16:9
4. Return final result

**Example:**
```javascript
const result = await window.go.main.App.ProcessImage(base64Data, "4K", true);
```

## Data Structures

### ImageResult

```typescript
interface ImageResult {
    id: string;          // Unique image ID
    url: string;         // Image page URL
    downloadURL: string; // Direct download URL
    previewURL: string;  // Thumbnail preview URL
    width: number;       // Original width in pixels
    height: number;      // Original height in pixels
    author: string;      // Author/photographer name
    source: string;      // Source platform (e.g., "Pixabay")
    tags: string[];      // Array of tags
}
```

## Environment Variables

### Required

- `PIXABAY_API_KEY`: Your Pixabay API key (get it from https://pixabay.com/api/docs/)

### Optional

- `SWEETDESK_DEBUG`: Enable debug logging (set to "1")
- `SUPABASE_URL`: Supabase project URL (for cloud storage)
- `SUPABASE_KEY`: Supabase anonymous key

## Error Handling

All methods return errors that should be handled in the frontend:

```javascript
try {
    const result = await window.go.main.App.ProcessImage(data, "4K", false);
    // Handle success
} catch (error) {
    console.error("Processing failed:", error);
    // Show error to user
}
```

## Performance Notes

### Processing Times (Estimated)

| Operation | Time (M1 Pro) | Time (Intel i7) |
|-----------|---------------|-----------------|
| Classification | ~0.5s | ~1s |
| Upscale 2x | ~15s | ~25s |
| Upscale 4x | ~45s | ~60s |
| Upscale 8x | ~120s | ~180s |
| Seam Carving | ~10s | ~15s |

### Memory Usage

- Original image: ~10-20 MB
- 4K upscaled: ~30-50 MB
- All processing is in-memory (no disk I/O until save)

## Adding New API Providers

To add support for additional image APIs (Unsplash, Pexels, etc.):

1. Implement the `APIProvider` interface in Go:

```go
type APIProvider interface {
    Search(query string, options SearchOptions) ([]ImageResult, error)
    Download(imageURL string) ([]byte, error)
    GetName() string
}
```

2. Create a new provider struct (e.g., `UnsplashProvider`)

3. Register it in `app.go` startup function

4. Update frontend to allow provider selection

## Python Scripts

### classify_image.py

Classifies images using heuristic-based analysis. For production, consider using:
- DeepGHS/imgutils
- TensorFlow models
- PyTorch classifiers

### seam_carving.py

Content-aware image resizing. For production, consider using:
- seam-carving library (faster)
- OpenCV implementations
- GPU-accelerated versions

## Binary Distribution

AI upscaler binaries should be downloaded from official sources:

- **Real-ESRGAN**: https://github.com/xinntao/Real-ESRGAN/releases
- **RealCUGAN**: https://github.com/nihui/realcugan-ncnn-vulkan/releases

Place binaries in:
```
binaries/
├── darwin/     # macOS
├── linux/      # Linux
└── windows/    # Windows
```

## License Compliance

| Component | License | Commercial Use |
|-----------|---------|----------------|
| Real-ESRGAN-ncnn-vulkan | MIT-like | ✅ Yes |
| RealCUGAN-ncnn-vulkan | MIT-like | ✅ Yes |
| DeepGHS/imgutils | MIT | ✅ Yes |
| seam-carving | MIT | ✅ Yes |
| Pixabay API | Free | ✅ Yes (with attribution) |

**Note:** If using Universal-NCNN-Upscaler instead, be aware of AGPL-3.0 requirements.

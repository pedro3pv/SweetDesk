package seamcarver

import (
"image"
"image/color"
"testing"
)

func TestResizeToExactDimensions(t *testing.T) {
// Create a test image (100x100 red square)
img := image.NewRGBA(image.Rect(0, 0, 100, 100))
for y := 0; y < 100; y++ {
for x := 0; x < 100; x++ {
img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
}
}

tests := []struct {
name         string
targetWidth  int
targetHeight int
wantWidth    int
wantHeight   int
}{
{
name:         "Upscale width, keep height",
targetWidth:  150,
targetHeight: 100,
wantWidth:    150,
wantHeight:   100,
},
{
name:         "Downscale width, upscale height",
targetWidth:  80,
targetHeight: 120,
wantWidth:    80,
wantHeight:   120,
},
{
name:         "Same dimensions",
targetWidth:  100,
targetHeight: 100,
wantWidth:    100,
wantHeight:   100,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
carver := NewSeamCarver(img)
result, err := carver.ResizeToExactDimensions(ResizeOptions{
TargetWidth:     tt.targetWidth,
TargetHeight:    tt.targetHeight,
MaxDeltaBySeams: 300,
})

if err != nil {
t.Fatalf("ResizeToExactDimensions() error = %v", err)
}

bounds := result.Bounds()
gotWidth := bounds.Dx()
gotHeight := bounds.Dy()

if gotWidth != tt.wantWidth {
t.Errorf("Width = %d, want %d", gotWidth, tt.wantWidth)
}
if gotHeight != tt.wantHeight {
t.Errorf("Height = %d, want %d", gotHeight, tt.wantHeight)
}
})
}
}

func TestAspectRatioAdjustment(t *testing.T) {
// Create a test image (100x100 square)
img := image.NewRGBA(image.Rect(0, 0, 100, 100))
for y := 0; y < 100; y++ {
for x := 0; x < 100; x++ {
img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
}

carver := NewSeamCarver(img)
result, err := carver.AdjustAspectRatio(AspectRatioOptions{
TargetRatio:     16.0 / 9.0,
MaxDeltaBySeams: 300,
})

if err != nil {
t.Fatalf("AdjustAspectRatio() error = %v", err)
}

bounds := result.Bounds()
width := bounds.Dx()
height := bounds.Dy()
actualRatio := float64(width) / float64(height)
expectedRatio := 16.0 / 9.0

// Check if ratio is within acceptable range
diff := actualRatio - expectedRatio
if diff < 0 {
diff = -diff
}

if diff > 0.1 {
t.Errorf("Aspect ratio = %.3f, want %.3f (diff: %.3f)", actualRatio, expectedRatio, diff)
}
}

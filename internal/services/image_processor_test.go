package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveToFile(t *testing.T) {
	ip := NewImageProcessor(context.Background())

	t.Run("saves file to specified path", func(t *testing.T) {
		tmpDir := t.TempDir()
		data := []byte("fake-image-data")
		fileName := "test-image.png"

		fullPath, err := ip.SaveToFile(data, tmpDir, fileName)
		if err != nil {
			t.Fatalf("SaveToFile failed: %v", err)
		}

		expected := filepath.Join(tmpDir, fileName)
		if fullPath != expected {
			t.Errorf("expected path %s, got %s", expected, fullPath)
		}

		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		if string(content) != string(data) {
			t.Errorf("saved content mismatch")
		}
	})

	t.Run("creates directory if it does not exist", func(t *testing.T) {
		tmpDir := filepath.Join(t.TempDir(), "nested", "dir")
		data := []byte("fake-image-data")
		fileName := "test-image.png"

		fullPath, err := ip.SaveToFile(data, tmpDir, fileName)
		if err != nil {
			t.Fatalf("SaveToFile failed: %v", err)
		}

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("file was not created at %s", fullPath)
		}

		_ = fullPath
	})

	t.Run("adds .png extension when missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		data := []byte("fake-image-data")
		fileName := "test-image"

		fullPath, err := ip.SaveToFile(data, tmpDir, fileName)
		if err != nil {
			t.Fatalf("SaveToFile failed: %v", err)
		}

		expected := filepath.Join(tmpDir, "test-image.png")
		if fullPath != expected {
			t.Errorf("expected path %s, got %s", expected, fullPath)
		}
	})

	t.Run("returns error for empty save path", func(t *testing.T) {
		_, err := ip.SaveToFile([]byte("data"), "", "file.png")
		if err == nil {
			t.Error("expected error for empty save path")
		}
	})
}

func TestConvertBase64RoundTrip(t *testing.T) {
	ip := NewImageProcessor(context.Background())

	original := []byte("test image data for base64 conversion")
	encoded := ip.ConvertToBase64(original)
	decoded, err := ip.ConvertFromBase64(encoded)
	if err != nil {
		t.Fatalf("ConvertFromBase64 failed: %v", err)
	}

	if string(decoded) != string(original) {
		t.Error("base64 round-trip failed: data mismatch")
	}
}

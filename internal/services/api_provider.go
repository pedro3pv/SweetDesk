package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// APIProvider defines the interface for image API providers
type APIProvider interface {
	Search(query string, options SearchOptions) ([]ImageResult, error)
	Download(imageURL string) ([]byte, error)
	GetName() string
}

// SearchOptions contains search parameters
type SearchOptions struct {
	Page     int
	PerPage  int
	Category string
	Orientation string // "horizontal", "vertical", "all"
	MinWidth int
	MinHeight int
}

// ImageResult represents a single image result
type ImageResult struct {
	ID          string
	URL         string
	DownloadURL string
	PreviewURL  string
	Width       int
	Height      int
	Author      string
	Source      string
	Tags        []string
}

// PixabayProvider implements APIProvider for Pixabay
type PixabayProvider struct {
	apiKey string
	client *http.Client
	ctx    context.Context
}

// NewPixabayProvider creates a new Pixabay provider
func NewPixabayProvider(ctx context.Context, apiKey string) *PixabayProvider {
	return &PixabayProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
		ctx:    ctx,
	}
}

// GetName returns the provider name
func (p *PixabayProvider) GetName() string {
	return "Pixabay"
}

// Search searches for images on Pixabay
func (p *PixabayProvider) Search(query string, options SearchOptions) ([]ImageResult, error) {
	baseURL := "https://pixabay.com/api/"
	
	params := url.Values{}
	params.Set("key", p.apiKey)
	params.Set("q", query)
	params.Set("image_type", "photo")
	params.Set("page", fmt.Sprintf("%d", options.Page))
	params.Set("per_page", fmt.Sprintf("%d", options.PerPage))
	
	if options.MinWidth > 0 {
		params.Set("min_width", fmt.Sprintf("%d", options.MinWidth))
	}
	if options.MinHeight > 0 {
		params.Set("min_height", fmt.Sprintf("%d", options.MinHeight))
	}
	if options.Orientation != "" && options.Orientation != "all" {
		params.Set("orientation", options.Orientation)
	}
	
	fullURL := baseURL + "?" + params.Encode()
	
	req, err := http.NewRequestWithContext(p.ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Total int `json:"total"`
		Hits  []struct {
			ID            int    `json:"id"`
			PageURL       string `json:"pageURL"`
			PreviewURL    string `json:"previewURL"`
			LargeImageURL string `json:"largeImageURL"`
			ImageWidth    int    `json:"imageWidth"`
			ImageHeight   int    `json:"imageHeight"`
			User          string `json:"user"`
			Tags          string `json:"tags"`
		} `json:"hits"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	images := make([]ImageResult, 0, len(result.Hits))
	for _, hit := range result.Hits {
		images = append(images, ImageResult{
			ID:          fmt.Sprintf("%d", hit.ID),
			URL:         hit.PageURL,
			DownloadURL: hit.LargeImageURL,
			PreviewURL:  hit.PreviewURL,
			Width:       hit.ImageWidth,
			Height:      hit.ImageHeight,
			Author:      hit.User,
			Source:      "Pixabay",
			Tags:        parseTags(hit.Tags),
		})
	}
	
	return images, nil
}

// Download downloads an image from URL
func (p *PixabayProvider) Download(imageURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(p.ctx, "GET", imageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned status %d", resp.StatusCode)
	}
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}
	
	return data, nil
}

// parseTags splits comma-separated tags into a slice
func parseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}
	result := []string{}
	for _, tag := range splitString(tags, ",") {
		trimmed := trim(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	current := ""
	for _, c := range s {
		if string(c) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

#!/usr/bin/env python3
"""
Image classifier using simple heuristics
Classifies images as anime, photo, or art based on visual characteristics
"""

import sys
import json
from PIL import Image
import numpy as np

def classify_image(image_path):
    """
    Classify an image as anime, photo, or art
    
    This is a simple heuristic-based classifier.
    For production, use DeepGHS/imgutils or similar ML models.
    """
    try:
        img = Image.open(image_path)
        img_array = np.array(img.convert('RGB'))
        
        # Calculate image statistics
        
        # Calculate color saturation
        hsv = rgb_to_hsv(img_array)
        saturation = np.mean(hsv[:, :, 1])
        
        # Calculate edge density (simple Sobel-like)
        edges = calculate_edge_density(img_array)
        
        # Heuristic classification
        # High saturation + high edge density = anime
        # Low saturation + moderate edges = photo
        # High variation = art
        
        if saturation > 0.5 and edges > 0.3:
            classification = "anime"
            confidence = 0.7
        elif saturation < 0.4 and edges < 0.4:
            classification = "photo"
            confidence = 0.8
        else:
            classification = "art"
            confidence = 0.6
        
        result = {
            "type": classification,
            "confidence": confidence,
            "model": "heuristic-v1"
        }
        
        print(json.dumps(result))
        return 0
        
    except Exception as e:
        error_result = {
            "type": "photo",  # Default fallback
            "confidence": 0.5,
            "model": "heuristic-v1",
            "error": str(e)
        }
        print(json.dumps(error_result))
        return 1

def rgb_to_hsv(rgb):
    """Convert RGB to HSV color space"""
    rgb = rgb.astype(float) / 255.0
    maxc = np.max(rgb, axis=2)
    minc = np.min(rgb, axis=2)
    
    v = maxc
    s = np.where(maxc != 0, (maxc - minc) / maxc, 0)
    
    rc = np.where(maxc != minc, (maxc - rgb[:,:,0]) / (maxc - minc), 0)
    gc = np.where(maxc != minc, (maxc - rgb[:,:,1]) / (maxc - minc), 0)
    bc = np.where(maxc != minc, (maxc - rgb[:,:,2]) / (maxc - minc), 0)
    
    h = np.zeros_like(v)
    h = np.where(rgb[:,:,0] == maxc, bc - gc, h)
    h = np.where(rgb[:,:,1] == maxc, 2.0 + rc - bc, h)
    h = np.where(rgb[:,:,2] == maxc, 4.0 + gc - rc, h)
    h = np.where(minc == maxc, 0, h)
    
    h = (h / 6.0) % 1.0
    
    hsv = np.dstack((h, s, v))
    return hsv

def calculate_edge_density(img_array):
    """Calculate edge density using simple gradient"""
    gray = np.mean(img_array, axis=2)
    
    # Simple horizontal and vertical gradients
    dx = np.abs(np.diff(gray, axis=1))
    dy = np.abs(np.diff(gray, axis=0))
    
    # Normalize and calculate density
    edge_density = (np.mean(dx) + np.mean(dy)) / 255.0
    return edge_density

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No image path provided"}), file=sys.stderr)
        sys.exit(1)
    
    image_path = sys.argv[1]
    sys.exit(classify_image(image_path))

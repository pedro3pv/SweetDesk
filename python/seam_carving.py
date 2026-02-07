#!/usr/bin/env python3
"""
Seam carving implementation for content-aware image resizing
Based on the seam carving algorithm
"""

import sys
import numpy as np
from PIL import Image

def calculate_energy(img):
    """Calculate energy map using gradient magnitude"""
    gray = np.mean(img, axis=2)
    
    # Calculate gradients
    dx = np.zeros_like(gray)
    dy = np.zeros_like(gray)
    
    dx[:, :-1] = np.abs(gray[:, 1:] - gray[:, :-1])
    dy[:-1, :] = np.abs(gray[1:, :] - gray[:-1, :])
    
    energy = dx + dy
    return energy

def find_vertical_seam(energy):
    """Find the vertical seam with minimum energy"""
    h, w = energy.shape
    dp = np.copy(energy)
    
    # Dynamic programming to find minimum energy path
    for i in range(1, h):
        for j in range(w):
            if j == 0:
                dp[i, j] += min(dp[i-1, j], dp[i-1, j+1])
            elif j == w - 1:
                dp[i, j] += min(dp[i-1, j-1], dp[i-1, j])
            else:
                dp[i, j] += min(dp[i-1, j-1], dp[i-1, j], dp[i-1, j+1])
    
    # Backtrack to find seam
    seam = np.zeros(h, dtype=int)
    seam[-1] = np.argmin(dp[-1])
    
    for i in range(h-2, -1, -1):
        j = seam[i+1]
        if j == 0:
            seam[i] = j + np.argmin(dp[i, j:j+2])
        elif j == w - 1:
            seam[i] = j - 1 + np.argmin(dp[i, j-1:j+1])
        else:
            seam[i] = j - 1 + np.argmin(dp[i, j-1:j+2])
    
    return seam

def remove_vertical_seam(img, seam):
    """Remove a vertical seam from the image"""
    h, w, c = img.shape
    new_img = np.zeros((h, w-1, c), dtype=img.dtype)
    
    for i in range(h):
        j = seam[i]
        new_img[i, :j] = img[i, :j]
        new_img[i, j:] = img[i, j+1:]
    
    return new_img

def add_vertical_seam(img, seam):
    """Add a vertical seam to the image by duplicating pixels"""
    h, w, c = img.shape
    new_img = np.zeros((h, w+1, c), dtype=img.dtype)
    
    for i in range(h):
        j = seam[i]
        new_img[i, :j] = img[i, :j]
        # Duplicate the seam pixel
        new_img[i, j] = img[i, j]
        new_img[i, j+1:] = img[i, j:]
    
    return new_img

def seam_carving_resize(img, target_width, target_height):
    """
    Resize image using seam carving
    
    This is a simplified implementation. For production, use
    a more robust library like `seam-carving` or `PySeam`.
    """
    h, w, c = img.shape
    
    # Resize width first
    while w > target_width:
        energy = calculate_energy(img)
        seam = find_vertical_seam(energy)
        img = remove_vertical_seam(img, seam)
        h, w, c = img.shape
    
    while w < target_width:
        energy = calculate_energy(img)
        seam = find_vertical_seam(energy)
        img = add_vertical_seam(img, seam)
        h, w, c = img.shape
    
    # Then resize height by rotating
    if h != target_height:
        img = np.transpose(img, (1, 0, 2))
        h, w, c = img.shape
        
        while w > target_height:
            energy = calculate_energy(img)
            seam = find_vertical_seam(energy)
            img = remove_vertical_seam(img, seam)
            h, w, c = img.shape
        
        while w < target_height:
            energy = calculate_energy(img)
            seam = find_vertical_seam(energy)
            img = add_vertical_seam(img, seam)
            h, w, c = img.shape
        
        img = np.transpose(img, (1, 0, 2))
    
    return img

def main(input_path, output_path, target_width, target_height, energy_mode="forward"):
    """Main function to perform seam carving"""
    try:
        # Load image
        img = Image.open(input_path)
        img_array = np.array(img)
        
        # Ensure RGB
        if len(img_array.shape) == 2:
            img_array = np.stack([img_array] * 3, axis=2)
        elif img_array.shape[2] == 4:
            img_array = img_array[:, :, :3]
        
        # Perform seam carving
        result = seam_carving_resize(img_array, int(target_width), int(target_height))
        
        # Save result
        output_img = Image.fromarray(result.astype(np.uint8))
        output_img.save(output_path)
        
        return 0
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1

if __name__ == "__main__":
    if len(sys.argv) < 5:
        print("Usage: seam_carving.py <input> <output> <width> <height> [energy_mode]", file=sys.stderr)
        sys.exit(1)
    
    input_path = sys.argv[1]
    output_path = sys.argv[2]
    target_width = sys.argv[3]
    target_height = sys.argv[4]
    energy_mode = sys.argv[5] if len(sys.argv) > 5 else "forward"
    
    sys.exit(main(input_path, output_path, target_width, target_height, energy_mode))

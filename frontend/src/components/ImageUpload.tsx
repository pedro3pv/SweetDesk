'use client';

import { useCallback } from 'react';

interface ImageUploadProps {
    onImageSelect: (data: string) => void;
}

export default function ImageUpload({ onImageSelect }: ImageUploadProps) {
    const handleFileSelect = useCallback(async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (event) => {
            const result = event.target?.result as string;
            // Remove data:image/...;base64, prefix
            const base64 = result.split(',')[1];
            onImageSelect(base64);
        };
        reader.readAsDataURL(file);
    }, [onImageSelect]);

    const handleDrop = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        const file = e.dataTransfer.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (event) => {
            const result = event.target?.result as string;
            const base64 = result.split(',')[1];
            onImageSelect(base64);
        };
        reader.readAsDataURL(file);
    }, [onImageSelect]);

    const handleDragOver = useCallback((e: React.DragEvent) => {
        e.preventDefault();
    }, []);

    return (
        <div className="space-y-4">
            <div
                onDrop={handleDrop}
                onDragOver={handleDragOver}
                className="border-2 border-dashed border-purple-300 dark:border-dark-border rounded-xl p-6 lg:p-8 text-center hover:border-purple-500 dark:hover:border-purple-600 hover:bg-purple-50/50 dark:hover:bg-purple-900/10 transition-all cursor-pointer"
            >
                <input
                    type="file"
                    accept="image/*"
                    onChange={handleFileSelect}
                    className="hidden"
                    id="file-upload"
                />
                <label htmlFor="file-upload" className="cursor-pointer">
                    <div className="space-y-2">
                        <div className="text-4xl">📁</div>
                        <p className="text-sm text-purple-700 dark:text-purple-300">
                            Drop an image here or click to browse
                        </p>
                        <p className="text-xs text-purple-500 dark:text-purple-400">
                            Supports: JPG, PNG, WebP
                        </p>
                    </div>
                </label>
            </div>

            <div className="text-center text-xs lg:text-sm text-purple-600 dark:text-purple-400">
                <p>Maximum file size: 10MB</p>
            </div>
        </div>
    );
}

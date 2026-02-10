'use client';

import React from "react"

import { useCallback, useState } from 'react';
import type { ImageResult } from '../lib/types';

interface ImageUploadProps {
    onImageSelect: (image: ImageResult) => void;
}

export default function ImageUpload({ onImageSelect }: ImageUploadProps) {
    const [isDragging, setIsDragging] = useState(false);

    const processFile = useCallback((file: File) => {
        const reader = new FileReader();
        reader.onload = (event) => {
            const dataUrl = event.target?.result as string;
            const img = new window.Image();
            img.crossOrigin = 'anonymous';
            img.onload = () => {
                const image: ImageResult = {
                    id: `upload-${Date.now()}`,
                    url: dataUrl,
                    downloadURL: dataUrl,
                    previewURL: dataUrl,
                    width: img.naturalWidth,
                    height: img.naturalHeight,
                    author: 'Local',
                    source: 'upload',
                    tags: ['upload'],
                    description: file.name,
                };
                onImageSelect(image);
            };
            img.src = dataUrl;
        };
        reader.readAsDataURL(file);
    }, [onImageSelect]);

    const handleFileSelect = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const files = e.target.files;
        if (files && files.length > 0) {
            Array.from(files).forEach(processFile);
        }
    }, [processFile]);
    
    const handleDrop = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
        const files = e.dataTransfer.files;
        if (files && files.length > 0) {
            Array.from(files).forEach(processFile);
        }
    }, [processFile]);

    const handleDragOver = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(true);
    }, []);

    const handleDragLeave = useCallback((e: React.DragEvent) => {
        e.preventDefault();
        setIsDragging(false);
    }, []);

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-center gap-3 px-6 py-4 border-b border-border">
                <h2 className="text-base lg:text-lg font-semibold text-foreground tracking-tight">Adicionar Arquivo</h2>
            </div>

            {/* Drop Zone */}
            <div className="flex-1 flex items-center justify-center p-8 lg:p-12">
                <div
                    onDrop={handleDrop}
                    onDragOver={handleDragOver}
                    onDragLeave={handleDragLeave}
                    className={`w-full max-w-2xl aspect-video rounded-xl border-2 border-dashed flex flex-col items-center justify-center gap-6 transition-colors cursor-pointer ${
                        isDragging
                            ? 'border-primary bg-primary/5'
                            : 'border-border hover:border-muted-foreground'
                    }`}
                >
                    <input
                        type="file"
                        accept="image/*"
                        multiple
                        onChange={handleFileSelect}
                        className="hidden"
                        id="file-upload"
                    />
                    <label htmlFor="file-upload" className="cursor-pointer flex flex-col items-center gap-6">
                        <svg className="w-16 h-16 lg:w-20 lg:h-20 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
                        </svg>
                        <div className="text-center">
                            <p className="text-base lg:text-lg text-foreground font-medium">Arraste para aqui</p>
                            <p className="text-sm lg:text-base text-muted-foreground mt-2">ou clique para selecionar</p>
                        </div>
                        <p className="text-sm text-muted-foreground">JPG, PNG, WebP</p>
                    </label>
                </div>
            </div>
        </div>
    );
}

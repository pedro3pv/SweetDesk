'use client';

import React, { createContext, useContext, useState, useCallback } from 'react';
import type { ImageResult, DownloadItem } from '../lib/types';

// ─── Download Context ───────────────────────────────────────────────

interface DownloadContextType {
    items: DownloadItem[];
    selectedImage: ImageResult | null;
    savePath: string;
    setSelectedImage: (image: ImageResult | null) => void;
    setSavePath: (path: string) => void;
    addItem: (image: ImageResult, overrides?: Partial<DownloadItem>) => void;
    addDownloadItem: (item: DownloadItem) => void;
    updateItem: (id: string, updates: Partial<DownloadItem>) => void;
    removeItem: (id: string) => void;
    clearAll: () => void;
}

const DownloadContext = createContext<DownloadContextType | null>(null);

export function useDownloads(): DownloadContextType {
    const ctx = useContext(DownloadContext);
    if (!ctx) throw new Error('useDownloads must be used within AppProvider');
    return ctx;
}

// ─── Provider ───────────────────────────────────────────────────────

export function AppProvider({ children }: { children: React.ReactNode }) {
    const [items, setItems] = useState<DownloadItem[]>([]);
    const [selectedImage, setSelectedImage] = useState<ImageResult | null>(null);
    const [savePath, setSavePath] = useState('');

    const addItem = useCallback((image: ImageResult, overrides?: Partial<DownloadItem>) => {
        setItems(prev => {
            if (prev.some(p => p.image.id === image.id)) return prev;
            return [...prev, {
                id: `dl-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
                image,
                name: image.description || `wallpaper-${image.id}`,
                upscale: true,
                dimension: '3840x2160',
                aspect: '16:9',
                selected: true,
                ...overrides,
            }];
        });
    }, []);

    const addDownloadItem = useCallback((item: DownloadItem) => {
        setItems(prev => {
            if (prev.some(p => p.image.id === item.image.id)) return prev;
            return [...prev, item];
        });
    }, []);

    const updateItem = useCallback((id: string, updates: Partial<DownloadItem>) => {
        setItems(prev => prev.map(item =>
            item.id === id ? { ...item, ...updates } : item
        ));
    }, []);

    const removeItem = useCallback((id: string) => {
        setItems(prev => prev.filter(item => item.id !== id));
    }, []);

    const clearAll = useCallback(() => {
        setItems([]);
        setSelectedImage(null);
        setSavePath('');
    }, []);

    return (
        <DownloadContext.Provider value={{
            items,
            selectedImage,
            savePath,
            setSelectedImage,
            setSavePath,
            addItem,
            addDownloadItem,
            updateItem,
            removeItem,
            clearAll,
        }}>
            {children}
        </DownloadContext.Provider>
    );
}

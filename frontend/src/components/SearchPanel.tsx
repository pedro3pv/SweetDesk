'use client';

import React from "react"

import { useState, useCallback } from 'react';
import type { ImageResult, FilterState } from '@/lib/types';
import { SearchImages } from '@wailsjs/go/main/App';

interface SearchPanelProps {
    onImageSelect: (image: ImageResult) => void;
    onAddToList: (image: ImageResult) => void;
}

// Mock data for preview (in Wails, this would call the Go backend)
const MOCK_IMAGES: ImageResult[] = Array.from({ length: 18 }, (_, i) => ({
    id: `img-${i + 1}`,
    url: `https://picsum.photos/seed/${i + 1}/800/600`,
    downloadURL: `https://picsum.photos/seed/${i + 1}/1920/1080`,
    previewURL: `https://picsum.photos/seed/${i + 1}/300/200`,
    width: 1920,
    height: 1080,
    author: ['unsplash', 'pinterest', 'unsplash', 'unsplash', 'pinterest', 'unsplash'][i % 6],
    source: ['unsplash', 'pinterest'][i % 2],
    tags: ['wallpaper', 'nature', 'landscape'],
    description: 'A beautiful high-resolution wallpaper perfect for your desktop.',
}));

export default function SearchPanel({ onImageSelect, onAddToList }: SearchPanelProps) {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<ImageResult[]>(MOCK_IMAGES);
    const [isSearching, setIsSearching] = useState(false);
    const [filters, setFilters] = useState<FilterState>({
        aspect: '',
        dimensions: '',
        type: '',
    });
    const [showFilters, setShowFilters] = useState(true);

    const handleSearch = useCallback(async () => {
        if (!query.trim()) return;

        setIsSearching(true);
        try {
            // Use Wails backend for search
            if (typeof SearchImages !== 'undefined') {
                const images = await SearchImages(query, 1, 18);
                if (images && images.length > 0) {
                    setResults(images.map((img) => ({
                        id: img.id,
                        url: img.url,
                        downloadURL: img.downloadURL,
                        previewURL: img.previewURL,
                        width: img.width,
                        height: img.height,
                        author: img.author,
                        source: img.source,
                        tags: img.tags || [],
                        description: '', // Add if available in backend
                    })));
                } else {
                    setResults([]);
                }
            } else {
                // Fallback to mock data for development
                await new Promise(r => setTimeout(r, 600));
                setResults(MOCK_IMAGES);
            }
        } catch (error) {
            console.error('Search failed:', error);
            // Fallback to mock on error
            setResults(MOCK_IMAGES);
        } finally {
            setIsSearching(false);
        }
    }, [query]);

    const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleSearch();
    }, [handleSearch]);

    const filteredResults = results.filter(img => {
        if (filters.type && img.source !== filters.type) return false;
        if (filters.aspect) {
            const ratio = img.width / img.height;
            if (filters.aspect === '16:9' && Math.abs(ratio - 16 / 9) > 0.2) return false;
            if (filters.aspect === '4:3' && Math.abs(ratio - 4 / 3) > 0.2) return false;
            if (filters.aspect === '1:1' && Math.abs(ratio - 1) > 0.2) return false;
            if (filters.aspect === '21:9' && Math.abs(ratio - 21 / 9) > 0.2) return false;
        }
        return true;
    });

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-center gap-3 px-5 py-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground tracking-tight">SWEETDESK</h2>
            </div>

            {/* Search Bar */}
            <div className="flex items-center gap-2 px-5 py-3 border-b border-border">
                <div className="flex-1 flex items-center gap-2 bg-input rounded-md px-3 py-2">
                    <svg className="w-4 h-4 text-muted-foreground shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <input
                        type="text"
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                        onKeyDown={handleKeyDown}
                        placeholder="Buscar wallpapers..."
                        className="flex-1 bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
                    />
                </div>
                <button
                    onClick={handleSearch}
                    disabled={isSearching}
                    className="p-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 disabled:opacity-50 transition-colors"
                    aria-label="Buscar"
                >
                    {isSearching ? (
                        <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
                        </svg>
                    ) : (
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    )}
                </button>
                <button
                    onClick={() => setShowFilters(!showFilters)}
                    className={`p-2 rounded-md transition-colors ${showFilters ? 'bg-primary text-primary-foreground' : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'}`}
                    aria-label="Filtros"
                >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
                    </svg>
                </button>
            </div>

            {/* Content area with optional filter sidebar */}
            <div className="flex flex-1 overflow-hidden">
                {/* Filters Sidebar */}
                {showFilters && (
                    <div className="w-40 shrink-0 border-r border-border p-3 overflow-y-auto">
                        {/* Filter: Filtros label */}
                        <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3">Filtros</p>

                        {/* Aspect Ratio */}
                        <div className="mb-4">
                            <p className="text-xs font-medium text-foreground mb-2">Aspecto</p>
                            <div className="flex flex-col gap-1">
                                {['', '16:9', '4:3', '1:1', '21:9'].map(a => (
                                    <button
                                        key={a}
                                        onClick={() => setFilters(f => ({ ...f, aspect: a }))}
                                        className={`text-left text-xs px-2 py-1.5 rounded transition-colors ${
                                            filters.aspect === a
                                                ? 'bg-primary text-primary-foreground'
                                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                        }`}
                                    >
                                        {a || 'Todos'}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Dimensions */}
                        <div className="mb-4">
                            <p className="text-xs font-medium text-foreground mb-2">Dimensoes</p>
                            <div className="flex flex-col gap-1">
                                {['', 'HD', 'FHD', '4K', '8K'].map(d => (
                                    <button
                                        key={d}
                                        onClick={() => setFilters(f => ({ ...f, dimensions: d }))}
                                        className={`text-left text-xs px-2 py-1.5 rounded transition-colors ${
                                            filters.dimensions === d
                                                ? 'bg-primary text-primary-foreground'
                                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                        }`}
                                    >
                                        {d || 'Todos'}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {/* Type */}
                        <div className="mb-4">
                            <p className="text-xs font-medium text-foreground mb-2">Tipo</p>
                            <div className="flex flex-col gap-1">
                                {['', 'unsplash', 'pinterest'].map(t => (
                                    <button
                                        key={t}
                                        onClick={() => setFilters(f => ({ ...f, type: t }))}
                                        className={`text-left text-xs px-2 py-1.5 rounded transition-colors capitalize ${
                                            filters.type === t
                                                ? 'bg-primary text-primary-foreground'
                                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                                        }`}
                                    >
                                        {t || 'Todos'}
                                    </button>
                                ))}
                            </div>
                        </div>
                    </div>
                )}

                {/* Image Grid */}
                <div className="flex-1 overflow-y-auto p-3">
                    {filteredResults.length > 0 ? (
                        <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-6 gap-2">
                            {filteredResults.map((image) => (
                                <button
                                    key={image.id}
                                    onClick={() => onImageSelect(image)}
                                    className="group relative aspect-square rounded-md overflow-hidden border border-border hover:border-primary transition-all focus:outline-none focus:ring-2 focus:ring-ring"
                                >
                                    <img
                                        src={image.previewURL || "/placeholder.svg"}
                                        alt={`By ${image.author}`}
                                        className="w-full h-full object-cover"
                                        crossOrigin="anonymous"
                                    />
                                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors" />
                                    <div className="absolute bottom-0 left-0 right-0 p-1 bg-gradient-to-t from-black/70 to-transparent opacity-0 group-hover:opacity-100 transition-opacity">
                                        <p className="text-[10px] text-foreground truncate">{image.source}</p>
                                    </div>
                                    {/* Download icon on hover */}
                                    <div className="absolute top-1 right-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <span
                                            onClick={(e) => { e.stopPropagation(); onAddToList(image); }}
                                            className="inline-flex items-center justify-center w-6 h-6 bg-primary/80 text-primary-foreground rounded cursor-pointer hover:bg-primary"
                                        >
                                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
                                            </svg>
                                        </span>
                                    </div>
                                </button>
                            ))}
                        </div>
                    ) : (
                        <div className="flex items-center justify-center h-full">
                            <p className="text-sm text-muted-foreground">Nenhum resultado encontrado</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

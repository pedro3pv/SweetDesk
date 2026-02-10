'use client';

import { useState } from 'react';

interface SearchPanelProps {
    onImageSelect: (data: string) => void;
}

interface ImageResult {
    id: string;
    url: string;
    downloadURL: string;
    previewURL: string;
    width: number;
    height: number;
    author: string;
    source: string;
    tags: string[];
}

export default function SearchPanel({ onImageSelect }: SearchPanelProps) {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<ImageResult[]>([]);
    const [isSearching, setIsSearching] = useState(false);
    const [selectedImage, setSelectedImage] = useState<string | null>(null);

    const handleSearch = async () => {
        if (!query.trim()) return;

        setIsSearching(true);
        try {
            // @ts-ignore - Wails runtime will be injected
            const images = await window.go.main.App.SearchImages(query, 1, 12);
            setResults(images || []);
        } catch (error) {
            console.error('Search failed:', error);
            setResults([]);
        } finally {
            setIsSearching(false);
        }
    };

    const handleImageClick = async (image: ImageResult) => {
        setSelectedImage(image.id);
        try {
            // @ts-ignore - Wails runtime will be injected
            const base64Data = await window.go.main.App.DownloadImage(image.downloadURL);
            onImageSelect(base64Data);
        } catch (error) {
            console.error('Download failed:', error);
        } finally {
            setSelectedImage(null);
        }
    };

    return (
        <div className="space-y-4">
            {/* Search Input */}
            <div className="flex space-x-2">
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                    placeholder="Search wallpapers..."
                    className="flex-1 px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <button
                    onClick={handleSearch}
                    disabled={isSearching}
                    className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400 transition-colors"
                >
                    {isSearching ? '⏳' : '🔍'}
                </button>
            </div>

            {/* Search Results */}
            {results.length > 0 && (
                <div className="grid grid-cols-2 gap-2 max-h-96 overflow-y-auto">
                    {results.map((image) => (
                        <div
                            key={image.id}
                            onClick={() => handleImageClick(image)}
                            className={`relative cursor-pointer rounded-lg overflow-hidden border-2 transition-all ${
                                selectedImage === image.id
                                    ? 'border-blue-500 opacity-50'
                                    : 'border-transparent hover:border-blue-300'
                            }`}
                        >
                            <img
                                src={image.previewURL}
                                alt={`By ${image.author}`}
                                className="w-full h-24 object-cover"
                            />
                            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-2">
                                <p className="text-xs text-white truncate">
                                    {image.author}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Empty State */}
            {!isSearching && results.length === 0 && query && (
                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    <p>No results found for "{query}"</p>
                </div>
            )}

            {/* Initial State */}
            {!query && (
                <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                    <p className="text-sm">Try searching for:</p>
                    <div className="flex flex-wrap gap-2 mt-2 justify-center">
                        {['nature', 'mountains', 'ocean', 'city', 'space'].map((tag) => (
                            <button
                                key={tag}
                                onClick={() => {
                                    setQuery(tag);
                                    setTimeout(handleSearch, 100);
                                }}
                                className="px-3 py-1 bg-gray-100 dark:bg-gray-700 rounded-full text-xs hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                            >
                                {tag}
                            </button>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}

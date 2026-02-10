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
                    className="flex-1 px-3 lg:px-4 py-2 text-sm lg:text-base rounded-lg border-2 border-purple-300 dark:border-dark-border bg-white dark:bg-dark-surface text-purple-900 dark:text-dark-text placeholder-purple-400 dark:placeholder-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                />
                <button
                    onClick={handleSearch}
                    disabled={isSearching}
                    className="px-3 lg:px-4 py-2 bg-gradient-to-r from-purple-500 to-indigo-600 text-white rounded-lg hover:from-purple-600 hover:to-indigo-700 disabled:from-gray-400 disabled:to-gray-500 transition-all shadow-md hover:shadow-lg transform hover:scale-105 disabled:transform-none"
                >
                    {isSearching ? '⏳' : '🔍'}
                </button>
            </div>

            {/* Search Results */}
            {results.length > 0 && (
                <div className="grid grid-cols-2 gap-2 max-h-96 overflow-y-auto custom-scrollbar">
                    {results.map((image) => (
                        <div
                            key={image.id}
                            onClick={() => handleImageClick(image)}
                            className={`relative cursor-pointer rounded-lg overflow-hidden border-2 transition-all transform hover:scale-105 ${
                                selectedImage === image.id
                                    ? 'border-purple-500 opacity-50 shadow-lg'
                                    : 'border-transparent hover:border-purple-400 dark:hover:border-purple-600 shadow hover:shadow-md'
                            }`}
                        >
                            <img
                                src={image.previewURL}
                                alt={`By ${image.author}`}
                                className="w-full h-24 object-cover"
                            />
                            <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-purple-900/90 to-transparent p-2">
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
                <div className="text-center py-8 text-purple-600 dark:text-purple-400">
                    <p>No results found for "{query}"</p>
                </div>
            )}

            {/* Initial State */}
            {!query && (
                <div className="text-center py-8 text-purple-600 dark:text-purple-400">
                    <p className="text-sm mb-2 lg:mb-3">Try searching for:</p>
                    <div className="flex flex-wrap gap-2 mt-2 justify-center">
                        {['nature', 'mountains', 'ocean', 'city', 'space'].map((tag) => (
                            <button
                                key={tag}
                                onClick={() => {
                                    setQuery(tag);
                                    setTimeout(handleSearch, 100);
                                }}
                                className="px-3 py-1 bg-gradient-to-r from-purple-100 to-indigo-100 dark:from-purple-900/30 dark:to-indigo-900/30 text-purple-700 dark:text-purple-300 rounded-full text-xs hover:from-purple-200 hover:to-indigo-200 dark:hover:from-purple-800/40 dark:hover:to-indigo-800/40 transition-all shadow-sm hover:shadow-md transform hover:scale-105"
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

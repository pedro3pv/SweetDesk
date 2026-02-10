'use client';

import { useState } from 'react';

interface ImagePreviewProps {
    originalImage: string | null;
    processedImage: string | null;
}

export default function ImagePreview({ originalImage, processedImage }: ImagePreviewProps) {
    const [activeView, setActiveView] = useState<'original' | 'processed'>('original');

    const saveImage = async () => {
        if (!processedImage) return;

        try {
            // Create a blob from base64
            const response = await fetch(`data:image/png;base64,${processedImage}`);
            const blob = await response.blob();

            // Create download link
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `sweetdesk-${Date.now()}.png`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        } catch (error) {
            console.error('Failed to save image:', error);
        }
    };

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 h-full">
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                    Preview
                </h3>

                {/* View Toggle */}
                {processedImage && (
                    <div className="flex space-x-2">
                        <button
                            onClick={() => setActiveView('original')}
                            className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                                activeView === 'original'
                                    ? 'bg-blue-500 text-white'
                                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                            }`}
                        >
                            Original
                        </button>
                        <button
                            onClick={() => setActiveView('processed')}
                            className={`px-3 py-1 rounded-lg text-sm font-medium transition-colors ${
                                activeView === 'processed'
                                    ? 'bg-blue-500 text-white'
                                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                            }`}
                        >
                            Processed
                        </button>
                    </div>
                )}
            </div>

            {/* Image Display */}
            <div className="relative bg-gray-100 dark:bg-gray-900 rounded-lg overflow-hidden" style={{ minHeight: '400px' }}>
                {(originalImage || processedImage) ? (
                    <div className="flex items-center justify-center p-4">
                        <img
                            src={`data:image/png;base64,${
                                activeView === 'processed' && processedImage
                                    ? processedImage
                                    : originalImage
                            }`}
                            alt="Preview"
                            className="max-w-full max-h-[600px] object-contain rounded-lg shadow-lg"
                        />
                    </div>
                ) : (
                    <div className="flex items-center justify-center h-[400px]">
                        <div className="text-center">
                            <div className="text-6xl mb-4">🖼️</div>
                            <p className="text-gray-500 dark:text-gray-400">
                                No image selected
                            </p>
                            <p className="text-sm text-gray-400 dark:text-gray-500 mt-2">
                                Upload or search for an image to get started
                            </p>
                        </div>
                    </div>
                )}
            </div>

            {/* Action Buttons */}
            {processedImage && (
                <div className="mt-4 flex space-x-3">
                    <button
                        onClick={saveImage}
                        className="flex-1 py-2 bg-green-500 text-white rounded-lg font-medium hover:bg-green-600 transition-colors"
                    >
                        💾 Save Image
                    </button>
                    <button
                        className="flex-1 py-2 bg-blue-500 text-white rounded-lg font-medium hover:bg-blue-600 transition-colors"
                        onClick={() => {
                            // Future: Set as wallpaper functionality
                            alert('Set as wallpaper feature coming soon!');
                        }}
                    >
                        🖼️ Set as Wallpaper
                    </button>
                </div>
            )}

            {/* Image Info */}
            {(originalImage || processedImage) && (
                <div className="mt-4 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
                    <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                            <span className="text-gray-600 dark:text-gray-400">View:</span>
                            <span className="ml-2 font-medium text-gray-900 dark:text-white">
                                {activeView === 'processed' ? 'Processed' : 'Original'}
                            </span>
                        </div>
                        <div>
                            <span className="text-gray-600 dark:text-gray-400">Format:</span>
                            <span className="ml-2 font-medium text-gray-900 dark:text-white">
                                PNG
                            </span>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

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
        <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-xl shadow-lg border border-purple-200/50 dark:border-dark-border p-4 lg:p-6 h-full">
            <div className="flex justify-between items-center mb-4">
                <h3 className="text-base lg:text-lg font-semibold text-purple-900 dark:text-purple-200">
                    Preview
                </h3>

                {/* View Toggle */}
                {processedImage && (
                    <div className="flex space-x-2">
                        <button
                            onClick={() => setActiveView('original')}
                            className={`px-3 py-1 rounded-lg text-xs lg:text-sm font-medium transition-all transform ${
                                activeView === 'original'
                                    ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-md scale-105'
                                    : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30'
                            }`}
                        >
                            Original
                        </button>
                        <button
                            onClick={() => setActiveView('processed')}
                            className={`px-3 py-1 rounded-lg text-xs lg:text-sm font-medium transition-all transform ${
                                activeView === 'processed'
                                    ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-md scale-105'
                                    : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30'
                            }`}
                        >
                            Processed
                        </button>
                    </div>
                )}
            </div>

            {/* Image Display */}
            <div className="relative bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-950/30 dark:to-indigo-950/30 rounded-lg overflow-hidden border border-purple-200/50 dark:border-purple-800/30" style={{ minHeight: '400px' }}>
                {(originalImage || processedImage) ? (
                    <div className="flex items-center justify-center p-4">
                        <img
                            src={`data:image/png;base64,${
                                activeView === 'processed' && processedImage
                                    ? processedImage
                                    : originalImage
                            }`}
                            alt="Preview"
                            className="max-w-full max-h-[500px] lg:max-h-[600px] object-contain rounded-lg shadow-xl"
                        />
                    </div>
                ) : (
                    <div className="flex items-center justify-center h-[400px]">
                        <div className="text-center">
                            <div className="text-6xl mb-4">🖼️</div>
                            <p className="text-purple-600 dark:text-purple-400 font-medium">
                                No image selected
                            </p>
                            <p className="text-xs lg:text-sm text-purple-500 dark:text-purple-500 mt-2">
                                Upload or search for an image to get started
                            </p>
                        </div>
                    </div>
                )}
            </div>

            {/* Action Buttons */}
            {processedImage && (
                <div className="mt-4 flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-3">
                    <button
                        onClick={saveImage}
                        className="flex-1 py-2 bg-gradient-to-r from-cyan-500 to-cyan-600 text-white rounded-lg text-sm lg:text-base font-medium hover:from-cyan-600 hover:to-cyan-700 transition-all shadow-md hover:shadow-lg transform hover:scale-105"
                    >
                        💾 Save Image
                    </button>
                    <button
                        className="flex-1 py-2 bg-gradient-to-r from-amber-400 to-amber-500 text-white rounded-lg text-sm lg:text-base font-medium hover:from-amber-500 hover:to-amber-600 transition-all shadow-md hover:shadow-lg transform hover:scale-105"
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
                <div className="mt-4 p-3 bg-gradient-to-r from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
                    <div className="grid grid-cols-2 gap-2 text-xs lg:text-sm">
                        <div>
                            <span className="text-purple-600 dark:text-purple-400">View:</span>
                            <span className="ml-2 font-medium text-purple-900 dark:text-purple-200">
                                {activeView === 'processed' ? 'Processed' : 'Original'}
                            </span>
                        </div>
                        <div>
                            <span className="text-purple-600 dark:text-purple-400">Format:</span>
                            <span className="ml-2 font-medium text-purple-900 dark:text-purple-200">
                                PNG
                            </span>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

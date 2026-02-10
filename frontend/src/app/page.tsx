'use client';

import { useState, useCallback } from 'react';
import ImageUpload from '@/components/ImageUpload';
import ImagePreview from '@/components/ImagePreview';
import ProcessingPanel from '@/components/ProcessingPanel';
import SearchPanel from '@/components/SearchPanel';

export default function Home() {
    const [imageData, setImageData] = useState<string | null>(null);
    const [processedImage, setProcessedImage] = useState<string | null>(null);
    const [isProcessing, setIsProcessing] = useState(false);
    const [activeTab, setActiveTab] = useState<'upload' | 'search'>('search');

    const handleImageSelect = useCallback((data: string) => {
        setImageData(data);
        setProcessedImage(null);
    }, []);

    const handleProcessComplete = useCallback((result: string) => {
        setProcessedImage(result);
        setIsProcessing(false);
    }, []);

    return (
        <div className="min-h-screen bg-gradient-to-br from-purple-50 via-indigo-50 to-purple-100 dark:from-dark-bg dark:via-purple-950 dark:to-dark-bg">
            <div className="container mx-auto px-4 py-6 lg:py-8 max-w-7xl">
                {/* Header */}
                <header className="text-center mb-6 lg:mb-8">
                    <h1 className="text-3xl lg:text-4xl font-bold text-purple-900 dark:text-purple-200 mb-2">
                        🍬 SweetDesk
                    </h1>
                    <p className="text-sm lg:text-base text-purple-700 dark:text-purple-300">
                        AI-Powered Wallpaper Processing in 4K
                    </p>
                </header>

                {/* Main Content */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 lg:gap-6">
                    {/* Left Panel - Input */}
                    <div className="lg:col-span-1 space-y-4 lg:space-y-6">
                        <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-xl shadow-lg border border-purple-200/50 dark:border-dark-border p-4 lg:p-6">
                            {/* Tabs */}
                            <div className="flex space-x-2 mb-4 lg:mb-6">
                                <button
                                    onClick={() => setActiveTab('search')}
                                    className={`flex-1 py-2 px-3 lg:px-4 rounded-lg font-medium text-sm lg:text-base transition-all transform ${
                                        activeTab === 'search'
                                            ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-lg scale-105'
                                            : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30'
                                    }`}
                                >
                                    🔍 Search
                                </button>
                                <button
                                    onClick={() => setActiveTab('upload')}
                                    className={`flex-1 py-2 px-3 lg:px-4 rounded-lg font-medium text-sm lg:text-base transition-all transform ${
                                        activeTab === 'upload'
                                            ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-lg scale-105'
                                            : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30'
                                    }`}
                                >
                                    📁 Upload
                                </button>
                            </div>

                            {/* Content */}
                            {activeTab === 'search' ? (
                                <SearchPanel onImageSelect={handleImageSelect} />
                            ) : (
                                <ImageUpload onImageSelect={handleImageSelect} />
                            )}
                        </div>

                        {/* Processing Panel */}
                        {imageData && (
                            <ProcessingPanel
                                imageData={imageData}
                                isProcessing={isProcessing}
                                onProcessStart={() => setIsProcessing(true)}
                                onProcessComplete={handleProcessComplete}
                            />
                        )}
                    </div>

                    {/* Right Panel - Preview */}
                    <div className="lg:col-span-2">
                        <ImagePreview
                            originalImage={imageData}
                            processedImage={processedImage}
                        />
                    </div>
                </div>

                {/* Footer */}
                <footer className="mt-8 lg:mt-12 text-center text-xs lg:text-sm text-purple-600 dark:text-purple-400">
                    <p>Made with ❤️ by Molasses Co. | Cross-platform wallpaper processing</p>
                </footer>
            </div>
        </div>
    );
}

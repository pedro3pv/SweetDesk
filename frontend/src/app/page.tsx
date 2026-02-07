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
        <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
            <div className="container mx-auto px-4 py-8">
                {/* Header */}
                <header className="text-center mb-8">
                    <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-2">
                        🍬 SweetDesk
                    </h1>
                    <p className="text-gray-600 dark:text-gray-400">
                        AI-Powered Wallpaper Processing in 4K
                    </p>
                </header>

                {/* Main Content */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    {/* Left Panel - Input */}
                    <div className="lg:col-span-1">
                        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
                            {/* Tabs */}
                            <div className="flex space-x-2 mb-6">
                                <button
                                    onClick={() => setActiveTab('search')}
                                    className={`flex-1 py-2 px-4 rounded-lg font-medium transition-colors ${
                                        activeTab === 'search'
                                            ? 'bg-blue-500 text-white'
                                            : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                                    }`}
                                >
                                    🔍 Search
                                </button>
                                <button
                                    onClick={() => setActiveTab('upload')}
                                    className={`flex-1 py-2 px-4 rounded-lg font-medium transition-colors ${
                                        activeTab === 'upload'
                                            ? 'bg-blue-500 text-white'
                                            : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
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
                <footer className="mt-12 text-center text-sm text-gray-500 dark:text-gray-400">
                    <p>Made with ❤️ by Molasses Co. | Cross-platform wallpaper processing</p>
                </footer>
            </div>
        </div>
    );
}

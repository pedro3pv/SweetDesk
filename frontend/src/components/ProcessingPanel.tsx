'use client';

import { useState } from 'react';

interface ProcessingPanelProps {
    imageData: string;
    isProcessing: boolean;
    onProcessStart: () => void;
    onProcessComplete: (result: string) => void;
}

export default function ProcessingPanel({
    imageData,
    isProcessing,
    onProcessStart,
    onProcessComplete,
}: ProcessingPanelProps) {
    const [targetResolution, setTargetResolution] = useState('4K');
    const [useSeamCarving, setUseSeamCarving] = useState(false);
    const [progress, setProgress] = useState<string>('');

    const handleProcess = async () => {
        if (!imageData) return;

        onProcessStart();
        setProgress('Starting processing...');

        try {
            // Step 1: Classify
            setProgress('🔍 Classifying image type...');
            // @ts-ignore
            const imageType = await window.go.main.App.ClassifyImage(imageData);
            
            setProgress(`📊 Detected: ${imageType === 'anime' ? '🎨 Anime' : '📷 Photo'}`);
            await new Promise(resolve => setTimeout(resolve, 500));

            // Step 2: Upscale
            setProgress('🚀 Upscaling to ' + targetResolution + '...');
            const scale = targetResolution === '8K' ? 8 : targetResolution === '5K' ? 5 : 4;
            
            // @ts-ignore
            const upscaled = await window.go.main.App.UpscaleImage(imageData, imageType, scale);
            
            setProgress('✨ Adjusting aspect ratio...');
            await new Promise(resolve => setTimeout(resolve, 500));

            // Step 3: Full processing
            setProgress('🎨 Finalizing...');
            // @ts-ignore
            const result = await window.go.main.App.ProcessImage(
                imageData,
                targetResolution,
                useSeamCarving
            );

            setProgress('✅ Processing complete!');
            onProcessComplete(result);
            
            setTimeout(() => setProgress(''), 2000);
        } catch (error) {
            console.error('Processing failed:', error);
            setProgress('❌ Processing failed: ' + (error as Error).message);
            setTimeout(() => setProgress(''), 3000);
        }
    };

    return (
        <div className="mt-6 bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Processing Options
            </h3>

            {/* Resolution Selection */}
            <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Target Resolution
                </label>
                <div className="grid grid-cols-3 gap-2">
                    {['4K', '5K', '8K'].map((res) => (
                        <button
                            key={res}
                            onClick={() => setTargetResolution(res)}
                            className={`py-2 px-4 rounded-lg font-medium transition-colors ${
                                targetResolution === res
                                    ? 'bg-blue-500 text-white'
                                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                            }`}
                        >
                            {res}
                        </button>
                    ))}
                </div>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {targetResolution === '4K' && '3840 × 2160'}
                    {targetResolution === '5K' && '5120 × 2880'}
                    {targetResolution === '8K' && '7680 × 4320'}
                </p>
            </div>

            {/* Aspect Ratio Method */}
            <div className="mb-4">
                <label className="flex items-center space-x-2 cursor-pointer">
                    <input
                        type="checkbox"
                        checked={useSeamCarving}
                        onChange={(e) => setUseSeamCarving(e.target.checked)}
                        className="w-4 h-4 text-blue-500 rounded focus:ring-blue-500"
                    />
                    <span className="text-sm text-gray-700 dark:text-gray-300">
                        Use Content-Aware Resize (Seam Carving)
                    </span>
                </label>
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400 ml-6">
                    {useSeamCarving
                        ? 'Intelligently preserves important content (slower)'
                        : 'Fast center crop to target aspect ratio'}
                </p>
            </div>

            {/* Process Button */}
            <button
                onClick={handleProcess}
                disabled={isProcessing || !imageData}
                className="w-full py-3 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-lg font-semibold hover:from-blue-600 hover:to-purple-600 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed transition-all transform hover:scale-105"
            >
                {isProcessing ? '⏳ Processing...' : '🚀 Process Image'}
            </button>

            {/* Progress Display */}
            {progress && (
                <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
                    <p className="text-sm text-blue-700 dark:text-blue-300 text-center">
                        {progress}
                    </p>
                </div>
            )}
        </div>
    );
}

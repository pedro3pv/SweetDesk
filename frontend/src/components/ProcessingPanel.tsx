'use client';

import { useState } from 'react';
import { ClassifyImage, UpscaleImage, ProcessImage } from '@wailsjs/go/main/App';

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
            const imageType = await ClassifyImage(imageData);
            
            setProgress(`📊 Detected: ${imageType === 'anime' ? '🎨 Anime' : '📷 Photo'}`);
            await new Promise(resolve => setTimeout(resolve, 500));

            // Step 2: Upscale
            setProgress('🚀 Upscaling to ' + targetResolution + '...');
            const scale = targetResolution === '8K' ? 8 : targetResolution === '5K' ? 5 : 4;
            
            const upscaled = await UpscaleImage(imageData, imageType, scale);
            
            setProgress('✨ Adjusting aspect ratio...');
            await new Promise(resolve => setTimeout(resolve, 500));

            // Step 3: Full processing
            setProgress('🎨 Finalizing...');
            const result = await ProcessImage(
                upscaled,
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
        <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-xl shadow-lg border border-purple-200/50 dark:border-dark-border p-4 lg:p-6">
            <h3 className="text-base lg:text-lg font-semibold text-purple-900 dark:text-purple-200 mb-4">
                Processing Options
            </h3>

            {/* Resolution Selection */}
            <div className="mb-4">
                <label className="block text-xs lg:text-sm font-medium text-purple-800 dark:text-purple-300 mb-2">
                    Target Resolution
                </label>
                <div className="grid grid-cols-3 gap-2">
                    {['4K', '5K', '8K'].map((res) => (
                        <button
                            key={res}
                            onClick={() => setTargetResolution(res)}
                            className={`py-2 px-3 lg:px-4 rounded-lg text-sm lg:text-base font-medium transition-all transform ${
                                targetResolution === res
                                    ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-lg scale-105'
                                    : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30 hover:scale-105'
                            }`}
                        >
                            {res}
                        </button>
                    ))}
                </div>
                <p className="mt-1 text-xs text-purple-600 dark:text-purple-400">
                    {targetResolution === '4K' && '3840 × 2160'}
                    {targetResolution === '5K' && '5120 × 2880'}
                    {targetResolution === '8K' && '7680 × 4320'}
                </p>
            </div>

            {/* Aspect Ratio Method */}
            <div className="mb-4">
                <label className="flex items-center space-x-2 cursor-pointer group">
                    <input
                        type="checkbox"
                        checked={useSeamCarving}
                        onChange={(e) => setUseSeamCarving(e.target.checked)}
                        className="w-4 h-4 text-purple-500 rounded focus:ring-purple-500 border-purple-300 dark:border-purple-600"
                    />
                    <span className="text-xs lg:text-sm text-purple-800 dark:text-purple-300 group-hover:text-purple-600 dark:group-hover:text-purple-200 transition-colors">
                        Use Content-Aware Resize (Seam Carving)
                    </span>
                </label>
                <p className="mt-1 text-xs text-purple-600 dark:text-purple-400 ml-6">
                    {useSeamCarving
                        ? 'Intelligently preserves important content (slower)'
                        : 'Fast center crop to target aspect ratio'}
                </p>
            </div>

            {/* Process Button */}
            <button
                onClick={handleProcess}
                disabled={isProcessing || !imageData}
                className="w-full py-2.5 lg:py-3 bg-gradient-to-r from-purple-500 via-indigo-500 to-purple-600 text-white rounded-lg text-sm lg:text-base font-semibold hover:from-purple-600 hover:via-indigo-600 hover:to-purple-700 disabled:from-gray-400 disabled:to-gray-500 disabled:cursor-not-allowed transition-all transform hover:scale-105 disabled:transform-none shadow-lg hover:shadow-xl"
            >
                {isProcessing ? '⏳ Processing...' : '🚀 Process Image'}
            </button>

            {/* Progress Display */}
            {progress && (
                <div className="mt-4 p-3 bg-gradient-to-r from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
                    <p className="text-xs lg:text-sm text-purple-700 dark:text-purple-300 text-center font-medium">
                        {progress}
                    </p>
                </div>
            )}
        </div>
    );
}

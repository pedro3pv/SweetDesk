'use client';

import { useState } from 'react';
import { ProcessImage } from '../../wailsjs/go/main/App';

interface ProcessingPanelProps {
    imageData: string;
    isProcessing: boolean;
    onProcessStart: () => void;
    onProcessComplete: (result: string) => void;
    onProcessEnd?: () => void;
    onProcessError?: (error: Error) => void;
}

export default function ProcessingPanel({
    imageData,
    isProcessing,
    onProcessStart,
    onProcessComplete,
    onProcessEnd,
    onProcessError,
}: ProcessingPanelProps) {
    const [targetResolution, setTargetResolution] = useState('4K');
    const [progress, setProgress] = useState<string>('');

    const handleProcess = async () => {
        if (!imageData || isProcessing) return;

        onProcessStart();
        setProgress('Starting processing...');

        try {
            // Map resolution names to width/height
            let targetWidth = 3840;
            let targetHeight = 2160;
            switch (targetResolution) {
                case '4K': targetWidth = 3840; targetHeight = 2160; break;
                case '5K': targetWidth = 5120; targetHeight = 2880; break;
                case '8K': targetWidth = 7680; targetHeight = 4320; break;
            }

            setProgress('🔍 Analyzing image...');
            await new Promise(resolve => setTimeout(resolve, 300));

            setProgress('🚀 Processing to ' + targetResolution + '...');
            await new Promise(resolve => setTimeout(resolve, 300));

            // SweetDesk-core handles classification + upscaling in one call
            const result = await ProcessImage(
                imageData,
                targetWidth,
                targetHeight,
                '',
                ''
            );

            setProgress('✅ Processing complete!');
            onProcessComplete(result);
            
            setTimeout(() => {
                setProgress('');
                onProcessEnd?.();
            }, 2000);
        } catch (error) {
            console.error('Processing failed:', error);
            const errorObj = error as Error;
            setProgress('❌ Processing failed: ' + errorObj.message);
            onProcessError?.(errorObj);
            setTimeout(() => setProgress(''), 3000);
        }
    };

    return (
        <div className="bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm rounded-xl shadow-lg border border-gray-200 dark:border-gray-700 p-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
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

            {/* Process Button */}
            <button
                onClick={handleProcess}
                disabled={isProcessing || !imageData}
                className="w-full py-3 bg-blue-500 text-white rounded-lg font-semibold hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
            >
                {isProcessing ? '⏳ Processing...' : '🚀 Process Image'}
            </button>

            {/* Progress Display */}
            {progress && (
                <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
                    <p className="text-sm text-blue-700 dark:text-blue-300 text-center font-medium">
                        {progress}
                    </p>
                </div>
            )}
        </div>
    );
}

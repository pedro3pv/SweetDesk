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
            setProgress('🔍 Analyzing image...');
            await new Promise(resolve => setTimeout(resolve, 300));

            setProgress('🚀 Processing to ' + targetResolution + '...');
            await new Promise(resolve => setTimeout(resolve, 300));

            // Map resolution names to width/height
            let targetWidth = 3840;
            let targetHeight = 2160;
            switch (targetResolution) {
                case '4K': targetWidth = 3840; targetHeight = 2160; break;
                case '5K': targetWidth = 5120; targetHeight = 2880; break;
                case '8K': targetWidth = 7680; targetHeight = 4320; break;
            }

            // Full processing (classification + upscale + adjustments in backend)
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
        <div className="bg-card rounded-xl shadow-lg border border-border p-4 lg:p-6">
            <h3 className="text-base lg:text-lg font-semibold text-foreground mb-4">
                Processing Options
            </h3>

            {/* Resolution Selection */}
            <div className="mb-4">
                <label className="block text-xs lg:text-sm font-medium text-muted-foreground mb-2">
                    Target Resolution
                </label>
                <div className="grid grid-cols-3 gap-2">
                    {['4K', '5K', '8K'].map((res) => (
                        <button
                            key={res}
                            onClick={() => setTargetResolution(res)}
                            className={`py-2 px-3 lg:px-4 rounded-lg text-sm lg:text-base font-medium transition-colors ${
                                targetResolution === res
                                    ? 'bg-primary text-primary-foreground'
                                    : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
                            }`}
                        >
                            {res}
                        </button>
                    ))}
                </div>
                <p className="mt-1 text-xs text-muted-foreground">
                    {targetResolution === '4K' && '3840 × 2160'}
                    {targetResolution === '5K' && '5120 × 2880'}
                    {targetResolution === '8K' && '7680 × 4320'}
                </p>
            </div>

            {/* Process Button */}
            <button
                onClick={handleProcess}
                disabled={isProcessing || !imageData}
                className="w-full py-2.5 lg:py-3 bg-primary text-primary-foreground rounded-lg text-sm lg:text-base font-semibold hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
                {isProcessing ? '⏳ Processing...' : '🚀 Process Image'}
            </button>

            {/* Progress Display */}
            {progress && (
                <div className="mt-4 p-3 bg-primary/10 rounded-lg border border-border">
                    <p className="text-xs lg:text-sm text-foreground text-center font-medium">
                        {progress}
                    </p>
                </div>
            )}
        </div>
    );
}

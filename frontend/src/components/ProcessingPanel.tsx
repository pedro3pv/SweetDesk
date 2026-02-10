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
    const [useSeamCarving, setUseSeamCarving] = useState(false);
    const [progress, setProgress] = useState<string>('');
    const [useCustomResolution, setUseCustomResolution] = useState(false);
    const [customWidth, setCustomWidth] = useState(3840);
    const [customHeight, setCustomHeight] = useState(2160);
    const [aspectRatio, setAspectRatio] = useState('16:9');
    const [customAspectW, setCustomAspectW] = useState(16);
    const [customAspectH, setCustomAspectH] = useState(9);

    const handleProcess = async () => {
        if (!imageData) return;

        onProcessStart();
        setProgress('Starting processing...');

        try {
            let result: string;

            if (useCustomResolution) {
                // Use custom resolution processing
                setProgress('🚀 Processing with custom resolution...');
                
                let targetW = customWidth;
                let targetH = customHeight;
                let aspectW = 0;
                let aspectH = 0;

                // Parse aspect ratio if specified
                if (aspectRatio === 'custom') {
                    aspectW = customAspectW;
                    aspectH = customAspectH;
                } else if (aspectRatio !== 'auto') {
                    const parts = aspectRatio.split(':');
                    if (parts.length === 2) {
                        const w = parseInt(parts[0]);
                        const h = parseInt(parts[1]);
                        if (!isNaN(w) && !isNaN(h)) {
                            aspectW = w;
                            aspectH = h;
                        }
                    }
                }

                // @ts-ignore
                result = await window.go.main.App.ProcessImageWithCustomResolution(
                    imageData,
                    targetW,
                    targetH,
                    aspectW,
                    aspectH,
                    useSeamCarving
                );
            } else {
                // Use standard resolution processing
                setProgress('🔍 Analyzing image...');
                await new Promise(resolve => setTimeout(resolve, 300));

                setProgress('🚀 Processing to ' + targetResolution + '...');
                await new Promise(resolve => setTimeout(resolve, 300));

                // Full processing (classification + upscale + adjustments in backend)
                result = await ProcessImage(
                    imageData,
                    targetResolution,
                    useSeamCarving
                );
            }

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
        <div className="bg-white/90 dark:bg-dark-surface/90 backdrop-blur-sm rounded-xl shadow-lg border border-purple-200/50 dark:border-dark-border p-4 lg:p-6">
            <h3 className="text-base lg:text-lg font-semibold text-purple-900 dark:text-purple-200 mb-4">
                Processing Options
            </h3>

            {/* Resolution Mode Toggle */}
            <div className="mb-4">
                <label className="flex items-center space-x-2 cursor-pointer group">
                    <input
                        type="checkbox"
                        checked={useCustomResolution}
                        onChange={(e) => setUseCustomResolution(e.target.checked)}
                        className="w-4 h-4 text-purple-500 rounded focus:ring-purple-500 border-purple-300 dark:border-purple-600"
                    />
                    <span className="text-xs lg:text-sm font-medium text-purple-800 dark:text-purple-300 group-hover:text-purple-600 dark:group-hover:text-purple-200 transition-colors">
                        Custom Resolution
                    </span>
                </label>
                <p className="mt-1 text-xs text-purple-600 dark:text-purple-400 ml-6">
                    {useCustomResolution
                        ? 'Specify exact dimensions and aspect ratio'
                        : 'Use standard 4K/5K/8K presets'}
                </p>
            </div>

            {useCustomResolution ? (
                <>
                    {/* Custom Width/Height */}
                    <div className="mb-4 grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-xs lg:text-sm font-medium text-purple-800 dark:text-purple-300 mb-2">
                                Width (px)
                            </label>
                            <input
                                type="number"
                                value={customWidth}
                                onChange={(e) => setCustomWidth(Math.max(1, parseInt(e.target.value) || customWidth))}
                                className="w-full px-3 py-2 border border-purple-300 dark:border-purple-600 rounded-lg bg-white dark:bg-dark-surface text-purple-900 dark:text-purple-100 focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                min="1"
                                max="15360"
                            />
                        </div>
                        <div>
                            <label className="block text-xs lg:text-sm font-medium text-purple-800 dark:text-purple-300 mb-2">
                                Height (px)
                            </label>
                            <input
                                type="number"
                                value={customHeight}
                                onChange={(e) => setCustomHeight(Math.max(1, parseInt(e.target.value) || customHeight))}
                                className="w-full px-3 py-2 border border-purple-300 dark:border-purple-600 rounded-lg bg-white dark:bg-dark-surface text-purple-900 dark:text-purple-100 focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                                min="1"
                                max="8640"
                            />
                        </div>
                    </div>

                    {/* Aspect Ratio Selection */}
                    <div className="mb-4">
                        <label className="block text-xs lg:text-sm font-medium text-purple-800 dark:text-purple-300 mb-2">
                            Aspect Ratio
                        </label>
                        <div className="grid grid-cols-4 gap-2 mb-2">
                            {['16:9', '21:9', '4:3', 'custom'].map((ratio) => (
                                <button
                                    key={ratio}
                                    onClick={() => setAspectRatio(ratio)}
                                    className={`py-2 px-3 rounded-lg text-sm font-medium transition-all transform ${
                                        aspectRatio === ratio
                                            ? 'bg-gradient-to-r from-purple-500 to-indigo-600 text-white shadow-lg scale-105'
                                            : 'bg-purple-100 dark:bg-dark-surface text-purple-700 dark:text-purple-300 hover:bg-purple-200 dark:hover:bg-purple-900/30 hover:scale-105'
                                    }`}
                                >
                                    {ratio === 'custom' ? 'Custom' : ratio}
                                </button>
                            ))}
                        </div>

                        {aspectRatio === 'custom' && (
                            <div className="grid grid-cols-2 gap-4 mt-2">
                                <div>
                                    <label className="block text-xs text-purple-700 dark:text-purple-400 mb-1">
                                        Aspect W
                                    </label>
                                    <input
                                        type="number"
                                        value={customAspectW}
                                        onChange={(e) => setCustomAspectW(Math.max(1, parseInt(e.target.value) || customAspectW))}
                                        className="w-full px-3 py-2 border border-purple-300 dark:border-purple-600 rounded-lg bg-white dark:bg-dark-surface text-purple-900 dark:text-purple-100 focus:ring-2 focus:ring-purple-500 focus:border-transparent text-sm"
                                        min="1"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs text-purple-700 dark:text-purple-400 mb-1">
                                        Aspect H
                                    </label>
                                    <input
                                        type="number"
                                        value={customAspectH}
                                        onChange={(e) => setCustomAspectH(Math.max(1, parseInt(e.target.value) || customAspectH))}
                                        className="w-full px-3 py-2 border border-purple-300 dark:border-purple-600 rounded-lg bg-white dark:bg-dark-surface text-purple-900 dark:text-purple-100 focus:ring-2 focus:ring-purple-500 focus:border-transparent text-sm"
                                        min="1"
                                    />
                                </div>
                            </div>
                        )}
                    </div>
                </>
            ) : (
                <>
                    {/* Standard Resolution Selection */}
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
                </>
            )}

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

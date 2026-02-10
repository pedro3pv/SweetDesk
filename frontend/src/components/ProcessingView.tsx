'use client';

import { useState, useEffect, useCallback } from 'react';
import type { DownloadItem } from '../lib/types';

// Extend the Window interface to include 'go'
declare global {
    interface Window {
        go?: {
            main?: {
                App?: {
                    ProcessImage?: (base64Data: string, dimension: string, useSeamCarving: boolean) => Promise<void>;
                    DownloadImage?: (url: string) => Promise<string>;
                };
            };
        };
    }
}

interface ProcessingViewProps {
    items: DownloadItem[];
    onComplete: () => void;
}

export default function ProcessingView({ items, onComplete }: ProcessingViewProps) {
    const [progress, setProgress] = useState(0);
    const [currentItem, setCurrentItem] = useState(0);
    const [status, setStatus] = useState<'processing' | 'complete' | 'error'>('processing');
    const [itemStatuses, setItemStatuses] = useState<Record<string, 'pending' | 'processing' | 'done' | 'error'>>(
        Object.fromEntries(items.map(i => [i.id, 'pending' as const]))
    );

    const selectedItems = items.filter(i => i.selected);

    const processItems = useCallback(async () => {
        if (selectedItems.length === 0) {
            setProgress(100);
            setStatus('complete');
            return;
        }
        for (let i = 0; i < selectedItems.length; i++) {
            const item = selectedItems[i];
            setCurrentItem(i);
            setItemStatuses(prev => ({ ...prev, [item.id]: 'processing' }));

            try {
                // Try Wails backend first
                // @ts-ignore
                if (typeof window !== 'undefined' && window.go?.main?.App?.ProcessImage) {
                    let base64Data: string | null = null;

                    if (item.image.downloadURL.startsWith('data:')) {
                        // Strip the data URL prefix and keep only the base64 payload
                        base64Data = item.image.downloadURL.split(',')[1] || null;
                    } else if (window.go.main.App.DownloadImage) {
                        // @ts-ignore - Download and convert remote image URL to base64 via Wails binding
                        base64Data = await window.go.main.App.DownloadImage(item.image.downloadURL);
                    }

                    if (base64Data) {
                        const useSeamCarving = false;
                        // @ts-ignore
                        await window.go.main.App.ProcessImage(base64Data, item.dimension, useSeamCarving);
                    } else {
                        // Fallback: simulate processing if we cannot obtain base64 data
                        await new Promise(r => setTimeout(r, 800 + Math.random() * 1200));
                    }
                } else {
                    // Simulate processing
                    await new Promise(r => setTimeout(r, 800 + Math.random() * 1200));
                }

                setItemStatuses(prev => ({ ...prev, [item.id]: 'done' }));
            } catch {
                setItemStatuses(prev => ({ ...prev, [item.id]: 'error' }));
            }

            setProgress(Math.round(((i + 1) / selectedItems.length) * 100));
        }
        setStatus('complete');
    }, [selectedItems]);

    useEffect(() => {
        processItems();
    }, [processItems]);

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground tracking-tight">Lista de Download</h2>
                <span className="text-xs text-muted-foreground">
                    {status === 'complete' ? 'Concluido' : `${currentItem + 1}/${selectedItems.length}`}
                </span>
            </div>

            {/* Items with status */}
            <div className="flex-1 overflow-y-auto">
                <div className="divide-y divide-border">
                    {selectedItems.map(item => (
                        <div key={item.id} className="flex items-center gap-3 px-5 py-3">
                            {/* Status icon */}
                            <div className="w-5 h-5 shrink-0 flex items-center justify-center">
                                {itemStatuses[item.id] === 'done' && (
                                    <svg className="w-5 h-5 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                                    </svg>
                                )}
                                {itemStatuses[item.id] === 'processing' && (
                                    <svg className="w-5 h-5 text-primary animate-spin" fill="none" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
                                    </svg>
                                )}
                                {itemStatuses[item.id] === 'pending' && (
                                    <div className="w-3 h-3 rounded-full bg-border" />
                                )}
                                {itemStatuses[item.id] === 'error' && (
                                    <svg className="w-5 h-5 text-destructive" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                )}
                            </div>

                            {/* Thumbnail */}
                            <div className="w-10 h-10 shrink-0 rounded overflow-hidden bg-muted">
                                <img
                                    src={item.image.previewURL || "/placeholder.svg"}
                                    alt={item.name}
                                    className="w-full h-full object-cover"
                                    crossOrigin="anonymous"
                                />
                            </div>

                            {/* Info */}
                            <div className="flex-1 min-w-0">
                                <p className="text-sm text-foreground truncate">{item.name}</p>
                                <p className="text-[10px] text-muted-foreground">{item.dimension} / {item.aspect}</p>
                            </div>

                            {/* Upscale badge */}
                            {item.upscale && (
                                <span className="text-[10px] px-2 py-0.5 bg-primary/10 text-primary rounded-full shrink-0">upscale</span>
                            )}

                            {/* Check */}
                            {itemStatuses[item.id] === 'done' && (
                                <svg className="w-4 h-4 text-accent shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                                </svg>
                            )}
                        </div>
                    ))}
                </div>
            </div>

            {/* Progress Bar */}
            <div className="px-5 py-4 border-t border-border">
                <div className="flex items-center justify-between mb-2">
                    <p className="text-xs text-muted-foreground">
                        {status === 'complete' ? 'Processamento concluido' : 'Processando...'}
                    </p>
                    <p className="text-xs text-muted-foreground font-mono">{progress}%</p>
                </div>
                <div className="w-full h-2 bg-muted rounded-full overflow-hidden">
                    <div
                        className="h-full bg-primary rounded-full transition-all duration-500 ease-out"
                        style={{ width: `${progress}%` }}
                    />
                </div>

                {status === 'complete' && (
                    <button
                        onClick={onComplete}
                        className="w-full mt-4 py-3 bg-accent text-accent-foreground rounded-md text-sm font-medium hover:bg-accent/90 transition-colors"
                    >
                        Concluido
                    </button>
                )}
            </div>
        </div>
    );
}

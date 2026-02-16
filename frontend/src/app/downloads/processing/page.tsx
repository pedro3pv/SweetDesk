'use client';

import { useEffect, useRef, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../../../contexts/AppContext';
import { useProcessing } from '../../../hooks/useProcessing';

export default function ProcessingPage() {
    const router = useRouter();
    const { items, savePath } = useDownloads();
    const { progress, currentItem, status, itemStatuses, startBatch } = useProcessing();
    const hasStartedRef = useRef(false);

    const selectedItems = items.filter(i => i.selected);

    // Start processing on mount
    useEffect(() => {
        if (hasStartedRef.current) return;
        if (selectedItems.length === 0 || !savePath) {
            router.push('/downloads');
            return;
        }
        hasStartedRef.current = true;
        startBatch(items, savePath);
    }, [items, savePath, selectedItems.length, startBatch, router]);

    const handleComplete = useCallback(() => {
        router.push('/downloads/complete');
    }, [router]);

    return (
        <div className="flex flex-col h-full flex-1">
            {/* Header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground tracking-tight">Processando</h2>
                <span className="text-xs text-muted-foreground">
                    {status === 'complete' ? 'Concluído' : `${currentItem + 1}/${selectedItems.length}`}
                </span>
            </div>

            {/* Items with status */}
            <div className="flex-1 overflow-y-auto">
                <div className="divide-y divide-border">
                    {selectedItems.map(item => {
                        const itemStatus = itemStatuses[item.id] || 'pending';
                        return (
                            <div key={item.id} className="flex items-center gap-3 px-5 py-3">
                                {/* Status icon */}
                                <div className="w-5 h-5 shrink-0 flex items-center justify-center">
                                    {itemStatus === 'done' && (
                                        <svg className="w-5 h-5 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                                        </svg>
                                    )}
                                    {itemStatus === 'processing' && (
                                        <svg className="w-5 h-5 text-primary animate-spin" fill="none" viewBox="0 0 24 24">
                                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
                                        </svg>
                                    )}
                                    {itemStatus === 'pending' && (
                                        <div className="w-3 h-3 rounded-full bg-border" />
                                    )}
                                    {itemStatus === 'error' && (
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

                                {/* Check for done */}
                                {itemStatus === 'done' && (
                                    <svg className="w-4 h-4 text-accent shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                                    </svg>
                                )}
                            </div>
                        );
                    })}
                </div>
            </div>

            {/* Progress Bar */}
            <div className="px-5 py-4 border-t border-border">
                <div className="flex items-center justify-between mb-2">
                    <p className="text-xs text-muted-foreground">
                        {status === 'complete' ? 'Processamento concluído' : 'Processando...'}
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
                        onClick={handleComplete}
                        className="w-full mt-4 py-3 bg-accent text-accent-foreground rounded-md text-sm font-medium hover:bg-accent/90 transition-colors"
                    >
                        Concluído
                    </button>
                )}
            </div>
        </div>
    );
}

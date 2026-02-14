'use client';

import { useState, useEffect, useRef } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { DownloadItem } from '../lib/types';

// --- Backend types (mirrors Go structs) ---
interface BatchItem {
    id: string;
    base64Data: string;
    downloadURL: string;
    name: string;
    dimension: string;
}

interface BatchItemStatus {
    id: string;
    status: 'pending' | 'processing' | 'done' | 'error';
    error?: string;
}

interface ProcessingStatus {
    isProcessing: boolean;
    total: number;
    current: number;
    progress: number;
    items: BatchItemStatus[];
    done: boolean;
}

interface ProcessingViewProps {
    items: DownloadItem[];
    savePath: string;
    onComplete: () => void;
}

export default function ProcessingView({ items, savePath, onComplete }: ProcessingViewProps) {
    const [progress, setProgress] = useState(0);
    const [currentItem, setCurrentItem] = useState(0);
    const [status, setStatus] = useState<'processing' | 'complete' | 'error'>('processing');
    const [itemStatuses, setItemStatuses] = useState<Record<string, 'pending' | 'processing' | 'done' | 'error'>>(
        () => Object.fromEntries(items.filter(i => i.selected).map(i => [i.id, 'pending' as const]))
    );

    const selectedItems = items.filter(i => i.selected);
    const hasStartedRef = useRef(false);

    // Apply a ProcessingStatus update from the backend
    function applyStatus(ps: ProcessingStatus) {
        if (!ps || !ps.items) return;
        setProgress(ps.progress);
        setCurrentItem(ps.current);
        const newStatuses: Record<string, 'pending' | 'processing' | 'done' | 'error'> = {};
        for (const item of ps.items) {
            newStatuses[item.id] = item.status as 'pending' | 'processing' | 'done' | 'error';
        }
        setItemStatuses(newStatuses);
        if (ps.done) {
            setStatus('complete');
        }
    }

    // On mount: recover existing progress OR start new batch
    useEffect(() => {
        async function recoverOrStart() {
            try {
                if (window.go?.main?.App?.GetProcessingStatus) {
                    const ps = await window.go.main.App.GetProcessingStatus();

                    if (ps && ps.items && ps.items.length > 0) {
                        // Compare backend item IDs with our current items
                        // to distinguish HMR re-mount from a brand-new batch
                        const backendIds = new Set(ps.items.map((i: BatchItemStatus) => i.id));
                        const currentIds = new Set(items.filter(i => i.selected).map(i => i.id));
                        const sameItems =
                            backendIds.size === currentIds.size &&
                            [...backendIds].every(id => currentIds.has(id));

                        if (sameItems && ps.isProcessing) {
                            // HMR during active processing — just recover state
                            applyStatus(ps);
                            hasStartedRef.current = true;
                            return;
                        }
                        if (sameItems && ps.done) {
                            // HMR after batch already completed — show finished state
                            applyStatus(ps);
                            hasStartedRef.current = true;
                            return;
                        }
                        // Items differ → fall through and start a new batch
                    }
                }
            } catch {
                // ignore — will start fresh
            }

            if (!hasStartedRef.current) {
                hasStartedRef.current = true;
                startBatch();
            }
        }

        function startBatch() {
            const selected = items.filter(i => i.selected);
            if (selected.length === 0) {
                setProgress(100);
                setStatus('complete');
                return;
            }

            // Build batch items for the backend
            const batchItems: BatchItem[] = selected.map(item => {
                let base64Data = '';
                let downloadURL = item.image.downloadURL || '';

                if (downloadURL.startsWith('data:')) {
                    base64Data = downloadURL.split(',')[1] || '';
                    downloadURL = '';
                }

                // Sanitize filename
                let name = (item.name || `wallpaper-${item.id}`).replace(/^.*[\/]/, '');
                name = name.replace(/[<>:"/\\|?*\x00-\x1F]/g, '_');
                name = name.replace(/^[.\s]+|[.\s]+$/g, '');
                if (!/\.[a-zA-Z0-9]{1,6}$/.test(name)) name += '.png';
                if (!name) name = `wallpaper_${Date.now()}.png`;

                return {
                    id: item.id,
                    base64Data,
                    downloadURL,
                    name,
                    dimension: item.dimension || '3840x2160',
                };
            });

            // Fire and forget — backend processes in a goroutine
            if (window.go?.main?.App?.ProcessBatch) {
                window.go.main.App.ProcessBatch(batchItems, savePath);
            }
        }

        recoverOrStart();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    // Listen for real-time progress events from backend
    useEffect(() => {
        function handleEvent(...args: unknown[]) {
            const ps = args[0] as ProcessingStatus;
            applyStatus(ps);
        }

        // EventsOn returns a cancel function
        const cancel = EventsOn('processing:status', handleEvent);
        return () => { cancel(); };
    }, []);

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

'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import type { DownloadItem } from '../lib/types';

// ─── Backend types (mirrors Go structs) ─────────────────────────────

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

export type ItemStatusMap = Record<string, 'pending' | 'processing' | 'done' | 'error'>;

export interface UseProcessingResult {
    progress: number;
    currentItem: number;
    status: 'idle' | 'processing' | 'complete' | 'error';
    itemStatuses: ItemStatusMap;
    startBatch: (items: DownloadItem[], savePath: string) => void;
    reset: () => void;
}

function clamp(value: number, min: number, max: number): number {
    return Math.max(min, Math.min(max, value));
}

function isValidStatus(s: string): s is 'pending' | 'processing' | 'done' | 'error' {
    return s === 'pending' || s === 'processing' || s === 'done' || s === 'error';
}

export function useProcessing(): UseProcessingResult {
    const [progress, setProgress] = useState(0);
    const [currentItem, setCurrentItem] = useState(0);
    const [status, setStatus] = useState<'idle' | 'processing' | 'complete' | 'error'>('idle');
    const [itemStatuses, setItemStatuses] = useState<ItemStatusMap>({});
    const hasStartedRef = useRef(false);

    // Apply a ProcessingStatus update from the backend — with validation
    const applyStatus = useCallback((ps: ProcessingStatus) => {
        if (!ps || !Array.isArray(ps.items)) return;

        const validProgress = clamp(
            typeof ps.progress === 'number' ? ps.progress : 0,
            0,
            100
        );
        const validCurrent = clamp(
            typeof ps.current === 'number' ? ps.current : 0,
            0,
            Math.max(0, ps.total - 1)
        );

        setProgress(validProgress);
        setCurrentItem(validCurrent);

        const newStatuses: ItemStatusMap = {};
        for (const item of ps.items) {
            if (item && item.id && isValidStatus(item.status)) {
                newStatuses[item.id] = item.status;
            }
        }
        setItemStatuses(newStatuses);

        if (ps.done) {
            setStatus('complete');
        } else if (ps.isProcessing) {
            setStatus('processing');
        }
    }, []);

    // Listen for progress events from backend
    useEffect(() => {
        const cancel = EventsOn('processing:status', (...args: unknown[]) => {
            const ps = args[0] as ProcessingStatus;
            applyStatus(ps);
        });
        return () => { cancel(); };
    }, [applyStatus]);

    // Start a new batch
    const startBatch = useCallback((items: DownloadItem[], savePath: string) => {
        if (hasStartedRef.current) return;
        hasStartedRef.current = true;

        const selected = items.filter(i => i.selected);
        if (selected.length === 0) {
            setProgress(100);
            setStatus('complete');
            return;
        }

        // Initialize item statuses
        setItemStatuses(
            Object.fromEntries(selected.map(i => [i.id, 'pending' as const]))
        );
        setProgress(0);
        setCurrentItem(0);
        setStatus('processing');

        // Build batch items for the backend
        const batchItems: BatchItem[] = selected.map(item => {
            let base64Data = '';
            let downloadURL = item.image.downloadURL || '';

            if (downloadURL.startsWith('data:')) {
                base64Data = downloadURL.split(',')[1] || '';
                downloadURL = '';
            }

            // Sanitize filename
            let name = (item.name || `wallpaper-${item.id}`).replace(/^.*[/]/, '');
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
    }, []);

    // Try to recover existing progress on mount
    const recoverState = useCallback(async (items: DownloadItem[]) => {
        try {
            if (window.go?.main?.App?.GetProcessingStatus) {
                const ps = await window.go.main.App.GetProcessingStatus();
                if (ps && ps.items && ps.items.length > 0) {
                    const backendIds = new Set(ps.items.map((i: BatchItemStatus) => i.id));
                    const currentIds = new Set(items.filter(i => i.selected).map(i => i.id));
                    const sameItems =
                        backendIds.size === currentIds.size &&
                        [...backendIds].every(id => currentIds.has(id));

                    if (sameItems && (ps.isProcessing || ps.done)) {
                        applyStatus(ps);
                        hasStartedRef.current = true;
                        return true;
                    }
                }
            }
        } catch {
            // ignore
        }
        return false;
    }, [applyStatus]);

    const reset = useCallback(() => {
        setProgress(0);
        setCurrentItem(0);
        setStatus('idle');
        setItemStatuses({});
        hasStartedRef.current = false;
    }, []);

    // Expose recoverState via startBatch wrapping
    const startBatchWithRecovery = useCallback(async (items: DownloadItem[], savePath: string) => {
        const recovered = await recoverState(items);
        if (!recovered) {
            startBatch(items, savePath);
        }
    }, [recoverState, startBatch]);

    return {
        progress,
        currentItem,
        status,
        itemStatuses,
        startBatch: startBatchWithRecovery as unknown as (items: DownloadItem[], savePath: string) => void,
        reset,
    };
}

'use client';

import { useState } from 'react';
import type { DownloadItem } from '../lib/types';

interface DownloadListProps {
    items: DownloadItem[];
    onUpdateItem: (id: string, updates: Partial<DownloadItem>) => void;
    onRemoveItem: (id: string) => void;
    onProcess: () => void;
    title?: string;
    processing?: boolean;
}

const DIMENSIONS = ['1920x1080', '2560x1440', '3840x2160', '5120x2880', 'custom'];
const ASPECTS = ['16:9', '4:3', '21:9', '1:1', 'custom'];

export default function DownloadList({
    items,
    onUpdateItem,
    onRemoveItem,
    onProcess,
    title = 'Lista de Download',
    processing = false,
}: DownloadListProps) {
    const selectedCount = items.filter(i => i.selected).length;
    const [customDimensions, setCustomDimensions] = useState<Record<string, string>>({});
    const [customAspects, setCustomAspects] = useState<Record<string, string>>({});

    const handleDimensionChange = (id: string, value: string) => {
        if (value === 'custom') {
            setCustomDimensions(prev => ({ ...prev, [id]: '' }));
        } else {
            setCustomDimensions(prev => { const n = { ...prev }; delete n[id]; return n; });
            onUpdateItem(id, { dimension: value });
        }
    };

    const handleCustomDimension = (id: string, value: string) => {
        setCustomDimensions(prev => ({ ...prev, [id]: value }));
        // Validate format WxH with reasonable bounds
        if (/^\d+x\d+$/.test(value)) {
            const [w, h] = value.split('x').map(Number);
            if (w >= 1 && h >= 1 && w <= 16384 && h <= 16384) {
                onUpdateItem(id, { dimension: value });

                // Auto-clear custom state if value matches a preset
                if (DIMENSIONS.slice(0, -1).includes(value)) {
                    setCustomDimensions(prev => {
                        const n = { ...prev };
                        delete n[id];
                        return n;
                    });
                }
            }
        }
    };

    const handleAspectChange = (id: string, value: string) => {
        if (value === 'custom') {
            setCustomAspects(prev => ({ ...prev, [id]: '' }));
        } else {
            setCustomAspects(prev => { const n = { ...prev }; delete n[id]; return n; });
            onUpdateItem(id, { aspect: value });
        }
    };

    const handleCustomAspect = (id: string, value: string) => {
        setCustomAspects(prev => ({ ...prev, [id]: value }));
        // Validate format W:H with reasonable bounds
        if (/^\d+:\d+$/.test(value)) {
            const [w, h] = value.split(':').map(Number);
            if (w >= 1 && h >= 1) {
                onUpdateItem(id, { aspect: value });

                // Auto-clear custom state if value matches a preset
                if (ASPECTS.slice(0, -1).includes(value)) {
                    setCustomAspects(prev => {
                        const n = { ...prev };
                        delete n[id];
                        return n;
                    });
                }
            }
        }
    };

    const isCustomDimension = (item: DownloadItem) => {
        return item.id in customDimensions || !DIMENSIONS.slice(0, -1).includes(item.dimension);
    };

    const isCustomAspect = (item: DownloadItem) => {
        return item.id in customAspects || !ASPECTS.slice(0, -1).includes(item.aspect);
    };

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-center justify-between px-5 py-4 border-b border-border">
                <h2 className="text-base font-semibold text-foreground tracking-tight">{title}</h2>
                <span className="text-xs text-muted-foreground">{selectedCount}/{items.length} selecionados</span>
            </div>

            {/* List */}
            <div className="flex-1 overflow-y-auto">
                {items.length === 0 ? (
                    <div className="flex items-center justify-center h-full">
                        <p className="text-sm text-muted-foreground">Nenhum item na lista</p>
                    </div>
                ) : (
                    <div className="divide-y divide-border">
                        {items.map(item => (
                            <div key={item.id} className="flex items-center gap-3 px-5 py-3 hover:bg-muted/50 transition-colors">
                                {/* Selection checkbox */}
                                <input
                                    type="checkbox"
                                    checked={item.selected}
                                    onChange={(e) => onUpdateItem(item.id, { selected: e.target.checked })}
                                    className="w-4 h-4 shrink-0 rounded border-border text-primary focus:ring-primary bg-input accent-primary"
                                />

                                {/* Image thumbnail */}
                                <div className="w-10 h-10 shrink-0 rounded overflow-hidden bg-muted">
                                    <img
                                        src={item.image.previewURL || "/placeholder.svg"}
                                        alt={item.name}
                                        className="w-full h-full object-cover"
                                        crossOrigin="anonymous"
                                    />
                                </div>

                                {/* Name (editable) */}
                                <div className="flex-1 min-w-0">
                                    <input
                                        type="text"
                                        value={item.name}
                                        onChange={(e) => onUpdateItem(item.id, { name: e.target.value })}
                                        className="w-full bg-transparent text-sm text-foreground focus:outline-none focus:bg-input focus:px-2 focus:rounded transition-all truncate"
                                        title="Clique para editar o nome"
                                    />
                                    <p className="text-[10px] text-muted-foreground">(pode editar clicando em cima)</p>
                                </div>

                                {/* Upscale */}
                                <label className="flex items-center gap-1.5 shrink-0 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={item.upscale}
                                        onChange={(e) => onUpdateItem(item.id, { upscale: e.target.checked })}
                                        className="w-3.5 h-3.5 rounded border-border text-primary focus:ring-primary bg-input accent-primary"
                                    />
                                    <span className="text-[10px] text-muted-foreground font-medium">upscale</span>
                                </label>

                                {/* Dimension select */}
                                {isCustomDimension(item) ? (
                                    <div className="flex items-center gap-1">
                                        <input
                                            type="text"
                                            value={customDimensions[item.id] ?? item.dimension}
                                            onChange={(e) => handleCustomDimension(item.id, e.target.value)}
                                            placeholder="WxH"
                                            className="w-20 text-[10px] bg-secondary text-secondary-foreground border border-border rounded px-1.5 py-1 focus:outline-none focus:ring-1 focus:ring-ring shrink-0"
                                        />
                                        <button
                                            onClick={() => handleDimensionChange(item.id, DIMENSIONS[0])}
                                            className="text-[10px] text-muted-foreground hover:text-foreground p-1"
                                            title="Voltar para predefinições"
                                            aria-label="Voltar para predefinições"
                                        >
                                            ↩
                                        </button>
                                    </div>
                                ) : (
                                    <select
                                        value={item.dimension}
                                        onChange={(e) => handleDimensionChange(item.id, e.target.value)}
                                        className="text-[10px] bg-secondary text-secondary-foreground border border-border rounded px-1.5 py-1 focus:outline-none focus:ring-1 focus:ring-ring shrink-0"
                                    >
                                        {DIMENSIONS.map(d => (
                                            <option key={d} value={d}>{d === 'custom' ? 'Personalizado' : d}</option>
                                        ))}
                                    </select>
                                )}

                                {/* Aspect select */}
                                {isCustomAspect(item) ? (
                                    <div className="flex items-center gap-1">
                                        <input
                                            type="text"
                                            value={customAspects[item.id] ?? item.aspect}
                                            onChange={(e) => handleCustomAspect(item.id, e.target.value)}
                                            placeholder="W:H"
                                            className="w-14 text-[10px] bg-secondary text-secondary-foreground border border-border rounded px-1.5 py-1 focus:outline-none focus:ring-1 focus:ring-ring shrink-0"
                                        />
                                        <button
                                            onClick={() => handleAspectChange(item.id, ASPECTS[0])}
                                            className="text-[10px] text-muted-foreground hover:text-foreground p-1"
                                            title="Voltar para predefinições"
                                            aria-label="Voltar para predefinições"
                                        >
                                            ↩
                                        </button>
                                    </div>
                                ) : (
                                    <select
                                        value={item.aspect}
                                        onChange={(e) => handleAspectChange(item.id, e.target.value)}
                                        className="text-[10px] bg-secondary text-secondary-foreground border border-border rounded px-1.5 py-1 focus:outline-none focus:ring-1 focus:ring-ring shrink-0"
                                    >
                                        {ASPECTS.map(a => (
                                            <option key={a} value={a}>{a === 'custom' ? 'Personalizado' : a}</option>
                                        ))}
                                    </select>
                                )}

                                {/* Remove */}
                                <button
                                    onClick={() => onRemoveItem(item.id)}
                                    className="p-1 text-muted-foreground hover:text-destructive transition-colors shrink-0"
                                    aria-label="Remover item"
                                >
                                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                                    </svg>
                                </button>
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Process Button */}
            {items.length > 0 && (
                <div className="px-5 py-4 border-t border-border">
                    <button
                        onClick={onProcess}
                        disabled={processing || selectedCount === 0}
                        className="w-full py-3 bg-primary text-primary-foreground rounded-md text-sm font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        {processing ? 'Processando...' : 'Processar e Salvar'}
                    </button>
                </div>
            )}
        </div>
    );
}

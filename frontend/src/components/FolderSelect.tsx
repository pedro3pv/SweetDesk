'use client';

import { useState } from 'react';

interface FolderSelectProps {
    onSelect: (path: string) => void;
    onCancel: () => void;
}

export default function FolderSelect({ onSelect, onCancel }: FolderSelectProps) {
    const [selectedPath, setSelectedPath] = useState('');
    const [customPath, setCustomPath] = useState('');

    const commonPaths = [
        { label: 'Desktop', path: '~/Desktop' },
        { label: 'Downloads', path: '~/Downloads' },
        { label: 'Imagens', path: '~/Pictures' },
        { label: 'Wallpapers', path: '~/Pictures/Wallpapers' },
        { label: 'Documentos', path: '~/Documents' },
    ];

    const handleConfirm = () => {
        const path = customPath || selectedPath;
        if (path) {
            onSelect(path);
        }
    };

    return (
        <div className="flex flex-col h-full">
            {/* Header */}
            <div className="flex items-center gap-3 px-5 py-4 border-b border-border">
                <button
                    onClick={onCancel}
                    className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
                    </svg>
                    Voltar
                </button>
            </div>

            <div className="flex-1 flex flex-col items-center justify-center p-8 gap-6">
                <div className="w-full max-w-md">
                    <h2 className="text-lg font-semibold text-foreground text-center mb-6">Selecionar Pasta para Salvar</h2>

                    {/* Folder icon */}
                    <div className="flex justify-center mb-6">
                        <svg className="w-16 h-16 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        </svg>
                    </div>

                    {/* Common paths */}
                    <div className="flex flex-col gap-2 mb-6">
                        {commonPaths.map(p => (
                            <button
                                key={p.path}
                                onClick={() => { setSelectedPath(p.path); setCustomPath(''); }}
                                className={`flex items-center gap-3 px-4 py-3 rounded-md border text-sm text-left transition-colors ${
                                    selectedPath === p.path && !customPath
                                        ? 'bg-primary/10 border-primary text-foreground'
                                        : 'bg-card border-border text-foreground hover:bg-muted'
                                }`}
                            >
                                <svg className="w-4 h-4 text-muted-foreground shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                                </svg>
                                <div className="flex-1 min-w-0">
                                    <p className="font-medium">{p.label}</p>
                                    <p className="text-xs text-muted-foreground truncate">{p.path}</p>
                                </div>
                            </button>
                        ))}
                    </div>

                    {/* Custom path */}
                    <div className="mb-6">
                        <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-2 block">
                            Ou digite um caminho personalizado
                        </label>
                        <input
                            type="text"
                            value={customPath}
                            onChange={(e) => { setCustomPath(e.target.value); setSelectedPath(''); }}
                            placeholder="/caminho/para/pasta"
                            className="w-full px-3 py-2 bg-input border border-border rounded-md text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                        />
                    </div>

                    {/* Actions */}
                    <div className="flex gap-3">
                        <button
                            onClick={onCancel}
                            className="flex-1 py-3 bg-secondary text-secondary-foreground rounded-md text-sm font-medium hover:bg-secondary/80 transition-colors"
                        >
                            Cancelar
                        </button>
                        <button
                            onClick={handleConfirm}
                            disabled={!selectedPath && !customPath}
                            className="flex-1 py-3 bg-primary text-primary-foreground rounded-md text-sm font-medium hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                        >
                            Confirmar
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

'use client';

import { useState, useEffect } from 'react';

interface FolderSelectProps {
    onSelect: (path: string) => void;
    onCancel: () => void;
}

export default function FolderSelect({ onSelect, onCancel }: FolderSelectProps) {
    const [selectedPath, setSelectedPath] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    // On mount, try to get the default save path from backend
    useEffect(() => {
        async function loadDefault() {
            try {
                if (window.go?.main?.App?.GetDefaultSavePath) {
                    const defaultPath = await window.go.main.App.GetDefaultSavePath();
                    if (defaultPath) setSelectedPath(defaultPath);
                }
            } catch {
                // ignore - will show empty
            }
        }
        loadDefault();
    }, []);

    const handleBrowse = async () => {
        setIsLoading(true);
        try {
            if (window.go?.main?.App?.SelectDirectory) {
                const result = await window.go.main.App.SelectDirectory();
                if (result) {
                    setSelectedPath(result);
                }
            } else {
                // Fallback: prompt for path if native dialog not available
                const path = prompt('Digite o caminho da pasta:', selectedPath);
                if (path) setSelectedPath(path);
            }
        } catch (err) {
            console.error('Failed to open directory dialog:', err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleConfirm = () => {
        if (selectedPath) {
            onSelect(selectedPath);
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

                    {/* Selected path display */}
                    <div className="mb-4">
                        <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-2 block">
                            Pasta selecionada
                        </label>
                        <div className="w-full px-3 py-3 bg-input border border-border rounded-md text-sm text-foreground min-h-[42px] flex items-center">
                            {selectedPath ? (
                                <span className="truncate">{selectedPath}</span>
                            ) : (
                                <span className="text-muted-foreground">Nenhuma pasta selecionada</span>
                            )}
                        </div>
                    </div>

                    {/* Browse button - opens native OS dialog */}
                    <button
                        onClick={handleBrowse}
                        disabled={isLoading}
                        className="w-full mb-6 py-3 bg-secondary text-secondary-foreground rounded-md text-sm font-medium hover:bg-secondary/80 disabled:opacity-50 transition-colors flex items-center justify-center gap-2"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        </svg>
                        {isLoading ? 'Abrindo...' : 'Escolher Pasta'}
                    </button>

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
                            disabled={!selectedPath}
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

'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../../../contexts/AppContext';

export default function CompletePage() {
    const router = useRouter();
    const { items, clearAll } = useDownloads();
    const itemCount = items.filter(i => i.selected).length;

    const handleBackToList = useCallback(() => {
        clearAll();
        router.push('/');
    }, [clearAll, router]);

    return (
        <div className="flex flex-col h-full flex-1 items-center justify-center p-8">
            <div className="flex flex-col items-center gap-6 max-w-sm text-center">
                {/* Success icon */}
                <div className="w-20 h-20 rounded-full bg-accent/10 flex items-center justify-center">
                    <svg className="w-10 h-10 text-accent" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                    </svg>
                </div>

                <div>
                    <h2 className="text-2xl font-bold text-foreground mb-2">Processo Finalizado</h2>
                    <p className="text-sm text-muted-foreground">
                        {itemCount} {itemCount === 1 ? 'imagem foi processada' : 'imagens foram processadas'} e salvas com sucesso.
                    </p>
                </div>

                <button
                    onClick={handleBackToList}
                    className="px-6 py-3 bg-secondary text-secondary-foreground rounded-md text-sm font-medium hover:bg-secondary/80 transition-colors border border-border"
                >
                    Voltar a lista de fotos
                </button>
            </div>
        </div>
    );
}

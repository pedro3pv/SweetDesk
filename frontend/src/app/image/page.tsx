'use client';

import { useCallback, useMemo, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useDownloads } from '../../contexts/AppContext';
import ImageDetail from '../../components/ImageDetail';
import type { ImageResult, DownloadItem } from '../../lib/types';

function ImagePageContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const { addDownloadItem } = useDownloads();

    // Reconstruct image from URL params
    const image = useMemo((): ImageResult | null => {
        const id = searchParams.get('id');
        if (!id) return null;
        return {
            id,
            url: searchParams.get('url') || '',
            downloadURL: searchParams.get('downloadURL') || '',
            previewURL: searchParams.get('url') || '',
            width: Number(searchParams.get('width')) || 0,
            height: Number(searchParams.get('height')) || 0,
            author: searchParams.get('author') || '',
            source: searchParams.get('source') || '',
            tags: (searchParams.get('tags') || '').split(',').filter(Boolean),
            description: searchParams.get('description') || '',
        };
    }, [searchParams]);

    const handleBack = useCallback(() => {
        router.back();
    }, [router]);

    const handleAddToList = useCallback((item: DownloadItem) => {
        addDownloadItem(item);
        router.push('/downloads');
    }, [addDownloadItem, router]);

    if (!image) {
        return (
            <div className="flex-1 flex items-center justify-center">
                <p className="text-muted-foreground">Imagem não encontrada</p>
            </div>
        );
    }

    return (
        <div className="flex-1 flex flex-col overflow-hidden">
            <ImageDetail
                image={image}
                onBack={handleBack}
                onAddToList={handleAddToList}
            />
        </div>
    );
}

export default function ImagePage() {
    return (
        <Suspense fallback={
            <div className="flex-1 flex items-center justify-center">
                <p className="text-muted-foreground">Carregando...</p>
            </div>
        }>
            <ImagePageContent />
        </Suspense>
    );
}

'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../contexts/AppContext';
import SearchPanel from '../components/SearchPanel';
import DownloadList from '../components/DownloadList';
import type { ImageResult } from '../lib/types';

export default function HomePage() {
    const router = useRouter();
    const { items, addItem, updateItem, removeItem } = useDownloads();

    const handleImageSelect = useCallback((image: ImageResult) => {
        router.push(`/image?id=${encodeURIComponent(image.id)}&url=${encodeURIComponent(image.previewURL || image.url)}&downloadURL=${encodeURIComponent(image.downloadURL)}&width=${image.width}&height=${image.height}&author=${encodeURIComponent(image.author)}&source=${encodeURIComponent(image.source)}&tags=${encodeURIComponent((image.tags || []).join(','))}&description=${encodeURIComponent(image.description || '')}`);
    }, [router]);

    const handleQuickAdd = useCallback((image: ImageResult) => {
        addItem(image);
    }, [addItem]);

    const handleProcess = useCallback(() => {
        router.push('/downloads/select-folder');
    }, [router]);

    return (
        <>
            <div className="flex-1 flex flex-col border-r border-border overflow-hidden">
                <SearchPanel
                    onImageSelect={handleImageSelect}
                    onAddToList={handleQuickAdd}
                />
            </div>
            {items.length > 0 && (
                <div className="w-96 flex flex-col border-l border-border overflow-hidden">
                    <DownloadList
                        items={items}
                        onUpdateItem={updateItem}
                        onRemoveItem={removeItem}
                        onProcess={handleProcess}
                    />
                </div>
            )}
        </>
    );
}

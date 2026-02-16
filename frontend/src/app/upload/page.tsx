'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../../contexts/AppContext';
import ImageUpload from '../../components/ImageUpload';
import DownloadList from '../../components/DownloadList';
import type { ImageResult } from '../../lib/types';

export default function UploadPage() {
    const router = useRouter();
    const { items, addItem, updateItem, removeItem } = useDownloads();

    const handleImageSelect = useCallback((image: ImageResult) => {
        // For uploads, add to list directly then navigate to image detail
        addItem(image);
        router.push('/downloads');
    }, [addItem, router]);

    const handleProcess = useCallback(() => {
        router.push('/downloads/select-folder');
    }, [router]);

    return (
        <>
            <div className="flex-1 flex flex-col border-r border-border overflow-hidden">
                <ImageUpload onImageSelect={handleImageSelect} />
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

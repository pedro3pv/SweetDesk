'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../../contexts/AppContext';
import DownloadList from '../../components/DownloadList';

export default function DownloadsPage() {
    const router = useRouter();
    const { items, updateItem, removeItem } = useDownloads();

    const handleProcess = useCallback(() => {
        router.push('/downloads/select-folder');
    }, [router]);

    return (
        <div className="flex-1 flex flex-col overflow-hidden">
            <DownloadList
                items={items}
                onUpdateItem={updateItem}
                onRemoveItem={removeItem}
                onProcess={handleProcess}
            />
        </div>
    );
}

'use client';

import { useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useDownloads } from '../../../contexts/AppContext';
import FolderSelect from '../../../components/FolderSelect';

export default function SelectFolderPage() {
    const router = useRouter();
    const { setSavePath } = useDownloads();

    const handleSelect = useCallback((path: string) => {
        setSavePath(path);
        router.push('/downloads/processing');
    }, [setSavePath, router]);

    const handleCancel = useCallback(() => {
        router.push('/downloads');
    }, [router]);

    return (
        <div className="flex-1 flex flex-col overflow-hidden">
            <FolderSelect
                onSelect={handleSelect}
                onCancel={handleCancel}
            />
        </div>
    );
}

'use client';

import { useState, useCallback } from 'react';
import type { ImageResult, DownloadItem, AppView } from '@/lib/types';
import SearchPanel from '@/components/SearchPanel';
import ImageUpload from '@/components/ImageUpload';
import ImageDetail from '@/components/ImageDetail';
import DownloadList from '@/components/DownloadList';
import FolderSelect from '@/components/FolderSelect';
import ProcessingView from '@/components/ProcessingView';
import CompleteView from '@/components/CompleteView';

export default function Home() {
    const [view, setView] = useState<AppView>('search');
    const [selectedImage, setSelectedImage] = useState<ImageResult | null>(null);
    const [downloadItems, setDownloadItems] = useState<DownloadItem[]>([]);
    const [savePath, setSavePath] = useState('');

    // Navigate to image detail
    const handleImageSelect = useCallback((image: ImageResult) => {
        setSelectedImage(image);
        setView('image-detail');
    }, []);

    // Add item to download list
    const handleAddToList = useCallback((itemOrImage: DownloadItem | ImageResult) => {
        let item: DownloadItem;
        if ('image' in itemOrImage) {
            item = itemOrImage as DownloadItem;
        } else {
            const image = itemOrImage as ImageResult;
            item = {
                id: `dl-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
                image,
                name: image.description || `wallpaper-${image.id}`,
                upscale: true,
                dimension: '3840x2160',
                aspect: '16:9',
                selected: true,
            };
        }
        setDownloadItems(prev => {
            // Don't add duplicates by image id
            if (prev.some(p => p.image.id === item.image.id)) return prev;
            return [...prev, item];
        });
        setView('download-list');
    }, []);

    // Quick add from search grid (doesn't navigate)
    const handleQuickAdd = useCallback((image: ImageResult) => {
        setDownloadItems(prev => {
            if (prev.some(p => p.image.id === image.id)) return prev;
            return [...prev, {
                id: `dl-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
                image,
                name: image.description || `wallpaper-${image.id}`,
                upscale: true,
                dimension: '3840x2160',
                aspect: '16:9',
                selected: true,
            }];
        });
    }, []);

    // Update a download item
    const handleUpdateItem = useCallback((id: string, updates: Partial<DownloadItem>) => {
        setDownloadItems(prev => prev.map(item =>
            item.id === id ? { ...item, ...updates } : item
        ));
    }, []);

    // Remove a download item
    const handleRemoveItem = useCallback((id: string) => {
        setDownloadItems(prev => prev.filter(item => item.id !== id));
    }, []);

    // Start processing
    const handleProcess = useCallback(() => {
        setView('folder-select');
    }, []);

    // Folder selected - start processing
    const handleFolderSelect = useCallback((path: string) => {
        setSavePath(path);
        setView('processing');
    }, []);

    // Processing complete
    const handleProcessComplete = useCallback(() => {
        setView('complete');
    }, []);

    // Back to search
    const handleBackToList = useCallback(() => {
        setDownloadItems([]);
        setSelectedImage(null);
        setView('search');
    }, []);

    // Determine which panels to show
    const showLeftPanel = view === 'search' || view === 'upload' || view === 'image-detail';
    const showCenterPanel = view === 'download-list' || view === 'folder-select' || view === 'processing';
    const showRightPanel = view === 'complete';

    return (
        <div className="h-screen flex flex-col bg-background text-foreground overflow-hidden">
            {/* Top Navigation */}
            <nav className="flex items-center justify-between px-5 py-3 border-b border-border bg-card">
                <div className="flex items-center gap-4">
                    <h1 className="text-sm font-bold text-foreground tracking-tight">SWEETDESK</h1>
                    <span className="text-xs text-muted-foreground">Wallpaper Processing</span>
                </div>
                <div className="flex items-center gap-2">
                    {/* Search Tab */}
                    <button
                        onClick={() => setView('search')}
                        className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
                            view === 'search'
                                ? 'bg-primary text-primary-foreground'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                        }`}
                    >
                        Buscar
                    </button>
                    {/* Upload Tab */}
                    <button
                        onClick={() => setView('upload')}
                        className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${
                            view === 'upload'
                                ? 'bg-primary text-primary-foreground'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                        }`}
                    >
                        Upload
                    </button>
                    {/* Download List Tab */}
                    <button
                        onClick={() => setView('download-list')}
                        className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors relative ${
                            view === 'download-list' || view === 'folder-select' || view === 'processing'
                                ? 'bg-primary text-primary-foreground'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                        }`}
                    >
                        Lista
                        {downloadItems.length > 0 && (
                            <span className="absolute -top-1 -right-1 w-4 h-4 flex items-center justify-center text-[9px] font-bold bg-accent text-accent-foreground rounded-full">
                                {downloadItems.length}
                            </span>
                        )}
                    </button>
                </div>
            </nav>

            {/* Main Content */}
            <main className="flex-1 flex overflow-hidden">
                {/* Left Panel: Search / Upload / Image Detail */}
                {showLeftPanel && (
                    <div className="flex-1 flex flex-col border-r border-border overflow-hidden">
                        {view === 'search' && (
                            <SearchPanel
                                onImageSelect={handleImageSelect}
                                onAddToList={handleQuickAdd}
                            />
                        )}
                        {view === 'upload' && (
                            <ImageUpload
                                onImageSelect={handleImageSelect}
                                onAddToList={handleQuickAdd}
                            />
                        )}
                        {view === 'image-detail' && selectedImage && (
                            <ImageDetail
                                image={selectedImage}
                                onBack={() => setView('search')}
                                onAddToList={handleAddToList}
                            />
                        )}
                    </div>
                )}

                {/* Center Panel: Download List / Folder Select / Processing */}
                {showCenterPanel && (
                    <div className="flex-1 flex flex-col overflow-hidden">
                        {view === 'download-list' && (
                            <DownloadList
                                items={downloadItems}
                                onUpdateItem={handleUpdateItem}
                                onRemoveItem={handleRemoveItem}
                                onProcess={handleProcess}
                            />
                        )}
                        {view === 'folder-select' && (
                            <FolderSelect
                                onSelect={handleFolderSelect}
                                onCancel={() => setView('download-list')}
                            />
                        )}
                        {view === 'processing' && (
                            <ProcessingView
                                items={downloadItems}
                                savePath={savePath}
                                onComplete={handleProcessComplete}
                                onCancel={() => setView('download-list')}
                            />
                        )}
                    </div>
                )}

                {/* Right Panel: Complete */}
                {showRightPanel && (
                    <div className="flex-1 flex flex-col overflow-hidden">
                        <CompleteView
                            itemCount={downloadItems.filter(i => i.selected).length}
                            onBackToList={handleBackToList}
                        />
                    </div>
                )}

                {/* Side panel: Show download list alongside search/upload when items exist */}
                {showLeftPanel && downloadItems.length > 0 && (
                    <div className="w-96 flex flex-col border-l border-border overflow-hidden">
                        <DownloadList
                            items={downloadItems}
                            onUpdateItem={handleUpdateItem}
                            onRemoveItem={handleRemoveItem}
                            onProcess={handleProcess}
                        />
                    </div>
                )}
            </main>
        </div>
    );
}

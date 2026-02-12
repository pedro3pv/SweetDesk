// Global type declarations for Wails Go bindings
// This file provides types for window.go.main.App.* methods

export {};

export interface BatchItem {
    id: string;
    base64Data: string;
    downloadURL: string;
    name: string;
    dimension: string;
}

export interface BatchItemStatus {
    id: string;
    status: 'pending' | 'processing' | 'done' | 'error';
    error?: string;
}

export interface ProcessingStatus {
    isProcessing: boolean;
    total: number;
    current: number;
    progress: number;
    items: BatchItemStatus[];
    done: boolean;
}

declare global {
    interface Window {
        go?: {
            main?: {
                App?: {
                    ProcessImage?: (base64Data: string, targetWidth: number, targetHeight: number, savePath: string, fileName: string) => Promise<string>;
                    DownloadImage?: (url: string) => Promise<string>;
                    SelectDirectory?: () => Promise<string>;
                    GetDefaultSavePath?: () => Promise<string>;
                    SearchImages?: (query: string, page: number, perPage: number) => Promise<unknown[]>;
                    UpscaleImage?: (base64Data: string, imageType: string, scale: number) => Promise<string>;
                    Greet?: (name: string) => Promise<string>;
                    ProcessBatch?: (items: BatchItem[], savePath: string) => Promise<void>;
                    GetProcessingStatus?: () => Promise<ProcessingStatus>;
                };
            };
        };
    }
}

export interface ImageResult {
    id: string;
    url: string;
    downloadURL: string;
    previewURL: string;
    width: number;
    height: number;
    author: string;
    source: string;
    tags: string[];
    description?: string;
}

export interface DownloadItem {
    id: string;
    image: ImageResult;
    name: string;
    upscale: boolean;
    dimension: string;
    aspect: string;
    selected: boolean;
}

export interface FilterState {
    aspect: string;
    dimensions: string;
    type: string;
}

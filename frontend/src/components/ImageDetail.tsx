'use client';

import { useState } from 'react';
import type { ImageResult, DownloadItem } from '../lib/types';

interface ImageDetailProps {
    image: ImageResult;
    onBack: () => void;
    onAddToList: (item: DownloadItem) => void;
}

const DIMENSIONS = ['1920x1080', '2560x1440', '3840x2160', '5120x2880', '7680x4320'];
const ASPECTS = ['16:9', '4:3', '21:9', '1:1'];

export default function ImageDetail({ image, onBack, onAddToList }: ImageDetailProps) {
    const [upscale, setUpscale] = useState(true);
    const [dimension, setDimension] = useState('3840x2160');
    const [aspect, setAspect] = useState('16:9');

    const handleAdd = () => {
        const item: DownloadItem = {
            id: `dl-${Date.now()}-${Math.random().toString(36).slice(2, 6)}`,
            image,
            name: image.description || `wallpaper-${image.id}`,
            upscale,
            dimension,
            aspect,
            selected: true,
        };
        onAddToList(item);
    };

    return (
        <div className="flex flex-col h-full overflow-y-auto">
            {/* Back button */}
            <div className="flex items-center gap-3 px-5 py-4 border-b border-border">
                <button
                    onClick={onBack}
                    className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
                    </svg>
                    Voltar
                </button>
            </div>

            {/* Image Preview */}
            <div className="flex-1 flex flex-col p-5 gap-5">
                <div className="relative rounded-lg overflow-hidden bg-muted flex items-center justify-center" style={{ maxHeight: '400px' }}>
                    <img
                        src={image.previewURL || image.url}
                        alt={image.description || 'Image preview'}
                        className="max-w-full max-h-[400px] object-contain"
                        crossOrigin="anonymous"
                    />
                    {/* Download icon overlay */}
                    <button
                        onClick={handleAdd}
                        className="absolute top-3 right-3 p-2 bg-primary/80 text-primary-foreground rounded-md hover:bg-primary transition-colors"
                        aria-label="Adicionar a lista"
                    >
                        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                        </svg>
                    </button>
                </div>

                {/* Description */}
                <div>
                    <h3 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Descricao</h3>
                    <p className="text-sm text-foreground leading-relaxed">
                        {image.description || 'Sem descricao disponivel para esta imagem. A imagem foi obtida de fontes externas e pode ser utilizada como wallpaper em alta resolucao.'}
                    </p>
                </div>

                {/* Tags */}
                <div>
                    <div className="flex flex-wrap gap-2">
                        {image.tags?.map(tag => (
                            <span key={tag} className="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded-md border border-border">
                                {tag}
                            </span>
                        ))}
                        <span className="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded-md border border-border">
                            {image.width}x{image.height}
                        </span>
                        <span className="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded-md border border-border">
                            {image.source}
                        </span>
                    </div>
                </div>

                {/* Upscale toggle */}
                <div>
                    <label className="flex items-center gap-2 cursor-pointer">
                        <input
                            type="checkbox"
                            checked={upscale}
                            onChange={(e) => setUpscale(e.target.checked)}
                            className="w-4 h-4 rounded border-border text-primary focus:ring-primary bg-input accent-primary"
                        />
                        <span className="text-sm text-foreground font-medium">Upscale</span>
                    </label>
                </div>

                {/* Dimension selection */}
                <div>
                    <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Dimensao</p>
                    <div className="flex flex-wrap gap-2">
                        {DIMENSIONS.map(d => (
                            <button
                                key={d}
                                onClick={() => setDimension(d)}
                                className={`text-xs px-3 py-1.5 rounded-md border transition-colors ${
                                    dimension === d
                                        ? 'bg-primary text-primary-foreground border-primary'
                                        : 'bg-secondary text-secondary-foreground border-border hover:border-muted-foreground'
                                }`}
                            >
                                {d}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Aspect selection */}
                <div>
                    <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">Aspecto</p>
                    <div className="flex flex-wrap gap-2">
                        {ASPECTS.map(a => (
                            <button
                                key={a}
                                onClick={() => setAspect(a)}
                                className={`text-xs px-3 py-1.5 rounded-md border transition-colors ${
                                    aspect === a
                                        ? 'bg-primary text-primary-foreground border-primary'
                                        : 'bg-secondary text-secondary-foreground border-border hover:border-muted-foreground'
                                }`}
                            >
                                {a}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Add to list button */}
                <button
                    onClick={handleAdd}
                    className="w-full py-3 bg-primary text-primary-foreground rounded-md text-sm font-medium hover:bg-primary/90 transition-colors"
                >
                    Adicionar a lista
                </button>
            </div>
        </div>
    );
}

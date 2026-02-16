'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useDownloads } from '../contexts/AppContext';

export default function Header() {
    const pathname = usePathname();
    const { items } = useDownloads();

    const navItems = [
        { href: '/', label: 'Buscar', match: (p: string) => p === '/' },
        { href: '/upload', label: 'Upload', match: (p: string) => p === '/upload' },
        {
            href: '/downloads',
            label: 'Lista',
            badge: items.length,
            match: (p: string) => p.startsWith('/downloads'),
        },
    ];

    return (
        <nav className="flex items-center justify-between px-6 py-4 border-b border-border bg-card">
            <div className="flex items-center gap-6">
                <h1 className="text-base lg:text-lg font-bold text-foreground tracking-tight">SWEETDESK</h1>
                <span className="text-sm lg:text-base text-muted-foreground">Wallpaper Processing</span>
            </div>
            <div className="flex items-center gap-3">
                {navItems.map(nav => (
                    <Link
                        key={nav.href}
                        href={nav.href}
                        className={`px-4 py-2 rounded-md text-sm lg:text-base font-medium transition-colors relative ${
                            nav.match(pathname)
                                ? 'bg-primary text-primary-foreground'
                                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                        }`}
                    >
                        {nav.label}
                        {nav.badge !== undefined && nav.badge > 0 && (
                            <span className="absolute -top-1 -right-1 w-5 h-5 flex items-center justify-center text-xs font-bold bg-accent text-accent-foreground rounded-full">
                                {nav.badge}
                            </span>
                        )}
                    </Link>
                ))}
            </div>
        </nav>
    );
}

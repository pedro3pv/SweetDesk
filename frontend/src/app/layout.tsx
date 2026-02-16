import React from "react"
import type { Metadata, Viewport } from "next";
import localFont from "next/font/local";
import "./globals.css";
import { AppProvider } from "../contexts/AppContext";
import Header from "../components/Header";

const geistSans = localFont({
    src: "./fonts/GeistVF.woff",
    variable: "--font-geist-sans",
    weight: "100 900",
});
const geistMono = localFont({
    src: "./fonts/GeistMonoVF.woff",
    variable: "--font-geist-mono",
    weight: "100 900",
});

export const metadata: Metadata = {
    title: "SweetDesk - Wallpaper Processing",
    description: "Busque, processe e salve wallpapers em alta resolução com upscale AI",
};

export const viewport: Viewport = {
    themeColor: "#1a1625",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="pt-BR" className="dark">
            <body
                className={`${geistSans.variable} ${geistMono.variable} antialiased`}
            >
                <AppProvider>
                    <div className="h-screen flex flex-col bg-background text-foreground overflow-hidden">
                        <Header />
                        <main className="flex-1 flex overflow-hidden">
                            {children}
                        </main>
                    </div>
                </AppProvider>
            </body>
        </html>
    );
}

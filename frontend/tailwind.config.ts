import type { Config } from "tailwindcss";

const config: Config = {
    content: [
        "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
        "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
        "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
    ],
    darkMode: 'class',
    theme: {
        extend: {
            colors: {
                // Purple Primary - Energetic and Motivating
                purple: {
                    50: '#f9f5ff',
                    100: '#f3e8ff',
                    200: '#e9d5ff',
                    300: '#d8b4fe',
                    400: '#c084fc',
                    500: '#a855f7', // PRIMARY
                    600: '#9333ea',
                    700: '#7e22ce',
                    800: '#6b21a8',
                    900: '#581c87',
                    950: '#3f0f5c',
                },
                // Indigo Secondary - Complementary (Blue-Purple)
                indigo: {
                    50: '#f0f4ff',
                    500: '#6366f1',
                    600: '#4f46e5',
                    700: '#4338ca',
                },
                // Accent Colors (Triadic Harmony)
                amber: {
                    400: '#fbbf24',
                    500: '#f59e0b',
                },
                cyan: {
                    400: '#22d3ee',
                    500: '#06b6d4',
                },
                // Dark Mode - Purple Dark Harmony
                dark: {
                    bg: '#1a1625',
                    surface: '#2d1f3d',
                    border: '#4a3361',
                    text: '#f0f0f0',
                    'text-secondary': '#b4a7c6',
                },
            },
            fontSize: {
                xs: '11px',
                sm: '13px',
                base: '14px',
                md: '16px',
                lg: '18px',
                xl: '20px',
                '2xl': '24px',
                '3xl': '30px',
                '4xl': '36px',
            },
        },
    },
    plugins: [],
};
export default config;

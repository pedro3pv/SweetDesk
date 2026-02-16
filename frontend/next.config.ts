import type { NextConfig } from "next";

const nextConfig: NextConfig = {
    /* config options here */
    output: "export",
    allowedDevOrigins: ["wails.localhost"],
    // Empty turbopack config to silence Next.js 16 warning
    turbopack: {},
    // Disable HMR in Wails to prevent WebSocket TLS errors
    webpack: (config, { dev, isServer }) => {
        if (dev && !isServer) {
            config.watchOptions = {
                ...config.watchOptions,
                poll: 1000,
                aggregateTimeout: 300,
            };
        }
        return config;
    },
};

export default nextConfig;

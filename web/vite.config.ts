import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
  plugins: [
    VitePWA({
      registerType: 'prompt',
      manifest: {
        name: 'Frivolin Viewer – Stadt NexVerse',
        short_name: 'Frivolin',
        description: 'Browserbasierter 3D-Viewer für Stadt NexVerse',
        lang: 'de',
        start_url: '/',
        scope: '/',
        display: 'standalone',
        theme_color: '#101327',
        background_color: '#101327',
        icons: [
          { src: '/icon.svg', sizes: 'any', type: 'image/svg+xml', purpose: 'any' }
        ]
      },
      workbox: {
        navigateFallback: '/index.html',
        maximumFileSizeToCacheInBytes: 2 * 1024 * 1024,
        // Babylon.js is currently a large bundle; do not precache it.
        // It is fetched normally and can be split into smaller chunks later.
        manifestEntries: undefined,
        manifestTransforms: [async (entries) => ({
          manifest: entries.filter(entry => !/\\/assets\\/index-[^/]+\\.js$/.test(entry.url)),
          warnings: []
        })],
        runtimeCaching: [],
        cleanupOutdatedCaches: true
      }
    })
  ],
  server: { port: 5173 }
});

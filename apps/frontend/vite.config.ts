import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { visualizer } from 'rollup-plugin-visualizer'
import { defineConfig, loadEnv } from 'vite'
import viteCompression from 'vite-plugin-compression'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    // 加载环境变量
    const env = loadEnv(mode, process.cwd(), 'VITE_')
    const apiTarget = env.VITE_API_TARGET || 'http://localhost:8080'

    return {
        plugins: [
            vue(),
            viteCompression({
                verbose: true,
                disable: false,
                threshold: 10240,
                algorithm: 'gzip',
                ext: '.gz',
            }),
            visualizer({
                open: false,
                gzipSize: true,
                brotliSize: true,
                filename: 'stats.html'
            })
        ],
        resolve: {
            alias: {
                '@': resolve(__dirname, './src')
            }
        },
        server: {
            host: '0.0.0.0',
            port: 3000,
            proxy: {
                '/api': {
                    target: apiTarget,
                    changeOrigin: true
                }
            }
        },
        build: {
            outDir: 'dist',
            sourcemap: false,
            minify: 'terser',
            rollupOptions: {
                output: {
                    manualChunks: {
                        'vue-vendor': ['vue', 'vue-router', 'pinia'],
                        'element-plus': ['element-plus'],
                        'axios': ['axios']
                    }
                }
            }
        }
    }
})

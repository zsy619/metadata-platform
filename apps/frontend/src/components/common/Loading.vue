<template>
    <div v-if="visible" class="loading-container" :class="{ 'loading-fullscreen': fullscreen, [`loading-size-${size}`]: true }">
        <div v-if="type === 'spinner'" class="loading-spinner" :style="spinnerStyle">
            <svg class="circular" viewBox="0 0 50 50">
                <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="4" />
            </svg>
            <p v-if="text" class="loading-text">{{ text }}</p>
        </div>
        <div v-else-if="type === 'dots'" class="loading-dots">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
            <p v-if="text" class="loading-text">{{ text }}</p>
        </div>
        <div v-else-if="type === 'ring'" class="loading-ring" :style="spinnerStyle">
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <p v-if="text" class="loading-text">{{ text }}</p>
        </div>
        <div v-else-if="type === '骨架屏'" class="loading-skeleton">
            <slot></slot>
        </div>
    </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'

interface Props {
    visible?: boolean
    type?: 'spinner' | 'dots' | 'ring' | '骨架屏'
    size?: 'large' | 'default' | 'small'
    text?: string
    fullscreen?: boolean
    color?: string
    strokeWidth?: number
}

const props = withDefaults(defineProps<Props>(), {
    visible: true,
    type: 'spinner',
    size: 'default',
    text: '',
    fullscreen: false,
    color: '#409eff',
    strokeWidth: 4
})

const spinnerStyle = computed(() => {
    const sizeMap = { large: 60, default: 40, small: 24 }
    const spinnerSize = sizeMap[props.size]
    return {
        width: `${spinnerSize}px`,
        height: `${spinnerSize}px`,
        '--loading-color': props.color,
        '--loading-stroke-width': `${props.strokeWidth}px`
    }
})
</script>
<style scoped>
.loading-container {
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.loading-container.loading-fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.9);
    z-index: 9999;
    flex-direction: column;
}

.loading-size-large .loading-text {
    font-size: 16px;
}

.loading-size-small .loading-text {
    font-size: 12px;
}

.loading-text {
    margin-top: 12px;
    color: #606266;
    font-size: 14px;
}

.loading-spinner {
    position: relative;
}

.loading-spinner .circular {
    width: 100%;
    height: 100%;
    animation: rotate 2s linear infinite;
}

.loading-spinner .path {
    stroke: var(--loading-color);
    stroke-width: var(--loading-stroke-width);
    stroke-linecap: round;
    animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
    100% {
        transform: rotate(360deg);
    }
}

@keyframes dash {
    0% {
        stroke-dasharray: 1, 150;
        stroke-dashoffset: 0;
    }
    50% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -35;
    }
    100% {
        stroke-dasharray: 90, 150;
        stroke-dashoffset: -124;
    }
}

.loading-dots {
    display: flex;
    gap: 6px;
}

.loading-dots .dot {
    width: 10px;
    height: 10px;
    background: var(--loading-color, #409eff);
    border-radius: 50%;
    animation: dot-bounce 1.4s infinite ease-in-out both;
}

.loading-dots .dot:nth-child(1) {
    animation-delay: -0.32s;
}

.loading-dots .dot:nth-child(2) {
    animation-delay: -0.16s;
}

@keyframes dot-bounce {
    0%, 80%, 100% {
        transform: scale(0);
    }
    40% {
        transform: scale(1);
    }
}

.loading-ring {
    position: relative;
}

.loading-ring div {
    position: absolute;
    width: 100%;
    height: 100%;
    border: 3px solid transparent;
    border-top-color: var(--loading-color);
    border-radius: 50%;
    animation: ring-spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
}

.loading-ring div:nth-child(1) {
    animation-delay: -0.45s;
}

.loading-ring div:nth-child(2) {
    animation-delay: -0.3s;
}

.loading-ring div:nth-child(3) {
    animation-delay: -0.15s;
}

@keyframes ring-spin {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}
</style>

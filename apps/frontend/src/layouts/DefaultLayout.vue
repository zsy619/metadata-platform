<template>
    <div :class="classObj" class="app-wrapper">
        <div v-if="device === 'mobile' && sidebar.opened" class="drawer-bg" @click="handleClickOutside" />
        <sidebar class="sidebar-container" />
        <div class="main-container">
            <div :class="{ 'fixed-header': fixedHeader }">
                <app-header />
            </div>
            <app-main />
        </div>
    </div>
</template>
<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { computed, onBeforeMount, onBeforeUnmount, onMounted } from 'vue'
import AppHeader from './components/AppHeader.vue'
import AppMain from './components/AppMain.vue'
import Sidebar from './components/Sidebar.vue'

const appStore = useAppStore()
const sidebar = computed(() => appStore.sidebar)
const device = computed(() => appStore.device)
const fixedHeader = true

const classObj = computed(() => {
    return {
        hideSidebar: !sidebar.value.opened,
        openSidebar: sidebar.value.opened,
        withoutAnimation: sidebar.value.withoutAnimation,
        mobile: device.value === 'mobile'
    }
})

const { body } = document
const WIDTH = 992 // refer to Bootstrap's responsive design

const isMobile = () => {
    const rect = body.getBoundingClientRect()
    return rect.width - 1 < WIDTH
}

const resizeHandler = () => {
    if (!document.hidden) {
        const mobile = isMobile()
        appStore.toggleDevice(mobile ? 'mobile' : 'desktop')

        if (mobile) {
            appStore.closeSidebar(true)
        }
    }
}

const handleClickOutside = () => {
    appStore.closeSidebar(false)
}

onBeforeMount(() => {
    window.addEventListener('resize', resizeHandler)
})

onMounted(() => {
    const mobile = isMobile()
    if (mobile) {
        appStore.toggleDevice('mobile')
        appStore.closeSidebar(true)
    }
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', resizeHandler)
})
</script>
<style scoped>
.app-wrapper {
    position: relative;
    height: 100%;
    width: 100%;
    display: flex;
}

.sidebar-container {
    transition: width 0.3s;
    width: 210px !important;
    background-color: #FFFFFF;
    height: 100%;
    position: fixed;
    font-size: 0;
    top: 0;
    bottom: 0;
    left: 0;
    z-index: 1001;
    overflow: hidden;
}

.main-container {
    min-height: 100%;
    transition: margin-left 0.3s;
    margin-left: 210px;
    width: 100%;
    position: relative;
    background-color: #f0f2f5;
}

.fixed-header {
    position: sticky;
    top: 0;
    z-index: 9;
    width: 100%;
    transition: width 0.3s;
}

.hideSidebar .sidebar-container {
    width: 64px !important;
}

.hideSidebar .main-container {
    margin-left: 64px;
}

/* mobile responsive */
.mobile .main-container {
    margin-left: 0px;
}

.mobile .sidebar-container {
    transition: transform 0.3s;
    width: 210px !important;
}

.mobile.hideSidebar .sidebar-container {
    pointer-events: none;
    transition-duration: 0.3s;
    transform: translate3d(-210px, 0, 0);
}

.drawer-bg {
    background: #000;
    opacity: 0.3;
    width: 100%;
    top: 0;
    height: 100%;
    position: absolute;
    z-index: 999;
}

.fixed-header {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 9;
    width: calc(100% - 210px);
    transition: width 0.3s;
}

.hideSidebar .fixed-header {
    width: calc(100% - 64px)
}

.mobile .fixed-header {
    width: 100%;
}
</style>
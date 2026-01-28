<template>
    <el-scrollbar wrap-class="scrollbar-wrapper">
        <el-menu :default-active="activeMenu" :collapse="!sidebar.opened" :unique-opened="false" :collapse-transition="true" background-color="#FFFFFF" text-color="#303133" active-text-color="#4051B5" mode="vertical" class="sidebar-menu">
            <div class="sidebar-logo-container" :class="{ 'collapse': !sidebar.opened }">
                <!-- Logo Logic kept same -->
                <transition name="sidebarLogoFade">
                    <router-link v-if="!sidebar.opened" key="collapse" class="sidebar-logo-link" to="/">
                        <div v-if="showLogo" class="sidebar-logo-container-inner">
                            <AnimatedLogo class="sidebar-logo" />
                        </div>
                        <h1 class="sidebar-title" v-else>元数据平台</h1>
                    </router-link>
                    <router-link v-else key="expand" class="sidebar-logo-link" to="/">
                        <div v-if="showLogo" class="sidebar-logo-container-inner">
                            <AnimatedLogo class="sidebar-logo" />
                        </div>
                        <h1 class="sidebar-title">元数据管理平台</h1>
                    </router-link>
                </transition>
            </div>
            <sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path" />
        </el-menu>
    </el-scrollbar>
</template>
<script setup lang="ts">
import AnimatedLogo from '@/components/AnimatedLogo.vue'
import { useAppStore } from '@/stores/app'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SidebarItem from './SidebarItem.vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const sidebar = computed(() => appStore.sidebar)
const showLogo = true

// 获取所有路由作为菜单，实际应从 permissionStore 获取
const routes = computed(() => router.options.routes)

const activeMenu = computed(() => {
    const { meta, path } = route
    if (meta.activeMenu) {
        return meta.activeMenu as string
    }
    return path
})
</script>
<style scoped>
.scrollbar-wrapper {
    height: 100%;
    overflow-x: hidden !important;
}

.scrollbar-wrapper :deep(.el-scrollbar__bar.is-vertical) {
    right: 0px;
    width: 6px;
}

.scrollbar-wrapper :deep(.el-scrollbar__thumb) {
    background-color: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
}

.scrollbar-wrapper :deep(.el-scrollbar__thumb:hover) {
    background-color: rgba(255, 255, 255, 0.3);
}

.sidebar-menu {
    border: none;
    height: 100%;
    width: 100% !important;
    padding: 8px 0;
}

.sidebar-menu:not(.el-menu--collapse) {
    width: 210px;
}

.sidebar-logo-container {
    position: relative;
    width: 100%;
    height: 60px;
    line-height: 60px;
    background: #4051B5;
    text-align: center;
    overflow: hidden;
    margin-bottom: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.sidebar-logo-container .sidebar-logo-link {
    height: 100%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
}

.sidebar-logo-container .sidebar-logo-link:hover {
    opacity: 0.9;
}

.sidebar-logo-container .sidebar-logo {
    width: 36px;
    height: 36px;
    vertical-align: middle;
    margin-right: 12px;
}

.sidebar-logo-container .sidebar-title {
    display: inline-block;
    margin: 0;
    color: #fff;
    font-weight: 600;
    line-height: 60px;
    font-size: 16px;
    font-family: 'PingFang SC', 'Microsoft YaHei', Arial, sans-serif;
    vertical-align: middle;
    letter-spacing: 0.5px;
}

.sidebar-logo-container.collapse .sidebar-logo {
    margin-right: 0px;
}

.sidebarLogoFade-enter-active,
.sidebarLogoFade-leave-active {
    transition: opacity 0.3s;
}

.sidebarLogoFade-enter-from,
.sidebarLogoFade-leave-to {
    opacity: 0;
}

/* 修复折叠状态下的图标样式 - 增强版 */
:deep(.el-menu--collapse .el-sub-menu__title),
:deep(.el-menu--collapse .el-menu-item) {
    justify-content: center !important;
    padding: 0 !important;
    margin: 4px 0 !important;
    /* 强制移除左右边距，保留上下边距 */
    width: 100% !important;
}

/* 专门针对图标的修复 */
:deep(.el-menu--collapse .el-sub-menu__title .el-icon),
:deep(.el-menu--collapse .el-menu-item .el-icon) {
    margin: 0 !important;
    /* 清除所有边距 */
    padding: 0 !important;
    font-size: 18px;
    width: auto !important;
    /* 防止宽度被撑开 */
}

/* 隐藏标题和箭头 */
:deep(.el-menu--collapse .el-sub-menu__title span),
:deep(.el-menu--collapse .el-menu-item .menu-title),
:deep(.el-menu--collapse .el-sub-menu__icon-arrow) {
    display: none !important;
    width: 0 !important;
    height: 0 !important;
    overflow: hidden;
    margin: 0 !important;
}

/* 修复 tooltip 触发层的样式（如果有） */
:deep(.el-menu--collapse .el-sub-menu__title > div),
:deep(.el-menu--collapse .el-menu-item > div) {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
}
</style>

<template>
    <el-scrollbar wrap-class="scrollbar-wrapper">
        <el-menu :default-active="activeMenu" :collapse="!sidebar.opened" :unique-opened="false" :collapse-transition="false" background-color="#001529" text-color="#bfcbd9" active-text-color="#409eff" mode="vertical" class="sidebar-menu">
            <div class="sidebar-logo-container" :class="{ 'collapse': !sidebar.opened }">
                <!-- Logo Logic kept same -->
                <transition name="sidebarLogoFade">
                    <router-link v-if="!sidebar.opened" key="collapse" class="sidebar-logo-link" to="/">
                        <img src="@/assets/logo.png" class="sidebar-logo" v-if="showLogo" />
                        <h1 class="sidebar-title" v-else>元数据平台</h1>
                    </router-link>
                    <router-link v-else key="expand" class="sidebar-logo-link" to="/">
                        <img src="@/assets/logo.png" class="sidebar-logo" v-if="showLogo" />
                        <h1 class="sidebar-title">元数据管理平台</h1>
                    </router-link>
                </transition>
            </div>
            <sidebar-item v-for="route in routes" :key="route.path" :item="route" :base-path="route.path" />
        </el-menu>
    </el-scrollbar>
</template>
<script setup lang="ts">
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
    width: 220px;
}

.sidebar-logo-container {
    position: relative;
    width: 100%;
    height: 60px;
    line-height: 60px;
    background: linear-gradient(135deg, #001529 0%, #002140 100%);
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
</style>

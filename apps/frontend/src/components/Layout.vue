<template>
    <div class="layout-container">
        <!-- 侧边栏 -->
        <aside class="sidebar" :class="{ 'collapsed': isSidebarCollapsed }">
            <div class="sidebar-header">
                <img src="../assets/logo.png" alt="Logo" class="logo-img" />
                <h2 class="sidebar-title" v-show="!isSidebarCollapsed">元数据管理平台</h2>
            </div>
            <nav class="sidebar-menu">
                <el-menu :default-active="activeMenu" class="sidebar-el-menu" :collapse="isSidebarCollapsed" :collapse-transition="true" background-color="#001529" text-color="#ffffff" active-text-color="#409eff" @select="handleMenuSelect">
                    <!-- 数据源管理 -->
                    <el-sub-menu index="1">
                        <template #title>
                            <el-icon>
                                <Connection />
                            </el-icon>
                            <span>数据源管理</span>
                        </template>
                        <el-menu-item index="/data-sources">数据源列表</el-menu-item>
                        <el-menu-item index="/data-sources/create">创建数据源</el-menu-item>
                        <!-- 元数据管理子菜单 -->
                        <el-sub-menu index="1-3">
                            <template #title>元数据管理</template>
                            <el-menu-item index="/metadata/tables">表与视图</el-menu-item>
                            <el-menu-item index="/metadata/fields">字段列表</el-menu-item>
                        </el-sub-menu>
                    </el-sub-menu>
                    <!-- 模型管理 -->
                    <el-sub-menu index="2">
                        <template #title>
                            <el-icon>
                                <Document />
                            </el-icon>
                            <span>模型管理</span>
                        </template>
                        <el-menu-item index="/models">模型列表</el-menu-item>
                        <el-menu-item index="/models/create">创建模型</el-menu-item>
                    </el-sub-menu>
                    <!-- 接口管理 -->
                    <el-sub-menu index="3">
                        <template #title>
                            <el-icon>
                                <Share />
                            </el-icon>
                            <span>接口管理</span>
                        </template>
                        <el-menu-item index="/apis">接口列表</el-menu-item>
                        <el-menu-item index="/apis/create">创建接口</el-menu-item>
                    </el-sub-menu>
                    <!-- 系统设置 -->
                    <el-menu-item index="/system/settings">
                        <el-icon>
                            <Setting />
                        </el-icon>
                        <span>系统设置</span>
                    </el-menu-item>
                </el-menu>
            </nav>
        </aside>
        <!-- 主内容区 -->
        <main class="main-content">
            <!-- 顶部导航 -->
            <header class="top-header">
                <div class="header-left">
                    <el-button type="text" @click="toggleSidebar" class="sidebar-toggle">
                        <el-icon v-if="!isSidebarCollapsed">
                            <Expand />
                        </el-icon>
                        <el-icon v-else>
                            <Fold />
                        </el-icon>
                    </el-button>
                    <el-breadcrumb separator="/">
                        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
                        <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
                    </el-breadcrumb>
                </div>
                <div class="header-right">
                    <el-dropdown>
                        <span class="user-info">
                            <el-avatar :size="32">
                                <el-icon>
                                    <User />
                                </el-icon>
                            </el-avatar>
                            <span class="username">管理员</span>
                            <el-icon class="el-icon--right">
                                <CaretBottom />
                            </el-icon>
                        </span>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item>
                                    <el-icon>
                                        <Setting />
                                    </el-icon>
                                    <span>个人设置</span>
                                </el-dropdown-item>
                                <el-dropdown-item divided>
                                    <el-icon>
                                        <SwitchButton />
                                    </el-icon>
                                    <span>退出登录</span>
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
            </header>
            <!-- 内容区域 -->
            <div class="content-wrapper">
                <router-view />
            </div>
        </main>
    </div>
</template>
<script setup lang="ts">
import {
    CaretBottom,
    Connection,
    Document,
    Expand,
    Fold,
    Setting,
    Share,
    SwitchButton,
    User
} from '@element-plus/icons-vue'
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// 响应式数据
const isSidebarCollapsed = ref(false)

// 计算属性
const activeMenu = computed(() => {
    return route.path
})

const currentTitle = computed(() => {
    return route.meta.title as string || '首页'
})

// 方法
const toggleSidebar = () => {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
}

const handleMenuSelect = (index: string) => {
    router.push(index)
}
</script>
<style scoped>
.layout-container {
    display: flex;
    width: 100%;
    height: 100vh;
    overflow: hidden;
    background-color: #f5f7fa;
}

/* 侧边栏样式 */
.sidebar {
    width: 240px;
    height: 100%;
    background-color: #001529;
    color: #fff;
    transition: width 0.3s cubic-bezier(0.2, 0, 0, 1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
}

.sidebar.collapsed {
    width: 64px;
}

.sidebar-header {
    height: 60px;
    padding: 0 16px;
    display: flex;
    align-items: center;
    gap: 12px;
    overflow: hidden;
    background-color: #002140;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-img {
    width: 32px;
    height: 32px;
    flex-shrink: 0;
    object-fit: contain;
}

.sidebar-title {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #fff;
    white-space: nowrap;
    opacity: 1;
    transition: opacity 0.2s;
}

.sidebar.collapsed .sidebar-title {
    opacity: 0;
    pointer-events: none;
}

.sidebar-menu {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
}

.sidebar-menu::-webkit-scrollbar {
    width: 0;
}

.sidebar-el-menu {
    border-right: none;
}

/* 修正折叠时的宽度 */
.sidebar-el-menu:not(.el-menu--collapse) {
    width: 240px;
}

/* 主内容区样式 */
.main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

/* 顶部导航样式 */
.top-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 64px;
    padding: 0 20px;
    background-color: #fff;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    z-index: 100;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 20px;
}

.sidebar-toggle {
    font-size: 20px;
    color: #333;
    cursor: pointer;
    padding: 8px;
    transition: background 0.3s;
    border-radius: 4px;
}

.sidebar-toggle:hover {
    background: rgba(0, 0, 0, 0.025);
}

.header-right {
    display: flex;
    align-items: center;
    gap: 20px;
}

.user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 0 8px;
    transition: background 0.3s;
    border-radius: 4px;
}

.user-info:hover {
    background: rgba(0, 0, 0, 0.025);
}

.username {
    font-size: 14px;
    color: #666;
}

/* 内容区域样式 */
.content-wrapper {
    flex: 1;
    padding: 0px;
    overflow-y: auto;
    background-color: #f0f2f5;
}

/* 滚动条样式 */
.content-wrapper::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

.content-wrapper::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 4px;
}

.content-wrapper::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 4px;
}

.content-wrapper::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
}
</style>
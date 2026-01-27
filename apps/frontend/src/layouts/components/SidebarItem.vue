<template>
    <div v-if="!item.meta?.hidden">
        <template v-if="hasOneShowingChild(item.children, item) && (!onlyOneChild.children || onlyOneChild.noShowingChildren) && !item.meta?.alwaysShow">
            <app-link v-if="onlyOneChild.meta" :to="resolvePath(onlyOneChild.path)">
                <el-menu-item :index="resolvePath(onlyOneChild.path)" :class="{ 'submenu-title-noDropdown': !isNest }">
                    <el-icon v-if="onlyOneChild.meta.icon || (item.meta && item.meta.icon)" class="menu-icon">
                        <component :is="onlyOneChild.meta.icon || item.meta.icon" />
                    </el-icon>
                    <template #title>
                        <span class="menu-title">{{ onlyOneChild.meta?.title }}</span>
                    </template>
                </el-menu-item>
            </app-link>
        </template>
        <el-sub-menu v-else :index="resolvePath(item.path)" popper-append-to-body>
            <template #title>
                <el-icon v-if="item.meta && item.meta.icon" class="menu-icon">
                    <component :is="item.meta.icon" />
                </el-icon>
                <span class="menu-title">{{ item.meta?.title }}</span>
            </template>
            <sidebar-item v-for="child in item.children" :key="child.path" :is-nest="true" :item="child" :base-path="resolvePath(child.path)" class="nest-menu" />
        </el-sub-menu>
    </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import AppLink from './Link.vue'

const props = defineProps({
    item: {
        type: Object,
        required: true
    },
    isNest: {
        type: Boolean,
        default: false
    },
    basePath: {
        type: String,
        default: ''
    }
})

const onlyOneChild = ref<any>(null)

const hasOneShowingChild = (children: any[] = [], parent: any) => {
    const showingChildren = children.filter(item => {
        if (item.meta?.hidden) {
            return false
        } else {
            onlyOneChild.value = item
            return true
        }
    })

    if (showingChildren.length === 1) {
        return true
    }

    if (showingChildren.length === 0) {
        onlyOneChild.value = { ...parent, path: '', noShowingChildren: true }
        return true
    }

    return false
}

const isExternal = (path: string) => /^(https?:|mailto:|tel:)/.test(path)

const resolvePath = (routePath: string) => {
    if (isExternal(routePath)) {
        return routePath
    }
    if (isExternal(props.basePath)) {
        return props.basePath
    }

    // Handle absolute paths
    if (routePath.startsWith('/')) {
        return routePath
    }

    // Join basePath and routePath, avoiding double slashes
    const base = props.basePath.replace(/\/+$/, '')
    const path = routePath.replace(/^\/+/, '')

    return base ? `${base}/${path}` : `/${path}`
}
</script>
<style scoped>
.menu-icon {
    font-size: 18px;
    margin-right: 8px;
    vertical-align: middle;
}

.menu-title {
    font-size: 14px;
    font-weight: 500;
}

:deep(.el-menu-item) {
    height: 48px;
    line-height: 48px;
    margin: 4px 8px;
    border-radius: 6px;
    transition: all 0.3s;
}

:deep(.el-menu-item:hover) {
    background-color: rgba(64, 158, 255, 0.1) !important;
}

:deep(.el-menu-item.is-active) {
    background-color: rgba(64, 158, 255, 0.15) !important;
    color: #409EFF;
    font-weight: 600;
}

:deep(.el-sub-menu__title) {
    height: 48px;
    line-height: 48px;
    margin: 4px 8px;
    border-radius: 6px;
    transition: all 0.3s;
}

:deep(.el-sub-menu__title:hover) {
    background-color: rgba(64, 158, 255, 0.1) !important;
}

:deep(.el-sub-menu.is-active > .el-sub-menu__title) {
    color: #409EFF;
    font-weight: 600;
}

.nest-menu :deep(.el-menu-item) {
    min-width: 180px;
    padding-left: 50px !important;
}
</style>

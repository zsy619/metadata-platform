<template>
    <div v-if="!item.meta?.hidden">
        <template v-if="hasOneShowingChild(item.children, item) && (!onlyOneChild.children || onlyOneChild.noShowingChildren) && !item.meta?.alwaysShow">
            <app-link v-if="onlyOneChild.meta" :to="resolvePath(onlyOneChild.path)">
                <el-menu-item :index="resolvePath(onlyOneChild.path)" :class="{ 'submenu-title-noDropdown': !isNest }">
                    <font-awesome-icon :icon="parseIcon(onlyOneChild.meta.icon || item.meta?.icon)" v-if="onlyOneChild.meta.icon || item.meta?.icon" class="menu-icon" />
                    <template #title>
                        <span class="menu-title">{{ onlyOneChild.meta?.title }}</span>
                    </template>
                </el-menu-item>
            </app-link>
        </template>
        <el-sub-menu v-else :index="resolvePath(item.path)" popper-append-to-body>
            <template #title>
                <font-awesome-icon :icon="parseIcon(item.meta?.icon)" v-if="item.meta?.icon" class="menu-icon" />
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

const parseIcon = (iconName: string) => {
    if (!iconName) return null
    
    if (iconName.startsWith('fa-solid ') || iconName.startsWith('fa-regular ') || iconName.startsWith('fa-brands ')) {
        return iconName
    }
    
    if (iconName.startsWith('fa-')) {
        return `fa-solid ${iconName}`
    }
    
    return null
}

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

    if (routePath.startsWith('/')) {
        return routePath
    }

    const base = props.basePath.replace(/\/+$/, '')
    const path = routePath.replace(/^\/+/, '')

    return base ? `${base}/${path}` : `/${path}`
}
</script>
<style scoped>
.menu-icon {
    font-size: 18px;
    width: 18px;
    height: 18px;
    margin-right: 8px;
    vertical-align: middle;
    display: inline-flex;
    align-items: center;
    justify-content: center;
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
    background-color: var(--el-color-primary-light-9) !important;
}

:deep(.el-menu-item.is-active) {
    background-color: var(--el-color-primary-light-9) !important;
    color: var(--el-color-primary);
    border-right: 3px solid var(--el-color-primary);
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
    background-color: var(--el-color-primary-light-9) !important;
}

:deep(.el-sub-menu.is-active > .el-sub-menu__title) {
    color: var(--el-color-primary);
    font-weight: 600;
}
</style>

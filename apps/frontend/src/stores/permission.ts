import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
// import { constantRoutes } from '@/router/routes' // 假设我们有静态路由定义

export const usePermissionStore = defineStore('permission', () => {
    const routes = ref<RouteRecordRaw[]>([])
    const addRoutes = ref<RouteRecordRaw[]>([])

    const setRoutes = (newRoutes: RouteRecordRaw[]) => {
        addRoutes.value = newRoutes
        // routes.value = constantRoutes.concat(newRoutes)
        routes.value = newRoutes // 暂时简单处理
    }

    const generateRoutes = async (_roles: string[]) => {
        // 这里可以添加基于角色的路由生成逻辑
        // 目前返回空数组，表示没有动态路由，或者全部路由已静态定义
        return new Promise<RouteRecordRaw[]>((resolve) => {
            const accessedRoutes: RouteRecordRaw[] = []
            setRoutes(accessedRoutes)
            resolve(accessedRoutes)
        })
    }

    return {
        routes,
        addRoutes,
        setRoutes,
        generateRoutes
    }
})

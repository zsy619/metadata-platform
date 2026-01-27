import { useUserStore } from '@/stores/user'
import storage from '@/utils/storage'
import type { Router } from 'vue-router'

export function setupGuards(router: Router) {
    // 路由前置守卫
    router.beforeEach(async (to, _from, next) => {
        // 设置页面标题
        document.title = (to.meta.title as string) || '元数据管理平台'

        // 登录检查
        const token = storage.get('token')
        const whiteList = ['/login']

        if (token) {
            const userStore = useUserStore()

            // 如果有 token 但没有用户信息，则获取用户信息
            if (!userStore.userInfo) {
                try {
                    await userStore.getInfo()
                } catch (error) {
                    // 获取用户信息失败，可能是 token 过期，清除 token 并跳转到登录页
                    console.error('获取用户信息失败:', error)
                    await userStore.logout()
                    next(`/login?redirect=${to.fullPath}`)
                    return
                }
            }

            if (to.path === '/login') {
                next('/')
            } else {
                next()
            }
        } else {
            if (whiteList.includes(to.path)) {
                next()
            } else {
                next(`/login?redirect=${to.fullPath}`)
            }
        }
    })

    // 路由后置守卫
    router.afterEach(() => {
        // 关闭进度条
        // NProgress.done()
    })

    // 错误处理
    router.onError((error) => {
        // NProgress.done()
        console.error('路由错误:', error)
    })
}

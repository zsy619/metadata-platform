import { getUserProfile, login as loginApi, logout as logoutApi } from '@/api/auth'
import type { LoginRequest, User } from '@/types/user'
import storage from '@/utils/storage'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useUserStore = defineStore('user', () => {
    // 状态
    const token = ref<string>(storage.get('token') || '')
    const refreshToken = ref<string>(storage.get('refreshToken') || '')
    // 从localStorage恢复用户信息
    const savedUserInfo = storage.get('userInfo')
    const userInfo = ref<User | null>(savedUserInfo ? (typeof savedUserInfo === 'string' ? JSON.parse(savedUserInfo) : savedUserInfo) : null)

    // Getters
    const isLoggedIn = computed(() => !!token.value)
    const userName = computed(() => {
        if (!userInfo.value) return ''
        return userInfo.value.name || userInfo.value.account || '管理员'
    })
    const userRoles = computed(() => userInfo.value?.roles || [])
    const avatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png')

    // Actions

    /**
     * 登录
     */
    const login = async (loginForm: LoginRequest) => {
        try {
            const res = await loginApi(loginForm)
            // 兼容直接返回token或返回包装对象的情况
            const data = (res as any).data || res
            const accessToken = data.token || data.access_token
            const user = data.user || {}

            // 保存状态
            token.value = accessToken
            userInfo.value = user

            // 持久化
            storage.set('token', accessToken)
            // 持久化用户信息，用于请求头传递
            storage.set('userInfo', JSON.stringify(user))
            // 持久化租户ID
            if (user.tenant_id) {
                storage.setTenantID(user.tenant_id)
            }
            // storage.set('refreshToken', res.data.refreshToken) // 如果后端返回 refreshToken

            return Promise.resolve(res)
        } catch (error) {
            return Promise.reject(error)
        }
    }

    /**
     * 获取用户信息
     */
    const getInfo = async () => {
        if (!token.value) return Promise.reject(new Error('No token'))

        try {
            const res = await getUserProfile()
            userInfo.value = res.data
            // 持久化用户信息，用于请求头传递
            storage.set('userInfo', JSON.stringify(res.data))
            return Promise.resolve(res.data)
        } catch (error) {
            return Promise.reject(error)
        }
    }

    /**
     * 退出登录
     */
    const logout = async () => {
        try {
            if (token.value) {
                await logoutApi()
            }
        } catch (error) {
            console.warn('Logout failed:', error)
        } finally {
            // 清除状态
            token.value = ''
            refreshToken.value = ''
            userInfo.value = null

            // 清除持久化
            storage.remove('token')
            storage.remove('refreshToken')
            storage.remove('userInfo')
            storage.remove('tenantID')
        }
    }

    /**
     * 重置Token (用于刷新Token后更新)
     */
    const setToken = (newToken: string) => {
        token.value = newToken
        storage.set('token', newToken)
    }

    return {
        token,
        refreshToken,
        userInfo,
        isLoggedIn,
        userName,
        userRoles,
        avatar,
        login,
        getInfo,
        logout,
        setToken
    }
})

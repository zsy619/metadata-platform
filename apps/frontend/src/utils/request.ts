import router from '@/router'
import { NProgress } from '@/utils/progress'
import storage from '@/utils/storage'
import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

// 请求计数器，用于控制 NProgress
let requestCount = 0
const startLoading = () => {
    if (requestCount === 0) {
        NProgress.start()
    }
    requestCount++
}
const endLoading = () => {
    if (requestCount <= 0) return
    requestCount--
    if (requestCount === 0) {
        NProgress.done()
    }
}

// 创建Axios实例
const request: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/',
    timeout: Number(import.meta.env.VITE_REQUEST_TIMEOUT) || 10000,
    headers: {
        'Content-Type': 'application/json;charset=utf-8'
    }
})

// 请求拦截器
request.interceptors.request.use(
    (config: any) => {
        startLoading()
        // 获取token
        const token = storage.get('token')
        if (token) {
            config.headers = config.headers || {}
            config.headers.Authorization = `Bearer ${token}`
        }

        // 获取用户信息并添加到请求头
        // 从localStorage获取用户信息（登录后存储）
        const userInfoStr = storage.get('userInfo')
        if (userInfoStr) {
            try {
                const userInfo = typeof userInfoStr === 'string' ? JSON.parse(userInfoStr) : userInfoStr
                // 添加用户ID和账户到请求头
                if (userInfo.id) {
                    config.headers['X-User-ID'] = userInfo.id
                }
                if (userInfo.account) {
                    config.headers['X-User-Account'] = userInfo.account
                }
            } catch (e) {
                console.warn('解析用户信息失败:', e)
            }
        }

        // 添加租户ID到请求头
        config.headers['X-Tenant-ID'] = storage.getTenantID()

        return config
    },
    (error: AxiosError) => {
        endLoading()
        console.error('请求错误:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    (response: AxiosResponse) => {
        endLoading()
        const { data } = response

        // 如果是二进制数据或文件下载，直接返回
        if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
            return response
        }

        // 处理业务错误
        if (data && typeof data === 'object' && data.code && data.code !== 200) {
            const msg = data.message || '请求失败'
            const url = response.config.url
            console.error(`API 业务错误 [${url}]:`, data)
            ElMessage.error(msg)
            return Promise.reject(new Error(msg))
        }

        return data
    },
    (error: any) => {
        endLoading()
        console.error('响应错误:', error)

        let errorMessage = '网络请求失败，请稍后重试'

        if (error.response) {
            const { status, data: errorData, config } = error.response
            const data = errorData as any
            const isLoginRequest = config.url?.includes('/api/auth/login')

            switch (status) {
                case 400:
                    errorMessage = data?.error || data?.message || '请求参数错误'
                    break
                case 401:
                    // 如果是登录接口返回401，说明是用户名密码错误，由业务层处理
                    if (isLoginRequest) {
                        errorMessage = data?.message || '用户名或密码错误'
                    } else {
                        errorMessage = '登录已过期，请重新登录'
                        // 清除token并跳转到登录页
                        storage.remove('token')
                        router.push('/login')
                    }
                    break
                case 403:
                    errorMessage = '没有权限访问该资源'
                    router.push('/403')
                    break
                case 404:
                    errorMessage = '请求的资源不存在'
                    break
                case 500:
                    errorMessage = '服务器内部错误'
                    router.push('/500')
                    break
                case 502:
                case 503:
                case 504:
                    errorMessage = '服务不可用或系统维护中'
                    router.push('/503')
                    break
                default:
                    errorMessage = data?.message || `请求失败，状态码：${status}`
            }
        } else if (error.code === 'ERR_NETWORK') {
            errorMessage = '网络连接失败'
            router.push('/network-error')
        } else if (error.request) {
            errorMessage = '网络请求超时，请检查网络连接'
        }

        ElMessage.error(errorMessage)
        return Promise.reject(error)
    }
)

export default request

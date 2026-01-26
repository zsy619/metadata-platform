/**
 * HTTP请求工具
 */
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import router from '@/router'

// 创建Axios实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 获取token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response
    
    // 如果是二进制数据或文件下载，直接返回
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response
    }
    
    // 处理业务错误
    if (data.code && data.code !== 200) {
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message || '请求失败'))
    }
    
    return data
  },
  (error: AxiosError) => {
    console.error('响应错误:', error)
    
    let errorMessage = '网络请求失败，请稍后重试'
    
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 400:
          errorMessage = data.message || '请求参数错误'
          break
        case 401:
          errorMessage = '登录已过期，请重新登录'
          // 清除token并跳转到登录页
          localStorage.removeItem('token')
          router.push('/login')
          break
        case 403:
          errorMessage = '没有权限访问该资源'
          break
        case 404:
          errorMessage = '请求的资源不存在'
          break
        case 500:
          errorMessage = '服务器内部错误'
          break
        case 502:
          errorMessage = '网关错误'
          break
        case 503:
          errorMessage = '服务暂时不可用'
          break
        case 504:
          errorMessage = '网关超时'
          break
        default:
          errorMessage = data.message || `请求失败，状态码：${status}`
      }
    } else if (error.request) {
      errorMessage = '网络请求超时，请检查网络连接'
    }
    
    ElMessage.error(errorMessage)
    return Promise.reject(error)
  }
)

export default request

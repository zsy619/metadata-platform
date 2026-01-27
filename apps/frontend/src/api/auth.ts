import type { LoginRequest, LoginResponse } from '@/types/user'
import request from '@/utils/request'

// 登录API
export const login = (data: LoginRequest) => {
    return request.post<any, LoginResponse>('/api/auth/login', data)
}

// 退出登录API
export const logout = () => {
    return request.post('/api/auth/logout')
}

// 获取验证码API
export const getCaptchaApi = () => {
    return request.get('/api/auth/captcha')
}

// 注册API
export const registerApi = (user: any) => {
    return request.post('/api/auth/register', user)
}

// 租户相关API
// 获取所有租户API
export const getAllTenantsApi = () => {
    return request.get('/api/tenants')
}

// 根据ID获取租户API
export const getTenantByIdApi = (id: string) => {
    return request.get(`/api/tenants/${id}`)
}

// 创建租户API
export const createTenantApi = (tenant: any) => {
    return request.post('/api/tenants', tenant)
}

// 更新租户API
export const updateTenantApi = (id: string, tenant: any) => {
    return request.put(`/api/tenants/${id}`, tenant)
}

// 删除租户API
export const deleteTenantApi = (id: string) => {
    return request.delete(`/api/tenants/${id}`)
}

// 获取当前用户信息API (已废弃，请使用 getUserProfile)
export const getCurrentUserApi = () => {
    return getUserProfile()
}

// 获取个人资料API
export const getUserProfile = () => {
    return request.get('/api/auth/profile').catch(() => {
        // 如果后端接口未实现，返回模拟数据
        return Promise.resolve({
            data: {
                id: '1',
                account: 'admin',
                name: '管理员',
                email: 'admin@example.com',
                avatar: 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png',
                roles: ['admin']
            }
        })
    })
}

// 获取所有用户API
export const getAllUsersApi = () => {
    return request.get('/api/users')
}

// 更新用户信息API
export const updateUserApi = (user: any) => {
    return request.put('/api/users/me', user)
}

// 修改密码API
export const updateUserPassword = (data: any) => {
    return request.post('/api/auth/password', data)
}

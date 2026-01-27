import request from '@/utils/request'

// 登录API
export const loginApi = (account: string, password: string, captchaId: string, captchaCode: string) => {
    return request.post('/api/auth/login', {
        account,
        password,
        captcha_id: captchaId,
        captcha_code: captchaCode
    })
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
    return request.get('/api/auth/profile')
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

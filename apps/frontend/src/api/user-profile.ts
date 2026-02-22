import request from '@/utils/request'

// ────── 用户档案 ──────
export const getUserProfile = async (userId: string): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/profile`, method: 'get' })

export const upsertUserProfile = async (userId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/profile`, method: 'put', data })

// ────── 地址簿 ──────
export const getUserAddresses = async (userId: string): Promise<any[]> =>
    request({ url: `/api/sso/user/${userId}/addresses`, method: 'get' })

export const createUserAddress = async (userId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/addresses`, method: 'post', data })

export const updateUserAddress = async (userId: string, addrId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/addresses/${addrId}`, method: 'put', data })

export const deleteUserAddress = async (userId: string, addrId: string): Promise<void> =>
    request({ url: `/api/sso/user/${userId}/addresses/${addrId}`, method: 'delete' })

export const setDefaultAddress = async (userId: string, addrId: string): Promise<void> =>
    request({ url: `/api/sso/user/${userId}/addresses/${addrId}/default`, method: 'put' })

// ────── 联系方式 ──────
export const getUserContacts = async (userId: string): Promise<any[]> =>
    request({ url: `/api/sso/user/${userId}/contacts`, method: 'get' })

export const createUserContact = async (userId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/contacts`, method: 'post', data })

export const updateUserContact = async (userId: string, contactId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/contacts/${contactId}`, method: 'put', data })

export const deleteUserContact = async (userId: string, contactId: string): Promise<void> =>
    request({ url: `/api/sso/user/${userId}/contacts/${contactId}`, method: 'delete' })

// ────── 第三方账号 ──────
export const getUserSocials = async (userId: string): Promise<any[]> =>
    request({ url: `/api/sso/user/${userId}/socials`, method: 'get' })

export const bindUserSocial = async (userId: string, data: any): Promise<any> =>
    request({ url: `/api/sso/user/${userId}/socials`, method: 'post', data })

export const unbindUserSocial = async (userId: string, socialId: string): Promise<void> =>
    request({ url: `/api/sso/user/${userId}/socials/${socialId}`, method: 'delete' })

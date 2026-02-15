/**
 * 本地存储封装
 */
class Storage {
    private prefix: string

    constructor(prefix = 'app_') {
        this.prefix = prefix
    }

    private getKey(key: string) {
        return `${this.prefix}${key}`
    }

    /**
     * 设置存储
     * @param key 键
     * @param value 值
     * @param expire 过期时间（秒）
     */
    set(key: string, value: any, expire?: number) {
        const data = {
            value,
            time: Date.now(),
            expire: expire ? Date.now() + expire * 1000 : null
        }
        localStorage.setItem(this.getKey(key), JSON.stringify(data))
    }

    /**
     * 获取存储
     * @param key 键
     * @returns 值
     */
    get<T = any>(key: string): T | null {
        const json = localStorage.getItem(this.getKey(key))
        if (!json) return null

        try {
            const data = JSON.parse(json)
            if (data.expire && data.expire < Date.now()) {
                this.remove(key)
                return null
            }
            return data.value
        } catch (e) {
            return null
        }
    }

    /**
     * 移除存储
     * @param key 键
     */
    remove(key: string) {
        localStorage.removeItem(this.getKey(key))
    }

    /**
     * 清除所有存储
     */
    clear() {
        Object.keys(localStorage).forEach(key => {
            if (key.startsWith(this.prefix)) {
                localStorage.removeItem(key)
            }
        })
    }

    /**
     * 设置Session存储
     */
    setSession(key: string, value: any) {
        sessionStorage.setItem(this.getKey(key), JSON.stringify(value))
    }

    /**
     * 获取Session存储
     */
    getSession<T = any>(key: string): T | null {
        const json = sessionStorage.getItem(this.getKey(key))
        if (!json) return null
        try {
            return JSON.parse(json)
        } catch (e) {
            return null
        }
    }

    /**
     * 移除Session存储
     */
    removeSession(key: string) {
        sessionStorage.removeItem(this.getKey(key))
    }

    /**
     * 清除Session存储
     */
    clearSession() {
        Object.keys(sessionStorage).forEach(key => {
            if (key.startsWith(this.prefix)) {
                sessionStorage.removeItem(key)
            }
        })
    }

    /**
     * 设置租户ID
     * @param tenantID 租户ID
     */
    setTenantID(tenantID: string) {
        this.set('tenantID', tenantID)
    }

    /**
     * 获取租户ID
     * @returns 租户ID，默认返回 "1"
     */
    getTenantID(): string {
        return this.get<string>('tenantID') || '1'
    }
}

export default new Storage()

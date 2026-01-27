/**
 * 常用验证工具函数
 */

/**
 * 验证手机号
 */
export const isMobile = (path: string): boolean => {
    const reg = /^1[3-9]\d{9}$/
    return reg.test(path)
}

/**
 * 验证邮箱
 */
export const isEmail = (path: string): boolean => {
    const reg = /^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$/
    return reg.test(path)
}

/**
 * 验证身份证
 */
export const isIdCard = (path: string): boolean => {
    const reg = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/
    return reg.test(path)
}

/**
 * 验证URL
 */
export const isURL = (path: string): boolean => {
    const reg = /^(https?|ftp):\/\/([a-zA-Z0-9.-]+(:[a-zA-Z0-9.&%$-]+)*@)*((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)){3}|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.(com|edu|gov|int|mil|net|org|biz|arpa|info|name|pro|aero|coop|museum|[a-zA-Z]{2}))(:[0-9]+)*(\/($|[a-zA-Z0-9.,?'\\+&%$#=~_-]+))*$/
    return reg.test(path)
}

/**
 * 验证密码强度 (8-20位，包含大小写字母+数字+特殊字符)
 */
export const isStrongPassword = (val: string): boolean => {
    const reg = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/
    return reg.test(val)
}

/**
 * 格式化工具函数
 */

/**
 * 格式化日期
 * @param dateStr 日期字符串/时间戳/Date对象
 * @param pattern 格式化模式，默认 'yyyy-MM-dd HH:mm:ss'
 */
export const formatDate = (dateStr: string | number | Date | undefined | null, pattern = 'yyyy-MM-dd HH:mm:ss'): string => {
    if (!dateStr || dateStr === '' || dateStr === '0001-01-01 00:00:00') return '-'

    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return String(dateStr)

    const opt: Record<string, string> = {
        'y+': date.getFullYear().toString(),
        'M+': (date.getMonth() + 1).toString(),
        'd+': date.getDate().toString(),
        'H+': date.getHours().toString(),
        'm+': date.getMinutes().toString(),
        's+': date.getSeconds().toString()
    }

    let res = pattern
    for (const k in opt) {
        const ret = new RegExp('(' + k + ')').exec(res)
        if (ret) {
            res = res.replace(ret[1], (ret[1].length === 1) ? (opt[k]) : (opt[k].padStart(ret[1].length, '0')))
        }
    }
    return res
}

/**
 * 格式化金额
 * @param amount 金额
 * @param decimals 小数位数
 * @param symbol 货币符号
 */
export const formatMoney = (amount: number | string, decimals = 2, symbol = '¥'): string => {
    if (amount === undefined || amount === null) return '-'
    const num = Number(amount)
    if (isNaN(num)) return '-'
    return symbol + num.toFixed(decimals).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

/**
 * 格式化百分比
 * @param val 数值 (0.1234)
 * @param decimals 小数位数
 */
export const formatPercent = (val: number | string, decimals = 2): string => {
    if (val === undefined || val === null) return '-'
    const num = Number(val)
    if (isNaN(num)) return '-'
    return (num * 100).toFixed(decimals) + '%'
}

/**
 * 格式化文件大小
 * @param bytes 字节数
 */
export const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

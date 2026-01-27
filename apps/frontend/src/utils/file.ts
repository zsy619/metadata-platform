/**
 * 文件处理工具
 */

/**
 * 获取文件后缀名
 */
export const getFileExtension = (filename: string): string => {
    return filename.slice((filename.lastIndexOf('.') - 1 >>> 0) + 2)
}

/**
 * 检查文件类型是否为图片
 */
export const isImage = (filename: string): boolean => {
    const ext = getFileExtension(filename).toLowerCase()
    return ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg'].includes(ext)
}

/**
 * 检查文件是否为Excel
 */
export const isExcel = (filename: string): boolean => {
    const ext = getFileExtension(filename).toLowerCase()
    return ['xlsx', 'xls', 'csv'].includes(ext)
}

/**
 * 下载文件流
 * @param data 文件流数据
 * @param filename 文件名
 */

export const downloadFile = (data: Blob, filename: string) => {
    const url = window.URL.createObjectURL(new Blob([data]))
    const link = document.createElement('a')
    link.style.display = 'none'
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
}

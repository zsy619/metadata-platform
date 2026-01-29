import { ElMessageBox } from 'element-plus'

/**
 * 通用确认对话框
 * @param content 提示内容
 * @param title 标题
 * @param type 消息类型
 * @returns Promise
 */
export const showConfirm = (
    content: string,
    title: string = '提示',
    type: 'warning' | 'info' | 'success' | 'error' = 'warning'
) => {
    return ElMessageBox.confirm(content, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: type,
        draggable: true,
    })
}

/**
 * 删除确认对话框
 * @param content 删除提示内容，默认为"确定要删除该数据吗？此操作不可恢复！"
 * @returns Promise
 */
export const showDeleteConfirm = (content: string = '确定要删除该数据吗？此操作不可恢复！') => {
    return showConfirm(content, '删除确认', 'error')
}

/**
 * 带有输入框的确认对话框
 * @param title 标题
 * @param message 提示信息
 * @param placeholder 输入框占位符
 * @returns Promise<{ value: string, action: string }>
 */
export const showPrompt = (title: string, message: string, placeholder: string = '') => {
    return ElMessageBox.prompt(message, title, {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPlaceholder: placeholder,
    })
}

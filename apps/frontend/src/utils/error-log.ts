import { version } from '../../package.json'

interface ErrorLog {
    message: string
    stack?: string
    url: string
    time: string
    userId?: string
    version: string
    type: 'js' | 'vue' | 'api'
}

/**
 * 错误上报工具
 */
export const errorLog = {
    // 错误日志队列
    logs: [] as ErrorLog[],

    /**
     * 添加错误日志
     * @param error 错误对象
     * @param type 错误类型
     */
    add(error: any, type: 'js' | 'vue' | 'api' = 'js') {
        const log: ErrorLog = {
            message: error.message || String(error),
            stack: error.stack,
            url: window.location.href,
            time: new Date().toISOString(),
            version: version,
            type
        }

        // 开发环境主要在控制台输出
        if (import.meta.env.MODE === 'development') {
            console.group('Error Log')
            console.log('Type:', type)
            console.log('Message:', log.message)
            console.log('URL:', log.url)
            console.log('Stack:', log.stack)
            console.groupEnd()
        }

        this.logs.push(log)
        // 这一步可以触发发送到后端的逻辑
        this.send()
    },

    /**
     * 发送错误日志到后端
     * (这里是模拟实现，实际需要后端提供接口)
     */
    send() {
        // 如果积累了超过5条，或者定时发送
        // 此处仅为示例，实际可调用 request.post('/api/monitor/error-log', this.logs)
        // console.log('Sending error logs to backend...', this.logs)
    }
}

/**
 * Vue 错误处理器
 */
export const setupErrorHandle = (app: any) => {
    app.config.errorHandler = (err: any, _vm: any, info: string) => {
        errorLog.add(err, 'vue')
        console.error('Vue Error:', err, info)
    }
}

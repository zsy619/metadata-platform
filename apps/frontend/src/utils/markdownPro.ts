/**
 * Markdown 编辑器 Pro 工具库
 * 包含文档对比、版本管理、批注系统等企业级功能
 */

/**
 * 文档版本接口
 */
export interface DocumentVersion {
    id: string
    version: number
    content: string
    author: string
    createdAt: string
    comment?: string
    changeType?: 'create' | 'update' | 'restore'
}

/**
 * 批注接口
 */
export interface Annotation {
    id: string
    userId: string
    userName: string
    content: string
    position: {
        line: number
        column: number
        text: string
    }
    createdAt: string
    status: 'pending' | 'resolved' | 'rejected'
    replies: AnnotationReply[]
}

/**
 * 批注回复接口
 */
export interface AnnotationReply {
    id: string
    userId: string
    userName: string
    content: string
    createdAt: string
}

/**
 * 文档对比结果
 */
export interface DocumentDiff {
    added: Array<{ line: number; content: string }>
    removed: Array<{ line: number; content: string }>
    modified: Array<{
        line: number
        oldContent: string
        newContent: string
    }>
    summary: {
        addedLines: number
        removedLines: number
        modifiedLines: number
        totalChanges: number
    }
}

/**
 * 文档统计信息
 */
export interface DocumentStatistics {
    wordCount: number
    characterCount: number
    paragraphCount: number
    lineCount: number
    headingCount: {
        h1: number
        h2: number
        h3: number
        h4: number
        h5: number
        h6: number
    }
    codeBlockCount: number
    linkCount: number
    imageCount: number
    tableCount: number
    readingTime: number
    writingTime: number
}

/**
 * 快捷键配置
 */
export interface ShortcutConfig {
    key: string
    description: string
    action: string
    category: string
}

/**
 * 快捷键列表
 */
export const shortcuts: ShortcutConfig[] = [
    { key: 'Ctrl+S', description: '保存文档', action: 'save', category: '文件操作' },
    { key: 'Ctrl+N', description: '新建文档', action: 'new', category: '文件操作' },
    { key: 'Ctrl+Z', description: '撤销', action: 'undo', category: '编辑' },
    { key: 'Ctrl+Y', description: '重做', action: 'redo', category: '编辑' },
    { key: 'Ctrl+B', description: '粗体', action: 'bold', category: '格式' },
    { key: 'Ctrl+I', description: '斜体', action: 'italic', category: '格式' },
    { key: 'Ctrl+K', description: '插入链接', action: 'link', category: '插入' },
    { key: 'Ctrl+Shift+I', description: '插入图片', action: 'image', category: '插入' },
    { key: 'Ctrl+Shift+T', description: '插入表格', action: 'table', category: '插入' },
    { key: 'Ctrl+Shift+C', description: '插入代码块', action: 'code', category: '插入' },
    { key: 'Ctrl+/', description: '切换注释', action: 'comment', category: '编辑' },
    { key: 'F1', description: '显示快捷键', action: 'help', category: '帮助' },
]

/**
 * 计算文档统计信息
 */
export function calculateStatistics(content: string): DocumentStatistics {
    const lines = content.split('\n')
    
    // 字数统计
    const chineseChars = (content.match(/[\u4e00-\u9fa5]/g) || []).length
    const englishWords = (content.match(/\b\w+\b/g) || []).length
    const wordCount = chineseChars + englishWords
    
    // 段落统计
    const paragraphs = content.split(/\n\n+/).filter(p => p.trim())
    
    // 标题统计
    const headings = { h1: 0, h2: 0, h3: 0, h4: 0, h5: 0, h6: 0 }
    lines.forEach(line => {
        const match = line.match(/^(#{1,6})\s+/)
        if (match) {
            const level = match[1].length
            headings[`h${level}` as keyof typeof headings]++
        }
    })
    
    // 代码块统计
    const codeBlocks = (content.match(/```[\s\S]*?```/g) || []).length
    
    // 链接统计
    const links = (content.match(/\[.*?\]\(.*?\)/g) || []).length
    
    // 图片统计
    const images = (content.match(/!\[.*?\]\(.*?\)/g) || []).length
    
    // 表格统计
    const tables = (content.match(/\|.*?\|.*?\|/g) || []).length
    
    return {
        wordCount,
        characterCount: content.length,
        paragraphCount: paragraphs.length,
        lineCount: lines.length,
        headingCount: headings,
        codeBlockCount: codeBlocks,
        linkCount: links,
        imageCount: images,
        tableCount: tables,
        readingTime: Math.ceil(wordCount / 300),
        writingTime: Math.ceil(wordCount / 60) // 假设写作速度 60 字/分钟
    }
}

/**
 * 文档对比（简化版 diff 算法）
 */
export function compareDocuments(oldContent: string, newContent: string): DocumentDiff {
    const oldLines = oldContent.split('\n')
    const newLines = newContent.split('\n')
    
    const added: Array<{ line: number; content: string }> = []
    const removed: Array<{ line: number; content: string }> = []
    const modified: Array<{ line: number; oldContent: string; newContent: string }> = []
    
    const maxLength = Math.max(oldLines.length, newLines.length)
    
    for (let i = 0; i < maxLength; i++) {
        const oldLine = oldLines[i]
        const newLine = newLines[i]
        
        if (oldLine === undefined) {
            added.push({ line: i + 1, content: newLine })
        } else if (newLine === undefined) {
            removed.push({ line: i + 1, content: oldLine })
        } else if (oldLine !== newLine) {
            modified.push({
                line: i + 1,
                oldContent: oldLine,
                newContent: newLine
            })
        }
    }
    
    return {
        added,
        removed,
        modified,
        summary: {
            addedLines: added.length,
            removedLines: removed.length,
            modifiedLines: modified.length,
            totalChanges: added.length + removed.length + modified.length
        }
    }
}

/**
 * 提取文档大纲
 */
export interface OutlineItem {
    id: string
    level: number
    text: string
    line: number
    children?: OutlineItem[]
}

export function extractOutline(content: string): OutlineItem[] {
    const lines = content.split('\n')
    const outline: OutlineItem[] = []
    const stack: Array<{ item: OutlineItem; level: number }> = []
    
    lines.forEach((line, index) => {
        const match = line.match(/^(#{1,6})\s+(.+)$/)
        if (match) {
            const level = match[1].length
            const text = match[2].trim()
            const item: OutlineItem = {
                id: `heading-${index}`,
                level,
                text,
                line: index + 1
            }
            
            // 找到合适的父节点
            while (stack.length > 0 && stack[stack.length - 1].level >= level) {
                stack.pop()
            }
            
            if (stack.length === 0) {
                outline.push(item)
            } else {
                const parent = stack[stack.length - 1].item
                if (!parent.children) {
                    parent.children = []
                }
                parent.children.push(item)
            }
            
            stack.push({ item, level })
        }
    })
    
    return outline
}

/**
 * 生成文档水印
 */
export function generateWatermark(text: string = '内部文档', options?: {
    fontSize?: number
    color?: string
    opacity?: number
    rotate?: number
    spacing?: number
}): string {
    const {
        fontSize = 24,
        color = '#000000',
        opacity = 0.1,
        rotate = -30,
        spacing = 200
    } = options || {}
    
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    
    if (!ctx) {
        return ''
    }
    
    canvas.width = 400
    canvas.height = 400
    
    ctx.clearRect(0, 0, canvas.width, canvas.height)
    ctx.save()
    
    // 移动到中心
    ctx.translate(canvas.width / 2, canvas.height / 2)
    ctx.rotate((rotate * Math.PI) / 180)
    
    // 设置样式
    ctx.fillStyle = color
    ctx.globalAlpha = opacity
    ctx.font = `${fontSize}px Arial`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    
    // 绘制文字
    ctx.fillText(text, 0, 0)
    
    ctx.restore()
    
    return canvas.toDataURL('image/png')
}

/**
 * 文档保护配置
 */
export interface DocumentProtection {
    isProtected: boolean
    password?: string
    isReadOnly: boolean
    allowPrint: boolean
    allowCopy: boolean
    allowExport: boolean
}

/**
 * 验证文档密码
 */
export function verifyDocumentPassword(content: string, password: string): boolean {
    // 简单实现，实际项目中应该使用更安全的加密方式
    const hash = btoa(password)
    return content.includes(`<!-- password:${hash} -->`)
}

/**
 * 设置文档密码
 */
export function setDocumentPassword(content: string, password: string): string {
    const hash = btoa(password)
    return `<!-- password:${hash} -->\n${content}`
}

/**
 * 批注管理
 */
export class AnnotationManager {
    private annotations: Annotation[] = []
    
    /**
     * 添加批注
     */
    addAnnotation(annotation: Omit<Annotation, 'id' | 'createdAt' | 'status' | 'replies'>): Annotation {
        const newAnnotation: Annotation = {
            ...annotation,
            id: `annotation-${Date.now()}`,
            createdAt: new Date().toISOString(),
            status: 'pending',
            replies: []
        }
        
        this.annotations.push(newAnnotation)
        return newAnnotation
    }
    
    /**
     * 获取批注列表
     */
    getAnnotations(): Annotation[] {
        return this.annotations
    }
    
    /**
     * 回复批注
     */
    replyAnnotation(annotationId: string, reply: Omit<AnnotationReply, 'id' | 'createdAt'>): boolean {
        const annotation = this.annotations.find(a => a.id === annotationId)
        if (!annotation) {
            return false
        }
        
        const newReply: AnnotationReply = {
            ...reply,
            id: `reply-${Date.now()}`,
            createdAt: new Date().toISOString()
        }
        
        annotation.replies.push(newReply)
        return true
    }
    
    /**
     * 解决批注
     */
    resolveAnnotation(annotationId: string): boolean {
        const annotation = this.annotations.find(a => a.id === annotationId)
        if (!annotation) {
            return false
        }
        
        annotation.status = 'resolved'
        return true
    }
    
    /**
     * 删除批注
     */
    deleteAnnotation(annotationId: string): boolean {
        const index = this.annotations.findIndex(a => a.id === annotationId)
        if (index === -1) {
            return false
        }
        
        this.annotations.splice(index, 1)
        return true
    }
    
    /**
     * 清除所有批注
     */
    clearAnnotations(): void {
        this.annotations = []
    }
}

/**
 * 版本历史管理
 */
export class VersionManager {
    private versions: DocumentVersion[] = []
    private maxVersions: number = 50
    
    constructor(maxVersions: number = 50) {
        this.maxVersions = maxVersions
    }
    
    /**
     * 保存版本
     */
    saveVersion(content: string, author: string, comment?: string): DocumentVersion {
        const version: DocumentVersion = {
            id: `version-${Date.now()}`,
            version: this.versions.length + 1,
            content,
            author,
            createdAt: new Date().toISOString(),
            comment,
            changeType: this.versions.length === 0 ? 'create' : 'update'
        }
        
        this.versions.push(version)
        
        // 限制版本数量
        if (this.versions.length > this.maxVersions) {
            this.versions.shift()
        }
        
        return version
    }
    
    /**
     * 获取版本列表
     */
    getVersions(): DocumentVersion[] {
        return this.versions
    }
    
    /**
     * 获取特定版本
     */
    getVersion(versionNumber: number): DocumentVersion | undefined {
        return this.versions.find(v => v.version === versionNumber)
    }
    
    /**
     * 恢复版本
     */
    restoreVersion(versionNumber: number): DocumentVersion | undefined {
        const version = this.getVersion(versionNumber)
        if (!version) {
            return undefined
        }
        
        // 创建恢复版本
        const restoredVersion: DocumentVersion = {
            id: `version-${Date.now()}`,
            version: this.versions.length + 1,
            content: version.content,
            author: version.author,
            createdAt: new Date().toISOString(),
            comment: `恢复到版本 v${version.version}`,
            changeType: 'restore'
        }
        
        this.versions.push(restoredVersion)
        return restoredVersion
    }
    
    /**
     * 对比版本
     */
    compareVersions(version1: number, version2: number): DocumentDiff | null {
        const v1 = this.getVersion(version1)
        const v2 = this.getVersion(version2)
        
        if (!v1 || !v2) {
            return null
        }
        
        return compareDocuments(v1.content, v2.content)
    }
    
    /**
     * 清除版本历史
     */
    clearVersions(): void {
        this.versions = []
    }
}

/**
 * 导出为 PDF（带水印）
 */
export async function exportToPdfWithWatermark(
    content: string,
    title: string,
    watermarkText?: string
): Promise<void> {
    const { exportToPdf } = await import('./documentExport')
    
    if (watermarkText) {
        const watermark = generateWatermark(watermarkText)
        const html = `
            <div style="position: relative;">
                <div style="background-image: url(${watermark}); opacity: 0.1; position: absolute; top: 0; left: 0; right: 0; bottom: 0; z-index: 0;"></div>
                <div style="position: relative; z-index: 1;">
                    ${content}
                </div>
            </div>
        `
        exportToPdf(html, title)
    } else {
        exportToPdf(content, title)
    }
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

/**
 * 格式化时间
 */
export function formatTime(seconds: number): string {
    const minutes = Math.floor(seconds / 60)
    const secs = seconds % 60
    
    if (minutes > 0) {
        return `${minutes}分${secs}秒`
    }
    return `${secs}秒`
}

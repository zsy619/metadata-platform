/**
 * Markdown 编辑器 Ultimate 企业版工具库
 * 包含 AI 助手、权限管理、工作流、审计日志、实时协作等高级功能
 */

/**
 * 用户角色枚举
 */
export enum UserRole {
    ADMIN = 'admin',           // 管理员
    EDITOR = 'editor',         // 编辑者
    REVIEWER = 'reviewer',     // 审核者
    VIEWER = 'viewer',         // 查看者
}

/**
 * 文档权限配置
 */
export interface DocumentPermission {
    documentId: string
    userId: string
    role: UserRole
    canRead: boolean
    canWrite: boolean
    canDelete: boolean
    canShare: boolean
    canApprove: boolean
    grantedAt: string
    grantedBy: string
}

/**
 * 工作流状态
 */
export enum WorkflowStatus {
    DRAFT = 'draft',               // 草稿
    PENDING_REVIEW = 'pending',    // 待审核
    IN_REVIEW = 'in_review',       // 审核中
    APPROVED = 'approved',         // 已通过
    REJECTED = 'rejected',         // 已拒绝
    PUBLISHED = 'published',       // 已发布
    ARCHIVED = 'archived',         // 已归档
}

/**
 * 工作流节点
 */
export interface WorkflowNode {
    id: string
    name: string
    type: 'start' | 'review' | 'approve' | 'end'
    reviewers: string[]
    status: WorkflowStatus
    comment?: string
    reviewedAt?: string
    reviewedBy?: string
}

/**
 * 工作流实例
 */
export interface WorkflowInstance {
    id: string
    documentId: string
    nodes: WorkflowNode[]
    currentNodeId: string
    status: WorkflowStatus
    createdAt: string
    createdBy: string
    completedAt?: string
}

/**
 * 审计日志类型
 */
export enum AuditLogType {
    CREATE = 'create',
    UPDATE = 'update',
    DELETE = 'delete',
    VIEW = 'view',
    SHARE = 'share',
    APPROVE = 'approve',
    REJECT = 'reject',
    PUBLISH = 'publish',
    DOWNLOAD = 'download',
    PRINT = 'print',
}

/**
 * 审计日志记录
 */
export interface AuditLog {
    id: string
    documentId: string
    userId: string
    userName: string
    action: AuditLogType
    details: string
    ipAddress?: string
    userAgent?: string
    timestamp: string
}

/**
 * AI 写作建议
 */
export interface AIWritingSuggestion {
    id: string
    type: 'grammar' | 'style' | 'clarity' | 'tone' | 'spelling'
    severity: 'error' | 'warning' | 'info'
    message: string
    suggestion?: string
    position: {
        start: number
        end: number
        line: number
        column: number
    }
    originalText: string
}

/**
 * 文档标签
 */
export interface DocumentTag {
    id: string
    name: string
    color: string
    count: number
    relatedTags?: string[]
}

/**
 * 文档关联
 */
export interface DocumentRelation {
    id: string
    sourceDocumentId: string
    targetDocumentId: string
    type: 'reference' | 'duplicate' | 'related' | 'parent' | 'child'
    strength: number  // 关联强度 0-1
    createdAt: string
}

/**
 * 实时协作用户
 */
export interface Collaborator {
    userId: string
    userName: string
    color: string
    cursorPosition?: {
        line: number
        column: number
    }
    selection?: {
        start: number
        end: number
    }
    isOnline: boolean
    lastActiveAt: string
}

/**
 * 操作转换 (OT) 操作
 */
export interface OTOperation {
    id: string
    documentId: string
    userId: string
    operation: 'insert' | 'delete' | 'replace'
    position: number
    text?: string
    length?: number
    timestamp: string
}

/**
 * 文档统计趋势
 */
export interface DocumentTrend {
    date: string
    views: number
    edits: number
    comments: number
    shares: number
}

/**
 * 权限管理器
 */
export class PermissionManager {
    private permissions: Map<string, DocumentPermission[]> = new Map()

    /**
     * 设置权限
     */
    setPermission(permission: DocumentPermission): void {
        const docPermissions = this.permissions.get(permission.documentId) || []
        const existingIndex = docPermissions.findIndex(p => p.userId === permission.userId)
        
        if (existingIndex >= 0) {
            docPermissions[existingIndex] = permission
        } else {
            docPermissions.push(permission)
        }
        
        this.permissions.set(permission.documentId, docPermissions)
    }

    /**
     * 检查权限
     */
    checkPermission(
        documentId: string,
        userId: string,
        action: 'read' | 'write' | 'delete' | 'share' | 'approve'
    ): boolean {
        const docPermissions = this.permissions.get(documentId) || []
        const permission = docPermissions.find(p => p.userId === userId)
        
        if (!permission) {
            return false
        }

        switch (action) {
            case 'read':
                return permission.canRead
            case 'write':
                return permission.canWrite
            case 'delete':
                return permission.canDelete
            case 'share':
                return permission.canShare
            case 'approve':
                return permission.canApprove
            default:
                return false
        }
    }

    /**
     * 获取用户角色
     */
    getUserRole(documentId: string, userId: string): UserRole | null {
        const docPermissions = this.permissions.get(documentId) || []
        const permission = docPermissions.find(p => p.userId === userId)
        return permission?.role || null
    }

    /**
     * 移除权限
     */
    removePermission(documentId: string, userId: string): void {
        const docPermissions = this.permissions.get(documentId) || []
        const filtered = docPermissions.filter(p => p.userId !== userId)
        this.permissions.set(documentId, filtered)
    }

    /**
     * 获取文档所有权限
     */
    getDocumentPermissions(documentId: string): DocumentPermission[] {
        return this.permissions.get(documentId) || []
    }
}

/**
 * 工作流管理器
 */
export class WorkflowManager {
    private workflows: Map<string, WorkflowInstance> = new Map()

    /**
     * 创建工作流
     */
    createWorkflow(documentId: string, nodes: WorkflowNode[], createdBy: string): WorkflowInstance {
        const workflow: WorkflowInstance = {
            id: `workflow-${Date.now()}`,
            documentId,
            nodes,
            currentNodeId: nodes[0]?.id || '',
            status: WorkflowStatus.DRAFT,
            createdAt: new Date().toISOString(),
            createdBy
        }
        
        this.workflows.set(workflow.id, workflow)
        return workflow
    }

    /**
     * 提交审核
     */
    submitForReview(workflowId: string): WorkflowInstance | null {
        const workflow = this.workflows.get(workflowId)
        if (!workflow) return null

        workflow.status = WorkflowStatus.PENDING_REVIEW
        if (workflow.nodes.length > 0) {
            workflow.nodes[0].status = WorkflowStatus.PENDING_REVIEW
            workflow.currentNodeId = workflow.nodes[0].id
        }
        
        this.workflows.set(workflowId, workflow)
        return workflow
    }

    /**
     * 审核通过
     */
    approve(
        workflowId: string,
        nodeId: string,
        userId: string,
        comment?: string
    ): WorkflowInstance | null {
        const workflow = this.workflows.get(workflowId)
        if (!workflow) return null

        const node = workflow.nodes.find(n => n.id === nodeId)
        if (!node) return null

        node.status = WorkflowStatus.APPROVED
        node.reviewedBy = userId
        node.reviewedAt = new Date().toISOString()
        node.comment = comment

        // 移动到下一个节点
        const currentIndex = workflow.nodes.findIndex(n => n.id === nodeId)
        if (currentIndex < workflow.nodes.length - 1) {
            workflow.currentNodeId = workflow.nodes[currentIndex + 1].id
            workflow.nodes[currentIndex + 1].status = WorkflowStatus.IN_REVIEW
        } else {
            workflow.status = WorkflowStatus.APPROVED
            workflow.completedAt = new Date().toISOString()
        }

        this.workflows.set(workflowId, workflow)
        return workflow
    }

    /**
     * 审核拒绝
     */
    reject(
        workflowId: string,
        nodeId: string,
        userId: string,
        comment: string
    ): WorkflowInstance | null {
        const workflow = this.workflows.get(workflowId)
        if (!workflow) return null

        const node = workflow.nodes.find(n => n.id === nodeId)
        if (!node) return null

        node.status = WorkflowStatus.REJECTED
        node.reviewedBy = userId
        node.reviewedAt = new Date().toISOString()
        node.comment = comment
        workflow.status = WorkflowStatus.REJECTED

        this.workflows.set(workflowId, workflow)
        return workflow
    }

    /**
     * 获取工作流
     */
    getWorkflow(workflowId: string): WorkflowInstance | undefined {
        return this.workflows.get(workflowId)
    }

    /**
     * 获取文档的工作流
     */
    getWorkflowByDocument(documentId: string): WorkflowInstance | undefined {
        for (const workflow of this.workflows.values()) {
            if (workflow.documentId === documentId) {
                return workflow
            }
        }
        return undefined
    }
}

/**
 * 审计日志管理器
 */
export class AuditLogManager {
    private logs: Map<string, AuditLog[]> = new Map()
    private maxLogsPerDocument: number = 1000

    constructor(maxLogs?: number) {
        if (maxLogs) {
            this.maxLogsPerDocument = maxLogs
        }
    }

    /**
     * 记录日志
     */
    log(
        documentId: string,
        userId: string,
        userName: string,
        action: AuditLogType,
        details: string,
        ipAddress?: string,
        userAgent?: string
    ): AuditLog {
        const log: AuditLog = {
            id: `log-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            documentId,
            userId,
            userName,
            action,
            details,
            ipAddress,
            userAgent,
            timestamp: new Date().toISOString()
        }

        const docLogs = this.logs.get(documentId) || []
        docLogs.unshift(log)

        // 限制日志数量
        if (docLogs.length > this.maxLogsPerDocument) {
            docLogs.pop()
        }

        this.logs.set(documentId, docLogs)
        return log
    }

    /**
     * 获取日志
     */
    getLogs(documentId: string, limit: number = 100): AuditLog[] {
        const docLogs = this.logs.get(documentId) || []
        return docLogs.slice(0, limit)
    }

    /**
     * 搜索日志
     */
    searchLogs(
        documentId: string,
        action?: AuditLogType,
        userId?: string,
        startDate?: string,
        endDate?: string
    ): AuditLog[] {
        const docLogs = this.logs.get(documentId) || []
        
        return docLogs.filter(log => {
            if (action && log.action !== action) return false
            if (userId && log.userId !== userId) return false
            if (startDate && log.timestamp < startDate) return false
            if (endDate && log.timestamp > endDate) return false
            return true
        })
    }

    /**
     * 清除日志
     */
    clearLogs(documentId: string): void {
        this.logs.set(documentId, [])
    }
}

/**
 * AI 写作助手（模拟）
 */
export class AIWritingAssistant {
    /**
     * 检查语法和风格
     */
    check(content: string): AIWritingSuggestion[] {
        const suggestions: AIWritingSuggestion[] = []
        const lines = content.split('\n')

        lines.forEach((line, lineIndex) => {
            // 检查长句子
            if (line.length > 100) {
                suggestions.push({
                    id: `suggestion-${Date.now()}-${lineIndex}`,
                    type: 'style',
                    severity: 'warning',
                    message: '句子过长，建议拆分',
                    suggestion: '考虑将长句拆分为多个短句',
                    position: {
                        start: 0,
                        end: line.length,
                        line: lineIndex + 1,
                        column: 1
                    },
                    originalText: line
                })
            }

            // 检查常见拼写错误（简单示例）
            const commonErrors: Record<string, string> = {
                'teh': 'the',
                'adn': 'and',
                'taht': 'that'
            }

            Object.entries(commonErrors).forEach(([error, correction]) => {
                const regex = new RegExp(`\\b${error}\\b`, 'gi')
                const match = regex.exec(line)
                if (match) {
                    suggestions.push({
                        id: `suggestion-${Date.now()}-${lineIndex}-${error}`,
                        type: 'spelling',
                        severity: 'error',
                        message: `可能的拼写错误`,
                        suggestion: correction,
                        position: {
                            start: match.index,
                            end: match.index + error.length,
                            line: lineIndex + 1,
                            column: match.index + 1
                        },
                        originalText: error
                    })
                }
            })
        })

        return suggestions
    }

    /**
     * 改进建议
     */
    suggestImprovements(content: string): string[] {
        const suggestions: string[] = []

        // 检查是否有标题
        if (!content.includes('#')) {
            suggestions.push('建议添加文档标题')
        }

        // 检查段落数量
        const paragraphs = content.split(/\n\n+/).filter(p => p.trim())
        if (paragraphs.length < 3) {
            suggestions.push('文档内容较少，建议补充更多细节')
        }

        // 检查代码块
        if (content.includes('```') && !content.includes('```typescript') && 
            !content.includes('```javascript') && !content.includes('```python')) {
            suggestions.push('建议为代码块指定编程语言')
        }

        return suggestions
    }

    /**
     * 生成摘要
     */
    generateSummary(content: string, maxLength: number = 200): string {
        const lines = content.split('\n')
        const headings = lines.filter(line => line.startsWith('#'))
        
        if (headings.length > 0) {
            return headings.slice(0, 5).join('\n').substring(0, maxLength)
        }

        return content.substring(0, maxLength)
    }
}

/**
 * 标签管理器
 */
export class TagManager {
    private tags: Map<string, DocumentTag> = new Map()

    constructor() {
        // 初始化常用标签
        const commonTags: DocumentTag[] = [
            { id: 'tag-1', name: '技术文档', color: '#409EFF', count: 0 },
            { id: 'tag-2', name: 'API', color: '#67C23A', count: 0 },
            { id: 'tag-3', name: '教程', color: '#E6A23C', count: 0 },
            { id: 'tag-4', name: '指南', color: '#F56C6C', count: 0 },
            { id: 'tag-5', name: '最佳实践', color: '#909399', count: 0 },
        ]
        
        commonTags.forEach(tag => this.tags.set(tag.id, tag))
    }

    /**
     * 添加标签
     */
    addTag(tag: DocumentTag): void {
        this.tags.set(tag.id, tag)
    }

    /**
     * 获取所有标签
     */
    getAllTags(): DocumentTag[] {
        return Array.from(this.tags.values())
    }

    /**
     * 推荐标签
     */
    recommendTags(content: string): string[] {
        const recommendations: string[] = []
        
        // 基于内容关键词推荐
        if (content.includes('API') || content.includes('interface')) {
            recommendations.push('API')
        }
        
        if (content.includes('function') || content.includes('class')) {
            recommendations.push('技术文档')
        }
        
        if (content.includes('步骤') || content.includes('教程')) {
            recommendations.push('教程')
        }

        return recommendations
    }

    /**
     * 增加标签计数
     */
    incrementTagCount(tagId: string): void {
        const tag = this.tags.get(tagId)
        if (tag) {
            tag.count++
            this.tags.set(tagId, tag)
        }
    }
}

/**
 * 文档关联管理器
 */
export class RelationManager {
    private relations: Map<string, DocumentRelation[]> = new Map()

    /**
     * 添加关联
     */
    addRelation(relation: DocumentRelation): void {
        const sourceRelations = this.relations.get(relation.sourceDocumentId) || []
        sourceRelations.push(relation)
        this.relations.set(relation.sourceDocumentId, sourceRelations)
    }

    /**
     * 获取关联文档
     */
    getRelatedDocuments(documentId: string): DocumentRelation[] {
        return this.relations.get(documentId) || []
    }

    /**
     * 移除关联
     */
    removeRelation(sourceId: string, targetId: string): void {
        const sourceRelations = this.relations.get(sourceId) || []
        const filtered = sourceRelations.filter(r => r.targetDocumentId !== targetId)
        this.relations.set(sourceId, filtered)
    }

    /**
     * 构建关联图谱
     */
    buildGraph(documentId: string): { nodes: any[], links: any[] } {
        const nodes = new Map<string, any>()
        const links: any[] = []

        const relations = this.getRelatedDocuments(documentId)
        
        // 添加中心节点
        nodes.set(documentId, { id: documentId, label: '当前文档' })

        relations.forEach(rel => {
            if (!nodes.has(rel.targetDocumentId)) {
                nodes.set(rel.targetDocumentId, {
                    id: rel.targetDocumentId,
                    label: `文档${rel.targetDocumentId}`
                })
            }

            links.push({
                source: documentId,
                target: rel.targetDocumentId,
                type: rel.type,
                strength: rel.strength
            })
        })

        return {
            nodes: Array.from(nodes.values()),
            links
        }
    }
}

/**
 * 实时协作管理器（简化版）
 */
export class CollaborationManager {
    private collaborators: Map<string, Collaborator[]> = new Map()
    private operations: Map<string, OTOperation[]> = new Map()

    /**
     * 添加协作者
     */
    addCollaborator(documentId: string, collaborator: Collaborator): void {
        const docCollaborators = this.collaborators.get(documentId) || []
        const existing = docCollaborators.find(c => c.userId === collaborator.userId)
        
        if (existing) {
            existing.isOnline = true
            existing.lastActiveAt = new Date().toISOString()
        } else {
            docCollaborators.push(collaborator)
        }
        
        this.collaborators.set(documentId, docCollaborators)
    }

    /**
     * 移除协作者
     */
    removeCollaborator(documentId: string, userId: string): void {
        const docCollaborators = this.collaborators.get(documentId) || []
        const filtered = docCollaborators.filter(c => c.userId !== userId)
        this.collaborators.set(documentId, filtered)
    }

    /**
     * 获取在线协作者
     */
    getOnlineCollaborators(documentId: string): Collaborator[] {
        const docCollaborators = this.collaborators.get(documentId) || []
        return docCollaborators.filter(c => c.isOnline)
    }

    /**
     * 更新光标位置
     */
    updateCursorPosition(
        documentId: string,
        userId: string,
        position: { line: number; column: number }
    ): void {
        const docCollaborators = this.collaborators.get(documentId) || []
        const collaborator = docCollaborators.find(c => c.userId === userId)
        
        if (collaborator) {
            collaborator.cursorPosition = position
            collaborator.lastActiveAt = new Date().toISOString()
        }
    }

    /**
     * 应用操作
     */
    applyOperation(operation: OTOperation): void {
        const docOperations = this.operations.get(operation.documentId) || []
        docOperations.push(operation)
        this.operations.set(operation.documentId, docOperations)
    }

    /**
     * 获取操作历史
     */
    getOperations(documentId: string): OTOperation[] {
        return this.operations.get(documentId) || []
    }
}

/**
 * 格式化日期
 */
export function formatDate(date: string | Date): string {
    const d = new Date(date)
    return d.toLocaleString('zh-CN')
}

/**
 * 生成唯一 ID
 */
export function generateId(prefix: string = 'id'): string {
    return `${prefix}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

/**
 * 防抖函数
 */
export function debounce<T extends (...args: any[]) => any>(
    func: T,
    wait: number
): (...args: Parameters<T>) => void {
    let timeout: NodeJS.Timeout | null = null
    return (...args: Parameters<T>) => {
        if (timeout) clearTimeout(timeout)
        timeout = setTimeout(() => func(...args), wait)
    }
}

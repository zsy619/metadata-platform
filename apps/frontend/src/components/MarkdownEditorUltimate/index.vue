<template>
    <div class="markdown-editor-ultimate" :class="{ 'fullscreen': isFullscreen }">
        <!-- 顶部工具栏 -->
        <div class="editor-toolbar">
            <div class="toolbar-left">
                <!-- 文件操作 -->
                <div class="toolbar-group">
                    <el-tooltip content="保存 (Ctrl+S)" placement="bottom">
                        <el-button size="small" @click="handleSave" :icon="Document" :loading="saving" />
                    </el-tooltip>
                    <el-dropdown trigger="click" size="small">
                        <el-button size="small" :icon="MoreFilled" />
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item @click="showWorkflow = true">
                                    <el-icon><Connection /></el-icon>
                                    工作流
                                </el-dropdown-item>
                                <el-dropdown-item @click="showPermissions = true">
                                    <el-icon><User /></el-icon>
                                    权限管理
                                </el-dropdown-item>
                                <el-dropdown-item @click="showAuditLogs = true">
                                    <el-icon><Document /></el-icon>
                                    审计日志
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- AI 助手 -->
                <div class="toolbar-group">
                    <el-tooltip content="AI 检查" placement="bottom">
                        <el-badge :value="aiSuggestions.length" :hidden="aiSuggestions.length === 0">
                            <el-button 
                                size="small" 
                                @click="runAICheck" 
                                :icon="Cpu"
                                :type="aiSuggestions.length > 0 ? 'warning' : ''"
                            />
                        </el-badge>
                    </el-tooltip>
                </div>
            </div>
            
            <div class="toolbar-right">
                <!-- 协作者 -->
                <el-tooltip content="在线协作者" placement="bottom">
                    <div class="collaborators">
                        <el-avatar
                            v-for="collab in onlineCollaborators"
                            :key="collab.userId"
                            :size="24"
                            :style="{ backgroundColor: collab.color }"
                            :title="collab.userName"
                        >
                            {{ collab.userName.charAt(0) }}
                        </el-avatar>
                        <span v-if="onlineCollaborators.length > 0" class="collaborator-count">
                            {{ onlineCollaborators.length }}
                        </span>
                    </div>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 工作流状态 -->
                <el-tag v-if="workflowInstance" :type="getWorkflowStatusType(workflowInstance.status)" size="small">
                    {{ getWorkflowStatusLabel(workflowInstance.status) }}
                </el-tag>
                <el-divider direction="vertical" v-if="workflowInstance" />
                
                <!-- 统计 -->
                <div class="stats-trigger" @click="showStats = true">
                    <el-icon><DataAnalysis /></el-icon>
                    <span>{{ statistics.wordCount }} 字</span>
                </div>
            </div>
        </div>
        
        <!-- 主体内容 -->
        <div class="editor-body">
            <!-- 编辑区 -->
            <div class="editor-pane">
                <textarea
                    ref="editorRef"
                    v-model="modelValue"
                    class="editor-textarea"
                    placeholder="开始编写..."
                    @input="handleInput"
                />
                
                <!-- AI 建议面板 -->
                <div v-if="showAISuggestions && aiSuggestions.length > 0" class="ai-panel">
                    <div class="ai-header">
                        <span>AI 建议 ({{ aiSuggestions.length }})</span>
                        <el-button link size="small" @click="showAISuggestions = false">
                            <el-icon><Close /></el-icon>
                        </el-button>
                    </div>
                    <el-scrollbar height="300px">
                        <div class="ai-suggestions">
                            <div
                                v-for="suggestion in aiSuggestions"
                                :key="suggestion.id"
                                class="ai-suggestion"
                                :class="suggestion.severity"
                            >
                                <div class="ai-suggestion-header">
                                    <el-tag size="small" :type="getSeverityType(suggestion.severity)">
                                        {{ getSuggestionTypeLabel(suggestion.type) }}
                                    </el-tag>
                                    <span class="ai-message">{{ suggestion.message }}</span>
                                </div>
                                <div class="ai-suggestion-content">
                                    <span class="original">{{ suggestion.originalText }}</span>
                                    <el-icon class="arrow"><DArrowRight /></el-icon>
                                    <span class="suggested" v-if="suggestion.suggestion">
                                        {{ suggestion.suggestion }}
                                    </span>
                                </div>
                                <div class="ai-suggestion-actions">
                                    <el-button 
                                        size="small" 
                                        text 
                                        type="primary"
                                        @click="applySuggestion(suggestion)"
                                    >
                                        应用
                                    </el-button>
                                    <el-button size="small" text @click="ignoreSuggestion(suggestion.id)">
                                        忽略
                                    </el-button>
                                </div>
                            </div>
                        </div>
                    </el-scrollbar>
                </div>
            </div>
            
            <!-- 预览区 -->
            <div class="preview-pane">
                <div class="preview-content" v-html="renderedContent"></div>
            </div>
        </div>
        
        <!-- 底部状态栏 -->
        <div class="editor-statusbar">
            <div class="status-left">
                <span class="status-item">
                    <el-icon><Position /></el-icon>
                    行 {{ cursorPosition.line }}, 列 {{ cursorPosition.column }}
                </span>
                <el-divider direction="vertical" />
                <span class="status-item" v-if="workflowInstance">
                    <el-icon><Connection /></el-icon>
                    工作流：{{ getWorkflowStatusLabel(workflowInstance.status) }}
                </span>
            </div>
            <div class="status-right">
                <el-button text size="small" @click="runAICheck">
                    <el-icon><Cpu /></el-icon>
                    AI 检查
                </el-button>
            </div>
        </div>
        
        <!-- 工作流对话框 -->
        <el-dialog v-model="showWorkflow" title="工作流" width="800px">
            <el-timeline v-if="workflowInstance">
                <el-timeline-item
                    v-for="node in workflowInstance.nodes"
                    :key="node.id"
                    :timestamp="node.reviewedAt || '待处理'"
                    :type="node.status === WorkflowStatus.APPROVED ? 'success' : node.status === WorkflowStatus.REJECTED ? 'danger' : ''"
                >
                    <el-card>
                        <div class="workflow-node">
                            <div class="node-header">
                                <strong>{{ node.name }}</strong>
                                <el-tag size="small" :type="getNodeStatusType(node.status)">
                                    {{ getNodeStatusLabel(node.status) }}
                                </el-tag>
                            </div>
                            <div class="node-info" v-if="node.reviewedBy">
                                审核人：{{ node.reviewedBy }}
                                <span v-if="node.comment"> - {{ node.comment }}</span>
                            </div>
                            <div class="node-actions" v-if="canApprove && node.status === WorkflowStatus.IN_REVIEW">
                                <el-button size="small" type="success" @click="approveNode(node.id)">通过</el-button>
                                <el-button size="small" type="danger" @click="rejectNode(node.id)">拒绝</el-button>
                            </div>
                        </div>
                    </el-card>
                </el-timeline-item>
            </el-timeline>
            <div v-else>
                <el-empty description="暂无工作流" />
            </div>
        </el-dialog>
        
        <!-- 权限管理对话框 -->
        <el-dialog v-model="showPermissions" title="权限管理" width="700px">
            <el-table :data="permissions" style="width: 100%">
                <el-table-column prop="userName" label="用户" width="150" />
                <el-table-column prop="role" label="角色" width="120">
                    <template #default="{ row }">
                        <el-tag size="small">{{ getRoleLabel(row.role) }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="权限">
                    <template #default="{ row }">
                        <el-checkbox-group v-model="row.permissions" disabled>
                            <el-checkbox label="读取" :disabled="!row.canRead" />
                            <el-checkbox label="编辑" :disabled="!row.canWrite" />
                            <el-checkbox label="删除" :disabled="!row.canDelete" />
                            <el-checkbox label="分享" :disabled="!row.canShare" />
                        </el-checkbox-group>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="150">
                    <template #default="{ row }">
                        <el-button size="small" text type="danger" @click="removePermission(row.userId)">
                            移除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-dialog>
        
        <!-- 审计日志对话框 -->
        <el-dialog v-model="showAuditLogs" title="审计日志" width="900px">
            <el-table :data="auditLogs" style="width: 100%">
                <el-table-column prop="timestamp" label="时间" width="180" />
                <el-table-column prop="userName" label="用户" width="120" />
                <el-table-column prop="action" label="操作" width="100">
                    <template #default="{ row }">
                        <el-tag size="small" :type="getActionType(row.action)">
                            {{ row.action }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="details" label="详情" />
            </el-table>
        </el-dialog>
        
        <!-- 统计面板 -->
        <el-dialog v-model="showStats" title="文档统计" width="600px">
            <el-descriptions :column="2" border>
                <el-descriptions-item label="总字数">{{ statistics.wordCount }}</el-descriptions-item>
                <el-descriptions-item label="字符数">{{ statistics.characterCount }}</el-descriptions-item>
                <el-descriptions-item label="段落数">{{ statistics.paragraphCount }}</el-descriptions-item>
                <el-descriptions-item label="行数">{{ statistics.lineCount }}</el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import {
    AIWritingAssistant,
    AIWritingSuggestion,
    AuditLogManager,
    AuditLogType,
    CollaborationManager,
    DocumentPermission,
    PermissionManager,
    UserRole,
    WorkflowInstance,
    WorkflowManager,
    WorkflowStatus
} from '@/utils/markdownUltimate'
import {
    Close,
    Connection,
    Cpu,
    DArrowRight,
    DataAnalysis,
    Document, MoreFilled,
    Position,
    User
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import { computed, onMounted, reactive, ref } from 'vue'

// Props
interface Props {
    modelValue: string
    title?: string
    documentId?: string
}

const props = withDefaults(defineProps<Props>(), {
    title: '',
    documentId: ''
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: string]
    'save': [content: string]
}>()

// 引用
const editorRef = ref<HTMLTextAreaElement | null>(null)

// 状态
const isFullscreen = ref(false)
const saving = ref(false)
const cursorPosition = reactive({ line: 1, column: 1 })
const showWorkflow = ref(false)
const showPermissions = ref(false)
const showAuditLogs = ref(false)
const showStats = ref(false)
const showAISuggestions = ref(false)
const aiSuggestions = ref<AIWritingSuggestion[]>([])

// 管理器实例
const permissionManager = new PermissionManager()
const workflowManager = new WorkflowManager()
const auditLogManager = new AuditLogManager()
const aiAssistant = new AIWritingAssistant()
const collaborationManager = new CollaborationManager()

// 计算属性
const modelValue = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const statistics = computed(() => ({
    wordCount: (modelValue.value.match(/\b\w+\b/g) || []).length + (modelValue.value.match(/[\u4e00-\u9fa5]/g) || []).length,
    characterCount: modelValue.value.length,
    paragraphCount: modelValue.value.split(/\n\n+/).filter(p => p.trim()).length,
    lineCount: modelValue.value.split('\n').length
}))

const renderedContent = computed(() => {
    if (!modelValue.value) return '<p style="color: #909399;">暂无内容</p>'
    try {
        return marked(modelValue.value)
    } catch {
        return '<p style="color: #f56c6c;">解析失败</p>'
    }
})

const workflowInstance = ref<WorkflowInstance | null>(null)
const permissions = ref<DocumentPermission[]>([])
const auditLogs = ref<any[]>([])
const onlineCollaborators = ref<any[]>([])

// 初始化
onMounted(() => {
    // 初始化示例数据
    initDemoData()
    
    // 记录访问日志
    auditLogManager.log(
        props.documentId || 'demo',
        'current-user',
        '当前用户',
        AuditLogType.VIEW,
        '查看文档'
    )
})

// 初始化示例数据
function initDemoData() {
    // 创建工作流
    const workflow = workflowManager.createWorkflow(
        props.documentId || 'demo',
        [
            { id: 'node-1', name: '技术审核', type: 'review', reviewers: ['tech-lead'], status: WorkflowStatus.PENDING_REVIEW },
            { id: 'node-2', name: '产品审核', type: 'approve', reviewers: ['product-manager'], status: WorkflowStatus.DRAFT }
        ],
        'current-user'
    )
    workflowInstance.value = workflow
    
    // 添加权限
    permissionManager.setPermission({
        documentId: props.documentId || 'demo',
        userId: 'current-user',
        role: UserRole.EDITOR,
        canRead: true,
        canWrite: true,
        canDelete: false,
        canShare: true,
        canApprove: true,
        grantedAt: new Date().toISOString(),
        grantedBy: 'admin'
    })
    
    // 添加协作者（模拟）
    collaborationManager.addCollaborator(props.documentId || 'demo', {
        userId: 'user-2',
        userName: '张三',
        color: '#409EFF',
        isOnline: true,
        lastActiveAt: new Date().toISOString()
    })
    
    updateCollaborators()
    updatePermissions()
    updateAuditLogs()
}

// 方法
const handleSave = () => {
    saving.value = true
    auditLogManager.log(
        props.documentId || 'demo',
        'current-user',
        '当前用户',
        AuditLogType.UPDATE,
        '保存文档'
    )
    emit('save', modelValue.value)
    setTimeout(() => {
        saving.value = false
        ElMessage.success('保存成功')
    }, 500)
}

const handleInput = () => {
    updateCursorPosition()
}

const updateCursorPosition = () => {
    if (!editorRef.value) return
    const pos = editorRef.value.selectionStart
    const text = modelValue.value.substring(0, pos)
    const lines = text.split('\n')
    cursorPosition.line = lines.length
    cursorPosition.column = lines[lines.length - 1].length + 1
}

const runAICheck = () => {
    const suggestions = aiAssistant.check(modelValue.value)
    aiSuggestions.value = suggestions
    showAISuggestions.value = true
    
    if (suggestions.length === 0) {
        ElMessage.success('文档质量良好，未发现问题')
    } else {
        ElMessage.warning(`发现 ${suggestions.length} 个建议`)
    }
}

const applySuggestion = (suggestion: AIWritingSuggestion) => {
    if (!suggestion.suggestion) return
    
    const text = modelValue.value
    const before = text.substring(0, suggestion.position.start)
    const after = text.substring(suggestion.position.end)
    modelValue.value = before + suggestion.suggestion + after
    
    ignoreSuggestion(suggestion.id)
    ElMessage.success('已应用建议')
}

const ignoreSuggestion = (id: string) => {
    aiSuggestions.value = aiSuggestions.value.filter(s => s.id !== id)
}

const approveNode = (nodeId: string) => {
    if (!workflowInstance.value) return
    
    const updated = workflowManager.approve(
        workflowInstance.value.id,
        nodeId,
        'current-user',
        '审核通过'
    )
    
    if (updated) {
        workflowInstance.value = updated
        auditLogManager.log(
            props.documentId || 'demo',
            'current-user',
            '当前用户',
            AuditLogType.APPROVE,
            `通过节点：${nodeId}`
        )
        ElMessage.success('审核通过')
    }
}

const rejectNode = (nodeId: string) => {
    if (!workflowInstance.value) return
    
    const updated = workflowManager.reject(
        workflowInstance.value.id,
        nodeId,
        'current-user',
        '需要修改'
    )
    
    if (updated) {
        workflowInstance.value = updated
        auditLogManager.log(
            props.documentId || 'demo',
            'current-user',
            '当前用户',
            AuditLogType.REJECT,
            `拒绝节点：${nodeId}`
        )
        ElMessage.warning('已拒绝')
    }
}

const updateCollaborators = () => {
    onlineCollaborators.value = collaborationManager.getOnlineCollaborators(props.documentId || 'demo')
}

const updatePermissions = () => {
    const perms = permissionManager.getDocumentPermissions(props.documentId || 'demo')
    permissions.value = perms.map(p => ({
        ...p,
        userName: p.userId,
        permissions: [
            p.canRead ? '读取' : '',
            p.canWrite ? '编辑' : '',
            p.canDelete ? '删除' : '',
            p.canShare ? '分享' : ''
        ].filter(Boolean)
    }))
}

const updateAuditLogs = () => {
    auditLogs.value = auditLogManager.getLogs(props.documentId || 'demo', 50)
}

const removePermission = (userId: string) => {
    permissionManager.removePermission(props.documentId || 'demo', userId)
    updatePermissions()
    ElMessage.success('已移除权限')
}

// 工具函数
const getWorkflowStatusType = (status: WorkflowStatus) => {
    const types: Record<WorkflowStatus, any> = {
        [WorkflowStatus.DRAFT]: 'info',
        [WorkflowStatus.PENDING_REVIEW]: 'warning',
        [WorkflowStatus.IN_REVIEW]: 'warning',
        [WorkflowStatus.APPROVED]: 'success',
        [WorkflowStatus.REJECTED]: 'danger',
        [WorkflowStatus.PUBLISHED]: 'success',
        [WorkflowStatus.ARCHIVED]: 'info'
    }
    return types[status] || 'info'
}

const getWorkflowStatusLabel = (status: WorkflowStatus) => {
    const labels: Record<WorkflowStatus, string> = {
        [WorkflowStatus.DRAFT]: '草稿',
        [WorkflowStatus.PENDING_REVIEW]: '待审核',
        [WorkflowStatus.IN_REVIEW]: '审核中',
        [WorkflowStatus.APPROVED]: '已通过',
        [WorkflowStatus.REJECTED]: '已拒绝',
        [WorkflowStatus.PUBLISHED]: '已发布',
        [WorkflowStatus.ARCHIVED]: '已归档'
    }
    return labels[status]
}

const getNodeStatusLabel = (status: WorkflowStatus) => getWorkflowStatusLabel(status)
const getNodeStatusType = (status: WorkflowStatus) => getWorkflowStatusType(status)

const getRoleLabel = (role: UserRole) => {
    const labels: Record<UserRole, string> = {
        [UserRole.ADMIN]: '管理员',
        [UserRole.EDITOR]: '编辑者',
        [UserRole.REVIEWER]: '审核者',
        [UserRole.VIEWER]: '查看者'
    }
    return labels[role]
}

const getActionType = (action: AuditLogType) => {
    const types: Record<AuditLogType, any> = {
        [AuditLogType.CREATE]: 'success',
        [AuditLogType.UPDATE]: 'warning',
        [AuditLogType.DELETE]: 'danger',
        [AuditLogType.VIEW]: 'info',
        [AuditLogType.APPROVE]: 'success',
        [AuditLogType.REJECT]: 'danger',
        [AuditLogType.PUBLISH]: 'success',
        [AuditLogType.SHARE]: 'warning',
        [AuditLogType.DOWNLOAD]: 'info',
        [AuditLogType.PRINT]: 'info'
    }
    return types[action] || 'info'
}

const getSeverityType = (severity: string) => {
    const types: Record<string, any> = {
        error: 'danger',
        warning: 'warning',
        info: 'info'
    }
    return types[severity] || 'info'
}

const getSuggestionTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
        grammar: '语法',
        style: '风格',
        clarity: '清晰度',
        tone: '语气',
        spelling: '拼写'
    }
    return labels[type] || type
}

const canApprove = computed(() => {
    return permissionManager.checkPermission(
        props.documentId || 'demo',
        'current-user',
        'approve'
    )
})

// 暴露方法
defineExpose({
    getContent: () => modelValue.value,
    setContent: (content: string) => {
        modelValue.value = content
    }
})
</script>

<style scoped lang="scss">
.markdown-editor-ultimate {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 200px);
    border: 1px solid #dcdfe6;
    background: #fff;
    
    &.fullscreen {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: 9999;
        height: 100vh;
    }
    
    .editor-toolbar {
        display: flex;
        justify-content: space-between;
        padding: 8px 12px;
        background: #f5f7fa;
        border-bottom: 1px solid #dcdfe6;
        
        .toolbar-group {
            display: flex;
            align-items: center;
            gap: 4px;
        }
        
        .collaborators {
            display: flex;
            align-items: center;
            gap: -8px;
            
            .el-avatar {
                border: 2px solid #fff;
                margin-left: -8px;
                
                &:first-child {
                    margin-left: 0;
                }
            }
            
            .collaborator-count {
                font-size: 12px;
                color: #909399;
                margin-left: 4px;
            }
        }
        
        .stats-trigger {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 12px;
            color: #909399;
            cursor: pointer;
            
            &:hover {
                color: #409EFF;
            }
        }
    }
    
    .editor-body {
        flex: 1;
        display: flex;
        overflow: hidden;
        
        .editor-pane {
            flex: 1;
            position: relative;
            display: flex;
            flex-direction: column;
            
            .editor-textarea {
                flex: 1;
                width: 100%;
                padding: 16px;
                border: none;
                resize: none;
                font-family: 'Consolas', monospace;
                font-size: 14px;
                line-height: 1.6;
            }
            
            .ai-panel {
                position: absolute;
                bottom: 0;
                left: 0;
                right: 0;
                background: #fff;
                border-top: 1px solid #dcdfe6;
                box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
                
                .ai-header {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 8px 12px;
                    background: #f5f7fa;
                    border-bottom: 1px solid #e4e7ed;
                }
                
                .ai-suggestions {
                    padding: 12px;
                    
                    .ai-suggestion {
                        padding: 12px;
                        margin-bottom: 12px;
                        border-radius: 4px;
                        border-left: 3px solid;
                        
                        &.error {
                            border-color: #f56c6c;
                            background: rgba(245, 108, 108, 0.05);
                        }
                        
                        &.warning {
                            border-color: #e6a23c;
                            background: rgba(230, 162, 60, 0.05);
                        }
                        
                        &.info {
                            border-color: #409eff;
                            background: rgba(64, 158, 255, 0.05);
                        }
                        
                        .ai-suggestion-header {
                            display: flex;
                            align-items: center;
                            gap: 8px;
                            margin-bottom: 8px;
                            
                            .ai-message {
                                font-size: 13px;
                                color: #606266;
                            }
                        }
                        
                        .ai-suggestion-content {
                            display: flex;
                            align-items: center;
                            gap: 8px;
                            margin-bottom: 8px;
                            font-size: 13px;
                            
                            .original {
                                color: #f56c6c;
                                text-decoration: line-through;
                            }
                            
                            .suggested {
                                color: #67c23a;
                                font-weight: 500;
                            }
                            
                            .arrow {
                                color: #909399;
                            }
                        }
                        
                        .ai-suggestion-actions {
                            display: flex;
                            gap: 8px;
                        }
                    }
                }
            }
        }
        
        .preview-pane {
            flex: 1;
            border-left: 1px solid #dcdfe6;
            
            .preview-content {
                height: 100%;
                padding: 16px;
                overflow-y: auto;
            }
        }
    }
    
    .editor-statusbar {
        display: flex;
        justify-content: space-between;
        padding: 4px 12px;
        background: #f5f7fa;
        border-top: 1px solid #dcdfe6;
        font-size: 12px;
        color: #909399;
        
        .status-left, .status-right {
            display: flex;
            align-items: center;
            gap: 12px;
        }
    }
}

.workflow-node {
    .node-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
    }
    
    .node-info {
        font-size: 13px;
        color: #606266;
        margin-bottom: 8px;
    }
}
</style>

<template>
    <div class="markdown-editor-pro" :class="{ 
        'fullscreen': isFullscreen,
        'dark-mode': darkMode,
        'read-mode': readMode
    }">
        <!-- 顶部工具栏 -->
        <div class="editor-toolbar">
            <div class="toolbar-left">
                <!-- 文件操作 -->
                <div class="toolbar-group">
                    <el-tooltip content="新建 (Ctrl+N)" placement="bottom">
                        <el-button size="small" @click="handleNew" :icon="DocumentAdd" />
                    </el-tooltip>
                    <el-tooltip content="保存 (Ctrl+S)" placement="bottom">
                        <el-button size="small" @click="handleSave" :icon="Document" :loading="saving" />
                    </el-tooltip>
                    <el-dropdown trigger="click" size="small">
                        <el-button size="small" :icon="MoreFilled" />
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item @click="showVersionHistory = true">
                                    <el-icon><Clock /></el-icon>
                                    版本历史
                                </el-dropdown-item>
                                <el-dropdown-item @click="showDocumentCompare = true">
                                    <el-icon><DArrowRight /></el-icon>
                                    文档对比
                                </el-dropdown-item>
                                <el-dropdown-item @click="showAnnotations = true">
                                    <el-icon><ChatDotRound /></el-icon>
                                    批注管理
                                </el-dropdown-item>
                                <el-dropdown-item divided @click="toggleProtection">
                                    <el-icon><Lock /></el-icon>
                                    {{ isProtected ? '取消保护' : '文档保护' }}
                                </el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 编辑操作 -->
                <div class="toolbar-group">
                    <el-tooltip content="撤销 (Ctrl+Z)" placement="bottom">
                        <el-button size="small" @click="handleUndo" :icon="RefreshLeft" :disabled="!canUndo" />
                    </el-tooltip>
                    <el-tooltip content="重做 (Ctrl+Y)" placement="bottom">
                        <el-button size="small" @click="handleRedo" :icon="RefreshRight" :disabled="!canRedo" />
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 视图模式 -->
                <div class="toolbar-group">
                    <el-tooltip content="阅读模式" placement="bottom">
                        <el-button 
                            size="small" 
                            @click="toggleReadMode" 
                            :icon="View"
                            :type="readMode ? 'primary' : ''"
                        />
                    </el-tooltip>
                    <el-tooltip content="全屏" placement="bottom">
                        <el-button size="small" @click="toggleFullscreen" :icon="FullScreen" />
                    </el-tooltip>
                </div>
            </div>
            
            <div class="toolbar-right">
                <!-- 统计信息 -->
                <el-popover
                    placement="bottom"
                    :width="300"
                    trigger="click"
                >
                    <template #reference>
                        <div class="stats-trigger">
                            <el-icon><DataAnalysis /></el-icon>
                            <span>{{ statistics.wordCount }} 字</span>
                        </div>
                    </template>
                    <div class="statistics-panel">
                        <h4>文档统计</h4>
                        <el-descriptions :column="2" size="small" border>
                            <el-descriptions-item label="总字数">{{ statistics.wordCount }}</el-descriptions-item>
                            <el-descriptions-item label="字符数">{{ statistics.characterCount }}</el-descriptions-item>
                            <el-descriptions-item label="段落数">{{ statistics.paragraphCount }}</el-descriptions-item>
                            <el-descriptions-item label="行数">{{ statistics.lineCount }}</el-descriptions-item>
                            <el-descriptions-item label="阅读时间">{{ statistics.readingTime }} 分钟</el-descriptions-item>
                            <el-descriptions-item label="写作时间">{{ statistics.writingTime }} 分钟</el-descriptions-item>
                        </el-descriptions>
                        
                        <h4 style="margin-top: 12px;">元素统计</h4>
                        <el-descriptions :column="2" size="small" border>
                            <el-descriptions-item label="标题数">{{ totalHeadings }}</el-descriptions-item>
                            <el-descriptions-item label="代码块">{{ statistics.codeBlockCount }}</el-descriptions-item>
                            <el-descriptions-item label="链接数">{{ statistics.linkCount }}</el-descriptions-item>
                            <el-descriptions-item label="图片数">{{ statistics.imageCount }}</el-descriptions-item>
                            <el-descriptions-item label="表格数">{{ statistics.tableCount }}</el-descriptions-item>
                        </el-descriptions>
                    </div>
                </el-popover>
                
                <el-divider direction="vertical" />
                
                <!-- 主题切换 -->
                <el-tooltip :content="darkMode ? '浅色模式' : '深色模式'" placement="bottom">
                    <el-button
                        size="small"
                        @click="toggleDarkMode"
                        :icon="darkMode ? 'Sunny' : 'Moon'"
                    />
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 导出 -->
                <el-dropdown trigger="click" size="small">
                    <el-button size="small" :icon="Download">导出</el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="handleExport('markdown')">Markdown</el-dropdown-item>
                            <el-dropdown-item @click="handleExport('html')">HTML</el-dropdown-item>
                            <el-dropdown-item @click="handleExport('pdf')">PDF</el-dropdown-item>
                            <el-dropdown-item @click="handleExport('word')">Word</el-dropdown-item>
                            <el-dropdown-item divided @click="showExportWithWatermark = true">
                                <el-icon><Lock /></el-icon>
                                带水印导出
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
                <el-divider direction="vertical" />
                
                <!-- 快捷键帮助 -->
                <el-tooltip content="快捷键 (F1)" placement="bottom">
                    <el-button size="small" @click="showShortcuts = true" :icon="QuestionFilled" />
                </el-tooltip>
            </div>
        </div>
        
        <!-- 主体内容 -->
        <div class="editor-body">
            <!-- 左侧大纲 -->
            <el-aside v-if="showOutline" width="260px" class="outline-panel">
                <div class="panel-header">
                    <span>文档大纲</span>
                    <el-button link size="small" @click="showOutline = false">
                        <el-icon><Close /></el-icon>
                    </el-button>
                </div>
                <el-scrollbar height="calc(100vh - 200px)">
                    <div class="outline-tree">
                        <div
                            v-for="item in outline"
                            :key="item.id"
                            class="outline-node"
                            :style="{ paddingLeft: (item.level - 1) * 16 + 8 + 'px' }"
                            @click="scrollToHeading(item.id, item.line)"
                        >
                            <el-icon class="outline-icon"><Bookmark /></el-icon>
                            <span class="outline-text">{{ item.text }}</span>
                            <span class="outline-line">L{{ item.line }}</span>
                        </div>
                        <el-empty v-if="outline.length === 0" description="暂无大纲" :image-size="60" />
                    </div>
                </el-scrollbar>
            </el-aside>
            
            <!-- 编辑区 -->
            <div 
                v-show="!readMode && (viewMode === 'edit' || viewMode === 'split')" 
                class="editor-pane"
                @dragover.prevent="handleDragOver"
                @drop.prevent="handleDrop"
            >
                <textarea
                    ref="editorRef"
                    v-model="modelValue"
                    class="editor-textarea"
                    placeholder="开始编写 Markdown 文档..."
                    @scroll="handleScroll"
                    @keydown.ctrl.s.prevent="handleSave"
                    @keydown.ctrl.b.prevent="insertFormat('bold')"
                    @keydown.ctrl.i.prevent="insertFormat('italic')"
                    @keydown.ctrl.k.prevent="insertLink"
                    @paste="handlePaste"
                />
                <div v-if="isDragging" class="drag-overlay">
                    <el-icon class="drag-icon"><UploadFilled /></el-icon>
                    <div class="drag-text">释放以上传图片</div>
                </div>
            </div>
            
            <!-- 预览区 -->
            <div v-show="viewMode === 'preview' || viewMode === 'split' || readMode" class="preview-pane">
                <div class="preview-header" v-if="!readMode">
                    <span>预览</span>
                    <div class="preview-actions">
                        <el-button size="small" text @click="zoomIn">
                            <el-icon><ZoomIn /></el-icon>
                        </el-button>
                        <el-button size="small" text @click="zoomOut">
                            <el-icon><ZoomOut /></el-icon>
                        </el-button>
                        <el-button size="small" text @click="zoomReset">
                            <el-icon><RefreshRight /></el-icon>
                        </el-button>
                    </div>
                </div>
                <div class="preview-content" v-html="renderedContent" :style="{ zoom: previewZoom }"></div>
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
                <span class="status-item" v-if="lastSaved">
                    <el-icon><Clock /></el-icon>
                    已保存：{{ lastSaved }}
                </span>
                <el-divider direction="vertical" />
                <span class="status-item" v-if="versionInfo">
                    <el-icon><Document /></el-icon>
                    版本：v{{ versionInfo.version }}
                </span>
            </div>
            <div class="status-right">
                <el-button text size="small" @click="showOutline = !showOutline">
                    <el-icon><Menu /></el-icon>
                    大纲
                </el-button>
            </div>
        </div>
        
        <!-- 版本历史对话框 -->
        <el-dialog
            v-model="showVersionHistory"
            title="版本历史"
            width="900px"
        >
            <el-timeline>
                <el-timeline-item
                    v-for="version in versions"
                    :key="version.id"
                    :timestamp="formatDate(version.createdAt)"
                    placement="top"
                    :type="version.changeType === 'create' ? 'success' : version.changeType === 'restore' ? 'warning' : ''"
                >
                    <el-card>
                        <div class="version-card">
                            <div class="version-header">
                                <el-tag size="small">v{{ version.version }}</el-tag>
                                <span class="version-author">{{ version.author }}</span>
                                <span class="version-comment">{{ version.comment }}</span>
                            </div>
                            <div class="version-actions">
                                <el-button size="small" @click="previewVersion(version.version)">
                                    预览
                                </el-button>
                                <el-button size="small" @click="compareVersion(version.version)">
                                    对比
                                </el-button>
                                <el-button size="small" type="primary" @click="restoreVersion(version.version)">
                                    恢复
                                </el-button>
                            </div>
                        </div>
                    </el-card>
                </el-timeline-item>
            </el-timeline>
        </el-dialog>
        
        <!-- 快捷键对话框 -->
        <el-dialog
            v-model="showShortcuts"
            title="快捷键"
            width="700px"
        >
            <el-table :data="shortcutsByCategory" :show-header="false" style="width: 100%">
                <el-table-column prop="category" label="类别" width="120" />
                <el-table-column prop="key" label="快捷键" width="150">
                    <template #default="{ row }">
                        <el-tag size="small">{{ row.key }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="description" label="说明" />
            </el-table>
        </el-dialog>
        
        <!-- 输入框引用 -->
        <input
            ref="imageUploadRef"
            type="file"
            accept="image/*"
            multiple
            style="display: none"
            @change="handleImageUpload"
        />
    </div>
</template>

<script setup lang="ts">
import {
    exportToHtml,
    exportToMarkdown,
    exportToPdf, exportToWord
} from '@/utils/documentExport'
import {
    AnnotationManager, VersionManager,
    calculateStatistics, extractOutline,
    shortcuts as shortcutList
} from '@/utils/markdownPro'
import {
    ChatDotRound,
    Clock,
    Close,
    DArrowRight,
    DataAnalysis,
    Document,
    DocumentAdd,
    Download,
    FullScreen,
    Lock,
    Menu,
    MoreFilled,
    Position,
    QuestionFilled, RefreshLeft,
    RefreshRight,
    UploadFilled,
    View,
    ZoomIn,
    ZoomOut
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

// Props
interface Props {
    modelValue: string
    title?: string
}

const props = withDefaults(defineProps<Props>(), {
    title: ''
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: string]
    'save': [content: string]
}>()

// 引用
const editorRef = ref<HTMLTextAreaElement | null>(null)
const imageUploadRef = ref<HTMLInputElement | null>(null)

// 状态
const viewMode = ref<'edit' | 'preview' | 'split'>('split')
const isFullscreen = ref(false)
const darkMode = ref(false)
const readMode = ref(false)
const saving = ref(false)
const lastSaved = ref('')
const isDragging = ref(false)
const showOutline = ref(true)
const previewZoom = ref(1)

// 对话框
const showVersionHistory = ref(false)
const showDocumentCompare = ref(false)
const showAnnotations = ref(false)
const showExportWithWatermark = ref(false)
const showShortcuts = ref(false)

// 保护
const isProtected = ref(false)

// 版本管理
const versionManager = new VersionManager(50)
const annotationManager = new AnnotationManager()

// 计算属性
const modelValue = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const statistics = computed(() => calculateStatistics(props.modelValue))
const outline = computed(() => extractOutline(props.modelValue))
const totalHeadings = computed(() => {
    const h = statistics.value.headingCount
    return h.h1 + h.h2 + h.h3 + h.h4 + h.h5 + h.h6
})

const renderedContent = computed(() => {
    if (!props.modelValue) {
        return '<p style="color: #909399; text-align: center; padding: 40px;">暂无内容</p>'
    }
    try {
        return marked(props.modelValue)
    } catch (error) {
        return '<p style="color: #f56c6c;">Markdown 解析失败</p>'
    }
})

const shortcutsByCategory = computed(() => {
    const categories = [...new Set(shortcutList.map(s => s.category))]
    return categories.map(category => ({
        category,
        key: shortcutList.find(s => s.category === category)?.key || '',
        description: shortcutList.filter(s => s.category === category).map(s => s.description).join('、')
    }))
})

const versions = computed(() => versionManager.getVersions())
const versionInfo = computed(() => versions.value[versions.value.length - 1])

// 监听内容变化
watch(() => props.modelValue, (newValue) => {
    // 自动保存草稿
    localStorage.setItem('markdown-pro-draft', newValue)
    
    // 保存版本（防抖）
    clearTimeout(autoSaveTimer)
    autoSaveTimer = setTimeout(() => {
        versionManager.saveVersion(newValue, '当前用户', '自动保存')
        lastSaved.value = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }, 5000)
}, { immediate: true })

let autoSaveTimer: NodeJS.Timeout | null = null

// 生命周期
onMounted(() => {
    const draft = localStorage.getItem('markdown-pro-draft')
    if (draft && !props.modelValue) {
        emit('update:modelValue', draft)
        ElMessage.info('已恢复草稿')
    }
    
    document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer)
    }
    document.removeEventListener('keydown', handleGlobalKeydown)
})

// 方法
const handleSave = () => {
    saving.value = true
    versionManager.saveVersion(props.modelValue, '当前用户', '手动保存')
    emit('save', props.modelValue)
    setTimeout(() => {
        saving.value = false
        lastSaved.value = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
        ElMessage.success('保存成功')
    }, 500)
}

const handleNew = () => {
    emit('update:modelValue', '')
    ElMessage.success('已新建')
}

const toggleFullscreen = () => {
    isFullscreen.value = !isFullscreen.value
}

const toggleDarkMode = () => {
    darkMode.value = !darkMode.value
}

const toggleReadMode = () => {
    readMode.value = !readMode.value
}

const zoomIn = () => {
    previewZoom.value = Math.min(2, previewZoom.value + 0.1)
}

const zoomOut = () => {
    previewZoom.value = Math.max(0.5, previewZoom.value - 0.1)
}

const zoomReset = () => {
    previewZoom.value = 1
}

const handleExport = (type: string) => {
    const filename = props.title || 'document'
    
    switch (type) {
        case 'markdown':
            exportToMarkdown(props.modelValue, filename)
            break
        case 'html':
            exportToHtml(marked(props.modelValue) as string, filename, props.title)
            break
        case 'pdf':
            exportToPdf(marked(props.modelValue) as string, props.title || 'Document')
            break
        case 'word':
            exportToWord(marked(props.modelValue) as string, filename, props.title)
            break
    }
    
    ElMessage.success('导出成功')
}

const restoreVersion = (version: number) => {
    versionManager.restoreVersion(version)
    ElMessage.success('版本已恢复')
}

const previewVersion = (version: number) => {
    ElMessage.info('预览功能开发中')
}

const compareVersion = (version: number) => {
    ElMessage.info('对比功能开发中')
}

const toggleProtection = () => {
    isProtected.value = !isProtected.value
    ElMessage.success(isProtected.value ? '文档已保护' : '已取消保护')
}

const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleString('zh-CN')
}

const scrollToHeading = (id: string, line: number) => {
    // 实现滚动定位
    console.log('滚动到:', id, line)
}

const handleGlobalKeydown = (e: KeyboardEvent) => {
    if (e.key === 'F1') {
        e.preventDefault()
        showShortcuts.value = true
    }
}

// 暴露方法
defineExpose({
    getContent: () => props.modelValue,
    setContent: (content: string) => emit('update:modelValue', content)
})
</script>

<style scoped lang="scss">
.markdown-editor-pro {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 200px);
    border: 1px solid #dcdfe6;
    border-radius: 4px;
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
    
    &.read-mode {
        .editor-toolbar, .editor-statusbar {
            display: none;
        }
        
        .preview-pane {
            max-width: 900px;
            margin: 0 auto;
            padding: 40px;
        }
    }
    
    .editor-toolbar {
        display: flex;
        justify-content: space-between;
        padding: 8px 12px;
        background: #f5f7fa;
        border-bottom: 1px solid #dcdfe6;
        flex-wrap: wrap;
        gap: 8px;
        
        .toolbar-group {
            display: flex;
            align-items: center;
            gap: 4px;
        }
        
        .stats-trigger {
            display: flex;
            align-items: center;
            gap: 4px;
            font-size: 12px;
            color: #909399;
            cursor: pointer;
            padding: 4px 8px;
            border-radius: 4px;
            
            &:hover {
                background: #e4e7ed;
            }
        }
    }
    
    .editor-body {
        flex: 1;
        display: flex;
        overflow: hidden;
        
        .outline-panel {
            border-right: 1px solid #dcdfe6;
            background: #fafafa;
            
            .panel-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 8px 12px;
                border-bottom: 1px solid #e4e7ed;
            }
            
            .outline-tree {
                .outline-node {
                    display: flex;
                    align-items: center;
                    gap: 8px;
                    padding: 6px 8px;
                    cursor: pointer;
                    border-radius: 4px;
                    
                    &:hover {
                        background: #f0f2f5;
                    }
                    
                    .outline-icon {
                        font-size: 12px;
                        color: #409EFF;
                    }
                    
                    .outline-text {
                        flex: 1;
                        font-size: 13px;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        white-space: nowrap;
                    }
                    
                    .outline-line {
                        font-size: 11px;
                        color: #909399;
                    }
                }
            }
        }
        
        .editor-pane {
            flex: 1;
            position: relative;
            
            .editor-textarea {
                width: 100%;
                height: 100%;
                padding: 16px;
                border: none;
                resize: none;
                font-family: 'Consolas', monospace;
                font-size: 14px;
                line-height: 1.6;
            }
        }
        
        .preview-pane {
            flex: 1;
            display: flex;
            flex-direction: column;
            
            .preview-header {
                display: flex;
                justify-content: space-between;
                padding: 8px 12px;
                background: #f5f7fa;
                border-bottom: 1px solid #e4e7ed;
            }
            
            .preview-content {
                flex: 1;
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
    
    .statistics-panel {
        h4 {
            margin: 0 0 8px 0;
            font-size: 14px;
            font-weight: 600;
        }
    }
}
</style>

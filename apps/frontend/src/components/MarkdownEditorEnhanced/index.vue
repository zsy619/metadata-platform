<template>
    <div class="markdown-editor-enhanced" :class="{ 
        'fullscreen': isFullscreen,
        'dark-mode': darkMode
    }">
        <!-- 顶部工具栏 -->
        <div class="editor-toolbar">
            <div class="toolbar-left">
                <!-- 文件操作组 -->
                <div class="toolbar-group">
                    <el-tooltip content="新建 (Ctrl+N)" placement="bottom">
                        <el-button size="small" @click="handleNew" :icon="DocumentAdd" />
                    </el-tooltip>
                    <el-tooltip content="打开模板" placement="bottom">
                        <el-button size="small" @click="showTemplateDialog = true" :icon="Document" />
                    </el-tooltip>
                    <el-tooltip content="保存 (Ctrl+S)" placement="bottom">
                        <el-button size="small" @click="handleSave" :icon="Document" :loading="saving" />
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 撤销/重做组 -->
                <div class="toolbar-group">
                    <el-tooltip content="撤销 (Ctrl+Z)" placement="bottom">
                        <el-button size="small" @click="handleUndo" :icon="RefreshLeft" :disabled="!canUndo" />
                    </el-tooltip>
                    <el-tooltip content="重做 (Ctrl+Y)" placement="bottom">
                        <el-button size="small" @click="handleRedo" :icon="RefreshRight" :disabled="!canRedo" />
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 格式组 -->
                <div class="toolbar-group">
                    <el-select
                        v-model="headingLevel"
                        size="small"
                        placeholder="标题"
                        style="width: 100px"
                        @change="insertHeading"
                    >
                        <el-option label="正文" value="p" />
                        <el-option label="标题 1" value="h1" />
                        <el-option label="标题 2" value="h2" />
                        <el-option label="标题 3" value="h3" />
                        <el-option label="标题 4" value="h4" />
                        <el-option label="标题 5" value="h5" />
                        <el-option label="标题 6" value="h6" />
                    </el-select>
                    
                    <el-tooltip content="粗体 (Ctrl+B)" placement="bottom">
                        <el-button size="small" @click="insertFormat('bold')" :icon="Edit" />
                    </el-tooltip>
                    <el-tooltip content="斜体 (Ctrl+I)" placement="bottom">
                        <el-button size="small" @click="insertFormat('italic')" :icon="ChatDotRound" />
                    </el-tooltip>
                    <el-tooltip content="删除线" placement="bottom">
                        <el-button size="small" @click="insertFormat('strikethrough')" :icon="Delete" />
                    </el-tooltip>
                    <el-tooltip content="行内代码" placement="bottom">
                        <el-button size="small" @click="insertFormat('code')" :icon="Terminal" />
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 列表组 -->
                <div class="toolbar-group">
                    <el-tooltip content="无序列表" placement="bottom">
                        <el-button size="small" @click="insertList('unordered')" :icon="List" />
                    </el-tooltip>
                    <el-tooltip content="有序列表" placement="bottom">
                    <el-button size="small" @click="insertList('ordered')" :icon="Tickets" />
                </el-tooltip>
                    <el-tooltip content="任务列表" placement="bottom">
                        <el-button size="small" @click="insertList('task')" :icon="Checked" />
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 插入组 -->
                <div class="toolbar-group">
                    <el-tooltip content="引用" placement="bottom">
                        <el-button size="small" @click="insertBlock('quote')" :icon="ChatDotRound" />
                    </el-tooltip>
                    <el-tooltip content="代码块" placement="bottom">
                    <el-button size="small" @click="insertBlock('code')" :icon="Terminal" />
                </el-tooltip>
                    <el-tooltip content="表格" placement="bottom">
                        <el-button size="small" @click="insertTable" :icon="Grid" />
                    </el-tooltip>
                    <el-tooltip content="链接 (Ctrl+K)" placement="bottom">
                        <el-button size="small" @click="insertLink" :icon="Link" />
                    </el-tooltip>
                    <el-tooltip content="图片" placement="bottom">
                        <el-button size="small" @click="imageUploadRef?.click()" :icon="Picture" />
                    </el-tooltip>
                    <el-tooltip content="表情符号" placement="bottom">
                        <el-popover
                            placement="bottom"
                            :width="320"
                            trigger="click"
                        >
                            <template #reference>
                                <el-button size="small" :icon="Odometer" />
                            </template>
                            <div class="emoji-picker">
                                <el-input
                                    v-model="emojiSearch"
                                    placeholder="搜索表情..."
                                    size="small"
                                    clearable
                                />
                                <div class="emoji-grid">
                                    <el-button
                                        v-for="emoji in filteredEmojis"
                                        :key="emoji.name"
                                        text
                                        size="large"
                                        @click="insertEmoji(emoji.emoji)"
                                    >
                                        {{ emoji.emoji }}
                                    </el-button>
                                </div>
                            </div>
                        </el-popover>
                    </el-tooltip>
                </div>
                
                <el-divider direction="vertical" />
                
                <!-- 高级功能组 -->
                <div class="toolbar-group">
                    <el-tooltip content="数学公式" placement="bottom">
                        <el-button size="small" @click="insertFormula" :icon="DataAnalysis" />
                    </el-tooltip>
                    <el-tooltip content="Mermaid 图表" placement="bottom">
                        <el-button size="small" @click="insertMermaid" :icon="Picture" />
                    </el-tooltip>
                </div>
            </div>
            
            <div class="toolbar-right">
                <!-- 主题切换 -->
                <el-tooltip :content="darkMode ? '浅色模式' : '深色模式'" placement="bottom">
                    <el-button
                        size="small"
                        @click="toggleDarkMode"
                        :icon="darkMode ? 'Sunny' : 'Moon'"
                    />
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 字数统计 -->
                <div class="word-stats">
                    <span class="stat-item" title="字数">
                        <el-icon><Document /></el-icon>
                        {{ wordCount }} 字
                    </span>
                    <span class="stat-item" title="阅读时间">
                        <el-icon><Timer /></el-icon>
                        {{ readingTime }} 分钟
                    </span>
                    <span class="stat-item" title="段落数">
                        <el-icon><Notebook /></el-icon>
                        {{ paragraphCount }} 段
                    </span>
                </div>
                <el-divider direction="vertical" />
                
                <!-- 导出菜单 -->
                <el-dropdown trigger="click" size="small">
                    <el-button size="small" :icon="Download">
                        导出
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="handleExport('markdown')">
                                <el-icon><DocumentAdd /></el-icon>
                                Markdown
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('html')">
                                <el-icon><DocumentCopy /></el-icon>
                                HTML
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('word')">
                                <el-icon><DocumentAdd /></el-icon>
                                Word
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('pdf')">
                                <el-icon><Printer /></el-icon>
                                PDF
                            </el-dropdown-item>
                            <el-dropdown-item divided @click="handleCopyContent">
                                <el-icon><DocumentCopy /></el-icon>
                                复制全文
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
                <el-divider direction="vertical" />
                
                <!-- 视图切换 -->
                <el-radio-group v-model="viewMode" size="small">
                    <el-tooltip content="编辑模式" placement="top">
                        <el-radio-button value="edit">
                            <el-icon><Edit /></el-icon>
                        </el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="预览模式" placement="top">
                        <el-radio-button value="preview">
                            <el-icon><View /></el-icon>
                        </el-radio-button>
                    </el-tooltip>
                    <el-tooltip content="分屏模式" placement="top">
                        <el-radio-button value="split">
                            <el-icon><FullScreen /></el-icon>
                        </el-radio-button>
                    </el-tooltip>
                </el-radio-group>
                <el-divider direction="vertical" />
                
                <!-- 大纲和全屏 -->
                <el-tooltip content="文档大纲" placement="bottom">
                    <el-button size="small" @click="showToc = !showToc" :icon="Menu" />
                </el-tooltip>
                <el-tooltip :content="isFullscreen ? '退出全屏' : '全屏'" placement="bottom">
                    <el-button size="small" @click="toggleFullscreen" :icon="isFullscreen ? 'FullScreen' : 'FullScreen'" />
                </el-tooltip>
            </div>
        </div>
        
        <!-- 主体内容区 -->
        <div class="editor-body">
            <!-- 左侧大纲面板 -->
            <el-aside v-if="showToc" width="220px" class="toc-panel">
                <div class="toc-header">
                    <span>文档大纲</span>
                    <el-button link size="small" @click="showToc = false">
                        <el-icon><Close /></el-icon>
                    </el-button>
                </div>
                <el-input
                    v-model="tocSearch"
                    placeholder="搜索标题..."
                    size="small"
                    clearable
                    class="toc-search"
                />
                <el-scrollbar height="calc(100vh - 320px)">
                    <div class="toc-content">
                        <el-collapse v-model="activeCollapse" accordion>
                            <el-collapse-item
                                v-for="(group, index) in groupedToc"
                                :key="index"
                                :name="index"
                            >
                                <template #title>
                                    <el-icon><Bookmark /></el-icon>
                                    {{ group.level === 1 ? '一级标题' : group.level === 2 ? '二级标题' : '其他标题' }}
                                    <el-tag size="small" type="info" style="margin-left: auto">
                                        {{ group.items.length }}
                                    </el-tag>
                                </template>
                                <div class="toc-items">
                                    <a
                                        v-for="heading in group.items"
                                        :key="heading.id"
                                        :href="`#${heading.id}`"
                                        class="toc-item"
                                        :class="`toc-level-${heading.level}`"
                                        @click="scrollToHeading(heading.id)"
                                    >
                                        <span class="toc-text">{{ heading.text }}</span>
                                    </a>
                                </div>
                            </el-collapse-item>
                        </el-collapse>
                        <el-empty 
                            v-if="tableOfContents.length === 0" 
                            description="暂无大纲" 
                            :image-size="60" 
                        />
                    </div>
                </el-scrollbar>
            </el-aside>
            
            <!-- 编辑区 -->
            <div 
                v-show="viewMode === 'edit' || viewMode === 'split'" 
                class="editor-pane"
                @dragover.prevent="handleDragOver"
                @dragleave.prevent
                @drop.prevent="handleDrop"
            >
                <textarea
                    ref="editorRef"
                    v-model="modelValue"
                    class="editor-textarea"
                    placeholder="开始编写 Markdown 文档... (支持拖拽图片到此处)"
                    @scroll="handleScroll"
                    @keydown.ctrl.s.prevent="handleSave"
                    @keydown.ctrl.b.prevent="insertFormat('bold')"
                    @keydown.ctrl.i.prevent="insertFormat('italic')"
                    @keydown.ctrl.k.prevent="insertLink"
                    @paste="handlePaste"
                />
                <!-- 拖拽上传提示 -->
                <div v-if="isDragging" class="drag-overlay">
                    <el-icon class="drag-icon"><UploadFilled /></el-icon>
                    <div class="drag-text">释放以上传图片</div>
                </div>
            </div>
            
            <!-- 预览区 -->
            <div v-show="viewMode === 'preview' || viewMode === 'split'" class="preview-pane">
                <div class="preview-header">
                    <span>预览</span>
                    <div class="preview-actions">
                        <el-tooltip content="放大" placement="top">
                            <el-button size="small" text @click="zoomIn">
                                <el-icon><ZoomIn /></el-icon>
                            </el-button>
                        </el-tooltip>
                        <el-tooltip content="缩小" placement="top">
                            <el-button size="small" text @click="zoomOut">
                                <el-icon><ZoomOut /></el-icon>
                            </el-button>
                        </el-tooltip>
                        <el-tooltip content="重置" placement="top">
                            <el-button size="small" text @click="zoomReset">
                                <el-icon><RefreshRight /></el-icon>
                            </el-button>
                        </el-tooltip>
                    </div>
                </div>
                <div class="preview-content" v-html="renderedContent" :style="{ zoom: previewZoom }"></div>
            </div>
        </div>
        
        <!-- 底部状态栏 -->
        <div class="editor-statusbar">
            <div class="status-left">
                <span class="status-item">
                    <el-icon><Document /></el-icon>
                    Markdown
                </span>
                <el-divider direction="vertical" />
                <span class="status-item">
                    <el-icon><Position /></el-icon>
                    行 {{ cursorPosition.line }}, 列 {{ cursorPosition.column }}
                </span>
                <el-divider direction="vertical" />
                <span class="status-item" v-if="autoSaveTime">
                    <el-icon><Clock /></el-icon>
                    上次保存：{{ autoSaveTime }}
                </span>
            </div>
            <div class="status-right">
                <el-tooltip content="拼写检查" placement="top">
                    <el-switch
                        v-model="spellCheck"
                        size="small"
                        inline-prompt
                        active-text="开"
                        inactive-text="关"
                    />
                </el-tooltip>
            </div>
        </div>
        
        <!-- 图片上传隐藏输入框 -->
        <input
            ref="imageUploadRef"
            type="file"
            accept="image/*"
            multiple
            style="display: none"
            @change="handleImageUpload"
        />
        
        <!-- 模板选择对话框 -->
        <el-dialog
            v-model="showTemplateDialog"
            title="选择文档模板"
            width="900px"
            :close-on-click-modal="false"
        >
            <div class="template-dialog-content">
                <el-tabs v-model="templateCategory">
                    <el-tab-pane
                        v-for="category in templateCategories"
                        :key="category"
                        :label="category"
                        :name="category"
                    >
                        <el-row :gutter="20">
                            <el-col
                                v-for="template in getTemplates(category)"
                                :key="template.id"
                                :span="8"
                            >
                                <el-card
                                    class="template-card"
                                    shadow="hover"
                                    @click="selectTemplate(template)"
                                >
                                    <template #header>
                                        <div class="template-card-header">
                                            <el-icon :size="24"><component :is="template.icon" /></el-icon>
                                            <span>{{ template.name }}</span>
                                        </div>
                                    </template>
                                    <div class="template-card-body">
                                        <p class="template-description">{{ template.description }}</p>
                                        <div class="template-preview">
                                            <pre>{{ template.content.substring(0, 100) }}...</pre>
                                        </div>
                                    </div>
                                </el-card>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
            </div>
            <template #footer>
                <el-button @click="showTemplateDialog = false">取消</el-button>
                <el-button type="primary" @click="applyTemplate">使用模板</el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
    DocumentAdd, Document, RefreshLeft, Plus,
    Delete, Edit, Check, View, Menu, Download,
    Printer
} from '@element-plus/icons-vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import {
    exportToMarkdown, exportToHtml, exportToPdf, exportToWord, copyToClipboard
} from '@/utils/documentExport'
import {
    emojiList, compressImage, validateImageFile,
    getTemplates, getTemplateById, DocumentTemplate
} from '@/utils/markdownEnhancer'

// Props
interface Props {
    modelValue: string
    title?: string
    path?: string
}

const props = withDefaults(defineProps<Props>(), {
    title: '',
    path: ''
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: string]
    'save': [content: string]
    'image-upload': [file: File]
}>()

// 编辑器引用
const editorRef = ref<HTMLTextAreaElement | null>(null)
const imageUploadRef = ref<HTMLInputElement | null>(null)

// 状态
const viewMode = ref<'edit' | 'preview' | 'split'>('split')
const isFullscreen = ref(false)
const darkMode = ref(false)
const saving = ref(false)
const autoSaveTime = ref('')
const showToc = ref(false)
const tocSearch = ref('')
const headingLevel = ref('p')
const spellCheck = ref(false)
const isDragging = ref(false)
const emojiSearch = ref('')

// 模板相关
const showTemplateDialog = ref(false)
const templateCategory = ref('技术文档')
const selectedTemplate = ref<DocumentTemplate | null>(null)

const templateCategories = computed(() => {
    const categories = new Set(documentTemplates.map(t => t.category))
    return Array.from(categories)
})

// 历史记录
const history = reactive({
    stack: [] as string[],
    index: -1,
    maxSize: 100
})

// 光标位置
const cursorPosition = reactive({
    line: 1,
    column: 1
})

// 预览缩放
const previewZoom = ref(1)

// 大纲折叠
const activeCollapse = ref<number | null>(0)

// 配置 marked
marked.setOptions({
    gfm: true,
    breaks: true,
    highlight: (code, lang) => {
        if (lang && hljs.getLanguage(lang)) {
            return hljs.highlight(code, { language: lang }).value
        }
        return hljs.highlightAuto(code).value
    }
})

// 计算属性
const modelValue = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const wordCount = computed(() => {
    const text = props.modelValue
    const chineseChars = (text.match(/[\u4e00-\u9fa5]/g) || []).length
    const englishWords = (text.match(/\b\w+\b/g) || []).length
    return chineseChars + englishWords
})

const readingTime = computed(() => {
    const minutes = Math.ceil(wordCount.value / 300)
    return Math.max(1, minutes)
})

const paragraphCount = computed(() => {
    return props.modelValue.split('\n\n').filter(p => p.trim()).length
})

const canUndo = computed(() => history.index > 0)
const canRedo = computed(() => history.index < history.stack.length - 1)

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

const tableOfContents = computed(() => {
    const headings: Array<{ id: string; level: number; text: string }> = []
    const lines = props.modelValue.split('\n')
    
    lines.forEach((line, index) => {
        const match = line.match(/^(#{1,6})\s+(.+)$/)
        if (match) {
            const level = match[1].length
            const text = match[2].replace(/[*_~`]/g, '')
            const id = `heading-${index}`
            headings.push({ id, level, text })
        }
    })
    
    return headings
})

const filteredEmojis = computed(() => {
    if (!emojiSearch.value) {
        return emojiList.slice(0, 50)
    }
    return emojiList.filter(e => 
        e.name.toLowerCase().includes(emojiSearch.value.toLowerCase())
    )
})

const groupedToc = computed(() => {
    const filtered = tableOfContents.value.filter(h =>
        h.text.toLowerCase().includes(tocSearch.value.toLowerCase())
    )
    
    const groups: Array<{ level: number; items: Array<typeof filtered[0]> }> = [
        { level: 1, items: [] },
        { level: 2, items: [] },
        { level: 3, items: [] }
    ]
    
    filtered.forEach(h => {
        if (h.level <= 3) {
            groups[h.level - 1].items.push(h)
        }
    })
    
    return groups.filter(g => g.items.length > 0)
})

// 监听内容变化
let autoSaveTimer: NodeJS.Timeout | null = null
watch(() => props.modelValue, (newValue) => {
    if (history.index < history.stack.length - 1) {
        history.stack = history.stack.slice(0, history.index + 1)
    }
    history.stack.push(newValue)
    if (history.stack.length > history.maxSize) {
        history.stack.shift()
    } else {
        history.index++
    }
    
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer)
    }
    autoSaveTimer = setTimeout(() => {
        autoSaveTime.value = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
        localStorage.setItem('markdown-draft-enhanced', props.modelValue)
    }, 2000)
}, { immediate: true })

// 生命周期
onMounted(() => {
    const draft = localStorage.getItem('markdown-draft-enhanced')
    if (draft && !props.modelValue) {
        emit('update:modelValue', draft)
        ElMessage.info('已恢复上次草稿')
    }
    
    if (props.modelValue) {
        history.stack.push(props.modelValue)
        history.index = 0
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
const updateCursorPosition = () => {
    if (!editorRef.value) return
    const pos = editorRef.value.selectionStart
    const text = props.modelValue.substring(0, pos)
    const lines = text.split('\n')
    cursorPosition.line = lines.length
    cursorPosition.column = lines[lines.length - 1].length + 1
}

const handleScroll = (event: Event) => {
    const target = event.target as HTMLTextAreaElement
    const previewPane = document.querySelector('.preview-content')
    if (previewPane) {
        const percentage = target.scrollTop / (target.scrollHeight - target.clientHeight)
        previewPane.scrollTop = percentage * (previewPane.scrollHeight - previewPane.clientHeight)
    }
}

const handleGlobalKeydown = (event: KeyboardEvent) => {
    if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
        return
    }
    
    if (event.ctrlKey || event.metaKey) {
        switch (event.key.toLowerCase()) {
            case 's':
                event.preventDefault()
                handleSave()
                break
            case 'n':
                event.preventDefault()
                handleNew()
                break
        }
    }
}

const insertText = (before: string, after: string = '', placeholder: string = '') => {
    if (!editorRef.value) return
    
    const start = editorRef.value.selectionStart
    const end = editorRef.value.selectionEnd
    const text = props.modelValue
    
    const selectedText = text.substring(start, end)
    let newText: string
    
    if (selectedText) {
        newText = text.substring(0, start) + before + selectedText + after + text.substring(end)
    } else {
        newText = text.substring(0, start) + before + placeholder + after + text.substring(end)
    }
    
    emit('update:modelValue', newText)
    
    setTimeout(() => {
        editorRef.value?.focus()
        const newStart = start + before.length
        const newEnd = selectedText ? newStart + selectedText.length : newStart + placeholder.length
        editorRef.value?.setSelectionRange(newStart, newEnd)
    }, 0)
}

// 工具栏方法
const handleNew = () => {
    emit('update:modelValue', '')
    ElMessage.success('已新建文档')
}

const handleSave = () => {
    saving.value = true
    emit('save', props.modelValue)
    setTimeout(() => {
        saving.value = false
        ElMessage.success('保存成功')
    }, 500)
}

const handleUndo = () => {
    if (canUndo.value) {
        history.index--
        emit('update:modelValue', history.stack[history.index])
    }
}

const handleRedo = () => {
    if (canRedo.value) {
        history.index++
        emit('update:modelValue', history.stack[history.index])
    }
}

const insertHeading = () => {
    if (headingLevel.value === 'p') {
        insertText('', '', '正文内容')
    } else {
        const level = parseInt(headingLevel.value.charAt(1))
        const hashes = '#'.repeat(level)
        insertText(`${hashes} `, '', '标题内容')
    }
    headingLevel.value = 'p'
}

const insertFormat = (format: string) => {
    switch (format) {
        case 'bold':
            insertText('**', '**', '粗体文本')
            break
        case 'italic':
            insertText('*', '*', '斜体文本')
            break
        case 'strikethrough':
            insertText('~~', '~~', '删除线文本')
            break
        case 'code':
            insertText('`', '`', '代码')
            break
    }
}

const insertList = (type: string) => {
    switch (type) {
        case 'unordered':
            insertText('- ', '', '列表项')
            break
        case 'ordered':
            insertText('1. ', '', '列表项')
            break
        case 'task':
            insertText('- [ ] ', '', '任务项')
            break
    }
}

const insertBlock = (type: string) => {
    switch (type) {
        case 'quote':
            insertText('> ', '', '引用内容')
            break
        case 'code':
            insertText('```\n', '\n```\n', '代码')
            break
    }
}

const insertLink = () => {
    const url = prompt('请输入链接地址:')
    if (url) {
        const text = prompt('请输入链接文本:', '链接') || '链接'
        insertText('[', `](${url})`, text)
    }
}

const insertImage = () => {
    imageUploadRef.value?.click()
}

const handleImageUpload = async (event: Event) => {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (file) {
        const validation = validateImageFile(file)
        if (!validation.valid) {
            ElMessage.error(validation.error)
            return
        }
        
        try {
            const compressed = await compressImage(file, {
                maxWidth: 1920,
                maxHeight: 1080,
                quality: 0.8
            })
            
            emit('image-upload', compressed)
            const url = URL.createObjectURL(compressed)
            insertText('![', `](${url})`, '图片描述')
            ElMessage.success('图片上传成功')
        } catch (error) {
            ElMessage.error('图片处理失败')
        }
    }
    input.value = ''
}

// 拖拽上传
const handleDragOver = () => {
    isDragging.value = true
}

const handleDrop = async (event: DragEvent) => {
    isDragging.value = false
    const files = event.dataTransfer?.files
    if (files && files.length > 0) {
        const file = files[0]
        if (file.type.startsWith('image/')) {
            await handleImageUpload({ target: { files: [file] } } as any)
        } else {
            ElMessage.warning('请上传图片文件')
        }
    }
}

// 粘贴上传
const handlePaste = async (event: ClipboardEvent) => {
    const items = event.clipboardData?.items
    if (items) {
        for (let i = 0; i < items.length; i++) {
            if (items[i].type.startsWith('image/')) {
                const file = items[i].getAsFile()
                if (file) {
                    event.preventDefault()
                    await handleImageUpload({ target: { files: [file] } } as any)
                    break
                }
            }
        }
    }
}

const insertTable = () => {
    const table = `
| 列 1 | 列 2 | 列 3 |
|------|------|------|
| 单元格 1 | 单元格 2 | 单元格 3 |
| 单元格 4 | 单元格 5 | 单元格 6 |
`.trim()
    insertText(table + '\n\n')
}

const insertFormula = () => {
    insertText('$$\n', '\n$$\n', 'E = mc^2')
}

const insertMermaid = () => {
    const mermaid = `
\`\`\`mermaid
graph LR
    A[开始] --> B{条件}
    B -->|是 | C[操作 1]
    B -->|否 | D[操作 2]
    C --> E[结束]
    D --> E
\`\`\`
`.trim()
    insertText(mermaid + '\n\n')
}

const insertEmoji = (emoji: string) => {
    insertText(emoji + ' ')
}

const toggleFullscreen = () => {
    isFullscreen.value = !isFullscreen.value
    if (isFullscreen.value) {
        document.documentElement.requestFullscreen?.()
    } else {
        document.exitFullscreen?.()
    }
}

const toggleDarkMode = () => {
    darkMode.value = !darkMode.value
    if (darkMode.value) {
        document.body.classList.add('dark-mode')
    } else {
        document.body.classList.remove('dark-mode')
    }
}

const scrollToHeading = (headingId: string) => {
    const element = document.getElementById(headingId)
    element?.scrollIntoView({ behavior: 'smooth' })
}

// 导出功能
const handleExport = (type: string) => {
    const filename = props.title || 'document'
    const content = props.modelValue
    
    switch (type) {
        case 'markdown':
            exportToMarkdown(content, filename)
            ElMessage.success('已导出为 Markdown')
            break
        case 'html':
            const htmlContent = marked(content)
            exportToHtml(htmlContent as string, filename, props.title || 'Document')
            ElMessage.success('已导出为 HTML')
            break
        case 'word':
            const wordContent = marked(content)
            exportToWord(wordContent as string, filename, props.title || 'Document')
            ElMessage.success('已导出为 Word')
            break
        case 'pdf':
            const pdfContent = marked(content)
            exportToPdf(pdfContent as string, props.title || 'Document')
            break
    }
}

const handleCopyContent = async () => {
    const success = await copyToClipboard(props.modelValue)
    if (success) {
        ElMessage.success('已复制到剪贴板')
    } else {
        ElMessage.error('复制失败')
    }
}

// 模板功能
const selectTemplate = (template: DocumentTemplate) => {
    selectedTemplate.value = template
}

const applyTemplate = () => {
    if (selectedTemplate.value) {
        emit('update:modelValue', selectedTemplate.value.content)
        showTemplateDialog.value = false
        ElMessage.success('模板已应用')
    }
}

// 预览缩放
const zoomIn = () => {
    previewZoom.value = Math.min(2, previewZoom.value + 0.1)
}

const zoomOut = () => {
    previewZoom.value = Math.max(0.5, previewZoom.value - 0.1)
}

const zoomReset = () => {
    previewZoom.value = 1
}

defineExpose({
    getContent: () => props.modelValue,
    setContent: (content: string) => emit('update:modelValue', content)
})
</script>

<style scoped lang="scss">
.markdown-editor-enhanced {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 200px);
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    background: #fff;
    transition: all 0.3s;
    
    &.fullscreen {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: 9999;
        height: 100vh;
        border-radius: 0;
    }
    
    &.dark-mode {
        background: #1e1e1e;
        color: #d4d4d4;
        
        .editor-toolbar {
            background: #2d2d2d;
            border-bottom-color: #404040;
        }
        
        .editor-textarea {
            background: #1e1e1e;
            color: #d4d4d4;
        }
        
        .preview-content {
            background: #1e1e1e;
            color: #d4d4d4;
        }
    }
    
    .editor-toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px 12px;
        background: #f5f7fa;
        border-bottom: 1px solid #dcdfe6;
        flex-wrap: wrap;
        gap: 8px;
        
        .toolbar-left, .toolbar-right {
            display: flex;
            align-items: center;
            gap: 8px;
            flex-wrap: wrap;
        }
        
        .toolbar-group {
            display: flex;
            align-items: center;
            gap: 4px;
        }
        
        .word-stats {
            display: flex;
            gap: 12px;
            font-size: 12px;
            color: #909399;
            
            .stat-item {
                display: flex;
                align-items: center;
                gap: 4px;
            }
        }
    }
    
    .editor-body {
        flex: 1;
        display: flex;
        overflow: hidden;
        position: relative;
        
        .toc-panel {
            border-right: 1px solid #dcdfe6;
            background: #fafafa;
            display: flex;
            flex-direction: column;
            
            .toc-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 8px 12px;
                border-bottom: 1px solid #e4e7ed;
                font-size: 14px;
                font-weight: 500;
            }
            
            .toc-search {
                padding: 8px;
                border-bottom: 1px solid #e4e7ed;
            }
            
            .toc-content {
                flex: 1;
                overflow-y: auto;
                padding: 8px;
                
                .toc-items {
                    display: flex;
                    flex-direction: column;
                    gap: 4px;
                    
                    .toc-item {
                        display: block;
                        padding: 6px 8px;
                        color: #606266;
                        text-decoration: none;
                        font-size: 13px;
                        border-radius: 4px;
                        cursor: pointer;
                        white-space: nowrap;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        transition: all 0.2s;
                        
                        &:hover {
                            background: #f0f2f5;
                            color: #409EFF;
                        }
                        
                        @for $i from 1 through 3 {
                            &.toc-level-#{$i} {
                                padding-left: #{($i - 1) * 16 + 8}px;
                            }
                        }
                    }
                }
            }
        }
        
        .editor-pane {
            flex: 1;
            overflow: hidden;
            position: relative;
            
            .editor-textarea {
                width: 100%;
                height: 100%;
                padding: 16px;
                border: none;
                resize: none;
                font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
                font-size: 14px;
                line-height: 1.6;
                outline: none;
                
                &:focus {
                    background: #fafafa;
                }
            }
            
            .drag-overlay {
                position: absolute;
                top: 0;
                left: 0;
                right: 0;
                bottom: 0;
                background: rgba(64, 158, 255, 0.1);
                border: 2px dashed #409EFF;
                display: flex;
                flex-direction: column;
                justify-content: center;
                align-items: center;
                z-index: 10;
                
                .drag-icon {
                    font-size: 64px;
                    color: #409EFF;
                    margin-bottom: 16px;
                }
                
                .drag-text {
                    font-size: 16px;
                    color: #409EFF;
                }
            }
        }
        
        .preview-pane {
            flex: 1;
            overflow: hidden;
            display: flex;
            flex-direction: column;
            
            .preview-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 8px 12px;
                background: #f5f7fa;
                border-bottom: 1px solid #e4e7ed;
                font-size: 14px;
                font-weight: 500;
                
                .preview-actions {
                    display: flex;
                    gap: 4px;
                }
            }
            
            .preview-content {
                flex: 1;
                padding: 16px;
                overflow-y: auto;
                font-size: 14px;
                line-height: 1.6;
                background: #fff;
            }
        }
    }
    
    .editor-statusbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 4px 12px;
        background: #f5f7fa;
        border-top: 1px solid #dcdfe6;
        font-size: 12px;
        color: #909399;
        
        .status-left, .status-right {
            display: flex;
            align-items: center;
            gap: 8px;
        }
    }
    
    .emoji-picker {
        .emoji-grid {
            display: grid;
            grid-template-columns: repeat(8, 1fr);
            gap: 4px;
            max-height: 300px;
            overflow-y: auto;
            margin-top: 8px;
        }
    }
}

.template-dialog-content {
    .template-card {
        margin-bottom: 16px;
        cursor: pointer;
        transition: all 0.3s;
        
        &:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
        
        .template-card-header {
            display: flex;
            align-items: center;
            gap: 8px;
            font-weight: 500;
        }
        
        .template-card-body {
            .template-description {
                font-size: 13px;
                color: #606266;
                margin-bottom: 8px;
            }
            
            .template-preview {
                background: #f5f7fa;
                padding: 8px;
                border-radius: 4px;
                font-size: 12px;
                
                pre {
                    margin: 0;
                    white-space: pre-wrap;
                    word-break: break-all;
                }
            }
        }
    }
}
</style>

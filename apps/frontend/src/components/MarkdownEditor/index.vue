<template>
    <div class="markdown-editor" :class="{ 'fullscreen': isFullscreen }">
        <!-- 顶部工具栏 -->
        <div class="editor-toolbar">
            <div class="toolbar-left">
                <!-- 文件操作 -->
                <el-tooltip content="新建 (Ctrl+N)" placement="bottom">
                    <el-button size="small" @click="handleNew" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'file-circle-plus']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="保存 (Ctrl+S)" placement="bottom">
                    <el-button size="small" @click="handleSave" :loading="saving" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'floppy-disk']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 撤销/重做 -->
                <el-tooltip content="撤销 (Ctrl+Z)" placement="bottom">
                    <el-button size="small" @click="handleUndo" :disabled="!canUndo" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'rotate-left']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="重做 (Ctrl+Y)" placement="bottom">
                    <el-button size="small" @click="handleRedo" :disabled="!canRedo" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'rotate-right']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 标题 -->
                <el-select
                    v-model="headingLevel"
                    size="small"
                    placeholder="正文"
                    style="width: 80px"
                    @change="insertHeading"
                    class="toolbar-select"
                >
                    <el-option label="正文" value="p" />
                    <el-option label="标题 1" value="h1" />
                    <el-option label="标题 2" value="h2" />
                    <el-option label="标题 3" value="h3" />
                    <el-option label="标题 4" value="h4" />
                    <el-option label="标题 5" value="h5" />
                    <el-option label="标题 6" value="h6" />
                </el-select>
                <el-divider direction="vertical" />
                
                <!-- 格式 -->
                <el-tooltip content="粗体 (Ctrl+B)" placement="bottom">
                    <el-button size="small" @click="insertFormat('bold')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'bold']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="斜体 (Ctrl+I)" placement="bottom">
                    <el-button size="small" @click="insertFormat('italic')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'italic']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="删除线" placement="bottom">
                    <el-button size="small" @click="insertFormat('strikethrough')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'strikethrough']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="行内代码" placement="bottom">
                    <el-button size="small" @click="insertFormat('code')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'code']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 列表 -->
                <el-tooltip content="无序列表" placement="bottom">
                    <el-button size="small" @click="insertList('unordered')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'list-ul']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="有序列表" placement="bottom">
                    <el-button size="small" @click="insertList('ordered')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'list-ol']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="任务列表" placement="bottom">
                    <el-button size="small" @click="insertList('task')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'square-check']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 引用与代码块 -->
                <el-tooltip content="引用" placement="bottom">
                    <el-button size="small" @click="insertBlock('quote')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'quote-left']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="代码块" placement="bottom">
                    <el-button size="small" @click="insertBlock('code')" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'code']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 插入 -->
                <el-tooltip content="链接 (Ctrl+K)" placement="bottom">
                    <el-button size="small" @click="insertLink" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'link']" />
                    </el-button>
                </el-tooltip>
                <el-tooltip content="图片" placement="bottom">
                    <el-button size="small" @click="insertImage" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'image']" />
                    </el-button>
                </el-tooltip>
                <div class="table-popover-wrapper">
                    <el-tooltip content="表格" placement="bottom">
                        <el-button 
                            size="small" 
                            class="toolbar-icon-btn"
                            @click="toggleTablePopover"
                            :class="{ 'active': showTablePopover }"
                        >
                            <font-awesome-icon :icon="['fas', 'table']" />
                        </el-button>
                    </el-tooltip>
                    
                    <!-- 自定义表格选择器弹出层 -->
                    <div 
                        v-show="showTablePopover" 
                        class="table-popover-content"
                        @mouseenter="showTablePopover = true"
                        @mouseleave="showTablePopover = false"
                    >
                        <div class="table-grid-selector">
                            <div class="table-grid" @mousemove="handleGridMove" @click="handleGridSelect">
                                <div 
                                    v-for="row in 8" 
                                    :key="row" 
                                    class="table-row"
                                >
                                    <div 
                                        v-for="col in 8" 
                                        :key="col"
                                        class="table-cell"
                                        :class="{ active: isCellActive(row, col) }"
                                        :data-row="row"
                                        :data-col="col"
                                    />
                                </div>
                            </div>
                            <div class="table-size-info">
                                {{ selectedRows }} × {{ selectedCols }} 表格
                            </div>
                        </div>
                    </div>
                </div>
                <el-tooltip content="数学公式" placement="bottom">
                    <el-button size="small" @click="insertFormula" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'square-root-variable']" />
                    </el-button>
                </el-tooltip>
                <el-divider direction="vertical" />
                
                <!-- 大纲 -->
                <el-tooltip content="文档大纲" placement="bottom">
                    <el-button size="small" @click="showToc = !showToc" class="toolbar-icon-btn">
                        <font-awesome-icon :icon="['fas', 'list']" />
                    </el-button>
                </el-tooltip>
            </div>
            
            <div class="toolbar-right">
                <!-- 字数统计 -->
                <span class="word-count">
                    {{ wordCount }} 字 | {{ readingTime }} 分钟
                </span>
                <el-divider direction="vertical" />
                
                <!-- 导出菜单 -->
                <el-dropdown trigger="click" size="small">
                    <el-button size="small">
                        <font-awesome-icon icon="download" />
                        导出
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item @click="handleExport('markdown')">
                                <font-awesome-icon icon="file-code" />
                                Markdown
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('html')">
                                <font-awesome-icon icon="file-code" />
                                HTML
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('word')">
                                <font-awesome-icon icon="file-word" />
                                Word
                            </el-dropdown-item>
                            <el-dropdown-item @click="handleExport('pdf')">
                                <font-awesome-icon icon="file-pdf" />
                                PDF
                            </el-dropdown-item>
                            <el-dropdown-item divided @click="handleCopyContent">
                                <font-awesome-icon icon="copy" />
                                复制全文
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
                <el-divider direction="vertical" />
                
                <!-- 视图切换 -->
                <el-radio-group v-model="viewMode" size="small">
                    <el-radio-button value="edit">
                        <font-awesome-icon icon="edit" />
                        编辑
                    </el-radio-button>
                    <el-radio-button value="preview">
                        <font-awesome-icon icon="eye" />
                        预览
                    </el-radio-button>
                    <el-radio-button value="split">
                        <font-awesome-icon icon="columns" />
                        分屏
                    </el-radio-button>
                </el-radio-group>
                <el-divider direction="vertical" />
                
                <!-- 全屏 -->
                <el-tooltip :content="isFullscreen ? '退出全屏' : '全屏'" placement="bottom">
                    <el-button size="small" @click="toggleFullscreen">
                        <font-awesome-icon :icon="isFullscreen ? 'compress' : 'expand'" />
                    </el-button>
                </el-tooltip>
            </div>
        </div>
        
        <!-- 主体内容区 -->
        <div class="editor-body">
            <!-- 左侧大纲面板 -->
            <el-aside v-if="showToc" width="200px" class="toc-panel">
                <div class="toc-header">
                    <span>文档大纲</span>
                    <el-button link size="small" @click="showToc = false">
                        <font-awesome-icon icon="times" />
                    </el-button>
                </div>
                <el-scrollbar height="calc(100vh - 300px)">
                    <div class="toc-content">
                        <a
                            v-for="heading in tableOfContents"
                            :key="heading.id"
                            :href="`#${heading.id}`"
                            class="toc-item"
                            :class="`toc-level-${heading.level}`"
                            @click="scrollToHeading(heading.id)"
                        >
                            {{ heading.text }}
                        </a>
                        <el-empty v-if="tableOfContents.length === 0" description="暂无大纲" :image-size="60" />
                    </div>
                </el-scrollbar>
            </el-aside>
            
            <!-- 编辑区 -->
            <div v-show="viewMode === 'edit' || viewMode === 'split'" class="editor-pane">
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
                />
            </div>
            
            <!-- 预览区 -->
            <div v-show="viewMode === 'preview' || viewMode === 'split'" class="preview-pane">
                <div class="preview-content" v-html="renderedContent"></div>
            </div>
        </div>
        
        <!-- 底部状态栏 -->
        <div class="editor-statusbar">
            <span>Markdown</span>
            <el-divider direction="vertical" />
            <span>Ln {{ cursorPosition.line }}, Col {{ cursorPosition.column }}</span>
            <el-divider direction="vertical" />
            <span v-if="autoSaveTime">上次保存：{{ autoSaveTime }}</span>
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
    </div>
</template>

<script setup lang="ts">
import { copyToClipboard, exportToHtml, exportToMarkdown, exportToPdf, exportToWord } from '@/utils/documentExport'
import { ElDropdown, ElDropdownItem, ElDropdownMenu, ElMessage } from 'element-plus'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import { marked } from 'marked'
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue'

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
const saving = ref(false)
const autoSaveTime = ref('')
const showToc = ref(false)
const headingLevel = ref('p')

// 表格选择器状态
const selectedRows = ref(1)
const selectedCols = ref(1)
const hoverRows = ref(1)
const hoverCols = ref(1)
const showTablePopover = ref(false)

// 显示/隐藏表格选择器
const toggleTablePopover = () => {
    showTablePopover.value = !showTablePopover.value
}

// 历史记录（撤销/重做）
const history = reactive({
    stack: [] as string[],
    index: -1,
    maxSize: 50
})

// 光标位置
const cursorPosition = reactive({
    line: 1,
    column: 1
})

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
    const text = props.modelValue.replace(/[#*_~`>\[\]()]/g, '').trim()
    return text.length
})

const readingTime = computed(() => {
    const minutes = Math.ceil(wordCount.value / 300)
    return Math.max(1, minutes)
})

const canUndo = computed(() => history.index > 0)
const canRedo = computed(() => history.index < history.stack.length - 1)

// 渲染 Markdown
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

// 文档大纲
const tableOfContents = computed(() => {
    const headings: Array<{ id: string; level: number; text: string }> = []
    const lines = props.modelValue.split('\n')
    
    lines.forEach((line, index) => {
        const match = line.match(/^(#{1,6})\s+(.+)$/)
        if (match) {
            const level = match[1].length
            const text = match[2]
            const id = `heading-${index}`
            headings.push({ id, level, text })
        }
    })
    
    return headings
})

// 监听内容变化，记录历史
let autoSaveTimer: NodeJS.Timeout | null = null
watch(() => props.modelValue, (newValue) => {
    // 添加到历史记录
    if (history.index < history.stack.length - 1) {
        history.stack = history.stack.slice(0, history.index + 1)
    }
    history.stack.push(newValue)
    if (history.stack.length > history.maxSize) {
        history.stack.shift()
    } else {
        history.index++
    }
    
    // 自动保存（防抖）
    if (autoSaveTimer) {
        clearTimeout(autoSaveTimer)
    }
    autoSaveTimer = setTimeout(() => {
        autoSaveTime.value = new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
        localStorage.setItem('markdown-draft', props.modelValue)
    }, 2000)
}, { immediate: true })

// 监听光标位置变化
watch(() => props.modelValue, () => {
    updateCursorPosition()
})

// 生命周期
onMounted(() => {
    // 恢复草稿
    const draft = localStorage.getItem('markdown-draft')
    if (draft && !props.modelValue) {
        emit('update:modelValue', draft)
        ElMessage.info('已恢复上次草稿')
    }
    
    // 初始化历史记录
    if (props.modelValue) {
        history.stack.push(props.modelValue)
        history.index = 0
    }
    
    // 添加键盘快捷键监听
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
    // 同步滚动逻辑
    const target = event.target as HTMLTextAreaElement
    const previewPane = document.querySelector('.preview-content')
    if (previewPane) {
        const percentage = target.scrollTop / (target.scrollHeight - target.clientHeight)
        previewPane.scrollTop = percentage * (previewPane.scrollHeight - previewPane.clientHeight)
    }
}

const handleGlobalKeydown = (event: KeyboardEvent) => {
    // 避免在输入框中触发
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
    
    // 恢复光标
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

const handleImageUpload = (event: Event) => {
    const input = event.target as HTMLInputElement
    const file = input.files?.[0]
    if (file) {
        emit('image-upload', file)
        // 这里可以处理图片上传逻辑
        const url = URL.createObjectURL(file)
        insertText('![', `](${url})`, '图片描述')
    }
    input.value = ''
}

const insertFormula = () => {
    insertText('$$\n', '\n$$\n', 'E = mc^2')
}

// 表格拖拽选择相关函数
const handleGridMove = (event: MouseEvent) => {
    const cell = (event.target as HTMLElement).closest('.table-cell') as HTMLElement
    if (cell) {
        const row = parseInt(cell.getAttribute('data-row') || '1')
        const col = parseInt(cell.getAttribute('data-col') || '1')
        hoverRows.value = row
        hoverCols.value = col
        selectedRows.value = row
        selectedCols.value = col
    }
}

const handleGridSelect = () => {
    if (hoverRows.value > 0 && hoverCols.value > 0) {
        insertTableWithSize(hoverRows.value, hoverCols.value)
    }
}

const isCellActive = (row: number, col: number) => {
    return row <= hoverRows.value && col <= hoverCols.value
}

const insertTableWithSize = (rows: number, cols: number) => {
    // 生成 Markdown 表格
    const headers: string[] = []
    const separator: string[] = []
    const rowData: string[] = []
    
    // 生成表头
    for (let c = 1; c <= cols; c++) {
        headers.push(`列 ${c}`)
        separator.push('------')
    }
    
    // 生成数据行
    for (let r = 1; r <= rows; r++) {
        const cells: string[] = []
        for (let c = 1; c <= cols; c++) {
            cells.push(`单元格 ${r}-${c}`)
        }
        rowData.push(`| ${cells.join(' | ')} |`)
    }
    
    const table = `| ${headers.join(' | ')} |\n| ${separator.join(' | ')} |\n${rowData.join('\n')}`
    insertText(table + '\n\n')
}

// 重写 insertTable 函数，使用新的拖拽选择
const insertTable = () => {
    insertTableWithSize(3, 3) // 默认 3x3
}

const toggleFullscreen = () => {
    isFullscreen.value = !isFullscreen.value
    if (isFullscreen.value) {
        document.documentElement.requestFullscreen?.()
    } else {
        document.exitFullscreen?.()
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

// 复制全文
const handleCopyContent = async () => {
    const success = await copyToClipboard(props.modelValue)
    if (success) {
        ElMessage.success('已复制到剪贴板')
    } else {
        ElMessage.error('复制失败')
    }
}

// 暴露方法给父组件
defineExpose({
    getContent: () => props.modelValue,
    setContent: (content: string) => emit('update:modelValue', content)
})
</script>

<style scoped lang="scss">
.markdown-editor {
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
    
    .editor-toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 3px 6px;
        background: #f5f7fa;
        border-bottom: 1px solid #dcdfe6;
        
        .toolbar-left, .toolbar-right {
            display: flex;
            align-items: center;
            gap: 1px;
            flex-wrap: wrap;
            
            .el-button {
                padding: 3px 5px;
                min-width: auto;
                height: 26px;
                font-size: 12px;
                border-radius: 3px;
                
                .font-awesome-icon {
                    font-size: 12px;
                }
            }
            
            .el-select {
                height: 26px;
                font-size: 12px;
                
                :deep(.el-input__wrapper) {
                    padding: 0 6px;
                    height: 26px;
                    border-radius: 3px;
                }
                
                :deep(.el-input__inner) {
                    font-size: 12px;
                }
            }
            
            .el-divider {
                margin: 0 2px;
                background-color: #dcdfe6;
                height: 18px;
            }
        }
        
        .word-count {
            font-size: 11px;
            color: #909399;
            margin: 0 4px;
            white-space: nowrap;
        }
        
        // 工具栏按钮样式
        .toolbar-icon-btn {
            padding: 3px 5px !important;
            min-width: auto !important;
            height: 26px !important;
            font-size: 12px !important;
            border-radius: 3px !important;
            
            .font-awesome-icon {
                font-size: 12px !important;
            }
        }
        
        .toolbar-select {
            height: 26px !important;
            font-size: 12px !important;
            
            :deep(.el-input__wrapper) {
                padding: 0 6px !important;
                height: 26px !important;
                border-radius: 3px !important;
            }
            
            :deep(.el-input__inner) {
                font-size: 12px !important;
            }
        }
        
        // 表格选择器样式
        .table-popover {
            padding: 0 !important;
        }
        
        .table-popover-wrapper {
            position: relative;
            display: inline-block;
            
            .table-popover-content {
                position: absolute;
                top: 100%;
                left: 0;
                margin-top: 4px;
                background: #fff;
                border: 1px solid #e4e7ed;
                border-radius: 4px;
                box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
                z-index: 2000;
                padding: 8px;
                white-space: nowrap;
                
                .table-grid-selector {
                    .table-grid {
                        display: inline-block;
                        border: 1px solid #e4e7ed;
                        border-radius: 3px;
                        overflow: hidden;
                        background: #fff;
                        
                        .table-row {
                            display: flex;
                            
                            .table-cell {
                                width: 20px;
                                height: 20px;
                                border: 1px solid #e4e7ed;
                                background: #fff;
                                cursor: pointer;
                                transition: all 0.15s;
                                box-sizing: border-box;
                                
                                &.active {
                                    background: #409eff;
                                    border-color: #409eff;
                                }
                                
                                &:hover {
                                    background: #ecf5ff;
                                    border-color: #409eff;
                                }
                            }
                        }
                    }
                    
                    .table-size-info {
                        margin-top: 8px;
                        text-align: center;
                        font-size: 12px;
                        color: #606266;
                        padding: 4px;
                        background: #f5f7fa;
                        border-radius: 3px;
                    }
                }
            }
        }
        
        .toolbar-icon-btn.active {
            background: #ecf5ff;
            border-color: #409eff;
            color: #409eff;
        }
    }
    
    .editor-body {
        flex: 1;
        display: flex;
        overflow: hidden;
        
        .toc-panel {
            border-right: 1px solid #dcdfe6;
            background: #fafafa;
            
            .toc-header {
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 8px 12px;
                border-bottom: 1px solid #e4e7ed;
                font-size: 14px;
                font-weight: 500;
            }
            
            .toc-content {
                padding: 8px;
                
                .toc-item {
                    display: block;
                    padding: 4px 8px;
                    color: #606266;
                    text-decoration: none;
                    font-size: 13px;
                    border-radius: 4px;
                    cursor: pointer;
                    white-space: nowrap;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    
                    &:hover {
                        background: #f0f2f5;
                        color: #409EFF;
                    }
                    
                    @for $i from 1 through 6 {
                        &.toc-level-#{$i} {
                            padding-left: #{($i - 1) * 16 + 8}px;
                        }
                    }
                }
            }
        }
        
        .editor-pane, .preview-pane {
            flex: 1;
            overflow: hidden;
            
            @media (min-width: 1200px) {
                &:nth-child(2):last-child {
                    border-left: 1px solid #dcdfe6;
                }
            }
        }
        
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
        
        .preview-content {
            padding: 16px;
            height: 100%;
            overflow-y: auto;
            font-size: 14px;
            line-height: 1.6;
            
            :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
                margin-top: 24px;
                margin-bottom: 16px;
                font-weight: 600;
                line-height: 1.25;
                scroll-margin-top: 20px;
            }
            
            :deep(h1) {
                font-size: 24px;
                border-bottom: 1px solid #eaecef;
                padding-bottom: 0.3em;
            }
            
            :deep(h2) {
                font-size: 20px;
                border-bottom: 1px solid #eaecef;
                padding-bottom: 0.3em;
            }
            
            :deep(code) {
                padding: 0.2em 0.4em;
                background-color: rgba(27, 31, 35, 0.05);
                border-radius: 3px;
                font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
            }
            
            :deep(pre) {
                padding: 16px;
                overflow: auto;
                background-color: #f6f8fa;
                border-radius: 3px;
                
                code {
                    padding: 0;
                    background: none;
                }
            }
            
            :deep(table) {
                border-collapse: collapse;
                width: 100%;
                margin: 16px 0;
                
                th, td {
                    padding: 6px 13px;
                    border: 1px solid #dfe2e5;
                }
                
                th {
                    background-color: #f6f8fa;
                    font-weight: 600;
                }
            }
        }
    }
    
    .editor-statusbar {
        display: flex;
        align-items: center;
        padding: 4px 12px;
        background: #f5f7fa;
        border-top: 1px solid #dcdfe6;
        font-size: 12px;
        color: #909399;
    }
}
</style>
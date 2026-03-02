<template>
    <div class="folder-tree-container">
        <!-- 工具栏 -->
        <div class="folder-toolbar">
            <div class="toolbar-title">
                <font-awesome-icon icon="fa-solid fa-folder" class="title-icon" />
                <span>文档目录</span>
            </div>
            <div class="toolbar-actions">
                <el-button type="primary" size="small" @click="handleCreateRootFolder" class="create-root-btn" circle title="新建根目录">
                    <font-awesome-icon icon="fa-solid fa-folder-plus" />
                </el-button>
                <el-button size="small" @click="refreshTree" class="refresh-btn" circle title="刷新">
                    <font-awesome-icon icon="fa-solid fa-arrows-rotate" />
                </el-button>
            </div>
        </div>

        <!-- 文件夹树 -->
        <div class="folder-tree-wrapper">
            <el-empty v-if="!loading && folderTree.length === 0" description="暂无文件夹" :image-size="60">
                <el-button type="primary" size="small" @click="handleCreateRootFolder">新建根目录</el-button>
            </el-empty>
            <el-tree
                v-else
                ref="treeRef"
                :data="folderTree"
                :props="treeProps"
                :expand-on-click-node="false"
                :highlight-current="true"
                node-key="id"
                default-expand-all
                v-loading="loading"
                :virtual-tree="true"
                @node-click="handleNodeClick"
                @node-contextmenu="handleNodeContextmenu"
            >
                <template #default="{ node, data }">
                    <div class="folder-tree-node">
                        <font-awesome-icon icon="fa-solid fa-folder" class="folder-icon" />
                        <span class="folder-name">{{ node.label }}</span>
                        <span v-if="data.docCount !== undefined" class="folder-doc-count">
                            {{ data.docCount }}
                        </span>
                    </div>
                </template>
            </el-tree>
        </div>

        <!-- 右键菜单 -->
        <div
            v-if="contextMenuVisible"
            class="context-menu"
            :style="{ top: contextMenuY + 'px', left: contextMenuX + 'px' }"
            @click.stop
        >
            <el-menu class="context-menu-list">
                <el-menu-item @click="handleCreateSubFolder">
                    <font-awesome-icon icon="fa-solid fa-folder-plus" />
                    <span>新建子文件夹</span>
                </el-menu-item>
                <el-menu-item @click="handleCreateSiblingFolder">
                    <font-awesome-icon icon="fa-solid fa-folder-plus" />
                    <span>新建同级文件夹</span>
                </el-menu-item>
                <el-menu-item @click="handleRenameFolder">
                    <font-awesome-icon icon="fa-solid fa-pen-to-square" />
                    <span>重命名</span>
                </el-menu-item>
                <el-menu-item @click="handleMoveFolder">
                    <font-awesome-icon icon="fa-solid fa-arrow-right" />
                    <span>移动</span>
                </el-menu-item>
                <el-menu-item @click="handleCopyFolder">
                    <font-awesome-icon icon="fa-solid fa-copy" />
                    <span>复制</span>
                </el-menu-item>
                <div class="context-menu-divider" />
                <el-menu-item @click="handleDeleteFolder" class="danger-item">
                    <font-awesome-icon icon="fa-solid fa-trash" />
                    <span>删除</span>
                </el-menu-item>
            </el-menu>
        </div>

        <!-- 新建/编辑文件夹对话框 -->
        <el-dialog
            v-model="dialogVisible"
            :title="dialogTitle"
            width="600px"
            :close-on-click-modal="false"
            @close="handleDialogClose"
        >
            <el-form
                ref="formRef"
                :model="formData"
                :rules="formRules"
                label-width="100px"
                label-position="left"
            >
                <el-form-item label="文件夹名称" prop="name">
                    <el-input
                        v-model="formData.name"
                        placeholder="请输入文件夹名称"
                        :maxlength="255"
                        show-word-limit
                        clearable
                    />
                </el-form-item>
                <el-form-item label="描述" prop="description">
                    <el-input
                        v-model="formData.description"
                        type="textarea"
                        :rows="3"
                        :placeholder="dialogMode === 'create' ? '请输入文件夹描述（可选）' : '留空则保留原描述'"
                        :maxlength="512"
                        show-word-limit
                        clearable
                    />
                </el-form-item>
                <el-form-item label="排序" prop="sortOrder">
                    <el-input-number
                        v-model="formData.sortOrder"
                        :min="0"
                        :max="9999"
                        controls-position="right"
                        style="width: 100%"
                    />
                </el-form-item>
                <el-form-item label="启用状态" prop="isEnabled">
                    <el-switch v-model="formData.isEnabled" />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="dialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="submitting" @click="handleSubmit">
                    确定
                </el-button>
            </template>
        </el-dialog>

        <!-- 移动文件夹对话框 -->
        <el-dialog
            v-model="moveDialogVisible"
            title="移动文件夹"
            width="600px"
            :close-on-click-modal="false"
            @close="handleMoveDialogClose"
        >
            <el-form label-width="100px" label-position="left">
                <el-form-item label="目标位置">
                    <el-tree-select
                        v-model="moveTargetId"
                        :data="moveTreeData"
                        :props="treeProps"
                        placeholder="请选择目标文件夹"
                        check-strictly
                        :render-after-expand="false"
                        style="width: 100%"
                    />
                </el-form-item>
            </el-form>
            <template #footer>
                <el-button @click="moveDialogVisible = false">取消</el-button>
                <el-button type="primary" :loading="moving" @click="handleMoveSubmit">
                    确定
                </el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import {
    copyFolder,
    createFolder,
    deleteFolder,
    getFolderTree,
    moveFolder,
    updateFolder
} from '@/api/document-folder'
import type { CreateFolderParams, DocumentFolderTree, UpdateFolderParams } from '@/types/document-folder'

import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'

// ==================== 状态定义 ====================

const loading = ref(false) // 树加载状态
const submitting = ref(false) // 提交状态
const moving = ref(false) // 移动状态
const folderTree = ref<DocumentFolderTree[]>([]) // 文件夹树数据
const treeRef = ref() // 树组件引用，用于调用树的方法（如 getCheckedNodes 等）

// 右键菜单状态
const contextMenuVisible = ref(false) // 右键菜单显示状态
const contextMenuX = ref(0) // 右键菜单 X 坐标
const contextMenuY = ref(0) // 右键菜单 Y 坐标
const currentFolder = ref<DocumentFolderTree | null>(null) // 当前右键选中的文件夹

// 对话框状态
const dialogVisible = ref(false) // 对话框显示状态
const dialogTitle = computed(() => dialogMode.value === 'create' ? '新建文件夹' : '重命名文件夹') // 对话框标题
const dialogMode = ref<'create' | 'edit'>('create') // 对话框模式：create=新建，edit=编辑
const createMode = ref<'child' | 'sibling'>('child') // 创建模式：child=子文件夹，sibling=同级文件夹
const editingFolderId = ref<string | null>(null) // 保存正在编辑/移动的文件夹 ID

// 表单数据
const formData = reactive<CreateFolderParams & UpdateFolderParams>({
    name: '', // 文件夹名称
    description: '', // 文件夹描述
    sortOrder: 0, // 排序值
    isEnabled: true // 是否启用
})

// 表单验证规则
const formRules = {
    name: [
        { required: true, message: '请输入文件夹名称', trigger: 'blur' },
        { max: 255, message: '文件夹名称不能超过 255 个字符', trigger: 'blur' }
    ]
}

// 树配置
const treeProps = {
    children: 'children', // 子节点字段
    label: 'name', // 节点标签字段
    value: 'id' // 节点 ID 字段
}

// 移动对话框
const moveDialogVisible = ref(false) // 移动对话框显示状态
const moveTargetId = ref('') // 目标文件夹 ID
const moveTreeData = ref<DocumentFolderTree[]>([]) // 移动对话框中的树数据

// ==================== 生命周期 ====================

onMounted(() => {
    loadFolderTree()
    document.addEventListener('click', closeContextMenu)
})

onUnmounted(() => {
    document.removeEventListener('click', closeContextMenu)
})

// ==================== 树加载 ====================

/**
 * 加载文件夹树
 */
const loadFolderTree = async () => {
    loading.value = true
    try {
        const res = await getFolderTree()
        folderTree.value = res
    } catch (error: any) {
        ElMessage.error('加载文件夹树失败：' + (error.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

/**
 * 刷新树
 */
const refreshTree = () => {
    loadFolderTree()
}

// ==================== 节点事件 ====================

/**
 * 节点点击
 */
const handleNodeClick = (data: DocumentFolderTree) => {
    console.log('节点点击:', data)
    // 可以在这里触发自定义事件，供父组件使用
    emit('node-click', data)
}

/**
 * 节点右键
 */
const handleNodeContextmenu = (event: MouseEvent, data: DocumentFolderTree) => {
    event.preventDefault()
    event.stopPropagation()
    
    currentFolder.value = data
    
    // 计算菜单位置，确保不会超出屏幕
    const menuWidth = 160 // 菜单宽度
    const menuPadding = 8 // 内边距
    const totalWidth = menuWidth + menuPadding * 2
    const totalHeight = 280 // 估算菜单总高度
    
    // 获取窗口尺寸
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    
    // 计算 X 坐标（确保不超出右边界）
    const x = event.clientX + totalWidth > windowWidth 
        ? event.clientX - totalWidth - 10 // 在鼠标左侧显示
        : event.clientX + 10 // 在鼠标右侧显示
    
    // 计算 Y 坐标（确保不超出下边界）
    const y = event.clientY + totalHeight > windowHeight
        ? windowHeight - totalHeight - 10 // 靠上对齐
        : event.clientY // 正常位置
    
    contextMenuX.value = Math.max(10, x) // 至少距离左边界 10px
    contextMenuY.value = Math.max(10, y) // 至少距离上边界 10px
    contextMenuVisible.value = true
}

/**
 * 关闭右键菜单
 */
const closeContextMenu = () => {
    contextMenuVisible.value = false
    currentFolder.value = null
}

// ==================== 右键菜单操作 ====================

/**
 * 新建根目录
 */
const handleCreateRootFolder = () => {
    dialogMode.value = 'create'
    formData.name = ''
    formData.description = ''
    formData.sortOrder = 0
    formData.isEnabled = true
    dialogVisible.value = true
}

/**
 * 新建子文件夹
 */
const handleCreateSubFolder = () => {
    if (!currentFolder.value) return
    
    dialogMode.value = 'create'
    createMode.value = 'child'
    formData.name = ''
    formData.description = ''
    formData.sortOrder = 0
    formData.isEnabled = true
    dialogVisible.value = true
    
    closeContextMenu()
}

/**
 * 新建同级文件夹
 */
const handleCreateSiblingFolder = () => {
    if (!currentFolder.value) return
    
    dialogMode.value = 'create'
    createMode.value = 'sibling'
    formData.name = ''
    formData.description = ''
    formData.sortOrder = 0
    formData.isEnabled = true
    dialogVisible.value = true
    
    closeContextMenu()
}

/**
 * 重命名文件夹
 */
const handleRenameFolder = () => {
    if (!currentFolder.value) return
    
    dialogMode.value = 'edit'
    editingFolderId.value = currentFolder.value.id // 保存正在编辑的文件夹 ID
    // 加载完整的文件夹信息，包括描述和排序等
    formData.name = currentFolder.value.name
    formData.description = currentFolder.value.description || ''
    formData.sortOrder = currentFolder.value.sortOrder ?? 0
    formData.isEnabled = currentFolder.value.isEnabled ?? true
    dialogVisible.value = true
    
    closeContextMenu()
}

/**
 * 移动文件夹
 */
const handleMoveFolder = () => {
    if (!currentFolder.value) return
    
    // 保存正在移动的文件夹 ID（在 closeContextMenu 之前）
    editingFolderId.value = currentFolder.value.id
    
    /**
     * 过滤掉当前文件夹及其子文件夹
     * 不能移动到自己的孩子下面，也不能移动到自己下面
     * @param nodes 树节点列表
     * @param targetId 要移动的文件夹 ID
     * @returns 过滤后的树节点列表
     */
    const filterSelfAndChildren = (nodes: DocumentFolderTree[], targetId: string): DocumentFolderTree[] => {
        const result: DocumentFolderTree[] = []
        for (const node of nodes) {
            if (node.id !== targetId) {
                // 检查是否是目标文件夹的子节点
                // 通过路径判断：如果节点路径以目标文件夹路径开头，则是其子节点
                const isChild = node.path.startsWith(currentFolder.value!.path + '/')
                if (!isChild) {
                    const newNode = { ...node }
                    if (node.children && node.children.length > 0) {
                        newNode.children = filterSelfAndChildren(node.children, targetId)
                    }
                    result.push(newNode)
                }
            }
        }
        return result
    }
    
    moveTreeData.value = filterSelfAndChildren(folderTree.value, currentFolder.value.id)
    moveTargetId.value = ''
    moveDialogVisible.value = true
    
    closeContextMenu()
}

/**
 * 复制文件夹
 */
const handleCopyFolder = async () => {
    if (!currentFolder.value) return
    
    try {
        await ElMessageBox.prompt('请输入新文件夹名称', '复制文件夹', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            inputValue: currentFolder.value.name + ' - 副本',
            inputPattern: /.+/,
            inputErrorMessage: '名称不能为空'
        }).then(async ({ value }) => {
            if (!currentFolder.value) return
            
            moving.value = true
            try {
                // 复制文件夹到当前父节点下
                // 从路径解析父节点 ID
                const pathParts = currentFolder.value.path.split('/').filter(p => p)
                let newParentId = ''
                if (pathParts.length > 1) {
                    // 取倒数第二个作为父 ID（最后一个是当前节点）
                    newParentId = pathParts[pathParts.length - 2]
                }
                // 根节点的 newParentId 保持空字符串
                
                await copyFolder(currentFolder.value.id, {
                    newParentId,
                    newName: value
                })
                ElMessage.success('复制成功')
                loadFolderTree()
            } catch (error: any) {
                const errorMsg = error?.response?.data?.message || error?.message || '未知错误'
                ElMessage.error(`复制失败：${errorMsg}`)
            } finally {
                moving.value = false
            }
        })
    } catch (error: any) {
        if (error !== 'cancel') {
            const errorMsg = error?.message || '未知错误'
            ElMessage.error(`复制失败：${errorMsg}`)
        }
    }
    
    closeContextMenu()
}

/**
 * 删除文件夹
 */
const handleDeleteFolder = () => {
    if (!currentFolder.value) return
    
    const hasChildren = currentFolder.value.hasChildren || (currentFolder.value.children && currentFolder.value.children.length > 0)
    const docCount = currentFolder.value.docCount || 0
    const folderId = currentFolder.value.id // 保存文件夹 ID
    const folderName = currentFolder.value.name // 保存文件夹名称
    
    let message = `确定要删除文件夹"${folderName}"吗？`
    if (hasChildren || docCount > 0) {
        message += '\n\n'
        if (hasChildren) {
            message += '⚠️ 该文件夹包含子文件夹，\n'
        }
        if (docCount > 0) {
            message += `⚠️ 该文件夹包含 ${docCount} 个文档，\n`
        }
        message += '\n删除后所有数据将无法恢复！'
    }
    
    ElMessageBox.confirm(
        message,
        '删除确认',
        {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
            dangerouslyUseHTMLString: false
        }
    ).then(async () => {
        try {
            await deleteFolder(folderId)
            ElMessage.success('删除成功')
            loadFolderTree()
        } catch (error: any) {
            const errorMsg = error?.response?.data?.message || error?.message || '未知错误'
            ElMessage.error(`删除失败：${errorMsg}`)
        } finally {
            closeContextMenu()
        }
    }).catch(() => {
        // 取消删除
        closeContextMenu()
    })
}

// ==================== 表单提交 ====================

const formRef = ref()

/**
 * 提交表单
 */
const handleSubmit = async () => {
    if (!formRef.value) {
        console.error('表单引用为空')
        ElMessage.error('表单初始化失败，请刷新页面')
        return
    }
    
    try {
        // 执行表单验证
        await formRef.value.validate(async (valid: boolean, fields?: any) => {
            console.log('表单验证结果:', valid, fields)
            
            if (!valid) {
                console.warn('表单验证失败:', fields)
                ElMessage.warning('请检查表单填写是否正确')
                return
            }
            
            submitting.value = true
            try {
                if (dialogMode.value === 'create') {
                    // 创建文件夹
                    let parentId = ''
                    if (createMode.value === 'child' && currentFolder.value) {
                        // 新建子文件夹，使用当前文件夹 ID 作为父 ID
                        parentId = currentFolder.value.id
                    } else if (createMode.value === 'sibling' && currentFolder.value) {
                        // 新建同级文件夹，使用当前文件夹的父 ID
                        // 从路径解析父 ID
                        const pathParts = currentFolder.value.path.split('/').filter(p => p)
                        if (pathParts.length > 1) {
                            parentId = pathParts[pathParts.length - 2]
                        }
                        // 根节点的 parentId 保持空字符串
                    }
                    
                    console.log('创建文件夹，parentId:', parentId, 'formData:', formData)
                    await createFolder({
                        ...formData,
                        parentId
                    })
                    ElMessage.success('创建成功')
                } else {
                    // 更新文件夹
                    if (!editingFolderId.value) {
                        ElMessage.error('文件夹信息丢失，请重试')
                        return
                    }
                    console.log('更新文件夹，id:', editingFolderId.value, 'formData:', formData)
                    await updateFolder(editingFolderId.value, formData)
                    ElMessage.success('更新成功')
                }
                
                dialogVisible.value = false
                loadFolderTree()
            } catch (error: any) {
                console.error('提交失败:', error)
                const action = dialogMode.value === 'create' ? '创建' : '更新'
                const errorMsg = error?.response?.data?.message || error?.message || '未知错误'
                ElMessage.error(`${action}失败：${errorMsg}`)
            } finally {
                submitting.value = false
            }
        })
    } catch (error: any) {
        console.error('验证过程出错:', error)
        ElMessage.error('操作失败：' + (error.message || '未知错误'))
        submitting.value = false
    }
}

/**
 * 对话框关闭
 */
const handleDialogClose = () => {
    formRef.value?.resetFields()
    currentFolder.value = null
    editingFolderId.value = null
}

/**
 * 移动对话框关闭
 */
const handleMoveDialogClose = () => {
    moveTargetId.value = ''
    editingFolderId.value = null
}

// ==================== 移动提交 ====================

/**
 * 移动提交
 */
const handleMoveSubmit = async () => {
    if (!editingFolderId.value || !moveTargetId.value) {
        if (!editingFolderId.value) {
            ElMessage.error('文件夹信息丢失，请重试')
        } else {
            ElMessage.error('请选择目标位置')
        }
        return
    }
    
    moving.value = true
    try {
        console.log('移动文件夹，id:', editingFolderId.value, 'newParentId:', moveTargetId.value)
        await moveFolder(editingFolderId.value, {
            newParentId: moveTargetId.value
        })
        ElMessage.success('移动成功')
        moveDialogVisible.value = false
        loadFolderTree()
    } catch (error: any) {
        console.error('移动失败:', error)
        const errorMsg = error?.response?.data?.message || error?.message || '未知错误'
        ElMessage.error(`移动失败：${errorMsg}`)
    } finally {
        moving.value = false
    }
}

// ==================== 事件定义 ====================

const emit = defineEmits<{
    (e: 'node-click', data: DocumentFolderTree): void
}>()
</script>

<style scoped lang="scss">
.folder-tree-container {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    background: #fff;
    overflow: hidden;
    
    .folder-toolbar {
        padding: 20px;
        border-bottom: 1px solid #e4e7ed;
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-shrink: 0;
        background: linear-gradient(to right, #f8faff, #fff);
        
        .toolbar-title {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 14px;
            font-weight: 600;
            color: #303133;
            
            .title-icon {
                color: #409EFF;
                font-size: 18px;
            }
        }
        
        .toolbar-actions {
            display: flex;
            gap: 8px;
            align-items: center;
            
            .el-button {
                font-size: 14px;
                transition: all 0.2s ease;
                width: 36px;
                height: 36px;
                display: flex;
                align-items: center;
                justify-content: center;
                
                &:hover {
                    transform: translateY(-1px);
                    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
                }
            }
            
            .create-root-btn {
                background: linear-gradient(135deg, #409EFF, #69b1ff);
                border: none;
                
                &:hover {
                    background: linear-gradient(135deg, #337ecc, #409EFF);
                }
            }
            
            .refresh-btn {
                background: #f5f7fa;
                border: 1px solid #dcdfe6;
                color: #606266;
                
                &:hover {
                    background: #ecf5ff;
                    border-color: #c6e2ff;
                    color: #409EFF;
                }
                
                &:active {
                    animation: rotate 1s ease-in-out;
                }
            }
        }
        
        @keyframes rotate {
            from {
                transform: rotate(0deg);
            }
            to {
                transform: rotate(360deg);
            }
        }
    }
    
    .folder-tree-wrapper {
        flex: 1;
        overflow: auto; // 允许滚动
        padding: 0 20px 20px 20px; // 上右下左：上边 0（与工具栏间距），右边 20px，下边 20px，左边 20px
        
        :deep(.el-tree) {
            background: transparent;
            max-height: 100%;
            overflow: visible; // 树本身不需要滚动，由 wrapper 处理
            
            // 隐藏滚动条但保持滚动功能（webkit 浏览器）
            &::-webkit-scrollbar {
                width: 4px;
            }
            
            &::-webkit-scrollbar-track {
                background: transparent;
            }
            
            &::-webkit-scrollbar-thumb {
                background: #ddd;
                border-radius: 2px;
                
                &:hover {
                    background: #ccc;
                }
            }
            
            // 树节点内容区域
            .el-tree-node__content {
                height: 36px;
                border-radius: 6px;
                margin-bottom: 2px;
                transition: all 0.2s ease;
                
                &:hover {
                    background-color: #f5f7fa;
                }
                
                &.is-current {
                    background-color: #ecf5ff;
                    
                    .folder-name {
                        color: #409EFF;
                        font-weight: 600;
                    }
                    
                    .folder-icon {
                        color: #409EFF;
                    }
                }
            }
            
            // 树节点内边距
            .el-tree-node__content {
                padding: 0 8px;
            }
            
            // 展开/折叠图标
            .el-tree-node__expand-icon {
                font-size: 14px;
                color: #909399;
                transition: transform 0.2s ease;
                
                &.is-leaf {
                    color: transparent;
                }
            }
        }
        
        .folder-tree-node {
            flex: 1;
            display: flex;
            align-items: center;
            gap: 8px;
            min-width: 0;
            
            .folder-icon {
                color: #606266;
                font-size: 16px;
                flex-shrink: 0;
                transition: color 0.2s ease;
            }
            
            .folder-name {
                flex: 1;
                font-size: 13px;
                color: #606266;
                white-space: nowrap;
                overflow: hidden;
                text-overflow: ellipsis;
                transition: color 0.2s ease;
            }
            
            .folder-doc-count {
                font-size: 11px;
                color: #909399;
                background: #f5f7fa;
                padding: 2px 6px;
                border-radius: 8px;
                flex-shrink: 0;
                font-weight: 500;
                transition: all 0.2s ease;
            }
        }
    }
}

// 右键菜单
.context-menu {
    position: fixed;
    z-index: 3000;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
    padding: 4px 0;
    animation: contextMenuFadeIn 0.15s ease-out;
    
    @keyframes contextMenuFadeIn {
        from {
            opacity: 0;
            transform: translateY(-4px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
    
    .context-menu-list {
        width: 160px;
        border: none;
        background: transparent;
        
        .el-menu-item {
            height: 36px;
            line-height: 36px;
            font-size: 13px;
            margin: 0 4px;
            border-radius: 4px;
            
            &:hover {
                background-color: #f5f7fa;
            }
            
            &.danger-item {
                color: #f56c6c;
                
                &:hover {
                    background-color: #fef0f0;
                }
            }
        }
        
        .context-menu-divider {
            height: 1px;
            background-color: #e4e7ed;
            margin: 4px 8px;
        }
    }
}
</style>

<template>
    <div class="document-create">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
            <el-skeleton :rows="10" animated />
        </div>

        <template v-else>
            <!-- 页面头部 -->
            <div class="page-header">
                <h1>{{ isEditMode ? '编辑文档' : '新建文档' }}</h1>
                <div class="header-actions">
                    <el-button @click="handleCancel">
                        <font-awesome-icon icon="fa-solid fa-arrow-left" />
                        <span>返回列表</span>
                    </el-button>
                    <el-button
                        type="primary"
                        @click="handleSubmit"
                        :loading="submitting"
                    >
                        <font-awesome-icon icon="fa-solid fa-check" />
                        <span>{{ isEditMode ? '保存修改' : '创建文档' }}</span>
                    </el-button>
                </div>
            </div>

            <!-- 使用标签页分离基本信息和内容编辑 -->
            <el-tabs v-model="activeTab" class="document-tabs" type="border-card">
                <!-- 基本信息标签页 -->
                <el-tab-pane label="基本信息" name="basic">
                    <el-form
                        ref="documentFormRef"
                        :model="documentForm"
                        :rules="formRules"
                        label-width="100px"
                        label-position="right"
                        class="basic-form"
                    >
                        <el-form-item label="文档标题" prop="title">
                            <el-input
                                v-model="documentForm.title"
                                placeholder="请输入文档标题"
                                clearable
                                maxlength="255"
                                show-word-limit
                            />
                        </el-form-item>

                        <el-form-item label="文档分类" prop="category">
                            <el-select
                                v-model="documentForm.category"
                                placeholder="请选择分类"
                                style="width: 100%"
                                filterable
                                allow-create
                                default-first-option
                            >
                                <el-option
                                    v-for="cat in categories"
                                    :key="cat.id"
                                    :label="cat.name"
                                    :value="cat.id"
                                />
                            </el-select>
                        </el-form-item>

                        <el-form-item label="文档路径" prop="path">
                        <el-tree-select
                            v-model="documentForm.path"
                            :data="folderTree"
                            :props="treeProps"
                            placeholder="请选择文档路径"
                            check-strictly
                            :render-after-expand="false"
                            clearable
                            filterable
                            style="width: 100%"
                            value-key="path"
                        >
                            <template #default="{ data }">
                                <span class="custom-tree-node">
                                    <font-awesome-icon icon="fa-solid fa-folder" />
                                    <span style="margin-left: 8px">{{ data.name }}</span>
                                </span>
                            </template>
                        </el-tree-select>
                    </el-form-item>

                        <el-form-item label="文档描述" prop="description">
                            <el-input
                                v-model="documentForm.description"
                                type="textarea"
                                :rows="3"
                                placeholder="请输入文档描述（可选）"
                                maxlength="1000"
                                show-word-limit
                            />
                        </el-form-item>

                        <el-form-item label="标签" prop="tags">
                            <el-select
                                v-model="documentForm.tags"
                                multiple
                                filterable
                                allow-create
                                default-first-option
                                placeholder="请输入标签，按回车键添加"
                                style="width: 100%"
                            >
                                <el-option
                                    v-for="tag in commonTags"
                                    :key="tag"
                                    :label="tag"
                                    :value="tag"
                                />
                            </el-select>
                        </el-form-item>

                        <el-form-item label="发布状态" prop="isPublished">
                            <el-radio-group v-model="documentForm.isPublished">
                                <el-radio :value="true">立即发布</el-radio>
                                <el-radio :value="false">草稿</el-radio>
                            </el-radio-group>
                        </el-form-item>
                    </el-form>
                </el-tab-pane>

                <!-- Markdown 编辑器标签页 -->
                <el-tab-pane label="内容编辑" name="editor">
                    <MarkdownEditor
                        ref="editorRef"
                        v-model="documentForm.content"
                        @save="handleSubmit"
                        @image-upload="handleImageUpload"
                    />
                </el-tab-pane>
            </el-tabs>
        </template>
    </div>
</template>

<script setup lang="ts">
import { createDocument, getDocumentById, getDocumentCategories, updateDocument } from '@/api/document'
import MarkdownEditor from '@/components/MarkdownEditor/index.vue'
import { getFolderTree } from '@/api/document-folder'
import type { DocumentCategory, DocumentDetail } from '@/types/document'
import type { DocumentFolderTree } from '@/types/document-folder'

import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// 路由
const router = useRouter()
const route = useRoute()

// 表单引用
const documentFormRef = ref<any>(null)
const editorRef = ref<InstanceType<typeof MarkdownEditor> | null>(null)

// 标签页
const activeTab = ref('basic')

// 状态
const submitting = ref(false)
const loading = ref(false)
const categories = ref<DocumentCategory[]>([])
const folderTree = ref<DocumentFolderTree[]>([])
const commonTags = ref<string[]>(['入门', '教程', 'API', '指南', '示例', '最佳实践'])

// 树形配置
const treeProps = {
    value: 'path',
    label: 'name',
    children: 'children'
}

// 判断是否是编辑模式
const isEditMode = computed(() => !!route.params.id)

// 表单数据
const documentForm = reactive({
    title: '',
    category: '',
    path: '',
    description: '',
    content: '',
    tags: [] as string[],
    isPublished: true
})

// 表单验证规则
const formRules = reactive({
    title: [
        { required: true, message: '请输入文档标题', trigger: 'blur' },
        { min: 2, max: 255, message: '标题长度在 2 到 255 个字符之间', trigger: 'blur' }
    ],
    category: [],
    path: [
        { required: true, message: '请输入文档路径', trigger: 'blur' },
        {
            pattern: /^\/[a-zA-Z0-9/_-]*$/,
            message: '路径必须以/开头，只能包含字母、数字、下划线和连字符',
            trigger: 'blur'
        }
    ]
})

// 生命周期
onMounted(() => {
    console.log('进入onMounted钩子')
    console.log('路由参数:', route.params)
    console.log('是否是编辑模式:', isEditMode.value)
    loadCategories()
    loadFolderTree()
    
    // 如果是编辑模式，加载文档数据
    if (isEditMode.value) {
        console.log('调用loadDocumentForEdit函数')
        loadDocumentForEdit()
    } else {
        // 如果是新建模式，检查是否有路由状态传递的文件夹信息
        const folderPath = route.query.folderPath as string
        if (folderPath) {
            documentForm.path = folderPath
        }
    }
})

// 加载分类列表
const loadCategories = async () => {
    try {
        const res: any = await getDocumentCategories()
        categories.value = res.data || res || []
    } catch (error: any) {
        console.error('加载分类失败:', error)
    }
}

// 加载文件夹树
const loadFolderTree = async () => {
    try {
        const res: any = await getFolderTree()
        folderTree.value = res.data || res || []
    } catch (error: any) {
        console.error('加载文件夹树失败:', error)
    }
}

// 加载文档用于编辑
const loadDocumentForEdit = async () => {
    console.log('进入loadDocumentForEdit函数')
    console.log('路由参数:', route.params)
    loading.value = true
    try {
        const docId = route.params.id as string
        console.log('文档ID:', docId)
        const doc: DocumentDetail = await getDocumentById(docId)
        console.log('获取到的文档数据:', doc)
        
        // 填充表单数据
        documentForm.title = doc.title
        documentForm.category = doc.category || ''
        documentForm.path = doc.path || ''
        documentForm.description = doc.description || ''
        documentForm.content = doc.content || ''
        
        // 处理标签：如果是 JSON 字符串，解析为数组
        if (doc.tags) {
            if (typeof doc.tags === 'string') {
                try {
                    documentForm.tags = JSON.parse(doc.tags)
                } catch {
                    // 如果解析失败，尝试按逗号分割
                    documentForm.tags = doc.tags.split(',').filter(t => t.trim())
                }
            } else if (Array.isArray(doc.tags)) {
                documentForm.tags = doc.tags
            }
        }
        
        documentForm.isPublished = true // 默认 true
        console.log('表单数据填充完成:', documentForm)
    } catch (error: any) {
        console.error('加载文档失败:', error)
        ElMessage.error('加载文档失败：' + (error.message || '未知错误'))
        router.push('/documents/list')
    } finally {
        loading.value = false
    }
}

// 取消
const handleCancel = async () => {
    // 检查表单是否有内容
    const hasUnsavedChanges = 
        documentForm.title || 
        documentForm.description || 
        documentForm.content || 
        (documentForm.tags && documentForm.tags.length > 0)

    if (hasUnsavedChanges) {
        try {
            await ElMessageBox.confirm(
                isEditMode.value 
                    ? '您有未保存的修改，确定要离开吗？' 
                    : '您有未保存的内容，确定要离开吗？',
                '确认离开',
                {
                    confirmButtonText: '确定离开',
                    cancelButtonText: '继续编辑',
                    type: 'warning'
                }
            )
            router.push('/documents/list')
        } catch (error) {
            // 用户选择继续编辑，不做任何操作
        }
    } else {
        router.push('/documents/list')
    }
}

// 保存处理
const handleSave = () => {
    handleSubmit()
}

// 图片上传处理
const handleImageUpload = async (file: File) => {
    // TODO: 实现图片上传到服务器的逻辑
    // 这里仅做本地预览
    ElMessage.info('图片上传功能待实现，当前仅支持本地预览')
}

// 表单提交
const handleSubmit = async () => {
    try {
        // 处理路径值：如果是对象，提取 path 属性
        if (documentForm.path && typeof documentForm.path === 'object') {
            documentForm.path = documentForm.path.path || ''
        }
        
        await documentFormRef.value.validate()
        
        // 验证内容
        if (!documentForm.content || documentForm.content.trim().length < 10) {
            ElMessage.warning('文档内容至少 10 个字符')
            activeTab.value = 'editor'
            return
        }
        
        submitting.value = true

        // 准备提交数据 - 标签转换为 JSON 字符串（与后端兼容）
        let tagsData: string = '[]'
        if (documentForm.tags.length > 0) {
            tagsData = JSON.stringify(documentForm.tags)
        }

        const submitData = {
            title: documentForm.title,
            category: documentForm.category,
            path: documentForm.path,
            description: documentForm.description || '',
            content: documentForm.content,
            tags: tagsData,
            isPublished: documentForm.isPublished,
            size: documentForm.content.length // 添加文档大小
        }

        console.log('提交数据:', submitData)

        // 根据模式调用不同的 API
        if (isEditMode.value) {
            const docId = route.params.id as string
            await updateDocument(docId, submitData)
            ElMessage.success('修改成功')
        } else {
            await createDocument(submitData)
            ElMessage.success('创建成功')
        }

        router.push('/documents/list')
    } catch (error: any) {
        console.error('表单提交失败:', error)
        if (error.response?.data?.message) {
            ElMessage.error('提交失败：' + error.response.data.message)
        } else if (error.message) {
            ElMessage.error(
                isEditMode.value 
                    ? '修改失败：' + error.message 
                    : '创建失败：' + error.message
            )
        } else if (error.title || error.path) {
            // 验证错误
            const messages: string[] = []
            if (error.title) messages.push('标题：' + error.title[0].message)
            if (error.path) messages.push('路径：' + error.path[0].message)
            ElMessage.error('请完善以下信息：\n' + messages.join('\n'))
        } else {
            ElMessage.error('提交失败，请检查数据格式')
        }
    } finally {
        submitting.value = false
    }
}
</script>

<style scoped lang="scss">
.document-create {
    padding: 20px;
    min-height: 100vh;
    background: #f5f7fa;

    .page-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
        background: #fff;
        padding: 16px 20px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

        h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 600;
            color: #303133;
        }

        .header-actions {
            display: flex;
            gap: 12px;
        }
    }

    .document-tabs {
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        min-height: calc(100vh - 180px);

        :deep(.el-tabs__content) {
            padding: 0;
            overflow: auto;
        }

        :deep(.el-tab-pane) {
            padding: 20px;
        }
    }

    .basic-form {
        max-width: 800px;
        margin: 0 auto;

        .el-form-item {
            margin-bottom: 24px;
        }

        .custom-tree-node {
            display: flex;
            align-items: center;
            gap: 8px;
        }
    }
}
</style>

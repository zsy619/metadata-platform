<template>
    <div class="document-list-page">
        <el-container class="document-container">
            <!-- 左侧文件夹树 -->
            <el-aside width="300px" class="document-aside">
                <FolderTree @node-click="handleFolderClick" />
            </el-aside>
            
            <!-- 右侧文档列表 -->
            <el-main class="document-main">
                <div class="document-list-container">
                    <!-- 面包屑导航 -->
                    <el-breadcrumb separator="/" class="document-breadcrumb">
                        <el-breadcrumb-item>全部文档</el-breadcrumb-item>
                        <el-breadcrumb-item v-for="item in breadcrumbList" :key="item.path">
                            {{ item.name }}
                        </el-breadcrumb-item>
                    </el-breadcrumb>
                    
                    <!-- 搜索和筛选 -->
                    <div class="document-toolbar">
                        <div class="toolbar-left">
                            <el-input
                                v-model="searchKeyword"
                                placeholder="搜索文档标题或内容..."
                                clearable
                                @clear="handleSearch"
                            >
                                <template #prefix>
                                    <el-icon><Search /></el-icon>
                                </template>
                                <template #append>
                                    <el-button @click="handleSearch">搜索</el-button>
                                </template>
                            </el-input>
                            
                            <el-select
                                v-model="currentCategory"
                                placeholder="全部分类"
                                clearable
                                style="width: 150px; margin-left: 12px"
                                @change="handleCategoryChange"
                            >
                                <el-option
                                    v-for="cat in categories"
                                    :key="cat.id"
                                    :label="cat.name"
                                    :value="cat.id"
                                />
                            </el-select>
                        </div>
                        
                        <div class="toolbar-actions">
                            <el-button type="primary" @click="handleCreateDocument">
                                <el-icon><DocumentAdd /></el-icon>
                                新建文档
                            </el-button>
                        </div>
                    </div>
                    
                    <!-- 文档列表 -->
                    <el-table
                        v-loading="loading"
                        :data="documentList"
                        style="width: 100%"
                        @row-click="handleRowClick"
                    >
                        <el-table-column prop="title" label="标题" min-width="300">
                            <template #default="{ row }">
                                <div class="document-title">
                                    <el-icon class="doc-icon"><Document /></el-icon>
                                    <span>{{ row.title }}</span>
                                </div>
                            </template>
                        </el-table-column>
                        <el-table-column prop="category" label="分类" width="120">
                            <template #default="{ row }">
                                <el-tag size="small">{{ row.category || '未分类' }}</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="path" label="路径" width="250" show-overflow-tooltip />
                        <el-table-column prop="version" label="版本" width="80" align="center">
                            <template #default="{ row }">
                                v{{ row.version || 1 }}
                            </template>
                        </el-table-column>
                        <el-table-column prop="updatedAt" label="更新时间" width="160">
                            <template #default="{ row }">
                                {{ formatDate(row.updatedAt) }}
                            </template>
                        </el-table-column>
                        <el-table-column label="操作" width="200" fixed="right">
                            <template #default="{ row }">
                                <el-button link type="primary" size="small" @click.stop="handleViewDocument(row)">
                                    查看
                                </el-button>
                                <el-button link type="primary" size="small" @click.stop="handleEditDocument(row)">
                                    编辑
                                </el-button>
                                <el-button link type="danger" size="small" @click.stop="handleDeleteDocument(row)">
                                    删除
                                </el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                    
                    <!-- 分页 -->
                    <div class="document-pagination">
                        <el-pagination
                            v-model:current-page="pagination.page"
                            v-model:page-size="pagination.pageSize"
                            :page-sizes="[10, 20, 50, 100]"
                            :total="pagination.total"
                            layout="total, sizes, prev, pager, next, jumper"
                            @size-change="handleSizeChange"
                            @current-change="handlePageChange"
                        />
                    </div>
                </div>
            </el-main>
        </el-container>
    </div>
</template>

<script setup lang="ts">
import { deleteDocument, getDocumentCategories, getDocumentList } from '@/api/document'
import FolderTree from '@/components/document/FolderTree.vue'
import type { DocumentCategory, DocumentInfo } from '@/types/document'
import type { DocumentFolderTree } from '@/types/document-folder'
import { Document, DocumentAdd, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

// 路由
const router = useRouter()

// ==================== 状态定义 ====================

const loading = ref(false)
const loadingText = ref('加载中...')
const documentList = ref<DocumentInfo[]>([])
const categories = ref<DocumentCategory[]>([])
const searchKeyword = ref('')
const currentCategory = ref('')
const currentFolder = ref<DocumentFolderTree | null>(null)

// 面包屑列表
const breadcrumbList = ref<{ name: string; path: string }[]>([])

// 分页
const pagination = reactive({
    page: 1,
    pageSize: 20,
    total: 0
})

// ==================== 生命周期 ====================

onMounted(() => {
    loadCategories()
    loadDocumentList()
})

// ==================== 数据加载 ====================

/**
 * 加载分类列表
 */
const loadCategories = async () => {
    try {
        const res: any = await getDocumentCategories()
        // 根据实际 API 响应结构调整
        categories.value = res.data || res || []
    } catch (error: any) {
        ElMessage.error('加载分类失败：' + (error.message || '未知错误'))
    }
}

/**
 * 加载文档列表
 */
const loadDocumentList = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const params: any = {
            page: pagination.page,
            pageSize: pagination.pageSize
        }
        
        if (searchKeyword.value) {
            params.keyword = searchKeyword.value
        }
        
        if (currentCategory.value) {
            params.category = currentCategory.value
        }
        
        const res: any = await getDocumentList(params)
        // 根据实际 API 响应结构调整
        documentList.value = res.data || res.list || []
        pagination.total = res.total || 0
    } catch (error: any) {
        ElMessage.error('加载文档列表失败：' + (error.message || '未知错误'))
    } finally {
        loading.value = false
    }
}

// ==================== 事件处理 ====================

/**
 * 文件夹点击
 */
const handleFolderClick = (data: DocumentFolderTree) => {
    currentFolder.value = data
    console.log('文件夹点击:', data)
    
    // 更新面包屑
    updateBreadcrumb(data)
    
    // 刷新文档列表（可以按文件夹筛选）
    loadDocumentList()
}

/**
 * 更新面包屑
 */
const updateBreadcrumb = (folder: DocumentFolderTree) => {
    const pathParts: { name: string; path: string }[] = []
    
    // 简单实现：直接使用文件夹路径
    // 实际项目中可能需要更复杂的路径解析
    pathParts.push({
        name: folder.name,
        path: folder.path
    })
    
    breadcrumbList.value = pathParts
}

/**
 * 搜索
 */
const handleSearch = () => {
    pagination.page = 1
    loadDocumentList()
}

/**
 * 分类变化
 */
const handleCategoryChange = () => {
    pagination.page = 1
    loadDocumentList()
}

/**
 * 新建文档
 */
const handleCreateDocument = () => {
    // 如果有选中的文件夹，传递文件夹路径
    if (currentFolder.value) {
        router.push({
            name: 'DocumentCreate',
            query: { folderPath: currentFolder.value.path }
        })
    } else {
        router.push({ name: 'DocumentCreate' })
    }
}

/**
 * 表格行点击
 */
const handleRowClick = (row: DocumentInfo) => {
    console.log('行点击:', row)
}

/**
 * 查看文档
 */
const handleViewDocument = (row: DocumentInfo) => {
    ElMessage.info('查看文档：' + row.title)
    // TODO: 跳转到文档详情页
}

/**
 * 编辑文档
 */
const handleEditDocument = (row: DocumentInfo) => {
    router.push(`/documents/${row.id}/edit`)
}

/**
 * 删除文档
 */
const handleDeleteDocument = async (row: DocumentInfo) => {
    try {
        await ElMessageBox.confirm(
            `确定要删除文档"${row.title}"吗？`,
            '删除确认',
            {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }
        )
        
        await deleteDocument(row.id)
        ElMessage.success('删除成功')
        loadDocumentList()
    } catch (error: any) {
        if (error !== 'cancel') {
            ElMessage.error('删除失败：' + (error.message || '未知错误'))
        }
    }
}

/**
 * 分页变化
 */
const handleSizeChange = () => {
    loadDocumentList()
}

const handlePageChange = () => {
    loadDocumentList()
}

/**
 * 格式化日期
 */
const formatDate = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    })
}
</script>

<style scoped lang="scss">
.document-list-page {
    height: 100vh;
    overflow: hidden;
    background: #f5f7fa;
    
    .document-container {
        height: 100vh;
        width: 100%;
        
        .document-aside {
            border-right: 1px solid #e4e7ed;
            background: #fff;
            overflow: hidden;
            display: flex;
            flex-direction: column;
            width: 280px !important; // 固定宽度
        }
        
        .document-main {
            background: #f5f7fa;
            padding: 20px;
            overflow: hidden;
            display: flex;
            flex-direction: column;
            gap: 16px;
            
            .document-list-container {
                height: 100%;
                display: flex;
                flex-direction: column;
                
                .document-breadcrumb {
                    padding: 12px 16px;
                    background: #fff;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
                    font-size: 14px;
                    margin-bottom: 16px;
                    
                    :deep(.el-breadcrumb__item) {
                        font-weight: 500;
                        
                        &:last-child {
                            color: #409EFF;
                        }
                    }
                }
                
                .document-toolbar {
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                    padding: 16px 20px;
                    background: #fff;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
                    margin-bottom: 16px;
                    
                    .toolbar-left {
                        display: flex;
                        align-items: center;
                        flex: 1;
                        gap: 12px;
                        
                        .el-input {
                            max-width: 400px;
                        }
                    }
                    
                    .toolbar-actions {
                        display: flex;
                        gap: 12px;
                        align-items: center;
                    }
                }
                
                :deep(.el-table) {
                    flex: 1;
                    background: #fff;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
                    overflow: hidden;
                    
                    th {
                        background: #fafafa;
                        font-weight: 600;
                        color: #606266;
                        font-size: 14px;
                        border-bottom: 1px solid #ebeef5;
                    }
                    
                    td {
                        padding: 14px 0;
                        border-bottom: 1px solid #f2f6fc;
                    }
                    
                    .el-table__body tr:hover {
                        background: #f5f7fa;
                    }
                    
                    .document-title {
                        display: flex;
                        align-items: center;
                        gap: 8px;
                        font-weight: 500;
                        color: #303133;
                        
                        .doc-icon {
                            color: #409EFF;
                            font-size: 16px;
                        }
                    }
                }
                
                .document-pagination {
                    padding: 16px 20px;
                    display: flex;
                    justify-content: flex-end;
                    background: #fff;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
                    margin-top: 16px;
                    
                    :deep(.el-pagination) {
                        padding: 0;
                    }
                }
            }
        }
    }
}
</style>

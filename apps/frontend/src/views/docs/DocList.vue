<template>
  <div class="doc-list-container">
    <!-- 页面标题 -->
    <el-card class="header-card" shadow="hover">
      <div class="page-header">
        <div class="header-left">
          <el-icon class="header-icon"><Document /></el-icon>
          <h2 class="page-title">系统文档</h2>
        </div>
        <div class="header-right">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文档..."
            clearable
            prefix-icon="Search"
            style="width: 300px"
            @input="handleSearch"
          />
        </div>
      </div>
    </el-card>

    <!-- 分类标签 -->
    <el-card class="category-card" shadow="hover">
      <div class="category-tabs">
        <el-tabs v-model="activeCategory" type="card" @tab-change="handleCategoryChange">
          <el-tab-pane label="全部分类" name="" />
          <el-tab-pane
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :name="category.id"
          />
        </el-tabs>
        <div class="category-info">
          <span class="total-count">共 {{ total }} 个文档</span>
        </div>
      </div>
    </el-card>

    <!-- 文档列表 - 使用表格布局 -->
    <el-card class="doc-table-card" shadow="hover">
      <el-table
        :data="documentList"
        style="width: 100%"
        v-loading="loading"
        :empty-text="''"
      >
        <el-table-column prop="title" label="文档标题" min-width="200">
          <template #default="{ row }">
            <div class="doc-title-cell" @click="openDocument(row)">
              <el-icon class="doc-icon"><Document /></el-icon>
              <span class="doc-title">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="getCategoryTagType(getCategoryName(row.category))">
              {{ getCategoryName(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="path" label="路径" width="200">
          <template #default="{ row }">
            <span class="doc-path">{{ row.path }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="updatedAt" label="更新时间" width="160">
          <template #default="{ row }">
            <span class="doc-date">{{ formatDate(row.updatedAt) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="size" label="大小" width="100" align="right">
          <template #default="{ row }">
            <span class="doc-size">{{ formatSize(row.size) }}</span>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 空状态 - 不显示 -->
      <div v-if="!loading && (!documentList || documentList.length === 0)" class="empty-placeholder">
      </div>
    </el-card>

    <!-- 分页 -->
    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      class="pagination"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
    
    <!-- 查看文档弹窗 -->
    <el-dialog
      v-model="viewDialogVisible"
      :title="currentDocument?.title"
      fullscreen
      :close-on-click-modal="false"
      class="document-view-dialog"
    >
      <div class="document-view-content" v-if="currentDocument">
        <div class="document-meta">
          <div class="meta-item">
            <el-tag size="small" :type="getCategoryTagType(getCategoryName(currentDocument.category || ''))">
              {{ getCategoryName(currentDocument.category || '') }}
            </el-tag>
          </div>
          <div class="meta-item">
            <el-icon><Clock /></el-icon>
            <span>创建：{{ formatDate(currentDocument.createdAt || '') }}</span>
          </div>
          <div class="meta-item">
            <el-icon><Clock /></el-icon>
            <span>更新：{{ formatDate(currentDocument.updatedAt || '') }}</span>
          </div>
          <div class="meta-item" v-if="currentDocument.tags && currentDocument.tags.length">
            <el-tag
              v-for="tag in currentDocument.tags"
              :key="tag"
              size="small"
              style="margin-right: 4px"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="document-description" v-if="currentDocument.description">
          <strong>文档描述：</strong>
          <p>{{ currentDocument.description }}</p>
        </div>
        
        <div class="document-body">
          <div v-html="renderMarkdown(currentDocument.content || '')" class="markdown-content"></div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="viewDialogVisible = false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { getDocumentById, getDocumentCategories, getDocumentList } from '@/api/document'
import type { DocumentCategory, DocumentInfo } from '@/types/document'
import { Clock, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

// 路由
const router = useRouter()

// 状态
const loading = ref(false)
const documentList = ref<DocumentInfo[]>([])
const categories = ref<DocumentCategory[]>([])
const activeCategory = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchTimeout = ref<ReturnType<typeof setTimeout> | null>(null)

// 查看文档弹窗
const viewDialogVisible = ref(false)
const currentDocument = ref<DocumentInfo | null>(null)

// 分类 ID 到名称的映射
const categoryMap = ref<Record<string, string>>({})

// 获取分类名称
const getCategoryName = (categoryId: string): string => {
  if (!categoryId) return ''
  return categoryMap.value[categoryId] || categoryId
}


// 方法
const fetchCategories = async () => {
  try {
    const cats = await getDocumentCategories()
    console.log('分类 API 响应:', cats)
    
    // 处理分类数据
    if (Array.isArray(cats)) {
      categories.value = cats
    } else if (cats && Array.isArray(cats.data)) {
      categories.value = cats.data
    } else {
      categories.value = []
    }
    
    console.log('处理后的分类列表:', categories.value)
    
    // 更新分类映射
    categoryMap.value = {}
    categories.value.forEach(cat => {
      categoryMap.value[cat.id] = cat.name
    })
    console.log('分类映射:', categoryMap.value)
  } catch (error) {
    console.error('获取文档分类失败:', error)
  }
}

const fetchDocuments = async () => {
  console.log('[页面] 开始加载文档列表...')
  console.log('[页面] 当前分类:', activeCategory.value)
  console.log('[页面] 当前关键词:', searchKeyword.value)
  console.log('[页面] 当前页码:', currentPage.value)
  console.log('[页面] 每页数量:', pageSize.value)
  
  loading.value = true
  try {
    const result = await getDocumentList({
      category: activeCategory.value || undefined,
      keyword: searchKeyword.value || undefined,
      page: currentPage.value,
      pageSize: pageSize.value
    })
    
    console.log('[页面] fetchDocuments 结果:', result)
    console.log('[页面] result 是否是数组:', Array.isArray(result))
    
    // 简化处理：如果是数组，直接使用
    if (Array.isArray(result)) {
      documentList.value = result
      total.value = result.length
      console.log('[页面] ✓ 使用数组，文档数量:', result.length)
    } 
    // 如果有 list 字段
    else if (result && result.list) {
      documentList.value = result.list || []
      total.value = result.total || (result.list || []).length
      console.log('[页面] ✓ 使用 result.list，文档数量:', (result.list || []).length)
    }
    // 其他情况
    else {
      console.warn('[页面] ⚠ 未知的响应结构，使用空数组')
      documentList.value = []
      total.value = 0
    }
    
    console.log('[页面] 最终文档列表:', documentList.value.length, '个文档')
    console.log('[页面] 最终总数量:', total.value)
  } catch (error) {
    console.error('[页面] ❌ 获取文档列表失败:', error)
    ElMessage.error('加载文档列表失败')
    documentList.value = []
    total.value = 0
  } finally {
    loading.value = false
    console.log('[页面] 加载完成，loading:', loading.value)
  }
}

const handleCategoryChange = (category: string) => {
  activeCategory.value = category
  currentPage.value = 1
  fetchDocuments()
}

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }
  searchTimeout.value = setTimeout(() => {
    currentPage.value = 1
    fetchDocuments()
  }, 300)
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchDocuments()
}

const handleCurrentChange = () => {
  fetchDocuments()
}

const openDocument = async (doc: DocumentInfo) => {
  console.log('[页面] 打开文档:', doc.id)
  try {
    // 获取完整文档详情
    const fullDoc = await getDocumentById(doc.id)
    currentDocument.value = fullDoc
    viewDialogVisible.value = true
    console.log('[页面] 文档加载成功')
  } catch (error) {
    console.error('[页面] 获取文档详情失败:', error)
    ElMessage.error('加载文档详情失败')
  }
}

// 渲染 Markdown
const renderMarkdown = (content: string): string => {
  // 简单实现，后续可以引入 marked 库
  return content.replace(/\n/g, '<br>')
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days === 1) {
    return '昨天'
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN')
  }
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const getCategoryTagType = (category: string): '' | 'success' | 'warning' | 'danger' | 'primary' | 'info' => {
  const typeMap: Record<string, '' | 'success' | 'warning' | 'danger' | 'primary' | 'info'> = {
    '系统概述': '',
    '技术架构': 'success',
    '核心功能': 'warning',
    '安全增强': 'danger',
    '性能优化': 'primary',
    '测试报告': 'info'
  }
  return typeMap[category] || ''
}

// 生命周期
onMounted(() => {
  console.log('[生命周期] onMounted 被调用')
  console.log('[生命周期] 开始加载分类和文档...')
  fetchCategories()
  fetchDocuments()
})
</script>

<style scoped>
/* 页面容器 */
.doc-list-container {
  padding: 24px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: calc(100vh - 84px);
}

/* 头部卡片 */
.header-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  background: #fff;
  transition: all 0.3s ease;
}

.header-card:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  font-size: 32px;
  color: #409EFF;
  filter: drop-shadow(0 2px 4px rgba(64, 158, 255, 0.3));
}

.page-title {
  margin: 0;
  font-size: 26px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  color: #303133;
}

/* 分类卡片 */
.category-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  background: #fff;
  transition: all 0.3s ease;
}

.category-card:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.category-tabs {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-info {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 0 16px;
}

.total-count {
  color: #909399;
  font-size: 14px;
  font-weight: 500;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 表格卡片 */
.doc-table-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  background: #fff;
  overflow: hidden;
  transition: all 0.3s ease;
}

.doc-table-card:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

/* 表格样式优化 */
:deep(.el-table) {
  background: #fff;
}

:deep(.el-table__header th) {
  background: linear-gradient(135deg, #f0f2f5 0%, #e8eaed 100%);
  color: #303133;
  font-weight: 600;
  font-size: 15px;
  padding: 16px 8px;
  border-bottom: 2px solid #ebeef5;
}

:deep(.el-table__row) {
  transition: all 0.3s ease;
}

:deep(.el-table__row:hover) {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  transform: scale(1.005);
}

:deep(.el-table__cell) {
  padding: 14px 8px;
  border-bottom: 1px solid #f0f0f0;
}

/* 文档标题单元格 */
.doc-title-cell {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.doc-title-cell:hover {
  transform: translateX(4px);
}

.doc-title-cell:hover .doc-title {
  color: #409EFF;
  text-decoration: underline;
}

.doc-icon {
  font-size: 20px;
  color: #409EFF;
  filter: drop-shadow(0 2px 4px rgba(64, 158, 255, 0.3));
}

.doc-title {
  font-weight: 500;
  font-size: 15px;
  color: #303133;
  transition: all 0.3s ease;
}

/* 分类标签 */
:deep(.el-tag) {
  border-radius: 6px;
  padding: 4px 12px;
  font-size: 13px;
  font-weight: 500;
  border: none;
  transition: all 0.3s ease;
}

:deep(.el-tag:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* 路径和日期 */
.doc-path,
.doc-date,
.doc-size {
  color: #606266;
  font-size: 14px;
  font-weight: 400;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  padding: 24px 0;
  flex-shrink: 0;
  background: #fff;
  border-radius: 12px;
  margin-top: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

:deep(.el-pagination) {
  padding: 16px 24px;
}

:deep(.el-pagination button) {
  border-radius: 6px;
  transition: all 0.3s ease;
}

:deep(.el-pagination button:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

:deep(.el-pagination .el-pager li) {
  border-radius: 6px;
  transition: all 0.3s ease;
}

:deep(.el-pagination .el-pager li:hover) {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

:deep(.el-pagination .is-active) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border: none;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

/* 查看文档对话框样式 */
.document-view-dialog {
  border-radius: 16px;
  overflow: hidden;
}

.document-view-dialog .el-dialog {
  margin-top: 0 !important;
  height: calc(100vh - 40px) !important;
  top: 20px !important;
  bottom: 20px !important;
  display: flex;
  flex-direction: column;
  border-radius: 16px;
}

.document-view-dialog .el-dialog__header {
  padding: 20px 32px;
  border-bottom: 2px solid #ebeef5;
  flex-shrink: 0;
  min-height: 70px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
}

.document-view-dialog .el-dialog__title {
  font-size: 22px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.document-view-dialog .el-dialog__body {
  padding: 0;
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
  position: relative;
}

.document-view-dialog .el-dialog__footer {
  padding: 20px 32px;
  border-top: 2px solid #ebeef5;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  min-height: 64px;
  z-index: 10;
}

.document-view-dialog .document-view-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding-bottom: 80px;
  background: #fff;
}

.document-view-dialog .document-view-content .document-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  padding: 24px 40px;
  border-bottom: 2px solid #f0f0f0;
  flex-shrink: 0;
  background: linear-gradient(135deg, #f9fafb 0%, #f5f7fa 100%);
}

.document-view-dialog .document-view-content .document-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.document-view-dialog .document-view-content .document-meta .meta-item .el-icon {
  font-size: 16px;
  color: #409EFF;
}

.document-view-dialog .document-view-content .document-description {
  padding: 24px 40px;
  margin: 0;
  border-bottom: 2px solid #f0f0f0;
  flex-shrink: 0;
  background: #fff;
}

.document-view-dialog .document-view-content .document-description strong {
  color: #303133;
  font-weight: 600;
  font-size: 16px;
  display: block;
  margin-bottom: 12px;
}

.document-view-dialog .document-view-content .document-description p {
  margin: 0;
  color: #606266;
  line-height: 1.8;
  font-size: 15px;
}

.document-view-dialog .document-view-content .document-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 0;
  background: #fff;
}

.document-view-dialog .document-view-content .document-body .markdown-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px 40px;
  line-height: 1.8;
  color: #303133;
  font-size: 15px;
}

/* Markdown 内容美化 */
.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 32px;
  margin-bottom: 20px;
  font-weight: 600;
  line-height: 1.4;
  color: #303133;
}

.markdown-content :deep(h1) {
  font-size: 2.2em;
  border-bottom: 3px solid #409EFF;
  padding-bottom: 0.4em;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  padding: 16px 20px;
  border-radius: 8px;
}

.markdown-content :deep(h2) {
  font-size: 1.8em;
  border-left: 5px solid #409EFF;
  padding-left: 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  padding: 12px 20px;
  border-radius: 6px;
}

.markdown-content :deep(p) {
  margin-bottom: 20px;
  line-height: 1.8;
}

.markdown-content :deep(code) {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.9em;
  color: #e74c3c;
  border: 1px solid #ebeef5;
}

.markdown-content :deep(pre) {
  background: linear-gradient(135deg, #2c3e50 0%, #34495e 100%);
  padding: 20px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 24px 0;
  border: 1px solid #ebeef5;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.markdown-content :deep(pre code) {
  background: transparent;
  padding: 0;
  color: #ecf0f1;
  border: none;
}

.markdown-content :deep(blockquote) {
  border-left: 5px solid #409EFF;
  padding-left: 20px;
  margin: 24px 0;
  color: #606266;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  padding: 16px 24px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 28px;
  margin-bottom: 20px;
}

.markdown-content :deep(li) {
  margin-bottom: 10px;
  line-height: 1.6;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 24px 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-radius: 8px;
  overflow: hidden;
}

.markdown-content :deep(table th),
.markdown-content :deep(table td) {
  border: 1px solid #ebeef5;
  padding: 14px 16px;
  text-align: left;
}

.markdown-content :deep(table th) {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  font-weight: 600;
  color: #606266;
  border-bottom: 2px solid #409EFF;
}

.markdown-content :deep(table tr:nth-child(even)) {
  background: linear-gradient(135deg, #fafafa 0%, #f5f7fa 100%);
}

.markdown-content :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 24px auto;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.markdown-content :deep(a) {
  color: #409EFF;
  text-decoration: none;
  border-bottom: 1px dashed #409EFF;
  transition: all 0.3s ease;
}

.markdown-content :deep(a:hover) {
  color: #667eea;
  border-bottom-style: solid;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eaed 100%);
  padding: 2px 6px;
  border-radius: 4px;
}
</style>

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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Document, Clock, Folder } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getDocumentList, getDocumentCategories } from '@/api/document'
import type { DocumentInfo, DocumentCategory } from '@/types/document'

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
  loading.value = true
  try {
    const result = await getDocumentList({
      category: activeCategory.value || undefined,
      keyword: searchKeyword.value || undefined,
      page: currentPage.value,
      pageSize: pageSize.value
    })
    
    console.log('=== fetchDocuments 结果 ===')
    console.log('result:', result)
    console.log('result 类型:', typeof result)
    console.log('result 是否是数组:', Array.isArray(result))
    console.log('result 长度:', Array.isArray(result) ? result.length : 'N/A')
    
    // 简化处理：如果是数组，直接使用
    if (Array.isArray(result)) {
      documentList.value = result
      total.value = result.length
      console.log('✓ 使用数组，文档数量:', result.length)
    } 
    // 如果有 list 字段
    else if (result && result.list) {
      documentList.value = result.list
      total.value = result.total || result.list.length
      console.log('✓ 使用 result.list，文档数量:', result.list.length)
    }
    // 其他情况
    else {
      console.warn('⚠ 未知的响应结构，使用空数组')
      documentList.value = []
      total.value = 0
    }
    
    console.log('最终文档列表:', documentList.value.length, '个文档')
    console.log('最终总数量:', total.value)
  } catch (error) {
    console.error('❌ 获取文档列表失败:', error)
    ElMessage.error('加载文档列表失败')
    documentList.value = []
    total.value = 0
  } finally {
    loading.value = false
    console.log('加载完成，loading:', loading.value)
  }
}

const handleCategoryChange = (category: string) => {
  activeCategory.value = category
  currentPage.value = 1
  fetchDocuments()
}

const handleSearch = () => {
  searchTimeout.value(() => {
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

const openDocument = (doc: DocumentInfo) => {
  router.push(`/docs/${doc.id}`)
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

const getCategoryTagType = (category: string) => {
  const typeMap: Record<string, any> = {
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
  fetchCategories()
  fetchDocuments()
})
</script>

<style scoped>
.doc-list-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 84px);
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.header-card {
  flex-shrink: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  font-size: 28px;
  color: #409EFF;
}

.page-title {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.category-card {
  flex-shrink: 0;
}

.category-tabs {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.category-info {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.total-count {
  color: #909399;
  font-size: 14px;
}

.doc-table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.doc-table-card :deep(.el-card__body) {
  flex: 1;
  padding: 0;
  display: flex;
  flex-direction: column;
}

.doc-table-card :deep(.el-table) {
  flex: 1;
}

.doc-title-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #409EFF;
}

.doc-title-cell:hover .doc-title {
  text-decoration: underline;
}

.doc-icon {
  font-size: 18px;
  color: #409EFF;
}

.doc-title {
  font-weight: 500;
  color: #303133;
}

.doc-path {
  color: #909399;
  font-size: 13px;
}

.doc-date {
  color: #909399;
  font-size: 13px;
}

.doc-size {
  color: #909399;
  font-size: 13px;
}

.empty-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.pagination {
  display: flex;
  justify-content: center;
  padding: 20px 0;
  flex-shrink: 0;
}
</style>

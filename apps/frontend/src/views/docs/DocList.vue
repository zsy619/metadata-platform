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

    <!-- 文档列表 -->
    <el-row :gutter="20" class="doc-grid">
      <el-col
        v-for="doc in documentList"
        :key="doc.id"
        :xs="24"
        :sm="12"
        :md="8"
        :lg="8"
        :xl="6"
        class="doc-col"
      >
        <el-card class="doc-card" shadow="hover" @click="openDocument(doc)">
          <template #header>
            <div class="doc-card-header">
              <el-icon class="doc-icon"><Document /></el-icon>
              <el-tag size="small" :type="getCategoryTagType(doc.category)">{{ doc.category }}</el-tag>
            </div>
          </template>
          
          <div class="doc-card-body">
            <h3 class="doc-title">{{ doc.title }}</h3>
            <p class="doc-description">{{ doc.description }}</p>
          </div>
          
          <div class="doc-card-footer">
            <div class="doc-meta">
              <el-icon><Clock /></el-icon>
              <span>{{ formatDate(doc.updatedAt) }}</span>
            </div>
            <div class="doc-size">
              <el-icon><Folder /></el-icon>
              <span>{{ formatSize(doc.size) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 空状态 -->
    <el-empty v-if="!loading && (!documentList || documentList.length === 0)" description="暂无文档" />

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="6" animated />
    </div>

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

// 计算属性
const searchTimeout = computed(() => {
  let timeout: NodeJS.Timeout | null = null
  return (fn: () => void, delay: number) => {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(fn, delay)
  }
})

// 方法
const fetchCategories = async () => {
  try {
    categories.value = await getDocumentCategories()
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
    documentList.value = result.list || []
    total.value = result.total || 0
  } catch (error) {
    console.error('获取文档列表失败:', error)
    ElMessage.error('加载文档列表失败')
    documentList.value = []
    total.value = 0
  } finally {
    loading.value = false
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
}

.header-card {
  margin-bottom: 20px;
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
  margin-bottom: 20px;
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

.doc-grid {
  margin-bottom: 20px;
}

.doc-col {
  margin-bottom: 20px;
}

.doc-card {
  height: 100%;
  cursor: pointer;
  transition: transform 0.2s;
}

.doc-card:hover {
  transform: translateY(-4px);
}

.doc-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.doc-icon {
  font-size: 20px;
  color: #409EFF;
}

.doc-card-body {
  padding: 8px 0;
}

.doc-title {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.doc-description {
  margin: 0;
  font-size: 13px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  line-height: 1.5;
}

.doc-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.doc-meta,
.doc-size {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;
}

.loading-container {
  padding: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}
</style>

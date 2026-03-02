<template>
  <div class="doc-detail-container">
    <!-- 返回按钮 -->
    <div class="back-button">
      <el-button @click="handleBack" :icon="ArrowLeft">返回文档列表</el-button>
    </div>

    <!-- 文档头部 -->
    <el-card class="doc-header-card" shadow="hover">
      <div class="doc-header">
        <h1 class="doc-title">{{ document.title }}</h1>
        <div class="doc-meta">
          <el-tag :type="getCategoryTagType(document.category)">{{ document.category }}</el-tag>
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            更新时间：{{ formatDate(document.updatedAt) }}
          </span>
          <span class="meta-item">
            <el-icon><Folder /></el-icon>
            大小：{{ formatSize(document.size) }}
          </span>
          <span class="meta-item" v-if="document.viewCount">
            <el-icon><View /></el-icon>
            阅读：{{ document.viewCount }}
          </span>
        </div>
      </div>
    </el-card>

    <!-- 文档内容 -->
    <el-card class="doc-content-card" shadow="hover">
      <div class="doc-content">
        <!-- Markdown 渲染区域 -->
        <div v-if="document.content" class="markdown-body" v-html="renderedContent"></div>
        <el-empty v-else description="文档内容为空" />
      </div>
    </el-card>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <el-skeleton :rows="10" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Clock, Folder, View } from '@element-plus/icons-vue'
import { getDocumentById } from '@/api/document'
import type { DocumentDetail } from '@/types/document'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const document = ref<DocumentDetail>({
  id: '',
  title: '',
  category: '',
  path: '',
  description: '',
  content: '',
  size: 0,
  createdAt: '',
  updatedAt: ''
})

// 简单的 Markdown 渲染（可以使用 marked 库来增强）
const renderedContent = computed(() => {
  if (!document.value.content) return ''
  
  // 简单的 Markdown 转换
  let content = document.value.content
    // 标题
    .replace(/^# (.*$)/gim, '<h1>$1</h1>')
    .replace(/^## (.*$)/gim, '<h2>$1</h2>')
    .replace(/^### (.*$)/gim, '<h3>$1</h3>')
    // 粗体
    .replace(/\*\*(.*)\*\*/gim, '<b>$1</b>')
    // 斜体
    .replace(/\*(.*)\*/gim, '<i>$1</i>')
    // 链接
    .replace(/\[([^\]]+)\]\(([^)]+)\)/gim, '<a href="$2" target="_blank">$1</a>')
    // 代码块
    .replace(/```([\s\S]*?)```/gim, '<pre><code>$1</code></pre>')
    // 行内代码
    .replace(/`([^`]+)`/gim, '<code>$1</code>')
    // 换行
    .replace(/\n/gim, '<br>')
  
  return content
})

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const formatSize = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const getCategoryTagType = (category: string) => {
  const typeMap: Record<string, string> = {
    '系统概述': 'info',
    '技术架构': 'primary',
    '核心功能': 'success',
    '安全增强': 'warning',
    '性能优化': 'danger',
    '测试报告': 'info',
    '部署指南': 'success',
    '最佳实践': 'warning'
  }
  return typeMap[category] || ''
}

const handleBack = () => {
  router.push('/docs')
}

const loadDocument = async () => {
  const docId = route.params.id as string
  if (!docId) return
  
  loading.value = true
  try {
    const res = await getDocumentById(docId)
    document.value = res.data || res
  } catch (error) {
    console.error('加载文档失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDocument()
})
</script>

<style scoped>
.doc-detail-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.back-button {
  margin-bottom: 20px;
}

.doc-header-card {
  margin-bottom: 20px;
}

.doc-header {
  padding: 10px 0;
}

.doc-title {
  margin: 0 0 15px 0;
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.doc-meta {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #909399;
  font-size: 14px;
}

.doc-content-card {
  min-height: 500px;
}

.doc-content {
  padding: 20px 0;
}

.markdown-body {
  line-height: 1.8;
  color: #303133;
}

.markdown-body :deep(h1) {
  font-size: 24px;
  margin-top: 30px;
  margin-bottom: 20px;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.markdown-body :deep(h2) {
  font-size: 20px;
  margin-top: 25px;
  margin-bottom: 15px;
  border-bottom: 1px solid #eaecef;
  padding-bottom: 0.3em;
}

.markdown-body :deep(h3) {
  font-size: 18px;
  margin-top: 20px;
  margin-bottom: 12px;
}

.markdown-body :deep(p) {
  margin: 15px 0;
}

.markdown-body :deep(code) {
  background-color: #f6f8fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 85%;
}

.markdown-body :deep(pre) {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
  margin: 15px 0;
}

.markdown-body :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.markdown-body :deep(a) {
  color: #409EFF;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.loading-overlay {
  margin-top: 20px;
}
</style>

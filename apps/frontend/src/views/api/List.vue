<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Connection />
                </el-icon>
                接口管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">
                    新增接口
                </el-button>
                <el-button type="warning" @click="handleBatchGenerate" :icon="DocumentAdd">
                    批量生成
                </el-button>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入接口名称、编码或路径搜索" clearable :prefix-icon="Search" style="width: 300px" @clear="handleSearch" @keyup.enter="handleSearch" />
                <el-select v-model="filterType" placeholder="筛选接口类型" style="width: 180px; margin-left: 10px" @change="handleSearch" clearable>
                    <el-option label="全部" value="" />
                    <el-option label="查询" :value="1" />
                    <el-option label="新增" :value="2" />
                    <el-option label="更新" :value="3" />
                    <el-option label="删除" :value="4" />
                    <el-option label="自定义" :value="5" />
                </el-select>
                <el-select v-model="filterState" placeholder="筛选状态" style="width: 120px; margin-left: 10px" @change="handleSearch" clearable>
                    <el-option label="全部" value="" />
                    <el-option label="启用" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
                    搜索
                </el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="apis" border style="width: 100%" height="100%">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关接口' : '暂无接口数据'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增接口</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="apiName" label="接口名称" width="200" />
                    <el-table-column prop="apiCode" label="接口编码" width="180" />
                    <el-table-column prop="apiPath" label="接口路径" width="250">
                        <template #default="scope">
                            <div class="api-path">
                                <el-tag size="small" :type="getMethodTagType(scope.row.apiMethod)">
                                    {{ scope.row.apiMethod }}
                                </el-tag>
                                <span>{{ scope.row.apiPath }}</span>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="modelName" label="所属模型" width="150" />
                    <el-table-column prop="apiType" label="接口类型" width="120">
                        <template #default="scope">
                            <el-tag :type="getAPITypeTagType(scope.row.apiType)">
                                {{ getAPITypeText(scope.row.apiType) }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="state" label="状态" width="100">
                        <template #default="scope">
                            <el-switch v-model="scope.row.state" :active-value="1" :inactive-value="0" @change="handleStateChange(scope.row)" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="needAuth" label="需要鉴权" width="120">
                        <template #default="scope">
                            <el-tag :type="scope.row.needAuth ? 'success' : 'info'">
                                {{ scope.row.needAuth ? '需要' : '不需要' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="createAt" label="创建时间" width="180" />
                    <el-table-column label="操作" width="250" fixed="right" class-name="action-column">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="View" @click="handleTest(scope.row)">
                                测试
                            </el-button>
                            <el-button type="success" size="small" :icon="Edit" @click="handleEdit(scope.row)">
                                编辑
                            </el-button>
                            <el-button type="info" size="small" :icon="Document" @click="handleDocument(scope.row)">
                                文档
                            </el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页区域 -->
            <div class="pagination-area">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" background layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import type { API } from '@/types/api'
import { showDeleteConfirm } from '@/utils/confirm'
import {
    Connection,
    Delete,
    Document,
    DocumentAdd,
    Edit,
    Plus,
    RefreshLeft,
    Search,
    View
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterType = ref('')
const filterState = ref('')
const apis = ref<API[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 生命周期钩子
onMounted(() => {
    fetchAPIs()
})

// 获取接口列表
const fetchAPIs = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        // TODO: 替换为真实 API 调用
        // const response: any = await getAPIs({ ... })
        // 目前保留模拟逻辑，但结构适配真实 API 返回

        // 模拟数据
        const mockData = [
            {
                apiID: 1,
                apiName: '获取用户列表',
                apiCode: 'user_list',
                apiPath: '/api/data/user_model/list',
                apiMethod: 'GET',
                modelID: 1,
                modelName: '用户模型',
                apiType: 1,
                state: 1,
                needAuth: true,
                needAudit: true,
                remark: '获取用户列表接口',
                createID: 1,
                createBy: 'admin',
                createAt: '2024-01-23 10:00:00',
                updateID: 1,
                updateBy: 'admin',
                updateAt: '2024-01-23 10:00:00'
            },
            {
                apiID: 2,
                apiName: '创建用户',
                apiCode: 'user_create',
                apiPath: '/api/data/user_model/create',
                apiMethod: 'POST',
                modelID: 1,
                modelName: '用户模型',
                apiType: 2,
                state: 1,
                needAuth: true,
                needAudit: true,
                remark: '创建用户接口',
                createID: 1,
                createBy: 'admin',
                createAt: '2024-01-23 10:00:00',
                updateID: 1,
                updateBy: 'admin',
                updateAt: '2024-01-23 10:00:00'
            }
        ]

        // 模拟请求延迟
        await new Promise(resolve => setTimeout(resolve, 300))

        apis.value = mockData
        total.value = mockData.length
    } catch (error) {
        console.error('获取接口列表失败:', error)
        ElMessage.error('获取接口列表失败')
    } finally {
        loading.value = false
    }
}

// 搜索
const handleSearch = () => {
    currentPage.value = 1
    fetchAPIs()
}

// 重置
const handleReset = () => {
    searchQuery.value = ''
    filterType.value = ''
    filterState.value = ''
    currentPage.value = 1
    fetchAPIs()
}

// 页码变化
const handleCurrentChange = (val: number) => {
    currentPage.value = val
    fetchAPIs()
}

// 每页条数变化
const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
    fetchAPIs()
}

// 状态变化
const handleStateChange = async (row: API) => {
    try {
        if (row.state === 1) {
            ElMessage.success('接口已启用')
        } else {
            ElMessage.success('接口已禁用')
        }
    } catch (error) {
        console.error('更新接口状态失败:', error)
        ElMessage.error('更新接口状态失败')
        row.state = row.state === 1 ? 0 : 1
    }
}

const handleCreate = () => router.push('/api/create')
const handleBatchGenerate = () => router.push('/api/batch-generate')
const handleEdit = (row: API) => router.push(`/api/${row.apiID}/edit`)

const handleDelete = (row: API) => {
    showDeleteConfirm(`确定要删除接口 "${row.apiName}" 吗？`).then(async () => {
        loadingText.value = '正在删除...'
        loading.value = true
        // Simulate async delete
        await new Promise(resolve => setTimeout(resolve, 500))
        ElMessage.success('删除成功')
        fetchAPIs() // resets loading
    })
}

const handleTest = (row: API) => router.push(`/api/${row.apiID}/test`)
const handleDocument = (row: API) => router.push(`/api/${row.apiID}/document`)

const getAPITypeText = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: '查询', 2: '新增', 3: '更新', 4: '删除', 5: '自定义'
    }
    return typeMap[type] || '未知'
}

const getAPITypeTagType = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: 'info', 2: 'success', 3: 'warning', 4: 'danger', 5: 'primary'
    }
    return typeMap[type] || 'info'
}

const getMethodTagType = (method: string): string => {
    const methodMap: Record<string, string> = {
        'GET': 'success', 'POST': 'primary', 'PUT': 'warning', 'DELETE': 'danger', 'PATCH': 'info'
    }
    return methodMap[method] || 'info'
}
</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
.container-padding {
    padding: 20px;
    padding-bottom: 0;
    height: calc(100vh - 84px);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    box-sizing: border-box;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.page-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin: 0;
}

.title-icon {
    font-size: 24px;
    color: #409eff;
}

.header-actions {
    display: flex;
    gap: 10px;
}

.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

:deep(.el-card__body) {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow: hidden;
    box-sizing: border-box;
}

.search-area {
    flex-shrink: 0;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
}

.table-area {
    flex: 1;
    overflow: hidden;
    margin-bottom: 20px;
}

.pagination-area {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
}

.api-path {
    display: flex;
    align-items: center;
    gap: 8px;
}

.api-path span {
    word-break: break-all;
}

:deep(.action-column .cell) {
    white-space: nowrap;
}
</style>

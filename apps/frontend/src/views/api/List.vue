<template>
    <div class="api-list">
        <div class="page-header">
            <h1>接口管理</h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">
                    新增接口
                </el-button>
                <el-button type="warning" @click="handleBatchGenerate" :icon="DocumentAdd">
                    批量生成
                </el-button>
            </div>
        </div>
        <el-card>
            <div class="card-header">
                <el-input v-model="searchQuery" placeholder="请输入接口名称、编码或路径搜索" clearable prefix-icon="Search" style="width: 300px" />
                <el-select v-model="filterType" placeholder="筛选接口类型" style="width: 180px; margin-left: 10px">
                    <el-option label="全部" value="" />
                    <el-option label="查询" :value="1" />
                    <el-option label="新增" :value="2" />
                    <el-option label="更新" :value="3" />
                    <el-option label="删除" :value="4" />
                    <el-option label="自定义" :value="5" />
                </el-select>
                <el-select v-model="filterState" placeholder="筛选状态" style="width: 120px; margin-left: 10px">
                    <el-option label="全部" value="" />
                    <el-option label="启用" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search">
                    搜索
                </el-button>
            </div>
            <el-table v-loading="loading" :data="apis" border style="width: 100%">
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
                        <el-switch v-model="scope.row.state" active-value="1" inactive-value="0" @change="handleStateChange(scope.row)" />
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
            <div class="pagination">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import type { API } from '@/types/api'
import {
    Delete,
    Document,
    DocumentAdd,
    Edit,
    Plus,
    Search,
    View
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
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
    loading.value = true
    try {
        // 模拟API调用，实际需要替换为真实API
        // const response = await getAPIs({
        //   page: currentPage.value,
        //   pageSize: pageSize.value,
        //   search: searchQuery.value,
        //   apiType: filterType.value ? Number(filterType.value) : undefined,
        //   state: filterState.value ? Number(filterState.value) : undefined
        // })
        // apis.value = response.data
        // total.value = response.total

        // 模拟数据
        apis.value = [
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
            },
            {
                apiID: 3,
                apiName: '更新用户',
                apiCode: 'user_update',
                apiPath: '/api/data/user_model/update',
                apiMethod: 'PUT',
                modelID: 1,
                modelName: '用户模型',
                apiType: 3,
                state: 0,
                needAuth: true,
                needAudit: true,
                remark: '更新用户接口',
                createID: 1,
                createBy: 'admin',
                createAt: '2024-01-23 10:00:00',
                updateID: 1,
                updateBy: 'admin',
                updateAt: '2024-01-23 10:00:00'
            },
            {
                apiID: 4,
                apiName: '删除用户',
                apiCode: 'user_delete',
                apiPath: '/api/data/user_model/delete',
                apiMethod: 'DELETE',
                modelID: 1,
                modelName: '用户模型',
                apiType: 4,
                state: 1,
                needAuth: true,
                needAudit: true,
                remark: '删除用户接口',
                createID: 1,
                createBy: 'admin',
                createAt: '2024-01-23 10:00:00',
                updateID: 1,
                updateBy: 'admin',
                updateAt: '2024-01-23 10:00:00'
            }
        ]
        total.value = apis.value.length
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

// 页码变化
const handleCurrentChange = (page: number) => {
    currentPage.value = page
    fetchAPIs()
}

// 每页条数变化
const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchAPIs()
}

// 状态变化
const handleStateChange = async (row: API) => {
    try {
        if (row.state === 1) {
            // 启用接口
            // await enableAPI(row.apiID)
            ElMessage.success('接口已启用')
        } else {
            // 禁用接口
            // await disableAPI(row.apiID)
            ElMessage.success('接口已禁用')
        }
    } catch (error) {
        console.error('更新接口状态失败:', error)
        ElMessage.error('更新接口状态失败')
        // 恢复原状态
        row.state = row.state === 1 ? 0 : 1
    }
}

// 新增接口
const handleCreate = () => {
    router.push('/apis/create')
}

// 批量生成接口
const handleBatchGenerate = () => {
    router.push('/apis/batch-generate')
}

// 编辑接口
const handleEdit = (row: API) => {
    router.push(`/apis/${row.apiID}/edit`)
}

// 删除接口
const handleDelete = (row: API) => {
    ElMessageBox.confirm(
        `确定要删除接口 "${row.apiName}" 吗？`,
        '删除确认',
        {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
        }
    ).then(async () => {
        try {
            // 模拟删除操作
            // await deleteAPI(row.apiID)
            ElMessage.success('删除成功')
            fetchAPIs()
        } catch (error) {
            console.error('删除接口失败:', error)
            ElMessage.error('删除失败')
        }
    }).catch(() => {
        // 取消删除
    })
}

// 测试接口
const handleTest = (row: API) => {
    router.push(`/apis/${row.apiID}/test`)
}

// 查看文档
const handleDocument = (row: API) => {
    router.push(`/apis/${row.apiID}/document`)
}

// 获取接口类型文本
const getAPITypeText = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: '查询',
        2: '新增',
        3: '更新',
        4: '删除',
        5: '自定义'
    }
    return typeMap[type] || '未知'
}

// 获取接口类型标签样式
const getAPITypeTagType = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: 'info',
        2: 'success',
        3: 'warning',
        4: 'danger',
        5: 'primary'
    }
    return typeMap[type] || 'info'
}

// 获取请求方法标签样式
const getMethodTagType = (method: string): string => {
    const methodMap: Record<string, string> = {
        'GET': 'success',
        'POST': 'primary',
        'PUT': 'warning',
        'DELETE': 'danger',
        'PATCH': 'info',
        'OPTIONS': 'info'
    }
    return methodMap[method] || 'info'
}
</script>
<style scoped>
.api-list {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.header-actions {
    display: flex;
    gap: 10px;
}

.card-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
}

.pagination {
    margin-top: 20px;
    text-align: right;
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

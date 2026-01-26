<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary">数据源管理</h1>
            <el-button type="primary" @click="handleCreate" :icon="Plus">
                新增数据源
            </el-button>
        </div>
        <el-card>
            <div class="flex-center m-b-lg">
                <el-input v-model="searchQuery" placeholder="请输入数据源名称搜索" clearable prefix-icon="Search" style="width: 300px" />
                <el-select v-model="filterType" placeholder="筛选数据源类型" style="width: 180px; margin-left: 10px">
                    <el-option label="全部" value="" />
                    <el-option label="MySQL" value="MySQL" />
                    <el-option label="PostgreSQL" value="PostgreSQL" />
                    <el-option label="SQL Server" value="SQL Server" />
                    <el-option label="Oracle" value="Oracle" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
                    搜索
                </el-button>
            </div>
            <el-table v-loading="loading" :data="dataSources" border stripe style="width: 100%">
                <el-table-column prop="connName" label="数据源名称" width="200" />
                <el-table-column prop="connKind" label="数据源类型" width="150" />
                <el-table-column prop="connVersion" label="版本" width="120" />
                <el-table-column prop="connHost" label="主机地址" width="200" />
                <el-table-column prop="connPort" label="端口" width="80" />
                <el-table-column prop="connDatabase" label="数据库" width="150" />
                <el-table-column prop="state" label="状态" width="120">
                    <template #default="scope">
                        <el-tag :type="scope.row.state === 1 ? 'success' :
                                scope.row.state === 0 ? 'danger' : 'warning'
                            ">
                            {{ scope.row.state === 1 ? '有效' :
                                scope.row.state === 0 ? '无效' : '未检测' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="createAt" label="创建时间" width="180" />
                <el-table-column label="操作" width="220" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text :icon="Connection" @click="handleTestConnection(scope.row)" :disabled="scope.row.state === 1">
                            测试
                        </el-button>
                        <el-button type="primary" size="small" text :icon="Edit" @click="handleEdit(scope.row)">
                            编辑
                        </el-button>
                        <el-button type="danger" size="small" text :icon="Delete" @click="handleDelete(scope.row)">
                            删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="flex-between m-t-lg">
                <div></div> <!-- Spacer -->
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import type { DataSource } from '@/types/data-source'
import {
    Connection,
    Delete,
    Edit,
    Plus,
    Search
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const filterType = ref('')
const dataSources = ref<DataSource[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 生命周期钩子
onMounted(() => {
    fetchDataSources()
})

const fetchDataSources = async () => {
    loading.value = true
    try {
        // 模拟数据
        dataSources.value = [
            {
                id: 1,
                connID: 800100001,
                parentID: 0,
                connName: '生产MySQL库',
                connKind: 'MySQL',
                connVersion: '8.0',
                connHost: '192.168.1.100',
                connPort: 3306,
                connUser: 'admin',
                connPassword: '***',
                connDatabase: 'metadata',
                connConn: '',
                isDeleted: false,
                state: 1,
                remark: '核心元数据库',
                sort: 0,
                createdAt: '2024-01-23 10:00:00',
                updatedAt: '2024-01-23 10:00:00'
            },
            {
                id: 2,
                connID: 800100002,
                parentID: 0,
                connName: '测试PG库',
                connKind: 'PostgreSQL',
                connVersion: '14.0',
                connHost: '192.168.1.101',
                connPort: 5432,
                connUser: 'postgres',
                connPassword: '***',
                connDatabase: 'test_db',
                connConn: '',
                isDeleted: false,
                state: -1,
                remark: '测试环境',
                sort: 0,
                createdAt: '2024-01-23 11:00:00',
                updatedAt: '2024-01-23 11:00:00'
            }
        ]
        total.value = dataSources.value.length
    } catch (error) {
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    currentPage.value = 1
    fetchDataSources()
}

const handleCurrentChange = (page: number) => {
    currentPage.value = page
    fetchDataSources()
}

const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchDataSources()
}

const handleCreate = () => router.push('/data-sources/create')
const handleEdit = (row: DataSource) => router.push(`/data-sources/${row.id}/edit`)

const handleDelete = (row: DataSource) => {
    ElMessageBox.confirm(`确定要删除数据源 "${row.connName}" 吗？`, '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        ElMessage.success('操作成功')
        fetchDataSources()
    })
}

const handleTestConnection = (row: DataSource) => {
    loading.value = true
    setTimeout(() => {
        ElMessage.success('连接成功')
        row.state = 1
        loading.value = false
    }, 800)
}
</script>
<style scoped>
/* 依赖全局 CSS，本地仅保留极少量布局微调 */
</style>
<template>
    <div class="data-source-list">
        <div class="page-header">
            <h1 class="text-primary">数据源管理</h1>
            <el-button type="primary" @click="handleCreate" :icon="Plus">
                新增数据源
            </el-button>
        </div>
        <el-card>
            <div class="flex-center m-b-lg">
                <el-input v-model="searchQuery" placeholder="请输入数据源名称搜索" clearable :prefix-icon="Search" style="width: 300px" />
                <el-select v-model="filterType" placeholder="筛选数据源类型" style="width: 180px; margin-left: 10px">
                    <el-option label="全部" value="" />
                    <el-option-group label="关系型">
                        <el-option label="MySQL" value="MySQL" />
                        <el-option label="PostgreSQL" value="PostgreSQL" />
                        <el-option label="SQL Server" value="SQL Server" />
                        <el-option label="Oracle" value="Oracle" />
                        <el-option label="TiDB" value="TiDB" />
                        <el-option label="OceanBase" value="OceanBase" />
                        <el-option label="SQLite" value="SQLite" />
                    </el-option-group>
                    <el-option-group label="大数据/分析">
                        <el-option label="ClickHouse" value="ClickHouse" />
                        <el-option label="Doris" value="Doris" />
                        <el-option label="StarRocks" value="StarRocks" />
                    </el-option-group>
                    <el-option-group label="国产化">
                        <el-option label="Dameng (DM)" value="DM" />
                        <el-option label="Kingbase" value="Kingbase" />
                        <el-option label="OpenGauss" value="OpenGauss" />
                    </el-option-group>
                    <el-option-group label="NoSQL">
                        <el-option label="MongoDB" value="MongoDB" />
                        <el-option label="Redis" value="Redis" />
                    </el-option-group>
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
                    搜索
                </el-button>
            </div>
            <el-table v-loading="loading" :data="dataSources" border stripe style="width: 100%">
                <el-table-column prop="conn_name" label="数据源名称" width="200" />
                <el-table-column prop="conn_kind" label="数据源类型" width="150" />
                <el-table-column prop="conn_version" label="版本" width="120" />
                <el-table-column prop="conn_host" label="主机地址" width="200" />
                <el-table-column prop="conn_port" label="端口" width="80" />
                <el-table-column prop="conn_database" label="数据库" width="150" />
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
                <el-table-column prop="create_at" label="创建时间" width="180">
                    <template #default="scope">
                        {{ formatDateTime(scope.row.create_at) }}
                    </template>
                </el-table-column>
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
import { deleteConn, getConns, testConn } from '@/api/metadata'
import type { MdConn } from '@/types/metadata'
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
const dataSources = ref<MdConn[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 时间格式化函数
const formatDateTime = (dateStr: string | undefined) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return '-'

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')

    return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 生命周期钩子
onMounted(() => {
    fetchDataSources()
})

const fetchDataSources = async () => {
    loading.value = true
    try {
        const response: any = await getConns()
        // 后端返回结构为 { code: 200, message: "...", data: [...] }
        const data = response?.data || []

        // 过滤和搜索(前端模拟,如果后端没有搜索接口)
        let result = [...data]
        if (searchQuery.value) {
            result = result.filter(item => item.conn_name.includes(searchQuery.value))
        }
        if (filterType.value) {
            result = result.filter(item => item.conn_kind === filterType.value)
        }

        dataSources.value = result
        total.value = result.length
    } catch (error) {
        console.error('加载数据源列表失败:', error)
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
const handleEdit = (row: MdConn) => router.push(`/data-sources/${row.id}/edit`)

const handleDelete = (row: MdConn) => {
    ElMessageBox.confirm(`确定要删除数据源 "${row.conn_name}" 吗?`, '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        try {
            await deleteConn(row.id)
            ElMessage.success('删除成功')
            fetchDataSources()
        } catch (error) {
            console.error('删除数据源失败:', error)
        }
    })
}

const handleTestConnection = async (row: MdConn) => {
    loading.value = true
    try {
        const res = await testConn(row.id)
        if (res && res.success) {
            ElMessage.success('连接成功')
            row.state = 1
        } else {
            ElMessage.error(res?.message || '连接失败')
            row.state = 0
        }
    } catch (error) {
        console.error('测试连接失败:', error)
    } finally {
        loading.value = false
    }
}
</script>
<style scoped>
.data-source-list {
    padding: 20px;
}

.list-card {
    border: none;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 8px;
}

.title-icon {
    font-size: 20px;
    color: var(--el-color-primary);
}

.title-text {
    font-size: 18px;
    font-weight: 600;
}

.filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.filter-left {
    display: flex;
    align-items: center;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}
</style>
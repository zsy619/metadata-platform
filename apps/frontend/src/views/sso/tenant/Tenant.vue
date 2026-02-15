<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <OfficeBuilding />
                </el-icon>
                租户管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">
                    新增租户
                </el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入租户名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="有效" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
                    搜索
                </el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">
                    重置
                </el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredData" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关租户' : '暂无租户'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增租户</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="tenant_name" label="租户名称" width="180" show-overflow-tooltip />
                    <el-table-column prop="tenant_code" label="租户编码" width="150" />
                    <el-table-column prop="status" label="状态" width="100">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="创建时间" width="220">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注" show-overflow-tooltip />
                    <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>
                编辑
              </el-button>
              <el-button v-if="!scope.row.is_system" type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>
                删除
              </el-button>
            </template>
          </el-table-column>
                </el-table>
            </div>
        </el-card>
        <TenantForm v-model="dialogVisible" :data="formData" @success="fetchData" />
    </div>
</template>
<script setup lang="ts">
import { deleteTenant, getTenants } from '@/api/user'
import { Delete, Edit, OfficeBuilding, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import TenantForm from './TenantForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const filteredData = computed(() => {
    let data = allData.value
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        data = data.filter(item =>
            (item.tenant_name || '').toLowerCase().includes(query) ||
            (item.tenant_code || '').toLowerCase().includes(query)
        )
    }
    if (filterStatus.value !== '') {
        data = data.filter(item => item.status === filterStatus.value)
    }
    return data
})

const dialogVisible = ref(false)
const formData = ref<any>({})

const formatDateTime = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const fetchData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getTenants()
        allData.value = res.data || res
    } catch (error) {
        console.error('加载租户列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
    formData.value = { status: 1 }
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    dialogVisible.value = true
}

const handleDelete = (row: any) => {
    console.log('删除租户:', row)
    ElMessageBox.confirm(`确定要删除租户 "${row.tenant_name}" 吗？`, '提示', {
        type: 'warning'
    }).then(async () => {
        try {
            console.log('调用删除API, id:', row.id)
            const res = await deleteTenant(row.id)
            console.log('删除结果:', res)
            ElMessage.success('删除成功')
            fetchData()
        } catch (error: any) {
            console.error('删除失败:', error)
            ElMessage.error(error.message || '删除失败')
        }
    }).catch(() => { })
}

onMounted(() => { fetchData() })
</script>
<style scoped>
.table-area {
    flex: 1;
    overflow: hidden;
    min-height: 0;
    display: flex;
    flex-direction: column;
}

.table-area :deep(.el-table) {
    flex: 1;
}
</style>

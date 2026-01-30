<template>
    <div class="model-preview">
        <div class="page-header">
            <el-button @click="handleBack" :icon="ArrowLeft">返回列表</el-button>
            <h1>数据预览: {{ modelName }}</h1>
            <div class="actions">
                <el-button type="success" :icon="Download" @click="handleExport">导出导出 JSON</el-button>
                <el-button type="primary" :icon="Upload" @click="importVisible = true">导入数据</el-button>
            </div>
        </div>
        <el-card>
            <div class="filter-bar">
                <!-- 动态生成查询表单 (暂时简化) -->
                <el-input v-model="search" placeholder="搜索内容" style="width: 200px" />
                <el-button type="primary" :icon="Search" @click="fetchData">查询</el-button>
            </div>
            <el-table :data="tableData" border v-loading="loading">
                <el-table-column v-for="col in columns" :key="col.columnName" :prop="col.columnName" :label="col.showTitle || col.columnTitle" />
            </el-table>
            <div class="pagination-container">
                <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="fetchData" @current-change="fetchData" />
            </div>
        </el-card>
        <!-- 导入对话框 -->
        <el-dialog v-model="importVisible" title="导入数据" width="400px" class="custom-dialog">
            <el-upload drag action="#" :auto-upload="false" :on-change="handleFileChange">
                <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                <div class="el-upload__text">拖拽文件或 <em>点击上传</em></div>
            </el-upload>
            <template #footer>
                <el-button @click="importVisible = false">取消</el-button>
                <el-button type="primary" @click="submitImport" :loading="importing">提交</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { exportData, importData } from '@/api/io'
import { previewModelData } from '@/api/model'
import { ArrowLeft, Download, Search, Upload, UploadFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const modelId = route.params.id as string
const modelName = ref('')
const columns = ref<any[]>([])
const tableData = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const search = ref('')

const importVisible = ref(false)
const importing = ref(false)
const selectedFile = ref(null)

onMounted(() => {
    loadMetadata()
    fetchData()
})

const loadMetadata = async () => {
    try {
        // const model = await getModelById(Number(modelId))
        // modelName.value = model.modelName
        modelName.value = '测试模型'

        // const fields = await getModelFields(Number(modelId))
        // columns.value = fields.filter(f => f.isShow === 1)
        columns.value = [{ columnName: 'id', columnTitle: 'ID' }, { columnName: 'name', columnTitle: '名称' }]
    } catch (e) {
        ElMessage.error('加载模型元数据失败')
    }
}

const fetchData = async () => {
    loading.value = true
    try {
        const res = await previewModelData(String(modelId), { page: page.value, pageSize: pageSize.value, search: search.value })
        tableData.value = res.data || []
        total.value = res.total || 0
    } catch (e) {
        // ElMessage.error('加载数据失败')
        tableData.value = [{ id: '1', name: '测试1' }]
        total.value = 1
    } finally {
        loading.value = false
    }
}

const handleBack = () => router.push('/metadata/model/list')

const handleExport = async () => {
    try {
        const blob = await exportData(modelName.value, 'json', { search: search.value })
        const url = window.URL.createObjectURL(new Blob([blob]))
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `${modelName.value}.json`)
        document.body.appendChild(link)
        link.click()
        ElMessage.success('导出已开始')
    } catch (e) {
        ElMessage.error('导出失败')
    }
}

const handleFileChange = (file: any) => {
    selectedFile.value = file.raw
}

const submitImport = async () => {
    if (!selectedFile.value) {
        ElMessage.warning('请选择文件')
        return
    }
    importing.value = true
    try {
        await importData(modelName.value, selectedFile.value)
        ElMessage.success('导入成功')
        importVisible.value = false
        fetchData()
    } catch (e) {
        ElMessage.error('导入失败')
    } finally {
        importing.value = false
    }
}
</script>
<style scoped>
.model-preview {
    padding: 10px;
}

.page-header {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 20px;
}

.actions {
    margin-left: auto;
}

.filter-bar {
    margin-bottom: 20px;
    display: flex;
    gap: 10px;
}

.pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}
</style>

<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary">编辑数据源</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <el-form ref="dataSourceFormRef" :model="dataSourceForm" :rules="formRules" label-width="120px" class="edit-form" style="max-width: 900px; margin: 0 auto;">
                <el-divider content-position="left">基础信息</el-divider>
                <el-row :gutter="32">
                    <el-col :span="12">
                        <el-form-item label="数据源名称" prop="conn_name">
                            <el-input v-model="dataSourceForm.conn_name" placeholder="请输入名称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="数据源类型" prop="conn_kind">
                            <el-select v-model="dataSourceForm.conn_kind" class="w-full">
                                <el-option label="MySQL" value="MySQL" />
                                <el-option label="PostgreSQL" value="PostgreSQL" />
                                <el-option label="SQL Server" value="SQL Server" />
                                <el-option label="Oracle" value="Oracle" />
                            </el-select>
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-divider content-position="left">连接配置</el-divider>
                <el-row :gutter="32">
                    <el-col :span="12">
                        <el-form-item label="主机地址" prop="conn_host">
                            <el-input v-model="dataSourceForm.conn_host" placeholder="例如: localhost 或 IP" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="数据库名" prop="conn_database">
                            <el-input v-model="dataSourceForm.conn_database" placeholder="Database Name" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row :gutter="32">
                    <el-col :span="12">
                        <el-form-item label="端口" prop="conn_port">
                            <el-input-number v-model="dataSourceForm.conn_port" :min="1" :max="65535" class="w-full" controls-position="right" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="版本" prop="conn_version">
                            <el-input v-model="dataSourceForm.conn_version" placeholder="如: 8.0" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-divider content-position="left">身份验证</el-divider>
                <el-row :gutter="32">
                    <el-col :span="12">
                        <el-form-item label="用户名" prop="conn_user">
                            <el-input v-model="dataSourceForm.conn_user" placeholder="Username" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="密码" prop="conn_password">
                            <el-input v-model="dataSourceForm.conn_password" type="password" placeholder="Password" show-password clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-form-item label="备注说明" prop="remark">
                    <el-input v-model="dataSourceForm.remark" type="textarea" :rows="4" placeholder="补充详细说明..." />
                </el-form-item>
                <div class="m-t-lg flex-center" style="gap: 16px;">
                    <el-button @click="handleCancel" style="width: 120px;">取消</el-button>
                    <el-button type="warning" @click="handleTestConnection" :loading="testingConnection" style="width: 120px;">
                        测试连接
                    </el-button>
                    <el-button type="primary" @click="handleSubmit" :loading="submitting" style="width: 120px;">
                        更新保存
                    </el-button>
                </div>
            </el-form>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { getConnById, testConn, updateConn } from '@/api/metadata'
import type { MdConn } from '@/types/metadata'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const dataSourceFormRef = ref()
const testingConnection = ref(false)
const submitting = ref(false)

const dataSourceForm = reactive<Partial<MdConn>>({
    conn_name: '',
    conn_kind: 'MySQL',
    conn_version: '',
    conn_host: '',
    conn_port: 3306,
    conn_user: '',
    conn_password: '',
    conn_database: '',
    remark: ''
})

const formRules = {
    conn_name: [{ required: true, message: '名不能为空', trigger: 'blur' }],
    conn_kind: [{ required: true, message: '必选', trigger: 'change' }],
    conn_host: [{ required: true, message: '地址必填', trigger: 'blur' }],
    conn_port: [{ required: true, message: '端口必填', trigger: 'blur' }],
    conn_user: [{ required: true, message: '用户必填', trigger: 'blur' }],
    conn_database: [{ required: true, message: '库名必填', trigger: 'blur' }]
}

onMounted(() => {
    if (route.params.id) {
        fetchData(route.params.id as string)
    }
})

const fetchData = async (id: string) => {
    try {
        const response: any = await getConnById(id)
        const data = response?.data || response
        Object.assign(dataSourceForm, data)
    } catch (err) {
        console.error('获取数据源详情失败:', err)
        ElMessage.error('获取详情失败')
    }
}

const handleCancel = () => router.push('/data-sources')

const handleTestConnection = async () => {
    if (!route.params.id) {
        ElMessage.warning('请先保存数据源')
        return
    }
    testingConnection.value = true
    try {
        await testConn(route.params.id as string)
        ElMessage.success('连接测试成功')
    } catch (error: any) {
        console.error('连接测试失败:', error)
        ElMessage.error(error?.message || '连接测试失败')
    } finally {
        testingConnection.value = false
    }
}

const handleSubmit = async () => {
    try {
        await dataSourceFormRef.value.validate()
        submitting.value = true
        await updateConn(route.params.id as string, dataSourceForm)
        ElMessage.success('更新成功')
        router.push('/data-sources')
    } catch (error: any) {
        console.error('更新数据源失败:', error)
        ElMessage.error(error?.message || '更新失败')
    } finally {
        submitting.value = false
    }
}
</script>
<style scoped>
/* 样式通过全局 base.css 和 components.css 驱动 */
</style>

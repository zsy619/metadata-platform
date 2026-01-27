<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary">创建数据源</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <el-form ref="dataSourceFormRef" :model="dataSourceForm" :rules="formRules" label-width="120px" class="create-form" style="max-width: 900px; margin: 0 auto;">
                <el-divider content-position="left">基础信息</el-divider>
                <el-row :gutter="32">
                    <el-col :span="12">
                        <el-form-item label="数据源名称" prop="conn_name">
                            <el-input v-model="dataSourceForm.conn_name" placeholder="请输入数据源名称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="数据源类型" prop="conn_kind">
                            <el-select v-model="dataSourceForm.conn_kind" placeholder="请选择类型" class="w-full" @change="handleKindChange">
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
                            <el-input v-model="dataSourceForm.conn_database" placeholder="数据库名称" clearable />
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
                <el-form-item label="备注说明">
                    <el-input v-model="dataSourceForm.remark" type="textarea" placeholder="补充详细说明..." :rows="4" clearable />
                </el-form-item>
                <div class="m-t-lg flex-center" style="gap: 16px;">
                    <el-button @click="handleCancel" style="width: 120px;">取消</el-button>
                    <el-button type="warning" @click="handleTestConnection" :loading="testingConnection" style="width: 120px;">
                        测试连接
                    </el-button>
                    <el-button type="primary" @click="handleSubmit" :loading="submitting" style="width: 120px;">
                        立即创建
                    </el-button>
                </div>
            </el-form>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { createConn, testRawConn } from '@/api/metadata'
import type { MdConn } from '@/types/metadata'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const dataSourceFormRef = ref()

// 响应式数据
const testingConnection = ref(false)
const submitting = ref(false)

const dataSourceForm = reactive<Partial<MdConn>>({
    conn_name: '',
    conn_kind: 'MySQL',
    conn_version: '8.0',
    conn_host: 'localhost',
    conn_port: 3306,
    conn_user: 'root',
    conn_password: '',
    conn_database: '',
    remark: ''
})

// 默认配置映射
const defaultConfigs: Record<string, any> = {
    'MySQL': { port: 3306, user: 'root', version: '8.0' },
    'PostgreSQL': { port: 5432, user: 'postgres', version: '16' },
    'SQL Server': { port: 1433, user: 'sa', version: '2022' },
    'Oracle': { port: 1521, user: 'system', version: '19c' },
    'TiDB': { port: 4000, user: 'root', version: 'v7.5' },
    'OceanBase': { port: 2881, user: 'root@sys', version: '4.0' },
    'SQLite': { port: 0, user: '', version: '3.0' },
    'ClickHouse': { port: 8123, user: 'default', version: '24.1' },
    'Doris': { port: 9030, user: 'root', version: '2.0' },
    'StarRocks': { port: 9030, user: 'root', version: '3.1' },
    'DM': { port: 5236, user: 'SYSDBA', version: '8.0' },
    'Kingbase': { port: 54321, user: 'SYSTEM', version: 'V8' },
    'OpenGauss': { port: 5432, user: 'omm', version: '5.0' },
    'MongoDB': { port: 27017, user: '', version: '7.0' },
    'Redis': { port: 6379, user: '', version: '7.0' }
}

const handleKindChange = (kind: string) => {
    const config = defaultConfigs[kind]
    if (config) {
        dataSourceForm.conn_port = config.port
        dataSourceForm.conn_user = config.user
        dataSourceForm.conn_version = config.version
    }
}

// 表单验证规则
const formRules = reactive({
    conn_name: [
        { required: true, message: '请输入数据源名称', trigger: 'blur' },
        { min: 2, max: 100, message: '数据源名称长度在 2 到 100 个字符', trigger: 'blur' }
    ],
    conn_kind: [
        { required: true, message: '请选择数据源类型', trigger: 'change' }
    ],
    conn_host: [
        { required: true, message: '请输入主机地址', trigger: 'blur' }
    ],
    conn_port: [
        { required: true, message: '请输入端口号', trigger: 'blur' },
        { type: 'number', message: '端口号必须是数字', trigger: 'blur' }
    ],
    conn_user: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    conn_database: [
        { required: true, message: '请输入数据库名', trigger: 'blur' }
    ]
})

// 返回列表
const handleCancel = () => {
    router.push('/data-sources')
}

// 测试连接
const handleTestConnection = async () => {
    try {
        await dataSourceFormRef.value.validate()
        testingConnection.value = true

        const res = await testRawConn(dataSourceForm)
        if (res) {
            ElMessage.success('连接测试成功')
        }
    } catch (error: any) {
        console.error('连接测试失败:', error)
        ElMessage.error(error?.message || '连接测试失败')
    } finally {
        testingConnection.value = false
    }
}

// 表单提交
const handleSubmit = async () => {
    try {
        await dataSourceFormRef.value.validate()
        submitting.value = true

        await createConn(dataSourceForm)
        ElMessage.success('创建成功')
        router.push('/data-sources')
    } catch (error: any) {
        console.error('提交失败:', error)
        ElMessage.error(error?.message || '提交失败')
    } finally {
        submitting.value = false
    }
}

// 验证字段列表
// const validateFields = ['connName', 'connKind', 'connHost', 'connPort', 'connUser', 'connDatabase']
</script>
<style scoped>
.data-source-create {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.create-form {
    padding: 0 20px;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 30px;
    gap: 10px;
}

.w-full {
    width: 100%;
}
</style>
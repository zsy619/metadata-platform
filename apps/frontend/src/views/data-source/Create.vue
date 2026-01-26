<template>
    <div class="data-source-create">
        <div class="page-header">
            <h1>创建数据源</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card class="create-card">
            <el-form ref="dataSourceFormRef" :model="dataSourceForm" :rules="formRules" label-width="120px" class="create-form">
                <el-row :gutter="20">
                    <el-col :span="12">
                        <el-form-item label="数据源名称" prop="connName">
                            <el-input v-model="dataSourceForm.connName" placeholder="请输入数据源名称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="数据源类型" prop="connKind">
                            <el-select v-model="dataSourceForm.connKind" placeholder="请选择数据源类型" class="w-full" @change="handleKindChange">
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
                <el-row :gutter="20">
                    <el-col :span="12">
                        <el-form-item label="数据库版本" prop="connVersion">
                            <el-input v-model="dataSourceForm.connVersion" placeholder="请输入数据库版本，如：8.0" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="主机地址" prop="connHost">
                            <el-input v-model="dataSourceForm.connHost" placeholder="请输入主机地址，如：localhost" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row :gutter="20">
                    <el-col :span="12">
                        <el-form-item label="端口号" prop="connPort">
                            <el-input-number v-model="dataSourceForm.connPort" :min="1" :max="65535" placeholder="请输入端口号" class="w-full" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="数据库名" prop="connDatabase">
                            <el-input v-model="dataSourceForm.connDatabase" placeholder="请输入数据库名称" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-row :gutter="20">
                    <el-col :span="12">
                        <el-form-item label="用户名" prop="connUser">
                            <el-input v-model="dataSourceForm.connUser" placeholder="请输入数据库用户名" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="密码" prop="connPassword">
                            <el-input v-model="dataSourceForm.connPassword" type="password" placeholder="请输入数据库密码" clearable show-password />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-form-item label="备注">
                    <el-input v-model="dataSourceForm.remark" type="textarea" placeholder="请输入备注信息" :rows="4" clearable />
                </el-form-item>
                <div class="form-actions">
                    <el-button @click="handleCancel">取消</el-button>
                    <el-button type="primary" @click="handleTestConnection" :loading="testingConnection">
                        {{ testingConnection ? '测试连接中...' : '测试连接' }}
                    </el-button>
                    <el-button type="success" @click="handleSubmit" :loading="submitting">
                        {{ submitting ? '提交中...' : '提交' }}
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
    connName: '',
    connKind: 'MySQL',
    connVersion: '8.0',
    connHost: 'localhost',
    connPort: 3306,
    connUser: 'root',
    connPassword: '',
    connDatabase: '',
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
        dataSourceForm.connPort = config.port
        dataSourceForm.connUser = config.user
        dataSourceForm.connVersion = config.version
    }
}

// 表单验证规则
const formRules = reactive({
    connName: [
        { required: true, message: '请输入数据源名称', trigger: 'blur' },
        { min: 2, max: 100, message: '数据源名称长度在 2 到 100 个字符', trigger: 'blur' }
    ],
    connKind: [
        { required: true, message: '请选择数据源类型', trigger: 'change' }
    ],
    connHost: [
        { required: true, message: '请输入主机地址', trigger: 'blur' }
    ],
    connPort: [
        { required: true, message: '请输入端口号', trigger: 'blur' },
        { type: 'number', message: '端口号必须是数字', trigger: 'blur' }
    ],
    connUser: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    connDatabase: [
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
    padding: 20px;
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
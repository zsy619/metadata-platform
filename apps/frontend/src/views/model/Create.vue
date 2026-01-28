<template>
    <div class="model-create">
        <div class="page-header">
            <h1>创建模型</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <el-steps :active="activeStep" space="300px">
                <el-step title="选择数据源" />
                <el-step title="选择表/视图" />
                <el-step title="配置模型信息" />
                <el-step title="完成" />
            </el-steps>
            <!-- 第一步：基础配置 -->
            <div v-if="activeStep === 0" class="step-content">
                <h3>选择数据源与模型类型</h3>
                <el-form :model="modelForm" label-width="120px">
                    <el-form-item label="数据源" prop="connID">
                        <el-select v-model="modelForm.connID" placeholder="请选择数据源" style="width: 100%" @change="handleDataSourceChange">
                            <el-option v-for="conn in dataSources" :key="conn.id" :label="conn.connName" :value="conn.id">
                                <div class="option-content">
                                    <span>{{ conn.connName }}</span>
                                    <span class="option-subtitle">{{ conn.connKind }} - {{ conn.connHost }}:{{ conn.connPort }}</span>
                                </div>
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="模型类型" prop="modelKind">
                        <el-select v-model="modelForm.modelKind" placeholder="请选择模型类型" style="width: 100%" @change="handleKindChange">
                            <el-option label="SQL语句" :value="1" />
                            <el-option label="视图/表" :value="2" />
                            <el-option label="存储过程" :value="3" />
                            <el-option label="关联" :value="4" />
                        </el-select>
                    </el-form-item>
                </el-form>
                <div class="step-actions">
                    <el-button type="primary" @click="nextStep" :disabled="!modelForm.connID || !modelForm.modelKind">
                        下一步
                    </el-button>
                </div>
            </div>
            <!-- 第二步：类型特定配置 -->
            <div v-if="activeStep === 1" class="step-content">
                <!-- 视图/表 (Kind 2) -->
                <template v-if="modelForm.modelKind === 2">
                    <h3>选择表/视图</h3>
                    <el-form-item label="对象类型">
                        <el-radio-group v-model="objectType" @change="handleObjectTypeChange">
                            <el-radio-button label="table">表</el-radio-button>
                            <el-radio-button label="view">视图</el-radio-button>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item :label="`选择${objectType === 'table' ? '表' : '视图'}`">
                        <el-tree v-loading="loadingTables" :data="tables" :props="tableTreeProps" node-key="id" ref="tableTree" @node-click="handleTableClick" style="max-height: 400px; overflow-y: auto">
                            <template #default="{ node, data }">
                                <span class="custom-tree-node">
                                    <span>{{ data.label }}</span>
                                    <span class="node-subtitle">{{ data.comment || '无描述' }}</span>
                                </span>
                            </template>
                        </el-tree>
                    </el-form-item>
                </template>
                <!-- SQL语句 (Kind 1) -->
                <template v-else-if="modelForm.modelKind === 1">
                    <h3>输入 SQL 语句</h3>
                    <el-form-item label="SQL 脚本">
                        <el-input v-model="modelForm.sqlContent" type="textarea" :rows="10" placeholder="请输入 SELECT SQL 语句" />
                    </el-form-item>
                </template>
                <!-- 存储过程 (Kind 3) -->
                <template v-else-if="modelForm.modelKind === 3">
                    <h3>选择存储过程</h3>
                    <el-form-item label="存储过程">
                        <el-select v-model="modelForm.spName" placeholder="请选择存储过程" style="width: 100%">
                            <el-option v-for="sp in storedProcedures" :key="sp.name" :label="sp.name" :value="sp.name" />
                        </el-select>
                    </el-form-item>
                </template>
                <!-- 关联 (Kind 4) -->
                <template v-else-if="modelForm.modelKind === 4">
                    <h3>配置关联关系 (选择主表)</h3>
                    <el-form-item label="主表">
                        <el-select v-model="modelForm.mainTable" placeholder="请选择主表" style="width: 100%">
                            <el-option v-for="t in tables" :key="t.id" :label="t.label" :value="t.label" />
                        </el-select>
                    </el-form-item>
                    <el-alert title="关联详细配置将在下一步或之后版本完善" type="info" show-icon :closable="false" style="margin-bottom: 20px" />
                </template>
                <div class="step-actions">
                    <el-button @click="prevStep">上一步</el-button>
                    <el-button type="primary" @click="nextStep" :disabled="!canProceedToInfo">
                        下一步
                    </el-button>
                </div>
            </div>
            <!-- 第三步：配置模型信息 -->
            <div v-if="activeStep === 2" class="step-content">
                <h3>配置模型信息</h3>
                <el-form ref="modelFormRef" :model="modelForm" :rules="formRules" label-width="120px" class="create-form">
                    <el-row :gutter="20">
                        <el-col :span="12">
                            <el-form-item label="模型名称" prop="modelName">
                                <el-input v-model="modelForm.modelName" placeholder="请输入模型名称" clearable />
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="模型编码" prop="modelCode">
                                <el-input v-model="modelForm.modelCode" placeholder="请输入模型编码（字母数字下划线组合）" clearable />
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row :gutter="20">
                        <el-col :span="12">
                            <el-form-item label="模型版本" prop="modelVersion">
                                <el-input v-model="modelForm.modelVersion" placeholder="请输入模型版本，如：1.0.0" clearable />
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="模型类型" prop="modelKind">
                                <el-select v-model="modelForm.modelKind" placeholder="请选择模型类型">
                                    <el-option label="SQL语句" :value="1" />
                                    <el-option label="视图/表" :value="2" />
                                    <el-option label="存储过程" :value="3" />
                                    <el-option label="关联" :value="4" />
                                </el-select>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-form-item label="是否公开">
                        <el-switch v-model="modelForm.isPublic" />
                    </el-form-item>
                    <el-form-item label="备注">
                        <el-input v-model="modelForm.remark" type="textarea" placeholder="请输入模型描述信息" :rows="4" clearable />
                    </el-form-item>
                    <div class="step-actions">
                        <el-button @click="prevStep">上一步</el-button>
                        <el-button type="primary" @click="handleSubmit">
                            提交
                        </el-button>
                    </div>
                </el-form>
            </div>
            <!-- 第四步：完成 -->
            <div v-if="activeStep === 3" class="step-content success-step">
                <el-result icon="success" title="模型创建成功" sub-title="您已成功创建了一个新的业务数据模型">
                    <template #extra>
                        <el-button type="primary" @click="handleFinish">
                            查看模型列表
                        </el-button>
                        <el-button @click="handleCreateAnother">
                            创建另一个模型
                        </el-button>
                    </template>
                </el-result>
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import type { Model } from '@/types/metadata'
import type { DataSource } from '@/types/metadata/datasource'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const modelFormRef = ref()
const tableTree = ref()

// 响应式数据
const activeStep = ref(0)
const loading = ref(false)
const loadingTables = ref(false)
const objectType = ref('table')
const storedProcedures = ref<any[]>([])
const modelForm = reactive<Partial<Model & {
    connName?: string,
    sqlContent?: string,
    spName?: string,
    mainTable?: string
}>>({
    connID: undefined,
    connName: '',
    modelName: '',
    modelCode: '',
    modelVersion: '1.0.0',
    modelKind: 2,
    isPublic: false,
    remark: '',
    sqlContent: '',
    spName: '',
    mainTable: ''
})

const dataSources = ref<DataSource[]>([])
const tables = ref<any[]>([])
const selectedTable = ref<any>(null)

// 表树配置
const tableTreeProps = {
    label: 'label',
    children: 'children',
    isLeaf: 'isLeaf'
}

// 表单验证规则
const formRules = reactive({
    connID: [
        { required: true, message: '请选择数据源', trigger: 'change' }
    ],
    modelName: [
        { required: true, message: '请输入模型名称', trigger: 'blur' },
        { min: 2, max: 100, message: '模型名称长度在 2 到 100 个字符', trigger: 'blur' }
    ],
    modelCode: [
        { required: true, message: '请输入模型编码', trigger: 'blur' },
        { pattern: /^[a-zA-Z0-9_]+$/, message: '模型编码只能包含字母、数字和下划线', trigger: 'blur' },
        { min: 2, max: 50, message: '模型编码长度在 2 到 50 个字符', trigger: 'blur' }
    ],
    modelVersion: [
        { required: true, message: '请输入模型版本', trigger: 'blur' }
    ],
    modelKind: [
        { required: true, message: '请选择模型类型', trigger: 'change' }
    ]
})

// 计算属性：是否可以进行下一步进入信息配置
const canProceedToInfo = computed(() => {
    if (modelForm.modelKind === 1) return !!modelForm.sqlContent
    if (modelForm.modelKind === 2) return !!selectedTable.value
    if (modelForm.modelKind === 3) return !!modelForm.spName
    if (modelForm.modelKind === 4) return !!modelForm.mainTable
    return false
})

// 生命周期钩子
const route = useRoute()

onMounted(() => {
    fetchDataSources()

    // 处理 URL 参数中的 kind
    const queryKind = route.query.kind
    if (queryKind) {
        modelForm.modelKind = Number(queryKind)
        handleKindChange(modelForm.modelKind)
    }
})

// 获取数据源列表
const fetchDataSources = async () => {
    try {
        // 模拟API调用，实际需要替换为真实API
        // const response = await getDataSources()
        // dataSources.value = response.data

        // 模拟数据
        dataSources.value = [
            { id: '1', parentID: '0', connName: '测试MySQL数据源', connKind: 'MySQL', connVersion: '8.0', connHost: 'localhost', connPort: 3306, connUser: 'root', connPassword: 'password', connDatabase: 'test_db', connConn: '', isDeleted: false, state: 1, remark: '', sort: 0, createdAt: '', updatedAt: '', tenantID: '0' }
        ] as any
    } catch (error) {
        console.error('获取数据源列表失败:', error)
        ElMessage.error('获取数据源列表失败')
    }
}

// 获取表/视图列表
const fetchTables = async () => {
    if (!modelForm.connID) return

    loadingTables.value = true
    try {
        // 模拟API调用，实际需要替换为真实API
        // const response = await getDataSourceTables(modelForm.connID, objectType.value)
        // tables.value = response.data

        // 模拟数据
        tables.value = [
            {
                id: 1,
                label: 'user',
                comment: '用户表',
                isLeaf: true
            },
            {
                id: 2,
                label: 'order',
                comment: '订单表',
                isLeaf: true
            },
            {
                id: 3,
                label: 'product',
                comment: '产品表',
                isLeaf: true
            }
        ]
    } catch (error) {
        console.error('获取表列表失败:', error)
        ElMessage.error('获取表列表失败')
    } finally {
        loadingTables.value = false
    }
}

// 数据源变化
const handleDataSourceChange = (connID: number) => {
    const conn = dataSources.value.find(item => item.id === connID)
    if (conn) {
        modelForm.connName = conn.connName
    }
    // 重置表选择
    selectedTable.value = null
    tables.value = []
    // 如果是表/视图或关联，获取表列表
    if (modelForm.modelKind === 2 || modelForm.modelKind === 4) {
        fetchTables()
    } else if (modelForm.modelKind === 3) {
        fetchStoredProcedures()
    }
}

// 类型变化
const handleKindChange = (kind: number) => {
    // 重置之前的选择
    selectedTable.value = null
    modelForm.sqlContent = ''
    modelForm.spName = ''
    modelForm.mainTable = ''

    if (kind === 2 || kind === 4) {
        fetchTables()
    } else if (kind === 3) {
        fetchStoredProcedures()
    }
}

// 获取存储过程
const fetchStoredProcedures = async () => {
    if (!modelForm.connID) return
    loadingTables.value = true
    try {
        // 模拟数据
        storedProcedures.value = [
            { name: 'sp_get_user_stats' },
            { name: 'sp_sync_orders' }
        ]
    } finally {
        loadingTables.value = false
    }
}

// 对象类型变化
const handleObjectTypeChange = () => {
    // 重置表选择
    selectedTable.value = null
    tables.value = []
    // 获取表列表
    fetchTables()
}

// 表点击
const handleTableClick = (data: any) => {
    selectedTable.value = data
}

// 下一步
const nextStep = () => {
    activeStep.value++
}

// 上一步
const prevStep = () => {
    activeStep.value--
}

// 表单提交
const handleSubmit = async () => {
    try {
        await modelFormRef.value.validate()
        loading.value = true

        // 模拟创建模型
        // let result
        // if (modelForm.modelKind === 1) {
        //   result = await buildModelFromSQL({ ... })
        // } else if (modelForm.modelKind === 2) {
        //   ... buildModelFromTable/View ...
        // }

        // 跳转到完成步骤
        nextStep()
    } catch (error) {
        console.error('创建模型失败:', error)
        ElMessage.error('创建模型失败')
    } finally {
        loading.value = false
    }
}

const handleCancel = () => {
    router.push('/metadata/model/list')
}

// 完成后查看模型列表
const handleFinish = () => {
    router.push('/metadata/model/list')
}

// 创建另一个模型
const handleCreateAnother = () => {
    // 重置表单
    activeStep.value = 0
    selectedTable.value = null
    Object.assign(modelForm, {
        connID: undefined,
        connName: '',
        modelName: '',
        modelCode: '',
        modelVersion: '1.0.0',
        modelKind: 2,
        isPublic: false,
        remark: ''
    })
}
</script>
<style scoped>
.model-create {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.step-content {
    margin-top: 40px;
}

.step-content h3 {
    margin-bottom: 20px;
    font-size: 18px;
    font-weight: 600;
}

.step-actions {
    margin-top: 40px;
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

.option-content {
    display: flex;
    flex-direction: column;
}

.option-subtitle {
    font-size: 12px;
    color: #909399;
}

.custom-tree-node {
    display: flex;
    flex-direction: column;
}

.node-subtitle {
    font-size: 12px;
    color: #909399;
}

.success-step {
    text-align: center;
    padding: 40px 0;
}
</style>

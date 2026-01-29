<template>
    <div class="create-sql-container">
        <div class="page-header">
            <h1>创建 SQL 模型</h1>
            <el-button @click="goBack" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <div class="content-body">
                <!-- 步骤条 -->
                <el-steps :active="activeStep" finish-status="success" simple class="mb-4">
                    <el-step title="基础信息" />
                    <el-step title="编写 SQL" />
                    <el-step title="参数配置" />
                    <el-step title="字段映射" />
                    <el-step title="权限配置" />
                </el-steps>
                <!-- 步骤 1: 基础信息 -->
                <div v-show="activeStep === 0" class="step-content">
                    <el-form ref="baseFormRef" :model="baseForm" :rules="baseRules" label-width="100px" class="max-w-800">
                        <el-form-item label="数&#8194;据&#8194;源" prop="connID">
                            <el-select v-model="baseForm.connID" placeholder="请选择数据源" class="w-full">
                                <el-option v-for="item in dataSources" :key="item.id" :label="item.conn_name" :value="item.id" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="模型名称" prop="modelName">
                            <el-input v-model="baseForm.modelName" placeholder="请输入模型名称" />
                        </el-form-item>
                        <el-form-item label="模型编码" prop="modelCode">
                            <el-input v-model="baseForm.modelCode" placeholder="自动生成或手动输入" />
                        </el-form-item>
                        <el-form-item label="描&#12288;&#12288;述" prop="remark">
                            <el-input v-model="baseForm.remark" type="textarea" :rows="3" placeholder="请输入模型描述" />
                        </el-form-item>
                    </el-form>
                </div>
                <!-- 步骤 2: SQL 编辑 -->
                <div v-show="activeStep === 1" class="step-content sql-step">
                    <el-alert title="提示：使用 :parameter_name 格式定义参数，例如 WHERE id = :user_id" type="info" show-icon class="mb-3" />
                    <div class="sql-editor-wrapper">
                        <Codemirror v-model="sqlContent" placeholder="请输入 SQL 查询语句..." :style="{ height: '400px' }" :autofocus="true" :indent-with-tab="true" :tab-size="2" :extensions="extensions" @change="detectParameters" />
                    </div>
                </div>
                <!-- 步骤 3: 参数配置 -->
                <div v-show="activeStep === 2" class="step-content">
                    <el-table :data="parameters" border style="width: 100%">
                        <el-table-column prop="name" label="参数名" width="180" />
                        <el-table-column label="数据类型" width="150">
                            <template #default="{ row }">
                                <el-select v-model="row.type" size="small">
                                    <el-option label="String" value="string" />
                                    <el-option label="Int" value="int" />
                                    <el-option label="Number" value="number" />
                                    <el-option label="Date" value="date" />
                                    <el-option label="Boolean" value="boolean" />
                                </el-select>
                            </template>
                        </el-table-column>
                        <el-table-column label="必填" width="100" align="center">
                            <template #default="{ row }">
                                <el-switch v-model="row.required" size="small" />
                            </template>
                        </el-table-column>
                        <el-table-column label="默认值">
                            <template #default="{ row }">
                                <el-input v-model="row.default" size="small" placeholder="无默认值" />
                            </template>
                        </el-table-column>
                    </el-table>
                    <div class="mt-3">
                        <el-button size="small" @click="detectParameters">重新识别参数</el-button>
                    </div>
                </div>
                <!-- 步骤 4: 字段映射 -->
                <div v-show="activeStep === 3" class="step-content max-w-1200">
                    <div class="mb-3 flex justify-between items-center">
                        <span>预览并配置结果集字段</span>
                        <el-button type="primary" link @click="fetchColumns">刷新字段</el-button>
                    </div>
                    <el-table :data="fieldMappings" border style="width: 100%">
                        <el-table-column prop="column_name" label="列名" width="180" />
                        <el-table-column label="显示名称">
                            <template #default="{ row }">
                                <el-input v-model="row.show_title" size="small" />
                            </template>
                        </el-table-column>
                        <el-table-column label="显示宽度" width="120">
                            <template #default="{ row }">
                                <el-input-number v-model="row.show_width" size="small" :min="0" :step="10" />
                            </template>
                        </el-table-column>
                        <el-table-column label="格式化" width="150">
                            <template #default="{ row }">
                                <el-select v-model="row.format" size="small" clearable placeholder="选择格式">
                                    <el-option label="无" value="" />
                                    <el-option label="千分位" value="number" />
                                    <el-option label="日期" value="date" />
                                    <el-option label="金额" value="currency" />
                                </el-select>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 步骤 5: 权限配置 -->
                <div v-show="activeStep === 4" class="step-content">
                    <el-form label-width="120px" class="max-w-800">
                        <el-form-item label="公开模型">
                            <el-switch v-model="permissions.isPublic" />
                            <span class="ml-2 text-gray-400 text-sm">开启后所有用户可见</span>
                        </el-form-item>
                        <el-form-item label="备&#12288;&#12288;注">
                            <el-input v-model="baseForm.remark" type="textarea" :rows="3" placeholder="请输入模型描述" />
                        </el-form-item>
                    </el-form>
                </div>
                <!-- 底部操作栏 -->
                <div class="footer-actions mt-6 flex justify-end">
                    <el-button @click="goBack">取消</el-button>
                    <el-button v-if="activeStep > 0" @click="handlePrev">上一步</el-button>
                    <el-button v-if="activeStep < 4" type="primary" @click="handleNext">下一步</el-button>
                    <el-button v-if="activeStep === 4" type="primary" :loading="submitting" @click="handleSubmit">
                        完成创建
                    </el-button>
                </div>
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { getConns, getDBTables } from '@/api/metadata'
import { createModelSql, generateModelCode, testSQL } from '@/api/model'
import type { FieldMapping, SQLParameter } from '@/types/metadata/model-params'
import { sql } from '@codemirror/lang-sql'
import { oneDark } from '@codemirror/theme-one-dark'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { Codemirror } from 'vue-codemirror'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeStep = ref(0)
const submitting = ref(false)

// 数据库元数据用于自动补全
const dbTables = ref<string[]>([])

// CodeMirror Extensions
const extensions = computed(() => {
    const schema: Record<string, string[]> = {}
    dbTables.value.forEach(table => {
        schema[table] = [] // 目前仅补全表名，列名补全需要额外 API
    })
    return [sql({ schema }), oneDark]
})

// 基础信息
const baseFormRef = ref<FormInstance>()
const dataSources = ref<any[]>([])
const baseForm = reactive({
    connID: '',
    modelName: '',
    modelCode: '',
    remark: ''
})
const baseRules = {
    connID: [{ required: true, message: '请选择数据源', trigger: 'change' }],
    modelName: [{ required: true, message: '请输入模型名称', trigger: 'blur' }],
    modelCode: [{ required: true, message: '请输入模型编码', trigger: 'blur' }]
}

// SQL 内容
const sqlContent = ref('')

// 参数列表
const parameters = ref<SQLParameter[]>([])

// 字段映射
const fieldMappings = ref<FieldMapping[]>([])

// 权限
const permissions = reactive({
    isPublic: false
})

onMounted(async () => {
    // 自动获取模型编码
    fetchGeneratedCode()

    try {
        const res: any = await getConns()
        dataSources.value = Array.isArray(res) ? res : (res.data || [])
    } catch (error) {
        console.error('Failed to load data sources', error)
    }
})

const fetchGeneratedCode = async () => {
    try {
        const res: any = await generateModelCode()
        if (res.data && res.data.code) {
            baseForm.modelCode = res.data.code
        }
    } catch (error) {
        console.error('Failed to auto generate model code', error)
    }
}

const goBack = () => {
    router.push('/metadata/model/list')
}

const handlePrev = () => {
    if (activeStep.value > 0) {
        activeStep.value--
    }
}

const handleNext = async () => {
    if (activeStep.value === 0) {
        if (!baseFormRef.value) return
        await baseFormRef.value.validate((valid) => {
            if (valid) activeStep.value++
        })
    } else if (activeStep.value === 1) {
        if (!sqlContent.value.trim()) {
            ElMessage.warning('请输入 SQL 语句')
            return
        }
        detectParameters()
        activeStep.value++
    } else if (activeStep.value === 2) {
        if (fieldMappings.value.length === 0) {
            await fetchColumns()
        }
        activeStep.value++
    } else {
        activeStep.value++
    }
}

// 获取库表元数据用于补全
const fetchDbMetadata = async () => {
    if (!baseForm.connID) return
    try {
        const res: any = await getDBTables(baseForm.connID)
        // 后端返回成功结构是 { code: 200, data: [...] }
        const list = res.data || []
        dbTables.value = list.map((t: any) => typeof t === 'string' ? t : (t.name || t.label || t.TableName))
    } catch (error) {
        console.error('Failed to fetch DB metadata for autocomplete', error)
    }
}

// 监听步骤变化，进入 SQL 编辑步时获取元数据
watch(activeStep, (newStep) => {
    if (newStep === 1) {
        fetchDbMetadata()
    }
})

// 简单的参数识别逻辑
const detectParameters = () => {
    const regex = /:([a-zA-Z_][a-zA-Z0-9_]*)/g
    const foundParams = new Set<string>()
    let match
    while ((match = regex.exec(sqlContent.value)) !== null) {
        foundParams.add(match[1])
    }

    // 更新参数列表，保留已配置的
    const newParams: SQLParameter[] = []
    foundParams.forEach(name => {
        const existing = parameters.value.find(p => p.name === name)
        if (existing) {
            newParams.push(existing)
        } else {
            newParams.push({
                name,
                type: 'string',
                required: true,
                default: ''
            })
        }
    })
    parameters.value = newParams
}

const fetchColumns = async () => {
    if (!sqlContent.value) return
    const loading = ElMessage.info({
        message: '正在分析 SQL 结果集...',
        duration: 0
    })

    try {
        const res = await testSQL({
            conn_id: baseForm.connID,
            sql_content: sqlContent.value,
            parameters: parameters.value.map(p => ({
                name: p.name,
                value: p.default
            }))
        })

        if (res && res.data && res.data.fields) {
            fieldMappings.value = res.data.fields.map((f: any) => ({
                column_name: f.column_name,
                show_title: f.show_title || f.column_name,
                show_width: f.show_width || 150,
                format: f.format || ''
            }))
            loading.close()
            ElMessage.success('字段解析成功')
        }
    } catch (error) {
        console.error(error)
        loading.close()
        ElMessage.error('字段解析失败，请检查 SQL 语句或参数')
    }
}

const handleSubmit = async () => {
    submitting.value = true
    try {
        const payload = {
            conn_id: baseForm.connID,
            model_name: baseForm.modelName,
            model_code: baseForm.modelCode,
            sql_content: sqlContent.value,
            parameters: parameters.value,
            field_mappings: fieldMappings.value,
            is_public: permissions.isPublic,
            remark: baseForm.remark
        }

        await createModelSql(payload)
        ElMessage.success('模型创建成功')
        router.push('/metadata/model/list')
    } catch (error) {
        console.error(error)
        ElMessage.error('创建失败')
    } finally {
        submitting.value = false
    }
}
</script>
<style scoped>
.create-sql-container {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.page-header h1 {
    font-size: 20px;
    font-weight: 600;
    margin: 0;
}

.content-body {
    padding: 20px 0;
}

.max-w-800 {
    max-width: 800px;
    margin: 0 auto;
}

.max-w-1200 {
    max-width: 1200px;
    margin: 0 auto;
}

.step-content {
    min-height: 400px;
    margin-bottom: 24px;
    padding-top: 20px;
}

.sql-editor-wrapper {
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    overflow: hidden;
}

.footer-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 24px;
    gap: 12px;
}

/* CodeMirror customization if needed */
:deep(.cm-editor) {
    height: 400px;
    font-size: 14px;
    font-family: 'Fira Code', 'Menlo', monospace;
}
</style>

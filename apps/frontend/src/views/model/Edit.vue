<template>
    <div class="model-edit">
        <div class="page-header">
            <h1>编辑模型</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card v-loading="loading">
            <el-tabs v-model:active-name="activeTab" type="border-card">
                <!-- 基本信息 -->
                <el-tab-pane label="基本信息" name="basic">
                    <el-form ref="modelFormRef" :model="modelForm" :rules="formRules" label-width="120px" class="edit-form">
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="模型名称" prop="modelName">
                                    <el-input v-model="modelForm.modelName" placeholder="请输入模型名称" clearable />
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="模型编码" prop="modelCode">
                                    <el-input v-model="modelForm.modelCode" placeholder="请输入模型编码" clearable :disabled="modelForm.isLocked" />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="模型版本" prop="modelVersion">
                                    <el-input v-model="modelForm.modelVersion" placeholder="请输入模型版本" clearable />
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="模型类型" prop="modelKind">
                                    <el-select v-model="modelForm.modelKind" placeholder="请选择模型类型" :disabled="modelForm.isLocked">
                                        <el-option label="SQL语句" :value="1" />
                                        <el-option label="视图/表" :value="2" />
                                        <el-option label="存储过程" :value="3" />
                                        <el-option label="关联" :value="4" />
                                    </el-select>
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="数据源" prop="connID">
                                    <el-select v-model="modelForm.connID" placeholder="请选择数据源" disabled>
                                        <el-option v-for="conn in dataSources" :key="conn.id" :label="conn.connName" :value="conn.id" />
                                    </el-select>
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="数据源表">
                                    <el-input v-model="modelForm.tableName" placeholder="数据源表" disabled />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="是否公开">
                                    <el-switch v-model="modelForm.isPublic" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="是否锁定">
                                    <el-switch v-model="modelForm.isLocked" />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-form-item label="备注">
                            <el-input v-model="modelForm.remark" type="textarea" placeholder="请输入模型描述信息" :rows="4" clearable />
                        </el-form-item>
                        <div class="form-actions">
                            <el-button @click="handleCancel">取消</el-button>
                            <el-button type="primary" @click="handleSubmit" :loading="submitting">
                                提交
                            </el-button>
                        </div>
                    </el-form>
                </el-tab-pane>
                <!-- 字段配置 -->
                <el-tab-pane label="字段配置" name="fields">
                    <div class="field-config">
                        <el-table v-loading="loadingFields" :data="modelFields" border style="width: 100%">
                            <el-table-column prop="columnName" label="字段名称" width="180" />
                            <el-table-column prop="columnTitle" label="字段标题" width="180" />
                            <el-table-column prop="showTitle" label="显示名称" width="180">
                                <template #default="scope">
                                    <el-input v-model="scope.row.showTitle" size="small" :disabled="modelForm.isLocked" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="isShow" label="显示" width="100">
                                <template #default="scope">
                                    <el-switch v-model="scope.row.isShow" :active-value="1" :inactive-value="0" :disabled="modelForm.isLocked" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="sort" label="排序" width="120">
                                <template #default="scope">
                                    <el-input-number v-model="scope.row.sort" size="small" :min="0" :disabled="modelForm.isLocked" />
                                </template>
                            </el-table-column>
                            <el-table-column label="高级" width="150" fixed="right" class-name="action-column">
                                <template #default="scope">
                                    <el-button type="primary" size="small" :icon="Operation" @click="handleEditFieldEnhancement(scope.row)">
                                        增强配置
                                    </el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <div class="field-actions">
                            <el-button type="primary" @click="handleSaveFields" :loading="savingFields" :disabled="modelForm.isLocked">
                                保存字段配置
                            </el-button>
                        </div>
                    </div>
                </el-tab-pane>
                <!-- 树形配置 -->
                <el-tab-pane label="树形配置" name="tree">
                    <div class="tree-config">
                        <el-form :model="modelForm" label-width="140px" :disabled="modelForm.isLocked">
                            <el-form-item label="启用树形结构">
                                <el-switch v-model="modelForm.isTree" />
                            </el-form-item>
                            <template v-if="modelForm.isTree">
                                <el-form-item label="父级字段名">
                                    <el-select v-model="modelForm.treeParentField">
                                        <el-option v-for="f in modelFields" :key="f.id" :label="f.columnName" :value="f.columnName" />
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="路径字段名">
                                    <el-select v-model="modelForm.treePathField">
                                        <el-option v-for="f in modelFields" :key="f.id" :label="f.columnName" :value="f.columnName" />
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="层级字段名">
                                    <el-select v-model="modelForm.treeLevelField">
                                        <el-option v-for="f in modelFields" :key="f.id" :label="f.columnName" :value="f.columnName" />
                                    </el-select>
                                </el-form-item>
                            </template>
                            <div class="form-actions">
                                <el-button type="primary" @click="handleSubmit" :loading="submitting">保存树形配置</el-button>
                            </div>
                        </el-form>
                    </div>
                </el-tab-pane>
                <!-- 查询模板 -->
                <el-tab-pane label="查询模板" name="templates">
                    <div class="template-management">
                        <div class="action-header">
                            <el-button type="primary" :icon="Plus" @click="handleAddTemplate">新增模板</el-button>
                        </div>
                        <el-table :data="queryTemplates" border>
                            <el-table-column prop="templateName" label="模板名称" />
                            <el-table-column prop="templateCode" label="模板编码" />
                            <el-table-column prop="isDefault" label="默认" width="80">
                                <template #default="scope">
                                    <el-tag :type="scope.row.isDefault ? 'success' : 'info'">{{ scope.row.isDefault ? '是' : '否' }}</el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="350" fixed="right" class-name="action-column">
                                <template #default="scope">
                                    <el-button type="success" size="small" :icon="Edit" @click="handleEditTemplate(scope.row)">
                                        编辑
                                    </el-button>
                                    <el-button type="primary" size="small" :icon="CopyDocument" @click="handleDuplicateTemplate(scope.row)">
                                        复制
                                    </el-button>
                                    <el-button type="info" size="small" :icon="View" @click="handlePreviewTemplate(scope.row)">
                                        SQL
                                    </el-button>
                                    <el-button type="danger" size="small" :icon="Delete" @click="handleDeleteTemplate(scope.row)">
                                        删除
                                    </el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-tab-pane>
                <!-- 高级配置 -->
                <el-tab-pane label="高级配置" name="advanced">
                    <el-form :model="advancedConfig" label-width="140px">
                        <el-form-item label="启用审计">
                            <el-switch v-model="advancedConfig.enableAudit" />
                        </el-form-item>
                        <div class="form-actions">
                            <el-button type="primary" @click="handleSaveAdvanced">保存高级配置</el-button>
                        </div>
                    </el-form>
                </el-tab-pane>
            </el-tabs>
        </el-card>
        <!-- 增强配置对话框 -->
        <el-dialog v-model="enhancementVisible" title="字段增强配置" width="600px" class="custom-dialog">
            <el-form :model="currentField" label-width="120px">
                <el-form-item label="允许为空"><el-switch v-model="currentField.isNullable" /></el-form-item>
                <el-form-item label="主键"><el-switch v-model="currentField.isPrimaryKey" /></el-form-item>
                <el-form-item label="正则校验"><el-input v-model="currentField.fieldRegex" /></el-form-item>
                <!-- 其他增强字段... -->
            </el-form>
            <template #footer>
                <el-button @click="enhancementVisible = false">取消</el-button>
                <el-button type="primary" @click="saveEnhancement">保存</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { getConns as getConnsApi } from '@/api/metadata'
import { duplicateQueryTemplate, previewQueryTemplate } from '@/api/model'
import type { Model, ModelField } from '@/types/metadata'
import { ArrowLeft, CopyDocument, Delete, Edit, Operation, Plus, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const modelFormRef = ref()

const activeTab = ref('basic')
const loading = ref(false)
const submitting = ref(false)
const loadingFields = ref(false)
const savingFields = ref(false)
const enhancementVisible = ref(false)

const modelForm = reactive<Partial<Model & { tableName?: string }>>({
    isTree: false,
    treeParentField: 'parent_id',
    treePathField: 'tree_path',
    treeLevelField: 'tree_level'
})

const modelFields = ref<ModelField[]>([])
const queryTemplates = ref<any[]>([])
const dataSources = ref<any[]>([])
const advancedConfig = reactive({ enableAudit: true })
const currentField = ref<any>({})

const formRules = {
    modelName: [{ required: true, message: '请输入模型名称', trigger: 'blur' }],
    modelCode: [{ required: true, message: '请输入模型编码', trigger: 'blur' }]
}

onMounted(() => {
    const modelId = route.params.id as string
    if (modelId) {
        loadData(modelId)
    }
})

const loadData = async (id: string) => {
    loading.value = true
    try {
        // const model = await getModelById(Number(id))
        // Object.assign(modelForm, model)
        Object.assign(modelForm, {
            id: id,
            modelName: '用户模型',
            modelCode: 'user_model',
            modelVersion: '1.0.0',
            modelKind: 2,
            isTree: false,
            tableName: 'user'
        })

        loadingFields.value = true
        // const fields = await getModelFields(Number(id))
        // modelFields.value = fields
        modelFields.value = [{
            id: '1',
            columnName: 'id',
            columnTitle: 'ID',
            showTitle: 'ID',
            isShow: 1,
            sort: 0,
            isNullable: false,
            isPrimaryKey: true
        } as any]
        loadingFields.value = false

        const conns: any = await getConnsApi()
        dataSources.value = conns?.data || []
    } catch (e) {
        ElMessage.error('加载数据失败')
    } finally {
        loading.value = false
    }
}

const handleCancel = () => router.push('/model/list')

const handleSubmit = async () => {
    submitting.value = true
    try {
        ElMessage.success('保存成功')
    } finally {
        submitting.value = false
    }
}

const handleSaveFields = async () => {
    savingFields.value = true
    try {
        ElMessage.success('字段配置保存成功')
    } finally {
        savingFields.value = false
    }
}

const handleSaveAdvanced = () => ElMessage.success('高级配置保存成功')

const handleAddTemplate = () => ElMessage.info('新增模板功能开发中')
const handleEditTemplate = (row: any) => ElMessage.info(`编辑模板: ${row.templateName}`)
const handleDuplicateTemplate = async (row: any) => {
    try {
        await duplicateQueryTemplate(String(modelForm.id), String(row.id))
        ElMessage.success('复制成功')
        loadData(String(modelForm.id))
    } catch (e) {
        ElMessage.error('复制失败')
    }
}
const handlePreviewTemplate = async (row: any) => {
    try {
        const res = await previewQueryTemplate(String(modelForm.id), String(row.id))
        ElMessageBox.alert(`<pre>${res.sql}</pre>`, 'SQL预览', { dangerouslyUseHTMLString: true })
    } catch (e) {
        ElMessage.error('预览失败')
    }
}
const handleDeleteTemplate = (row: any) => {
    ElMessageBox.confirm(`确定要删除模板 ${row.templateName} 吗?`, '删除确认').then(() => {
        ElMessage.success('删除成功')
    })
}

const handleEditFieldEnhancement = (field: any) => {
    currentField.value = { ...field }
    enhancementVisible.value = true
}

const saveEnhancement = () => {
    const index = modelFields.value.findIndex(f => f.id === currentField.value.id)
    if (index > -1) {
        modelFields.value[index] = { ...currentField.value }
    }
    enhancementVisible.value = false
    ElMessage.success('增强配置已暂存')
}
</script>
<style scoped>
.model-edit {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.edit-form {
    max-width: 1000px;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 10px;
}

.form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
}

:deep(.action-column .cell) {
    white-space: nowrap;
}
</style>

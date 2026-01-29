<template>
    <div class="field-config-panel" :class="{ 'is-collapsed': collapsed }">
        <div class="panel-header" @click="$emit('toggle-collapse')">
            <div class="header-content">
                <el-icon>
                    <Operation />
                </el-icon>
                <span>配置详情</span>
            </div>
            <el-icon class="collapse-icon">
                <ArrowUp v-if="!collapsed" />
                <ArrowDown v-else />
            </el-icon>
        </div>
        <div v-show="!collapsed" class="panel-body">
            <el-tabs v-model="activeTab" class="config-tabs">
                <!-- 1. 表管理 (md_model_table) -->
                <el-tab-pane name="tables">
                    <template #label>
                        <el-tooltip content="表管理" placement="top">
                            <el-icon>
                                <Fold />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <div class="flex-between mb-4">
                            <h4 class="content-title">已选表清单</h4>
                        </div>
                        <div class="list-container">
                            <div v-for="node in tableNodes" :key="node.id" class="list-item" :class="{ 'is-active': selectedElement?.id === node.id }">
                                <div class="item-main" @click="handleSelectNode">
                                    <el-icon class="mr-2">
                                        <Grid />
                                    </el-icon>
                                    <span class="item-name">{{ node.data.label }}</span>
                                    <el-tag v-if="node.data.isMain" size="small" type="warning" class="ml-2">主</el-tag>
                                </div>
                                <div class="item-actions">
                                    <el-tooltip content="设置为主表" v-if="!node.data.isMain" placement="top">
                                        <el-button icon="Star" circle size="small" @click.stop="setMainTable(node.id)" />
                                    </el-tooltip>
                                    <el-button icon="Delete" circle size="small" type="danger" @click.stop="$emit('remove-table', node.id)" />
                                </div>
                            </div>
                            <el-empty v-if="tableNodes.length === 0" description="暂无已选表" :image-size="40" />
                        </div>
                        <div v-if="selectedElement?.type === 'table'" class="detail-section mt-4">
                            <h4 class="content-title">节点详情: {{ selectedElement.label }}</h4>
                            <el-form label-position="top" size="small">
                                <el-form-item label="对象别名 (alias)">
                                    <el-input v-model="selectedElement.data.alias" placeholder="设置表别名" />
                                </el-form-item>
                            </el-form>
                        </div>
                    </div>
                </el-tab-pane>
                <!-- 2. 字段选择 (md_model_field) -->
                <el-tab-pane name="fields">
                    <template #label>
                        <el-tooltip content="输出字段" placement="top">
                            <el-icon>
                                <List />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <div class="flex-between mb-2">
                            <h4 class="content-title">已选字段列表</h4>
                            <span class="count-tag">{{ selectedFieldsCount }}</span>
                        </div>
                        <el-input v-model="fieldSearch" placeholder="搜索字段名或别名" prefix-icon="Search" size="small" clearable class="mb-2" />
                        <div class="flex gap-2 mb-2">
                            <el-button size="small" type="primary" link @click="handleCheckAll">全选</el-button>
                            <el-button size="small" type="info" link @click="handleUncheckAll">清空</el-button>
                        </div>
                        <div class="list-container">
                            <template v-for="node in tableNodes" :key="'f-' + node.id">
                                <div v-if="hasSelectedFields(node)" class="table-group">
                                    <div class="group-header">{{ node.data.label }}</div>
                                    <div v-for="field in getSelectedFields(node)" :key="field.id" class="list-item">
                                        <div class="item-main">
                                            <span class="item-name">{{ field.name }}</span>
                                            <span class="item-type">{{ field.type }}</span>
                                        </div>
                                        <div class="field-edit">
                                            <el-input v-model="field.alias" size="small" placeholder="别名" class="field-alias" />
                                            <el-button icon="Close" circle size="small" @click="$emit('remove-field', node.id, field.id)" />
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <el-empty v-if="selectedFieldsCount === 0" description="暂未勾选字段" :image-size="40">
                                <p class="text-xs text-secondary">请双击画布中的表节点来勾选字段</p>
                            </el-empty>
                        </div>
                    </div>
                </el-tab-pane>
                <!-- 3. 关联配置 (md_model_join) -->
                <el-tab-pane name="joins">
                    <template #label>
                        <el-tooltip content="关联关系" placement="top">
                            <el-icon>
                                <Connection />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div v-if="selectedElement?.type === 'edge' || activeTab === 'joins'" class="tab-content">
                        <h4 class="content-title">关联设置 (JOIN)</h4>
                        <div v-if="selectedElement?.type === 'edge'">
                            <el-form label-position="top" size="small">
                                <el-form-item label="关联类型">
                                    <el-select v-model="selectedElement.data.joinType" style="width: 100%">
                                        <el-option label="LEFT JOIN" value="LEFT JOIN" />
                                        <el-option label="INNER JOIN" value="INNER JOIN" />
                                        <el-option label="RIGHT JOIN" value="RIGHT JOIN" />
                                    </el-select>
                                </el-form-item>
                                <el-form-item label="关联条件 (ON)">
                                    <div v-for="(cond, index) in (selectedElement.data.conditions || [])" :key="index" class="config-item-row mb-2">
                                        <el-select v-model="cond.operator1" placeholder="逻辑" size="small" style="width: 70px" v-if="Number(index) > 0">
                                            <el-option label="AND" value="AND" />
                                            <el-option label="OR" value="OR" />
                                        </el-select>
                                        <div v-else style="width: 70px"></div>
                                        <el-select v-model="cond.leftField" placeholder="源表字段" size="small" style="flex: 1">
                                            <el-option v-for="f in sourceFields" :key="f.value" :label="f.label" :value="f.value" />
                                        </el-select>
                                        <span style="font-size: 12px; margin: 0 4px">=</span>
                                        <el-select v-model="cond.rightField" placeholder="目标字段" size="small" style="flex: 1">
                                            <el-option v-for="f in targetFields" :key="f.value" :label="f.label" :value="f.value" />
                                        </el-select>
                                        <el-button type="danger" link :icon="Delete" @click="removeJoinCondition(Number(index))" />
                                    </div>
                                    <el-button type="primary" link icon="Plus" size="small" @click="addJoinCondition">添加条件</el-button>
                                    <div class="mt-2">
                                        <el-input v-model="selectedElement.data.joinCondition" type="textarea" :rows="2" placeholder="或者手动输入 SQL 条件" />
                                    </div>
                                </el-form-item>
                            </el-form>
                        </div>
                        <el-empty v-else description="在画布中选择连线以配置关联" :image-size="40" />
                    </div>
                </el-tab-pane>
                <!-- 4. 查询条件 (wheres) -->
                <el-tab-pane name="wheres">
                    <template #label>
                        <el-tooltip content="查询条件" placement="top">
                            <el-icon>
                                <Filter />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">查询条件 (WHERE)</h4>
                        <div v-for="(item, index) in modelConfig.wheres" :key="index" class="config-item-row">
                            <el-select v-model="item.operator1" placeholder="逻辑" size="small" style="width: 70px" v-if="Number(index) > 0">
                                <el-option label="AND" value="AND" />
                                <el-option label="OR" value="OR" />
                            </el-select>
                            <div v-else style="width: 70px"></div>
                            <el-select v-model="item.field" placeholder="选择字段" size="small" style="width: 140px">
                                <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                            </el-select>
                            <el-select v-model="item.operator" placeholder="操作符" size="small" style="width: 80px">
                                <el-option label="=" value="=" />
                                <el-option label="!=" value="!=" />
                                <el-option label=">" value=">" />
                                <el-option label=">=" value=">=" />
                                <el-option label="<" value="<" />
                                <el-option label="<=" value="<=" />
                                <el-option label="LIKE" value="LIKE" />
                                <el-option label="IN" value="IN" />
                                <el-option label="IS NULL" value="IS NULL" />
                            </el-select>
                            <el-input v-model="item.value" placeholder="值" size="small" style="flex: 1" v-if="item.operator !== 'IS NULL'" />
                            <el-button type="danger" link :icon="Delete" @click="removeWhere(Number(index))" />
                        </div>
                        <el-button type="primary" link icon="Plus" size="small" @click="addWhere">添加条件</el-button>
                    </div>
                </el-tab-pane>
                <!-- 5. 排序管理 (orders) -->
                <el-tab-pane name="orders">
                    <template #label>
                        <el-tooltip content="排序管理" placement="top">
                            <el-icon>
                                <Sort />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">排序 (ORDER BY)</h4>
                        <div v-for="(item, index) in modelConfig.orders" :key="index" class="config-item-row">
                            <el-select v-model="item.field" placeholder="选择字段" size="small" style="flex: 1">
                                <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                            </el-select>
                            <el-select v-model="item.direction" placeholder="排序" size="small" style="width: 80px">
                                <el-option label="ASC" value="ASC" />
                                <el-option label="DESC" value="DESC" />
                            </el-select>
                            <el-button type="danger" link :icon="Delete" @click="removeOrder(Number(index))" />
                        </div>
                        <el-button type="primary" link icon="Plus" size="small" @click="addOrder">添加排序</el-button>
                    </div>
                </el-tab-pane>
                <!-- 6. 分组聚合 (groups) -->
                <el-tab-pane name="groups">
                    <template #label>
                        <el-tooltip content="分组聚合" placement="top">
                            <el-icon>
                                <Histogram />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">分组与聚合 (GROUP / HAVING)</h4>
                        <el-form label-position="top" size="small">
                            <el-form-item label="分组字段 (GROUP BY)">
                                <div v-for="(item, index) in modelConfig.groupBy" :key="index" class="config-item-row mb-2">
                                    <el-select v-model="item.field" placeholder="选择字段" size="small" style="flex: 1">
                                        <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                                    </el-select>
                                    <el-button type="danger" link :icon="Delete" @click="removeGroup(Number(index))" />
                                </div>
                                <el-button type="primary" link icon="Plus" size="small" @click="addGroup">添加分组字段</el-button>
                            </el-form-item>
                            <el-form-item label="聚合过滤 (HAVING)">
                                <div v-for="(item, index) in (modelConfig.havings || [])" :key="index" class="config-item-row mb-2">
                                    <el-select v-model="item.operator1" placeholder="逻辑" size="small" style="width: 70px" v-if="Number(index) > 0">
                                        <el-option label="AND" value="AND" />
                                        <el-option label="OR" value="OR" />
                                    </el-select>
                                    <div v-else style="width: 70px"></div>
                                    <el-select v-model="item.func" placeholder="函数" size="small" style="width: 90px">
                                        <el-option label="COUNT" value="COUNT" />
                                        <el-option label="SUM" value="SUM" />
                                        <el-option label="AVG" value="AVG" />
                                        <el-option label="MAX" value="MAX" />
                                        <el-option label="MIN" value="MIN" />
                                    </el-select>
                                    <el-select v-model="item.field" placeholder="字段" size="small" style="width: 120px">
                                        <el-option label="*" value="*" />
                                        <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                                    </el-select>
                                    <el-select v-model="item.operator" placeholder="op" size="small" style="width: 70px">
                                        <el-option label="=" value="=" />
                                        <el-option label=">" value=">" />
                                        <el-option label="<" value="<" />
                                        <el-option label=">=" value=">=" />
                                        <el-option label="<=" value="<=" />
                                    </el-select>
                                    <el-input v-model="item.value" placeholder="值" size="small" style="flex: 1" />
                                    <el-button type="danger" link :icon="Delete" @click="removeHaving(Number(index))" />
                                </div>
                                <el-button type="primary" link icon="Plus" size="small" @click="addHaving">添加聚合过滤</el-button>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-tab-pane>
                <!-- 7. 通用设置 (settings) -->
                <el-tab-pane name="settings">
                    <template #label>
                        <el-tooltip content="通用设置" placement="top">
                            <el-icon>
                                <Odometer />
                            </el-icon>
                        </el-tooltip>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">数据限制 (LIMIT / OFFSET)</h4>
                        <el-divider border-style="dashed" style="margin: 12px 0;" />
                        <el-form label-position="top" size="small">
                            <div class="flex gap-4">
                                <el-form-item label="跳过行数 (OFFSET)" style="flex: 1">
                                    <el-input-number v-model="modelConfig.offset" :min="0" style="width: 100%" controls-position="right" />
                                </el-form-item>
                                <el-form-item label="最大显示 (LIMIT)" style="flex: 1">
                                    <el-input-number v-model="modelConfig.limit" :min="1" :max="10000" style="width: 100%" controls-position="right" />
                                </el-form-item>
                            </div>
                        </el-form>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<script setup lang="ts">
import {
    ArrowDown,
    ArrowUp,
    Connection,
    Delete,
    Fold, Grid,
    List,
    Operation
} from '@element-plus/icons-vue';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    elements: any[]
    selectedElement: any
    collapsed: boolean
    modelConfig: any
}>()

const emit = defineEmits(['toggle-collapse', 'remove-table', 'remove-field'])

const activeTab = ref('tables')
const fieldSearch = ref('')

// 数据过滤
const tableNodes = computed(() => props.elements.filter(el => el.type === 'table'))

// 选中的字段计数
const selectedFieldsCount = computed(() => {
    let count = 0
    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            count += node.data.fields.filter((f: any) => f.selected).length
        }
    })
    return count
})

const hasSelectedFields = (node: any) => {
    // 基础检查：是否有选中字段
    const hasSelected = node.data.fields && node.data.fields.some((f: any) => f.selected)
    if (!hasSelected) return false

    // 搜索过滤检查
    if (!fieldSearch.value) return true
    return node.data.fields.some((f: any) => f.selected && (
        f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
        (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
    ))
}

const getSelectedFields = (node: any) => {
    let fields = node.data.fields.filter((f: any) => f.selected)
    if (fieldSearch.value) {
        fields = fields.filter((f: any) =>
            f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
            (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
        )
    }
    return fields
}

const handleCheckAll = () => {
    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            node.data.fields.forEach((f: any) => {
                // If search is active, only check visible ones? 
                // Creating a robust filter check here
                if (fieldSearch.value) {
                    const match = f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
                        (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
                    if (match) f.selected = true
                } else {
                    f.selected = true
                }
            })
        }
    })
}

const handleUncheckAll = () => {
    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            node.data.fields.forEach((f: any) => {
                if (fieldSearch.value) {
                    const match = f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
                        (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
                    if (match) f.selected = false
                } else {
                    f.selected = false
                }
            })
        }
    })
}

// 获取所有可用字段 (用于 Where/Order 选择)
const allFields = computed(() => {
    const fields: any[] = []
    props.elements.forEach(el => {
        if (el.type === 'table') {
            const tableName = el.data.tableName || el.data.label
            if (el.data.fields) {
                el.data.fields.forEach((f: any) => {
                    fields.push({
                        label: `${tableName}.${f.name}`,
                        // 使用 节点ID::字段名 作为唯一标识，方便后续解析
                        value: `${el.id}::${f.name}`,
                        nodeId: el.id,
                        tableName: tableName,
                        columnName: f.name
                    })
                })
            }
        }
    })
    return fields
})

const addWhere = () => {
    if (!props.modelConfig.wheres) props.modelConfig.wheres = []
    props.modelConfig.wheres.push({ operator1: 'AND', field: '', operator: '=', value: '' })
}

const removeWhere = (index: number) => {
    props.modelConfig.wheres.splice(index, 1)
}

const addOrder = () => {
    if (!props.modelConfig.orders) props.modelConfig.orders = []
    props.modelConfig.orders.push({ field: '', direction: 'ASC' })
}

const removeOrder = (index: number) => {
    props.modelConfig.orders.splice(index, 1)
}

const addHaving = () => {
    if (!props.modelConfig.havings) props.modelConfig.havings = []
    props.modelConfig.havings.push({ operator1: 'AND', func: 'COUNT', field: '*', operator: '>', value: '' })
}

const removeHaving = (index: number) => {
    props.modelConfig.havings.splice(index, 1)
}

const addGroup = () => {
    if (!props.modelConfig.groupBy) props.modelConfig.groupBy = []
    props.modelConfig.groupBy.push({ field: '' })
}

const removeGroup = (index: number) => {
    props.modelConfig.groupBy.splice(index, 1)
}

const sourceFields = computed(() => {
    if (props.selectedElement?.type !== 'edge') return []
    const sourceNode = props.elements.find(el => el.id === props.selectedElement.source)
    if (!sourceNode || !sourceNode.data.fields) return []
    return sourceNode.data.fields.map((f: any) => ({
        label: f.name,
        value: f.name
    }))
})

const targetFields = computed(() => {
    if (props.selectedElement?.type !== 'edge') return []
    const targetNode = props.elements.find(el => el.id === props.selectedElement.target)
    if (!targetNode || !targetNode.data.fields) return []
    return targetNode.data.fields.map((f: any) => ({
        label: f.name,
        value: f.name
    }))
})

const addJoinCondition = () => {
    if (!props.selectedElement.data.conditions) {
        props.selectedElement.data.conditions = []
    }
    props.selectedElement.data.conditions.push({ operator1: 'AND', leftField: '', operator: '=', rightField: '' })
}

const removeJoinCondition = (index: number) => {
    props.selectedElement.data.conditions.splice(index, 1)
}

// 交互逻辑
const handleSelectNode = () => {
    // 这里我们保持选中逻辑一致，通过选中状态触发
    activeTab.value = 'tables'
}

const setMainTable = (tableId: string) => {
    tableNodes.value.forEach(node => {
        node.data.isMain = (node.id === tableId)
    })
}

// 监听选择
watch(() => props.selectedElement, (newEl) => {
    if (!newEl) return
    if (newEl.type === 'table') {
        activeTab.value = 'tables'
    } else if (newEl.type === 'edge') {
        activeTab.value = 'joins'
    }
}, { immediate: true })
</script>
<style scoped>
.field-config-panel {
    display: flex;
    flex-direction: column;
    background-color: #fff;
    transition: width 0.1s linear, height 0.1s linear;
    /* 缩短过渡时间使拖拽更跟手 */
    position: relative;
    box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
    border-radius: 4px;
}

/* Resize logic removed */

.panel-header {
    padding: 10px 16px;
    border-bottom: 1px solid #f2f6fc;
    font-weight: 600;
    color: #303133;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: move;
    user-select: none;
}

.header-content {
    display: flex;
    align-items: center;
    gap: 8px;
}

.collapse-icon {
    cursor: pointer;
    font-size: 14px;
    color: #909399;
}

.panel-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.config-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    /* Important for flex child to shrink/grow correctly */
}

:deep(.el-tabs__header) {
    margin: 0;
    background-color: #fcfcfc;
}

:deep(.el-tabs__content) {
    flex: 1;
    overflow-y: auto;
}

.tab-content {
    padding: 16px;
}

.config-item-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
    align-items: center;
}

.content-title {
    margin-top: 0;
    margin-bottom: 12px;
    font-size: 13px;
    color: #303133;
    border-left: 3px solid #409eff;
    padding-left: 8px;
    font-weight: 600;
}

.list-container {
    border: 1px solid #f0f2f5;
    border-radius: 4px;
    overflow: hidden;
}

.list-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid #fafafa;
}

.list-item:hover {
    background-color: #f5f7fa;
}

.list-item.is-active {
    background-color: #ecf5ff;
}

.item-main {
    flex: 1;
    display: flex;
    align-items: center;
    cursor: pointer;
}

.item-name {
    font-size: 13px;
    color: #303133;
    font-weight: 500;
}

.item-type {
    font-size: 11px;
    color: #909399;
    margin-left: 8px;
}

.item-actions {
    display: flex;
    gap: 4px;
}

.table-group {
    background-color: #fff;
    margin-bottom: 8px;
    border: 1px solid #f0f2f5;
    border-radius: 4px;
}

.group-header {
    background-color: #f8f9fb;
    padding: 6px 12px;
    font-size: 12px;
    font-weight: bold;
    color: #409EFF;
    border-bottom: 1px solid #f0f2f5;
}

.field-edit {
    display: flex;
    align-items: center;
    gap: 8px;
}

.field-alias {
    width: 110px;
}

.count-tag {
    background-color: #409EFF;
    color: #fff;
    padding: 2px 8px;
    border-radius: 10px;
    font-size: 11px;
}

.flex-between {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.detail-section {
    padding: 16px;
    background-color: #fcfcfc;
    border-radius: 4px;
    border: 1px solid #f0f2f5;
}

.mt-4 {
    margin-top: 16px;
}

.mb-2 {
    margin-bottom: 8px;
}

.mb-4 {
    margin-bottom: 16px;
}

.mr-2 {
    margin-right: 8px;
}

.ml-2 {
    margin-left: 8px;
}

.text-xs {
    font-size: 12px;
}

.text-secondary {
    color: #909399;
}

.is-collapsed .panel-header {
    border-bottom: none;
}

:deep(.el-input--small .el-input__wrapper) {
    padding: 1px 8px;
}
</style>

<template>
    <div class="field-config-panel" :class="{ 'is-collapsed': collapsed }">
        <div class="panel-header" @click="$emit('toggle-collapse')">
            <div class="header-content">
                <el-icon>
                    <Operation />
                </el-icon>
                <span>配置详情</span>
            </div>
            <div class="header-actions" @click.stop>
                <el-icon class="header-icon" title="最大化/还原" @click="$emit('toggle-maximize')">
                    <CopyDocument v-if="maximized" />
                    <FullScreen v-else />
                </el-icon>
                <el-icon class="header-icon collapse-icon" @click="$emit('toggle-collapse')">
                    <ArrowUp v-if="!collapsed" />
                    <ArrowDown v-else />
                </el-icon>
            </div>
        </div>
        <div v-show="!collapsed" class="panel-body">
            <el-tabs v-model="activeTab" class="config-tabs">
                <!-- 1. 表管理 (md_model_table) -->
                <el-tab-pane name="tables">
                    <template #label>
                        <span class="custom-tab-label">
                            <el-icon>
                                <Fold />
                            </el-icon>
                            <span>表/视图</span>
                        </span>
                    </template>
                    <div class="tab-content">
                        <div class="flex-between mb-4">
                            <h4 class="content-title">已选表清单</h4>
                        </div>
                        <el-table ref="modelTableRef" :data="tableNodes" size="small" border style="width: 100%" highlight-current-row @row-click="handleSelectNode" class="config-table">
                            <el-table-column label="别名" min-width="120">
                                <template #default="{ row }">
                                    <div class="flex-align-center">
                                        <el-icon class="mr-1" v-if="row.data.isMain" :style="{ verticalAlign: 'middle' }">
                                            <StarFilled style="color: #E6A23C" />
                                        </el-icon>
                                        <el-input v-model="row.data.tableAlias" size="small" placeholder="输入别名" @input="handleAliasChange" :class="{ 'border-danger': isAliasDuplicate(row.data.tableAlias) }" />
                                        <el-tooltip v-if="isAliasDuplicate(row.data.tableAlias)" content="别名重复" placement="top">
                                            <el-icon class="text-danger ml-1">
                                                <Warning />
                                            </el-icon>
                                        </el-tooltip>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column label="标题" min-width="120">
                                <template #default="{ row }">
                                    <el-input v-model="row.data.tableTitle" size="small" placeholder="设置易读标题" />
                                </template>
                            </el-table-column>
                            <el-table-column label="表/视图" min-width="100" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <span class="text-secondary text-xs">{{ row.data.tableName }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="85" align="center">
                                <template #default="{ row }">
                                    <div class="table-actions">
                                        <el-tooltip content="设为主表" v-if="!row.data.isMain" placement="top">
                                            <el-button icon="Star" circle size="small" @click.stop="setMainTable(row.id)" />
                                        </el-tooltip>
                                        <el-button icon="Delete" circle size="small" type="danger" @click.stop="$emit('remove-table', row.id)" />
                                    </div>
                                </template>
                            </el-table-column>
                            <template #empty>
                                <el-empty description="暂无已选表" :image-size="40" />
                            </template>
                        </el-table>
                    </div>
                </el-tab-pane>
                <!-- 2. 字段选择 (md_model_field) -->
                <el-tab-pane name="fields">
                    <template #label>
                        <span class="custom-tab-label">
                            <el-icon>
                                <List />
                            </el-icon>
                            <span>字段</span>
                        </span>
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
                        <div class="list-container no-border">
                            <template v-for="node in tableNodes" :key="'f-' + node.id">
                                <div v-if="hasSelectedFields(node)" class="table-group mb-4">
                                    <div class="group-header">
                                        <el-icon class="mr-1">
                                            <Grid />
                                        </el-icon>
                                        {{ node.data.tableAlias || node.data.label }} ({{ node.data.tableName }})
                                    </div>
                                    <el-table :data="getSelectedFields(node)" size="small" border style="width: 100%">
                                        <el-table-column label="字段" align="center">
                                            <el-table-column prop="name" label="名称" min-width="120" show-overflow-tooltip />
                                            <el-table-column prop="type" label="类型" width="80" show-overflow-tooltip />
                                        </el-table-column>
                                        <el-table-column label="配制" align="center">
                                            <el-table-column label="别名(alias)" min-width="120">
                                                <template #default="{ row }">
                                                    <div class="flex-align-center">
                                                        <el-input v-model="row.alias" size="small" :placeholder="row.name" @input="handleFieldAliasChange" :class="{ 'border-danger': isFieldAliasDuplicate(row.alias || row.name) }" />
                                                        <el-tooltip v-if="isFieldAliasDuplicate(row.alias || row.name)" content="字段别名全局冲突" placement="top">
                                                            <el-icon class="text-danger ml-1">
                                                                <Warning />
                                                            </el-icon>
                                                        </el-tooltip>
                                                    </div>
                                                </template>
                                            </el-table-column>
                                            <el-table-column label="函数" width="100">
                                                <template #default="{ row }">
                                                    <el-input v-model="row.func" size="small" placeholder="函数" />
                                                </template>
                                            </el-table-column>
                                            <el-table-column label="聚合" width="90">
                                                <template #default="{ row }">
                                                    <el-select v-model="row.aggFunc" size="small" clearable placeholder="无">
                                                        <el-option label="SUM" value="sum" />
                                                        <el-option label="COUNT" value="count" />
                                                        <el-option label="AVG" value="avg" />
                                                        <el-option label="MAX" value="max" />
                                                        <el-option label="MIN" value="min" />
                                                        <el-option label="DISTINCT" value="distinct" />
                                                    </el-select>
                                                </template>
                                            </el-table-column>
                                            <el-table-column label="标题" min-width="120">
                                                <template #default="{ row }">
                                                    <el-input v-model="row.showTitle" size="small" :placeholder="row.comment || row.name" />
                                                </template>
                                            </el-table-column>
                                            <el-table-column label="宽度" width="100">
                                                <template #default="{ row }">
                                                    <el-input-number v-model="row.showWidth" :min="40" :max="1000" size="small" controls-position="right" style="width: 100%" />
                                                </template>
                                            </el-table-column>
                                        </el-table-column>
                                        <el-table-column label="操作" width="50" align="center" fixed="right">
                                            <template #default="{ row }">
                                                <el-button circle size="small" type="danger" plain @click="$emit('remove-field', node.id, row.id)">
                                                    <el-icon>
                                                        <Close />
                                                    </el-icon>
                                                </el-button>
                                            </template>
                                        </el-table-column>
                                    </el-table>
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
                        <span class="custom-tab-label">
                            <el-icon>
                                <Connection />
                            </el-icon>
                            <span>关联</span>
                        </span>
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
                        <span class="custom-tab-label">
                            <el-icon>
                                <Filter />
                            </el-icon>
                            <span>查询</span>
                        </span>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">查询条件 (WHERE)</h4>
                        <div v-for="(item, index) in modelConfig.wheres" :key="index" class="config-item-row">
                            <el-select v-model="item.operator1" placeholder="逻辑" size="small" style="width: 70px" v-if="Number(index) > 0">
                                <el-option label="AND" value="AND" />
                                <el-option label="OR" value="OR" />
                            </el-select>
                            <div v-else style="width: 70px"></div>
                            <el-select v-model="item.brackets1" placeholder="(" size="small" style="width: 60px" clearable>
                                <el-option label="" value="" />
                                <el-option label="(" value="(" />
                                <el-option label="((" value="((" />
                                <el-option label="(((" value="(((" />
                            </el-select>
                            <el-input v-model="item.func" placeholder="函数" size="small" style="width: 80px" @input="handleWhereFuncChange(item)" />
                            <el-select v-model="item.field" placeholder="选择字段" size="small" style="width: 140px" @change="val => handleWhereFieldChange(item, val)">
                                <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                            </el-select>
                            <el-select v-model="item.operator" placeholder="操作符" size="small" style="width: 80px" @change="val => handleWhereOperatorChange(item, val)">
                                <el-option label="=" value="=" />
                                <el-option label="!=" value="!=" />
                                <el-option label=">" value=">" />
                                <el-option label=">=" value=">=" />
                                <el-option label="<" value="<" />
                                <el-option label="<=" value="<=" />
                                <el-option label="LIKE" value="LIKE" />
                                <el-option label="IN" value="IN" />
                                <el-option label="IS NULL" value="IS NULL" />
                                <el-option label="BETWEEN" value="BETWEEN" />
                            </el-select>
                            <el-input v-model="item.value" placeholder="值1" size="small" style="flex: 1" v-if="item.operator !== 'IS NULL'" />
                            <el-input v-model="item.value2" placeholder="值2" size="small" style="flex: 1" v-if="item.operator === 'BETWEEN'" />
                            <el-select v-model="item.brackets2" placeholder=")" size="small" style="width: 60px" clearable>
                                <el-option label="" value="" />
                                <el-option label=")" value=")" />
                                <el-option label="))" value="))" />
                                <el-option label=")))" value=")))" />
                            </el-select>
                            <el-button-group class="ml-2">
                                <el-button size="small" link :icon="ArrowUp" @click="moveUp(modelConfig.wheres, Number(index))" :disabled="index === 0" />
                                <el-button size="small" link :icon="ArrowDown" @click="moveDown(modelConfig.wheres, Number(index))" :disabled="index === modelConfig.wheres.length - 1" />
                            </el-button-group>
                            <el-button type="danger" link :icon="Delete" @click="removeWhere(Number(index))" />
                        </div>
                        <el-button class="add-btn-block" icon="Plus" border-style="dashed" @click="addWhere">添加条件</el-button>
                    </div>
                </el-tab-pane>
                <!-- 5. 排序管理 (orders) -->
                <el-tab-pane name="orders">
                    <template #label>
                        <span class="custom-tab-label">
                            <el-icon>
                                <Sort />
                            </el-icon>
                            <span>排序</span>
                        </span>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">排序 (ORDER BY)</h4>
                        <div v-for="(item, index) in modelConfig.orders" :key="index" class="config-item-row">
                            <el-input v-model="item.func" placeholder="函数" size="small" style="width: 80px" />
                            <el-select v-model="item.field" placeholder="选择字段" size="small" style="flex: 1">
                                <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                            </el-select>
                            <el-select v-model="item.direction" placeholder="排序" size="small" style="width: 80px">
                                <el-option label="ASC" value="ASC" />
                                <el-option label="DESC" value="DESC" />
                            </el-select>
                            <el-button-group class="ml-2">
                                <el-button size="small" link :icon="ArrowUp" @click="moveUp(modelConfig.orders, Number(index))" :disabled="index === 0" />
                                <el-button size="small" link :icon="ArrowDown" @click="moveDown(modelConfig.orders, Number(index))" :disabled="index === modelConfig.orders.length - 1" />
                            </el-button-group>
                            <el-button type="danger" link :icon="Delete" @click="removeOrder(Number(index))" />
                        </div>
                        <el-button class="add-btn-block" icon="Plus" border-style="dashed" @click="addOrder">添加排序</el-button>
                    </div>
                </el-tab-pane>
                <!-- 6. 分组聚合 (groups) -->
                <el-tab-pane name="groups">
                    <template #label>
                        <span class="custom-tab-label">
                            <el-icon>
                                <Histogram />
                            </el-icon>
                            <span>分组</span>
                        </span>
                    </template>
                    <div class="tab-content">
                        <h4 class="content-title">分组与聚合 (GROUP / HAVING)</h4>
                        <el-form label-position="top" size="small">
                            <el-form-item label="分组字段 (GROUP BY)">
                                <div v-for="(item, index) in modelConfig.groupBy" :key="index" class="config-item-row mb-2">
                                    <el-select v-model="item.agg_func" placeholder="聚合" size="small" style="width: 80px" clearable>
                                        <el-option label="" value="" />
                                        <el-option label="SUM" value="sum" />
                                        <el-option label="COUNT" value="count" />
                                        <el-option label="AVG" value="avg" />
                                        <el-option label="MAX" value="max" />
                                        <el-option label="MIN" value="min" />
                                    </el-select>
                                    <el-input v-model="item.func" placeholder="函数" size="small" style="width: 80px" />
                                    <el-select v-model="item.field" placeholder="选择字段" size="small" style="flex: 1">
                                        <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                                    </el-select>
                                    <el-button-group class="ml-2">
                                        <el-button size="small" link :icon="ArrowUp" @click="moveUp(modelConfig.groupBy, Number(index))" :disabled="index === 0" />
                                        <el-button size="small" link :icon="ArrowDown" @click="moveDown(modelConfig.groupBy, Number(index))" :disabled="index === modelConfig.groupBy.length - 1" />
                                    </el-button-group>
                                    <el-button type="danger" link :icon="Delete" @click="removeGroup(Number(index))" />
                                </div>
                                <el-button class="add-btn-block" icon="Plus" border-style="dashed" @click="addGroup">添加分组字段</el-button>
                            </el-form-item>
                            <el-form-item label="聚合过滤 (HAVING)">
                                <div v-for="(item, index) in (modelConfig.havings || [])" :key="index" class="config-item-row mb-2">
                                    <el-select v-model="item.operator1" placeholder="逻辑" size="small" style="width: 70px" v-if="Number(index) > 0">
                                        <el-option label="AND" value="AND" />
                                        <el-option label="OR" value="OR" />
                                    </el-select>
                                    <div v-else style="width: 70px"></div>
                                    <el-select v-model="item.brackets1" placeholder="(" size="small" style="width: 60px" clearable>
                                        <el-option label="" value="" />
                                        <el-option label="(" value="(" />
                                        <el-option label="((" value="((" />
                                        <el-option label="(((" value="(((" />
                                    </el-select>
                                    <el-select v-model="item.agg_func" placeholder="聚合" size="small" style="width: 90px" clearable @change="handleHavingAggFuncChange(item)">
                                        <el-option label="" value="" />
                                        <el-option label="COUNT" value="COUNT" />
                                        <el-option label="SUM" value="SUM" />
                                        <el-option label="AVG" value="AVG" />
                                        <el-option label="MAX" value="MAX" />
                                        <el-option label="MIN" value="MIN" />
                                    </el-select>
                                    <el-input v-model="item.func" placeholder="函数" size="small" style="width: 80px" @input="handleHavingFuncChange(item)" />
                                    <el-select v-model="item.field" placeholder="字段" size="small" style="width: 120px" @change="(val: string) => handleHavingFieldChange(item, val)">
                                        <el-option label="*" value="*" />
                                        <el-option v-for="f in allFields" :key="f.value" :label="f.label" :value="f.value" />
                                    </el-select>
                                    <el-select v-model="item.operator" placeholder="比较" size="small" style="width: 80px" @change="(val: string) => handleHavingOperatorChange(item, val)">
                                        <el-option label="=" value="=" />
                                        <el-option label=">" value=">" />
                                        <el-option label="<" value="<" />
                                        <el-option label=">=" value=">=" />
                                        <el-option label="<=" value="<=" />
                                        <el-option label="<>" value="<>" />
                                        <el-option label="LIKE" value="LIKE" />
                                        <el-option label="BETWEEN" value="BETWEEN" />
                                        <el-option label="IS NULL" value="IS NULL" />
                                        <el-option label="IS NOT NULL" value="IS NOT NULL" />
                                    </el-select>
                                    <el-input v-model="item.value" placeholder="值" size="small" style="flex: 1" v-if="!['IS NULL', 'IS NOT NULL'].includes(item.operator)" />
                                    <span v-if="item.operator === 'BETWEEN'" class="mx-1 text-gray-400">AND</span>
                                    <el-input v-model="item.value2" placeholder="值2" size="small" style="flex: 1" v-if="item.operator === 'BETWEEN'" />
                                    <el-select v-model="item.brackets2" placeholder=")" size="small" style="width: 60px" clearable>
                                        <el-option label="" value="" />
                                        <el-option label=")" value=")" />
                                        <el-option label="))" value="))" />
                                        <el-option label=")))" value=")))" />
                                    </el-select>
                                    <el-button-group class="ml-2">
                                        <el-button size="small" link :icon="ArrowUp" @click="moveUp(modelConfig.havings, Number(index))" :disabled="index === 0" />
                                        <el-button size="small" link :icon="ArrowDown" @click="moveDown(modelConfig.havings, Number(index))" :disabled="index === (modelConfig.havings || []).length - 1" />
                                    </el-button-group>
                                    <el-button type="danger" link :icon="Delete" @click="removeHaving(Number(index))" />
                                </div>
                                <el-button class="add-btn-block" icon="Plus" border-style="dashed" @click="addHaving">添加聚合过滤</el-button>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-tab-pane>
                <!-- 7. 通用设置 (settings) -->
                <el-tab-pane name="settings">
                    <template #label>
                        <span class="custom-tab-label">
                            <el-icon>
                                <Odometer />
                            </el-icon>
                            <span>设置</span>
                        </span>
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
                                    <el-input-number v-model="modelConfig.limit" :min="0" :max="10000" style="width: 100%" controls-position="right" />
                                </el-form-item>
                            </div>
                        </el-form>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
        <!-- 进度条弹窗 -->
        <el-dialog v-model="processing" title="处理中" width="300px" :close-on-click-modal="false" :show-close="false" align-center>
            <div class="flex flex-col items-center py-4">
                <el-progress type="circle" :percentage="progressPercentage" :status="progressStatus" />
                <div class="mt-4 text-gray-600">{{ progressText }}</div>
            </div>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import {
    ArrowDown, ArrowUp,
    Close,
    Connection, CopyDocument, Delete,
    Filter,
    Fold, FullScreen, Grid,
    Histogram,
    List,
    Odometer,
    Operation,
    Sort,
    StarFilled,
    Warning
} from '@element-plus/icons-vue';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    elements: any[]
    selectedElement: any
    collapsed: boolean
    modelConfig: any
    maximized: boolean
}>()

const emit = defineEmits(['toggle-collapse', 'remove-table', 'remove-field', 'select-node', 'toggle-maximize'])

const modelTableRef = ref<any>(null)
const activeTab = ref('tables')
const fieldSearch = ref('')

// 进度条状态
const processing = ref(false)
const progressPercentage = ref(0)
const progressStatus = ref('')
const progressText = ref('')

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

// 列表排序辅助函数
const moveUp = (list: any[], index: number) => {
    if (index > 0) {
        const item = list[index]
        list.splice(index, 1)
        list.splice(index - 1, 0, item)
    }
}

const moveDown = (list: any[], index: number) => {
    if (index < list.length - 1) {
        const item = list[index]
        list.splice(index, 1)
        list.splice(index + 1, 0, item)
    }
}

// 批量处理工具函数
const processBatch = async (items: any[], updateFn: (item: any) => void, actionName: string) => {
    if (items.length === 0) return

    processing.value = true
    progressPercentage.value = 0
    progressStatus.value = ''
    progressText.value = `${actionName}中...`

    const batchSize = 500 // 每批处理数量
    const total = items.length
    let current = 0

    return new Promise<void>((resolve) => {
        const nextBatch = () => {
            const end = Math.min(current + batchSize, total)

            for (let i = current; i < end; i++) {
                updateFn(items[i])
            }

            current = end
            progressPercentage.value = Math.floor((current / total) * 100)

            if (current < total) {
                // 使用 setTimeout 让出主线程，允许 UI 渲染
                setTimeout(nextBatch, 10)
            } else {
                progressStatus.value = 'success'
                progressText.value = '完成'
                setTimeout(() => {
                    processing.value = false
                    resolve()
                }, 500)
            }
        }

        nextBatch()
    })
}

const handleCheckAll = async () => {
    const allFields: any[] = []

    // 1. 收集所有需要操作的字段
    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            node.data.fields.forEach((f: any) => {
                if (fieldSearch.value) {
                    const match = f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
                        (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
                    if (match) allFields.push(f)
                } else {
                    allFields.push(f)
                }
            })
        }
    })

    // 2. 批量处理
    await processBatch(allFields, (f) => { f.selected = true }, '全选')
}

const handleUncheckAll = async () => {
    const allFields: any[] = []

    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            node.data.fields.forEach((f: any) => {
                if (fieldSearch.value) {
                    const match = f.name.toLowerCase().includes(fieldSearch.value.toLowerCase()) ||
                        (f.alias && f.alias.toLowerCase().includes(fieldSearch.value.toLowerCase()))
                    if (match) allFields.push(f)
                } else {
                    allFields.push(f)
                }
            })
        }
    })

    await processBatch(allFields, (f) => { f.selected = false }, '清空')
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
    props.modelConfig.wheres.push({ operator1: 'AND', brackets1: '', func: '', field: '', operator: '=', value: '', value2: '', brackets2: '', param_key: '' })
}

// 自动填充字段参数 helper
// 自动填充字段参数 helper
const autoFillParams = (item: any, fieldVal?: string) => {
    const fieldValue = fieldVal || item.field
    if (!fieldValue) return

    // 从 allFields 中查找完整信息以构建唯一参数名
    const targetField = allFields.value.find((f: any) => f.value === fieldValue)
    if (targetField) {
        let paramName = `${targetField.tableName}_${targetField.columnName}`

        // 追加函数前缀
        if (item.func) {
            paramName = `${item.func}_${paramName}`
        }
        if (item.agg_func) {
            paramName = `${item.agg_func}_${paramName}`
        }

        // 转小写并用于生成 param_key
        paramName = paramName.toLowerCase().replace(/[^a-z0-9_]/g, '_')

        // 自动设置 param_key
        item.param_key = paramName

        // 如果值为空，或者已经是参数格式（以:开头），则自动更新
        if (!item.value || item.value.startsWith(':')) {
            item.value = `:${paramName}`
        }
        if (item.operator === 'BETWEEN') {
            if (!item.value2 || item.value2.startsWith(':')) {
                item.value2 = `:${paramName}_end`
            }
        }
    } else {
        // Fallback checks if simple parsing is needed, but lookup should typically succeed
        const parts = fieldValue.split('::')
        if (parts.length === 2) {
            const fieldName = parts[1]
            item.param_key = fieldName
            if (!item.value || item.value.startsWith(':')) {
                item.value = `:${fieldName}`
            }
            if (item.operator === 'BETWEEN') {
                if (!item.value2 || item.value2.startsWith(':')) {
                    item.value2 = `:${fieldName}_end`
                }
            }
        }
    }
}

const handleWhereFieldChange = (item: any, val: any) => {
    autoFillParams(item, val)
}

const handleWhereOperatorChange = (item: any, val: any) => {
    if (val === 'BETWEEN') { // use val instead of item.operator
        autoFillParams(item)
    }
}

const handleWhereFuncChange = (item: any) => {
    // 如果输入了函数，也可以触发自动填充
    if (item.func) {
        autoFillParams(item)
    }
}

const removeWhere = (index: number) => {
    props.modelConfig.wheres.splice(index, 1)
}

const addOrder = () => {
    if (!props.modelConfig.orders) props.modelConfig.orders = []
    props.modelConfig.orders.push({ field: '', direction: 'ASC', func: '' })
}

const removeOrder = (index: number) => {
    props.modelConfig.orders.splice(index, 1)
}

const addHaving = () => {
    if (!props.modelConfig.havings) props.modelConfig.havings = []
    props.modelConfig.havings.push({
        field: '',
        func: '',
        agg_func: '',
        operator: '=',
        value: '',
        operator1: 'AND',
        brackets1: '',
        brackets2: '',
        value2: '',
        param_key: ''
    })
}

const handleHavingOperatorChange = (item: any, val: string) => {
    if (['IS NULL', 'IS NOT NULL'].includes(val)) {
        item.value = ''
        item.value2 = ''
    }
}

const handleHavingFieldChange = (item: any, fieldVal: string) => {
    autoFillParams(item, fieldVal)
}

const handleHavingAggFuncChange = (item: any) => {
    autoFillParams(item)
}

const handleHavingFuncChange = (item: any) => {
    autoFillParams(item)
}

const removeHaving = (index: number) => {
    props.modelConfig.havings.splice(index, 1)
}

const addGroup = () => {
    if (!props.modelConfig.groupBy) props.modelConfig.groupBy = []
    props.modelConfig.groupBy.push({ field: '', func: '', agg_func: '' })
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
// 交互逻辑
const handleSelectNode = (node: any) => {
    emit('select-node', node.id)
    activeTab.value = 'tables'
}

// 别名唯一性校验
const isAliasDuplicate = (alias: string) => {
    if (!alias) return false
    const count = tableNodes.value.filter(el => el.data.tableAlias === alias).length
    return count > 1
}

const isFieldAliasDuplicate = (alias: string) => {
    if (!alias) return false
    let count = 0
    tableNodes.value.forEach(node => {
        if (node.data.fields) {
            count += node.data.fields.filter((f: any) => f.selected && (f.alias || f.name) === alias).length
        }
    })
    return count > 1
}

const handleAliasChange = (val: string) => {
    if (props.selectedElement && props.selectedElement.type === 'table') {
        // 同步 label 用于画布显示联动逻辑 (如果 TableNode 逻辑没改，这里保持)
        // 但现在 TableNode 优先显示 tableTitle，所以主要改 data.tableAlias
        props.selectedElement.label = val
        props.selectedElement.data.label = val
    }
}

const handleFieldAliasChange = () => {
    // 可以在这里做一些实时处理
}

const setMainTable = (tableId: string) => {
    tableNodes.value.forEach(node => {
        node.data.isMain = (node.id === tableId)
    })
}

// 监听选择
watch(() => props.selectedElement, (newEl) => {
    if (!newEl) {
        if (modelTableRef.value) modelTableRef.value.setCurrentRow(null)
        return
    }
    if (newEl.type === 'table') {
        activeTab.value = 'tables'
        // 延迟执行以确保 el-table 已渲染并能找到行
        setTimeout(() => {
            const row = tableNodes.value.find(n => n.id === newEl.id)
            if (row && modelTableRef.value) {
                modelTableRef.value.setCurrentRow(row)
            }
        }, 100)
    } else if (newEl.type === 'edge') {
        activeTab.value = 'joins'
    }
}, { immediate: true })
</script>
<style scoped>
.field-config-panel {
    height: 100%;
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

.header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
}

.header-icon {
    cursor: pointer;
    font-size: 14px;
    color: #909399;
    transition: color 0.2s;
}

.header-icon:hover {
    color: #409EFF;
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
    height: 100%;
}

:deep(.el-tabs__header) {
    margin: 0;
    background-color: #fcfcfc;
}

:deep(.el-tabs__content) {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
}

/* 优化滚动条样式 */
:deep(.el-tabs__content::-webkit-scrollbar) {
    width: 6px;
}

:deep(.el-tabs__content::-webkit-scrollbar-thumb) {
    background: #e4e7ed;
    border-radius: 3px;
}

:deep(.el-tabs__content::-webkit-scrollbar-thumb:hover) {
    background: #c0c4cc;
}

:deep(.el-tab-pane) {
    height: 100%;
}

.tab-content {
    padding: 16px;
    display: flex;
    flex-direction: column;
    height: 100%;
    box-sizing: border-box;
}

.config-item-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
    align-items: center;
    width: 100%;
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
    overflow-y: auto;
    flex: 1;
    min-height: 0;
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

.item-main {
    flex: 1;
    display: flex;
    align-items: center;
    cursor: pointer;
}

.flex-align-center {
    display: flex;
    align-items: center;
}

.is-main-text {
    font-weight: bold;
    color: #E6A23C;
}

.table-actions {
    display: flex;
    justify-content: center;
    gap: 4px;
}

.no-border {
    border: none !important;
}

.config-table {
    --el-table-header-bg-color: #f8f9fb;
    margin-bottom: 4px;
    flex: 1;
}

:deep(.el-table .el-table__cell) {
    padding: 4px 0;
}

.table-group {
    background-color: #fff;
    margin-bottom: 12px;
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

.ml-1 {
    margin-left: 4px;
}

.ml-2 {
    margin-left: 8px;
}

.item-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: hidden;
}

.item-name {
    font-size: 13px;
    font-weight: 500;
    color: #303133;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.item-sub {
    font-size: 11px;
    color: #909399;
    margin-top: 2px;
}

.text-danger {
    color: #F56C6C;
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

.custom-tab-label {
    display: flex;
    align-items: center;
    gap: 4px;
}

.add-btn-block {
    width: 100%;
    border-style: dashed;
    margin-top: 4px;
    color: #409EFF;
}
</style>

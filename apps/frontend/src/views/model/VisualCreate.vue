<template>
    <div class="visual-model-builder">
        <!-- 顶部操作栏 -->
        <div class="builder-header">
            <div class="header-left">
                <el-button @click="handleBack" :icon="ArrowLeft" circle />
                <span class="title">可视化模型构建器</span>
            </div>
            <div class="header-actions">
                <el-button @click="addTable" type="primary" :icon="Plus">添加表</el-button>
                <el-button @click="handlePreviewSQL" :icon="View">预览 SQL</el-button>
                <el-button @click="handleFitView" :icon="FullScreen">适应画布</el-button>
                <el-button @click="handleAutoLayout" :icon="Rank">自动布局</el-button>
                <el-button-group class="ml-2">
                    <el-button @click="undo" :icon="RefreshLeft" title="撤销 (Ctrl+Z)" :disabled="historyIndex <= 0" />
                    <el-button @click="redo" :icon="RefreshRight" title="重做 (Ctrl+Shift+Z)" :disabled="historyIndex >= historyStack.length - 1" />
                </el-button-group>
                <el-button @click="handleAutoAlias" type="warning" :icon="Check">智能别名</el-button>
                <div class="divider"></div>
                <!-- 节点搜索 -->
                <el-select v-model="nodeSearchQuery" placeholder="搜索表..." filterable size="small" style="width: 150px" @change="handleNodeSearch" clearable>
                    <el-option v-for="node in elements.filter(el => el.type === 'table')" :key="node.id" :label="node.data.label" :value="node.id" />
                </el-select>
                <el-button-group class="ml-2">
                    <el-button @click="leftPanelVisible = !leftPanelVisible" :type="leftPanelVisible ? 'primary' : ''" :icon="View" title="显示/隐藏左面板" />
                    <el-button @click="rightPanelVisible = !rightPanelVisible" :type="rightPanelVisible ? 'primary' : ''" :icon="Setting" title="显示/隐藏右面板" />
                </el-button-group>
                <el-button @click="handleSave" type="success" :icon="Check" :loading="saving" class="ml-2">保存模型</el-button>
            </div>
        </div>
        <!-- 主容器 -->
        <div class="builder-container" ref="containerRef">
            <!-- 左栏：基础信息 -->
            <div v-if="leftPanelVisible" class="side-panel left-panel" :style="{
                top: leftPanelPos.y + 'px',
                left: leftPanelPos.x + 'px',
                width: leftPanelSize.width + 'px',
                height: leftCollapsed ? 'auto' : (leftPanelSize.height + 'px')
            }">
                <BaseInfoPanel :model="modelForm" :collapsed="leftCollapsed" @toggle-collapse="leftCollapsed = !leftCollapsed" @refresh-code="fetchGeneratedCode" @mousedown="startDrag($event, 'left')" />
                <!-- 调整大小手柄（右侧） -->
                <div v-if="!leftCollapsed" class="panel-resize-handle right" @mousedown="startResize($event, 'left')"></div>
                <!-- 底部调整大小手柄 -->
                <div v-if="!leftCollapsed" class="panel-resize-handle bottom" @mousedown="startResize($event, 'left')"></div>
            </div>
            <!-- 中栏：画布 -->
            <div class="canvas-panel">
                <VueFlow v-model="elements" :node-types="nodeTypes" @node-click="handleNodeClick" @edge-click="handleEdgeClick" @node-context-menu="onNodeContextMenu" @pane-click="closeMenu" @node-drag-start="closeMenu" @node-drag-stop="pushHistory" @edge-update-end="pushHistory" fit-view-on-init>
                    <Background pattern-color="#aaa" :gap="16" />
                    <Controls />
                    <MiniMap />
                </VueFlow>
                <!-- 右键菜单 -->
                <div v-show="menuVisible" class="context-menu" :style="{ left: menuPos.x + 'px', top: menuPos.y + 'px' }">
                    <div class="menu-item" @click="handleSetMainFromMenu">设为主表</div>
                    <div class="menu-item" @click="handleSmartSelectFromMenu">智能选字段</div>
                    <div class="menu-item delete" @click="handleDeleteFromMenu">删除节点</div>
                </div>
            </div>
            <!-- 右栏：配置面板 -->
            <div v-if="rightPanelVisible" class="side-panel right-panel" :style="{
                top: rightPanelPos.y + 'px',
                right: rightPanelPos.x + 'px',
                width: rightPanelSize.width + 'px',
                height: rightCollapsed ? 'auto' : (rightPanelSize.height + 'px')
            }">
                <FieldConfigPanel :elements="elements" :selected-element="selectedElement" :model-config="modelConfig" :collapsed="rightCollapsed" @toggle-collapse="rightCollapsed = !rightCollapsed" @mousedown="startDrag($event, 'right')" @remove-table="handleRemoveTable" @remove-field="handleRemoveField" />
                <!-- 调整大小手柄（左侧） -->
                <div v-if="!rightCollapsed" class="panel-resize-handle left" @mousedown="startResize($event, 'right')"></div>
                <!-- 底部调整大小手柄 -->
                <div v-if="!rightCollapsed" class="panel-resize-handle bottom" @mousedown="startResize($event, 'right')"></div>
            </div>
        </div>
        <!-- 表选择弹窗 -->
        <TableSelectDialog ref="tableSelectDialogRef" :conn-id="modelForm.connID?.toString()" @confirm="handleTableSelect" />
        <!-- SQL 预览弹窗 -->
        <el-dialog v-model="sqlPreviewVisible" title="SQL & 数据预览" width="80%">
            <el-tabs v-model="previewActiveTab">
                <el-tab-pane label="SQL 预览" name="sql">
                    <div class="sql-preview-content">
                        <div class="toolbar mb-2">
                            <el-button type="primary" :icon="VideoPlay" @click="handleExecuteSQL" :loading="executing">运行查询</el-button>
                        </div>
                        <el-input v-model="previewContent.sql" type="textarea" :rows="15" readonly resize="none" />
                        <div v-if="previewContent.args && previewContent.args.length > 0" class="sql-args mt-4">
                            <h4>参数:</h4>
                            <div class="args-list">
                                <span v-for="(arg, i) in previewContent.args" :key="i" class="arg-item">{{ arg }}</span>
                            </div>
                        </div>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="数据预览" name="data">
                    <div class="data-preview-content" v-loading="executing">
                        <el-table :data="previewContent.data" border height="500px" style="width: 100%">
                            <template v-if="previewContent.columns.length > 0">
                                <el-table-column v-for="col in previewContent.columns" :key="col" :prop="col" :label="col" min-width="120" show-overflow-tooltip />
                            </template>
                            <template v-else>
                                <el-table-column label="暂无数据" />
                            </template>
                        </el-table>
                        <div class="mt-2 text-gray-500">共 {{ previewContent.total }} 条结果</div>
                    </div>
                </el-tab-pane>
            </el-tabs>
            <template #footer>
                <el-button @click="sqlPreviewVisible = false">关闭</el-button>
                <el-button type="primary" @click="handleSaveFromPreview">保存模型</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { showConfirm } from '@/utils/confirm'
import { ArrowLeft, Check, FullScreen, Plus, Rank, RefreshLeft, RefreshRight, Setting, VideoPlay, View } from '@element-plus/icons-vue'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MarkerType, VueFlow, useVueFlow } from '@vue-flow/core'
import { MiniMap } from '@vue-flow/minimap'
import { ElMessage } from 'element-plus'
import { nextTick, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import TableNode from '@/components/model/TableNode.vue'
import BaseInfoPanel from './components/BaseInfoPanel.vue'
import FieldConfigPanel from './components/FieldConfigPanel.vue'
import TableSelectDialog from './components/TableSelectDialog.vue'

import { getFieldsByTableId } from '@/api/metadata'
import { generateModelCode, getModelById, previewVisualModelSQL, saveVisualModel } from '@/api/model'
import type { MdTable } from '@/types/metadata'

// Vue Flow 样式
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/minimap/dist/style.css'

const router = useRouter()
const route = useRoute()
const { onConnect, addEdges, removeNodes, removeEdges, fitView, getSelectedElements, setCenter } = useVueFlow()

// 定义自定义节点类型
const nodeTypes = {
    table: TableNode,
}

// 引用
const tableSelectDialogRef = ref()

// 模型数据
const modelForm = reactive({
    connID: undefined as number | undefined,
    modelName: '',
    modelCode: '',
    modelVersion: '1.0.0',
    remark: ''
})

// 面板大小固定为 600 * 600
// 移除硬编码的 reactive，后面在面板状态部分统一定义

// 模型高级配置 (Global Config)
const modelConfig = reactive({
    limit: 100,
    offset: 0,
    groupBy: [] as any[],
    havings: [] as any[],
    wheres: [] as any[],
    orders: [] as any[]
})

// 自动获取编码
const fetchGeneratedCode = async () => {
    // 只有在新建模式下才自动生成
    if (route.params.id) return
    try {
        const res: any = await generateModelCode()
        if (res.data && res.data.code) {
            modelForm.modelCode = res.data.code
        }
    } catch (error) {
        console.error('Failed to auto generate model code', error)
    }
}

// 面板状态
const leftCollapsed = ref(false)
const rightCollapsed = ref(false)

const leftPanelPos = reactive({ x: 20, y: 40 })
const rightPanelPos = reactive({ x: 20, y: 40 })
const leftPanelSize = reactive({ width: 300, height: 400 })
const rightPanelSize = reactive({ width: 600, height: 600 })
const containerRef = ref<HTMLElement | null>(null)

// Context Menu State
const menuVisible = ref(false)
const menuPos = reactive({ x: 0, y: 0 })
const contextNode = ref<any>(null)
const nodeSearchQuery = ref('')

// Undo/Redo History
const historyStack = ref<any[]>([])
const historyIndex = ref(-1)

const pushHistory = () => {
    const snapshot = {
        elements: JSON.parse(JSON.stringify(elements.value)),
        modelConfig: JSON.parse(JSON.stringify(modelConfig))
    }
    // Only push if different from last snapshot
    if (historyIndex.value >= 0) {
        const last = JSON.stringify(historyStack.value[historyIndex.value])
        if (last === JSON.stringify(snapshot)) return
    }

    if (historyIndex.value < historyStack.value.length - 1) {
        historyStack.value = historyStack.value.slice(0, historyIndex.value + 1)
    }
    historyStack.value.push(snapshot)
    if (historyStack.value.length > 50) {
        historyStack.value.shift()
    } else {
        historyIndex.value++
    }
}

const undo = () => {
    if (historyIndex.value > 0) {
        historyIndex.value--
        const state = historyStack.value[historyIndex.value]
        elements.value = JSON.parse(JSON.stringify(state.elements))
        Object.assign(modelConfig, JSON.parse(JSON.stringify(state.modelConfig)))
        ElMessage.info('已撤销')
    }
}

const redo = () => {
    if (historyIndex.value < historyStack.value.length - 1) {
        historyIndex.value++
        const state = historyStack.value[historyIndex.value]
        elements.value = JSON.parse(JSON.stringify(state.elements))
        Object.assign(modelConfig, JSON.parse(JSON.stringify(state.modelConfig)))
        ElMessage.info('已重做')
    }
}

// 面板展开状态
const leftPanelVisible = ref(true)
const rightPanelVisible = ref(true)

// 拖拽与调整大小逻辑
let draggingPanel = ref<string | null>(null)
let resizingPanel = ref<string | null>(null)
let startPos = { x: 0, y: 0 }
let initialPanelPos = { x: 0, y: 0 }
let initialPanelSize = { width: 0, height: 0 }

const startDrag = (e: MouseEvent, panel: string) => {
    const target = e.target as HTMLElement
    // 如果点击的是 resize-handle，则不触发拖拽
    if (target.classList.contains('resize-handle')) return
    if (!target.closest('.panel-header')) return

    draggingPanel.value = panel
    startPos = { x: e.clientX, y: e.clientY }
    const currentPos = panel === 'left' ? leftPanelPos : rightPanelPos
    initialPanelPos = { x: currentPos.x, y: currentPos.y }

    document.addEventListener('mousemove', handleDrag)
    document.addEventListener('mouseup', stopDrag)
}

const handleDrag = (e: MouseEvent) => {
    if (!draggingPanel.value) return

    const container = containerRef.value
    if (!container) return

    const containerRect = container.getBoundingClientRect()
    const containerWidth = containerRect.width
    const containerHeight = containerRect.height

    const dx = e.clientX - startPos.x
    const dy = e.clientY - startPos.y

    if (draggingPanel.value === 'left') {
        let newX = initialPanelPos.x + dx
        let newY = initialPanelPos.y + dy

        // Clamping Left Panel
        // x: [0, containerWidth - 300] (assuming width 300)
        // y: [0, containerHeight - header]
        newX = Math.max(0, Math.min(newX, containerWidth - 300))
        newY = Math.max(0, Math.min(newY, containerHeight - 50)) // simple clamp

        leftPanelPos.x = newX
        leftPanelPos.y = newY
    } else {
        // rightPanelPos is "right: x", "top: y"
        // dy is same (y increases downwards)
        // dx: if dragging right (positive dx), "right" value should decrease?
        // Wait, rightPanelPos.x is passed to style check:
        // right: rightPanelPos.x + 'px'
        // dragging right (increase mouse X) -> should decrease right value (closer to edge) -> wait, dragging right means moving TOWARDS right edge.
        // If I drag mouse right (+dx), the panel moves right. The distance from right edge (right property) decreases.
        // So: newRight = initialRight - dx

        let newX = initialPanelPos.x - dx
        let newY = initialPanelPos.y + dy

        // Clamping Right Panel
        newX = Math.max(0, Math.min(newX, containerWidth - 600)) // Width 600
        newY = Math.max(0, Math.min(newY, containerHeight - 50))

        rightPanelPos.x = newX
        rightPanelPos.y = newY
    }
}

const stopDrag = () => {
    draggingPanel.value = null
    document.removeEventListener('mousemove', handleDrag)
    document.removeEventListener('mouseup', stopDrag)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
}

// 调整大小逻辑
const startResize = (e: MouseEvent, panel: string) => {
    e.stopPropagation()
    resizingPanel.value = panel
    startPos = { x: e.clientX, y: e.clientY }
    const currentSize = panel === 'left' ? leftPanelSize : rightPanelSize
    initialPanelSize = { width: currentSize.width, height: currentSize.height || 0 }

    document.addEventListener('mousemove', handleResize)
    document.addEventListener('mouseup', stopResize)
    document.body.style.cursor = panel === 'left' ? 'e-resize' : 'w-resize'
    document.body.style.userSelect = 'none'
}

const handleResize = (e: MouseEvent) => {
    if (!resizingPanel.value) return

    const dx = e.clientX - startPos.x
    const dy = e.clientY - startPos.y

    if (resizingPanel.value === 'left') {
        const newWidth = initialPanelSize.width + dx
        if (newWidth > 200 && newWidth < 800) {
            leftPanelSize.width = newWidth
        }
    } else if (resizingPanel.value === 'right') {
        // 右边面板的手柄在左侧，向左拖拽（dx 为负）应当增加宽度
        const newWidth = initialPanelSize.width - dx
        const newHeight = initialPanelSize.height + dy

        if (newWidth > 300 && newWidth < 1000) {
            rightPanelSize.width = newWidth
        }
        if (newHeight > 200 && newHeight < 900) {
            rightPanelSize.height = newHeight
        }
    }
}

const stopResize = () => {
    resizingPanel.value = null
    document.removeEventListener('mousemove', handleResize)
    document.removeEventListener('mouseup', stopResize)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
}

// 画布元素
const elements = ref<any[]>([])
const selectedElement = ref<any>(null)

// SQL Preview State
const sqlPreviewVisible = ref(false)
const previewActiveTab = ref('sql')
const executing = ref(false)
const previewContent = reactive({
    sql: '',
    args: [] as any[],
    data: [] as any[],
    columns: [] as string[],
    total: 0
})
const saving = ref(false)

// 添加表入口
const addTable = () => {
    if (!modelForm.connID) {
        ElMessage.warning('请先选择数据源')
        return
    }
    tableSelectDialogRef.value?.show()
}

// 处理表选择确认
const handleTableSelect = async (selectedTables: MdTable[]) => {
    const newNodes: any[] = []

    for (const table of selectedTables) {
        // 检查是否已存在
        if (elements.value.find(el => el.id === table.id)) {
            continue
        }

        try {
            // 获取字段信息
            const fieldsRes = await getFieldsByTableId(table.id)
            const fields = (Array.isArray(fieldsRes) ? fieldsRes : (fieldsRes as any).data || []).map((f: any) => ({
                id: f.id,
                name: f.column_name,
                type: f.column_type,
                comment: f.column_comment,
                length: f.data_length,
                selected: false
            }))

            // 判断是否为主表（如果画布为空，则第一张选中的表为主表）
            const isMain = elements.value.length === 0 && newNodes.length === 0

            const newNode = {
                id: table.id,
                type: 'table',
                label: table.table_name,
                position: { x: 100 + elements.value.length * 50 + newNodes.length * 50, y: 100 + elements.value.length * 50 + newNodes.length * 50 },
                data: {
                    label: table.table_name,
                    schema: table.table_schema,
                    tableName: table.table_name,
                    isMain: isMain,
                    fields: fields
                }
            }
            newNodes.push(newNode)
        } catch (error) {
            console.error(`Failed to load fields for table ${table.table_name}`, error)
            ElMessage.error(`加载表 ${table.table_name} 字段失败`)
        }
    }

    if (newNodes.length > 0) {
        elements.value = [...elements.value, ...newNodes]
        nextTick(() => {
            fitView({ padding: 0.2, duration: 500 })
        })
        pushHistory()
    }
}

// 连线逻辑
onConnect((params) => {
    const edgeId = `e-${params.source}-${params.target}`

    // Check for Type Mismatch
    let edgeStyle = {}
    let edgeLabel = undefined
    let edgeAnimated = false
    let joinCondition = ''

    const sourceNode = elements.value.find(n => n.id === params.source)
    const targetNode = elements.value.find(n => n.id === params.target)

    // Explicit handle connection (Field to Field)
    if (sourceNode && targetNode && params.sourceHandle && params.targetHandle) {
        const sField = sourceNode.data.fields?.find((f: any) => f.name === params.sourceHandle)
        const tField = targetNode.data.fields?.find((f: any) => f.name === params.targetHandle)

        if (sField && tField) {
            // Very simple type check: string vs int/number
            const isString = (t: string) => /char|text|string/i.test(t)
            const isNum = (t: string) => /int|number|float|double|decimal/i.test(t)

            const sType = sField.type || ''
            const tType = tField.type || ''

            if ((isString(sType) && isNum(tType)) || (isNum(sType) && isString(tType))) {
                edgeStyle = { stroke: '#F56C6C', strokeWidth: 2 }
                edgeLabel = 'Type Mismatch'
                edgeAnimated = true
                ElMessage.warning(`类型警告: ${sField.name}(${sType}) 与 ${tField.name}(${tType}) 不匹配`)
            }
        }
    }

    if (sourceNode && targetNode) {
        // 策略1: 目标表有 source_id (例如: users.id = orders.user_id)
        // 从字段列表中查找
        // 字段存在于 node.data.fields
        const sFields = sourceNode.data.fields || []
        const tFields = targetNode.data.fields || []

        // 假设源表主键为 id (常见)
        const sId = sFields.find((f: any) => f.name.toLowerCase() === 'id')
        if (sId) {
            // 找目标表是否有 sourceName_id or sourceNameId
            const sName = sourceNode.data.tableName || sourceNode.data.label
            const targetKey = `${sName}_id`.toLowerCase()
            const tFk = tFields.find((f: any) => f.name.toLowerCase() === targetKey || f.name.toLowerCase() === sName.toLowerCase() + 'id')

            if (tFk) {
                joinCondition = `${sourceNode.data.label}.${sId.name} = ${targetNode.data.label}.${tFk.name}`
            }
        }

        // 策略2: 源表有 target_id (例如: orders.user_id = users.id)
        if (!joinCondition) {
            const tId = tFields.find((f: any) => f.name.toLowerCase() === 'id')
            if (tId) {
                const tName = targetNode.data.tableName || targetNode.data.label
                const sourceKey = `${tName}_id`.toLowerCase()
                const sFk = sFields.find((f: any) => f.name.toLowerCase() === sourceKey || f.name.toLowerCase() === tName.toLowerCase() + 'id')

                if (sFk) {
                    joinCondition = `${sourceNode.data.label}.${sFk.name} = ${targetNode.data.label}.${tId.name}`
                }
            }
        }

        // 策略3: 同名字段 (除了 id, create* update* 等通用字段)
        if (!joinCondition) {
            const common = sFields.filter((s: any) =>
                !['id', 'created_at', 'updated_at', 'create_by', 'update_by', 'is_deleted'].includes(s.name.toLowerCase()) &&
                tFields.some((t: any) => t.name === s.name)
            )
            if (common.length > 0) {
                // 取第一个
                joinCondition = `${sourceNode.data.label}.${common[0].name} = ${targetNode.data.label}.${common[0].name}`
            }
            if (!joinCondition && params.sourceHandle && params.targetHandle) {
                joinCondition = `${sourceNode.data.label}.${params.sourceHandle} = ${targetNode.data.label}.${params.targetHandle}`
            }
        }
    }

    // 提取自动推断出的字段
    let leftField = ''
    let rightField = ''
    const match = (joinCondition || '').match(/([^=<>!]+)\s*(=|!=|<>|>|<|>=|<=)\s*([^=<>!]+)/);
    if (match) {
        const leftParts = match[1].trim().split('.');
        const rightParts = match[3].trim().split('.');
        leftField = leftParts[leftParts.length - 1];
        rightField = rightParts[rightParts.length - 1];
    }

    const newEdge = {
        id: edgeId,
        ...params,
        type: 'default',
        markerEnd: MarkerType.ArrowClosed,
        style: edgeStyle,
        label: edgeLabel,
        animated: edgeAnimated,
        data: {
            joinType: 'LEFT JOIN',
            joinCondition: joinCondition,
            conditions: [
                { operator1: 'AND', leftField: leftField || (params.sourceHandle || ''), operator: '=', rightField: rightField || (params.targetHandle || '') }
            ]
        }
    }
    addEdges([newEdge])

    // 延迟一下确保 Vue Flow 内部状态更新后选中
    setTimeout(() => {
        const edge = elements.value.find(el => el.id === edgeId)
        if (edge) {
            selectedElement.value = edge
        }
    }, 50)
    pushHistory()
})

// 点击处理
const handleNodeClick = ({ node }: any) => {
    // 查找原引用以确保双向绑定生效
    const originalNode = elements.value.find(el => el.id === node.id)
    if (originalNode) {
        selectedElement.value = originalNode
    }
}

const handleEdgeClick = ({ edge }: any) => {
    const originalEdge = elements.value.find(el => el.id === edge.id)
    if (originalEdge) {
        selectedElement.value = originalEdge
    }
}


// 删除表处理
const handleRemoveTable = (tableId: string) => {
    showConfirm('确定要删除该表及其所有的字段配置和关联关系吗？').then(() => {
        // 1. 画布中删除节点及相关连线
        removeNodes([tableId])

        // 2. 状态中删除对应元素
        elements.value = elements.value.filter(el => el.id !== tableId)

        // 3. 如果是当前选中，清空选中态
        if (selectedElement.value?.id === tableId) {
            selectedElement.value = null
        }
    }).catch(() => { })
}

// 移除字段处理
const handleRemoveField = (tableId: string, fieldId: string) => {
    const node = elements.value.find(el => el.id === tableId)
    if (node && node.data.fields) {
        const field = node.data.fields.find((f: any) => f.id === fieldId)
        if (field) {
            field.selected = false
            pushHistory()
        }
    }
}

const handleBack = () => {
    router.push('/metadata/model/list')
}

const handleFitView = () => {
    fitView({ padding: 0.2, duration: 500 })
}

const handleAutoLayout = () => {
    const nodes = elements.value.filter(el => el.type === 'table')
    const edges = elements.value.filter(el => el.type === 'default')

    if (nodes.length === 0) return

    // 1. Find root (Main table or first table)
    let root = nodes.find(n => n.data.isMain) || nodes[0]

    // 2. BFS for levels
    const levels: Record<string, number> = {}
    const visited = new Set<string>()
    const queue: { id: string, level: number }[] = []

    queue.push({ id: root.id, level: 0 })
    levels[root.id] = 0
    visited.add(root.id)

    while (queue.length > 0) {
        const { id, level } = queue.shift()!

        // Find neighbors via edges
        // Only target nodes for now to flow left-to-right (assuming directed, but edges are undirected connections usually)
        // We check edges where this node is source OR target
        const neighbors = edges
            .filter(e => e.source === id || e.target === id)
            .map(e => e.source === id ? e.target : e.source)

        neighbors.forEach(nid => {
            if (!visited.has(nid)) {
                visited.add(nid)
                levels[nid] = level + 1
                queue.push({ id: nid, level: level + 1 })
            }
        })
    }

    // Handle disconnected nodes (put them at level 0 or separate?)
    // Put them at level 0 at bottom
    nodes.forEach(n => {
        if (!visited.has(n.id)) {
            levels[n.id] = 0
        }
    })

    // 3. Assign positions
    const LEVEL_WIDTH = 350
    const NODE_HEIGHT = 100 // Estimate
    const LEVEL_Gap = 50

    // Group by level
    const rows: Record<number, any[]> = {}
    Object.entries(levels).forEach(([id, level]) => {
        if (!rows[level]) rows[level] = []
        rows[level].push(nodes.find(n => n.id === id))
    })

    Object.keys(rows).forEach(lvlStr => {
        const lvl = Number(lvlStr)
        const items = rows[lvl]
        items.forEach((node, idx) => {
            node.position = {
                x: lvl * LEVEL_WIDTH + 50,
                y: idx * (NODE_HEIGHT + LEVEL_Gap) + 50
            }
        })
    })

    // Trigger update
    // elements.value = [...nodes, ...edges] // VueFlow handles reactive change?
    // It should work if we modify node.position directly inside elements array if it is reactive.
    // Let's force update

    nextTick(() => {
        fitView({ padding: 0.2, duration: 500 })
    })
    pushHistory()
}

// Context Menu Handlers
const onNodeContextMenu = (e: any) => {
    e.event.preventDefault()
    menuPos.x = e.event.clientX
    menuPos.y = e.event.clientY
    contextNode.value = e.node
    menuVisible.value = true
}

const closeMenu = () => {
    menuVisible.value = false
}

const handleSetMainFromMenu = () => {
    if (contextNode.value) {
        elements.value.forEach(el => {
            if (el.type === 'table') {
                el.data.isMain = (el.id === contextNode.value.id)
            }
        })
        pushHistory()
        closeMenu()
    }
}

const handleDeleteFromMenu = () => {
    if (contextNode.value) {
        removeNodes([contextNode.value])
        pushHistory()
        closeMenu()
    }
}

const handleNodeSearch = (nodeId: string) => {
    if (!nodeId) return
    const node = elements.value.find(el => el.id === nodeId)
    if (node) {
        // 聚焦节点
        setCenter(node.position.x + 100, node.position.y + 50, { zoom: 1.2, duration: 800 })

        // 闪烁效果 (临时样式)
        node.data.isSearching = true
        setTimeout(() => {
            node.data.isSearching = false
        }, 2000)

        // 选中该节点
        elements.value.forEach(el => el.selected = (el.id === nodeId))
    }
}

const handleAutoAlias = () => {
    const tableNodes = elements.value.filter(el => el.type === 'table')
    const fieldCounts: Record<string, number> = {}

    // 1. 统计选中的字段名频率
    tableNodes.forEach(node => {
        node.data.fields?.forEach((f: any) => {
            if (f.selected) {
                const name = (f.name || '').toLowerCase()
                fieldCounts[name] = (fieldCounts[name] || 0) + 1
            }
        })
    })

    // 2. 对于冲突字段，自动生成别名
    let resolvedCount = 0
    tableNodes.forEach(node => {
        node.data.fields?.forEach((f: any) => {
            if (f.selected) {
                const name = (f.name || '').toLowerCase()
                if (fieldCounts[name] > 1) {
                    if (!f.alias) {
                        const tableName = node.data.tableName || node.data.label
                        f.alias = `${tableName}_${f.name}`
                        resolvedCount++
                    }
                }
            }
        })
    })

    if (resolvedCount > 0) {
        ElMessage.success(`自动解决了 ${resolvedCount} 个字段名冲突`)
        pushHistory()
    } else {
        ElMessage.info('未发现明显的名称冲突')
    }
}

const handleSmartSelectFromMenu = () => {
    if (contextNode.value && contextNode.value.data && contextNode.value.data.fields) {
        const fields = contextNode.value.data.fields
        const keywords = /name|title|code|status|type|名称|标题|状态|类型/i
        const excludes = /^(id|create_|update_|is_del)/i

        let count = 0
        fields.forEach((f: any) => {
            const name = f.name.toLowerCase()
            const comment = (f.comment || '').toLowerCase()
            // 如果已经在排除列表，且未被手动选中？不管，智能选就是重置为智能状态
            // 策略：排除系统字段，包含关键词的选中。

            if (excludes.test(name)) {
                f.selected = false
            } else if (keywords.test(name) || keywords.test(comment)) {
                f.selected = true
                count++
            }
        })

        // 如果没选中任何字段，至少选中第一个非ID字段？
        if (count === 0 && fields.length > 0) {
            const firstUseful = fields.find((f: any) => !excludes.test(f.name))
            if (firstUseful) firstUseful.selected = true
        }

        ElMessage.success('已智能勾选推荐字段')
        pushHistory()
        closeMenu()
    }
}

// 快捷键处理
const onKeyDown = (e: KeyboardEvent) => {
    // Ignore if inputting text (e.target is input/textarea)
    const target = e.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') return

    // Delete / Backspace
    if (e.key === 'Delete' || e.key === 'Backspace') {
        const selected = getSelectedElements.value
        if (selected && selected.length > 0) {
            // 类型断言修复
            const nodes = selected.filter(el => !('source' in el))
            const edges = selected.filter(el => 'source' in el)

            if (nodes.length > 0) removeNodes(nodes as any)
            if (edges.length > 0) removeEdges(edges as any)
        }
    }

    // Undo: Ctrl+Z / Cmd+Z
    if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) {
        e.preventDefault()
        undo()
    }

    // Redo: Ctrl+Shift+Z / Cmd+Shift+Z or Ctrl+Y
    if ((e.ctrlKey || e.metaKey) && ((e.key === 'z' && e.shiftKey) || e.key === 'y')) {
        e.preventDefault()
        redo()
    }

    // Ctrl+S / Cmd+S
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault()
        handleSave()
    }

    // Ctrl+Enter / Cmd+Enter (Preview/Run)
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault()
        if (sqlPreviewVisible.value) {
            handleExecuteSQL()
        } else {
            handlePreviewSQL()
        }
    }
}

// 提取组装逻辑
const validateAndAssemble = () => {
    // 1. 校验必填项
    if (!modelForm.connID) {
        ElMessage.warning('请选择数据连接')
        return null
    }
    if (!modelForm.modelName) {
        ElMessage.warning('请输入模型名称')
        return null
    }
    if (!modelForm.modelCode) {
        ElMessage.warning('请输入模型编码')
        return null
    }

    // 2. 找到主表
    // 默认第一个表节点为主表，或者需要用户指定 (暂未实现指定主表UI, 假设逻辑是第一个)
    const tableNodes = elements.value.filter(el => el.type === 'table')
    if (tableNodes.length === 0) {
        ElMessage.warning('请至少添加一个表')
        return null
    }

    // 3. 组装 Tables & Fields
    const tables: any[] = []
    const fields: any[] = []

    tableNodes.forEach((node, index) => {
        const isMain = index === 0 // 简单逻辑：第一个是主表
        tables.push({
            id: node.id, // 使用节点ID作为关联ID (后端会重新生成如果需要，或者这里传 unique ID)
            // 注意：后端保存逻辑中，如果 ID 是 uuid 且已存在则更新。
            // 这里我们是新建，或者覆盖。
            // 使用临时 ID 只要保证关联正确即可。
            // node.id is uuid generated by VueFlow or us.
            // 传递给后端作为 TableID。
            is_main: isMain,
            table_schema: node.data.schema,
            table_name: node.data.tableName,
            table_alias: node.data.alias || '', // 别名支持?
            conn_id: modelForm.connID?.toString()
        })

        if (node.data.fields) {
            node.data.fields.forEach((f: any) => {
                // 只添加选中的字段? 还是全部?
                // 通常只添加选中的字段进入模型字段列表
                if (f.selected) {
                    fields.push({
                        model_id: '', // filled by backend
                        table_schema: node.data.schema,
                        table_name: node.data.tableName,
                        column_name: f.name,
                        column_title: f.comment || f.name,
                        show_title: f.comment || f.name, // 默认展示标题
                        field_type: f.type, // 需要映射?
                        // ...
                    })
                }
            })
        }
    })

    // Validate if any fields selected
    if (fields.length === 0) {
        ElMessage.warning('请至少选择一个字段')
        return null
    }

    // 4. 组装 Joins
    const joins: any[] = []
    for (const el of elements.value) {
        if (el.type === 'edge') {
            // Source & Target
            const sourceNode = elements.value.find(n => n.id === el.source)
            const targetNode = elements.value.find(n => n.id === el.target)
            if (!sourceNode || !targetNode) continue

            // 处理多个关联条件
            const conditions = el.data.conditions || []
            if (conditions.length > 0) {
                conditions.forEach((cond: any, index: number) => {
                    joins.push({
                        table_id: sourceNode.id,
                        table_schema: sourceNode.data.schema,
                        table_name: sourceNode.data.tableName,
                        join_table_id: targetNode.id,
                        join_table_schema: targetNode.data.schema,
                        join_table_name: targetNode.data.tableName,
                        join_type: el.data.joinType,
                        operator1: index === 0 ? 'AND' : (cond.operator1 || 'AND'),
                        operator2: cond.operator || '=',
                        column_name: cond.leftField,
                        join_column_name: cond.rightField,
                        remark: el.data.joinCondition // 保留以防万一
                    })
                })
            } else {
                // 如果没有条件，添加一个空的或者尝试解析 legacy joinCondition
                joins.push({
                    table_id: sourceNode.id,
                    table_schema: sourceNode.data.schema,
                    table_name: sourceNode.data.tableName,
                    join_table_id: targetNode.id,
                    join_table_schema: targetNode.data.schema,
                    join_table_name: targetNode.data.tableName,
                    join_type: el.data.joinType,
                    remark: el.data.joinCondition,
                    column_name: '',
                    join_column_name: '',
                    operator2: '='
                })

                const condition = el.data.joinCondition || '';
                const match = condition.match(/([^=<>!]+)\s*(=|!=|<>|>|<|>=|<=)\s*([^=<>!]+)/);
                if (match) {
                    const operator = match[2].trim();
                    const leftParts = match[1].trim().split('.');
                    const rightParts = match[3].trim().split('.');
                    const lastJoin = joins[joins.length - 1];
                    lastJoin.column_name = leftParts[leftParts.length - 1];
                    lastJoin.join_column_name = rightParts[rightParts.length - 1];
                    lastJoin.operator2 = operator;
                }
            }
        }
    }

    // 5. 组装 Wheres
    const wheres = (modelConfig.wheres || []).map((item: any, index: number) => {
        if (!item.field) return null
        const [nodeId, colName] = item.field.split('::')
        const node = elements.value.find(el => el.id === nodeId)
        if (!node) return null

        return {
            table_id: nodeId,
            table_schema: node.data.schema,
            table_name: node.data.tableName,
            column_name: colName,
            column_title: colName,
            operator1: index === 0 ? 'AND' : (item.operator1 || 'AND'),
            operator2: item.operator,
            value1: item.value
        }
    }).filter(Boolean)

    // 6. 组装 Orders
    const orders = (modelConfig.orders || []).map((item: any) => {
        if (!item.field) return null
        const [nodeId, colName] = item.field.split('::')
        const node = elements.value.find(el => el.id === nodeId)
        if (!node) return null

        return {
            table_id: nodeId,
            table_schema: node.data.schema,
            table_name: node.data.tableName,
            column_name: colName,
            order_type: item.direction
        }
    }).filter(Boolean)

    // 7. 组装 Groups
    const groups = (modelConfig.groupBy || []).map((item: any) => {
        if (!item.field) return null
        const [nodeId, colName] = item.field.split('::')
        const node = elements.value.find(el => el.id === nodeId)
        if (!node) return null

        return {
            table_id: nodeId,
            table_schema: node.data.schema,
            table_name: node.data.tableName,
            column_name: colName
        }
    }).filter(Boolean)

    // 8. 组装 Havings
    const havings = (modelConfig.havings || []).map((item: any, index: number) => {
        if (!item.field) return null
        const [nodeId, colName] = item.field.split('::')
        const node = elements.value.find(el => el.id === nodeId)
        if (!node) return null

        return {
            table_id: nodeId,
            table_schema: node.data.schema,
            table_name: node.data.tableName,
            column_name: colName === '*' ? '*' : colName,
            aggregate_func: item.func,
            operator1: index === 0 ? 'AND' : (item.operator1 || 'AND'),
            operator2: item.operator,
            value1: item.value
        }
    }).filter(Boolean)

    // 9. 组装请求
    const payload = {
        model_id: route.params.id ? route.params.id.toString() : undefined, // Edit mode
        conn_id: modelForm.connID.toString(),
        model_name: modelForm.modelName,
        model_code: modelForm.modelCode,
        model_version: modelForm.modelVersion,
        model_kind: 2, // 2 = Visual Model
        remark: modelForm.remark,
        is_public: false,
        parameters: JSON.stringify({
            visual: {
                elements: elements.value,
                config: modelConfig
            }
        }),
        tables,
        fields,
        joins,
        wheres,
        orders,
        groups,
        havings
    }
    return payload
}

const handleSave = async () => {
    const payload = validateAndAssemble()
    if (!payload) return

    saving.value = true
    try {
        const res: any = await saveVisualModel(payload)
        ElMessage.success('模型保存成功')

        // 如果是新建，保存后跳转到编辑模式
        if (!route.params.id && res.data && res.data.id) {
            router.replace(`/metadata/model/visual-edit/${res.data.id}`)
        }
    } catch (e) {
        console.error(e)
    } finally {
        saving.value = false
    }
}

const handlePreviewSQL = async () => {
    const payload = validateAndAssemble()
    if (!payload) return

    try {
        const res = await previewVisualModelSQL(payload, false)
        previewContent.sql = res.sql
        previewContent.args = res.args
        // 重置数据
        previewContent.data = []
        previewContent.columns = []
        previewContent.total = 0
        previewActiveTab.value = 'sql'
        sqlPreviewVisible.value = true
    } catch (e: any) {
        console.error(e)
        // 解析后端错误
        const msg = e.response?.data?.msg || e.message || '预览失败'
        if (msg.includes('ambiguous')) {
            ElMessage.error(`SQL生成错误: 字段不明确 (${msg})。请尝试为字段或表添加别名。`)
        } else {
            ElMessage.error(msg)
        }
    }
}

const handleExecuteSQL = async () => {
    const payload = validateAndAssemble()
    if (!payload) return

    executing.value = true
    previewActiveTab.value = 'data'

    try {
        const res = await previewVisualModelSQL(payload, true)
        if (res.data) {
            previewContent.data = res.data
            previewContent.total = res.total || 0
            if (res.data.length > 0) {
                previewContent.columns = Object.keys(res.data[0])
            } else {
                // Try to guess from fields if data is empty?
                // payload.fields has column_name / show_title
                previewContent.columns = payload.fields.map((f: any) => f.show_title || f.column_name)
            }
        }
    } catch (e) {
        console.error(e)
        ElMessage.error('执行查询失败')
        previewActiveTab.value = 'sql' // switch back to see error if needed?
    } finally {
        executing.value = false
    }
}

const handleSaveFromPreview = () => {
    sqlPreviewVisible.value = false
    handleSave()
}

onMounted(async () => {
    if (route.params.id) {
        // Edit Mode: Load Logic
        try {
            const res: any = await getModelById(route.params.id as string)
            if (res.data) {
                const model = res.data
                modelForm.modelName = model.model_name
                modelForm.modelCode = model.model_code
                modelForm.connID = Number(model.conn_id)
                modelForm.remark = model.remark
                modelForm.modelVersion = model.model_version

                // Restore visual state
                if (model.parameters) {
                    let params: any = {}
                    try {
                        params = JSON.parse(model.parameters)
                    } catch (e) {
                        console.error('Failed to parse parameters', e)
                    }

                    if (params.visual) {
                        if (params.visual.config) {
                            Object.assign(modelConfig, params.visual.config)
                        }
                        if (params.visual.elements) {
                            // Restore VueFlow elements
                            const restoredElements = params.visual.elements
                            elements.value = restoredElements
                            // Fit view after load
                            nextTick(() => {
                                fitView({ padding: 0.2 })
                            })
                        }
                    }
                }
            }
        } catch (e) {
            console.error(e)
            ElMessage.error('加载模型失败')
        }
    } else {
        fetchGeneratedCode()
    }

    // Initial snapshot
    nextTick(() => {
        pushHistory()
    })

    window.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
    window.removeEventListener('keydown', onKeyDown)
    document.removeEventListener('mousemove', handleDrag)
    document.removeEventListener('mouseup', stopDrag)
})
</script>
<style scoped>
.context-menu {
    position: fixed;
    background: #fff;
    border: 1px solid #e4e7ed;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    border-radius: 4px;
    padding: 5px 0;
    z-index: 3000;
    min-width: 120px;
}

.divider {
    width: 1px;
    height: 20px;
    background: #dcdfe6;
    margin: 0 10px;
}

.menu-item {
    padding: 8px 16px;
    font-size: 13px;
    color: #606266;
    cursor: pointer;
    transition: background-color 0.2s;
}

.menu-item:hover {
    background-color: #ecf5ff;
    color: #409eff;
}

.menu-item.table-node.is-main {
    border-color: #e6a23c;
}

.table-node.is-searching {
    border-color: #409eff;
    box-shadow: 0 0 10px #409eff;
    animation: flash 1s infinite alternate;
}

@keyframes flash {
    from {
        opacity: 1;
        transform: scale(1);
    }

    to {
        opacity: 0.8;
        transform: scale(1.05);
    }
}

.menu-item.delete {
    color: #f56c6c;
}

.menu-item.delete:hover {
    background-color: #fef0f0;
}

/* 连线动画相关 */
:deep(.vue-flow__edge-path) {
    stroke-width: 2;
    transition: stroke 0.3s, stroke-width 0.3s;
}

:deep(.vue-flow__edge.selected .vue-flow__edge-path) {
    stroke: #409eff;
    stroke-width: 3;
}

:deep(.vue-flow__edge.animated .vue-flow__edge-path) {
    stroke-dasharray: 5;
    animation: dashdraw 0.5s linear infinite;
}

@keyframes dashdraw {
    from {
        stroke-dashoffset: 10;
    }

    to {
        stroke-dashoffset: 0;
    }
}

.visual-model-builder {
    height: calc(100vh - 50px);
    /* 减去页眉高度 50px */
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background-color: var(--el-bg-color-page);
}

.builder-header {
    height: 50px;
    background-color: #fff;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 16px;
    flex-shrink: 0;
    z-index: 200;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 12px;
}

.header-left .title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
}

.header-actions {
    display: flex;
    gap: 10px;
}

.builder-container {
    flex: 1;
    position: relative;
    overflow: hidden;
}

.side-panel {
    position: absolute;
    width: 300px;
    z-index: 100;
    background-color: #fff;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    border-radius: 8px;
    overflow: hidden;
    height: auto;
    max-height: calc(100% - 80px);
    top: 40px;
    border: 1px solid #ebeef5;
    display: flex;
    flex-direction: column;
}

.left-panel {
    left: 20px;
}

.right-panel {
    right: 20px;
}

.canvas-panel {
    width: 100%;
    height: 100%;
    background-color: #fafafa;
}

.panel-resize-handle {
    position: absolute;
    background: transparent;
    z-index: 10;
}

.panel-resize-handle.right {
    top: 0;
    right: 0;
    width: 5px;
    height: 100%;
    cursor: e-resize;
}

.panel-resize-handle.left {
    top: 0;
    left: 0;
    width: 5px;
    height: 100%;
    cursor: w-resize;
}

.panel-resize-handle.bottom {
    bottom: 0;
    left: 0;
    width: 100%;
    height: 5px;
    cursor: s-resize;
}

.panel-resize-handle:hover {
    background: rgba(64, 158, 255, 0.3);
}

/* Vue Flow 容器样式 */
:deep(.vue-flow) {
    width: 100%;
    height: 100%;
}
</style>

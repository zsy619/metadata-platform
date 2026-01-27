<template>
    <div class="object-browser">
        <div class="browser-header">
            <el-input v-model="filterText" placeholder="搜索对象..." prefix-icon="Search" clearable />
            <el-button-group class="m-l-sm">
                <el-button :icon="Refresh" circle @click="refresh" title="刷新" />
            </el-button-group>
        </div>
        <div class="browser-content" v-loading="loading">
            <el-tree ref="treeRef" :data="treeData" :props="defaultProps" :filter-node-method="filterNode" node-key="id" highlight-current accordion @node-expand="handleNodeExpand" @node-click="handleNodeClick" class="filter-tree">
                <template #default="{ node, data }">
                    <span class="custom-tree-node">
                        <el-icon class="m-r-xs">
                            <component :is="getNodeIcon(data.type)" />
                        </el-icon>
                        <span :title="data.label">{{ node.label }}</span>
                        <span v-if="data.count" class="node-count">({{ data.count }})</span>
                    </span>
                </template>
            </el-tree>
        </div>
    </div>
</template>
<script setup lang="ts">
import { getConns, getDBTables, getDBViews } from '@/api/metadata'
import { Coin, Connection, Document, Folder, Refresh, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref, watch } from 'vue'

const props = defineProps({
    dataSourceId: {
        type: String,
        default: ''
    }
})

const emit = defineEmits(['select-table', 'select-view'])

const filterText = ref('')
const treeRef = ref()
const loading = ref(false)
const treeData = ref<any[]>([])

const defaultProps = {
    children: 'children',
    label: 'label',
    isLeaf: 'leaf'
}

watch(filterText, (val) => {
    treeRef.value!.filter(val)
})

const filterNode = (value: string, data: any) => {
    if (!value) return true
    return data.label.includes(value)
}

const getNodeIcon = (type: string) => {
    switch (type) {
        case 'root': return Folder
        case 'conn': return Connection
        case 'schema': return Folder
        case 'table-folder': return Folder
        case 'view-folder': return Folder
        case 'table': return Coin
        case 'view': return View
        default: return Document
    }
}

// 加载数据源列表（如果未指定ID）或加载指定数据源的Schema/Tables
const loadData = async () => {
    loading.value = true
    try {
        if (!props.dataSourceId) {
            // 加载所有数据源
            const res: any = await getConns()
            const conns = res.data || res
            treeData.value = conns.map((conn: any) => ({
                id: `conn_${conn.id}`,
                label: conn.conn_name,
                type: 'conn',
                data: conn,
                leaf: false
            }))
        } else {
            // 已指定数据源，直接加载该数据源的结构（通常作为子组件使用时）
            // 这里简化逻辑，暂只处理从根加载的情况，如果需要单独使用，可扩展
        }
    } catch (error) {
        console.error('加载失败', error)
        ElMessage.error('加载失败')
    } finally {
        loading.value = false
    }
}

const handleNodeExpand = async (data: any, node: any) => {
    if (node.loaded) return

    // 如果是数据源节点，加载Schema或直接加载Table/View(视数据库类型而定)
    if (data.type === 'conn') {
        node.loading = true
        try {
            // 简单处理：直接分为 Table 和 View 两个文件夹
            // 实际可能需要处理Schema，这里暂简化
            node.data.children = [
                { id: `${data.id}_tables`, label: 'Tables', type: 'table-folder', parentId: data.data.id, leaf: false },
                { id: `${data.id}_views`, label: 'Views', type: 'view-folder', parentId: data.data.id, leaf: false }
            ]
            node.loaded = true
            // 手动触发展开效果
            treeRef.value.append(node.data.children[0], node)
            treeRef.value.append(node.data.children[1], node)
            // Remove original placeholder if any
        } catch (error) {
            // ...
        } finally {
            node.loading = false
        }
    } else if (data.type === 'table-folder') {
        node.loading = true
        try {
            // Load Tables
            const res: any = await getDBTables(data.parentId)
            const tables = res.data || res
            // 替换children
            const children = tables.map((t: string) => ({
                id: `${data.id}_${t}`,
                label: t,
                type: 'table',
                connId: data.parentId,
                leaf: true
            }))
            treeRef.value.updateKeyChildren(data.id, children)
            node.loaded = true
        } catch (error: any) {
            ElMessage.error('加载表失败: ' + error.message)
        } finally {
            node.loading = false
        }
    } else if (data.type === 'view-folder') {
        node.loading = true
        try {
            // Load Views
            const res: any = await getDBViews(data.parentId)
            const views = res.data || res
            const children = views.map((v: string) => ({
                id: `${data.id}_${v}`,
                label: v,
                type: 'view',
                connId: data.parentId,
                leaf: true
            }))
            treeRef.value.updateKeyChildren(data.id, children)
            node.loaded = true
        } catch (error: any) {
            ElMessage.error('加载视图失败: ' + error.message)
        } finally {
            node.loading = false
        }
    }
}

const handleNodeClick = (data: any) => {
    if (data.type === 'table') {
        emit('select-table', { connId: data.connId, tableName: data.label })
    } else if (data.type === 'view') {
        emit('select-view', { connId: data.connId, viewName: data.label })
    }
}

const refresh = () => {
    loadData()
}

onMounted(() => {
    loadData()
})
</script>
<style scoped>
.object-browser {
    height: 100%;
    display: flex;
    flex-direction: column;
    border-right: 1px solid var(--el-border-color-light);
    background-color: var(--el-bg-color);
}

.browser-header {
    padding: 10px;
    display: flex;
    align-items: center;
    border-bottom: 1px solid var(--el-border-color-light);
}

.browser-content {
    flex: 1;
    overflow: auto;
    padding: 10px 0;
}

.custom-tree-node {
    display: flex;
    align-items: center;
    font-size: 14px;
}

.node-count {
    color: var(--el-text-color-secondary);
    font-size: 12px;
    margin-left: 4px;
}

.m-l-sm {
    margin-left: 8px;
}

.m-r-xs {
    margin-right: 4px;
}
</style>

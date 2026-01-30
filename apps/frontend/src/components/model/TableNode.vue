<template>
    <div class="table-node" :class="{ 'is-main': data.isMain, 'is-searching': data.isSearching }">
        <div class="table-header">
            <el-icon>
                <Grid />
            </el-icon>
            <span class="table-name">{{ data.tableTitle || data.tableName || data.label }}</span>
            <el-tag v-if="data.isMain" size="small" type="warning" effect="dark">主</el-tag>
        </div>
        <div class="table-fields">
            <div v-for="field in data.fields" :key="field.id" class="field-item-wrapper" @click.stop="toggleField(field)">
                <el-tooltip effect="dark" placement="right" :show-after="500">
                    <template #content>
                        <div>类型: {{ field.type }}</div>
                        <div v-if="field.comment">注释: {{ field.comment }}</div>
                    </template>
                    <div class="field-item">
                        <el-checkbox v-model="field.selected" @click.stop />
                        <span class="field-name">{{ field.name }}</span>
                        <!-- <span class="field-type">{{ field.type }}</span> -->
                    </div>
                </el-tooltip>
            </div>
        </div>
        <Handle type="source" :position="Position.Right" />
        <Handle type="target" :position="Position.Left" />
    </div>
</template>
<script setup lang="ts">
import { Grid } from '@element-plus/icons-vue';
import { Handle, type NodeProps, Position } from '@vue-flow/core';

interface TableData {
    label: string
    tableName: string
    tableTitle: string
    tableAlias: string
    isMain: boolean
    isSearching?: boolean
    fields: Array<{ id: string; name: string; type: string; selected: boolean; comment?: string }>
}

const props = defineProps<NodeProps<TableData>>()

const toggleField = (field: any) => {
    field.selected = !field.selected
}
</script>
<style scoped>
.table-node {
    background: white;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    min-width: 180px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.table-node.is-main {
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

.table-header {
    padding: 8px 12px;
    background: #f5f7fa;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: bold;
}

.table-name {
    flex: 1;
    font-size: 13px;
    color: #303133;
}

.table-fields {
    padding: 4px 0;
    max-height: 200px;
    overflow-y: auto;
}

.field-item {
    padding: 4px 12px;
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
}

.field-item:hover {
    background: #f0f2f5;
}

.field-name {
    flex: 1;
    font-size: 12px;
    color: #606266;
}

.field-type {
    font-size: 11px;
    color: #909399;
}
</style>

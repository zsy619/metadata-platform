<template>
    <div class="tree-name-wrapper">
        <el-icon class="tree-type-icon" :size="16">
            <Folder v-if="hasChildren" color="#409EFF" />
            <Document v-else color="#909399" />
        </el-icon>
        <span class="tree-name-text">{{ name }}</span>
    </div>
</template>

<script setup lang="ts">
import { Document, Folder } from '@element-plus/icons-vue';
import { computed } from 'vue';

const props = defineProps<{
    row: any
    nameField?: string
}>()

const nameField = computed(() => props.nameField || 'name')

const name = computed(() => props.row?.[nameField.value] || '')

const hasChildren = computed(() => {
    const children = props.row?.children
    return children && Array.isArray(children) && children.length > 0
})
</script>

<style scoped>
/* 关键：让展开按钮和自定义内容在同一行 */
::v-deep(.el-table__cell) .cell {
    display: flex;
    align-items: center;
}

/* 展开按钮容器 */
::v-deep(.el-table__expand-icon) {
    display: flex;
    align-items: center;
    height: 20px;
    margin-right: 4px;
}

/* 叶子节点（无子节点）的缩进对齐 */
::v-deep(.el-table__placeholder) {
    width: 20px;
    margin-right: 4px;
}
.tree-name-wrapper {
    display: flex;
    align-items: center;
    flex: 1;
}

.tree-type-icon {
    margin-right: 6px;
    display: flex;
    align-items: center;
}

.tree-name-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>

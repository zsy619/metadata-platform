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

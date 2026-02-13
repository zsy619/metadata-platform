<template>
    <div class="data-table-wrapper">
        <el-table ref="tableRef" v-loading="loading" :data="tableData" :stripe="stripe" :border="border" :size="size" :height="height" :max-height="maxHeight" :row-key="rowKey" :expand-row-keys="expandedRowKeys" @selection-change="handleSelectionChange" @sort-change="handleSortChange" @row-click="handleRowClick" @expand-change="handleExpandChange" class="data-table">
            <el-table-column v-if="showSelection" type="selection" width="50" align="center" :reserve-selection="reserveSelection" />
            <el-table-column v-if="showIndex" type="index" label="序号" width="60" align="center" :index="indexMethod" />
            <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" :min-width="col.minWidth" :align="col.align || 'left'" :fixed="col.fixed" :sortable="col.sortable" :formatter="col.formatter" :show-overflow-tooltip="col.showOverflowTooltip !== false" :resizable="col.resizable !== false">
                <template #header v-if="col.required || col.rules">
                    <span class="required-mark" v-if="col.required">* </span>{{ col.label }}
                </template>
                <template #default="scope">
                    <slot :name="col.slot || col.prop" :row="scope.row" :column="scope.column" :$index="scope.$index" :value="scope.row[col.prop]">
                        <template v-if="col.formatter">
                            {{ col.formatter(scope.row, scope.column, scope.row[col.prop], scope.$index) }}
                        </template>
                        <template v-else-if="col.boolTag">
                            <el-tag :type="scope.row[col.prop] ? col.boolTag.trueType || 'success' : col.boolTag.falseType || 'info'" :effect="col.boolTag.effect || 'light'">
                                {{ scope.row[col.prop] ? (col.boolTag.trueText || '是') : (col.boolTag.falseText || '否') }}
                            </el-tag>
                        </template>
                        <template v-else>
                            {{ scope.row[col.prop] ?? '-' }}
                        </template>
                    </slot>
                </template>
            </el-table-column>
            <el-table-column v-if="showAction" label="操作" :width="actionWidth" :fixed="actionFixed" align="center">
                <template #default="scope">
                    <slot name="action" :row="scope.row" :$index="scope.$index"></slot>
                </template>
            </el-table-column>
        </el-table>
        <div v-if="showPagination" class="pagination-wrapper">
            <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="pageSizes" :total="total" :layout="paginationLayout" :background="paginationBackground" :small="paginationSmall" @current-change="handlePageChange" @size-change="handleSizeChange" />
        </div>
    </div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { ElTable, TableColumnCtx } from 'element-plus'

interface Column {
    prop: string
    label: string
    width?: string | number
    minWidth?: string | number
    align?: 'left' | 'center' | 'right'
    fixed?: boolean | 'left' | 'right'
    sortable?: boolean | 'custom'
    formatter?: (row: any, column: TableColumnCtx, cellValue: any, index: number) => any
    showOverflowTooltip?: boolean
    resizable?: boolean
    slot?: string
    required?: boolean
    rules?: any[]
    boolTag?: {
        trueType?: string
        falseType?: string
        trueText?: string
        falseText?: string
        effect?: 'light' | 'dark' | 'plain'
    }
}

interface Props {
    data?: any[]
    columns?: Column[]
    loading?: boolean
    stripe?: boolean
    border?: boolean
    size?: 'large' | 'default' | 'small'
    height?: string | number
    maxHeight?: string | number
    rowKey?: string
    reserveSelection?: boolean
    showSelection?: boolean
    showIndex?: boolean
    showAction?: boolean
    actionWidth?: string | number
    actionFixed?: boolean | 'left' | 'right'
    showPagination?: boolean
    total?: number
    page?: number
    limit?: number
    pageSizes?: number[]
    paginationLayout?: string
    paginationBackground?: boolean
    paginationSmall?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    data: () => [],
    columns: () => [],
    loading: false,
    stripe: true,
    border: true,
    size: 'default',
    rowKey: 'id',
    reserveSelection: true,
    showSelection: false,
    showIndex: false,
    showAction: false,
    actionWidth: 150,
    actionFixed: 'right',
    showPagination: false,
    total: 0,
    page: 1,
    limit: 20,
    pageSizes: () => [10, 20, 50, 100],
    paginationLayout: 'total, sizes, prev, pager, next, jumper',
    paginationBackground: true,
    paginationSmall: false
})

const emit = defineEmits(['update:page', 'update:limit', 'selection-change', 'sort-change', 'row-click', 'expand-change', 'refresh'])

const tableRef = ref<InstanceType<typeof ElTable>>()
const currentPage = ref(props.page)
const pageSize = ref(props.limit)
const expandedRowKeys = ref<string[]>([])

const tableData = computed(() => props.data)

watch(() => props.page, (val) => {
    currentPage.value = val
})

watch(() => props.limit, (val) => {
    pageSize.value = val
})

const indexMethod = (index: number) => {
    return (currentPage.value - 1) * pageSize.value + index + 1
}

const handleSelectionChange = (selection: any[]) => {
    emit('selection-change', selection)
}

const handleSortChange = ({ prop, order }: any) => {
    emit('sort-change', { prop, order })
}

const handleRowClick = (row: any) => {
    emit('row-click', row)
}

const handleExpandChange = (row: any, expanded: boolean) => {
    emit('expand-change', row, expanded)
}

const handlePageChange = (page: number) => {
    emit('update:page', page)
    emit('refresh', { page, limit: pageSize.value })
}

const handleSizeChange = (size: number) => {
    pageSize.value = size
    emit('update:limit', size)
    emit('refresh', { page: currentPage.value, limit: size })
}

const toggleRowExpansion = (row: any, expanded?: boolean) => {
    tableRef.value?.toggleRowExpansion(row, expanded)
}

const toggleRowSelection = (row: any, selected?: boolean) => {
    tableRef.value?.toggleRowSelection(row, selected)
}

const clearSelection = () => {
    tableRef.value?.clearSelection()
}

const clearSort = () => {
    tableRef.value?.clearSort()
}

const setCurrentRow = (row: any) => {
    tableRef.value?.setCurrentRow(row)
}

const expandRow = (rowKey: string) => {
    if (!expandedRowKeys.value.includes(rowKey)) {
        expandedRowKeys.value.push(rowKey)
    }
}

const collapseRow = (rowKey: string) => {
    const index = expandedRowKeys.value.indexOf(rowKey)
    if (index > -1) {
        expandedRowKeys.value.splice(index, 1)
    }
}

defineExpose({
    tableRef,
    toggleRowExpansion,
    toggleRowSelection,
    clearSelection,
    clearSort,
    setCurrentRow,
    expandRow,
    collapseRow
})
</script>
<style scoped>
.data-table-wrapper {
    width: 100%;
}

.data-table {
    width: 100%;
}

.pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding: 16px 0;
}

.required-mark {
    color: #f56c6c;
    margin-right: 2px;
}

:deep(.el-table__header th) {
    background-color: #f5f7fa;
}

:deep(.el-table__body tr:hover > td) {
    background-color: #f5f7fa !important;
}

:deep(.el-table .el-table__cell) {
    padding: 10px 0;
}
</style>

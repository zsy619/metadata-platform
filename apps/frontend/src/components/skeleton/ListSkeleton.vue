<template>
    <div class="list-skeleton">
        <!-- 搜索栏骨架 -->
        <div class="search-skeleton" v-if="showSearch">
            <el-skeleton-item variant="rect" style="width: 200px; height: 32px; margin-right: 10px" />
            <el-skeleton-item variant="rect" style="width: 120px; height: 32px; margin-right: 10px" />
            <el-skeleton-item variant="button" style="width: 80px; height: 32px" />
        </div>
        <!-- 表格骨架 -->
        <div class="table-skeleton">
            <!-- 表头 -->
            <div style="display: flex; margin-bottom: 10px; border-bottom: 1px solid #ebeef5; padding-bottom: 10px;">
                <el-skeleton-item variant="text" style="width: 10%; margin-right: 2%" />
                <el-skeleton-item variant="text" style="width: 20%; margin-right: 2%" />
                <el-skeleton-item variant="text" style="width: 30%; margin-right: 2%" />
                <el-skeleton-item variant="text" style="width: 15%; margin-right: 2%" />
                <el-skeleton-item variant="text" style="width: 15%;" />
            </div>
            <!-- 表格行 -->
            <el-skeleton :rows="rows" animated />
        </div>
    </div>
</template>
<script setup lang="ts">
defineProps({
    rows: {
        type: Number,
        default: 10
    },
    showSearch: {
        type: Boolean,
        default: true
    }
})
</script>
<style scoped>
.list-skeleton {
    padding: 20px;
    background: #fff;
}

.search-skeleton {
    display: flex;
    margin-bottom: 20px;
}

.table-skeleton {
    border: 1px solid #ebeef5;
    padding: 20px;
    border-radius: 4px;
}
</style>

<template>
    <div class="result-container">
        <el-result icon="warning" title="503" sub-title="服务暂时不可用，系统可能正在维护中">
            <template #extra>
                <el-button type="primary" @click="handleRetry">重试连接</el-button>
            </template>
        </el-result>
    </div>
</template>
<script setup lang="ts">
const handleRetry = () => {
    window.location.reload()
}
</script>
<style scoped>
.result-container {
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: #fff;
}
</style>

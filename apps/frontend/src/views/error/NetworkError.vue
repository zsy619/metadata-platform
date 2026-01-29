<template>
    <div class="result-container">
        <el-result icon="error" title="网络连接失败" sub-title="请检查您的网络设置或稍后重试">
            <template #extra>
                <el-button type="primary" @click="handleRetry">刷新页面</el-button>
                <el-button @click="handleBack">返回上页</el-button>
            </template>
        </el-result>
    </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const handleRetry = () => {
    window.location.reload()
}

const handleBack = () => {
    router.back()
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

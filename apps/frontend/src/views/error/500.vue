<template>
    <div class="server-error">
        <el-result icon="error" title="500" sub-title="抱歉，服务器出错了">
            <template #extra>
                <el-button type="primary" @click="goHome">返回首页</el-button>
                <el-button @click="retry">重试</el-button>
            </template>
        </el-result>
    </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'

defineOptions({
    name: 'ServerError'
})

const router = useRouter()

const goHome = () => {
    router.push('/')
}

const retry = () => {
    window.location.reload()
}
</script>
<style scoped>
.server-error {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
}
</style>

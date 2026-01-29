<template>
    <div class="forbidden">
        <el-result icon="warning" title="403" sub-title="抱歉，您没有权限访问此页面">
            <template #extra>
                <el-button type="primary" @click="goHome">返回首页</el-button>
                <el-button @click="goBack">返回上页</el-button>
            </template>
        </el-result>
    </div>
</template>
<script setup lang="ts">
import { useRouter } from 'vue-router'

defineOptions({
    name: 'Forbidden'
})

const router = useRouter()

const goHome = () => {
    router.push('/')
}

const goBack = () => {
    router.back()
}
</script>
<style scoped>
.forbidden {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
}
</style>

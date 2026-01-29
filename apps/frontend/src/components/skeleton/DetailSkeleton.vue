<template>
    <div class="detail-skeleton">
        <!-- 头部信息 -->
        <div class="header-skeleton">
            <el-skeleton animated>
                <template #template>
                    <el-skeleton-item variant="h1" style="width: 30%; margin-bottom: 20px" />
                    <div style="display: flex; gap: 20px; margin-bottom: 20px">
                        <el-skeleton-item variant="text" style="width: 15%" />
                        <el-skeleton-item variant="text" style="width: 15%" />
                        <el-skeleton-item variant="text" style="width: 15%" />
                    </div>
                    <el-skeleton-item variant="text" style="width: 60%; margin-bottom: 40px" />
                </template>
            </el-skeleton>
        </div>
        <!-- 内容区域 -->
        <div class="content-skeleton">
            <el-skeleton :rows="5" animated />
        </div>
    </div>
</template>
<style scoped>
.detail-skeleton {
    padding: 20px;
    background: #fff;
    height: 100%;
}
</style>

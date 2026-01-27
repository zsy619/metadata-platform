<template>
    <page-container title="仪表盘" description="系统运行状态概览">
        <template #header>
            <el-button type="primary" :icon="Refresh" @click="handleRefresh">刷新数据</el-button>
        </template>
        <!-- 统计卡片 -->
        <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="6" v-for="item in stats" :key="item.title">
                <el-card shadow="hover" class="stat-card">
                    <div class="stat-content">
                        <div class="stat-icon" :style="{ backgroundColor: item.color }">
                            <el-icon>
                                <component :is="item.icon" />
                            </el-icon>
                        </div>
                        <div class="stat-info">
                            <div class="stat-value">{{ item.value }}</div>
                            <div class="stat-title">{{ item.title }}</div>
                        </div>
                    </div>
                </el-card>
            </el-col>
        </el-row>
        <!-- 快捷操作与最近活动 -->
        <el-row :gutter="20" class="m-t-20">
            <el-col :xs="24" :lg="16">
                <el-card shadow="hover" header="快捷操作">
                    <div class="quick-actions">
                        <el-button type="default" size="large" class="action-btn" @click="$router.push('/metadata/datasource/create')">
                            <el-icon>
                                <Connection />
                            </el-icon>
                            <span>新建数据源</span>
                        </el-button>
                        <el-button type="default" size="large" class="action-btn" @click="$router.push('/model/create')">
                            <el-icon>
                                <DocumentAdd />
                            </el-icon>
                            <span>创建模型</span>
                        </el-button>
                        <el-button type="default" size="large" class="action-btn" @click="$router.push('/api/create')">
                            <el-icon>
                                <MagicStick />
                            </el-icon>
                            <span>发布接口</span>
                        </el-button>
                        <el-button type="default" size="large" class="action-btn" @click="$router.push('/system/users')">
                            <el-icon>
                                <UserFilled />
                            </el-icon>
                            <span>用户管理</span>
                        </el-button>
                    </div>
                </el-card>
            </el-col>
            <el-col :xs="24" :lg="8">
                <el-card shadow="hover" header="系统公告">
                    <el-timeline>
                        <el-timeline-item timestamp="2026-01-27" placement="top">
                            <el-card>
                                <h4>系统升级通知</h4>
                                <p>元数据平台 v1.0.0 正式发布</p>
                            </el-card>
                        </el-timeline-item>
                        <el-timeline-item timestamp="2026-01-20" placement="top">
                            <p>完成基础框架搭建</p>
                        </el-timeline-item>
                    </el-timeline>
                </el-card>
            </el-col>
        </el-row>
    </page-container>
</template>
<script setup lang="ts">
import PageContainer from '@/components/common/PageContainer.vue'
import {
    Connection,
    DataLine,
    Document as DocumentAdd,
    Files,
    MagicStick,
    Refresh,
    Share,
    User,
    UserFilled
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { markRaw, ref } from 'vue'

const stats = ref([
    { title: '数据源', value: 12, icon: markRaw(DataLine), color: '#409EFF' },
    { title: '数据模型', value: 45, icon: markRaw(Files), color: '#67C23A' },
    { title: 'API接口', value: 128, icon: markRaw(Share), color: '#E6A23C' },
    { title: '注册用户', value: 89, icon: markRaw(User), color: '#F56C6C' },
])


const handleRefresh = () => {
    ElMessage.success('数据已刷新')
}
</script>
<style scoped>
.m-t-20 {
    margin-top: 20px;
}

.stat-card {
    margin-bottom: 20px;
    cursor: pointer;
    transition: transform 0.3s;
}

.stat-card:hover {
    transform: translateY(-5px);
}

.stat-content {
    display: flex;
    align-items: center;
}

.stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 20px;
    color: #fff;
    font-size: 24px;
}

.stat-info {
    flex: 1;
}

.stat-value {
    font-size: 24px;
    font-weight: bold;
    color: #303133;
}

.stat-title {
    font-size: 14px;
    color: #909399;
    margin-top: 5px;
}

.quick-actions {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
}

.action-btn {
    display: flex;
    flex-direction: column;
    height: auto;
    padding: 20px;
    width: 120px;
    gap: 10px;
}

.action-btn :deep(span) {
    margin-left: 0 !important;
    margin-top: 5px;
}
</style>

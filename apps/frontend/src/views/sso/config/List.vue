<template>
    <div class="sso-page">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Key />
                </el-icon>
                SSO单点登录配置
            </h1>
        </div>
        
        <el-card class="main-card">
            <el-tabs v-model="activeTab" type="card" class="config-tabs">
                <el-tab-pane label="协议配置" name="protocol">
                    <ProtocolConfig />
                </el-tab-pane>
                <el-tab-pane label="客户端配置" name="client">
                    <ClientConfig />
                </el-tab-pane>
                <el-tab-pane label="密钥管理" name="key">
                    <KeyManager />
                </el-tab-pane>
                <el-tab-pane label="字段映射" name="mapping">
                    <FieldMapping />
                </el-tab-pane>
                <el-tab-pane label="会话管理" name="session">
                    <SessionManager />
                </el-tab-pane>
            </el-tabs>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import { Key } from '@element-plus/icons-vue'
import { ref } from 'vue'
import ProtocolConfig from './ProtocolConfig.vue'
import ClientConfig from './ClientConfig.vue'
import KeyManager from './KeyManager.vue'
import FieldMapping from './FieldMapping.vue'
import SessionManager from './SessionManager.vue'

const activeTab = ref('protocol')
</script>

<style scoped>
.sso-page {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

:deep(.el-card__body) {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow: hidden;
    box-sizing: border-box;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.page-title {
    font-size: 20px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 8px;
}

.title-icon {
    font-size: 24px;
    color: var(--el-color-primary);
}

.config-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.config-tabs :deep(.el-tabs__content) {
    flex: 1;
    overflow: hidden;
}

.config-tabs :deep(.el-tab-pane) {
    height: 100%;
}

.text-primary {
    color: var(--el-text-color-primary);
}
</style>

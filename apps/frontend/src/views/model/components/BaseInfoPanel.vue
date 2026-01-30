<template>
    <div class="base-info-panel" :class="{ 'is-collapsed': collapsed }">
        <div class="panel-header" @click="$emit('toggle-collapse')">
            <div class="header-content">
                <el-icon>
                    <InfoFilled />
                </el-icon>
                <span>基础信息</span>
            </div>
            <el-icon class="collapse-icon">
                <ArrowUp v-if="!collapsed" />
                <ArrowDown v-else />
            </el-icon>
        </div>
        <div v-show="!collapsed" class="panel-body">
            <el-form :model="model" label-position="top" size="small">
                <el-form-item label="数据源" required>
                    <el-select v-model="model.connID" placeholder="请选择数据源" style="width: 100%" @change="handleConnChange">
                        <el-option v-for="item in dataSources" :key="item.id" :label="item.conn_name" :value="item.id" />
                    </el-select>
                </el-form-item>
                <el-form-item label="模型名称" required>
                    <el-input v-model="model.modelName" placeholder="请输入模型名称" />
                </el-form-item>
                <el-form-item label="模型编码" required>
                    <el-input v-model="model.modelCode" placeholder="自动生成或手动输入">
                        <template #append>
                            <el-button :icon="Refresh" @click="$emit('refresh-code')" title="重新获取编码" />
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item label="版本">
                    <el-input v-model="model.modelVersion" placeholder="1.0.0" />
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="model.remark" type="textarea" :rows="4" placeholder="模型描述信息" />
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>
<script setup lang="ts">
import { getConns } from '@/api/metadata';
import type { MdConn } from '@/types/metadata';
import { ArrowDown, ArrowUp, InfoFilled, Refresh } from '@element-plus/icons-vue';
import { onMounted, ref } from 'vue';

const props = defineProps<{
    model: any
    collapsed: boolean
}>()

const emit = defineEmits(['toggle-collapse', 'conn-change', 'refresh-code'])

const dataSources = ref<MdConn[]>([])

onMounted(async () => {
    try {
        const res: any = await getConns()
        dataSources.value = Array.isArray(res) ? res : (res.data || [])

        // 如果没有选择数据源且列表不为空，默认选择第一个 (对应“自动获取数据源”)
        if (!props.model.connID && dataSources.value.length > 0) {
            props.model.connID = dataSources.value[0].id
            handleConnChange(props.model.connID)
        }
    } catch (e) {
        console.error('Failed to fetch data sources', e)
    }
})

const handleConnChange = (val: any) => {
    const conn = dataSources.value.find(c => c.id === val)
    emit('conn-change', conn)
}
</script>
<style scoped>
.base-info-panel {
    display: flex;
    flex-direction: column;
    background-color: #fff;
    transition: all 0.3s ease;
    width: 100%;
    height: 100%;
}

.panel-header {
    padding: 10px 16px;
    border-bottom: 1px solid #f2f6fc;
    font-weight: 600;
    color: #303133;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: move;
    /* Indicate draggable */
    user-select: none;
    flex-shrink: 0;
}

.header-content {
    display: flex;
    align-items: center;
    gap: 8px;
}

.collapse-icon {
    cursor: pointer;
    font-size: 14px;
    color: #909399;
}

.panel-body {
    padding: 16px;
    flex: 1;
    overflow-y: auto;
    /* max-height removed to allow filling container */
}

.is-collapsed .panel-header {
    border-bottom: none;
}

.is-collapsed {
    height: auto !important;
}

:deep(.el-input-group__append) {
    padding: 0 10px;
}
</style>

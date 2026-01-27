<template>
    <el-dialog v-bind="$attrs" v-model="visible" :title="title" :width="width" :fullscreen="isFullscreen" :destroy-on-close="destroyOnClose" :close-on-click-modal="closeOnClickModal" :show-close="false" class="base-dialog" :class="{ 'is-fullscreen': isFullscreen }" @close="handleClose">
        <template #header>
            <div class="dialog-header">
                <div class="header-left">
                    <slot name="header">
                        <span class="dialog-title">{{ title }}</span>
                    </slot>
                </div>
                <div class="dialog-header-actions">
                    <el-icon class="action-icon" @click="toggleFullscreen" :title="isFullscreen ? '还原' : '全屏'">
                        <FullScreen v-if="!isFullscreen" />
                        <CopyDocument v-else />
                    </el-icon>
                    <el-icon class="action-icon close-icon" @click="handleClose" title="关闭">
                        <Close />
                    </el-icon>
                </div>
            </div>
        </template>
        <div class="dialog-body" :style="bodyStyle">
            <slot></slot>
        </div>
        <template #footer v-if="showFooter">
            <slot name="footer">
                <div class="dialog-footer">
                    <el-button @click="handleCancel">{{ cancelText }}</el-button>
                    <el-button type="primary" :loading="loading" @click="handleConfirm">
                        {{ confirmText }}
                    </el-button>
                </div>
            </slot>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { Close, CopyDocument, FullScreen } from '@element-plus/icons-vue'
import { computed, ref, watch } from 'vue'

interface Props {
    modelValue: boolean
    title?: string
    width?: string | number
    loading?: boolean
    confirmText?: string
    cancelText?: string
    showFooter?: boolean
    destroyOnClose?: boolean
    closeOnClickModal?: boolean
    maxHeight?: string | number
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    title: '提示',
    width: '50%',
    loading: false,
    confirmText: '确定',
    cancelText: '取消',
    showFooter: true,
    destroyOnClose: true,
    closeOnClickModal: false,
    maxHeight: '60vh'
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const isFullscreen = ref(false)

const bodyStyle = computed(() => {
    if (isFullscreen.value) {
        return {
            height: 'calc(100vh - 110px)',
            overflow: 'auto'
        }
    }
    return {
        maxHeight: typeof props.maxHeight === 'number' ? `${props.maxHeight}px` : props.maxHeight,
        overflow: 'auto'
    }
})

const toggleFullscreen = () => {
    isFullscreen.value = !isFullscreen.value
}

const handleConfirm = () => {
    emit('confirm')
}

const handleCancel = () => {
    visible.value = false
    emit('cancel')
}

const handleClose = () => {
    emit('close')
}

// Reset fullscreen when closed
watch(() => props.modelValue, (val) => {
    if (!val) {
        setTimeout(() => {
            isFullscreen.value = false
        }, 300)
    }
})
</script>
<style scoped>
.base-dialog :deep(.el-dialog) {
    display: flex;
    flex-direction: column;
    margin: 0 !important;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    border-radius: 8px;
    overflow: hidden;
}

.base-dialog :deep(.el-dialog.is-fullscreen) {
    width: 100% !important;
    height: 100% !important;
    top: 0 !important;
    left: 0 !important;
    transform: none !important;
    border-radius: 0;
}

.base-dialog :deep(.el-dialog__header) {
    padding: 0;
    margin-right: 0;
    border-bottom: 1px solid #f0f0f0;
}

.base-dialog :deep(.el-dialog__body) {
    padding: 20px;
    flex: 1;
    overflow: auto;
}

.base-dialog :deep(.el-dialog__footer) {
    padding: 10px 20px;
    border-top: 1px solid #f0f0f0;
}

.dialog-header {
    display: flex;
    align-items: center;
    height: 50px;
    padding: 0 20px;
    width: 100%;
}

.header-left {
    flex: 1;
    overflow: hidden;
}

.dialog-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.dialog-header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
}

.action-icon {
    font-size: 18px;
    cursor: pointer;
    color: #909399;
    transition: color 0.3s;
}

.action-icon:hover {
    color: #409eff;
}

.close-icon:hover {
    color: #f56c6c;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}
</style>

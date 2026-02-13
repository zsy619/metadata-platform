<template>
    <div class="empty-state" :class="{ 'is-center': center, [`empty-size-${size}`]: true }">
        <div class="empty-image">
            <img v-if="imageUrl" :src="imageUrl" class="empty-img" :style="{ width: imageSize + 'px' }" />
            <component v-else :is="emptyImage" class="empty-icon" :style="{ fontSize: imageSize + 'px' }" />
        </div>
        <div class="empty-description">
            <p v-if="description" class="empty-desc-text">{{ description }}</p>
            <p v-else class="empty-desc-text">{{ defaultDescription }}</p>
        </div>
        <div v-if="$slots.default || showAction" class="empty-action">
            <slot>
                <el-button v-if="showAction" :type="buttonType" :size="buttonSize" @click="handleClick">
                    {{ buttonText }}
                </el-button>
            </slot>
        </div>
    </div>
</template>
<script setup lang="ts">
import { computed } from 'vue'
import { WarningFilled } from '@element-plus/icons-vue'

interface Props {
    image?: string
    imageSize?: number
    description?: string
    center?: boolean
    showAction?: boolean
    buttonText?: string
    buttonType?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
    buttonSize?: 'large' | 'default' | 'small'
    size?: 'large' | 'default' | 'small'
}

const props = withDefaults(defineProps<Props>(), {
    image: '',
    imageSize: 120,
    description: '',
    center: true,
    showAction: false,
    buttonText: '刷新',
    buttonType: 'primary',
    buttonSize: 'default',
    size: 'default'
})

const emit = defineEmits(['click'])

const imageUrl = computed(() => props.image)

const emptyImage = WarningFilled

const defaultDescription = '暂无数据'

const handleClick = () => {
    emit('click')
}
</script>
<style scoped>
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
    text-align: center;
    box-sizing: border-box;
}

.empty-state.is-center {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
}

.empty-image {
    margin-bottom: 16px;
}

.empty-img {
    display: block;
    width: 120px;
    height: auto;
}

.empty-icon {
    color: #c0c4cc;
}

.empty-description {
    margin-bottom: 24px;
}

.empty-desc-text {
    color: #909399;
    font-size: 14px;
    margin: 0;
    line-height: 1.5;
}

.empty-size-large .empty-desc-text {
    font-size: 16px;
}

.empty-size-small .empty-desc-text {
    font-size: 12px;
}

.empty-action {
    display: flex;
    justify-content: center;
    gap: 12px;
}
</style>

<template>
    <div class="page-container">
        <div class="page-header" v-if="$slots.header || title">
            <div class="header-content">
                <div class="title" v-if="title">{{ title }}</div>
                <div class="actions">
                    <slot name="header"></slot>
                </div>
            </div>
            <div class="description" v-if="description">{{ description }}</div>
        </div>
        <div class="page-content" :class="{ 'has-padding': padding }">
            <slot></slot>
        </div>
    </div>
</template>
<script setup lang="ts">
defineProps({
    title: {
        type: String,
        default: ''
    },
    description: {
        type: String,
        default: ''
    },
    padding: {
        type: Boolean,
        default: true
    }
})
</script>
<style scoped>
.page-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 100%;
    background-color: #f0f2f5;
}

.page-header {
    background: #fff;
    padding: 16px 24px;
    border-bottom: 1px solid #e8e8e8;
}

.header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.title {
    font-size: 20px;
    font-weight: 500;
    color: #1f2f3d;
}

.description {
    margin-top: 12px;
    color: #909399;
    font-size: 14px;
}

.page-content {
    flex: 1;
    overflow: auto;
}

.page-content.has-padding {
    padding: 10px;
}
</style>

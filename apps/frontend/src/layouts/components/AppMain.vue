<template>
    <section class="app-main">
        <router-view v-slot="{ Component, route }">
            <transition name="fade-transform" mode="out-in">
                <keep-alive :include="cachedViews">
                    <component :is="Component" :key="route.path" />
                </keep-alive>
            </transition>
        </router-view>
    </section>
</template>
<script setup lang="ts">
import { computed } from 'vue';

// Future enhancement: get from tagsView store
const cachedViews = computed(() => [])
</script>
<style scoped>
.app-main {
    /* 50 = navbar  */
    min-height: calc(100vh - 50px);
    width: 100%;
    position: relative;
    overflow: hidden;
    padding: 20px;
    background-color: var(--el-bg-color-page);
}

.fixed-header+.app-main {
    padding-top: 50px;
}

/* fade-transform */
.fade-transform-leave-active,
.fade-transform-enter-active {
    transition: all .5s;
}

.fade-transform-enter-from {
    opacity: 0;
    transform: translateX(-30px);
}

.fade-transform-leave-to {
    opacity: 0;
    transform: translateX(30px);
}
</style>

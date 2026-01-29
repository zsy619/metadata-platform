<template>
    <el-config-provider :locale="locale">
        <div class="app-container">
            <router-view />
        </div>
    </el-config-provider>
</template>
<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import { ElConfigProvider } from 'element-plus'
import en from 'element-plus/es/locale/lang/en'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const appStore = useAppStore()
const i18n = useI18n()

// 动态语言配置
const locale = computed(() => {
    return appStore.language === 'en-US' ? en : zhCn
})

watch(() => appStore.language, (val: string) => {
    i18n.locale.value = val
}, { immediate: true })

onMounted(() => {
    // 初始化主题
    appStore.setTheme(appStore.theme)
})
</script>
<style scoped>
.app-container {
    width: 100%;
    height: 100vh;
    overflow: hidden;
}
</style>
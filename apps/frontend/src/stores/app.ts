import storage from '@/utils/storage'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
    // 状态
    const sidebar = ref({
        opened: storage.get('sidebarStatus') ?? true,
        withoutAnimation: false,
        width: storage.get('sidebarWidth') || 310,
        isResizing: false
    })

    const device = ref('desktop')
    const theme = ref(storage.get('theme') || 'light')
    const language = ref(storage.get('language') || 'zh-CN')

    // Actions
    const toggleSidebar = () => {
        sidebar.value.opened = !sidebar.value.opened
        sidebar.value.withoutAnimation = false
        if (sidebar.value.opened) {
            storage.set('sidebarStatus', 1)
        } else {
            storage.set('sidebarStatus', 0)
        }
    }

    const closeSidebar = (withoutAnimation: boolean) => {
        storage.set('sidebarStatus', 0)
        sidebar.value.opened = false
        sidebar.value.withoutAnimation = withoutAnimation
    }

    const setSidebarWidth = (width: number) => {
        sidebar.value.width = width
        storage.set('sidebarWidth', width)
    }

    const setResizing = (isResizing: boolean) => {
        sidebar.value.isResizing = isResizing
    }

    const toggleDevice = (d: string) => {
        device.value = d
    }

    const setTheme = (t: string) => {
        theme.value = t
        storage.set('theme', t)
        if (t === 'dark') {
            document.documentElement.classList.add('dark')
            document.documentElement.style.colorScheme = 'dark'
        } else {
            document.documentElement.classList.remove('dark')
            document.documentElement.style.colorScheme = 'light'
        }
    }

    const setLanguage = (lang: string) => {
        language.value = lang
        storage.set('language', lang)
    }

    return {
        sidebar,
        device,
        theme,
        language,
        toggleSidebar,
        closeSidebar,
        toggleDevice,
        setTheme,
        setLanguage,
        setSidebarWidth,
        setResizing
    }
})

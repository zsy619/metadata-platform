import storage from '@/utils/storage'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
    // 状态
    const sidebar = ref({
        opened: storage.get('sidebarStatus') ?? true,
        withoutAnimation: false
    })

    const device = ref('desktop')
    const theme = ref(storage.get('theme') || 'light')

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

    const toggleDevice = (d: string) => {
        device.value = d
    }

    const setTheme = (t: string) => {
        theme.value = t
        storage.set('theme', t)
        // 这里可以添加切换CSS变量或类的逻辑
        document.documentElement.className = t
    }

    return {
        sidebar,
        device,
        theme,
        toggleSidebar,
        closeSidebar,
        toggleDevice,
        setTheme
    }
})

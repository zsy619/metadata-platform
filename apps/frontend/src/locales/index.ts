import storage from '@/utils/storage'
import { createI18n } from 'vue-i18n'
import enUS from './en-US'
import zhCN from './zh-CN'

const i18n = createI18n({
    legacy: false, // Use Composition API
    locale: storage.get('language') || 'zh-CN',
    fallbackLocale: 'zh-CN',
    messages: {
        'zh-CN': zhCN,
        'en-US': enUS
    }
})

export default i18n

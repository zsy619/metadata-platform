import i18n from '@/locales'
import { setupErrorHandle } from '@/utils/error-log'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import './assets/styles/index.css'
import router from './router'

Object.values(fas).forEach((icon: any) => {
  library.add(icon)
})
Object.values(far).forEach((icon: any) => {
  library.add(icon)
})
Object.values(fab).forEach((icon: any) => {
  library.add(icon)
})

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.component('FontAwesomeIcon', FontAwesomeIcon)

// 安装插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(i18n)

// 注册全局错误处理
setupErrorHandle(app)

// 挂载应用
app.mount('#app')
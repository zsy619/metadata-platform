import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import './assets/styles/index.css'
import router from './router'

const app = createApp(App)

// 安装插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus)

// 挂载应用
app.mount('#app')
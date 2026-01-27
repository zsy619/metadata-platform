import { createRouter, createWebHistory } from 'vue-router'
import { setupGuards } from './guards'
import routes from './routes'

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior(_to, _from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { top: 0 }
        }
    }
})

setupGuards(router)

export default router
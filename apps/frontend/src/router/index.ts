import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Layout',
        component: () => import('@/components/Layout.vue'),
        redirect: '/data-sources',
        children: [
            // 数据源管理
            {
                path: '/data-sources',
                name: 'DataSourceList',
                component: () => import('@/views/data-source/List.vue'),
                meta: { title: '数据源管理' }
            },
            {
                path: '/data-sources/create',
                name: 'DataSourceCreate',
                component: () => import('@/views/data-source/Create.vue'),
                meta: { title: '创建数据源' }
            },
            {
                path: '/data-sources/:id/edit',
                name: 'DataSourceEdit',
                component: () => import('@/views/data-source/Edit.vue'),
                meta: { title: '编辑数据源' }
            },

            // 元数据管理
            {
                path: '/metadata/tables',
                name: 'MetadataTableList',
                component: () => import('@/views/metadata/TableList.vue'),
                meta: { title: '表与视图' }
            },
            {
                path: '/metadata/fields',
                name: 'MetadataFieldList',
                component: () => import('@/views/metadata/FieldList.vue'),
                meta: { title: '字段列表' }
            },

            // 模型管理
            {
                path: '/models',
                name: 'ModelList',
                component: () => import('@/views/model/List.vue'),
                meta: { title: '模型管理' }
            },
            {
                path: '/models/create',
                name: 'ModelCreate',
                component: () => import('@/views/model/Create.vue'),
                meta: { title: '创建模型' }
            },
            {
                path: '/models/:id/edit',
                name: 'ModelEdit',
                component: () => import('@/views/model/Edit.vue'),
                meta: { title: '编辑模型' }
            },
            {
                path: '/models/:id/preview',
                name: 'ModelPreview',
                component: () => import('@/views/model/Preview.vue'),
                meta: { title: '数据预览' }
            },

            // 接口管理
            {
                path: '/apis',
                name: 'APIList',
                component: () => import('@/views/api/List.vue'),
                meta: { title: '接口管理' }
            },
            {
                path: '/apis/create',
                name: 'APICreate',
                component: () => import('@/views/api/Create.vue'),
                meta: { title: '创建接口' }
            },
            {
                path: '/apis/:id/edit',
                name: 'APIEdit',
                component: () => import('@/views/api/Edit.vue'),
                meta: { title: '编辑接口' }
            },

            // 系统设置
            {
                path: '/system/settings',
                name: 'SystemSettings',
                component: () => import('@/views/system/Settings.vue'),
                meta: { title: '系统设置' }
            },
            {
                path: '/system/audit',
                name: 'AuditLog',
                component: () => import('@/views/system/AuditLog.vue'),
                meta: { title: '审计日志' }
            },

            // 个人设置
            {
                path: '/profile',
                name: 'UserProfile',
                component: () => import('@/views/profile/Index.vue'),
                meta: { title: '个人设置' }
            }
        ]
    },

    // 登录页面
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { title: '登录' }
    },

    // 404页面
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFound.vue'),
        meta: { title: '页面不存在' }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// 路由前置守卫
router.beforeEach((to, _from, next) => {
    // 设置页面标题
    document.title = (to.meta.title as string) || '元数据管理平台'

    // 登录检查
    const token = localStorage.getItem('token')
    if (to.path !== '/login' && !token) {
        next('/login')
    } else {
        next()
    }
})

export default router
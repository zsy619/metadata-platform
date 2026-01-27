import Layout from '@/layouts/DefaultLayout.vue'
import { type RouteRecordRaw } from 'vue-router'

// Wrapper component for nested routes
const RouterView = { template: '<router-view />' }

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: Layout,
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/home/Dashboard.vue'),
                meta: { title: '首页', icon: 'Odometer', affix: true }
            },
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/profile/Index.vue'),
                meta: { title: '个人设置', hidden: true }
            }
        ]
    },
    {
        path: '/data-sources',
        component: Layout,
        meta: { title: '数据源管理', icon: 'Connection' },
        redirect: '/data-sources/list',
        children: [
            {
                path: 'list',
                name: 'DataSourceList',
                component: () => import('@/views/data-source/List.vue'),
                meta: { title: '数据源列表', icon: 'List' }
            },
            {
                path: 'create',
                name: 'DataSourceCreate',
                component: () => import('@/views/data-source/Create.vue'),
                meta: { title: '添加数据源', icon: 'Plus' }
            },
            {
                path: ':id/edit',
                name: 'DataSourceEdit',
                component: () => import('@/views/data-source/Edit.vue'),
                meta: { title: '编辑数据源', hidden: true }
            },
            {
                path: 'metadata',
                name: 'Metadata',
                meta: { title: '元数据管理', icon: 'FolderOpened' },
                component: RouterView, // Nested group
                children: [
                    {
                        path: 'tables',
                        name: 'MetadataTableList',
                        component: () => import('@/views/metadata/TableList.vue'),
                        meta: { title: '表与视图', icon: 'Grid' }
                    },
                    {
                        path: 'fields',
                        name: 'MetadataFieldList',
                        component: () => import('@/views/metadata/FieldList.vue'),
                        meta: { title: '字段列表', icon: 'Tickets' }
                    }
                ]
            }
        ]
    },
    {
        path: '/models',
        component: Layout,
        meta: { title: '模型管理', icon: 'Document' },
        redirect: '/models/list',
        children: [
            {
                path: 'list',
                name: 'ModelList',
                component: () => import('@/views/model/List.vue'),
                meta: { title: '模型列表', icon: 'List' }
            },
            {
                path: 'create',
                name: 'ModelCreate',
                component: () => import('@/views/model/Create.vue'),
                meta: { title: '创建模型', hidden: true }
            },
            {
                path: ':id/edit',
                name: 'ModelEdit',
                component: () => import('@/views/model/Edit.vue'),
                meta: { title: '编辑模型', hidden: true }
            },
            {
                path: ':id/preview',
                name: 'ModelPreview',
                component: () => import('@/views/model/Preview.vue'),
                meta: { title: '数据预览', hidden: true }
            }
        ]
    },
    {
        path: '/apis',
        component: Layout,
        meta: { title: '接口管理', icon: 'Share' },
        redirect: '/apis/list',
        children: [
            {
                path: 'list',
                name: 'APIList',
                component: () => import('@/views/api/List.vue'),
                meta: { title: '接口列表', icon: 'List' }
            },
            {
                path: 'create',
                name: 'APICreate',
                component: () => import('@/views/api/Create.vue'),
                meta: { title: '创建接口', hidden: true }
            },
            {
                path: ':id/edit',
                name: 'APIEdit',
                component: () => import('@/views/api/Edit.vue'),
                meta: { title: '编辑接口', hidden: true }
            }
        ]
    },
    {
        path: '/system',
        component: Layout,
        meta: { title: '系统设置', icon: 'Setting' },
        redirect: '/system/settings',
        children: [
            {
                path: 'settings',
                name: 'SystemSettings',
                component: () => import('@/views/system/Settings.vue'),
                meta: { title: '系统配置', icon: 'Tools' }
            },
            {
                path: 'audit',
                name: 'Audit',
                meta: { title: '审计日志', icon: 'Document' },
                component: RouterView,
                redirect: '/system/audit/login',
                children: [
                    {
                        path: 'login',
                        name: 'LoginLog',
                        component: () => import('@/views/system/audit/LoginLog.vue'),
                        meta: { title: '登录日志', icon: 'User' }
                    },
                    {
                        path: 'operation',
                        name: 'OperationLog',
                        component: () => import('@/views/system/audit/OperationLog.vue'),
                        meta: { title: '操作日志', icon: 'Edit' }
                    },
                    {
                        path: 'data',
                        name: 'DataLog',
                        component: () => import('@/views/system/audit/DataLog.vue'),
                        meta: { title: '数据日志', icon: 'DataLine' }
                    }
                ]
            }
        ]
    },
    {
        path: '/login',
        component: () => import('@/layouts/AuthLayout.vue'),
        children: [
            {
                path: '',
                name: 'Login',
                component: () => import('@/views/auth/Login.vue'),
                meta: { title: '登录' }
            }
        ]
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/error/404.vue'),
        meta: { title: '页面不存在', hidden: true }
    }
]

export default routes

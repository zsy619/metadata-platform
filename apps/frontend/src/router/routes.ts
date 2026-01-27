import Layout from '@/layouts/DefaultLayout.vue'
import { h } from 'vue'
import { RouterView as VueRouterView, type RouteRecordRaw } from 'vue-router'

// Wrapper component for nested routes to avoid "router-view inside transition" warning
const RouterView = {
    name: 'RouterViewWrapper',
    render: () => h(VueRouterView)
}


const routes: RouteRecordRaw[] = [
    {
        path: '/',
        component: Layout,
        redirect: '/home/dashboard',
        children: [
            {
                path: 'home/dashboard',
                name: 'Dashboard',
                component: () => import('@/views/home/Dashboard.vue'),
                meta: { title: '首页', icon: 'Odometer', affix: true }
            }
        ]
    },
    {
        path: '/metadata',
        component: Layout,
        meta: { title: '元数据管理', icon: 'FolderOpened' },
        redirect: '/metadata/datasource/list',
        children: [
            {
                path: 'datasource',
                component: RouterView,
                meta: { title: '数据源管理', icon: 'Connection' },
                children: [
                    {
                        path: 'list',
                        name: 'DataSourceList',
                        component: () => import('@/views/metadata/datasource/List.vue'),
                        meta: { title: '数据源列表', icon: 'List' }
                    },
                    {
                        path: 'create',
                        name: 'DataSourceCreate',
                        component: () => import('@/views/metadata/datasource/Create.vue'),
                        meta: { title: '添加数据源', icon: 'Plus' }
                    },
                    {
                        path: ':id/edit',
                        name: 'DataSourceEdit',
                        component: () => import('@/views/metadata/datasource/Edit.vue'),
                        meta: { title: '编辑数据源', hidden: true }
                    }
                ]
            },
            {
                path: 'table',
                component: RouterView,
                meta: { title: '表与视图', icon: 'Grid' },
                children: [
                    {
                        path: 'list',
                        name: 'MetadataTableList',
                        component: () => import('@/views/metadata/table/List.vue'),
                        meta: { title: '表列表', icon: 'List' }
                    }
                ]
            },
            {
                path: 'field',
                component: RouterView,
                meta: { title: '字段列表', icon: 'Tickets' },
                children: [
                    {
                        path: 'list',
                        name: 'MetadataFieldList',
                        component: () => import('@/views/metadata/field/List.vue'),
                        meta: { title: '字段列表', icon: 'List' }
                    }
                ]
            }
        ]
    },
    {
        path: '/user',
        component: Layout,
        redirect: '/user/profile',
        meta: { hidden: true },
        children: [
            {
                path: 'profile',
                name: 'Profile',
                component: () => import('@/views/user/Profile.vue'),
                meta: { title: '个人设置' }
            }
        ]
    },
    {
        path: '/model',
        component: Layout,
        meta: { title: '模型管理', icon: 'Document' },
        redirect: '/model/list',
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
        path: '/api',
        component: Layout,
        meta: { title: '接口管理', icon: 'Share' },
        redirect: '/api/list',
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
                component: RouterView,
                redirect: '/system/audit/login',
                meta: { title: '审计日志', icon: 'Document' },
                children: [
                    {
                        path: 'login',
                        name: 'AuditLoginLog',
                        component: () => import('@/views/system/audit/LoginLog.vue'),
                        meta: { title: '登录日志', icon: 'UserFilled' }
                    },
                    {
                        path: 'operation',
                        name: 'AuditOperationLog',
                        component: () => import('@/views/system/audit/OperationLog.vue'),
                        meta: { title: '操作日志', icon: 'List' }
                    },
                    {
                        path: 'data',
                        name: 'AuditDataChangeLog',
                        component: () => import('@/views/system/audit/DataChangeLog.vue'),
                        meta: { title: '数据变更', icon: 'Edit' }
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
                component: () => import('@/views/login/Index.vue'),
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

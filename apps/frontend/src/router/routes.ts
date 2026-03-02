import Layout from '@/layouts/DefaultLayout.vue'
import { h } from 'vue'
import { RouterView as VueRouterView, type RouteRecordRaw } from 'vue-router'

const RouterView = {
    name: 'RouterViewWrapper',
    render: () => h(VueRouterView)
}

const routes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/Index.vue'),
        meta: { title: '登录', hidden: true }
    },
    {
        path: '/',
        component: Layout,
        redirect: '/home/dashboard',
        children: [
            {
                path: 'home/dashboard',
                name: 'Dashboard',
                component: () => import('@/views/home/Dashboard.vue'),
                meta: { title: '首页', icon: 'fa-gauge-high', affix: true }
            }
        ]
    },
    {
        path: '/metadata',
        component: Layout,
        meta: { title: '元数据管理', icon: 'fa-folder-open' },
        redirect: '/metadata/datasource/list',
        children: [
            {
                path: 'datasource',
                component: RouterView,
                meta: { title: '数据源管理', icon: 'fa-network-wired' },
                children: [
                    {
                        path: 'list',
                        name: 'DataSourceList',
                        component: () => import('@/views/metadata/datasource/List.vue'),
                        meta: { title: '数据源列表', icon: 'fa-list' }
                    },
                    {
                        path: 'create',
                        name: 'DataSourceCreate',
                        component: () => import('@/views/metadata/datasource/Create.vue'),
                        meta: { title: '添加数据源', icon: 'fa-plus' }
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
                path: 'maintenance',
                component: RouterView,
                meta: { title: '元数据维护', icon: 'fa-wrench' },
                children: [
                    {
                        path: 'table',
                        name: 'MetadataTableList',
                        component: () => import('@/views/metadata/table/List.vue'),
                        meta: { title: '表列表', icon: 'fa-table' }
                    },
                    {
                        path: 'view',
                        name: 'MetadataViewList',
                        component: () => import('@/views/metadata/table/List.vue'),
                        meta: { title: '视图列表', icon: 'fa-table-columns' }
                    },
                    {
                        path: 'procedure',
                        name: 'MetadataProcedureList',
                        component: () => import('@/views/metadata/procedure/List.vue'),
                        meta: { title: '存储过程', icon: 'fa-database' }
                    },
                    {
                        path: 'function',
                        name: 'MetadataFunctionList',
                        component: () => import('@/views/metadata/procedure/List.vue'),
                        meta: { title: '函数', icon: 'fa-calculator' }
                    },
                    {
                        path: 'field',
                        component: RouterView,
                        meta: { title: '字段列表', icon: 'fa-ticket' },
                        children: [
                            {
                                path: 'list',
                                name: 'MetadataFieldList',
                                component: () => import('@/views/metadata/field/List.vue'),
                                meta: { title: '字段列表', icon: 'fa-list' }
                            }
                        ]
                    }
                ]
            },
            {
                path: 'model',
                component: RouterView,
                meta: { title: '模型维护', icon: 'fa-file' },
                redirect: '/metadata/model/list',
                children: [
                    {
                        path: 'list',
                        name: 'ModelList',
                        component: () => import('@/views/model/List.vue'),
                        meta: { title: '模型列表', icon: 'fa-list' }
                    },
                    {
                        path: 'sql-test',
                        name: 'ModelSqlTest',
                        component: () => import('@/views/model/SqlTest.vue'),
                        meta: { title: 'SQL模型测试', icon: 'fa-desktop' }
                    },
                    {
                        path: 'table-view-test',
                        name: 'ModelTableViewTest',
                        component: () => import('@/views/model/TableViewTest.vue'),
                        meta: { title: '表视图模型测试', icon: 'fa-table' }
                    },
                    {
                        path: 'create',
                        name: 'ModelCreate',
                        component: () => import('@/views/model/Create.vue'),
                        meta: { title: '创建模型', hidden: true }
                    },
                    {
                        path: 'visual-create',
                        name: 'ModelVisualCreate',
                        component: () => import('@/views/model/VisualCreate.vue'),
                        meta: { title: '可视化创建模型', hidden: true }
                    },
                    {
                        path: 'visual-edit/:id',
                        name: 'ModelVisualEdit',
                        component: () => import('@/views/model/VisualCreate.vue'),
                        meta: { title: '可视化编辑模型', hidden: true }
                    },
                    {
                        path: 'create-sql',
                        name: 'ModelCreateSql',
                        component: () => import('@/views/model/CreateSql.vue'),
                        meta: { title: '创建SQL模型', hidden: true }
                    },
                    {
                        path: 'edit-sql/:id',
                        name: 'ModelEditSql',
                        component: () => import('@/views/model/CreateSql.vue'),
                        meta: { title: '编辑SQL模型', hidden: true }
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
            }
        ]
    },
    {
        path: '/monitor',
        component: Layout,
        meta: { title: '系统监控', icon: 'fa-desktop' },
        redirect: '/monitor/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'MonitorDashboard',
                component: () => import('@/views/monitor/Dashboard.vue'),
                meta: { title: '监控仪表盘', icon: 'fa-gauge-high' }
            },
            {
                path: 'performance',
                name: 'MonitorPerformance',
                component: () => import('@/views/monitor/Performance.vue'),
                meta: { title: '性能分析', icon: 'fa-chart-line' }
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
        path: '/api',
        component: Layout,
        meta: { title: '接口管理', icon: 'fa-share-nodes' },
        redirect: '/api/list',
        children: [
            {
                path: 'list',
                name: 'APIList',
                component: () => import('@/views/api/List.vue'),
                meta: { title: '接口列表', icon: 'fa-list' }
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
        path: '/sso',
        component: Layout,
        meta: { title: 'SSO管理', icon: 'fa-key' },
        redirect: '/sso/config/protocol',
        children: [
            {
                path: 'config',
                component: RouterView,
                meta: { title: 'SSO配置', icon: 'fa-gear' },
                redirect: '/sso/config/protocol',
                children: [
                    {
                        path: 'protocol',
                        name: 'SSOProtocol',
                        component: () => import('@/views/sso/config/ProtocolConfig.vue'),
                        meta: { title: '协议配置', icon: 'fa-plug' }
                    },
                    {
                        path: 'client',
                        name: 'SSOClient',
                        component: () => import('@/views/sso/config/ClientConfig.vue'),
                        meta: { title: '客户端配置', icon: 'fa-desktop' }
                    },
                    {
                        path: 'key',
                        name: 'SSOKey',
                        component: () => import('@/views/sso/config/KeyManager.vue'),
                        meta: { title: '密钥管理', icon: 'fa-lock' }
                    },
                    {
                        path: 'mapping',
                        name: 'SSOMapping',
                        component: () => import('@/views/sso/config/FieldMapping.vue'),
                        meta: { title: '字段映射', icon: 'fa-exchange-alt' }
                    },
                    {
                        path: 'session',
                        name: 'SSOSession',
                        component: () => import('@/views/sso/config/SessionManager.vue'),
                        meta: { title: '会话管理', icon: 'fa-clock' }
                    }
                ]
            },
            {
                path: 'tenant',
                name: 'SSOTenant',
                component: () => import('@/views/sso/tenant/Tenant.vue'),
                meta: { title: '租户管理', icon: 'fa-building' }
            },
            {
                path: 'app',
                name: 'SSOApp',
                component: () => import('@/views/sso/app/App.vue'),
                meta: { title: '应用列表', icon: 'fa-th-large' }
            },
            {
                path: 'org',
                name: 'SSOOrg',
                component: () => import('@/views/sso/org/Org.vue'),
                meta: { title: '组织管理', icon: 'fa-sitemap' }
            },
            {
                path: 'orgKind',
                name: 'SSOOrgKind',
                component: () => import('@/views/sso/org/OrgKind.vue'),
                meta: { title: '组织类型', icon: 'fa-layer-group' }
            },
            {
                path: 'menu',
                name: 'SSOMenu',
                component: () => import('@/views/sso/menu/Menu.vue'),
                meta: { title: '菜单管理', icon: 'fa-bars' }
            },
            {
                path: 'pos',
                name: 'SSOPos',
                component: () => import('@/views/sso/pos/Pos.vue'),
                meta: { title: '职位管理', icon: 'fa-briefcase' }
            },
            {
                path: 'role',
                name: 'SSORole',
                component: () => import('@/views/sso/role/Role.vue'),
                meta: { title: '角色管理', icon: 'fa-user-shield' }
            },
            {
                path: 'roleGroup',
                name: 'SSORoleGroup',
                component: () => import('@/views/sso/role/RoleGroup.vue'),
                meta: { title: '角色组', icon: 'fa-folder' }
            },
            {
                path: 'userGroup',
                name: 'SSOUserGroup',
                component: () => import('@/views/sso/user/UserGroup.vue'),
                meta: { title: '用户组', icon: 'fa-users' }
            },
            {
                path: 'user',
                name: 'SSOUser',
                component: () => import('@/views/sso/user/User.vue'),
                meta: { title: '用户管理', icon: 'fa-user-astronaut' }
            }
        ]
    },
    {
        path: '/system',
        component: Layout,
        meta: { title: '系统设置', icon: 'fa-gear' },
        redirect: '/system/settings',
        children: [
            {
                path: 'settings',
                name: 'SystemSettings',
                component: () => import('@/views/system/Settings.vue'),
                meta: { title: '系统配置', icon: 'fa-wrench' }
            },
            {
                path: 'audit',
                component: RouterView,
                redirect: '/system/audit/login',
                meta: { title: '审计日志', icon: 'fa-clock' },
                children: [
                    {
                        path: 'login',
                        name: 'AuditLoginLog',
                        component: () => import('@/views/system/audit/LoginLog.vue'),
                        meta: { title: '登录日志', icon: 'fa-user' }
                    },
                    {
                        path: 'operation',
                        name: 'AuditOperationLog',
                        component: () => import('@/views/system/audit/OperationLog.vue'),
                        meta: { title: '操作日志', icon: 'fa-list' }
                    },
                    {
                        path: 'data',
                        name: 'AuditDataChangeLog',
                        component: () => import('@/views/system/audit/DataChangeLog.vue'),
                        meta: { title: '数据变更', icon: 'fa-pen-to-square' }
                    },
                    {
                        path: 'access',
                        name: 'AuditAccessLog',
                        component: () => import('@/views/system/audit/AccessLog.vue'),
                        meta: { title: '访问日志', icon: 'fa-desktop' }
                    }
                ]
            }
        ]
    },
    {
        path: '/docs',
        component: Layout,
        meta: { title: '系统文档', icon: 'fa-book' },
        redirect: '/docs/list',
        children: [
            {
                path: 'list',
                name: 'DocumentList',
                component: () => import('@/views/docs/DocList.vue'),
                meta: { title: '文档列表', icon: 'fa-list' }
            },
            {
                path: ':id',
                name: 'DocumentDetail',
                component: () => import('@/views/docs/DocDetail.vue'),
                meta: { title: '文档详情', hidden: true }
            }
        ]
    },
    {
        path: '/documents',
        component: Layout,
        meta: { title: '文档管理', icon: 'fa-file-lines' },
        redirect: '/documents/list',
        children: [
            {
                path: 'list',
                name: 'DocumentManageList',
                component: () => import('@/views/document/List.vue'),
                meta: { title: '文档目录', icon: 'fa-folder-tree' }
            },
            {
                path: 'create',
                name: 'DocumentCreate',
                component: () => import('@/views/document/Create.vue'),
                meta: { title: '新建文档', icon: 'fa-plus', hidden: true }
            },
            {
                path: ':id/edit',
                name: 'DocumentEdit',
                component: () => import('@/views/document/Create.vue'),
                meta: { title: '编辑文档', hidden: true }
            }
        ]
    },
    {
        path: '/403',
        name: 'Forbidden',
        component: () => import('@/views/error/403.vue'),
        meta: { title: '无权限', hidden: true }
    },
    {
        path: '/500',
        name: 'ServerError',
        component: () => import('@/views/error/500.vue'),
        meta: { title: '服务器错误', hidden: true }
    },
    {
        path: '/503',
        name: 'ServiceUnavailable',
        component: () => import('@/views/error/503.vue'),
        meta: { title: '服务维护中', hidden: true }
    },
    {
        path: '/network-error',
        name: 'NetworkError',
        component: () => import('@/views/error/NetworkError.vue'),
        meta: { title: '网络错误', hidden: true }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/error/404.vue'),
        meta: { title: '页面不存在', hidden: true }
    }
]

export default routes

/*
 菜单初始化数据
 
 基于前端路由配置生成
 对应模型: SsoMenu (sso_menu表)
 
 菜单类型说明:
 - M: 目录
 - C: 菜单
 - F: 按钮
 - Z: 资源
 
 数据权限范围(data_range):
 - 1: 全部数据权限
 - 2: 自定义数据权限
 - 3: 本部门数据权限
 - 4: 本部门及以下数据权限
 
 Date: 2026-02-27
*/

SET NAMES utf8mb4;

-- 清空现有菜单数据（可选，谨慎使用）
-- TRUNCATE TABLE `sso_menu`;

-- ============================================
-- 一级菜单
-- ============================================

-- 首页 (一级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_HOME', '', 'metadata-platform', '首页', 'home',
    1, '1', 1, 'C', 'fa-gauge-high',
    '/home/dashboard', 'GET', '', '首页仪表盘', 1, 1,
    0, 1, '1', '0', 'system'
);

-- 元数据管理 (一级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_METADATA', '', 'metadata-platform', '元数据管理', 'metadata',
    1, '1', 1, 'M', 'fa-folder-open',
    '/metadata', 'GET', '', '元数据管理模块', 2, 1,
    0, 1, '1', '0', 'system'
);

-- 系统监控 (一级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MONITOR', '', 'metadata-platform', '系统监控', 'monitor',
    1, '1', 1, 'M', 'fa-desktop',
    '/monitor', 'GET', '', '系统监控模块', 3, 1,
    0, 1, '1', '0', 'system'
);

-- 接口管理 (一级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_API', '', 'metadata-platform', '接口管理', 'api',
    1, '1', 1, 'M', 'fa-share-nodes',
    '/api', 'GET', '', '接口管理模块', 4, 1,
    0, 1, '1', '0', 'system'
);

-- SSO管理 (一级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO', '', 'metadata-platform', 'SSO管理', 'sso',
    1, '1', 1, 'M', 'fa-key',
    '/sso', 'GET', '', '单点登录管理模块', 5, 1,
    0, 1, '1', '0', 'system'
);

-- 系统设置 (一级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SYSTEM', '', 'metadata-platform', '系统设置', 'system',
    1, '1', 1, 'M', 'fa-gear',
    '/system', 'GET', '', '系统设置模块', 6, 1,
    0, 1, '1', '0', 'system'
);

-- ============================================
-- 二级菜单
-- ============================================

-- 元数据管理 > 数据源管理 (二级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_METADATA_DATASOURCE', 'MENU_METADATA', 'metadata-platform', '数据源管理', 'metadata:datasource',
    1, '1', 1, 'M', 'fa-network-wired',
    '/metadata/datasource', 'GET', '', '数据源管理', 1, 2,
    0, 1, '1', '0', 'system'
);

-- 元数据管理 > 元数据维护 (二级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_METADATA_MAINTENANCE', 'MENU_METADATA', 'metadata-platform', '元数据维护', 'metadata:maintenance',
    1, '1', 1, 'M', 'fa-wrench',
    '/metadata/maintenance', 'GET', '', '元数据维护', 2, 2,
    0, 1, '1', '0', 'system'
);

-- 元数据管理 > 模型维护 (二级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_METADATA_MODEL', 'MENU_METADATA', 'metadata-platform', '模型维护', 'metadata:model',
    1, '1', 1, 'M', 'fa-file',
    '/metadata/model', 'GET', '', '模型维护', 3, 2,
    0, 1, '1', '0', 'system'
);

-- 系统监控 > 监控仪表盘 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MONITOR_DASHBOARD', 'MENU_MONITOR', 'metadata-platform', '监控仪表盘', 'monitor:dashboard',
    1, '1', 1, 'C', 'fa-gauge-high',
    '/monitor/dashboard', 'GET', '', '监控仪表盘', 1, 2,
    0, 1, '1', '0', 'system'
);

-- 系统监控 > 性能分析 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MONITOR_PERFORMANCE', 'MENU_MONITOR', 'metadata-platform', '性能分析', 'monitor:performance',
    1, '1', 1, 'C', 'fa-chart-line',
    '/monitor/performance', 'GET', '', '性能分析', 2, 2,
    0, 1, '1', '0', 'system'
);

-- 接口管理 > 接口列表 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_API_LIST', 'MENU_API', 'metadata-platform', '接口列表', 'api:list',
    1, '1', 1, 'C', 'fa-list',
    '/api/list', 'GET', '', '接口列表', 1, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 租户管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_TENANT', 'MENU_SSO', 'metadata-platform', '租户管理', 'sso:tenant',
    1, '1', 1, 'C', 'fa-building',
    '/sso/tenant', 'GET', '', '租户管理', 1, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 应用列表 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_APP', 'MENU_SSO', 'metadata-platform', '应用列表', 'sso:app',
    1, '1', 1, 'C', 'fa-desktop',
    '/sso/app', 'GET', '', '应用列表', 2, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 组织管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_ORG', 'MENU_SSO', 'metadata-platform', '组织管理', 'sso:org',
    1, '1', 1, 'C', 'fa-school',
    '/sso/org', 'GET', '', '组织管理', 3, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 组织类型 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_ORG_KIND', 'MENU_SSO', 'metadata-platform', '组织类型', 'sso:orgKind',
    1, '1', 1, 'C', 'fa-layer-group',
    '/sso/orgKind', 'GET', '', '组织类型', 4, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 菜单管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_MENU', 'MENU_SSO', 'metadata-platform', '菜单管理', 'sso:menu',
    1, '1', 1, 'C', 'fa-bars',
    '/sso/menu', 'GET', '', '菜单管理', 5, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 职位管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_POS', 'MENU_SSO', 'metadata-platform', '职位管理', 'sso:pos',
    1, '1', 1, 'C', 'fa-briefcase',
    '/sso/pos', 'GET', '', '职位管理', 6, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 角色管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_ROLE', 'MENU_SSO', 'metadata-platform', '角色管理', 'sso:role',
    1, '1', 1, 'C', 'fa-user',
    '/sso/role', 'GET', '', '角色管理', 7, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 角色组 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_ROLE_GROUP', 'MENU_SSO', 'metadata-platform', '角色组', 'sso:roleGroup',
    1, '1', 1, 'C', 'fa-folder',
    '/sso/roleGroup', 'GET', '', '角色组', 8, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 用户组 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_USER_GROUP', 'MENU_SSO', 'metadata-platform', '用户组', 'sso:userGroup',
    1, '1', 1, 'C', 'fa-users',
    '/sso/userGroup', 'GET', '', '用户组', 9, 2,
    0, 1, '1', '0', 'system'
);

-- SSO管理 > 用户管理 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SSO_USER', 'MENU_SSO', 'metadata-platform', '用户管理', 'sso:user',
    1, '1', 1, 'C', 'fa-user-astronaut',
    '/sso/user', 'GET', '', '用户管理', 10, 2,
    0, 1, '1', '0', 'system'
);

-- 系统设置 > 系统配置 (二级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SYSTEM_SETTINGS', 'MENU_SYSTEM', 'metadata-platform', '系统配置', 'system:settings',
    1, '1', 1, 'C', 'fa-wrench',
    '/system/settings', 'GET', '', '系统配置', 1, 2,
    0, 1, '1', '0', 'system'
);

-- 系统设置 > 审计日志 (二级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_SYSTEM_AUDIT', 'MENU_SYSTEM', 'metadata-platform', '审计日志', 'system:audit',
    1, '1', 1, 'M', 'fa-clock',
    '/system/audit', 'GET', '', '审计日志', 2, 2,
    0, 1, '1', '0', 'system'
);

-- ============================================
-- 三级菜单
-- ============================================

-- 数据源管理 > 数据源列表 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_DATASOURCE_LIST', 'MENU_METADATA_DATASOURCE', 'metadata-platform', '数据源列表', 'metadata:datasource:list',
    1, '1', 1, 'C', 'fa-list',
    '/metadata/datasource/list', 'GET', '', '数据源列表', 1, 3,
    0, 1, '1', '0', 'system'
);

-- 数据源管理 > 添加数据源 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_DATASOURCE_CREATE', 'MENU_METADATA_DATASOURCE', 'metadata-platform', '添加数据源', 'metadata:datasource:create',
    1, '1', 1, 'C', 'fa-plus',
    '/metadata/datasource/create', 'GET', '', '添加数据源', 2, 3,
    0, 1, '1', '0', 'system'
);

-- 元数据维护 > 表列表 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MAINTENANCE_TABLE', 'MENU_METADATA_MAINTENANCE', 'metadata-platform', '表列表', 'metadata:maintenance:table',
    1, '1', 1, 'C', 'fa-table',
    '/metadata/maintenance/table', 'GET', '', '表列表', 1, 3,
    0, 1, '1', '0', 'system'
);

-- 元数据维护 > 视图列表 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MAINTENANCE_VIEW', 'MENU_METADATA_MAINTENANCE', 'metadata-platform', '视图列表', 'metadata:maintenance:view',
    1, '1', 1, 'C', 'fa-table-columns',
    '/metadata/maintenance/view', 'GET', '', '视图列表', 2, 3,
    0, 1, '1', '0', 'system'
);

-- 元数据维护 > 存储过程 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MAINTENANCE_PROCEDURE', 'MENU_METADATA_MAINTENANCE', 'metadata-platform', '存储过程', 'metadata:maintenance:procedure',
    1, '1', 1, 'C', 'fa-database',
    '/metadata/maintenance/procedure', 'GET', '', '存储过程', 3, 3,
    0, 1, '1', '0', 'system'
);

-- 元数据维护 > 函数 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MAINTENANCE_FUNCTION', 'MENU_METADATA_MAINTENANCE', 'metadata-platform', '函数', 'metadata:maintenance:function',
    1, '1', 1, 'C', 'fa-calculator',
    '/metadata/maintenance/function', 'GET', '', '函数', 4, 3,
    0, 1, '1', '0', 'system'
);

-- 元数据维护 > 字段列表 (三级目录)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MAINTENANCE_FIELD', 'MENU_METADATA_MAINTENANCE', 'metadata-platform', '字段列表', 'metadata:maintenance:field',
    1, '1', 1, 'M', 'fa-ticket',
    '/metadata/maintenance/field', 'GET', '', '字段列表', 5, 3,
    0, 1, '1', '0', 'system'
);

-- 模型维护 > 模型列表 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MODEL_LIST', 'MENU_METADATA_MODEL', 'metadata-platform', '模型列表', 'metadata:model:list',
    1, '1', 1, 'C', 'fa-list',
    '/metadata/model/list', 'GET', '', '模型列表', 1, 3,
    0, 1, '1', '0', 'system'
);

-- 模型维护 > SQL模型测试 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_MODEL_SQL_TEST', 'MENU_METADATA_MODEL', 'metadata-platform', 'SQL模型测试', 'metadata:model:sql-test',
    1, '1', 1, 'C', 'fa-desktop',
    '/metadata/model/sql-test', 'GET', '', 'SQL模型测试', 2, 3,
    0, 1, '1', '0', 'system'
);

-- 审计日志 > 登录日志 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_AUDIT_LOGIN_LOG', 'MENU_SYSTEM_AUDIT', 'metadata-platform', '登录日志', 'system:audit:login',
    1, '1', 1, 'C', 'fa-user',
    '/system/audit/login', 'GET', '', '登录日志', 1, 3,
    0, 1, '1', '0', 'system'
);

-- 审计日志 > 操作日志 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_AUDIT_OPERATION_LOG', 'MENU_SYSTEM_AUDIT', 'metadata-platform', '操作日志', 'system:audit:operation',
    1, '1', 1, 'C', 'fa-list',
    '/system/audit/operation', 'GET', '', '操作日志', 2, 3,
    0, 1, '1', '0', 'system'
);

-- 审计日志 > 数据变更 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_AUDIT_DATA_CHANGE_LOG', 'MENU_SYSTEM_AUDIT', 'metadata-platform', '数据变更', 'system:audit:data',
    1, '1', 1, 'C', 'fa-pen-to-square',
    '/system/audit/data', 'GET', '', '数据变更日志', 3, 3,
    0, 1, '1', '0', 'system'
);

-- 审计日志 > 访问日志 (三级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_AUDIT_ACCESS_LOG', 'MENU_SYSTEM_AUDIT', 'metadata-platform', '访问日志', 'system:audit:access',
    1, '1', 1, 'C', 'fa-desktop',
    '/system/audit/access', 'GET', '', '访问日志', 4, 3,
    0, 1, '1', '0', 'system'
);

-- ============================================
-- 四级菜单
-- ============================================

-- 字段列表 > 字段列表 (四级菜单)
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'MENU_FIELD_LIST', 'MENU_MAINTENANCE_FIELD', 'metadata-platform', '字段列表', 'metadata:maintenance:field:list',
    1, '1', 1, 'C', 'fa-list',
    '/metadata/maintenance/field/list', 'GET', '', '字段列表', 1, 4,
    0, 1, '1', '0', 'system'
);

-- ============================================
-- 按钮权限 (F类型)
-- ============================================

-- 数据源管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_DATASOURCE_ADD', 'MENU_DATASOURCE_LIST', 'metadata-platform', '新增数据源', 'metadata:datasource:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增数据源按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_DATASOURCE_EDIT', 'MENU_DATASOURCE_LIST', 'metadata-platform', '编辑数据源', 'metadata:datasource:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑数据源按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_DATASOURCE_DELETE', 'MENU_DATASOURCE_LIST', 'metadata-platform', '删除数据源', 'metadata:datasource:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除数据源按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_DATASOURCE_TEST', 'MENU_DATASOURCE_LIST', 'metadata-platform', '测试连接', 'metadata:datasource:test',
    1, '1', 1, 'F', '',
    '', 'POST', '', '测试数据源连接按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

-- 表列表按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_TABLE_SYNC', 'MENU_MAINTENANCE_TABLE', 'metadata-platform', '同步表', 'metadata:table:sync',
    1, '1', 1, 'F', '',
    '', 'POST', '', '同步表按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_TABLE_SAVE', 'MENU_MAINTENANCE_TABLE', 'metadata-platform', '保存表', 'metadata:table:save',
    1, '1', 1, 'F', '',
    '', 'POST', '', '保存表按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_TABLE_EDIT', 'MENU_MAINTENANCE_TABLE', 'metadata-platform', '编辑表', 'metadata:table:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑表按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_TABLE_DELETE', 'MENU_MAINTENANCE_TABLE', 'metadata-platform', '删除表', 'metadata:table:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除表按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_TABLE_VIEW', 'MENU_MAINTENANCE_TABLE', 'metadata-platform', '查看详情', 'metadata:table:view',
    1, '1', 1, 'F', '',
    '', 'GET', '', '查看表详情按钮', 5, 4,
    0, 1, '1', '0', 'system'
);

-- 视图列表按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_VIEW_SYNC', 'MENU_MAINTENANCE_VIEW', 'metadata-platform', '同步视图', 'metadata:view:sync',
    1, '1', 1, 'F', '',
    '', 'POST', '', '同步视图按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_VIEW_SAVE', 'MENU_MAINTENANCE_VIEW', 'metadata-platform', '保存视图', 'metadata:view:save',
    1, '1', 1, 'F', '',
    '', 'POST', '', '保存视图按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_VIEW_EDIT', 'MENU_MAINTENANCE_VIEW', 'metadata-platform', '编辑视图', 'metadata:view:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑视图按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_VIEW_DELETE', 'MENU_MAINTENANCE_VIEW', 'metadata-platform', '删除视图', 'metadata:view:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除视图按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_VIEW_VIEW', 'MENU_MAINTENANCE_VIEW', 'metadata-platform', '查看详情', 'metadata:view:view',
    1, '1', 1, 'F', '',
    '', 'GET', '', '查看视图详情按钮', 5, 4,
    0, 1, '1', '0', 'system'
);

-- 存储过程按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_PROCEDURE_SYNC', 'MENU_MAINTENANCE_PROCEDURE', 'metadata-platform', '同步存储过程', 'metadata:procedure:sync',
    1, '1', 1, 'F', '',
    '', 'POST', '', '同步存储过程按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_PROCEDURE_SAVE', 'MENU_MAINTENANCE_PROCEDURE', 'metadata-platform', '保存存储过程', 'metadata:procedure:save',
    1, '1', 1, 'F', '',
    '', 'POST', '', '保存存储过程按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_PROCEDURE_EDIT', 'MENU_MAINTENANCE_PROCEDURE', 'metadata-platform', '编辑存储过程', 'metadata:procedure:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑存储过程按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_PROCEDURE_DELETE', 'MENU_MAINTENANCE_PROCEDURE', 'metadata-platform', '删除存储过程', 'metadata:procedure:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除存储过程按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_PROCEDURE_VIEW', 'MENU_MAINTENANCE_PROCEDURE', 'metadata-platform', '查看详情', 'metadata:procedure:view',
    1, '1', 1, 'F', '',
    '', 'GET', '', '查看存储过程详情按钮', 5, 4,
    0, 1, '1', '0', 'system'
);

-- 函数按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_FUNCTION_SYNC', 'MENU_MAINTENANCE_FUNCTION', 'metadata-platform', '同步函数', 'metadata:function:sync',
    1, '1', 1, 'F', '',
    '', 'POST', '', '同步函数按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_FUNCTION_SAVE', 'MENU_MAINTENANCE_FUNCTION', 'metadata-platform', '保存函数', 'metadata:function:save',
    1, '1', 1, 'F', '',
    '', 'POST', '', '保存函数按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_FUNCTION_EDIT', 'MENU_MAINTENANCE_FUNCTION', 'metadata-platform', '编辑函数', 'metadata:function:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑函数按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_FUNCTION_DELETE', 'MENU_MAINTENANCE_FUNCTION', 'metadata-platform', '删除函数', 'metadata:function:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除函数按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_FUNCTION_VIEW', 'MENU_MAINTENANCE_FUNCTION', 'metadata-platform', '查看详情', 'metadata:function:view',
    1, '1', 1, 'F', '',
    '', 'GET', '', '查看函数详情按钮', 5, 4,
    0, 1, '1', '0', 'system'
);

-- 模型管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_MODEL_ADD', 'MENU_MODEL_LIST', 'metadata-platform', '新增模型', 'metadata:model:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增模型按钮', 1, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_MODEL_EDIT', 'MENU_MODEL_LIST', 'metadata-platform', '编辑模型', 'metadata:model:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑模型按钮', 2, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_MODEL_DELETE', 'MENU_MODEL_LIST', 'metadata-platform', '删除模型', 'metadata:model:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除模型按钮', 3, 4,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_MODEL_PUBLISH', 'MENU_MODEL_LIST', 'metadata-platform', '发布模型', 'metadata:model:publish',
    1, '1', 1, 'F', '',
    '', 'POST', '', '发布模型按钮', 4, 4,
    0, 1, '1', '0', 'system'
);

-- 接口管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_API_ADD', 'MENU_API_LIST', 'metadata-platform', '新增接口', 'api:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增接口按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_API_EDIT', 'MENU_API_LIST', 'metadata-platform', '编辑接口', 'api:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑接口按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_API_DELETE', 'MENU_API_LIST', 'metadata-platform', '删除接口', 'api:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除接口按钮', 3, 3,
    0, 1, '1', '0', 'system'
);

-- SSO用户管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_USER_ADD', 'MENU_SSO_USER', 'metadata-platform', '新增用户', 'sso:user:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增用户按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_USER_EDIT', 'MENU_SSO_USER', 'metadata-platform', '编辑用户', 'sso:user:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑用户按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_USER_DELETE', 'MENU_SSO_USER', 'metadata-platform', '删除用户', 'sso:user:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除用户按钮', 3, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_USER_RESET_PWD', 'MENU_SSO_USER', 'metadata-platform', '重置密码', 'sso:user:resetPwd',
    1, '1', 1, 'F', '',
    '', 'POST', '', '重置用户密码按钮', 4, 3,
    0, 1, '1', '0', 'system'
);

-- SSO角色管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ROLE_ADD', 'MENU_SSO_ROLE', 'metadata-platform', '新增角色', 'sso:role:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增角色按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ROLE_EDIT', 'MENU_SSO_ROLE', 'metadata-platform', '编辑角色', 'sso:role:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑角色按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ROLE_DELETE', 'MENU_SSO_ROLE', 'metadata-platform', '删除角色', 'sso:role:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除角色按钮', 3, 3,
    0, 1, '1', '0', 'system'
);

-- SSO菜单管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_MENU_ADD', 'MENU_SSO_MENU', 'metadata-platform', '新增菜单', 'sso:menu:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增菜单按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_MENU_EDIT', 'MENU_SSO_MENU', 'metadata-platform', '编辑菜单', 'sso:menu:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑菜单按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_MENU_DELETE', 'MENU_SSO_MENU', 'metadata-platform', '删除菜单', 'sso:menu:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除菜单按钮', 3, 3,
    0, 1, '1', '0', 'system'
);

-- SSO组织管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ORG_ADD', 'MENU_SSO_ORG', 'metadata-platform', '新增组织', 'sso:org:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增组织按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ORG_EDIT', 'MENU_SSO_ORG', 'metadata-platform', '编辑组织', 'sso:org:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑组织按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_ORG_DELETE', 'MENU_SSO_ORG', 'metadata-platform', '删除组织', 'sso:org:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除组织按钮', 3, 3,
    0, 1, '1', '0', 'system'
);

-- SSO租户管理按钮权限
INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_TENANT_ADD', 'MENU_SSO_TENANT', 'metadata-platform', '新增租户', 'sso:tenant:add',
    1, '1', 1, 'F', '',
    '', 'POST', '', '新增租户按钮', 1, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_TENANT_EDIT', 'MENU_SSO_TENANT', 'metadata-platform', '编辑租户', 'sso:tenant:edit',
    1, '1', 1, 'F', '',
    '', 'PUT', '', '编辑租户按钮', 2, 3,
    0, 1, '1', '0', 'system'
);

INSERT INTO `sso_menu` (
    `id`, `parent_id`, `app_code`, `menu_name`, `menu_code`, 
    `status`, `data_range`, `is_visible`, `menu_type`, `icon`, 
    `url`, `method`, `target`, `remark`, `sort`, `tier`, 
    `is_deleted`, `is_system`, `tenant_id`, `create_id`, `create_by`
) VALUES (
    'BTN_SSO_TENANT_DELETE', 'MENU_SSO_TENANT', 'metadata-platform', '删除租户', 'sso:tenant:delete',
    1, '1', 1, 'F', '',
    '', 'DELETE', '', '删除租户按钮', 3, 3,
    0, 1, '1', '0', 'system'
);
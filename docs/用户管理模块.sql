/*
 Navicat Premium Dump SQL
 
 Source Server         : 【本机】mysql
 Source Server Type    : MySQL
 Source Server Version : 90500 (9.5.0)
 Source Host           : localhost:3306
 Source Schema         : metadata_sso
 
 Target Server Type    : MySQL
 Target Server Version : 90500 (9.5.0)
 File Encoding         : 65001
 
 Date: 24/01/2026 10:24:27
 */
SET
    NAMES utf8mb4;

SET
    FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sso_tenant
-- ----------------------------
DROP TABLE IF EXISTS `sso_tenant`;
CREATE TABLE `sso_tenant` (
  `id` varchar(64) NOT NULL COMMENT '租户id',
  `tenant_name` varchar(128) DEFAULT '' COMMENT '租户名称',
  `tenant_code` varchar(64) DEFAULT '' COMMENT '租户编码',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态: 1 有效 0 无效',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
  `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
  `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
  `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
  `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
  `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='租户表';

-- ----------------------------
-- Table structure for sso_app
-- ----------------------------
DROP TABLE IF EXISTS `sso_applicaiton`;

CREATE TABLE `sso_applicaiton` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `parent_id` varchar(64) NOT NULL DEFAULT '' COMMENT '父级id',
    `application_name` varchar(32) DEFAULT '' COMMENT '子系统名字（中文名：比如 教务系统）',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `host` varchar(512) DEFAULT '' COMMENT '系统运行机器的域名或ip',
    `logo` varchar(256) DEFAULT '' COMMENT '系统logo',
    `remark` varchar(512) DEFAULT '' COMMENT '描述信息',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_app_application_code` (`application_code`) USING BTREE,
    KEY `idx_sso_application_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '应用服务模块表';

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
    `id` varchar(64) NOT NULL COMMENT '主键',
    `ptype` varchar(100) DEFAULT NULL,
    `v0` varchar(100) DEFAULT NULL,
    `v1` varchar(100) DEFAULT NULL,
    `v2` varchar(100) DEFAULT NULL,
    `v3` varchar(100) DEFAULT NULL,
    `v4` varchar(100) DEFAULT NULL,
    `v5` varchar(100) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for sso_menu
-- ----------------------------
DROP TABLE IF EXISTS `sso_menu`;

CREATE TABLE `sso_menu` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `menu_name` varchar(128) DEFAULT '' COMMENT '菜单名称',
    `menu_code` varchar(128) DEFAULT '' COMMENT '标识权限',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用（1 可用 0 禁用）',
    `data_scope` char(1) DEFAULT '1' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）',
    `visible` int NOT NULL DEFAULT '1' COMMENT '菜单状态（1 显示 0 隐藏）',
    `menu_type` char(1) DEFAULT '' COMMENT '菜单类型（M目录 C菜单 F按钮 Z资源）',
    `icon` varchar(128) DEFAULT '' COMMENT '菜单图标',
    `url` varchar(512) DEFAULT '#' COMMENT '请求地址',
    `method` varchar(16) DEFAULT '' COMMENT 'HTTP方法 (GET, POST, PUT, DELETE)',
    `target` varchar(36) DEFAULT '' COMMENT '打开方式',
    `remark` varchar(512) DEFAULT '' COMMENT '描述',
    `sort` int DEFAULT '0' COMMENT '排序',
    `tier` int DEFAULT '0' COMMENT '树层级',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_menu_menu_code` (`menu_code`) USING BTREE,
    KEY `idx_sso_menu_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '菜单权限表';

-- ----------------------------
-- Table structure for sso_organization_menu
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization_menu`;

CREATE TABLE `sso_organization_menu` (
    `id` varchar(64) NOT NULL,

    `menu_id` varchar(64) DEFAULT '0',
    `organization_id` varchar(64) DEFAULT '0',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_menu_organization_menu` (`menu_id`, `organization_id`) USING BTREE,
    KEY `idx_sso_organization_menu_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 914000000 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '菜单和部门关联表';

-- ----------------------------
-- Table structure for sso_position
-- ----------------------------
DROP TABLE IF EXISTS `sso_position`;

CREATE TABLE `sso_position` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `organization_id` varchar(64) DEFAULT '0' COMMENT '组织专属岗位',
    `kind_code` varchar(64) DEFAULT '' COMMENT '组织类型专属岗位',
    `pos_name` varchar(128) DEFAULT '' COMMENT '职位名称',
    `pos_code` varchar(64) DEFAULT '' COMMENT '职位编码',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `remark` varchar(512) DEFAULT '' COMMENT '描述信息',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_position_pos_code` (`pos_code`) USING BTREE,
    KEY `idx_sso_position_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '职位表';

-- ----------------------------
-- Table structure for sso_position_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_position_role`;

CREATE TABLE `sso_position_role` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `pos_id` varchar(64) DEFAULT '0' COMMENT '职位id',
    `role_id` varchar(64) DEFAULT '0' COMMENT '角色id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_position_role_pos_role` (`pos_id`, `role_id`) USING BTREE,
    KEY `idx_sso_position_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 3 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '职位角色表';

-- ----------------------------
-- Table structure for sso_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_role`;

CREATE TABLE `sso_role` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `organization_id` varchar(64) DEFAULT '0' COMMENT '组织专属岗位',
    `kind_code` varchar(64) DEFAULT '' COMMENT '组织类型专属岗位',
    `role_name` varchar(128) DEFAULT '' COMMENT '角色名称',
    `role_code` varchar(128) DEFAULT '' COMMENT '角色代码',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `data_scope` char(1) DEFAULT '1' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）',
    `remark` varchar(512) DEFAULT '' COMMENT '描述',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_role_role_code` (`role_code`) USING BTREE,
    KEY `idx_sso_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 10 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色管理表';

-- ----------------------------
-- Table structure for sso_role_group
-- ----------------------------
DROP TABLE IF EXISTS `sso_role_group`;

CREATE TABLE `sso_role_group` (
    `id` varchar(64) NOT NULL,

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `organization_id` varchar(64) DEFAULT '0' COMMENT '组织专属岗位',
    `kind_code` varchar(64) DEFAULT '' COMMENT '组织类型专属岗位',
    `group_name` varchar(128) DEFAULT '' COMMENT '角色组名称',
    `group_code` varchar(128) DEFAULT '' COMMENT '角色组编码',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `remark` varchar(512) DEFAULT '' COMMENT '备注',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_role_group_group_code` (`group_code`) USING BTREE,
    KEY `idx_sso_role_group_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色组表';

-- ----------------------------
-- Table structure for sso_role_group_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_role_group_role`;

CREATE TABLE `sso_role_group_role` (
    `id` varchar(64) NOT NULL,

    `group_id` varchar(64) DEFAULT '0',
    `role_id` varchar(64) DEFAULT '0',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_role_group_role_group_role` (`group_id`, `role_id`) USING BTREE,
    KEY `idx_sso_role_group_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色组关联角色表';

-- ----------------------------
-- Table structure for sso_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sso_role_menu`;

CREATE TABLE `sso_role_menu` (
    `id` varchar(64) NOT NULL COMMENT '物理id',

    `role_id` varchar(64) DEFAULT '0' COMMENT '角色id',
    `menu_id` varchar(64) DEFAULT '0' COMMENT '菜单id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_role_menu_role_menu` (`role_id`, `menu_id`) USING BTREE,
    KEY `idx_sso_role_menu_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 6919 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色菜单表';

-- ----------------------------
-- Table structure for sso_organization_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization_role`;

CREATE TABLE `sso_organization_role` (
    `id` varchar(64) NOT NULL,

    `organization_id` varchar(64) DEFAULT '0',
    `role_id` varchar(64) DEFAULT '0',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_role_organization_role` (`role_id`, `organization_id`) USING BTREE,
    KEY `idx_sso_organization_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 914000061 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '角色和部门关联表';

-- ----------------------------
-- Table structure for sso_organization
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization`;

CREATE TABLE `sso_organization` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `from_id` varchar(64) DEFAULT '' COMMENT '来源id，数据同步使用',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `unit_name` varchar(128) DEFAULT '' COMMENT '组织名称',
    `unit_short` varchar(128) DEFAULT '' COMMENT '组织简称',
    `unit_en` varchar(128) DEFAULT '' COMMENT '英文全称',
    `unit_en_short` varchar(128) DEFAULT '' COMMENT '英文简称',
    `unit_code` varchar(64) DEFAULT '' COMMENT '组织编码',
    `kind_code` varchar(64) DEFAULT '' COMMENT '类型编码',
    `logo` varchar(256) DEFAULT '' COMMENT '系统logo',
    `host` varchar(512) DEFAULT '' COMMENT '域名或ip',
    `contact` varchar(128) DEFAULT '' COMMENT '联系人姓名',
    `phone` varchar(128) DEFAULT '' COMMENT '联系电话',
    `address` varchar(256) DEFAULT '' COMMENT '联系地址',
    `postcode` varchar(16) DEFAULT '' COMMENT '邮编',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `remark` varchar(512) DEFAULT '' COMMENT '描述信息',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_unit_code` (`unit_code`) USING BTREE,
    KEY `idx_sso_organization_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 57 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '组织表';

-- ----------------------------
-- Table structure for sso_organization_kind
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization_kind`;

CREATE TABLE `sso_organization_kind` (
    `id` varchar(64) NOT NULL,

    `parent_id` varchar(64) NOT NULL DEFAULT '' COMMENT '父级id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `kind_name` varchar(100) DEFAULT '' DEFAULT '' COMMENT '名称',
    `kind_code` varchar(64) DEFAULT '' COMMENT '编码',
    `kind_tag` varchar(64) DEFAULT '' COMMENT '等级标识',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `remark` varchar(512) DEFAULT '' COMMENT '描述信息',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_kind_kind_code` (`kind_code`) USING BTREE,
    KEY `idx_sso_organization_kind_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 901000013 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '组织类型表';

-- ----------------------------
-- Table structure for sso_organization_kind_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization_kind_role`;

CREATE TABLE `sso_organization_kind_role` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `kind_code` varchar(64) DEFAULT '' COMMENT '编码',
    `role_id` varchar(64) DEFAULT '0' COMMENT '角色id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_kind_role_kind_role` (`kind_code`, `role_id`) USING BTREE,
    KEY `idx_sso_organization_kind_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '组织类型角色表';

-- ----------------------------
-- Table structure for sso_user
-- ----------------------------
DROP TABLE IF EXISTS `sso_user`;

CREATE TABLE `sso_user` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `account_id` varchar(64) DEFAULT '' COMMENT '用户编号，系统自动生成',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `account` varchar(128) DEFAULT NULL COMMENT '帐号',
    `password` varchar(64) DEFAULT NULL COMMENT '密码',
    `salt` varchar(32) DEFAULT '' COMMENT '盐加密',
    `name` varchar(128) DEFAULT NULL COMMENT '用户姓名',
    `code` varchar(64) DEFAULT NULL COMMENT '扩展编号：如教师工号',
    `sex` varchar(6) DEFAULT '男' COMMENT '性别（男、女）',
    `idcard` varchar(32) DEFAULT '' COMMENT '身份证件',
    `mobile` varchar(32) DEFAULT '' COMMENT '手机号',
    `email` varchar(128) DEFAULT '' COMMENT '电子邮箱',
    `avatar` varchar(128) DEFAULT '' COMMENT '头像',
    `organization_id` varchar(64) DEFAULT '' COMMENT '组织id',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `end_time` datetime DEFAULT NULL COMMENT '结束时间',
    `kind` int DEFAULT '2' COMMENT '分类（1：超级管理员 2：子管理员 99：其他）',
    `remark` varchar(512) DEFAULT '' COMMENT '描述信息',
    `sort` int DEFAULT '0' COMMENT '排序',
    `first_login` int DEFAULT '0' COMMENT '第一次登陆:0',
    `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
    `last_ip` varchar(32) DEFAULT '' COMMENT '最后登录IP',
    `login_error_count` int DEFAULT '0' COMMENT '登录次数',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_account` (`account`) USING BTREE,
    KEY `idx_sso_user_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 906000067 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户表';

-- ----------------------------
-- Table structure for sso_user_group
-- ----------------------------
DROP TABLE IF EXISTS `sso_user_group`;

CREATE TABLE `sso_user_group` (
    `id` varchar(64) NOT NULL,

    `parent_id` varchar(64) DEFAULT '' COMMENT '父id',
    `application_code` varchar(64) DEFAULT '' COMMENT '子系统编码（比如：zuul）',
    `organization_id` varchar(64) DEFAULT '0' COMMENT '组织专属岗位',
    `kind_code` varchar(64) DEFAULT '' COMMENT '组织类型专属岗位',
    `group_name` varchar(128) DEFAULT '' COMMENT '角色组名称',
    `group_code` varchar(128) DEFAULT '' COMMENT '角色组编码',
    `status` int NOT NULL DEFAULT '1' COMMENT '是否可用:1 可用 0 禁用',
    `remark` varchar(512) DEFAULT '' COMMENT '备注',
    `sort` int DEFAULT '0' COMMENT '排序',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_group_group_code` (`group_code`) USING BTREE,
    KEY `idx_sso_user_group_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 916000006 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户组表';

-- ----------------------------
-- Table structure for sso_user_group_user
-- ----------------------------
DROP TABLE IF EXISTS `sso_user_group_user`;

CREATE TABLE `sso_user_group_user` (
    `id` varchar(64) NOT NULL,

    `group_id` varchar(64) DEFAULT '',
    `user_id` varchar(64) DEFAULT '',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_group_user_group_user` (`group_id`, `user_id`) USING BTREE,
    KEY `idx_sso_user_group_user_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户组关联用户表';

-- ----------------------------
-- Table structure for sso_user_position
-- ----------------------------
DROP TABLE IF EXISTS `sso_user_position`;

CREATE TABLE `sso_user_position` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `user_id` varchar(64) DEFAULT '' COMMENT '用户id',
    `position_id` varchar(64) DEFAULT '' COMMENT '职位id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_position_user_position` (`user_id`, `position_id`) USING BTREE,
    KEY `idx_sso_user_position_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 3 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户职位表';


-- ----------------------------
-- Table structure for sso_user_role
-- ----------------------------
DROP TABLE IF EXISTS `sso_user_role`;

CREATE TABLE `sso_user_role` (
    `id` varchar(64) NOT NULL COMMENT '主键',

    `user_id` varchar(64) DEFAULT '' COMMENT '用户id',
    `role_id` varchar(64) DEFAULT '' COMMENT '角色id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_role_user_role` (`user_id`, `role_id`) USING BTREE,
    KEY `idx_sso_user_role_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 61 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户角色表';

-- ----------------------------
-- Table structure for sso_user_role_group
-- ----------------------------
DROP TABLE IF EXISTS `sso_user_role_group`;

CREATE TABLE `sso_user_role_group` (
    `id` varchar(64) NOT NULL,

    `user_id` varchar(64) DEFAULT '0' COMMENT '用户id',
    `group_id` varchar(64) DEFAULT '0' COMMENT '角色组id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_user_role_group_user_role_group` (`user_id`, `group_id`) USING BTREE,
    KEY `idx_sso_user_role_group_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户角色组关联表';

-- ----------------------------
-- Table structure for sso_organization_user
-- ----------------------------
DROP TABLE IF EXISTS `sso_organization_user`;

CREATE TABLE `sso_organization_user` (
    `id` varchar(64) NOT NULL,

    `organization_id` varchar(64) DEFAULT '0' COMMENT '部门id',
    `user_id` varchar(64) DEFAULT '0' COMMENT '用户id',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `index_sso_organization_user_user_unit` (`user_id`, `organization_id`) USING BTREE,
    KEY `idx_sso_organization_user_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 18 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '组织包含用户表';

SET
    FOREIGN_KEY_CHECKS = 1;
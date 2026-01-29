/*
 Navicat Premium Dump SQL
 
 Source Server         : 【本机】mysql
 Source Server Type    : MySQL
 Source Server Version : 90500 (9.5.0)
 Source Host           : localhost:3306
 Source Schema         : metadata_platform
 
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
-- Table structure for md_conn
-- ----------------------------
DROP TABLE IF EXISTS `md_conn`;

CREATE TABLE `md_conn` (
    `id` varchar(64) NOT NULL COMMENT '数据连接主键',

    `parent_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '父级id',
    `conn_name` varchar(256) DEFAULT '' COMMENT '数据连接名称',
    `conn_kind` varchar(64) DEFAULT '' COMMENT '数据连接类型（例如MySQL, Oracle, SQLServer, DB2, DM, KingbaseES）',
    `conn_version` varchar(64) DEFAULT '' COMMENT '数据库版本（例如8.0, 12c, 2019）',
    `conn_host` varchar(128) DEFAULT '' COMMENT '数据连接主机地址',
    `conn_port` int NOT NULL DEFAULT '0' COMMENT '数据连接端口号',
    `conn_user` varchar(128) DEFAULT '' COMMENT '用户名',
    `conn_password` varchar(128) DEFAULT '' COMMENT '密码',
    `conn_database` varchar(128) DEFAULT '' COMMENT '数据库',
    `conn_conn` varchar(1024) DEFAULT '' COMMENT '链接地址：自动生成',
    `state` int NOT NULL DEFAULT '0' COMMENT '连接状态: 0=未检测, 1=有效',
    `remark` varchar(512) DEFAULT '' COMMENT '备注',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 800100017 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '数据连接';

-- ----------------------------
-- Table structure for md_table
-- ----------------------------
DROP TABLE IF EXISTS `md_table`;

CREATE TABLE `md_table` (
    `id` varchar(64) NOT NULL COMMENT '字段主键',

    `conn_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '数据连接ID',
    `conn_name` varchar(256) DEFAULT '' COMMENT '数据连接名称',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `table_type` varchar(64) DEFAULT '' COMMENT '表类型',
    `table_comment` varchar(256) DEFAULT '' COMMENT '表描述',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 800200323 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '数据连接表';

-- ----------------------------
-- Table structure for md_table_field
-- ----------------------------
DROP TABLE IF EXISTS `md_table_field`;

CREATE TABLE `md_table_field` (
    `id` varchar(64) NOT NULL COMMENT '字段主键',

    `conn_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '数据连接ID',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_title` varchar(256) DEFAULT '' COMMENT '字段标题',
    `column_type` varchar(64) DEFAULT '' COMMENT '数据类型，例如INT、VARCHAR(255)、TIMESTAMP等',
    `column_length` int DEFAULT '0' COMMENT '字段长度',
    `column_comment` varchar(256) DEFAULT '' COMMENT '字段描述',
    `is_nullable` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否允许为空',
    `is_primary_key` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为主键',
    `is_auto_increment` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否自增',
    `default_value` varchar(256) DEFAULT '' COMMENT '默认值',
    `extra_info` varchar(1024) DEFAULT '' COMMENT '额外信息（如auto_increment, unique等）',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 800305548 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '数据连接表字段';

-- ----------------------------
-- Table structure for md_model
-- ----------------------------
DROP TABLE IF EXISTS `md_model`;

CREATE TABLE `md_model` (
    `id` varchar(64) NOT NULL COMMENT '模型ID',

    `parent_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '父级id',
    `conn_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '数据连接ID',
    `conn_name` varchar(256) DEFAULT '' COMMENT '数据连接名称',
    `model_name` varchar(128) DEFAULT '' COMMENT '模型名称',
    `model_code` varchar(128) DEFAULT '' COMMENT '模型编码',
    `model_version` varchar(64)  DEFAULT '1.0.0' COMMENT '模型版本',
    `model_logo` varchar(512) DEFAULT '' COMMENT '模型图片',
    `model_kind` int NOT NULL DEFAULT '0' COMMENT '模型类型：1sql语句、2视图/表、3存储过程、4关联',
    `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开',
    `is_locked` tinyint(1) DEFAULT '0' COMMENT '是否锁定',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uix_md_model_title_creator` (`model_code`, `create_by`, `tenant_id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型定义';

-- ----------------------------
-- Table structure for md_model_group
-- ----------------------------
DROP TABLE IF EXISTS `md_model_group`;

CREATE TABLE `md_model_group` (
    `id` varchar(64) NOT NULL COMMENT '模型分组ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_title` varchar(256) DEFAULT '' COMMENT '字段标题',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `agg_func` varchar(64) DEFAULT '' COMMENT '聚合函数：sum/count/avg/max/min',
    `is_show` int NOT NULL DEFAULT '0' COMMENT '展示: 1 有效 0 无效',
    `show_title` varchar(128) DEFAULT '' COMMENT '字段显示名称',
    `show_width` int NOT NULL DEFAULT '100' COMMENT '字段显示宽度',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-选择字段';

-- ----------------------------
-- Table structure for md_model_having
-- ----------------------------
DROP TABLE IF EXISTS `md_model_having`;

CREATE TABLE `md_model_having` (
    `id` varchar(64) NOT NULL COMMENT '模型having条件ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `operator1` varchar(64) DEFAULT '' COMMENT '操作符：and/or',
    `brackets1` varchar(64) DEFAULT '' COMMENT '括号：(',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_title` varchar(256) DEFAULT '' COMMENT '字段标题',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `operator2` varchar(64) DEFAULT '' COMMENT '运算符：=/</>/<=/>=/like/between',
    `having_table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `having_table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '比较表ID',
    `having_table_name` varchar(256) DEFAULT '' COMMENT '比较表名称',
    `having_table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `having_column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '比较字段ID',
    `having_column_name` varchar(256) DEFAULT '' COMMENT '比较字段名称',
    `having_column_title` varchar(256) DEFAULT '' COMMENT '比较字段标题',
    `having_func` varchar(64) DEFAULT '' COMMENT '比较字段函数',
    `value1` varchar(128) DEFAULT '' COMMENT '值1',
    `value2` varchar(128) DEFAULT '' COMMENT '值2',
    `brackets2` varchar(64) DEFAULT '' COMMENT '括号：)',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-having条件';

-- ----------------------------
-- Table structure for md_model_join
-- ----------------------------
DROP TABLE IF EXISTS `md_model_join`;

CREATE TABLE `md_model_join` (
    `id` varchar(64) NOT NULL COMMENT '模型连接ID',

    `parent_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '父级id',
    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `join_type` varchar(64) DEFAULT '' COMMENT '关联类型：Left Join/Right Join/Inner Join',
    `operator1` varchar(64) DEFAULT '' COMMENT '操作符：and/or',
    `brackets1` varchar(64) DEFAULT '' COMMENT '括号：(',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_title` varchar(256) DEFAULT '' COMMENT '字段标题',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `operator2` varchar(64)  DEFAULT '=' COMMENT '运算符：=/</>/<=/>=/like/between',
    `join_table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `join_table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '关联表ID',
    `join_table_name` varchar(256) DEFAULT '' COMMENT '关联表',
    `join_table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `join_column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '关联字段ID',
    `join_column_name` varchar(256) DEFAULT '' COMMENT '关联字段名称',
    `join_column_title` varchar(256) DEFAULT '' COMMENT '关联字段标题',
    `join_func` varchar(256) DEFAULT '' COMMENT '关联字段函数',
    `value1` varchar(128) DEFAULT '' COMMENT '值1',
    `value2` varchar(128) DEFAULT '' COMMENT '值2',
    `brackets2` varchar(64) DEFAULT '' COMMENT '括号：)',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-关联';

-- ----------------------------
-- Table structure for md_model_limit
-- ----------------------------
DROP TABLE IF EXISTS `md_model_limit`;

CREATE TABLE `md_model_limit` (
    `id` varchar(64) NOT NULL COMMENT '模型限制ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `page` int NOT NULL DEFAULT '0' COMMENT '页码',
    `limit` int NOT NULL DEFAULT '0' COMMENT '每页数量',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-获取条数';

-- ----------------------------
-- Table structure for md_model_order
-- ----------------------------
DROP TABLE IF EXISTS `md_model_order`;

CREATE TABLE `md_model_order` (
    `id` varchar(64) NOT NULL COMMENT '模型排序ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_alias` varchar(256) DEFAULT '' COMMENT '字段别名',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `order_type` varchar(128)  DEFAULT 'Asc' COMMENT '排序：Asc/Desc',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-排序字段';

-- ----------------------------
-- Table structure for md_model_field
-- ----------------------------
DROP TABLE IF EXISTS `md_model_field`;

CREATE TABLE `md_model_field` (
    `id` varchar(64) NOT NULL COMMENT '模型查询ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_alias` varchar(256) DEFAULT '' COMMENT '字段别名',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `agg_func` varchar(64) DEFAULT '' COMMENT '聚合函数：sum/count/avg/max/min',
    `show_title` varchar(128) DEFAULT '' COMMENT '字段显示名称',
    `show_width` int NOT NULL DEFAULT '100' COMMENT '字段显示宽度',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-选择字段';

-- ----------------------------
-- Table structure for md_model_sql
-- ----------------------------
DROP TABLE IF EXISTS `md_model_sql`;

CREATE TABLE `md_model_sql` (
    `id` varchar(64) NOT NULL COMMENT '模型SQLID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COMMENT '查询内容：sql语句',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-sql';

-- ----------------------------
-- Table structure for md_model_table
-- ----------------------------
DROP TABLE IF EXISTS `md_model_table`;

CREATE TABLE `md_model_table` (
    `id` varchar(64) NOT NULL COMMENT '模型表ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `is_main` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否主表',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-表';

-- ----------------------------
-- Table structure for md_model_where
-- ----------------------------
DROP TABLE IF EXISTS `md_model_where`;

CREATE TABLE `md_model_where` (
    `id` varchar(64) NOT NULL COMMENT '模型条件ID',

    `model_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '模型ID',
    `operator1` varchar(64) DEFAULT '' COMMENT '操作符：and/or',
    `brackets1` varchar(64) DEFAULT '' COMMENT '括号：(',
    `table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '表ID',
    `table_name` varchar(256) DEFAULT '' COMMENT '表名称',
    `table_title` varchar(256) DEFAULT '' COMMENT '表标题',
    `column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '字段ID',
    `column_name` varchar(256) DEFAULT '' COMMENT '字段名称',
    `column_title` varchar(256) DEFAULT '' COMMENT '字段标题',
    `func` varchar(256) DEFAULT '' COMMENT '字段函数',
    `operator2` varchar(64) DEFAULT '' COMMENT '运算符：=/</>/<=/>=/like/between',
    `where_table_schema` varchar(64) DEFAULT '' COMMENT '表模式',
    `where_table_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '比较表ID',
    `where_table_name` varchar(256) DEFAULT '' COMMENT '比较表名称',
    `where_table_title` varchar(256) DEFAULT '' COMMENT '比较表标题',
    `where_column_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '比较字段ID',
    `where_column_name` varchar(256) DEFAULT '' COMMENT '比较字段名称',
    `where_column_title` varchar(256) DEFAULT '' COMMENT '比较字段标题',
    `where_func` varchar(64) DEFAULT '' COMMENT '比较字段函数',
    `value1` varchar(128) DEFAULT '' COMMENT '值1',
    `value2` varchar(128) DEFAULT '' COMMENT '值2',
    `brackets2` varchar(64) DEFAULT '' COMMENT '括号：)',
    `is_deleted` tinyint(1) DEFAULT '0' COMMENT '删除标识',
    `tenant_id` varchar(64) NOT NULL DEFAULT '0' COMMENT '租户ID',
    `create_id` varchar(64) DEFAULT '0' COMMENT '创建人id',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建人',
    `create_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_id` varchar(64) DEFAULT '0' COMMENT '修改人id',
    `update_by` varchar(64) DEFAULT '' COMMENT '修改人',
    `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`) USING BTREE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_bin ROW_FORMAT = DYNAMIC COMMENT = '模型-where条件';

SET
    FOREIGN_KEY_CHECKS = 1;
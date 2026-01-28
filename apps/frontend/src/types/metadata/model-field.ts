import type { BaseEntity } from './model-base';

/**
 * 模型字段信息
 */
export interface ModelField extends BaseEntity {
    /** 模型ID */
    modelID: string;
    /** 表模式 */
    tableSchema: string;
    /** 表ID */
    tableID: string;
    /** 表名称 */
    tableName: string;
    /** 表标题 */
    tableTitle: string;
    /** 字段ID */
    columnID: string;
    /** 字段名称 */
    columnName: string;
    /** 字段标题 */
    columnTitle: string;
    /** 字段函数 */
    func: string;
    /** 聚合函数：sum/count/avg/max/min */
    aggFunc: string;
    /** 展示: 1 有效 0 无效 */
    isShow: number;
    /** 字段显示名称 */
    showTitle: string;
    /** 字段显示宽度 */
    showWidth: number;
    /** 状态: 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;

    /** 是否允许为空 */
    isNullable: boolean;
    /** 是否为主键 */
    isPrimaryKey: boolean;
    /** 默认值 */
    defaultValue: string;
    /** 字段正则表达式校验 */
    fieldRegex: string;
    /** 字段最小长度 */
    fieldMinLen: number;
    /** 字段最大长度 */
    fieldMaxLen: number;
    /** 字段最小值 */
    fieldMin: number;
    /** 字段最大值 */
    fieldMax: number;

    /** 是否可搜索 */
    isSearchable: boolean;
    /** 是否可排序 */
    isSortable: boolean;
    /** 是否可筛选 */
    isFilterable: boolean;
    /** 组件类型 */
    componentType: string;
}

/**
 * 租户信息
 */
export interface Tenant extends BaseEntity {
    /** 租户名称 */
    tenantName: string;
    /** 租户编码 */
    tenantCode: string;
    /** 联系人 */
    linkman: string;
    /** 联系电话 */
    contact: string;
    /** 联系地址 */
    address: string;
    /** 状态: 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
}

/**
 * 用户信息
 */
export interface User extends BaseEntity {
    /** 账号ID */
    accountID: string;
    /** 子系统编码 */
    svcCode: string;
    /** 账号 */
    account: string;
    /** 姓名 */
    name: string;
    /** 扩展编号 */
    code: string;
    /** 性别 */
    sex: string;
    /** 身份证件 */
    idcard: string;
    /** 手机号 */
    mobile: string;
    /** 电子邮箱 */
    email: string;
    /** 头像 */
    avatar: string;
    /** 组织ID */
    unitID: string;
    /** 学校 */
    school: string;
    /** 班级 */
    class: string;
    /** 状态: 1 有效 0 无效 */
    state: number;
    /** 结束时间 */
    endTime: string;
    /** 分类：1超级管理员 2子管理员 3教师 4学生 99其他 */
    kind: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
    /** 第一次登录 */
    firstLogin: number;
    /** 最后登录时间 */
    lastLoginTime: string;
    /** 最后登录IP */
    lastIP: string;
    /** 登录错误次数 */
    loginErrorCount: number;
}

/**
 * 菜单信息
 */
export interface Menu extends BaseEntity {
    /** 父级ID */
    parentID: string;
    /** 子系统编码 */
    svcCode: string;
    /** 菜单名称 */
    menuName: string;
    /** 菜单编码 */
    menuCode: string;
    /** 数据范围 */
    dataScope: string;
    /** 可见性：1显示 0隐藏 */
    visible: number;
    /** 菜单类型：M目录 C菜单 F按钮 Z资源 */
    menuType: string;
    /** 菜单图标 */
    icon: string;
    /** 请求地址 */
    url: string;
    /** HTTP方法 */
    method: string;
    /** 打开方式 */
    target: string;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
    /** 树层级 */
    tier: number;
    /** 状态: 1 有效 0 无效 */
    state: number;
}

/**
 * 角色信息
 */
export interface Role extends BaseEntity {
    /** 父级ID */
    parentID: string;
    /** 子系统编码 */
    svcCode: string;
    /** 组织ID */
    unitID: string;
    /** 类型编码 */
    kindCode: string;
    /** 角色名称 */
    roleName: string;
    /** 角色编码 */
    roleCode: string;
    /** 数据范围 */
    dataScope: string;
    /** 状态: 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
}
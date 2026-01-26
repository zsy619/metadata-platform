/**
 * 用户管理模块类型定义
 */

/**
 * 基础实体类型
 */
export interface BaseEntity {
  /** 主键ID */
  id: string;
  /** 租户ID */
  tenantID: string;
  /** 创建时间 */
  createAt: string;
  /** 更新时间 */
  updateAt: string;
  /** 删除时间 */
  deleteAt?: string;
}

/**
 * 用户模型
 */
export interface User extends BaseEntity {
  /** 用户编号，系统自动生成 */
  accountID: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 帐号 */
  account: string;
  /** 密码 */
  password?: string;
  /** 盐加密 */
  salt: string;
  /** 用户姓名 */
  name: string;
  /** 扩展编号：如教师工号 */
  code: string;
  /** 性别（男、女） */
  sex: string;
  /** 身份证件 */
  idcard: string;
  /** 手机号 */
  mobile: string;
  /** 电子邮箱 */
  email: string;
  /** 头像 */
  avatar: string;
  /** 组织id */
  unitID: string;
  /** 学校 */
  school: string;
  /** 班级 */
  class: string;
  /** 是否可用:1 可用 0 禁用 */
  state: number;
  /** 结束时间 */
  endTime: string;
  /** 分类（1：超级管理员 2：子管理员 3：教师 4：学生 99：其他） */
  kind: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 第一次登陆:0 */
  firstLogin: number;
  /** 最后登录时间 */
  lastLoginTime: string;
  /** 最后登录IP */
  lastIP: string;
  /** 登录次数 */
  loginErrorCount: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 租户模型
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
  /** 备注 */
  remark: string;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 应用服务模块模型
 */
export interface App extends BaseEntity {
  /** 父级id */
  parentID: string;
  /** 子系统名字（中文名：比如 教务系统） */
  svcName: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 是否可用:1 可用 0 禁用 */
  state: number;
  /** 系统运行机器的域名或ip */
  host: string;
  /** 系统logo */
  logo: string;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 菜单权限模型
 */
export interface Menu extends BaseEntity {
  /** 父id */
  parentID: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 菜单名称 */
  menuName: string;
  /** 标识权限 */
  menuCode: string;
  /** 是否可用（1 可用 0 禁用） */
  state: number;
  /** 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限） */
  dataScope: string;
  /** 菜单状态（1 显示 0 隐藏） */
  visible: number;
  /** 菜单类型（M目录 C菜单 F按钮 Z资源） */
  menuType: string;
  /** 菜单图标 */
  icon: string;
  /** 请求地址 */
  url: string;
  /** HTTP方法 (GET, POST, PUT, DELETE) */
  method: string;
  /** 打开方式 */
  target: string;
  /** 描述 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 树层级 */
  tier: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 角色管理模型
 */
export interface Role extends BaseEntity {
  /** 父id */
  parentID: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 组织专属岗位 */
  unitID: string;
  /** 组织类型专属岗位 */
  kindCode: string;
  /** 角色名称 */
  roleName: string;
  /** 角色代码 */
  roleCode: string;
  /** 是否可用:1 可用 0 禁用 */
  state: number;
  /** 数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限） */
  dataScope: string;
  /** 描述 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 组织模型
 */
export interface Unit extends BaseEntity {
  /** 父id */
  parentID: string;
  /** 来源id，数据同步使用 */
  fromID: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 组织名称 */
  unitName: string;
  /** 组织简称 */
  unitShort: string;
  /** 英文全称 */
  unitEn: string;
  /** 英文简称 */
  unitEnShort: string;
  /** 组织编码 */
  unitCode: string;
  /** 类型编码 */
  kindCode: string;
  /** 系统logo */
  logo: string;
  /** 域名或ip */
  host: string;
  /** 联系人姓名 */
  contact: string;
  /** 联系电话 */
  phone: string;
  /** 联系地址 */
  address: string;
  /** 邮编 */
  postcode: string;
  /** 是否可用:1 可用 0 禁用 */
  state: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 职位模型
 */
export interface Pos extends BaseEntity {
  /** 父id */
  parentID: string;
  /** 子系统编码（比如：zuul） */
  svcCode: string;
  /** 组织专属岗位 */
  unitID: string;
  /** 组织类型专属岗位 */
  kindCode: string;
  /** 职位名称 */
  posName: string;
  /** 职位编码 */
  posCode: string;
  /** 是否可用:1 可用 0 禁用 */
  state: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人 */
  createBy: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 用户角色关联模型
 */
export interface UserRole {
  /** 主键 */
  id: string;
  /** 租户ID */
  tenantID: string;
  /** 用户id */
  userID: string;
  /** 角色id */
  roleID: string;
  /** 创建人 */
  createBy: string;
  /** 创建时间 */
  createAt: string;
  /** 修改人 */
  updateBy: string;
  /** 修改时间 */
  updateAt: string;
}

/**
 * 用户职位关联模型
 */
export interface UserPos {
  /** 主键 */
  id: string;
  /** 租户ID */
  tenantID: string;
  /** 用户id */
  userID: string;
  /** 职位id */
  posID: string;
  /** 创建人 */
  createBy: string;
  /** 创建时间 */
  createAt: string;
  /** 修改人 */
  updateBy: string;
  /** 修改时间 */
  updateAt: string;
}

/**
 * 角色菜单关联模型
 */
export interface RoleMenu {
  /** 主键 */
  id: string;
  /** 租户ID */
  tenantID: string;
  /** 角色id */
  roleID: string;
  /** 菜单id */
  menuID: string;
  /** 创建人 */
  createBy: string;
  /** 创建时间 */
  createAt: string;
  /** 修改人 */
  updateBy: string;
  /** 修改时间 */
  updateAt: string;
}

/**
 * 职位角色关联模型
 */
export interface PosRole {
  /** 主键 */
  id: string;
  /** 租户ID */
  tenantID: string;
  /** 职位id */
  posID: string;
  /** 角色id */
  roleID: string;
  /** 创建人 */
  createBy: string;
  /** 创建时间 */
  createAt: string;
  /** 修改人 */
  updateBy: string;
  /** 修改时间 */
  updateAt: string;
}

/**
 * 登录请求
 */
export interface LoginRequest {
  /** 帐号 */
  account: string;
  /** 密码 */
  password: string;
}

/**
 * 登录响应
 */
export interface LoginResponse {
  /** 令牌 */
  token: string;
  /** 用户信息 */
  user: User;
}

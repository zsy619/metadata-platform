/**
 * 基础实体类型
 */
export interface BaseEntity {
    /** 主键ID */
    id: string;
    /** 租户ID */
    tenantID: string;
    /** 创建时间 */
    createdAt: string;
    /** 更新时间 */
    updatedAt: string;
}

/**
 * 模型基础信息
 */
export interface Model extends BaseEntity {
    /** 父级ID */
    parentID: string;
    /** 数据连接ID */
    connID: string;
    /** 数据连接名称 */
    connName: string;
    /** 模型名称 */
    modelName: string;
    /** 模型编码 */
    modelCode: string;
    /** 模型版本 */
    modelVersion: string;
    /** 模型图片 */
    modelLogo: string;
    /** 模型类型：1sql语句、2视图/表、3存储过程、4关联 */
    modelKind: number;
    /** 是否公开 */
    isPublic: boolean;
    /** 是否锁定 */
    isLocked: boolean;
    /** 删除标识 */
    isDeleted: boolean;
    /** 状态: 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;

    /** 是否为树形结构 */
    isTree: boolean;
    /** 树形父级字段名 */
    treeParentField: string;
    /** 树形路径字段名 */
    treePathField: string;
    /** 树形层级字段名 */
    treeLevelField: string;
}
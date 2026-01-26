/**
 * API接口类型定义
 */

/**
 * 接口基础信息
 */
export interface API {
  /** 接口ID */
  apiID: number;
  /** 接口名称 */
  apiName: string;
  /** 接口编码 */
  apiCode: string;
  /** 接口路径 */
  apiPath: string;
  /** 请求方法 */
  apiMethod: string;
  /** 模型ID */
  modelID: number;
  /** 模型名称 */
  modelName: string;
  /** 接口类型：1-查询，2-新增，3-更新，4-删除，5-自定义 */
  apiType: number;
  /** 接口状态：1-启用，0-禁用 */
  state: number;
  /** 是否需要鉴权：1-需要，0-不需要 */
  needAuth: boolean;
  /** 是否需要审计：1-需要，0-不需要 */
  needAudit: boolean;
  /** 接口描述 */
  remark: string;
  /** 创建人ID */
  createID: number;
  /** 创建人 */
  createBy: string;
  /** 创建时间 */
  createAt: string;
  /** 更新人ID */
  updateID: number;
  /** 更新人 */
  updateBy: string;
  /** 更新时间 */
  updateAt: string;
}

/**
 * 接口查询参数
 */
export interface APIQueryParams {
  /** 页码 */
  page?: number;
  /** 每页条数 */
  pageSize?: number;
  /** 搜索关键词 */
  search?: string;
  /** 接口类型 */
  apiType?: number;
  /** 模型ID */
  modelID?: number;
  /** 状态 */
  state?: number;
}

/**
 * 接口响应结果
 */
export interface APIResponse {
  /** 数据列表 */
  data: API[];
  /** 总数 */
  total: number;
  /** 页码 */
  page: number;
  /** 每页条数 */
  pageSize: number;
}

/**
 * 接口请求参数配置
 */
export interface APIRequestParam {
  /** 参数ID */
  paramID: number;
  /** 接口ID */
  apiID: number;
  /** 参数名称 */
  paramName: string;
  /** 参数类型：1-路径参数，2-查询参数，3-请求体参数 */
  paramType: number;
  /** 数据类型：string, number, boolean, array, object */
  dataType: string;
  /** 是否必填：1-必填，0-可选 */
  isRequired: boolean;
  /** 默认值 */
  defaultValue: string;
  /** 参数描述 */
  remark: string;
  /** 排序 */
  sort: number;
}

/**
 * 接口响应数据配置
 */
export interface APIResponseData {
  /** 响应ID */
  responseID: number;
  /** 接口ID */
  apiID: number;
  /** 响应字段名称 */
  fieldName: string;
  /** 数据类型 */
  dataType: string;
  /** 是否必填 */
  isRequired: boolean;
  /** 示例值 */
  exampleValue: string;
  /** 字段描述 */
  remark: string;
  /** 排序 */
  sort: number;
}

/**
 * 接口测试结果
 */
export interface APITestResult {
  /** 是否成功 */
  success: boolean;
  /** 状态码 */
  statusCode: number;
  /** 响应时间（ms） */
  responseTime: number;
  /** 响应数据 */
  responseData: any;
  /** 错误信息 */
  errorMessage?: string;
}

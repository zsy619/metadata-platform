/**
 * 模型响应结果
 */
import { Model } from './model-base';

export interface ModelResponse {
  /** 数据列表 */
  data: Model[];
  /** 总数 */
  total: number;
  /** 页码 */
  page: number;
  /** 每页条数 */
  pageSize: number;
}
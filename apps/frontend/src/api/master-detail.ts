/**
 * 主子表事务API服务
 */
import request from '@/utils/request'

/**
 * 创建主子表事务
 * @param masterModelName 主表模型名
 * @param detailModelName 子表模型名
 * @param data 事务提交数据 { master: {...}, details: [...] }
 * @returns 创建结果
 */
export const createMasterDetail = async (masterModelName: string, detailModelName: string, data: any): Promise<any> => {
    return request({
        url: `/api/metadata/master-detail/${masterModelName}/${detailModelName}`,
        method: 'post',
        data
    })
}

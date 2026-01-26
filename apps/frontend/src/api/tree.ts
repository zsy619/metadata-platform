/**
 * 树形结构API服务
 */
import request from '@/utils/request'

/**
 * 获取树形结构数据
 * @param modelName 模型名称
 * @param params 查询参数
 * @returns 树形数据
 */
export const getTreeData = async (modelName: string, params?: any): Promise<any[]> => {
    return request({
        url: `/api/metadata/tree/${modelName}`,
        method: 'get',
        params
    })
}

/**
 * 添加树节点
 * @param modelName 模型名称
 * @param data 节点数据
 * @returns 添加结果
 */
export const addTreeNode = async (modelName: string, data: any): Promise<any> => {
    return request({
        url: `/api/metadata/tree/${modelName}/node`,
        method: 'post',
        data
    })
}

/**
 * 移动树节点
 * @param modelName 模型名称
 * @param id 节点ID
 * @param newParentID 新父节点ID
 * @returns 移动结果
 */
export const moveTreeNode = async (modelName: string, id: string, newParentID: string): Promise<any> => {
    return request({
        url: `/api/metadata/tree/${modelName}/node/${id}/move`,
        method: 'put',
        data: { parent_id: newParentID }
    })
}

/**
 * 删除树节点
 * @param modelName 模型名称
 * @param id 节点ID
 * @returns 删除结果
 */
export const deleteTreeNode = async (modelName: string, id: string): Promise<any> => {
    return request({
        url: `/api/metadata/tree/${modelName}/node/${id}`,
        method: 'delete'
    })
}

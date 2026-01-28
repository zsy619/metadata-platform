/**
 * 数据连接详情 (对应后端 MdConn 模型)
 */
export interface MdConn {
    id: string;
    tenant_id: string;
    parent_id: string;
    conn_name: string;
    conn_kind: string;
    conn_version: string;
    conn_host: string;
    conn_port: number;
    conn_user: string;
    conn_password: string;
    conn_database: string;
    conn_conn: string;
    state: number;
    remark: string;
    is_deleted: boolean;
    create_id: string;
    create_by: string;
    create_at: string;
    update_id: string;
    update_by: string;
    update_at: string;
}

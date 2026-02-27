export const DATA_RANGE = {
  ALL: '1',
  CUSTOM: '2',
  DEPARTMENT: '3',
  DEPT_AND_BELOW: '4'
} as const

export const DATA_RANGE_OPTIONS = [
  { value: DATA_RANGE.ALL, label: '全部数据权限' },
  { value: DATA_RANGE.CUSTOM, label: '自定数据权限' },
  { value: DATA_RANGE.DEPARTMENT, label: '本部门数据权限' },
  { value: DATA_RANGE.DEPT_AND_BELOW, label: '本部门及以下' }
]

export const DATA_RANGE_LABELS: Record<string, string> = {
  [DATA_RANGE.ALL]: '全部',
  [DATA_RANGE.CUSTOM]: '自定义',
  [DATA_RANGE.DEPARTMENT]: '本部门',
  [DATA_RANGE.DEPT_AND_BELOW]: '本部门及以下'
}

export const STATUS = {
  ENABLED: 1,
  DISABLED: 0
} as const

export const USER_KIND = {
  ADMIN: 1,
  USER: 2
} as const

export const MENU_TYPE = {
  DIRECTORY: 'M',
  MENU: 'C',
  BUTTON: 'F'
} as const

export const MENU_TYPE_OPTIONS = [
  { value: MENU_TYPE.DIRECTORY, label: '目录' },
  { value: MENU_TYPE.MENU, label: '菜单' },
  { value: MENU_TYPE.BUTTON, label: '按钮' }
]

/**
 * 数据库类型常量定义（与后端保持一致）
 */
export const DB_TYPE = {
  MYSQL: 'MySQL',
  POSTGRESQL: 'PostgreSQL',
  SQL_SERVER: 'SQL Server',
  ORACLE: 'Oracle',
  SQLITE: 'SQLite',
  CLICKHOUSE: 'ClickHouse',
  DM: 'DM',
  MONGODB: 'MongoDB',
  REDIS: 'Redis',
  TIDB: 'TiDB',
  OCEANBASE: 'OceanBase',
  DORIS: 'Doris',
  STARROCKS: 'StarRocks',
  OPENGauss: 'OpenGauss',
  KINGBASE: 'Kingbase'
} as const

/**
 * 数据库类型别名映射（用于标准化输入的数据库类型）
 */
const DB_TYPE_ALIASES: Record<string, string> = {
  'mysql': DB_TYPE.MYSQL,
  'mariadb': DB_TYPE.MYSQL,
  'postgres': DB_TYPE.POSTGRESQL,
  'postgresql': DB_TYPE.POSTGRESQL,
  'pg': DB_TYPE.POSTGRESQL,
  'sqlserver': DB_TYPE.SQL_SERVER,
  'sql server': DB_TYPE.SQL_SERVER,
  'mssql': DB_TYPE.SQL_SERVER,
  'oracle': DB_TYPE.ORACLE,
  'sqlite': DB_TYPE.SQLITE,
  'sqlite3': DB_TYPE.SQLITE,
  'clickhouse': DB_TYPE.CLICKHOUSE,
  'dm': DB_TYPE.DM,
  'dameng': DB_TYPE.DM,
  'mongodb': DB_TYPE.MONGODB,
  'mongo': DB_TYPE.MONGODB,
  'redis': DB_TYPE.REDIS,
  'tidb': DB_TYPE.TIDB,
  'oceanbase': DB_TYPE.OCEANBASE,
  'doris': DB_TYPE.DORIS,
  'starrocks': DB_TYPE.STARROCKS,
  'opengauss': DB_TYPE.OPENGauss,
  'kingbase': DB_TYPE.KINGBASE
}

/**
 * 标准化数据库类型（不区分大小写）
 * @param dbType 输入的数据库类型
 * @returns 标准化后的数据库类型
 */
export function normalizeDBType(dbType: string): string {
  const normalized = dbType?.trim().toLowerCase() || ''
  return DB_TYPE_ALIASES[normalized] || dbType
}

/**
 * 判断是否为 MySQL 兼容数据库
 * @param dbType 数据库类型
 * @returns 是否为 MySQL 兼容数据库
 */
export function isMySQL(dbType: string): boolean {
  const normalized = normalizeDBType(dbType)
  return [
    DB_TYPE.MYSQL,
    DB_TYPE.TIDB,
    DB_TYPE.OCEANBASE,
    DB_TYPE.DORIS,
    DB_TYPE.STARROCKS
  ].includes(normalized as any)
}

/**
 * 判断是否为 PostgreSQL 兼容数据库
 * @param dbType 数据库类型
 * @returns 是否为 PostgreSQL 兼容数据库
 */
export function isPostgreSQL(dbType: string): boolean {
  const normalized = normalizeDBType(dbType)
  return [
    DB_TYPE.POSTGRESQL,
    DB_TYPE.OPENGauss,
    DB_TYPE.KINGBASE
  ].includes(normalized as any)
}

/**
 * 获取数据库类型的显示名称列表（用于下拉选择等）
 */
export const DB_TYPE_OPTIONS = [
  { value: DB_TYPE.MYSQL, label: 'MySQL' },
  { value: DB_TYPE.POSTGRESQL, label: 'PostgreSQL' },
  { value: DB_TYPE.SQL_SERVER, label: 'SQL Server' },
  { value: DB_TYPE.ORACLE, label: 'Oracle' },
  { value: DB_TYPE.SQLITE, label: 'SQLite' },
  { value: DB_TYPE.CLICKHOUSE, label: 'ClickHouse' },
  { value: DB_TYPE.DM, label: '达梦 (DM)' },
  { value: DB_TYPE.MONGODB, label: 'MongoDB' },
  { value: DB_TYPE.REDIS, label: 'Redis' },
  { value: DB_TYPE.TIDB, label: 'TiDB' },
  { value: DB_TYPE.OCEANBASE, label: 'OceanBase' },
  { value: DB_TYPE.DORIS, label: 'Doris' },
  { value: DB_TYPE.STARROCKS, label: 'StarRocks' },
  { value: DB_TYPE.OPENGauss, label: 'OpenGauss' },
  { value: DB_TYPE.KINGBASE, label: '人大金仓' }
]

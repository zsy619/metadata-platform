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

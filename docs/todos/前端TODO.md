# 前端开发TODO

> **项目**: 元数据管理平台  
> **技术栈**: Vue3 + TypeScript + Element Plus + Vite  
> **最后更新**: 2026-01-24

---

## 阶段1：基础框架 (预计4周)

### 项目初始化
- [ ] 创建Vue3+TypeScript项目（Vite）
- [ ] 安装依赖（vue-router@4 + pinia + element-plus + axios + sass）
- [ ] 配置tsconfig.json
- [ ] 配置vite.config.ts（路径别名@+代理+打包优化）
- [ ] 配置环境变量（.env.development/.env.production）
- [ ] 配置ESLint和Prettier
- [ ] 创建目录结构（api/assets/components/layouts/router/stores/types/utils/views）

### 工具类封装
- [ ] Axios封装（request.ts）
  - [ ] 创建axios实例（baseURL+timeout）
  - [ ] 请求拦截器（添加Token+Content-Type+Request-ID）
  - [ ] 响应拦截器（错误处理+Token过期自动刷新+响应解包）
  - [ ] 请求重试机制
  - [ ] 请求取消机制
  - [ ] 错误统一处理（网络错误/业务错误/系统错误）
- [ ] 本地存储封装（storage.ts）
  - [ ] Storage类（支持localStorage/sessionStorage）
  - [ ] 过期时间支持
  - [ ] 自动JSON序列化
- [ ] 工具函数库
  - [ ] 日期格式化（format.ts）
  - [ ] 数字格式化（千分位/百分比/金额）
  - [ ] 表单验证函数（validate.ts：邮箱/手机号/身份证/URL）
  - [ ] 文件处理（FileUtils：大小格式化/类型判断）

### 路由配置
- [ ] 路由定义（routes.ts）
  - [ ] 登录页路由
  - [ ] 首页仪表盘
  - [ ] 数据源管理路由
  - [ ] 模型管理路由
  - [ ] 接口管理路由
  - [ ] 404路由
- [ ] 路由守卫（guards.ts）
  - [ ] 前置守卫（登录验证+权限验证+用户信息加载）
  - [ ] 后置守卫（NProgress+页面标题）
  - [ ] 错误处理
- [ ] 路由懒加载配置
- [ ] 路由缓存配置（keep-alive）

### 状态管理（Pinia）
- [ ] useUserStore（user.ts）
  - [ ] 状态：token/refreshToken/userInfo
  - [ ] Getters：isLoggedIn/userName/userRoles
  - [ ] Actions：login/logout/getUserInfo/refreshToken
  - [ ] 状态持久化（localStorage）
- [ ] useAppStore（app.ts）
  - [ ] 状态：菜单折叠/主题配置/面包屑
  - [ ] Actions：toggleSidebar/setTheme/setBreadcrumb
- [ ] usePermissionStore（permission.ts，可选）
  - [ ] 动态路由管理
  - [ ] 权限验证

### 布局组件
- [ ] DefaultLayout（layouts/DefaultLayout.vue）
  - [ ] 顶部导航栏（用户信息+退出+主题切换）
  - [ ] 侧边菜单（多级菜单+折叠）
  - [ ] 主内容区（router-view+keep-alive）
  - [ ] 面包屑导航
  - [ ] 响应式适配
- [ ] Sidebar组件（layouts/components/Sidebar.vue）
  - [ ] 菜单渲染（递归组件）
  - [ ] 菜单高亮
  - [ ] 菜单展开/折叠
  - [ ] 图标显示
- [ ] AppHeader组件（layouts/components/AppHeader.vue）
  - [ ] Logo+系统名称
  - [ ] 用户下拉菜单
  - [ ] 全屏按钮
  - [ ] 主题切换
- [ ] Breadcrumb组件（layouts/components/Breadcrumb.vue）
  - [ ] 自动生成面包屑
  - [ ] 点击跳转
- [ ] AuthLayout（登录页布局）

### 公共组件
- [ ] DataTable组件（components/table/DataTable.vue）
  - [ ] props配置化（columns/data/loading/pagination）
  - [ ] 分页支持
  - [ ] 排序支持
  - [ ] 筛选支持
  - [ ] 操作列插槽
  - [ ] 空数据提示
  - [ ] 加载状态
  - [ ] 刷新功能
- [ ] DynamicForm组件（components/form/DynamicForm.vue）
  - [ ] JSON配置驱动
  - [ ] 支持各种输入类型（input/select/date-picker/upload等）
  - [ ] 表单验证
  - [ ] 重置功能
  - [ ] 只读/禁用模式
- [ ] BaseDialog组件（components/dialog/BaseDialog.vue）
  - [ ] 统一样式
  - [ ] 多种尺寸
  - [ ] 确认/取消按钮
  - [ ] Loading状态
- [ ] PageContainer组件（components/common/PageContainer.vue）
  - [ ] 页面头部（标题+操作按钮）
  - [ ] 内容区
  - [ ] 统一间距
- [ ] EmptyState组件（components/common/EmptyState.vue）
  - [ ] 图标+文字
  - [ ] 操作按钮
- [ ] Loading组件（全局Loading/局部Loading/骨架屏）

### API接口定义
- [ ] Auth API（api/auth/index.ts）
  - [ ] login/logout/getProfile/refresh/changePassword
- [ ] DataSource API（api/datasource/index.ts）
  - [ ] getList/create/update/delete/testConnection
  - [ ] getTables/getViews/getTableStructure/previewData
- [ ] Model API（api/model/index.ts）
  - [ ] getList/create/update/delete
  - [ ] buildFromTable/buildFromView
  - [ ] getFields/addField/updateField/deleteField
- [ ] API Management API（api/api/index.ts）
  - [ ] getList/batchGenerate/update/delete/enable
- [ ] QueryTemplate API（api/query-template/index.ts）
  - [ ] getList/create/update/delete/setDefault
- [ ] TypeScript类型定义（types/*.ts）
  - [ ] User/UserInfo/LoginParams
  - [ ] DataSource/DataSourceForm
  - [ ] Model/ModelField/BuildModelParams
  - [ ] API/APIParam
  - [ ] QueryTemplate/QueryCondition

### 基础页面
- [ ] Login页面（views/auth/Login.vue）
  - [ ] 登录表单（用户名+密码）
  - [ ] 表单验证
  - [ ] 记住密码
  - [ ] 调用登录API
  - [ ] 保存Token
  - [ ] 跳转首页
  - [ ] 错误提示
  - [ ] Loading状态
- [ ] Dashboard页面（views/home/Dashboard.vue）
  - [ ] 统计卡片（数据源数/模型数/接口数/用户数）
  - [ ] 图表展示（可选）
  - [ ] 快捷操作
- [ ] 404页面（views/error/404.vue）
  - [ ] 友好提示
  - [ ] 返回首页按钮

---

## 阶段2：核心功能 (预计5周)

### 数据源管理
- [ ] 数据源列表页（views/datasource/List.vue）
  - [ ] 列表展示（表格）
  - [ ] 分页
  - [ ] 搜索（关键字）
  - [ ] 筛选（类型）
  - [ ] 操作按钮（新建/编辑/删除/测试连接）
  - [ ] 连接状态显示
  - [ ] 批量操作
- [ ] 数据源表单对话框（components/DataSourceForm.vue）
  - [ ] 表单字段（名称/类型/主机/端口/用户名/密码/数据库）
  - [ ] 表单验证（必填/格式）
  - [ ] 连接测试按钮
  - [ ] 测试成功自动填充
  - [ ] 密码显示/隐藏
  - [ ] 保存/取消
- [ ] 数据库对象浏览器（components/ObjectBrowser.vue）
  - [ ] 树形结构（数据源/Schema/表/视图）
  - [ ] 懒加载
  - [ ] 搜索功能
  - [ ] 刷新功能
  - [ ] 右键菜单（预览/构建模型）
- [ ] 数据预览对话框（components/DataPreview.vue）
  - [ ] 表格展示
  - [ ] 分页
  - [ ] 列信息显示

### 模型管理⭐⭐⭐
- [ ] 模型列表页（views/model/List.vue）
  - [ ] 列表展示
  - [ ] 分页+搜索
  - [ ] 筛选（按数据源/类型/状态）
  - [ ] 操作（新建/从表构建/编辑/删除/克隆）
  - [ ] 模型状态显示（锁定/公开）
  - [ ] 批量操作
- [ ] 模型构建向导（components/BuildWizard.vue）
  - [ ] 步骤器（6步）
  - [ ] 向导流程控制
  - [ ] 数据验证
  - [ ] 进度保存（可选）
- [ ] 步骤1：数据源选择器（components/DataSourceSelector.vue）
  - [ ] 数据源列表
  - [ ] 连接状态显示
  - [ ] 测试连接
- [ ] 步骤2：表/视图选择器（components/TableViewSelector.vue）
  - [ ] 标签页切换（表/视图）
  - [ ] 表列表展示
  - [ ] 搜索功能
  - [ ] 预览表结构按钮
  - [ ] 预览表数据按钮
  - [ ] 记录数显示
- [ ] 步骤3：字段配置器（components/FieldConfigurator.vue）
  - [ ] 字段列表表格
  - [ ] 全选/全不选
  - [ ] 字段选择（checkbox）
  - [ ] 字段属性编辑（显示名称/字段代码/数据类型/长度/必填/主键/默认值）
  - [ ] 添加自定义字段
  - [ ] 删除自定义字段
  - [ ] 字段拖拽排序
- [ ] 步骤4：验证规则配置器（components/ValidationConfigurator.vue）
  - [ ] 为每个字段配置验证规则
  - [ ] 验证类型选择（必填/长度/范围/格式/自定义）
  - [ ] 验证规则模板（邮箱/手机号/身份证/URL）
  - [ ] 错误提示信息配置
  - [ ] 验证规则测试
- [ ] 步骤5：显示配置器（components/DisplayConfigurator.vue）
  - [ ] 显示名称配置
  - [ ] 显示顺序调整（拖拽）
  - [ ] 显示宽度配置
  - [ ] 组件类型选择（input/select/textarea/date-picker/upload等）
  - [ ] 组件选项配置（下拉选项/上传限制等）
  - [ ] 占位符/帮助文本
  - [ ] 可搜索/可排序/可筛选配置
- [ ] 步骤6：配置预览（components/BuildPreview.vue）
  - [ ] 模型信息摘要
  - [ ] 字段列表预览
  - [ ] 验证规则预览
  - [ ] 显示配置预览
  - [ ] 确认创建按钮

### 接口管理
- [ ] 接口列表页（views/api/List.vue）
  - [ ] 列表展示
  - [ ] 筛选（按模型/方法/状态）
  - [ ] 操作（批量生成/新建/编辑/删除/测试）
  - [ ] 启用/禁用开关
  - [ ] 接口路径显示
  - [ ] 请求方法标签（GET/POST/PUT/DELETE）
- [ ] 批量生成对话框（components/BatchGenerateDialog.vue）
  - [ ] 模型选择（下拉+搜索）
  - [ ] 接口类型选择（Create/Read/Update/Delete/List）
  - [ ] 路径前缀配置
  - [ ] 接口预览列表
  - [ ] 冲突检测提示
  - [ ] 确认生成
- [ ] 接口编辑页面（views/api/Editor.vue）
  - [ ] 基本信息（名称/编码/路径/方法/关联模型/描述）
  - [ ] 路径参数配置
  - [ ] 查询参数配置
  - [ ] 请求体配置
  - [ ] 响应配置
  - [ ] 保存/取消
- [ ] 接口测试工具（components/APITester.vue）
  - [ ] 请求构造器（自动填充参数）
  - [ ] 参数输入（路径/查询/body）
  - [ ] 发送请求按钮
  - [ ] 请求信息显示（URL/Method/Headers/Body）
  - [ ] 响应展示（JSON格式化+状态码+响应时间）
  - [ ] 历史记录（保存+快速复用）
  - [ ] 错误提示

### 查询模板管理⭐⭐
- [ ] 查询模板列表（在模型详情页添加标签）
  - [ ] 模板列表展示
  - [ ] 默认模板标识
  - [ ] 操作（新建/编辑/删除/设为默认/复制/预览）
- [ ] 可视化条件构建器（components/QueryBuilder.vue）
  - [ ] 条件行组件
  - [ ] 字段选择下拉（显示字段类型）
  - [ ] 操作符选择（根据字段类型显示可用操作符）
  - [ ] 值输入组件
    - [ ] 文本输入
    - [ ] 数字输入
    - [ ] 日期选择器
    - [ ] 下拉选择
    - [ ] 多选输入（用于IN）
    - [ ] 范围输入（用于BETWEEN）
  - [ ] 逻辑运算符选择（AND/OR）
  - [ ] 添加条件按钮
  - [ ] 删除条件按钮
  - [ ] 条件组支持（括号）
  - [ ] 添加条件组
  - [ ] 嵌套条件支持
- [ ] SQL预览组件（components/SQLPreview.vue）
  - [ ] 显示生成的SQL
  - [ ] 语法高亮
  - [ ] 复制按钮
  - [ ] 格式化按钮
- [ ] 查询模板编辑器（views/query-template/Editor.vue）
  - [ ] 模板基本信息（名称/编码/描述）
  - [ ] 集成条件构建器
  - [ ] 条件配置
  - [ ] SQL预览
  - [ ] 结果预览（执行查询+表格展示）
  - [ ] 保存/测试/取消

### 数据查询界面
- [ ] 数据列表页（views/data/List.vue）
  - [ ] 动态表格（根据模型字段生成列）
  - [ ] 高级筛选器
    - [ ] 快速筛选（常用字段）
    - [ ] 自定义筛选（调用条件构建器）
    - [ ] 模板筛选（应用查询模板）
    - [ ] 筛选保存为模板
  - [ ] 字段选择器（勾选显示字段+调整顺序+固定列）
  - [ ] 分页
  - [ ] 排序
  - [ ] 操作列（查看/编辑/删除）
  - [ ] 批量操作
  - [ ] 导出按钮
- [ ] 数据详情页（views/data/Detail.vue）
  - [ ] 字段值展示
  - [ ] 格式化显示
  - [ ] 编辑按钮
  - [ ] 返回列表
- [ ] 数据表单对话框（components/DataForm.vue）
  - [ ] 动态表单（根据模型字段配置生成）
  - [ ] 组件类型渲染（input/select/date-picker等）
  - [ ] 表单验证（应用验证规则）
  - [ ] 默认值填充
  - [ ] 必填标识
  - [ ] 帮助文本显示
  - [ ] 保存/取消

---

## 阶段3：高级功能 (预计3周)

### 树形数据展示
- [ ] 树形表格组件（components/table/TreeTable.vue）
  - [ ] 树形结构渲染
  - [ ] 展开/折叠节点
  - [ ] 层级缩进显示
  - [ ] 连接线显示（可选）
  - [ ] 懒加载子节点
  - [ ] 全部展开/折叠
- [ ] 树形节点操作
  - [ ] 添加同级节点对话框
  - [ ] 添加子节点对话框
  - [ ] 编辑节点对话框
  - [ ] 删除节点确认（提示子节点处理方式）
  - [ ] 移动节点对话框或拖拽
- [ ] 树形筛选器组件（components/TreeFilter.vue）
  - [ ] 关键字搜索
  - [ ] 自动展开匹配节点
  - [ ] 高亮匹配文本
  - [ ] 只显示匹配节点及其路径
- [ ] 树形导出功能
  - [ ] 导出为层级JSON
  - [ ] 导出为扁平JSON
  - [ ] 导出为Excel（层级缩进）

### 主子表管理
- [ ] 主子表表单（components/form/MasterDetailForm.vue）
  - [ ] 主表表单区域
  - [ ] 子表表格（嵌入式可编辑表格）
  - [ ] 添加子表行按钮
  - [ ] 删除子表行按钮
  - [ ] 编辑子表行（行内编辑或对话框）
  - [ ] 子表验证
  - [ ] 至少一条子表验证
  - [ ] 保存整体数据
- [ ] 主子表详情（components/MasterDetailView.vue）
  - [ ] 主表信息展示
  - [ ] 子表列表展示（可分页）
  - [ ] 子表排序筛选
  - [ ] 编辑按钮
  - [ ] 删除按钮
  - [ ] 添加子表数据按钮

### 导入导出
- [ ] 导出对话框（components/ImportExport/ExportDialog.vue）
  - [ ] 格式选择（Excel/CSV/JSON）
  - [ ] 字段选择器（勾选要导出的字段）
  - [ ] 筛选条件配置（应用当前筛选或自定义）
  - [ ] 查询模板选择
  - [ ] 导出设置（编码/分隔符/样式等）
  - [ ] 进度显示（进度条）
  - [ ] 完成后自动下载
  - [ ] 错误提示
- [ ] 导入对话框（components/ImportExport/ImportDialog.vue）
  - [ ] 模板下载按钮
  - [ ] 文件上传区（拖拽或点击）
  - [ ] 文件类型验证
  - [ ] 数据预览（显示前10行）
  - [ ] 字段映射确认
  - [ ] 导入设置（是否忽略错误/更新已存在数据）
  - [ ] 导入执行
  - [ ] 进度显示
  - [ ] 成功/失败数量统计
  - [ ] 错误详情展示
  - [ ] 下载错误记录
- [ ] 文件上传组件（components/FileUpload.vue）
  - [ ] 拖拽上传
  - [ ] 点击选择
  - [ ] 文件类型验证
  - [ ] 文件大小验证
  - [ ] 上传进度
  - [ ] 上传成功/失败提示
  - [ ] 文件列表显示
  - [ ] 删除文件

### 审计日志查看
- [ ] 操作日志页面（views/logs/OperationLogs.vue）
  - [ ] 日志列表
  - [ ] 筛选（按用户/操作类型/模块/时间范围/状态）
  - [ ] 日志详情对话框（完整请求/响应信息+JSON格式化）
  - [ ] 导出日志
- [ ] 数据变更日志页面（views/logs/DataChangeLogs.vue）
  - [ ] 变更记录列表
  - [ ] 筛选（按表名/记录ID/用户/时间范围）
  - [ ] 变更对比展示（并排对比+高亮变更字段）
  - [ ] 版本回溯功能
  - [ ] 变更统计图表
- [ ] 访问日志页面（views/logs/AccessLogs.vue）
  - [ ] 访问记录列表
  - [ ] 筛选（按用户/路径/状态码/时间范围）
  - [ ] 可视化统计
    - [ ] 访问量趋势图
    - [ ] 响应时间统计
    - [ ] 错误率统计
    - [ ] TOP10路径
  - [ ] 异常访问高亮

---

## 阶段4：系统完善 (预计2周)

### 用户体验优化
- [ ] 加载状态优化
  - [ ] 全局Loading（路由切换+接口请求）
  - [ ] 骨架屏组件（列表页+详情页）
  - [ ] 进度条组件（文件上传+数据导入导出）
  - [ ] Loading文案提示
- [ ] 错误处理优化
  - [ ] 全局错误边界组件
  - [ ] 友好错误页面（500/403等）
  - [ ] 错误提示优化（分类+建议）
  - [ ] 错误上报（收集前端错误发送到后端）
- [ ] 操作反馈优化
  - [ ] 操作确认对话框（删除/批量操作/重要修改）
  - [ ] 二次确认（高危操作）
  - [ ] 成功提示（Toast）
  - [ ] 失败提示（详细错误信息）
  - [ ] 警告提示
- [ ] 空状态优化
  - [ ] 空数据提示（图标+文字+'引导操作'）
  - [ ] 无权限提示
  - [ ] 无搜索结果提示
  - [ ] 操作引导（新手引导/功能提示/帮助文档链接）

### 性能优化
- [ ] 代码分割
  - [ ] 路由懒加载（全部页面）
  - [ ] 组件懒加载（大组件+非首屏组件）
  - [ ] 第三方库按需引入（Element Plus自动导入）
  - [ ] Tree Shaking优化
- [ ] 资源优化
  - [ ] 图片压缩
  - [ ] 图片懒加载
  - [ ] 启用Gzip压缩
  - [ ] CDN加速（可选）
- [ ] 运行时优化
  - [ ] 虚拟列表（大数据量表格+长列表）
  - [ ] 防抖节流（搜索输入+窗口resize+滚动）
  - [ ] 组件缓存（Keep-Alive+缓存策略）
  - [ ] 计算属性优化
- [ ] 请求优化
  - [ ] 请求缓存（相同请求+缓存时间配置）
  - [ ] 请求合并（批量查询）
  - [ ] 请求取消（页面切换+重复请求）
  - [ ] 接口防抖
- [ ] 打包优化
  - [ ] Chunk分包策略
  - [ ] 资源预加载（preload/prefetch）
  - [ ] 构建缓存
  - [ ] 分析打包体积（rollup-plugin-visualizer）

### 主题与国际化
- [ ] 主题切换
  - [ ] 亮色主题
  - [ ] 暗色主题
  - [ ] 自定义主题色
  - [ ] 主题配置器（颜色选择+预览）
  - [ ] 主题持久化
  - [ ] 一键切换
  - [ ] CSS变量管理
- [ ] 国际化支持（可选）
  - [ ] 集成vue-i18n
  - [ ] 中文语言包
  - [ ] 英文语言包
  - [ ] 语言切换
  - [ ] 语言持久化
  - [ ] 日期/数字本地化
- [ ] 移动端适配（可选）
  - [ ] 响应式布局优化（媒体查询）
  - [ ] 侧边栏改为抽屉
  - [ ] 表格横向滚动
  - [ ] 表单单列布局
  - [ ] 触摸优化（滑动/长按/手势）
  - [ ] 移动端菜单（底部导航）

### 监控与分析
- [ ] 系统监控面板（views/monitor/Dashboard.vue）
  - [ ] 数据统计卡片
  - [ ] 性能指标展示（接口响应时间/QPS/错误率）
  - [ ] 趋势图表
  - [ ] 实时监控（WebSocket）
- [ ] 性能监控面板（views/monitor/Performance.vue）
  - [ ] 接口响应时间分布
  - [ ] 慢接口TOP10
  - [ ] 错误接口TOP10
  - [ ] 性能趋势图
- [ ] 用户行为分析（可选）
  - [ ] 页面访问统计
  - [ ] 功能使用频率
  - [ ] 用户路径分析

### 文档与帮助
- [ ] 帮助文档集成
  - [ ] 在线帮助中心（内嵌或链接）
  - [ ] 功能说明
  - [ ] 操作指南
  - [ ] 常见问题
  - [ ] 视频教程（可选）
- [ ] 组件文档（开发用）
  - [ ] 公共组件使用文档
  - [ ] 示例代码
- [ ] 更新日志展示
  - [ ] 版本更新记录
  - [ ] 新功能说明

---

## 验收标准

### 阶段1验收
- 登录功能正常
- 数据源页面可用
- 模型列表可展示
- 布局美观，响应式正常

### 阶段2验收
- 模型构建向导流程完整
- 接口管理功能可用
- 查询模板编辑器易用
- 数据查询界面功能完整

### 阶段3验收
- 树形数据展示和操作正常
- 主子表管理功能完整
- 导入导出功能可用
- 日志查看界面清晰

### 阶段4验收
- 首屏加载时间<3s
- 大数据量不卡顿
- 主题切换流畅
- 用户体验优秀

---

**预计总工期**: 14周（3.5个月）  
**关键依赖**: 后端接口完成后才能进行前后端联调

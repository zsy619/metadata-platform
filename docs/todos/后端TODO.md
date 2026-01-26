# 后端开发TODO

> **项目**: 元数据管理平台  
> **技术栈**: Golang + Hertz + GORM + MySQL/PostgreSQL  
> **最后更新**: 2026-01-24

---

## 阶段1：基础框架 (预计4周)

### 项目初始化
- [ ] 创建项目目录结构（Clean Architecture）
- [ ] 配置管理模块（Viper，支持YAML/JSON/环境变量）
- [ ] 日志系统（Zap + Lumberjack，多级别+轮转）
- [ ] 数据库连接池（MySQL/PostgreSQL，健康检查）
- [ ] 雪花算法ID生成器（并发安全，性能>100k/s）
- [ ] 数据库迁移脚本（用户管理+元数据模块11张表）

### Web框架
- [ ] Hertz服务器初始化（优雅关闭）
- [ ] 统一响应封装（Response/PageResponse/TraceID）
- [ ] 全局异常处理中间件（panic捕获+日志）
- [ ] CORS中间件（跨域配置）
- [ ] 请求日志中间件（记录请求/响应）

### 数据源管理模块（对应md_conn/md_table/md_table_field）
- [ ] DataSource Model定义（密码加密BeforeCreate/BeforeUpdate钩子）
- [ ] DataSource Repository（CRUD+分页+搜索+软删除）
- [ ] DataSource Service（连接测试+密码加密处理）
- [ ] 数据库驱动适配器
  - [ ] 定义MetadataExtractor接口
  - [ ] MySQLExtractor实现（GetTables/GetViews/GetColumns/GetTableStructure）
  - [ ] PostgreSQLExtractor实现
  - [ ] 连接池管理器（动态创建+缓存+关闭）
- [ ] DataSource API Handler
  - [ ] POST /api/data-sources（创建）
  - [ ] GET /api/data-sources（列表+分页+搜索）
  - [ ] GET /api/data-sources/{id}（详情，密码脱敏）
  - [ ] PUT /api/data-sources/{id}（更新）
  - [ ] DELETE /api/data-sources/{id}（软删除+检查依赖）
  - [ ] POST /api/data-sources/{id}/test（连接测试）
  - [ ] GET /api/data-sources/{id}/tables（获取表列表）
  - [ ] GET /api/data-sources/{id}/views（获取视图列表）
  - [ ] GET /api/data-sources/{id}/tables/{table}/structure（表结构）
  - [ ] GET /api/data-sources/{id}/tables/{table}/preview（数据预览）

### 模型管理模块（对应md_model/md_model_field/md_model_table等）
- [ ] Model定义（md_model，支持model_kind：1SQL/2视图表/3存储过程/4关联）
- [ ] ModelField定义（md_model_field，字段配置+函数+聚合）
- [ ] ModelTable定义（md_model_table，主表标记）
- [ ] ModelJoin定义（md_model_join，JOIN配置+树形结构）
- [ ] ModelWhere定义（md_model_where，WHERE条件+括号+逻辑运算）
- [ ] ModelGroup定义（md_model_group，GROUP BY）
- [ ] ModelHaving定义（md_model_having，HAVING条件）
- [ ] ModelOrder定义（md_model_order，ORDER BY+ASC/DESC）
- [ ] ModelLimit定义（md_model_limit，分页）
- [ ] ModelSQL定义（md_model_sql，原始SQL）
- [ ] Model Repository（CRUD+级联查询+版本管理）
- [ ] Model Service
  - [ ] Create/Get/Update/Delete/List方法
  - [ ] BuildFromTable方法（从表构建模型+字段映射）
  - [ ] BuildFromView方法（从视图构建+只读标记）
  - [ ] 字段数据类型映射（MySQL/PostgreSQL类型映射表）
  - [ ] 验证规则自动生成（基于NOT NULL/长度/类型）
  - [ ] 默认值处理（固定值+动态值CURRENT_TIMESTAMP）
  - [ ] 模型配置验证（唯一性+主键+类型）
- [ ] Model API Handler
  - [ ] POST /api/models（创建）
  - [ ] GET /api/models（列表+分页+按数据源筛选）
  - [ ] GET /api/models/{id}（详情+包含字段）
  - [ ] PUT /api/models/{id}（更新）
  - [ ] DELETE /api/models/{id}（软删除+检查依赖）
  - [ ] POST /api/models/build-from-table（从表构建）
  - [ ] POST /api/models/build-from-view（从视图构建）
  - [ ] GET /api/models/{id}/fields（字段列表）
  - [ ] POST /api/models/{id}/fields（添加字段）
  - [ ] PUT /api/models/{id}/fields/{fieldId}（更新字段）
  - [ ] DELETE /api/models/{id}/fields/{fieldId}（删除字段）
  - [ ] POST /api/models/{id}/fields/batch-update（批量更新）
  - [ ] POST /api/models/{id}/fields/reorder（字段排序）

### 认证与权限
- [ ] JWT工具封装（GenerateToken/ParseToken/RefreshToken）
- [ ] Auth Service（Login/Logout/GetUserInfo/ChangePassword）
- [ ] Auth Middleware（Token验证+用户信息提取+过期处理）
- [ ] Auth API Handler
  - [ ] POST /api/auth/login
  - [ ] POST /api/auth/logout
  - [ ] GET /api/auth/profile
  - [ ] POST /api/auth/refresh
  - [ ] POST /api/auth/change-password
- [ ] Casbin权限集成（RBAC模型+权限验证中间件+基础策略）

---

## 阶段2：核心业务 (预计5周)

### SQL生成引擎⭐⭐⭐
- [ ] SQLBuilder架构
  - [ ] 创建SQLBuilder主类
  - [ ] 定义ModelData结构体（聚合所有模型数据）
  - [ ] loadModelData方法（加载模型+表+字段+JOIN+WHERE+GROUP等）
- [ ] SELECT子句构建器
  - [ ] buildSelectClause方法
  - [ ] buildFieldExpression方法（字段+表名前缀）
  - [ ] 支持字段函数（UPPER/LOWER/DATE_FORMAT/SUBSTRING等）
  - [ ] 支持聚合函数（SUM/COUNT/AVG/MAX/MIN）
  - [ ] 支持字段别名（show_title）
- [ ] FROM子句构建器
  - [ ] buildFromClause方法
  - [ ] 识别主表（is_main标记）
  - [ ] 支持schema.table格式
- [ ] JOIN子句构建器
  - [ ] buildJoinClause方法
  - [ ] buildJoinTree方法（支持嵌套JOIN，parent_id树形结构）
  - [ ] generateJoinSQL方法（递归生成）
  - [ ] buildJoinConditions方法（ON条件+括号+AND/OR）
  - [ ] 支持LEFT/RIGHT/INNER/FULL OUTER JOIN
  - [ ] 支持字段与字段比较
  - [ ] 支持字段与值比较
  - [ ] 支持JOIN字段函数
- [ ] WHERE子句构建器
  - [ ] buildWhereClause方法
  - [ ] buildSingleCondition方法
  - [ ] 支持所有运算符（=/!=/>/</>=/<=/LIKE/IN/BETWEEN/IS NULL/IS NOT NULL）
  - [ ] 支持复杂条件（括号+AND/OR组合）
  - [ ] 支持字段与字段比较
  - [ ] 支持字段与值比较
  - [ ] 参数化查询（防SQL注入）
  - [ ] 支持字段函数
  - [ ] 支持动态参数替换
- [ ] GROUP BY和HAVING子句构建器
  - [ ] buildGroupByClause方法
  - [ ] buildHavingClause方法
  - [ ] 支持分组字段函数
  - [ ] 支持HAVING聚合条件（类似WHERE逻辑）
- [ ] ORDER BY和LIMIT子句构建器
  - [ ] buildOrderByClause方法（多字段+ASC/DESC）
  - [ ] buildLimitClause方法（LIMIT+OFFSET）
  - [ ] 支持动态分页参数
- [ ] SQL组装与执行
  - [ ] buildFromMetadata方法（组装完整SQL）
  - [ ] buildFromSQL方法（处理原始SQL类型）
  - [ ] SQL语法验证
  - [ ] SQL注入检测

### SQL执行器
- [ ] SQLExecutor类
  - [ ] Execute方法（执行模型查询+动态参数）
  - [ ] ExecuteCount方法（统计总数）
  - [ ] parseRows方法（结果解析为map切片）
  - [ ] SQL日志记录
  - [ ] 执行时间监控
  - [ ] 慢查询告警
- [ ] 数据源连接管理
  - [ ] GetConnection方法（根据conn_id获取连接）
  - [ ] 连接池复用
  - [ ] 连接健康检查

### 字段增强配置
- [ ] FieldEnhancement Model（需新建表md_model_field_enhancements）
  - [ ] display_name（显示名称）
  - [ ] display_order（显示顺序）
  - [ ] display_width（显示宽度）
  - [ ] is_searchable/is_sortable/is_filterable（查询配置）
  - [ ] placeholder/help_text（提示信息）
  - [ ] component_type（组件类型：input/select/date-picker等）
  - [ ] component_config（JSON组件配置：选项/限制等）
- [ ] FieldEnhancement Service
  - [ ] UpdateEnhancements方法（单个更新）
  - [ ] BatchUpdateEnhancements方法（批量更新）
  - [ ] GetEnhancements方法（按模型ID查询）
- [ ] FieldEnhancement API
  - [ ] GET /api/models/{id}/fields/enhancements
  - [ ] PUT /api/models/{id}/fields/enhancements
  - [ ] POST /api/models/{id}/fields/batch-enhancements

### CRUD接口自动生成
- [ ] API Model（基于现有表）
- [ ] APIParam Model
- [ ] API Generator Service
  - [ ] BatchGenerate方法（批量生成CRUD接口配置）
  - [ ] 生成Create/Read/Update/Delete/List接口
  - [ ] 生成接口路径（/api/data/{model_code}）
  - [ ] 生成参数配置
  - [ ] 路径冲突检测
- [ ] Dynamic Router
  - [ ] RegisterAPI方法（动态注册路由）
  - [ ] handleCreate/handleGet/handleUpdate/handleDelete/handleList
  - [ ] 路由参数绑定
  - [ ] 请求验证
- [ ] CRUD Service
  - [ ] Create方法（验证+默认值+插入+返回ID）
  - [ ] Get方法（按ID查询）
  - [ ] Update方法（验证+更新）
  - [ ] Delete方法（软删除）
  - [ ] List方法（分页+排序+筛选+应用查询模板）
  - [ ] BatchCreate/BatchUpdate/BatchDelete
  - [ ] 数据验证器（基于字段配置）
  - [ ] 默认值应用
- [ ] API Management API
  - [ ] GET /api/apis（接口列表）
  - [ ] POST /api/apis/batch-generate（批量生成）
  - [ ] PUT /api/apis/{id}（更新配置）
  - [ ] DELETE /api/apis/{id}（删除）
  - [ ] POST /api/apis/{id}/enable（启用/禁用）
  - [ ] POST /api/apis/{id}/test（测试接口）

### 查询模板管理
- [ ] QueryTemplate Model（基于md_model相关表）
- [ ] QueryCondition Model
- [ ] QueryTemplate Service
  - [ ] Create/Get/Update/Delete方法
  - [ ] SetDefault方法（设置默认模板）
  - [ ] ApplyTemplate方法（应用模板到查询）
  - [ ] Duplicate方法（复制模板）
- [ ] QueryTemplate API
  - [ ] GET /api/models/{id}/query-templates
  - [ ] POST /api/models/{id}/query-templates
  - [ ] GET /api/query-templates/{id}
  - [ ] PUT /api/query-templates/{id}
  - [ ] DELETE /api/query-templates/{id}
  - [ ] POST /api/query-templates/{id}/set-default
  - [ ] POST /api/query-templates/{id}/duplicate
  - [ ] GET /api/query-templates/{id}/preview（预览SQL+结果）

### 数据查询API
- [ ] POST /api/data/{model}/query（通用查询接口）
- [ ] GET /api/data/{model}（列表查询+查询模板）
- [ ] GET /api/data/{model}/{id}（详情查询）
- [ ] POST /api/data/{model}（创建数据）
- [ ] PUT /api/data/{model}/{id}（更新数据）
- [ ] DELETE /api/data/{model}/{id}（删除数据）
- [ ] POST /api/data/{model}/batch-create（批量创建）
- [ ] PUT /api/data/{model}/batch-update（批量更新）
- [ ] DELETE /api/data/{model}/batch-delete（批量删除）
- [ ] GET /api/data/{model}/statistics（数据统计）
- [ ] POST /api/data/{model}/aggregate（聚合查询）

---

## 阶段3：高级功能 (预计3周)

### 树形结构支持
- [ ] 扩展Model表（添加is_tree/parent_field/path_field/level_field，需migration）
- [ ] TreeConfig配置管理
- [ ] Tree Service
  - [ ] GetTree方法（完整树+层级限制+懒加载）
  - [ ] GetChildren方法（直接子节点）
  - [ ] GetPath方法（根到节点路径）
  - [ ] GetDescendants方法（所有后代）
  - [ ] AddNode方法（自动计算path+level）
  - [ ] MoveNode方法（移动+重算路径+循环引用检测）
  - [ ] DeleteNode方法（级联删除或移动子节点）
  - [ ] 路径计算器（格式：/1/3/5/）
  - [ ] 层级计算器
  - [ ] 循环引用检测器
- [ ] Tree API
  - [ ] GET /api/tree/{model}（获取树）
  - [ ] POST /api/tree/{model}/node（添加节点）
  - [ ] PUT /api/tree/{model}/node/{id}（更新节点）
  - [ ] DELETE /api/tree/{model}/node/{id}（删除节点）
  - [ ] POST /api/tree/{model}/node/{id}/move（移动节点）
  - [ ] GET /api/tree/{model}/node/{id}/path（节点路径）
  - [ ] GET /api/tree/{model}/node/{id}/children（子节点）
  - [ ] GET /api/tree/{model}/node/{id}/descendants（后代）

### 主子表支持
- [ ] ModelRelation Model（需新建表，主子表关系配置）
- [ ] MasterDetail Service
  - [ ] CreateMasterDetail方法（事务：主表+子表批量）
  - [ ] GetMasterDetail方法（主表+子表列表）
  - [ ] UpdateMasterDetail方法（事务：主表+子表增删改）
  - [ ] DeleteMasterDetail方法（级联删除或设置NULL）
  - [ ] 事务管理器增强（REQUIRED/REQUIRES_NEW/NESTED）
  - [ ] 级联操作处理
- [ ] MasterDetail API
  - [ ] POST /api/master-detail/{master}/{detail}/create
  - [ ] GET /api/master-detail/{master}/{detail}/{masterId}
  - [ ] PUT /api/master-detail/{master}/{detail}/{masterId}
  - [ ] DELETE /api/master-detail/{master}/{detail}/{masterId}
  - [ ] POST /api/master-detail/{master}/{detail}/{masterId}/details（添加子表）
  - [ ] PUT /api/master-detail/{master}/{detail}/{masterId}/details/{detailId}
  - [ ] DELETE /api/master-detail/{master}/{detail}/{masterId}/details/{detailId}

### 数据导入导出
- [ ] Excel Service
  - [ ] ExportToExcel方法（样式配置+大数据量流式写入）
  - [ ] GenerateTemplate方法（根据模型生成模板）
  - [ ] ParseExcel方法（解析+验证）
  - [ ] ImportFromExcel方法（批量插入+错误收集）
  - [ ] 数据验证（类型+长度+必填+唯一性+业务规则）
- [ ] CSV Service
  - [ ] ExportToCSV方法（分隔符配置+编码UTF-8/GBK）
  - [ ] ImportFromCSV方法（编码识别+解析+导入）
- [ ] JSON Service
  - [ ] ExportToJSON方法（完整结构+关联数据可选）
  - [ ] ImportFromJSON方法
- [ ] Import/Export API
  - [ ] POST /api/data/{model}/export（导出+格式+字段+筛选）
  - [ ] GET /api/data/{model}/import-template（下载模板）
  - [ ] POST /api/data/{model}/import（导入+进度+错误返回）
  - [ ] GET /api/data/{model}/import/{taskId}/status（导入进度）

### 审计日志
- [ ] OperationLog Model（sys_operation_logs表）
- [ ] DataChangeLog Model（sys_data_change_logs表）
- [ ] AccessLog Model（sys_access_logs表）
- [ ] Log Middleware
  - [ ] OperationLog中间件（记录请求/响应+异步写入）
  - [ ] DataChangeLog拦截器（记录变更前后+对比字段）
  - [ ] AccessLog中间件（记录访问+响应时间）
- [ ] Log Service
  - [ ] QueryOperationLogs方法（按用户/模块/操作/时间筛选）
  - [ ] QueryDataChangeLogs方法（按表/记录/用户/时间）
  - [ ] QueryAccessLogs方法（按用户/路径/状态/时间）
  - [ ] 异常行为检测（暴力破解/异常登录/异常操作）
- [ ] Log API
  - [ ] GET /api/logs/operations（操作日志）
  - [ ] GET /api/logs/data-changes（数据变更日志）
  - [ ] GET /api/logs/access（访问日志）
  - [ ] GET /api/logs/statistics（日志统计）

---

## 阶段4：系统完善 (预计2周)

### 性能优化
- [ ] Redis缓存集成
  - [ ] 连接池配置
  - [ ] 缓存管理器（Get/Set/Delete/Clear）
  - [ ] 查询结果缓存（键策略+失效策略）
  - [ ] 模型配置缓存
  - [ ] 字典数据缓存
  - [ ] 缓存穿透防护（布隆过滤器+空值缓存）
  - [ ] 缓存击穿防护（互斥锁+热点永不过期）
  - [ ] 缓存雪崩防护（过期时间随机化+限流）
- [ ] 查询优化
  - [ ] 分析慢查询日志
  - [ ] 添加必要索引
  - [ ] 优化N+1查询（批量查询+关联加载）
  - [ ] 添加查询超时控制
- [ ] 批量操作优化
  - [ ] 批量插入（合理批次大小+内存控制）
  - [ ] 批量更新（事务批次提交）
  - [ ] 批量删除（分批处理）
- [ ] 连接池优化
  - [ ] 调整连接池参数
  - [ ] 连接池监控（活跃/等待/创建销毁速率）
  - [ ] 连接泄漏检测

### 监控与告警
- [ ] Prometheus集成
  - [ ] 集成客户端库
  - [ ] HTTP请求指标（总数/耗时分布/错误率）
  - [ ] 数据库指标（查询数/慢查询/连接数）
  - [ ] 缓存指标（命中率/容量/使用率）
  - [ ] 系统指标（CPU/内存/goroutine）
  - [ ] 暴露/metrics端点
- [ ] 健康检查
  - [ ] GET /api/health（数据库+Redis+磁盘空间）
  - [ ] GET /api/ping（探活）
- [ ] 日志增强
  - [ ] 链路追踪ID（请求级别）
  - [ ] 结构化日志（JSON格式）
  - [ ] 日志级别动态调整
- [ ] 告警管理
  - [ ] 告警规则配置（CPU>80%/内存>80%/错误率>5%/响应时间>1s）
  - [ ] 告警通知（邮件/钉钉/企业微信）
  - [ ] 告警抑制（合并+频率限制）

### 安全加固
- [ ] 输入验证增强
  - [ ] SQL注入防护（参数化查询+禁止动态拼接）
  - [ ] XSS防护（输出HTML编码+CSP头）
  - [ ] CSRF防护（Token验证+SameSite Cookie）
  - [ ] 文件上传安全（类型验证+大小限制+病毒扫描）
- [ ] 访问控制增强
  - [ ] 接口限流（基于IP+基于用户+滑动窗口）
  - [ ] IP黑白名单（配置+自动封禁）
  - [ ] 接口权限细化（字段级+数据范围）
- [ ] 数据安全
  - [ ] 敏感数据加密（AES-256）
  - [ ] 数据脱敏（手机号/身份证/银行卡）
  - [ ] SSL/TLS配置（HTTPS+证书管理）
- [ ] 安全审计强化
  - [ ] 敏感操作审计
  - [ ] 权限变更审计
  - [ ] 配置修改审计

### 文档与测试
- [ ] Swagger API文档
  - [ ] 集成Swagger
  - [ ] 添加API注释
  - [ ] 自动生成文档
  - [ ] 接口示例
  - [ ] 错误码说明
- [ ] 单元测试
  - [ ] SQL生成引擎测试
  - [ ] 模型构建器测试
  - [ ] CRUD Service测试
  - [ ] 测试覆盖率>60%
- [ ] 集成测试
  - [ ] API端到端测试
  - [ ] 数据库集成测试
- [ ] 数据字典生成
  - [ ] 自动生成Markdown数据字典
  - [ ] 自动生成Excel数据字典
- [ ] 部署文档
  - [ ] 环境要求
  - [ ] 配置说明
  - [ ] 启动步骤
  - [ ] 常见问题

---

## 验收标准

### 阶段1验收
- 可完整登录认证
- 可创建数据源并测试连接
- 可从数据库表构建模型
- 所有基础API可用

### 阶段2验收
- SQL生成引擎可生成复杂SQL
- 可批量生成CRUD接口
- 动态路由正常工作
- 查询模板功能完整

### 阶段3验收
- 树形数据管理正常
- 主子表操作正确
- 导入导出功能可用
- 审计日志完整记录

### 阶段4验收
- 系统性能达标（响应<200ms）
- 监控指标可查看
- 安全测试通过
- 文档完整

---

**预计总工期**: 14周（3.5个月）  
**关键路径**: 阶段1基础框架 → 阶段2 SQL生成引擎 → 阶段2 CRUD接口生成

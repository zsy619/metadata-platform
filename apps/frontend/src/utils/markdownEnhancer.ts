/**
 * Markdown 编辑器增强工具
 * 包含图片处理、模板系统、表情符号等企业级功能
 */

/**
 * 文档模板接口
 */
export interface DocumentTemplate {
    id: string
    name: string
    description: string
    category: string
    content: string
    icon: string
}

/**
 * 预定义文档模板
 */
export const documentTemplates: DocumentTemplate[] = [
    {
        id: 'api-doc',
        name: 'API 文档',
        description: '标准的 API 接口文档模板',
        category: '技术文档',
        icon: 'Document',
        content: `# {{API 名称}}

## 接口描述

{{简要描述 API 的功能和用途}}

## 请求 URL

\`\`\`
{{HTTP_METHOD}} {{URL_PATH}}
\`\`\`

## 请求参数

### 请求头 (Headers)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| Content-Type | string | 是 | application/json |
| Authorization | string | 是 | Bearer {{token}} |

### 请求体 (Body)

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| {{param1}} | string | 是 | {{描述}} |
| {{param2}} | number | 否 | {{描述}} |

## 响应结果

### 成功响应

\`\`\`json
{
    "code": 200,
    "message": "success",
    "data": {
        {{response_fields}}
    }
}
\`\`\`

### 错误响应

\`\`\`json
{
    "code": {{error_code}},
    "message": "{{error_message}}"
}
\`\`\`

## 请求示例

\`\`\`bash
curl -X {{METHOD}} {{URL}} \\
  -H 'Content-Type: application/json' \\
  -H 'Authorization: Bearer {{token}}' \\
  -d '{{request_body}}'
\`\`\`

## 注意事项

- {{注意事项 1}}
- {{注意事项 2}}

## 相关链接

- [相关文档 1]({{link1}})
- [相关文档 2]({{link2}})
`
    },
    {
        id: 'meeting-notes',
        name: '会议纪要',
        description: '标准的会议纪要模板',
        category: '会议文档',
        icon: 'Calendar',
        content: `# 会议纪要

## 基本信息

- **会议主题**: {{会议主题}}
- **会议时间**: {{YYYY-MM-DD HH:mm}}
- **会议地点**: {{会议室/线上链接}}
- **主持人**: {{主持人姓名}}
- **记录人**: {{记录人姓名}}
- **参会人员**: {{参会人员名单}}
- **缺席人员**: {{缺席人员名单}}

## 会议议程

### 1. {{议题一}}

**讨论内容**:
- {{讨论要点 1}}
- {{讨论要点 2}}

**决议**:
- [ ] {{决议事项 1}}
- [ ] {{决议事项 2}}

**负责人**: {{负责人}}
**截止时间**: {{YYYY-MM-DD}}

### 2. {{议题二}}

**讨论内容**:
- {{讨论要点}}

**决议**:
- {{决议事项}}

## 待办事项

| 序号 | 待办事项 | 负责人 | 截止时间 | 状态 |
|------|----------|--------|----------|------|
| 1 | {{待办 1}} | {{负责人}} | {{日期}} | ☐ 未开始 |
| 2 | {{待办 2}} | {{负责人}} | {{日期}} | ☐ 未开始 |

## 下次会议

- **时间**: {{YYYY-MM-DD HH:mm}}
- **议题**: {{主要议题}}

## 附件

- [{{附件名称 1}}]({{链接 1}})
- [{{附件名称 2}}]({{链接 2}})
`
    },
    {
        id: 'tech-proposal',
        name: '技术方案',
        description: '技术方案设计文档模板',
        category: '技术文档',
        icon: 'Setting',
        content: `# 技术方案设计

## 文档信息

| 项目 | 内容 |
|------|------|
| 文档版本 | v1.0.0 |
| 撰写人 | {{姓名}} |
| 撰写日期 | {{YYYY-MM-DD}} |
| 审核人 | {{姓名}} |
| 审核日期 | {{YYYY-MM-DD}} |

## 1. 背景与目标

### 1.1 项目背景

{{描述项目的背景、业务需求和痛点}}

### 1.2 技术目标

- **性能目标**: {{如：QPS 达到 10000}}
- **可用性目标**: {{如：99.9%}}
- **扩展性目标**: {{如：支持水平扩展}}

## 2. 需求分析

### 2.1 功能性需求

| 编号 | 需求描述 | 优先级 |
|------|----------|--------|
| FR-1 | {{需求 1}} | P0 |
| FR-2 | {{需求 2}} | P1 |

### 2.2 非功能性需求

| 编号 | 需求描述 | 指标 |
|------|----------|------|
| NFR-1 | 性能要求 | {{具体指标}} |
| NFR-2 | 安全要求 | {{具体要求}} |

## 3. 技术方案

### 3.1 架构设计

#### 3.1.1 系统架构图

\`\`\`mermaid
graph TB
    A[客户端] --> B[负载均衡]
    B --> C[服务层]
    C --> D[数据层]
\`\`\`

#### 3.1.2 技术选型

| 层级 | 技术 | 说明 |
|------|------|------|
| 前端 | {{技术栈}} | {{理由}} |
| 后端 | {{技术栈}} | {{理由}} |
| 数据库 | {{技术栈}} | {{理由}} |

### 3.2 核心流程

\`\`\`mermaid
sequenceDiagram
    participant A as 用户
    participant B as 服务 A
    participant C as 服务 B
    
    A->>B: 请求
    B->>C: 调用
    C-->>B: 响应
    B-->>A: 返回
\`\`\`

## 4. 数据设计

### 4.1 数据模型

\`\`\`sql
CREATE TABLE {{table_name}} (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    {{field1}} VARCHAR(255) NOT NULL COMMENT '{{注释}}',
    {{field2}} INT DEFAULT 0 COMMENT '{{注释}}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT '{{表注释}}';
\`\`\`

### 4.2 ER 图

\`\`\`mermaid
erDiagram
    {{ENTITY1}} ||--o{ {{ENTITY2}} : contains
\`\`\`

## 5. 接口设计

### 5.1 接口列表

| 接口名 | 方法 | URL | 说明 |
|--------|------|-----|------|
| {{接口 1}} | POST | {{/api/xxx}} | {{说明}} |

## 6. 部署方案

### 6.1 部署架构

\`\`\`mermaid
graph LR
    A[用户] --> B[CDN]
    B --> C[Web 服务器]
    C --> D[应用服务器]
    D --> E[数据库]
\`\`\`

### 6.2 资源配置

| 资源 | 规格 | 数量 | 说明 |
|------|------|------|------|
| 应用服务器 | {{4C8G}} | {{2}} | {{说明}} |
| 数据库 | {{8C16G}} | {{1}} | {{主从配置}} |

## 7. 监控与告警

### 7.1 监控指标

- **应用监控**: {{指标列表}}
- **系统监控**: {{指标列表}}
- **业务监控**: {{指标列表}}

### 7.2 告警策略

| 指标 | 阈值 | 告警级别 | 通知方式 |
|------|------|----------|----------|
| CPU 使用率 | > 80% | 警告 | 邮件 |
| 错误率 | > 1% | 严重 | 短信 + 电话 |

## 8. 风险评估

| 风险项 | 影响程度 | 发生概率 | 应对措施 |
|--------|----------|----------|----------|
| {{风险 1}} | 高 | 中 | {{应对措施}} |

## 9. 项目计划

### 9.1 里程碑

| 阶段 | 时间 | 交付物 |
|------|------|--------|
| 需求评审 | {{日期}} | 需求文档 |
| 开发 | {{日期范围}} | 代码 |
| 测试 | {{日期范围}} | 测试报告 |
| 上线 | {{日期}} | 线上环境 |

## 10. 参考资料

- [{{参考资料 1}}]({{链接}})
- [{{参考资料 2}}]({{链接}})
`
    },
    {
        id: 'project-plan',
        name: '项目计划',
        description: '项目规划和管理模板',
        category: '项目管理',
        icon: 'Management',
        content: `# 项目计划

## 项目概述

| 项目名称 | {{项目名称}} |
|----------|--------------|
| 项目经理 | {{姓名}} |
| 开始日期 | {{YYYY-MM-DD}} |
| 预计结束 | {{YYYY-MM-DD}} |
| 当前状态 | {{规划中/进行中}} |

## 项目目标

### SMART 目标

- **Specific（具体的）**: {{具体目标}}
- **Measurable（可衡量的）**: {{衡量指标}}
- **Achievable（可实现的）**: {{可行性说明}}
- **Relevant（相关的）**: {{业务价值}}
- **Time-bound（有时限的）**: {{完成时间}}

## 项目范围

### 包含内容

- ✅ {{范围内的工作 1}}
- ✅ {{范围内的工作 2}}

### 不包含内容

- ❌ {{范围外的工作 1}}
- ❌ {{范围外的工作 2}}

## 项目团队

| 角色 | 姓名 | 职责 | 投入比例 |
|------|------|------|----------|
| 项目经理 | {{姓名}} | {{职责}} | 100% |
| 技术负责人 | {{姓名}} | {{职责}} | 100% |
| 开发人员 | {{姓名}} | {{职责}} | 100% |

## 项目里程碑

| 里程碑 | 计划日期 | 实际日期 | 状态 | 交付物 |
|--------|----------|----------|------|--------|
| 项目启动 | {{日期}} | - | ☐ 未开始 | 项目章程 |
| 需求确认 | {{日期}} | - | ☐ 未开始 | 需求文档 |
| 设计评审 | {{日期}} | - | ☐ 未开始 | 设计文档 |
| 开发完成 | {{日期}} | - | ☐ 未开始 | 代码 |
| 测试完成 | {{日期}} | - | ☐ 未开始 | 测试报告 |
| 项目上线 | {{日期}} | - | ☐ 未开始 | 线上环境 |

## 项目进度

### 甘特图

\`\`\`mermaid
gantt
    title 项目进度计划
    dateFormat  YYYY-MM-DD
    section 阶段一
    需求分析 :a1, 2024-01-01, 14d
    section 阶段二
    系统设计 :a2, after a1, 14d
    section 阶段三
    开发实现 :a3, after a2, 30d
    section 阶段四
    测试验收 :a4, after a3, 14d
\`\`\`

## 风险管理

| 风险编号 | 风险描述 | 影响程度 | 发生概率 | 应对策略 | 负责人 |
|----------|----------|----------|----------|----------|--------|
| RISK-001 | {{风险描述}} | 高 | 中 | {{应对策略}} | {{姓名}} |

## 沟通计划

| 会议类型 | 频率 | 时间 | 参与人 | 内容 |
|----------|------|------|--------|------|
| 项目例会 | 每周 | 周一 10:00 | 全体成员 | 进度同步 |
| 技术评审 | 按需 | - | 技术人员 | 方案评审 |

## 预算估算

| 类别 | 预算金额 | 实际支出 | 说明 |
|------|----------|----------|------|
| 人力成本 | {{金额}} | - | {{说明}} |
| 设备成本 | {{金额}} | - | {{说明}} |
| 其他成本 | {{金额}} | - | {{说明}} |
| **合计** | **{{总金额}}** | - | - |

## 变更管理

| 变更编号 | 变更内容 | 变更日期 | 影响评估 | 审批状态 |
|----------|----------|----------|----------|----------|
| CHG-001 | {{变更内容}} | {{日期}} | {{影响}} | {{状态}} |

## 项目报告

### 周报

#### 本周完成

- [x] {{完成事项 1}}
- [x] {{完成事项 2}}

#### 下周计划

- [ ] {{计划事项 1}}
- [ ] {{计划事项 2}}

#### 风险与问题

- ⚠️ {{风险/问题 1}}
- ⚠️ {{风险/问题 2}}
`
    }
]

/**
 * 表情符号列表
 */
export const emojiList = [
    // 笑脸表情
    { emoji: '😀', name: '大笑' },
    { emoji: '😃', name: '笑脸' },
    { emoji: '😄', name: '开心' },
    { emoji: '😁', name: '嘻嘻' },
    { emoji: '😊', name: '微笑' },
    { emoji: '😂', name: '笑哭' },
    { emoji: '🤣', name: '大笑' },
    { emoji: '😍', name: '花痴' },
    { emoji: '😘', name: '飞吻' },
    { emoji: '😜', name: '鬼脸' },
    
    // 手势表情
    { emoji: '👍', name: '赞' },
    { emoji: '👎', name: '踩' },
    { emoji: '👏', name: '鼓掌' },
    { emoji: '🙏', name: '祈祷' },
    { emoji: '🤝', name: '握手' },
    { emoji: '💪', name: '肌肉' },
    { emoji: '✊', name: '拳头' },
    { emoji: '✌️', name: '胜利' },
    { emoji: '👌', name: 'OK' },
    
    // 爱心表情
    { emoji: '❤️', name: '红心' },
    { emoji: '💕', name: '双心' },
    { emoji: '💖', name: '闪光心' },
    { emoji: '💗', name: '心跳' },
    { emoji: '💓', name: '心动' },
    { emoji: '💞', name: '旋转心' },
    { emoji: '💘', name: '爱神之箭' },
    { emoji: '💝', name: '心形礼盒' },
    
    // 星星表情
    { emoji: '⭐', name: '星星' },
    { emoji: '🌟', name: '闪星' },
    { emoji: '✨', name: '火花' },
    { emoji: '💫', name: '头晕' },
    
    // 天气自然
    { emoji: '☀️', name: '太阳' },
    { emoji: '🌙', name: '月亮' },
    { emoji: '⭐', name: '星星' },
    { emoji: '🌈', name: '彩虹' },
    { emoji: '🔥', name: '火焰' },
    { emoji: '💧', name: '水滴' },
    { emoji: '❄️', name: '雪花' },
    
    // 动物
    { emoji: '🐶', name: '狗' },
    { emoji: '🐱', name: '猫' },
    { emoji: '🐭', name: '老鼠' },
    { emoji: '🐹', name: '仓鼠' },
    { emoji: '🐰', name: '兔子' },
    { emoji: '🦊', name: '狐狸' },
    { emoji: '🐻', name: '熊' },
    { emoji: '🐼', name: '熊猫' },
    
    // 食物
    { emoji: '🍎', name: '苹果' },
    { emoji: '🍌', name: '香蕉' },
    { emoji: '🍇', name: '葡萄' },
    { emoji: '🍉', name: '西瓜' },
    { emoji: '🍊', name: '橘子' },
    { emoji: '🍋', name: '柠檬' },
    { emoji: '🍓', name: '草莓' },
    { emoji: '🍒', name: '樱桃' },
    
    // 运动
    { emoji: '⚽', name: '足球' },
    { emoji: '🏀', name: '篮球' },
    { emoji: '🏈', name: '橄榄球' },
    { emoji: '⚾', name: '棒球' },
    { emoji: '🎾', name: '网球' },
    { emoji: '🏐', name: '排球' },
    { emoji: '🏓', name: '乒乓球' },
    
    // 办公用品
    { emoji: '📝', name: '备忘录' },
    { emoji: '📌', name: '图钉' },
    { emoji: '📎', name: '回形针' },
    { emoji: '📊', name: '图表' },
    { emoji: '📈', name: '上涨' },
    { emoji: '📉', name: '下跌' },
    { emoji: '✅', name: '完成' },
    { emoji: '❌', name: '错误' },
    { emoji: '⚠️', name: '警告' },
    { emoji: '💡', name: '灯泡' },
    { emoji: '🎯', name: '靶心' },
    { emoji: '📅', name: '日历' },
    { emoji: '📆', name: '撕页日历' },
    { emoji: '🔔', name: '铃铛' },
    { emoji: '🔖', name: '书签' },
    { emoji: '📑', name: '标签' }
]

/**
 * 图片压缩配置
 */
export interface ImageCompressOptions {
    maxWidth?: number
    maxHeight?: number
    quality?: number
    format?: 'image/jpeg' | 'image/png' | 'image/webp'
}

/**
 * 压缩图片
 */
export function compressImage(
    file: File,
    options: ImageCompressOptions = {}
): Promise<File> {
    const {
        maxWidth = 1920,
        maxHeight = 1080,
        quality = 0.8,
        format = 'image/jpeg'
    } = options

    return new Promise((resolve, reject) => {
        const reader = new FileReader()
        
        reader.onload = (e) => {
            const img = new Image()
            img.onload = () => {
                // 计算缩放比例
                let width = img.width
                let height = img.height
                
                if (width > maxWidth || height > maxHeight) {
                    const ratio = Math.min(maxWidth / width, maxHeight / height)
                    width = Math.floor(width * ratio)
                    height = Math.floor(height * ratio)
                }
                
                // 创建 canvas
                const canvas = document.createElement('canvas')
                canvas.width = width
                canvas.height = height
                
                const ctx = canvas.getContext('2d')
                if (!ctx) {
                    reject(new Error('无法获取 canvas 上下文'))
                    return
                }
                
                // 绘制图片
                ctx.drawImage(img, 0, 0, width, height)
                
                // 转换为 Blob
                canvas.toBlob(
                    (blob) => {
                        if (!blob) {
                            reject(new Error('图片压缩失败'))
                            return
                        }
                        
                        const compressedFile = new File([blob], file.name, {
                            type: format
                        })
                        resolve(compressedFile)
                    },
                    format,
                    quality
                )
            }
            
            img.onerror = () => {
                reject(new Error('图片加载失败'))
            }
            
            img.src = e.target?.result as string
        }
        
        reader.onerror = () => {
            reject(new Error('文件读取失败'))
        }
        
        reader.readAsDataURL(file)
    })
}

/**
 * 验证图片文件
 */
export function validateImageFile(file: File): { valid: boolean; error?: string } {
    // 检查文件类型
    const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
    if (!validTypes.includes(file.type)) {
        return {
            valid: false,
            error: '不支持的图片格式，仅支持 JPG、PNG、GIF、WebP'
        }
    }
    
    // 检查文件大小（最大 10MB）
    const maxSize = 10 * 1024 * 1024
    if (file.size > maxSize) {
        return {
            valid: false,
            error: '图片大小不能超过 10MB'
        }
    }
    
    return { valid: true }
}

/**
 * 获取模板列表
 */
export function getTemplates(category?: string): DocumentTemplate[] {
    if (!category) {
        return documentTemplates
    }
    return documentTemplates.filter(t => t.category === category)
}

/**
 * 根据 ID 获取模板
 */
export function getTemplateById(id: string): DocumentTemplate | undefined {
    return documentTemplates.find(t => t.id === id)
}

/**
 * 替换模板中的占位符
 */
export function renderTemplate(
    template: string,
    variables: Record<string, string>
): string {
    let result = template
    Object.entries(variables).forEach(([key, value]) => {
        result = result.replace(new RegExp(`{{${key}}}`, 'g'), value)
    })
    return result
}

/**
 * 提取模板中的占位符
 */
export function extractTemplateVariables(template: string): string[] {
    const matches = template.matchAll(/{{([^}]+)}}/g)
    return Array.from(matches, match => match[1])
}

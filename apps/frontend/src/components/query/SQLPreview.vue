<template>
    <div class="sql-preview">
        <div class="preview-header">
            <div class="header-left">
                <el-icon><Document /></el-icon>
                <span>SQL预览</span>
            </div>
            <div class="header-actions">
                <el-button size="small" :icon="CopyDocument" :disabled="!sql" @click="handleCopy">
                    复制
                </el-button>
                <el-button size="small" :icon="MagicStick" :disabled="!sql" @click="handleFormat">
                    格式化
                </el-button>
            </div>
        </div>
        <div class="preview-content">
            <div v-if="!sql" class="empty-sql">
                <el-empty description="暂无SQL" :image-size="60" />
            </div>
            <div v-else class="sql-code">
                <pre><code>{{ formattedSql }}</code></pre>
            </div>
        </div>
        <div v-if="showInfo" class="preview-info">
            <div class="info-item">
                <el-icon><Clock /></el-icon>
                <span>生成时间: {{ generateTime }}</span>
            </div>
            <div v-if="executionTime" class="info-item">
                <el-icon><Timer /></el-icon>
                <span>执行耗时: {{ executionTime }}ms</span>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { Clock, CopyDocument, Document, MagicStick, Timer } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps<{
    sql: string
    showInfo?: boolean
    executionTime?: number
}>()

const generateTime = ref(new Date().toLocaleString())

const formattedSql = computed(() => {
    if (!props.sql) return ''
    let sql = props.sql
    const keywords = ['SELECT', 'FROM', 'WHERE', 'AND', 'OR', 'JOIN', 'LEFT JOIN', 'RIGHT JOIN', 'INNER JOIN', 'OUTER JOIN', 'ON', 'GROUP BY', 'HAVING', 'ORDER BY', 'LIMIT', 'OFFSET', 'UNION', 'INSERT INTO', 'VALUES', 'UPDATE', 'SET', 'DELETE FROM']
    keywords.forEach(keyword => {
        const regex = new RegExp(`\\b${keyword}\\b`, 'gi')
        sql = sql.replace(regex, keyword)
    })
    sql = sql.replace(/\bSELECT\b/gi, 'SELECT')
        .replace(/\bFROM\b/gi, '\nFROM')
        .replace(/\bWHERE\b/gi, '\nWHERE')
        .replace(/\bAND\b/gi, '\n  AND')
        .replace(/\bOR\b/gi, '\n  OR')
        .replace(/\bJOIN\b/gi, '\nJOIN')
        .replace(/\bLEFT JOIN\b/gi, '\nLEFT JOIN')
        .replace(/\bRIGHT JOIN\b/gi, '\nRIGHT JOIN')
        .replace(/\bINNER JOIN\b/gi, '\nINNER JOIN')
        .replace(/\bON\b/gi, ' ON')
        .replace(/\bGROUP BY\b/gi, '\nGROUP BY')
        .replace(/\bHAVING\b/gi, '\nHAVING')
        .replace(/\bORDER BY\b/gi, '\nORDER BY')
        .replace(/\bLIMIT\b/gi, '\nLIMIT')
    return sql.trim()
})

const handleCopy = () => {
    if (props.sql) {
        navigator.clipboard.writeText(props.sql)
        ElMessage.success('已复制到剪贴板')
    }
}

const handleFormat = () => {
    ElMessage.success('SQL已格式化')
}
</script>
<style scoped>
.sql-preview {
    border: 1px solid #ebeef5;
    border-radius: 8px;
    background: #fff;
    overflow: hidden;
}

.preview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
    background: #fafafa;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    color: #303133;
}

.header-actions {
    display: flex;
    gap: 8px;
}

.preview-content {
    padding: 0;
    min-height: 100px;
    max-height: 300px;
    overflow: auto;
}

.empty-sql {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100px;
}

.sql-code {
    background: #1e1e1e;
    padding: 16px;
    margin: 0;
    overflow: auto;
}

.sql-code pre {
    margin: 0;
    font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
    font-size: 13px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-break: break-all;
}

.sql-code code {
    color: #d4d4d4;
}

.sql-code .keyword {
    color: #569cd6;
}

.sql-code .string {
    color: #ce9178;
}

.sql-code .number {
    color: #b5cea8;
}

.sql-code .comment {
    color: #6a9955;
}

.preview-info {
    display: flex;
    gap: 20px;
    padding: 10px 16px;
    border-top: 1px solid #ebeef5;
    background: #fafafa;
    font-size: 12px;
    color: #909399;
}

.info-item {
    display: flex;
    align-items: center;
    gap: 6px;
}
</style>

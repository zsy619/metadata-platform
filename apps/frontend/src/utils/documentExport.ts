/**
 * 文档导出工具
 * 支持导出为 PDF、Word、Markdown、HTML 等格式
 */

/**
 * 导出为 Markdown 文件
 */
export function exportToMarkdown(content: string, filename: string = 'document') {
    const blob = new Blob([content], { type: 'text/markdown' })
    const url = URL.createObjectURL(blob)
    downloadFile(url, `${filename}.md`)
    URL.revokeObjectURL(url)
}

/**
 * 导出为 HTML 文件
 */
export function exportToHtml(content: string, filename: string = 'document', title: string = 'Document') {
    const html = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title}</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css">
    <style>
        body {
            max-width: 900px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        @media print {
            body {
                padding: 0;
            }
        }
    </style>
</head>
<body class="markdown-body">
    ${content}
</body>
</html>
`.trim()
    
    const blob = new Blob([html], { type: 'text/html' })
    const url = URL.createObjectURL(blob)
    downloadFile(url, `${filename}.html`)
    URL.revokeObjectURL(url)
}

/**
 * 导出为 PDF（使用浏览器打印功能）
 */
export function exportToPdf(content: string, title: string = 'Document') {
    // 创建一个临时窗口用于打印
    const printWindow = window.open('', '_blank')
    if (!printWindow) {
        console.error('无法打开打印窗口')
        return
    }
    
    const html = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${title}</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.2.0/github-markdown.min.css">
    <style>
        @media print {
            body {
                padding: 0;
                margin: 0;
            }
            .markdown-body {
                padding: 20px;
            }
            @page {
                margin: 2cm;
            }
        }
    </style>
</head>
<body class="markdown-body">
    ${content}
</body>
</html>
`.trim()
    
    printWindow.document.write(html)
    printWindow.document.close()
    
    // 等待内容加载后打印
    printWindow.onload = () => {
        printWindow.print()
    }
}

/**
 * 导出为 Word 文档（.docx）
 */
export function exportToWord(content: string, filename: string = 'document', title: string = 'Document') {
    // 将 Markdown 转换为简单的 HTML
    const html = `
<html xmlns:o='urn:schemas-microsoft-com:office:office' 
      xmlns:w='urn:schemas-microsoft-com:office:word' 
      xmlns='http://www.w3.org/TR/REC-html40'>
<head>
    <meta charset='utf-8'>
    <title>${title}</title>
    <style>
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            font-size: 14px;
            line-height: 1.6;
        }
        h1, h2, h3, h4, h5, h6 {
            color: #333;
            margin-top: 24px;
            margin-bottom: 16px;
        }
        h1 { font-size: 24px; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h2 { font-size: 20px; border-bottom: 1px solid #eaecef; padding-bottom: 0.3em; }
        h3 { font-size: 16px; }
        code {
            background-color: #f6f8fa;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: Consolas, Monaco, monospace;
        }
        pre {
            background-color: #f6f8fa;
            padding: 16px;
            border-radius: 3px;
            overflow: auto;
        }
        pre code {
            background: none;
            padding: 0;
        }
        blockquote {
            padding: 0 1em;
            color: #6a737d;
            border-left: 0.25em solid #dfe2e5;
            margin: 0;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 16px 0;
        }
        th, td {
            border: 1px solid #dfe2e5;
            padding: 6px 13px;
        }
        th {
            background-color: #f6f8fa;
            font-weight: 600;
        }
        img {
            max-width: 100%;
        }
    </style>
</head>
<body>
    ${content}
</body>
</html>
`.trim()
    
    const blob = new Blob([html], { type: 'application/msword' })
    const url = URL.createObjectURL(blob)
    downloadFile(url, `${filename}.doc`)
    URL.revokeObjectURL(url)
}

/**
 * 下载文件的辅助函数
 */
function downloadFile(url: string, filename: string) {
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
}

/**
 * 复制到剪贴板
 */
export async function copyToClipboard(text: string): Promise<boolean> {
    try {
        await navigator.clipboard.writeText(text)
        return true
    } catch (error) {
        console.error('复制到剪贴板失败:', error)
        return false
    }
}

/**
 * 打印文档
 */
export function printDocument(content: string, title: string = 'Document') {
    const printWindow = window.open('', '_blank')
    if (!printWindow) {
        console.error('无法打开打印窗口')
        return
    }
    
    const html = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>${title}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            padding: 20px;
        }
        @media print {
            body {
                padding: 0;
            }
        }
        pre {
            white-space: pre-wrap;
            word-wrap: break-word;
        }
    </style>
</head>
<body>
    ${content}
</body>
</html>
`.trim()
    
    printWindow.document.write(html)
    printWindow.document.close()
    printWindow.onload = () => {
        printWindow.print()
    }
}

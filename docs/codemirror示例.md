## 1. 参考示例

https://www.cnblogs.com/Lu-Lu/p/16265815.html

## 2.示例代码

```html
<textarea
   :id="'mycode'+(index*1+1)"
   :ref="'mycode'+(index*1+1)"
   v-model="item.sqlContent"
   class="CodeMirror-hints"
   :class="'mycode'+(index*1+1)"
   placeholder="按Ctrl键进行代码提示"
/>
```

```javascript
import CodeMirror from 'codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/solarized.css'
import 'codemirror/addon/edit/closebrackets.js'
import 'codemirror/mode/sql/sql.js'
import 'codemirror/addon/display/autorefresh'

import 'codemirror/addon/hint/show-hint.js'
import 'codemirror/addon/hint/show-hint.css'

import 'codemirror/addon/hint/sql-hint.js'
```

表及字段处理：

```javascript
init (ind, sqlId) {
      const _this = this
      let code = 'mycode'+ind
      this.$nextTick(() => {
        // 实例初始化
        this.editor = CodeMirror.fromTextArea(this.$refs[code][0], {
          tabSize: 4,
          mode: 'text/x-mysql', //语言sql
          theme: 'solarized', // 主题
          autoCloseBrackets: true, // 在键入时自动关闭括号和引号
          autoRefresh: true,
          styleActiveLine: true,
          lineNumbers: true,
          line: true,
          lineWrapping: true,
          readOnly: "nocursor", //默认只读，无光标，如果是true好像有bug
          hintOptions: { //自定义提示内容
            completeSingle: false,
            hint: this.handleShowHint2,
          },
          // extraKeys: {
          //   'Ctrl-Space': editor => {
          //     editor.showHint()
          //   }
          // }
        })
        // 下方代码基本不用改动
        this.editor.on('keypress', editor => {
          const __Cursor = editor.getDoc().getCursor()
          const __Token = editor.getTokenAt(__Cursor)
          // console.log(__Cursor)
          // console.log(__Token)
          if (
            __Token.type &&
            __Token.type !== 'string' &&
            __Token.type !== 'punctuation' &&
            __Token.string.indexOf('.') === -1
          ) {
            // 把输入的关键字统一变成大写字母
            editor.replaceRange(
              __Token.string.toUpperCase(),
              {
                line: __Cursor.line,
                ch: __Token.start
              },
              {
                line: __Cursor.line,
                ch: __Token.end
              },
              __Token.string
            )
          }
          if (this.timeout) {
            clearTimeout(this.timeout)
          }
          // this.timeout = setTimeout(() => {
          //   console.log('调用接口')
          // }, 500)
          editor.showHint()
        })
        // 用户实时输入监听
        this.editor.on('cursorActivity', function (editor) {
          const __Cursor = editor.getDoc().getCursor()
          const __Token = editor.getTokenAt(__Cursor)
          const { string } = __Token
          // console.log(__Cursor, __Token, string)
          // console.log(string)
          _this.formatHint(string)

          if (string.charAt(string.length - 1) === '.') {
            const curIndex = __Token.start
            const curLine = _this.editor.getLine(__Cursor.line)
            const key = curLine.slice(curLine.slice(0, curIndex).lastIndexOf(' ') + 1, curIndex) // 点前的关键字
            console.log('keykeykeykeykey', key)
          }
        })
        console.log('让编辑器每次在调用的时候进行自动刷新', this.editor)
        this.editor.on('change', editor => {
          this.editableTabs[ind-1].sql = editor.getValue()
          if (this.$emit) {
            this.$emit('input', this.editableTabs[ind-1].sql)
          }
        })
        // 将所有codemirror存放对象中
        this.codeMirrorObj[sqlId] = this.editor
          console.log(this.codeMirrorObj)
      })
  },
  // 获取关键词/库名/表名接口
    getKeyWordList(data) {
      this.hintOp = []
      api.getKeyWordList(data.datasourceId).then((req) => {
        if(req.data.state) {
          req.data.value.forEach((item) => {
            let obj = {
              text: item.wordname,
              displayText: item.wordname,
              displayInfo: item.wordfrom,
              type: item.wordtype,
              render: this.hintRender,
            }
            this.hintOp.push(obj)
          })
        }
        this.hintOpAll = this.hintOp
      })
    },
  // 过滤hintOption
  formatHint(val) {
    this.hintOp = []
    this.hintOpAll.forEach((item) => {
      //统一变成大写去匹配，否则无法匹配大小写
      if(item.displayText.toUpperCase().indexOf(val.toUpperCase()) == 0) {
        this.hintOp.push(item)
      }
    })
  },
  handleShowHint2(cmInstance, hintOptions) {
      let cursor = cmInstance.getCursor();
      let cursorLine = cmInstance.getLine(cursor.line);
      let end = cursor.ch;
      let start = end;

      let token = cmInstance.getTokenAt(cursor)
      // console.log(cmInstance, cursor, cursorLine, end, token)
      return {
        list: this.hintOp,
        // [{
        //   text: "abcd",
        //   displayText: "abcd",
        //   displayInfo: "提示信息1",
        //   render: this.hintRender
        //       }, {
        //   text: "qwer",
        //   displayText: "qwer",
        //   displayInfo: "提示信息2",
        //   render: this.hintRender
        // }],
        from: {
          ch: token.start, line: cursor.line
        },
        to: {
          ch: token.end, line: cursor.line
        }
      }
    },
    // 提示框显示及样式
    hintRender(element, self, data) {
      let div = document.createElement("div");
      div.setAttribute("class", "autocomplete-div");

      // 添加弹出框表/字段图标
      let divIcon = ''
      if(data.type == 'table') {
        divIcon = document.createElement("i");
        divIcon.setAttribute("class", "el-icon-date");
        divIcon.style.color = 'blue'
      } else if(data.type == 'field') {
        divIcon = document.createElement("i");
        divIcon.setAttribute("class", "el-icon-film");
        divIcon.style.color = 'blue'
      } else {
        divIcon = document.createElement("i");
      }

      let divText = document.createElement("span");
      divText.setAttribute("class", "autocomplete-name");
      divText.innerText = data.displayText;
      divText.style.display = 'inline-block'
      divText.style.overflow = 'hidden'
      divText.style.whiteSpace = 'nowrap'
      divText.style.textOverflow = 'ellipsis'
      divText.style.marginRight = '10px'

      let divInfo = document.createElement("span");
      divInfo.setAttribute("class", "autocomplete-hint");
      divInfo.innerText = data.displayInfo;
      divInfo.style.display = 'inline-block'
      divInfo.style.float = 'right'
      divInfo.style.color = 'gray'
      divInfo.style.maxWidth = '150px'
      divInfo.style.overflow = 'hidden'
      divInfo.style.whiteSpace = 'nowrap'
      divInfo.style.textOverflow = 'ellipsis'

      div.appendChild(divIcon);
      div.appendChild(divText);
      div.appendChild(divInfo);
      element.appendChild(div);
    },
```

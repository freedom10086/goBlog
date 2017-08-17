const icons = [
    ["undo", "撤销(Ctrl+Z)", 'undo'],
    ["redo", "重做(Ctrl+Y)", "redo"],
    ["bold", "粗体", "format_bold"],
    ["del", "删除线", "strikethrough_s"],
    ["italic", "斜体", "format_italic"],
    ["quote", "引用", "format_quote"],
    ["h1", "标题1", "h1"],
    ["h2", "标题2", "h2"],
    ["h3", "标题3", "h3"],
    ["list-ul", "无序列表", "format_list_bulleted"],
    ["list-ol", "有序列表", "format_list_numbered"],
    ["hr", "横线", "remove"],
    ["link", "添加链接", "insert_link"],
    ["image", "添加图片", "insert_photo"],
    ["code", "代码", "code"],
    ["table", "添加表格", "border_all"],
    ["emoji", "表情", "insert_emoticon"],
    ["watch", "关闭实时预览", "visibility_off"],
    ["fullscreen", "进入全屏模式(ESC还原)", "fullscreen"],
    ["preview", "预览", "desktop_windows"],
];

class Editor {
    constructor(editorBox) {
        this.editor = editorBox;
        this.setHandlers();
        this.init();
    }

    init() {
        let _this = this;
        let leftDiv = document.createElement("div");
        leftDiv.setAttribute("id", "editor-editor");
        leftDiv.setAttribute("class", "editor-left");
        //let isHeader = (/h(\d)/.test(name));
        leftDiv.innerHTML = `
        <ul class="editor-toolbar">
        ${icons.map(item => `
            <li data-name="${item[0]}" title="${item[1]}"><i class="icon-md">${item[2]}</i></li>`).join('')}
        </ul>
        <textarea id="editor-input" title="content"></textarea>`;
        this.editor.appendChild(leftDiv);
        let right = document.createElement("div");
        right.setAttribute("id", "markdown-preview");
        right.setAttribute("class", "editor-preview");
        this.editor.appendChild(right);
        this.input = this.editor.querySelector("#editor-input");

        [...this.editor.querySelectorAll(".editor-toolbar li")].forEach(item => {
            let name = item.getAttribute("data-name");
            console.log(name);
            if (name == 'emoji') {//处理表情事件
                SmileyBox.initSmiley(item, (src, title) => {
                    let sel = _this.getSelection();
                    let str = "![" + title + "](" + src + ")\n";
                    if (!_this.isLineStart()) {
                        str = '\n' + str;
                    }
                    _this.replaceText(str);
                    _this.setSelection(sel.start + str.length);
                    _this.input.focus();
                });
            } else {
                let h = this.handlers[name];
                if (h) {
                    item.addEventListener("click", function () {
                        console.log("click", _this.handlers[name]);
                        _this.handlers[name]();
                        if (name !== "watch" && name !== "preview" && name !== "fullscreen") {
                            _this.input.focus();
                        }
                    })
                }
            }
        });
    }

    //获取选择域位置，如果未选择便是光标位置
    getSelection() {
        let start = this.input.selectionStart;
        let len = this.input.selectionEnd - start;
        return {
            start: start,
            end: this.input.selectionEnd,
            length: len,
            text: this.input.value.substr(start, len)
        };
    }

    //替换选择
    replaceText(text, start = this.input.selectionStart, end = this.input.selectionEnd) {
        console.log("replace>>", start, end, text);
        this.input.value = this.input.value.substr(0, start) +
            text + this.input.value.substr(end, this.input.value.length);
    }

    //设置选择
    setSelection(start, end = start) {
        this.input.setSelectionRange(start, end);
    }

    //是否是一行的开始
    isLineStart() {
        let start = getSelection().start;
        if (start <= 0) {
            return true;
        } else if (this.input.value.charAt(start - 1) == '\n') {
            return true;
        }
        return false;
    }

    //获得一行数据
    getLineText() {
        let sel = this.getSelection();
        let start = sel.start, end = sel.end;
        if (sel.start <= 0) {
            start = 0;
        } else {
            while (start > 0) {
                if (this.input.value.charAt(start - 1) == '\n') {
                    break;
                }
                start--
            }
            if (start < 0) {
                start = 0;
            }
        }

        while (end < this.input.value.length) {
            console.log(this.input.value.charAt(end));
            if (this.input.value.charAt(end) == '\n') {
                break;
            }
            end++
        }
        console.log("start:", start, "end:", end);
        return {
            start: start,
            end: end,
            length: end - start,
            text: this.input.value.substr(start, end - start),
            position: sel.start
        };
    }

    //处理要在一行开头的位置插入的元素
    //如 h1 h2
    insertHeadLike(tagname) {
        let line = this.getLineText();
        let startText = tagname + ' ';
        //判断是不是要取消或者加上
        console.log(">>" + line.text + "<<", line.start, line.end);
        if (line.text.startsWith(startText)) {
            //取消
            this.replaceText(line.text.substr(startText.length), line.start, line.end);
            this.setSelection(line.end - startText.length);
        } else {
            this.replaceText(startText + line.text, line.start, line.end);
            this.setSelection(line.end + startText.length);
        }
    }

    //处理 **?** __?__ ~~?~~等
    insertBoldLike(tagname, isNewLineStart = false, isNewLineEnd = false) {
        let sel = this.getSelection();
        if (sel.length >= tagname.length * 2) {
            //已经选中要取消
            if (sel.text.startsWith(tagname) && sel.text.endsWith(tagname)) {
                this.replaceText(sel.text.substr(tagname.length, sel.length - tagname.length * 2));
                this.setSelection(sel.end - tagname.length * 2);
                return
            }
        }

        let before = sel.start <= 0 ? '' : this.input.value.charAt(sel.start - 1);
        let startText = tagname;
        if (isNewLineStart && before != '\n' && before != '') {
            startText = '\n' + tagname;
        } else if (before != ' ' && before != '\n' && before != '') {
            startText = ' ' + tagname;
        }
        let endText = tagname;
        if (isNewLineEnd) {
            endText += '\n';
        }
        // **** or **+??+** or \n+**+??+**
        this.replaceText(startText + sel.text + endText);
        this.setSelection(sel.start + startText.length + sel.length);
    }

    //处理 ------
    insertHrlike(tagname, isNewLineStart = false, isNewLineEnd = false) {
        let sel = this.getSelection();
        let before = sel.start <= 0 ? '' : this.input.value.charAt(sel.start - 1);
        let startText = tagname;
        if (isNewLineStart && before != '\n' && before != '') {
            startText = '\n' + tagname;
        }

        let endText = '';
        if (isNewLineEnd) {
            endText = '\n';
        }

        this.replaceText(startText + endText);
        this.setSelection(sel.start + startText.length + endText.length);
    }

    //处理ul ol 引用
    //1. 2. - - >
    insertUlLike(tagname) {
        let isOl = (tagname == '1.');
        let sel = this.getSelection();
        let line = this.getLineText();

        let index = 1;
        if (isOl) {
            let content = this.input.value.substr(0, line.start - 1);
            let i = content.lastIndexOf('\n');
            if (i >= 0 && i < this.input.value.length - 3) {
                let res = this.input.value.substr(i + 1).match(/^(\d+)\.\s+/);
                console.log(res);
                if (res != null) {
                    index = Number.parseInt(res[1]) + 1;
                }
            }
        }
        if (sel.text.length == 0 || !sel.text.includes('\n')) {//没有选择多行
            let ins = isOl ? index + '.' : tagname;
            this.insertHeadLike(ins);
        } else {
            //选择了多行
            let sels = sel.text.split("\n");
            let i = 0, len = sels.length;
            for (; i < len && sels[i] != ''; i++) {
                if (sels[i].startsWith('- ')) {
                    sels[i] = sels[i].substr(2)
                } else if (isOl && sels[i].match(/(\d+)\.\s+/) != null) {
                    sels[i] = sels[i].substr(sels[i].match(/(\d+)\.\s+/)[0].length)
                } else {
                    sels[i] = (isOl ? (index + i) + '. ' : tagname + ' ') + sels[i];
                }
            }
            this.replaceText(sels.join('\n'));
        }
    }

    setHandlers() {
        const _this = this;
        this.handlers = {
            undo: function () {
                document.execCommand("undo", false, null);
            },
            redo: function () {
                document.execCommand("redo", false, null);
            },
            bold: function () {
                _this.insertBoldLike('**');
            },
            del: function () {
                _this.insertBoldLike('~~');
            },
            italic: function () {
                _this.insertBoldLike('*');
            },
            quote: function () {
                _this.insertHeadLike('>');
            },

            h1: function () {
                _this.insertHeadLike('#');
            },

            h2: function () {
                _this.insertHeadLike('##');
            },

            h3: function () {
                _this.insertHeadLike('###');
            },

            "list-ul": function () {
                _this.insertUlLike('-')
            },

            "list-ol": function () {
                _this.insertUlLike('1.')
            },

            hr: function () {
                _this.insertHrlike('\n------------', true, true);
            },

            link: function () {
                let sel = _this.getSelection();
                let link = "http://链接地址";
                if (sel.text.startsWith("http")) {
                    link = sel.text;
                }
                let str = "[这儿是链接文字](" + link + ")";
                _this.replaceText(str);
                _this.setSelection(sel.start + str.length)
            },

            image: function () {
                let sel = _this.getSelection();
                let title = "这儿是图片描述";
                let link = "http://图片地址";
                let str = "![" + title + "](" + link + ")";
                _this.replaceText(str);
                _this.setSelection(sel.start + str.length)
            },

            code: function () {
                _this.insertBoldLike('`');
            },

            "code-block": function () {
                _this.insertBoldLike('```', true, true);
            },

            watch: function () {
                let node = _this.editor.querySelector("li[data-name=watch]");
                if (node.getAttribute("title") == "关闭实时预览") {
                    node.innerHTML = `<i class="icon-md">visibility_off</i>`;
                    node.setAttribute("title", '开启实时预览');
                    _this.editor.querySelector("#editor-editor").style.width = '100%';
                    _this.editor.querySelector("#markdown-preview").style.display = 'none';
                } else {
                    node.innerHTML = `<i class="icon-md">visibility</i>`;
                    node.setAttribute("title", '关闭实时预览');
                    _this.editor.querySelector("#editor-editor").style.width = '50%';
                    _this.editor.querySelector("#markdown-preview").style.display = 'block';
                }
            },

            fullscreen: function () {
                let clickEvent = function (event) {
                    if (!event.shiftKey && event.keyCode === 27) {
                        if (_this.editor.getAttribute("class").includes("full-screen")) {
                            _this.fullscreen();
                        }
                    }
                };
                let node = _this.editor.querySelector("li[data-name=fullscreen]");
                //当前是全屏 退出
                if (_this.editor.getAttribute("class").includes("full-screen")) {
                    _this.editor.setAttribute("class", 'editor');
                    node.innerHTML = `<i class="icon-md">fullscreen</i>`;
                    node.setAttribute("title", '进入全屏模式(ESC还原)');
                    window.removeEventListener('keyup', clickEvent);
                } else {
                    //进入全屏
                    _this.editor.setAttribute("class", 'editor full-screen');
                    node.innerHTML = `<i class="icon-md">fullscreen_exit</i>`;
                    node.setAttribute("title", '退出全屏模式');
                    window.addEventListener('keyup', clickEvent)
                }
            },

            preview: function () {
                let __this = this;
                let closebtn = _this.editor.querySelector(".editor-btn-close");
                if (!closebtn) {
                    closebtn = document.createElement("i");
                    closebtn.setAttribute("class", "icon-md editor-btn-close");
                    closebtn.setAttribute("title", "退出预览");
                    closebtn.innerHTML = "cancel";
                    closebtn.addEventListener('click', function () {
                        _this.editor.querySelector("#editor-editor").style.display = 'flex';
                        _this.editor.querySelector("#markdown-preview").style.width = __this.savedPreviewWidth;
                        _this.editor.querySelector(".editor-btn-close").style.display = 'none';
                    });
                    _this.editor.appendChild(closebtn);
                } else if (closebtn.style.display == 'none') {
                    closebtn.style.display = 'block'
                }

                __this.savedPreviewWidth = _this.editor.querySelector("#markdown-preview").style.width;
                _this.editor.querySelector("#editor-editor").style.display = 'none';
                _this.editor.querySelector("#markdown-preview").style.width = '100%';
            },
            savedPreviewWidth: '50%',
        }
    }
}

//同步滚动
function syncScroll(item1, item2) {
    let scrool = item1.scrollTop;
    let scroolh = item1.scrollHeight;
    let nDivHight = item1.offsetHeight;
    let persent = scrool / (scroolh - nDivHight);

    let scroolh_r = item2.scrollHeight;
    let nDivHight_r = item2.offsetHeight;
    item2.scrollTop = persent * (scroolh_r - nDivHight_r);
}
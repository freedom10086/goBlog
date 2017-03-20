let editor, editorHelper, preview, input;

function initEditor(box) {
    editor = box;
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
    editor.appendChild(leftDiv);
    let right = document.createElement("div");
    right.setAttribute("id", "markdown-preview");
    right.setAttribute("class", "editor-preview");
    editor.appendChild(right);

    input = document.querySelector("#editor-input");
    editorHelper = new Editor(input);

    [...editor.querySelectorAll(".editor-toolbar li")].forEach(item => {
        let name = item.getAttribute("data-name");
        console.log(name);
        let h = format_functions[name];
        if (h) {
            item.addEventListener("click", function () {
                format_functions[name]();
                if (name !== "watch" && name !== "preview" && name !== "fullscreen") {
                    input.focus();
                }
            })
        }
    });
}


class Editor {
    constructor(el) {
        this.el = el;
    }

    //获取选择域位置，如果未选择便是光标位置
    getSelection() {
        let start = this.el.selectionStart;
        let len = this.el.selectionEnd - start;
        return {
            start: start,
            end: this.el.selectionEnd,
            length: len,
            text: this.el.value.substr(start, len)
        };
    }

    //替换选择
    replaceText(text, start = this.el.selectionStart, end = this.el.selectionEnd) {
        console.log("replace>>", start, end, text);
        this.el.value = this.el.value.substr(0, start) + text + this.el.value.substr(end, this.el.value.length);
    }

    //设置选择
    setSelection(start, end = start) {
        this.el.setSelectionRange(start, end);
    }

    //是否是一行的开始
    isLineStart() {
        let start = getSelection().start;
        if (start <= 0) {
            return true;
        } else if (this.el.value.charAt(start - 1) == '\n') {
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
                if (this.el.value.charAt(start - 1) == '\n') {
                    break;
                }
                start--
            }
            if (start < 0) {
                start = 0;
            }
        }

        while (end < this.el.value.length) {
            console.log(this.el.value.charAt(end));
            if (this.el.value.charAt(end) == '\n') {
                break;
            }
            end++
        }
        console.log("start:", start, "end:", end);
        return {
            start: start,
            end: end,
            length: end - start,
            text: this.el.value.substr(start, end - start),
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

        let before = sel.start <= 0 ? '' : this.el.value.charAt(sel.start - 1);
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
        editorHelper.replaceText(startText + sel.text + endText);
        editorHelper.setSelection(sel.start + startText.length + sel.length);
    }

    //处理 ------
    insertHrlike(tagname, isNewLineStart = false, isNewLineEnd = false) {
        let sel = this.getSelection();
        let before = sel.start <= 0 ? '' : this.el.value.charAt(sel.start - 1);
        let startText = tagname;
        if (isNewLineStart && before != '\n' && before != '') {
            startText = '\n' + tagname;
        }

        let endText = '';
        if (isNewLineEnd) {
            endText = '\n';
        }

        editorHelper.replaceText(startText + endText);
        editorHelper.setSelection(sel.start + startText.length + endText.length);
    }

    //处理ul ol 引用
    //1. 2. - - >
    insertUlLike(tagname) {
        let isOl = (tagname == '1.');
        let sel = this.getSelection();
        let line = this.getLineText();

        let index = 1;
        if (isOl) {
            let content = this.el.value.substr(0, line.start - 1);
            let i = content.lastIndexOf('\n');
            if (i >= 0 && i < this.el.value.length - 3) {
                let res = this.el.value.substr(i + 1).match(/^(\d+)\.\s+/);
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
}


const format_functions = {
    redo: function () {
        document.execCommand("redo", false, null);
    },
    bold: function () {
        editorHelper.insertBoldLike('**');
    },
    del: function () {
        editorHelper.insertBoldLike('~~');
    },
    italic: function () {
        editorHelper.insertBoldLike('*');
    },
    quote: function () {
        editorHelper.insertHeadLike('>');
    },

    h1: function () {
        editorHelper.insertHeadLike('#');
    },

    h2: function () {
        editorHelper.insertHeadLike('##');
    },

    h3: function () {
        editorHelper.insertHeadLike('###');
    },

    "list-ul": function () {
        editorHelper.insertUlLike('-')
    },

    "list-ol": function () {
        editorHelper.insertUlLike('1.')
    },

    hr: function () {
        editorHelper.insertHrlike('------------', true, true);
    },

    link: function () {
        let selection = editorHelper.getSelection();
        let link = "http://链接地址";
        if (selection.text.startsWith("http")) {
            link = selection.text;
        }
        let str = "[这儿是链接文字](" + link + ")";
        editorHelper.replaceText(str);
    },

    image: function () {
        let title = "这儿是图片描述";
        let link = "http://图片地址";
        let str = "![" + title + "](" + link + ")";
        editorHelper.replaceText(str);
    },

    code: function () {
        editorHelper.insertBoldLike('`');
    },

    "code-block": function () {
        editorHelper.insertBoldLike('```', true, true);
    },

    emoji: function () {
        let paddingsize = 8;
        let smileybox = editor.querySelector("smiley_box");
        if (!smileybox) {
            smileybox = document.createElement("div");
            smileybox.setAttribute("id", "smiley_box");
            smileybox.setAttribute("class", 'smiley_box');
            let smileyContent = `<img src="smiley/tb/tb' + i + '.png"/>`;
            smileybox.innerHTML = smileyContent;
            editor.appendChild(smileybox);
        }

        const iconbtn = editor.querySelector("li[data-name=emoji]");
        if (smileybox.style.display == "none") {
            smileybox.style.display = 'block';
            const tops = iconbtn.offset().top + iconbtn.height();
            const lefts = iconbtn.offset().left - smileybox.offsetWidth / 2;
            smileybox.style.left = (lefts + paddingsize * 2) + 'px';
            smileybox.style.top = (tops + paddingsize * 2) + 'px';

            //拖拽
            smileybox.onmousedown = function (ev) {
                const oEvent = ev || event;
                const mx = oEvent.clientX;
                const my = oEvent.clientY;
                const disX = mx - smileybox.offsetLeft;
                const disY = my - smileybox.offsetTop;
                document.onmousemove = function (ev) {
                    const oEvent = ev || event;
                    const mx = oEvent.clientX;
                    const my = oEvent.clientY;
                    let x = mx - disX;
                    let y = my - disY;
                    if (x < 0) {
                        x = 0;
                    } else if (x > document.documentElement.clientWidth - smileybox.offsetWidth) {
                        x = document.documentElement.clientWidth - smileybox.offsetWidth;
                    }
                    if (y < 0) {
                        y = 0;
                    } else if (y > document.documentElement.clientHeight - smileybox.offsetHeight) {
                        y = document.documentElement.clientHeight - smileybox.offsetHeight;
                    }
                    smileybox.style.left = x + 'px';
                    smileybox.style.top = y + 'px';
                };
                document.onmouseup = function () {
                    document.onmousemove = null;
                    document.onmouseup = null;
                };
                return false;
            };

            //失去焦点消失
            input.onmousedown = function () {
                //iconbtn.parent().toggleClass("active");
                smileybox.style.display = 'none';
                input.onmousedown = null;
                document.onmousedown = null;
            }
        } else {
            smileybox.style.display = 'none';
            smileybox.onmousedown = null;
        }
    },

    watch: function () {
        let node = editor.querySelector("li[data-name=watch]");
        if (node.getAttribute("title") == "关闭实时预览") {
            node.innerHTML = `<i class="icon-md">visibility_off</i>`;
            node.setAttribute("title", '开启实时预览');
            editor.querySelector("#editor-editor").style.width = '100%';
            editor.querySelector("#markdown-preview").style.display = 'none';
        } else {
            node.innerHTML = `<i class="icon-md">visibility</i>`;
            node.setAttribute("title", '关闭实时预览');
            editor.querySelector("#editor-editor").style.width = '50%';
            editor.querySelector("#markdown-preview").style.display = 'block';
        }
    },

    fullscreen: function () {
        let _this = this;
        let clickEvent = function (event) {
            if (!event.shiftKey && event.keyCode === 27) {
                if (editor.getAttribute("class").includes("full-screen")) {
                    _this.fullscreen();
                }
            }
        };
        let node = editor.querySelector("li[data-name=fullscreen]");
        //当前是全屏 退出
        if (editor.getAttribute("class").includes("full-screen")) {
            editor.setAttribute("class", 'editor');
            node.innerHTML = `<i class="icon-md">fullscreen</i>`;
            node.setAttribute("title", '进入全屏模式(ESC还原)');
            window.removeEventListener('keyup', clickEvent);
        } else {
            //进入全屏
            editor.setAttribute("class", 'editor full-screen');
            node.innerHTML = `<i class="icon-md">fullscreen_exit</i>`;
            node.setAttribute("title", '退出全屏模式');
            window.addEventListener('keyup', clickEvent)
        }
    },

    preview: function () {
        let closebtn = editor.querySelector(".editor-btn-close");
        if (!closebtn) {
            let _this = this;
            closebtn = document.createElement("i");
            closebtn.setAttribute("class", "icon-md editor-btn-close");
            closebtn.setAttribute("title", "退出预览");
            closebtn.innerHTML = "cancel";
            closebtn.addEventListener('click', function () {
                editor.querySelector("#editor-editor").style.display = 'flex';
                editor.querySelector("#markdown-preview").style.width = _this.savedPreviewWidth;
                editor.querySelector(".editor-btn-close").style.display = 'none';
            });
            editor.appendChild(closebtn);
        } else if (closebtn.style.display == 'none') {
            closebtn.style.display = 'block'
        }

        this.savedPreviewWidth = editor.querySelector("#markdown-preview").style.width;
        editor.querySelector("#editor-editor").style.display = 'none';
        editor.querySelector("#markdown-preview").style.width = '100%';
    },
    savedPreviewWidth: '50%',
};

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

function insertSmiley(name) {
    let currurl = window.location.href;
    currurl = currurl.substring(0, currurl.lastIndexOf("/"));
    let imgurl = currurl + "/smiley/tb/" + name + ".png";
    let str = "![表情](" + imgurl + ")";
    let nextline = "";
    if (!inputaera.isLineStart()) {
        nextline = "\n";
    }
    inputaera.insertAtCousor(nextline + str);

}
//表情
const SmileyBox = {
    smiley: null,
    callback: null,//表情点击callback('url','name');
    smileys: [
        [0, "tb", '贴吧', 73, ".png", ['惊哭', '呼~', '喷', '生气', '不高兴', '惊讶', '酷', '疑问', '吐舌', '泪', '哈哈', '怒', '鄙视', '勉强', '花心',
            '啊', '乖', '钱', '委屈', '黑线', '睡觉', '呵呵', '太开心', '真棒', '咦', '汗', '阴险', '狂汗', '滑稽', '开心', '笑眼', '吐', '冷']],
        [1, "xds", '小电视', 15, ".webp", ['笑', '发愁', '赞', '差评', '嘟嘴', '汗', '害羞', '吃惊', '哭泣', '太喜欢', '好怒啊', '困惑', '我好兴奋', '思索', '无语']],
        [2, "bz", '2233娘', 15, ".webp", ['大笑', '吃惊', '大哭', '耶', '卖萌', '疑问', '汗', '困惑', '怒', '委屈', '郁闷', '第一', '喝茶', '吐魂', '无言']],
        [3, "qy", '蛆音娘', 20, ".webp", ['卖萌', '吃瓜群众', '吃惊', '害怕', '扶额', '滑稽', '哼', '机智', '哭泣', '睡觉觉', '生气', '偷看', '吐血', '无语', '摇头', '疑问', 'die', 'OK', '肥皂', '大笑']]
    ],
    initSmiley: function (button, callback) {
        if (this.smiley) {
            return
        }
        this.callback = callback;
        console.log("callback", callback);
        let _this = this;
        console.log("init");
        this.btn = button;

        button.addEventListener('click', () => {
            _this.smiley = document.getElementById("smiley");
            if (!_this.smiley) {
                SmileyBox.createSmiley();
            } else {
                SmileyBox.switchState();
            }
        });
    },
    createSmiley: function () {
        this.smiley = document.createElement("div");
        this.smiley.id = "smiley";
        this.smiley.className = "smiley";
        this.smiley.innerHTML = `
                <div class="smiley-tab-box">
                    ${this.smileys.map(item => `
                    <div class = "smiley-tab">
                    <img title="${item[2]}" src="smiley/${item[1]}/${item[1]}${2}${item[4]}">
                    </div>`).join('')}
                    <button type="button" class="close">
                        <span>×</span>
                    </button>
                </div>
        <div class="smiley-container">

        </div>`;

        document.body.appendChild(this.smiley);
        this.smiley.style.top = this.btn.offsetTop + 32 + 'px';
        this.smiley.style.left = this.btn.offsetLeft - 50 + 'px';
        this.smiley.style.display = "block";

        [...this.smiley.querySelectorAll(".smiley-tab")].forEach((v, k) => {
            v.addEventListener("click", function () {
                SmileyBox.switchSmiley(k);
            })
        });

        this.switchSmiley(0);
        this.smiley.querySelector(".close").addEventListener('click', function () {
            SmileyBox.smiley.style.display = "none";
        });
        this.setDrag(true);
    },
    switchSmiley: function (index) {
        let smiley = this.smiley;
        let lis = [...smiley.querySelectorAll(".smiley-tab")];
        if (lis[index].className.includes("on")) {
            //do nothing
            return
        }
        let curr = smiley.querySelector(".smiley-tab.on");
        if (curr) {
            curr.className = "smiley-tab"
        }
        lis[index].className = "smiley-tab on";
        let smileyContainer = smiley.querySelector(".smiley-container");
        let str = "";
        for (let i = 1; i <= this.smileys[index][3]; i++) {
            str += `<img title="${this.smileys[index][5][i - 1]}" src="smiley/${this.smileys[index][1]}/${this.smileys[index][1] + i}${this.smileys[index][4]}"/>`;
        }
        smileyContainer.innerHTML = str;
        let _this = this;
        [...smileyContainer.querySelectorAll("img")].forEach((v, k) => {
            v.addEventListener('click', () => {
                console.log("smiley click", k);
                if (_this.callback) {
                    //表情被点击
                    console.log("smiley click", 'callback');
                    _this.callback(v.getAttribute("src"), v.getAttribute("title"));
                }
            });
        });
    },
    switchState: function (callback) {
        console.log("switch state");
        let state = 0;
        let smileybox = this.smiley;
        if (smileybox.style.display === "none") {
            smileybox.style.display = "block";
            state = 1;
            //拖拽
            this.setDrag(true);
        } else {
            smileybox.style.display = "none";
            state = 0;
            this.setDrag(false);
        }
        if (callback) {
            callback(state);
        }
    },

    setDrag: function (state = false) {
        if (state) {
            this.smiley.querySelector(".smiley-tab-box").onmousedown = function (ev) {
                let oEvent = ev || event;
                let mx = oEvent.clientX;
                let my = oEvent.clientY;
                let disX = mx - SmileyBox.smiley.offsetLeft;
                let disY = my - SmileyBox.smiley.offsetTop;
                document.onmousemove = function (ev) {
                    let oEvent = ev || event;
                    let mx = oEvent.clientX;
                    let my = oEvent.clientY;
                    let x = mx - disX;
                    let y = my - disY;
                    if (x < 0) {
                        x = 0;
                    } else if (x > document.documentElement.clientWidth - SmileyBox.smiley.offsetWidth) {
                        x = document.documentElement.clientWidth - SmileyBox.smiley.offsetWidth;
                    }
                    if (y < 0) {
                        y = 0;
                    } else if (y > document.documentElement.clientHeight - SmileyBox.smiley.offsetHeight) {
                        y = document.documentElement.clientHeight - SmileyBox.smiley.offsetHeight;
                    }
                    SmileyBox.smiley.style.left = x + 'px';
                    SmileyBox.smiley.style.top = y + 'px';
                };
                document.onmouseup = function () {
                    document.onmousemove = null;
                    document.onmouseup = null;
                };
                return false;
            };
        } else {
            this.smiley.onmousedown = null;
        }
    }
};
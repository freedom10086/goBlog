//返回顶部
function backTop(acceleration, time) {
    acceleration = acceleration || 0.1;
    time = time || 13;
    let x1 = 0;
    let y1 = 0;
    let x2 = 0;
    let y2 = 0;

    if (document.documentElement) {
        x1 = document.documentElement.scrollLeft || 0;
        y1 = document.documentElement.scrollTop || 0;
    }
    if (document.body) {
        x2 = document.body.scrollLeft || 0;
        y2 = document.body.scrollTop || 0;
    }
    let x3 = window.scrollX || 0;
    let y3 = window.scrollY || 0;

    // 滚动条到页面顶部的水平距离
    const x = Math.max(x1, Math.max(x2, x3));
    // 滚动条到页面顶部的垂直距离
    const y = Math.max(y1, Math.max(y2, y3));
    const speed = 1 + acceleration;
    window.scrollTo(Math.floor(x / speed), Math.floor(y / speed));

    if (x > 0 || y > 0) {
        const invokeFunction = "goTop(" + acceleration + ", " + time + ")";
        window.setTimeout(invokeFunction, time);
    }
}

//设置用户信息登陆以及未登陆状态
function setProfileView() {
    if (window.location.pathname === "/login" || window.location.pathname === "/register") return;
    const loginView = document.querySelector("#nav-login-view");
    const profileView = document.querySelector("#nav-profile-view");
    if (token !== null && profile !== null) {
        if (loginView) loginView.style.display = "none";
        if (profileView) {
            profileView.querySelector("#nav-profile-username").innerHTML = profile.username;
            profileView.style.display = "block";
        }
    } else {
        profile = null;
        if (loginView) loginView.style.display = "block";
        if (profileView) profileView.style.display = "none";
    }
}

//退出登陆
function exitLogin() {
    Modal.confirm({
        "title": "退出登陆", "content": "你要确认退出登陆吗?", "btnPrimary": {"css": "btn-danger"},
        "btnCancle": {"text": "取消"}
    }, () => {
        localStorage.clear();
        location.reload();
        Toast.show("退出登陆成功");
    });
}

//验证么创造/验证类
class Yzm {
    constructor() {
        this.code = [];
    }

    //创建验证码
    createYZM(canvas) {
        let code = this.code;
        let _this = this;
        let i;
        code.length = 0;
        const context = canvas.getContext("2d");
        context.clearRect(0, 0, canvas.width, canvas.height);
        const random = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
            'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'];
        const colors = ["red", "green", "brown", "blue", "orange", "purple", "black"];
        for (i = 0; i < 4; i++) {
            const index = Math.floor(Math.random() * 36);
            code.push(random[index]);
        }
        context.beginPath();
        // Sprinkle in some random dots
        for (i = 0; i < 10; i++) {
            let px = Math.floor(Math.random() * canvas.width);
            let py = Math.floor(Math.random() * canvas.height);
            context.moveTo(px, py);
            context.lineTo(px + 1, py + 1);
            context.strokeStyle = colors[Math.floor(Math.random() * colors.length)];
            context.lineWidth = Math.floor(Math.random() * 2);
            context.stroke();
        }

        for (i = 0; i < 2; i++) {
            //随机线条
            context.moveTo(0, Math.floor(Math.random() * canvas.height));//随机线的起点x坐标是画布x坐标0位置，y坐标是画布高度的随机数
            context.lineTo(canvas.width, Math.floor(Math.random() * canvas.height));//随机线的终点x坐标是画布宽度，y坐标是画布高度的随机数
            context.lineWidth = 0.3;//随机线宽
            context.strokeStyle = colors[Math.floor(Math.random() * colors.length)];
            context.stroke();//描边，即起点描到终点
        }

        let deg, cos, sin, dg;
        context.font = "20px Arial";
        let cx = (canvas.width - 30) / 3;
        for (i = 0; i < 4; i++) {
            context.fillStyle = colors[Math.floor(Math.random() * colors.length)];
            //产生一个正负30度以内的角度值以及一个用于变形的dg值
            dg = Math.random() * 4.5 / 10;
            deg = Math.floor(Math.random() * 60);
            deg = deg > 30 ? (30 - deg) : deg;
            cos = Math.cos(deg * Math.PI / 180);
            sin = Math.sin(deg * Math.PI / 180);

            context.save();
            context.setTransform(cos, sin + dg, -sin + dg, cos, cx * (i + 1) - 12, 18);
            context.fillText(code[i], 0, 0);
            context.restore();
        }

        canvas.onclick = function () {
            context.clearRect(0, 0, canvas.width, canvas.height);
            _this.createYZM(canvas);
        }
    }

    //验证验证码
    checkValid(valid) {
        console.log(valid, this.code);
        if (!valid || valid.length !== 4) return false;
        for (let i = 0; i < 4; i++) {
            if (valid[i].toUpperCase() !== this.code[i]) {
                break;
            }

            if (i === 3) {
                return true
            }
        }
        return false;
    }
}

//模态确认框
class Modal {
    bindConfirm(btn, options) {
        if (!btn) return;
        btn.addEventListener('click', () => {
            Modal.confirm(options);
        })
    }

    // callback 如果没有返回值(或者返回值为true)，则直接拿关闭,为false取消关闭，为非空字符串，则为错误提示。
    // inputs 是否有输入框{"type"输入类型,"hint","value"默认值}
    static create(title, content, btnPrimary, btnCancel, callback, input) {
        btnPrimary = btnPrimary || {};
        btnCancel = btnCancel || {};
        let modal = document.querySelector(`.modal`);
        let btnPrimaryText = btnPrimary && btnPrimary.text || "确认";
        let btnPrimaryCss = btnPrimary && btnPrimary.css || "btn-primary";
        let btnCancelText = btnCancel && btnCancel.text || "关闭";
        let btnCancelCss = btnCancel && btnCancel.css || "btn-info";
        let inputContent = typeof (input) !== "undefined" ? (`<input type="${input.type || "text"}"
        class="form-control" placeholder="${input.hint || ""}" value="${input.value || ""}"/>`) : "";
        if (!modal) {
            modal = document.createElement("div");
            modal.className = "modal";
            modal.innerHTML = `
            <div class="modal-dialog">
            <div class="modal-content">
            <form>
            <div class="modal-header">
            <h5 class="modal-title">${title}</h5>
            <button type="button" class="close" data-type="close">
            <span aria-hidden="true">&times;</span>
            </button>
            </div>
            <div class="modal-body">
                <p data-type="content">${content}</p>
                <div data-type="input" style="display: ${inputContent === "" ? "none" : "block"}">${inputContent}</div>
            </div>
            <div class="modal-footer">
            <button type="button" data-type="confirm" class="btn ${btnPrimaryCss}">${btnPrimaryText}</button>
            <button type="button" class="btn ${btnCancelCss}" data-type="close">${btnCancelText}</button>
            </div></form></div></div>`;
            this.modal = modal;
            document.body.appendChild(modal);
            [...modal.querySelectorAll("[data-type=close]")].forEach((v, k) => {
                v.addEventListener('click', () => {
                    modal.style.display = "none";
                });
            });

            let inputBox = modal.querySelector("[data-type=input] input");
            if (inputBox !== null) {
                inputBox.addEventListener("change", function () {
                    if (inputBox.value.length > 0 && inputBox.className.includes("is-invalid")) {
                        inputBox.classList.remove("is-invalid");
                        inputBox.setCustomValidity("");
                    }
                });
            }

            modal.querySelector("[data-type=confirm]").addEventListener('click', () => {
                let content = "";
                //有输入内容
                if (modal.querySelector("[data-type=input]").style.display !== 'none') {
                    content = inputBox.value;
                    if (typeof (content) === "undefined" || content.length === 0) {
                        if (!inputBox.classList.contains("is-invalid")) {
                            inputBox.classList.add("is-invalid");
                            inputBox.checkValidity();
                        }
                        return
                    }

                    if (typeof (callback) !== "undefined" && callback !== null) {
                        let res = callback(content);
                        if (res !== null && res !== true) {
                            inputBox.classList.add("is-invalid");
                            inputBox.setCustomValidity(res === false ? "输入不合法" : res);
                            return;
                        }
                    }

                    modal.style.display = "none";
                }
            })
        } else {
            modal.querySelector("p[data-type=content]").innerHTML = content;
            modal.querySelector(".modal-title").innerHTML = title;
            if (inputContent !== "") {
                let input = modal.querySelector("[data-type=input]");
                input.style.display = "block";
                input.innerHTML = inputContent;
            }
        }
        return modal;
    }

    //确认对话框 title, content, btnPrimary, btnCancle, callback
    static confirm(options, callback) {
        options = options || {};
        let modal = document.querySelector(`.modal`);
        let title = options.title || "提示";
        if (!modal) {
            modal = Modal.create(title, options.content || "", options.btnPrimary || {"text": "确认"},
                options.btnCancle || {"text": "关闭"}, callback);
        } else {
            modal.querySelector(".modal-title").innerHTML = title;
            modal.querySelector("p[data-type=content]").innerHTML = options.content || "";
        }

        modal.querySelector("[data-type=input]").style.display = 'none';
        if (modal && modal.style.display !== "block") {
            modal.querySelector(".modal-dialog").className = "modal-dialog slide-down";
            modal.style.display = "block";
        }
    }

    //填写对话框 //{title, callback, input{type,value,hint}}
    static promote(options, callback) {
        options = options || {};
        let title = options.title || "提示";
        let input = options.input || {type: "text"};
        let modal = document.querySelector(`.modal`);
        if (!modal) {
            modal = Modal.create(title, "", options.btnPrimary || {"text": "提交"},
                options.btnCancle || {"text": "取消"}, callback, options.input);
        } else {
            modal.querySelector(".modal-title").innerHTML = title;
            modal.querySelector("p[data-type=content]").innerHTML = "";
            modal.querySelector("[data-type=input]").innerHTML = `<input type="${input.type || "text"}" 
                 class="form-control" placeholder="${input.hint || ""}" value="${input.value || ""}"/>`;
        }
        modal.querySelector("[data-type=input]").style.display = 'block';
        if (modal && modal.style.display !== "block") {
            modal.querySelector(".modal-dialog").className = "modal-dialog slide-down";
            modal.style.display = "block";
        }
    }

    static hide() {
        let modal = document.querySelector(`.modal`);
        if (modal && modal.style.display !== "none") {
            modal.style.display = "none";
        }
    }
}

//用户卡片
const UserCard = {
    card: null,
    tmpl: `
    <img class="card-img-top" src="images/card_img.webp" alt="Card image cap">
    <div class="card-body">
        <div class="row">
            <a target="_blank" href="#" class="face"><img
                    src="images/avater.jpg"></a>
            <div class="user-info">
                <a href="#">悬崖边缘的猫</a>
                <span class="badge badge-info">男</span>
                <span class="badge badge-success">13</span>
            </div>
        </div>
        <p>
            <a href="#" target="_blank">关注: 15</a>
            <a href="#" target="_blank" class="ml-1">粉丝: 999</a>
        </p>
        <p class="sign">冷冷的猫粮在脸上胡乱的拍O__O</p>
        <div class="d-flex justify-content-end">
            <a class="btn btn-primary" href="#">+关注</a>
            <a class="btn btn-info ml-2" href="#" target="_blank">私信</a>
        </div>
    </div>`,
    mousemove: function (e) {
        let x = e.clientX;
        let y = e.clientY;
        UserCard.hide(x, y);
    },
    isin: false,//鼠标是否在元素上
    init: function (as) { //参数标签数组
        [...as].forEach((v, k) => {
            v.addEventListener("mouseover", function (e) {
                //确定要显示的位置
                UserCard.isin = true;
                UserCard.show(e.clientX + 5, e.clientY + 5)
            });

            v.addEventListener("mouseout", function (e) {
                UserCard.isin = false;
                UserCard.show(e.clientX + 5, e.clientY + 5)
            });
        });
    },

    show: function (x, y) {
        if (!UserCard.isin) {
            return
        }
        if (this.card == null) {
            this.card = document.createElement("div");
            this.card.id = "user-card";
            this.card.className = "user-card";
            this.card.innerHTML = this.tmpl;
            document.body.appendChild(this.card)
        }

        if (x && y && this.card.style.display.indexOf("none") == -1) {
            this.card.style.left = x + 'px';
            this.card.style.top = y + 'px';
        }
        document.addEventListener("mousemove", this.mousemove);
        this.card.style.display = "";
    },

    hide: function (x, y) {
        if (this.card == null) {
            return
        }

        if (UserCard.isin) {
            return
        }

        let width = this.card.offsetWidth;
        let height = this.card.offsetHeight;
        //移动到卡片上
        if (x >= this.card.offsetLeft - 15 && x <= this.card.offsetLeft + width + 15
            && y >= this.card.offsetTop - 15 && y <= this.card.offsetTop + height + 15) {
            return
        } else {
            document.removeEventListener('mousemove', this.mousemove);
            this.card.style.display = "none";
        }
    }
};

//toast提示加载
//let toast = new Loading(timeout);
//toast.show();

class Toast {
    static show(text, timeout = 1.5) {
        let toast = document.querySelector("#messageToast");
        if (toast === null) {
            toast = document.createElement("div");
            toast.id = "messageToast";
            toast.innerHTML = `<div class="toast fade-in">${text || "加载中..."}</div>`;
            document.body.appendChild(toast)
        }

        toast.style = "";
        setTimeout(Loading.dismiss, (timeout || 1.5) * 1000);
        return this;
    }

    static dismiss() {
        let toast = document.querySelector("#messageToast");
        if (toast) {
            toast.style.display = "none";
        }
    }
}

class Loading {
    static show(timeout, text) {
        let toast = document.querySelector("#loadingToast");
        if (toast === null) {
            toast = document.createElement("div");
            toast.id = "loadingToast";
            toast.innerHTML = `<div class="loading-toast fade-in"><i class="loading"></i><p class="toast-content">${text || "加载中..."}</p></div>`;
            document.body.appendChild(toast)
        }

        toast.style = "";
        setTimeout(Loading.dismiss, (timeout || 8) * 1000);
        return this;
    }

    static dismiss() {
        let toast = document.querySelector("#loadingToast");
        if (toast) {
            toast.style.display = "none";
        }
    }
}

//下拉框
class DropDown {
    constructor() {
        let drops = document.querySelectorAll(".dropdown");
        [...drops].forEach((v, k) => {
            const toggle = v.querySelector(".dropdown-toggle");
            const dropContent = v.querySelector(".dropdown-menu");
            toggle.addEventListener('click', () => {
                if (dropContent.className.includes("show")) {
                    dropContent.className = "dropdown-menu";
                } else {
                    dropContent.className = "dropdown-menu show";
                }
            });

            [...v.querySelectorAll('.dropdown-item')].forEach((v, k) => {
                v.addEventListener('click', () => {
                    dropContent.className = "dropdown-menu";
                })
            });
        });
    }
}

//可切换的tab
class TabBox {
    constructor(id, callback) { //callback 点击后的回掉 值index
        this.callback = callback;
        this.box = document.getElementById(id);
        if (this.box === null) return;
        this.tabs = this.box.querySelectorAll('.nav-link');
        this.panels = this.box.querySelectorAll('.tab-content');
        for (let i = 0; i < this.tabs.length; i++) {
            const tab = this.tabs[i];
            this.setTabHandler(tab, i);
        }

        this.tabs[0].className = 'nav-link active';
        this.panels[0].className = 'tab-content active';
    }

    setTabHandler(tab, tabPos) {
        let that = this;
        tab.onclick = function () {
            if (that.tabs[tabPos].className.includes("active")) {
                return
            }
            for (let i = 0; i < that.tabs.length; i++) {
                if (i !== tabPos) {
                    that.tabs[i].className = 'nav-link';
                }
            }

            tab.className = 'nav-link active';

            for (let i = 0; i < that.panels.length; i++) {
                if (i !== tabPos) {
                    that.panels[i].className = 'tab-content';
                }
            }

            that.panels[tabPos].className = 'tab-content active';
            if (that.callback) {
                that.callback(tabPos);
            }
        }
    }
}

const ajaxPromise = param => {
    return new Promise((resovle, reject) => {
        var xhr = new XMLHttpRequest();
        xhr.open(param.type || "get", param.url, true);
        xhr.send(param.data || null);

        xhr.onreadystatechange = () => {
            var DONE = 4; // readyState 4 代表已向服务器发送请求
            var OK = 200; // status 200 代表服务器返回成功
            if (xhr.readyState === DONE) {
                if (xhr.status === OK) {
                    resovle(JSON.parse(xhr.responseText));
                } else {
                    reject(JSON.parse(xhr.responseText));
                }
            }
        }
    })
}

// fetch api 返回promise对象
// https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
// fetch text
function fetchText(url, success) {
    const myHeaders = new Headers({
        "Content-Type": "text/plain",
        "Accept": "text/plain",
    });
    const myInit = {
        method: 'GET',
        headers: myHeaders,
        mode: 'cors',
        cache: 'default'
    };

    const myRequest = new Request(url, myInit);
    fetch(myRequest).then(function (res) {
        console.log(res);
        if (res.ok) {
            return res.text();
        }
        throw new Error(res.status);
    }).then(function (text) {
        if (success) {
            success(text);
        }
        console.log(text)
    }).catch(function (error) {
        console.log('There has been a problem with your fetch operation: ' + error.message);
    });
}

//========Ajax================//
class Ajax {
    static send(url, method = 'GET', data, success, fail) {
        const x = new XMLHttpRequest();
        x.open(method, url, /*async*/ true);
        x.onreadystatechange = function () {
            if (x.readyState === 4) {
                let status = x.status;
                if (status >= 200 && status < 300) {
                    success && success(status, x.responseText)
                } else {
                    fail && fail(status, x.responseText);
                }
            }
        };
        if (method === 'POST') {
            x.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        }
        x.send(data)
    }

    static get(url, data, success, fail) {
        let query = [];
        for (let key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }
        url = url + (query.length ? '?' + query.join('&') : '');
        Ajax.send(url, 'GET', null, success, fail)
    }

    static post(url, data, success, fail) {
        let query = [];
        for (let key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }
        Ajax.send(url, 'POST', query.join('&'), success, fail)
    }

    static delete(url, data, success, fail) {
        let query = [];
        for (let key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }
        Ajax.send(url, 'DELETE', query.join('&'), success, fail)
    }
}

//========API================//
class ApiClass {
    constructor() {
        this.version = 1;
    }

    checkEmail(email, result) {
        console.log("check email:", email);
        Ajax.get('/register', {mod: 'checkEmail', email},
            (status, res) => {
                result(true, res, status)
            },
            (status, res) => {
                result(false, res, status)
            });
    }

    checkUsername(username, result) {
        console.log("check username:", username);
        Ajax.get('/register', {mod: 'checkUsername', username},
            function (status, res) {
                result(true, res, status)
            },
            function (status, res) {
                result(false, res, status)
            });
    }

    login(username, password, result) {
        console.log("login:", username);
        Ajax.post("/login", {username, password}, (status, res) => {
            result(true, res, status);
        }, (status, res) => {
            result(false, res, status);
        })
    }

    regist(username, email, result) {
        console.log("register:", username);
        Ajax.post("/register", {username, email}, (status, res) => {
            result(true, res, status);
        }, (status, res) => {
            result(false, res, status);
        })
    }

    registDone(token, password, sex, result) {
        console.log("register2:", token);
        Ajax.post("/users", {token, password, sex}, (status, res) => {
            result(true, res, status);
        }, (status, res) => {
            result(false, res, status);
        })
    }

    forgetPassword(email, result) {
        console.log("/account/forgetpassword");
        Ajax.post("/account/forgetpassword", {email}, (status, res) => {
            result(true, res, status);
        }, (status, res) => {
            result(false, res, status);
        })
    }
}

let Api = new ApiClass();

//========Fetch==============//
class Api2 {
    constructor() {
        this.version = 1;
    }

    async post(url = '', data = {}) {
        const formData = new URLSearchParams();
        for (let key in data) {
            // TODO check if need encodeURIComponent
            formData.append(encodeURIComponent(key), encodeURIComponent(data[key]));
        }
        return fetch(url, {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            mode: 'cors', // no-cors, *cors, same-origin
            cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
            credentials: 'same-origin', // include, *same-origin, omit
            headers: {
                //"Content-Type": "text/plain",
                'Content-Type': 'application/x-www-form-urlencoded',
                //'Content-Type': 'application/json',
                "Accept": "application/json",
            },
            redirect: 'follow', // manual, *follow, error
            referrer: 'no-referrer', // no-referrer, *client
            body: formData // body data type must match "Content-Type" header
            // if Content-Type -> JSON.stringify(data), if
        });
    }

    async login(username, password) {
        console.log("/login:", username);
        const response = await this.post("/login", {username, password});
        if (!response.ok) {
            throw Error(response.status + ":" + response.statusText);
        }
        return response.json();
    }

    async forgetPassword(email) {
        console.log("/account/forgetpassword");
        return await this.post("/account/forgetpassword", {email});
    }
}

let api2 = new Api2();

window.onload = function () {
    console.log("======init js====");
    //读取变量
    profile = JSON.parse(localStorage.getItem("profile"));
    token = localStorage.getItem("token");
    console.log(token, profile);
    //设置状页
    setProfileView();

    new DropDown();
    UserCard.init(document.querySelectorAll("a[href^=users]"));
    document.querySelector(".navbar-toggler").addEventListener('click', function () {
        let content = document.querySelector(".navbar-collapse");
        //collapse navbar-collapse
        if (content.className.includes("collapse ")) {
            content.className = "navbar-collapse";
        } else {
            content.className = 'collapse navbar-collapse';
        }
    });
    if (typeof initPage !== 'undefined') {
        initPage()
    }
};
class Ajax {
    constructor(url, async = true) {
        this.url = url;
        this.async = async;
    }

    send(method = 'GET', data, success, fail) {
        const x = new XMLHttpRequest();
        x.open(method, this.url, this.async);
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

    get(data, success, fail) {
        let query = [];
        for (let key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }
        this.url = this.url + (query.length ? '?' + query.join('&') : '');
        this.send('GET', null, success, fail)
    }

    post(data, success, fail) {
        let query = [];
        for (let key in data) {
            query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
        }
        this.send('POST', query.join('&'), success, fail)
    }
}

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

//promise 版本getJSON
function getJSON(url) {
    return new Promise(function (resolve, reject) {
        const client = new XMLHttpRequest();
        client.open("GET", url);
        client.onreadystatechange = handler;
        client.responseType = "json";
        client.setRequestHeader("Accept", "application/json");
        client.send();

        function handler() {
            if (this.readyState !== 4) {
                return;
            }
            if (this.status >= 200 && this.status < 300) {
                resolve(this.response);
            } else {
                reject(new Error(this.statusText));
            }
        }
    });
}

/*
 fetch api
 https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
 */
function fetchJSON(url, success) {
    const myHeaders = new Headers({
        "Content-Type": "text/plain",
        "Accept": "application/json",
    });
    const myInit = {
        method: 'GET',
        headers: myHeaders,
        mode: 'cors',
        cache: 'default'
    };

    const myRequest = new Request(url, myInit);
    fetch(myRequest).then(function (res) {
        if (res.ok) {
            return res.json();
        }
        throw new Error(res.status);
    }).then(function (json) {
        if (success) {
            success(json);
        }
        console.log(json)
    }).catch(function (error) {
        console.log('There has been a problem with your fetch operation: ' + error.message);
    });
}

//fetch text
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

//模态确认框
class Modal {
    constructor(title, content, callback, btnPrimary, btnCancle) {
        let modal = document.querySelector(`.modal[title=${title}]`);
        let btnPrimaryText = btnPrimary && btnPrimary["text"] || "确认";
        let btnPrimaryCss = btnPrimary && btnPrimary["css"] || "btn-primary";
        let btnCancleText = btnCancle && btnCancle["text"] || "关闭";
        let btnCancleCss = btnCancle && btnCancle["css"] || "btn-info";
        if (!modal) {
            modal = document.createElement("div");
            modal.title = title;
            modal.className = "modal";
            modal.innerHTML = `
            <div class="modal-dialog">
            <div class="modal-content">
            <div class="modal-header">
            <h5 class="modal-title">${title}</h5>
            <button type="button" class="close" data-type="close">
            <span aria-hidden="true">&times;</span>
            </button>
            </div>
            <div class="modal-body">
                <p data-type="content">${content}</p>
            </div>
            <div class="modal-footer">
            <button type="button" class="btn ${btnPrimaryCss}">${btnPrimaryText}</button>
            <button type="button" class="btn ${btnCancleCss}" data-type="close">${btnCancleText}</button>
            </div></div></div>`;

            this.modal = modal;
            document.body.appendChild(modal);

            [...modal.querySelectorAll("[data-type=close]")].forEach((v, k) => {
                v.addEventListener('click', () => {
                    modal.style.display = "none";
                    modal.querySelector(".modal-dialog").className = "modal-dialog";
                });
            });
        } else {
            modal.querySelector("p[data-type=content]").innerHTML = contents;
            modal.querySelector(".modal-title").innerHTML = title;
        }
    }

    bind() {
        if (!this.modal) return;
        btn.addEventListener('click', () => {
            this.show();
        })
    }

    //参数可选
    show(title, content) {
        if (this.modal && this.modal.style.display !== "block") {
            if (title) {
                this.modal.querySelector(".modal-title").innerHTML = title;
            }

            if (content) {
                this.modal.querySelector("p[data-type=content]").innerHTML = content;
            }

            this.modal.querySelector(".modal-dialog").className = "modal-dialog slide-down";
            this.modal.style.display = "block";
        }
    }

    hide() {
        if (this.modal && this.modal.style.display !== "none") {
            this.modal.querySelector(".modal-dialog").className = "modal-dialog";
            this.modal.style.display = "none";
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
class Loading {
    constructor(timeout = 5) {
        this.timeout = timeout;
        this.tmpl = `<div class="toast fade-in">
            <i class="loading"></i>
            <p class="toast-content">加载中...</p>
            </div>`
    }

    show() {
        let toast = document.querySelector("#loadingToast");
        if (toast == null) {
            toast = document.createElement("div");
            toast.id = "loadingToast";
            toast.innerHTML = this.tmpl;
            document.body.appendChild(toast)
        }

        toast.style = "";
        setTimeout(Loading.dismiss, this.timeout * 1000);

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

/*
 //利用对话框返回的值 （true 或者 false）
 if (confirm("确定删除？")) {
 console.log("ok");
 //location.href="http://blog.csdn.net/fengyifei11228/";
 } else {
 }


 var name=prompt("请输入您的名字","");//将输入的内容赋给变量 name ，
 //这里需要注意的是，prompt有两个参数，前面是提示的话，后面是当对话框出来后，在对话框里的默认值
 if(name)//如果返回的有内容
 {alert("欢迎您："+ name)}
 */


//========API================//
const Api = {
    version: '1.0'
};

Api.checkEmail = function (email, result) {
    console.log("check email:", email);
    new Ajax('/register').get(
        {mod: 'checkEmail', email},
        function (status, res) {
            result(true, res)
        },
        function (status, res) {
            result(false, res)
        })
};

Api.checkUsername = function (username, result) {
    console.log("check username:", username);
    new Ajax('/register').get(
        {mod: 'checkUsername', username},
        function (status, res) {
            result(true, res)
        },
        function (status, res) {
            result(false, res)
        })
};
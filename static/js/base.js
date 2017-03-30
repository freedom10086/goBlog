console.log("==in base==");
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

function fetchJSON(url) {
    fetch(url).then((res) => {
        //console.log(res);
        if (res.ok && res.status >= 200 && res.status < 300) {
            return res.json();
        }
        throw new Error(res.status);
    })
        .then((json) => {
            console.log(json);
        })
        .catch(function (error) {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
}

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
            </div>
            </div>`;

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

/*
 usage
 getJSON("/posts.json").then(function (json) {
 console.log('Contents: ' + json);
 }, function (error) {
 console.error('出错了', error);
 });
 */

//api
//function checkEmail(email, success, fail) {
//    new Ajax('/register').get({mod: '', email: email}, success, fail)
//}

/*
 fetch api
 https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
 */
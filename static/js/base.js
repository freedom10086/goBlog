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
            if (x.readyState == 4) {
                let status = x.status;
                if (status >= 200 && status < 300) {
                    success && success(status, x.responseText)
                } else {
                    fail && fail(status, x.responseText);
                }
            }
        };
        if (method == 'POST') {
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
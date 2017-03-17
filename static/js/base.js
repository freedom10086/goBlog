const ajax = {};
ajax.x = function () {
    if (typeof XMLHttpRequest !== 'undefined') {
        return new XMLHttpRequest();
    }
    const versions = [
        "MSXML2.XmlHttp.6.0",
        "MSXML2.XmlHttp.5.0",
        "MSXML2.XmlHttp.4.0",
        "MSXML2.XmlHttp.3.0",
        "MSXML2.XmlHttp.2.0",
        "Microsoft.XmlHttp"
    ];

    let xhr;
    for (let i = 0; i < versions.length; i++) {
        try {
            xhr = new ActiveXObject(versions[i]);
            break;
        } catch (e) {
        }
    }
    return xhr;
};

ajax.send = function (url, method, data, success, fail, async) {
    if (async === undefined) {
        async = true;
    }
    let x = ajax.x();
    x.open(method, url, async);
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
};

ajax.get = function (url, data, success, fail, async) {
    let query = [];
    for (let key in data) {
        query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
    }
    ajax.send(url + (query.length ? '?' + query.join('&') : ''), 'GET', null, success, fail, async)
};

ajax.post = function (url, data, success, fail, async) {
    let query = [];
    for (let key in data) {
        query.push(encodeURIComponent(key) + '=' + encodeURIComponent(data[key]));
    }
    ajax.send(url, 'POST', query.join('&'), success, fail, async)
};

//promise 版本ajax
// Return a new promise.
function ajax2(method, url, data) {
    return new Promise((resolve, reject) => {
        const req = new XMLHttpRequest();
        req.open(method, url);
        req.onload = function () {
            if (req.status >= 200 && req.status < 300) {
                resolve(status, req.response);
            } else {
                reject(Error(req.statusText));
            }
        };

        req.onerror = function (e) {
            reject("NetWork Error!");
        };
        // Make the request
        req.send(data);
    });
}

/*
 usage
 ajax('GET','http://www.baidu.com')
 .then((status,response)=>{
 console.log("Success!",response);
 }).catch((err)=>{
 console.error("Failed!", err);
 });
 */
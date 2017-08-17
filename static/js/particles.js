//仿知乎背景 https://github.com/sunweiling/zhihu-canvas
class Circle {
    constructor(x, y) {
        this.x = x;
        this.y = y;
        this.r = Math.random() * 10 + 2;
        this._mx = Math.random() * 2 - 1;
        this._my = Math.random() * 2 - 1;
    }

    drawCircle(ctx) {
        ctx.beginPath();
        ctx.arc(this.x, this.y, this.r, 0, 360);
        ctx.closePath();
        ctx.fillStyle = 'rgba(204, 204, 204, 0.2)';
        ctx.fill();
    }

    drawLine(ctx, _circle) {
        let dx = this.x - _circle.x;
        let dy = this.y - _circle.y;
        let d = Math.sqrt(dx * dx + dy * dy);
        if (d < 150) {
            ctx.beginPath();
            ctx.moveTo(this.x, this.y);//起始点
            ctx.lineTo(_circle.x, _circle.y);//终点
            ctx.closePath();
            ctx.strokeStyle = 'rgba(204, 204, 204, 0.1)';
            ctx.stroke();
        }
    }

    move(w, h) {
        this._mx = (this.x < w && this.x > 0) ? this._mx : ( -this._mx);
        this._my = (this.y < h && this.y > 0) ? this._my : ( -this._my);
        this.x += this._mx / 2;
        this.y += this._my / 2;
    }
}
class currentCircle extends Circle {
    constructor(x, y) {
        super(x, y);
    }

    drawCircle(ctx) {
        ctx.beginPath();
        this.r = 5;
        ctx.arc(this.x, this.y, this.r, 0, 360);
        ctx.closePath();
        ctx.fillStyle = 'rgba(204, 204, 204, 0.2)';
        ctx.fill();
    }
}
window.requestAnimationFrame = window.requestAnimationFrame
    || window.mozRequestAnimationFrame
    || window.webkitRequestAnimationFrame
    || window.msRequestAnimationFrame;
//=====vars=====
let canvas, ctx, w, h;
let circles = [];
let current_circle = new currentCircle(0, 0);
let init = function (num) {
    for (let i = 0; i < num; i++) {
        circles.push(new Circle(Math.random() * w, Math.random() * h));
    }
    draw();
};
let draw = function () {
    ctx.clearRect(0, 0, w, h);
    for (let i = 0; i < circles.length; i++) {
        circles[i].move(w, h);
        circles[i].drawCircle(ctx);
        for (let j = i + 1; j < circles.length; j++) {
            circles[i].drawLine(ctx, circles[j])
        }
    }

    if (current_circle.x) {
        current_circle.drawCircle(ctx);
        for (let k = 1; k < circles.length; k++) {
            current_circle.drawLine(ctx, circles[k]);
        }
    }
    requestAnimationFrame(draw);
};
window.addEventListener('load', () => {
    //1638 819
    canvas = document.createElement('canvas');
    canvas.style.position = 'absolute';
    canvas.style.width = '100%';
    canvas.style.height = '100%';
    document.body.insertBefore(canvas, document.body.firstChild);
    ctx = canvas.getContext("2d");
    w = canvas.width = canvas.offsetWidth;
    h = canvas.height = canvas.offsetHeight;
    init(w * h / 12000);
});
window.onmousemove = function (e) {
    e = e || window.event;
    current_circle.x = e.clientX;
    current_circle.y = e.clientY;
};
window.onmouseout = function () {
    current_circle.x = null;
    current_circle.y = null;
};
window.onresize = function () {
    w = canvas.width = canvas.offsetWidth;
    h = canvas.height = canvas.offsetHeight;
};
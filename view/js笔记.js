//1 鼠标点击坐标
div.onclick = function(ev){
	var oEvent = ev||event;
	var x = oEvent.clientX;
	var y = oEvent.clientY;
}

//2 设置某div位置
div.style.display = 'block';
div.style.left = x +'px';
div.style.top = y+'px';

document.onclick = function(){

}

//keydown keyup
ev.keyCode 获得keycode 判断按键

//拖拽
window.onload = function(){
	div = xxx
	div.onmousedown = function(ev){
		var mx = oEvent.clientX;
		var my = oEvent.clientY;

		var disX = mx - div.offsetLeft;
		var disY = my - div.offsetTop;

		document.onmousemove = function(ev){
			var mx = oEvent.clientX;
			var my = oEvent.clientY;

			var x = mx-disX;
			var y = my - disY;

			if(x<0){
				x = 0;
			}else if(x>document.documentElement.clientWidth - div.offsetWidth){
				x = document.documentElement.clientWidth - div.offsetWidth;
			}

			if(y<0){
				y = 0;
			}else if(y>document.documentElement.clientHeight - div.offsetHeight){
				y = document.documentElement.clientHeight - div.offsetHeight;
			}


			div.style.left = x+'px';
			div.style.top = y+'px';
		}

		document.onmouseup = function(){
			document.onmousemove = null;
			document.onmouseup = null;
		}

		return false;
	}



    document.ready和onload的区别——JavaScript文档加载完成事件
    页面加载完成有两种事件：

一是ready，表示文档结构已经加载完成（不包含图片等非文字媒体文件）；

二是onload，指示页面包含图片等文件在内的所有元素都加载完成。


	js判断函数参数是否传入
    typeof(variableName) == "undefined"
	
}

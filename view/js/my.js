/**
 * Created by yang on 2016/10/7.
 */

function textchange(){
  var input=document.getElementById("input").value;
  marked(input, function (err, content) {
    if (err) throw err;
    document.getElementById('preview').innerHTML = content;
  });
}

function scrollbar_l(){
  var scrool = document.getElementById('input').scrollTop;
  var scroolh = document.getElementById('input').scrollHeight;
  var nDivHight = document.getElementById('input').offsetHeight;
  var persent = scrool/(scroolh-nDivHight);

  var scroolh_r = document.getElementById('preview').scrollHeight;
  var nDivHight_r = document.getElementById('preview').offsetHeight;

  var scrool_top_r = persent*(scroolh_r-nDivHight_r);
  document.getElementById('preview').scrollTop = scrool_top_r;
}


$.fn.getCursorPosition = function(){
    if(this.lengh == 0) return -1;
    return $(this).getSelectionStart();
  }
$.fn.setCursorPosition = function(position){
  if(this.lengh == 0) return this;
  return $(this).setSelection(position, position);
}
$.fn.getSelection = function(){
  if(this.lengh == 0) return -1;
  var s = $(this).getSelectionStart();
  var e = $(this).getSelectionEnd();
  return this[0].value.substring(s,e);
}
$.fn.getSelectionStart = function(){
  if(this.lengh == 0) return -1;
  input = this[0];

  var pos = input.value.length;

  if (input.createTextRange) {
    var r = document.selection.createRange().duplicate();
    r.moveEnd('character', input.value.length);
    if (r.text == '') 
    pos = input.value.length;
    pos = input.value.lastIndexOf(r.text);
  } else if(typeof(input.selectionStart)!="undefined")
  pos = input.selectionStart;

  return pos;
}
$.fn.getSelectionEnd = function(){
  if(this.lengh == 0) return -1;
  input = this[0];

  var pos = input.value.length;

  if (input.createTextRange) {
    var r = document.selection.createRange().duplicate();
    r.moveStart('character', -input.value.length);
    if (r.text == '') 
    pos = input.value.length;
    pos = input.value.lastIndexOf(r.text);
  } else if(typeof(input.selectionEnd)!="undefined")
  pos = input.selectionEnd;

  return pos;
}
$.fn.setSelection = function(selectionStart, selectionEnd) {
    if(this.lengh == 0) return this;
    input = this[0];

    if (input.createTextRange) {
        var range = input.createTextRange();
        range.collapse(true);
        range.moveEnd('character', selectionEnd);
        range.moveStart('character', selectionStart);
        range.select();
    } else if (input.setSelectionRange) {
        input.focus();
        input.setSelectionRange(selectionStart, selectionEnd);
    }

    return this;
}
$.fn.insertAtCousor = function(myValue){
  var $t=$(this)[0];
  if (document.selection) {
    this.focus();
    sel = document.selection.createRange();
    sel.text = myValue;
    this.focus();
  }else if($t.selectionStart || $t.selectionStart == '0'){
    var startPos = $t.selectionStart;
    var endPos = $t.selectionEnd;
    var scrollTop = $t.scrollTop;
    $t.value = $t.value.substring(0, startPos) + myValue + $t.value.substring(endPos, $t.value.length);
    this.focus();
    $t.selectionStart = startPos + myValue.length;
    $t.selectionEnd = startPos + myValue.length;
    $t.scrollTop = scrollTop;
  }else{
    this.value += myValue;
    this.focus();
  }

  //刷新preview
  textchange();
}

$.fn.isLineStart = function(){
  input = this[0];
  var length = input.value.length;
  var start = $(this).getSelectionStart();

  if(start<=0){
    return true;
  }else if(input.value.charAt(start-1) == '\n'){
    return true;
  }else{
    return false;
  }
}

$.fn.lineStartPos = function(){
  input = this[0];
  var length = input.value.length;
  var start = $(this).getSelectionStart();

  if(start<=0){
    return 0;
  }

  var i=start;
  for(;i>0;i--){
    if(input.value.charAt(i)=='\n'){
      i++;
      break;
    }
  }
  return i;
}

$.fn.lineEndPos = function(){
  if(this.lengh == 0) return -1;
  input = this[0];
  var length = input.value.length;
  var end = $(this).getSelectionEnd();

  var i=end;
  for(;i<length;i++){
    if(input.value.charAt(i)=='\n'){
      break;
    }
  }
  return i;
}

var editor           = $("#myeditor");
var preview          = $("#preview");
var inputholder      = $("#inputholder");
var toolbar          = $("#toolbar");
var inputaera            = $("#input");
var fullscreen = false;
var isOnlyshowInput = false;

var toolbarHandlers = {
  undo : function() {
    document.execCommand("undo", false, null);
  },

  redo : function() {
    document.execCommand("redo", false, null);
  },

  bold : function() {
    var selection = inputaera.getSelection();
    if(selection==""){
      inputaera.insertAtCousor("**** ");
      var pos =  inputaera.getCursorPosition();
      inputaera.setCursorPosition(pos-3);
    }else{
      inputaera.insertAtCousor("**" + selection + "** ");
    }
  },

  del : function() {
    var selection = inputaera.getSelection();
    if(selection==""){
      inputaera.insertAtCousor("~~~~ ");
      var pos =  inputaera.getCursorPosition();
      inputaera.setCursorPosition(pos-3);
    }else{
      inputaera.insertAtCousor("~~" + selection + "~~ ");
    }
  },

  italic : function() {
    var selection = inputaera.getSelection();
    if(selection==""){
      inputaera.insertAtCousor("** ");
      var pos =  inputaera.getCursorPosition();
      inputaera.setCursorPosition(pos-2);
    }else{
      inputaera.insertAtCousor("*" + selection + "* ");
    }
  },

  quote : function() {
    var selection = inputaera.getSelection();
    if (!inputaera.isLineStart()){
      inputaera.insertAtCousor("\n> " + selection);
    }
    else{
      inputaera.insertAtCousor("> " + selection);
    }
  },

  h1 : function() {
    var lstart = inputaera.lineStartPos();
    var lend = inputaera.lineEndPos();
    inputaera.setSelection(lstart,lend);
    var selection = inputaera.getSelection();
    if(selection.indexOf("# ")==0){
      inputaera.insertAtCousor(selection.substring(2));
    }else{
      inputaera.insertAtCousor("# " + selection);
    }
  },

  h2 : function() {
    var lstart = inputaera.lineStartPos();
    var lend = inputaera.lineEndPos();
    inputaera.setSelection(lstart,lend);
    var selection = inputaera.getSelection();
    if(selection.indexOf("## ")==0){
      inputaera.insertAtCousor(selection.substring(3));
    }else{
      inputaera.insertAtCousor("## " + selection);
    }
  },

  h3 : function() {
    var lstart = inputaera.lineStartPos();
    var lend = inputaera.lineEndPos();
    inputaera.setSelection(lstart,lend);
    var selection = inputaera.getSelection();
    if(selection.indexOf("### ")==0){
      inputaera.insertAtCousor(selection.substring(4));
    }else{
      inputaera.insertAtCousor("### " + selection);
    }
  },

  "list-ul" : function() {
    var selection = inputaera.getSelection();
    var nextline = "";
    if(!inputaera.isLineStart()){
      nextline = "\n";
    }
    if (selection === ""){
      inputaera.insertAtCousor(nextline+"- ");
    }else{
      var selectionText = selection.split("\n");
      for (var i = 0, len = selectionText.length; i < len; i++)
      {
        selectionText[i] = (selectionText[i] === "") ? "" : "- " + selectionText[i];
      }
      inputaera.insertAtCousor(nextline+selectionText.join("\n"));
    }
  },

  "list-ol" : function() {
    var selection = inputaera.getSelection();
    var nextline = "";
    if(!inputaera.isLineStart()){
      nextline = "\n";
    }
    if (selection === "")
    {
      inputaera.insertAtCousor(nextline+"1. ");
    }
    else
    {
      var selectionText = selection.split("\n");

      for (var i = 0, len = selectionText.length; i < len; i++)
      {
        selectionText[i] = (selectionText[i] === "") ? "" : (i+1) + ". " + selectionText[i];
      }

      inputaera.insertAtCousor(nextline+selectionText.join("\n"));

    }
  },

  hr : function() {
    var position = inputaera.getCursorPosition();
    var selstart = inputaera.getSelectionStart();
    var selection = inputaera.getSelection();

    inputaera.insertAtCousor((inputaera.isLineStart() ? "\n" : "\n\n") + "------------\n\n");
  },

  link: function (){
    var selection = inputaera.getSelection();

    var title = "这儿是链接文字"
    var link = "http://链接地址"

    if(selection!=""&&selection.indexOf("http://")!=-1){
        link = selection;
    }else if(selection!=""){
        title = selection;
    }
    
    var  str = "[" + title + "](" + link + ")";                             
    inputaera.insertAtCousor(str);
  },

  image : function() {
    
    var title = "这儿是图片描述"
    var link = "http://图片地址"
    
    var  str = "![" + title + "](" + link + ")";                             
    inputaera.insertAtCousor(str);
  },

  code : function() {
    var selection = inputaera.getSelection();
    
    inputaera.insertAtCousor("`" + selection + "`");
    if(selection==""){
      var pos =  inputaera.getCursorPosition();
      inputaera.setCursorPosition(pos-1);
    }
  },

  "code-block" : function() {
    var selection = inputaera.getSelection();
    var nextline = "";
    if(!inputaera.isLineStart()){
      nextline = "\n";
    }
    inputaera.insertAtCousor(nextline+"```\n" + selection + "\n```\n");
    
    if(selection==""){
      var pos =  inputaera.getCursorPosition();
      inputaera.setCursorPosition(pos-5);
    }
  },

  table : function() {
    var nextline = "";
    if(!inputaera.isLineStart()){
      nextline = "\n";
    }
    var str = nextline+"|  项目  |  价格   |  数量  |\n| :------: | :-----: | :-----: |\n| 计算机 | $1600  |   5   |\n|  手机  |  $12   |   12   |"
    inputaera.insertAtCousor(str);
  },

  datetime : function() {
    var addZero = function(d) {
            return (d < 10) ? "0" + d : d;
    };

    var selection = inputaera.getSelection();
    var date      = new Date();
    var weekDay = date.getDay();
    var year    = date.getFullYear();
    var month   = addZero(date.getMonth() + 1);
    var day     = addZero(date.getDate());
    var hour    = addZero(date.getHours());
    var min     = addZero(date.getMinutes());
    var second  = addZero(date.getSeconds());

    var fymd    = year  + "-" + month + "-" + day;
    var hms     = hour  + ":" + min   + ":" + second;
    var cnWeekDays = ["日", "一", "二", "三", "四", "五", "六"];
    var datefmt = "星期" + cnWeekDays[weekDay];
    var str = fymd+" "+hms+" "+datefmt;
    if(!inputaera.isLineStart()){
      str = "\n"+str;
    }
    inputaera.insertAtCousor(str);
  },

  emoji : function() {
    var smiley_container =  $("#smiley_container");
    if(smiley_container.length <= 0){
      var strFace, labFace;
      strFace = '<div id="smiley_container" style="position:absolute;display:none;z-index:999;" class="smiley">' +
              '<table border="0" cellspacing="0" cellpadding="0"><tr>';
      for(var i=1; i<=33; i++){
        labFace = 'tb'+i;
        strFace += '<td><img src="smiley/tieba/tb'+i+'.png" onclick="insertSmiley(\'' + labFace + '\');"/></td>';
        if( i % 9 == 0 ) strFace += '</tr><tr>';
      }

      strFace += '</tr></table></div>';
      editor.append(strFace);
      $("#smiley_container").hide();
    }
    smiley_container =  $("#smiley_container")
    var paddingsize = 8;
    var iconbtn =  toolbar.find(".fa[name=emoji]");

    var tops = iconbtn.offset().top+iconbtn.height();
    var lefts = iconbtn.offset().left - smiley_container.width()/2;
  
    smiley_container.css({
        left :lefts+paddingsize*2,
        top :tops+paddingsize*2,
        position:"absolute"
    });

    iconbtn.parent().toggleClass("active");
    smiley_container.toggle();
  },

  watch : function() {
    showhideInput();
  },

  preview : function() {
    showhidepreview();
  },

  fullscreen : function() {
    //this.fullscreen();
    enterFullscreen();
  },
};


/**
 * 鼠标和触摸事件的判断/选择方法
 * MouseEvent or TouchEvent type switch
 *
 * @param   {String} [mouseEventType="click"]    供选择的鼠标事件
 * @param   {String} [touchEventType="touchend"] 供选择的触摸事件
 * @returns {String} EventType                   返回事件类型名称
 */

function mouseOrTouch(mouseEventType, touchEventType) {
  mouseEventType = mouseEventType || "click";
  touchEventType = touchEventType || "touchend";

  var eventType  = mouseEventType;

  try {
    document.createEvent("TouchEvent");
    eventType = touchEventType;
  } catch(e) {}

  return eventType;
};


/**
 * 工具栏图标事件处理器
 * Bind toolbar icons event handle
 *
 * @returns {editormd}  返回editormd的实例对象
 */
function  setToolbarHandler() {
  console.log("setToolbarHandler");
  var toolbarIcons        = toolbar.find(".editormd-menu > li > a");

  toolbarIcons.bind(mouseOrTouch("click", "touchend"), function(event) {
        var _this               = this;
        var icon                = $(this).children(".fa");
        var name                = icon.attr("name");

        $.proxy(toolbarHandlers[name], _this)(inputaera);

        if (name !== "watch" && name !== "preview"  && name !== "fullscreen"){
            inputaera.focus();
        }

        return false;

    });
};


function setViewSizes(){
    var paddingsize = 8;

    var windowheight = $(window).height();
    var windowwidth = $(window).width();

    if(!fullscreen){
      editor.css({
            width    : $(window).width()*0.9,
            height   : $(window).height()*0.75
        });
    }

    var editorh = editor.height();
    var editorw = editor.width();

    inputholder.css({
        width    : editorw/2,
        height   : editorh
    })

    inputaera.css({
        width    : editorw/2-paddingsize*2,
        height   : editorh-toolbar.height()-paddingsize*2,
        padding  : paddingsize
    })

    preview.css({
        width    : editorw/2-paddingsize*2,
        height   : editorh-paddingsize*2,
        padding  : paddingsize
    })
}

/**
 * 编辑器全屏显示
 * Fullscreen show
 *
 * @returns {editormd}         返回editormd的实例对象
 */

function enterFullscreen(){
    var fullscreenClass  = "editormd-fullscreen";

    toolbar.find(".fa[name=fullscreen]").parent().toggleClass("active");

    var escHandle = function(event) {
        if (!event.shiftKey && event.keyCode === 27)
        {
            if (fullscreen)
            {
                fullscreenExit();
            }
        }
    };

    if (!editor.hasClass(fullscreenClass)){
        fullscreen = true;
        $("html,body").css("overflow", "hidden");
        editor.css({
            width    : $(window).width(),
            height   : $(window).height()
        }).addClass(fullscreenClass);

        $(window).bind("keyup", escHandle);
    }else{
        $(window).unbind("keyup", escHandle);
        fullscreenExit();
    }
    setViewSizes();
}

/**
 * 编辑器退出全屏显示
 * Exit fullscreen state
 *
 * @returns {editormd}         返回editormd的实例对象
 */

function fullscreenExit() {
    var fullscreenClass  = "editormd-fullscreen";
    fullscreen = false;
    toolbar.find(".fa[name=fullscreen]").parent().removeClass("active");
    $("html,body").css("overflow", "");
    editor.css({
        width    : editor.data("oldWidth"),
        height   : editor.data("oldHeight")
    }).removeClass(fullscreenClass);
    setViewSizes();
}


function showhideInput(){
    var watchIcon   = "fa-eye-slash";
    var unWatchIcon = "fa-eye";
    var icon    = toolbar.find(".fa[name=watch]");
    if(isOnlyshowInput){
      //关闭只显示编辑框
      //显示预览
      preview.show();
      isOnlyshowInput = false;
      icon.parent().attr("title", "关闭实时预览");
      icon.removeClass(unWatchIcon).addClass(watchIcon);
      setViewSizes();
    }else{
      //打开只显示编辑框
      //关闭预览
      preview.hide();
      isOnlyshowInput =  true;
      var paddingsize = 8;
      var editorh = editor.height();
      var editorw = editor.width();

      inputholder.css({
          width    : editorw,
          height   : editorh
      })

      inputaera.css({
          width    : editorw-paddingsize*2,
          height   : editorh-toolbar.height()-paddingsize*2,
          padding  : paddingsize
      })

      icon.parent().attr("title", "开启实时预览");
      icon.removeClass(watchIcon).addClass(unWatchIcon);
    }
}

var isOnlyShowPreview = false;
function showhidepreview(){
    var closeBtn =  editor.find(".editormd-preview-close-btn");

    if(isOnlyShowPreview){
      //正常状态
      inputholder.show();
      isOnlyShowPreview = false;
      closeBtn.hide().unbind(mouseOrTouch("click", "touchend"));
      setViewSizes();
    }else{
      //只显示预览状态
      inputholder.hide()
      isOnlyShowPreview =  true;
      var paddingsize = 8;
      var editorh = editor.height();
      var editorw = editor.width();

      preview.css({
          width    : editorw-paddingsize*2,
          height   : editorh-paddingsize*2,
          padding  : paddingsize
      })

     var tops = editor.offset().top;

     var lefts = editor.offset().left;

     closeBtn.css({
        left :lefts+paddingsize*2,
        top :tops+paddingsize*2
     })
     closeBtn.show().bind(
        mouseOrTouch("click", "touchend"), function(){
            showhidepreview();
      });
    }
}

function insertSmiley(name){
  var currurl = window.location.href;
  currurl = currurl.substring(0,currurl.lastIndexOf("/"))
  var imgurl = currurl+"/smiley/tieba/"+name+".png";
  var str = "![表情]("+imgurl+")"
  var nextline = "";
    if(!inputaera.isLineStart()){
      nextline = "\n";
    }
  inputaera.insertAtCousor(nextline+str);
  
}
# 我的博客系统文档


## 1.安装mysql
```shell
sudo apt install mysql-server
CREATE DATABASE MyBlog DEFAULT CHARACTER SET utf8mb4;
```
## 2. 数据表定义
```sql
CREATE TABLE `user` (
  `uid` int NOT NULL AUTO_INCREMENT COMMENT '用户表用户id',
  `username` varchar(25) NOT NULL COMMENT '用户名',
  `password` varchar(50) NOT NULL COMMENT '密码',
  `email` varchar(25) NOT NULL COMMENT '邮件',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0-正常 1-禁止访问',
  `sex` tinyint NOT NULL DEFAULT 0 COMMENT '性别0-未知，1-男，2-女',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '个人描述',
  `exp` int NOT NULL DEFAULT 0 COMMENT '经验值',
  `sites` varchar(50) NOT NULL DEFAULT '' COMMENT '个人网站',
  `birthday` date NOT NULL DEFAULT '0000-00-00' COMMENT '生日',
  `messages` int NOT NULL DEFAULT 0 COMMENT '新消息数目',
  `posts` int NOT NULL DEFAULT 0 COMMENT '发帖数',
  `replys` int NOT NULL DEFAULT 0 COMMENT '回复数',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号码',
  `regtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  PRIMARY KEY (`uid`) ,
  UNIQUE INDEX `user_username` (`username`),
  UNIQUE INDEX `user_email` (`email`)
);

CREATE TABLE `category` (
  `cid` tinyint NOT NULL AUTO_INCREMENT COMMENT '版块表id',
  `title` varchar(25) NOT NULL COMMENT '版块名字',
  `description` varchar(150) NOT NULL DEFAULT '' COMMENT '描述',
  `posts` int NOT NULL DEFAULT 0 COMMENT '版块帖子计数',
  `todayposts` int NOT NULL DEFAULT 0 COMMENT '今日新帖数目',
  `lastpost` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后发表时间',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
  PRIMARY KEY (`cid`) ,
  UNIQUE INDEX `cat_name` (`title`)
);

CREATE TABLE `post` (
  `tid` int NOT NULL AUTO_INCREMENT COMMENT '文章表id',
  `cid` tinyint NOT NULL  COMMENT '版块id',
  `uid` int NOT NULL  COMMENT '用户id',
  `author` varchar(25) NOT NULL DEFAULT '' COMMENT '用户名',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
  `content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
  `type` tinyint NOT NULL DEFAULT 0 COMMENT '类型0-一般，1-管理员加精华',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0-正常，1-不可回复2不可查看',
  `views` int NOT NULL DEFAULT 0 COMMENT '查看数',
  `replys` int NOT NULL DEFAULT 0 COMMENT '回复数',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
  `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '编辑时间',
  `lastreply` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后回复时间',
  PRIMARY KEY (`tid`) ,
  INDEX `post_new` (`cid`,`lastreply`),
  INDEX `post_who` (`uid`),
  FOREIGN KEY (`cid`) REFERENCES `category` (`cid`) ON DELETE CASCADE,
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);


DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '评论表id',
  `tid` int NOT NULL COMMENT '帖子id',
  `pid` int NOT NULL DEFAULT 0 COMMENT '父评论id 0-一般回复0+回复某个回复（楼中楼）',
  `uid` int NOT NULL COMMENT '用户id',
  `tuid` int NOT NULL COMMENT '回复对象uid',
  `author` varchar(25) NOT NULL DEFAULT '' COMMENT '用户名',
  `content` varchar(2000) NOT NULL DEFAULT '' COMMENT '内容',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
  `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '编辑时间',
  `isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0未读1已读',
  `replys` int NOT NULL DEFAULT 0 COMMENT '回复计数大于0表示有楼中楼回复',
  PRIMARY KEY (`id`) ,
  INDEX `comment_tid` (`tid`),
  INDEX `comment_who` (`uid`),
  INDEX `comment_pid` (`pid`),
  FOREIGN KEY (`tid`) REFERENCES `post` (`tid`) ON DELETE CASCADE,
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
  FOREIGN KEY (`tuid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `star` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '收藏表id',
  `uid` int NOT NULL COMMENT '用户id',
  `tid` int NOT NULL COMMENT '帖子id',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '帖子标题',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  PRIMARY KEY (`id`) ,
  UNIQUE INDEX `star_unique` (`uid`, `tid`),
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

DROP TABLE IF EXISTS `atmessage`;
CREATE TABLE `atmessage` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '回复和@表 id',
  `uid` int NOT NULL COMMENT '我的uid',
  `fuid` int NOT NULL COMMENT '来自用户uid',
  `fauthor` varchar(25) NOT NULL DEFAULT '' COMMENT '来自用户名',
  `tid` int NOT NULL COMMENT '来自文章id'
  `cid` int NOT NULL COMMENT '来自文章评论id'
  `title` varchar(50) NOT NULL DEFAULT '',
  `content` varchar(200) NOT NULL DEFAULT '' COMMENT '消息内容',
  `isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0未读1已读',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  PRIMARY KEY (`id`) ,
  INDEX `msg_uid` (`uid`),
  INDEX `msg_from` (`fuid`),
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
  FOREIGN KEY (`fuid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `follow` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '我的uid',
  `tuid` int NOT NULL COMMENT '对方uid',
  `name` varchar(25) NOT NULL COMMENT '我的名字，方便查询关注我的';
  `tname` varchar(25) NOT NULL DEFAULT '' COMMENT '对方用户名',
  `note` varchar(25) NOT NULL DEFAULT '' COMMENT '备注名',
  `relation` tinyint NOT NULL DEFAULT 0 COMMENT '关注表还没想好冗余',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  PRIMARY KEY (`id`) ,
  UNIQUE INDEX `follow_me_t` (`uid`,`tuid`),
  INDEX `friend_rev` (`tuid`),
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
  FOREIGN KEY (`tuid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `chat` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '聊天表id',
  `send` int NOT NULL COMMENT '我的uid',
  `recieve` int NOT NULL COMMENT '来自uid',
  `content` varchar(500) NOT NULL DEFAULT '' COMMENT '内容',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
  `isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0-未读，1-已读，只有recieve读取才置为已读',
  PRIMARY KEY (`id`) ,
  INDEX `chat_my_s` (`send`),
  INDEX `chat_my_r` (`recieve`),
  FOREIGN KEY (`send`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
  FOREIGN KEY (`recieve`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);
```

## 2. 触发器定义
```sql
#发表文章 更新category计数 更新用户经验值+3
CREATE TRIGGER `t_post_add` After INSERT ON `post` FOR EACH ROW BEGIN
  update `user` set `exp` = `exp` + 3,`posts` = `posts` + 1 where `uid` = new.uid;
  update `category` set `posts` = `posts` + 1,`lastpost` = new.created, `todayposts` = `todayposts` + 1 where `cid` = new.cid;
END;

#删除文章 更新category计数 更新用户经验值-3
DROP TRIGGER IF EXISTS `t_post_del`;
CREATE TRIGGER `t_post_del` After DELETE ON `post` FOR EACH ROW BEGIN
  update `category` set `posts` = `posts` -1 ,`todayposts` = `todayposts` - 1 where `cid` = old.cid;
  update `user` set `exp` = `exp` - 3,`posts` = `posts` -1 where `uid` = old.uid;
END;

#增加之前
CREATE TRIGGER `t_comment_add_b` Before INSERT ON `comment` FOR EACH ROW BEGIN
  if new.uid = new.tuid then
    set new.isread = '1';
  end if;
END;

#增加评论
CREATE TRIGGER `t_comment_add_a` After INSERT ON `comment` FOR EACH ROW BEGIN
  update `user` set `exp` = `exp` + 1,`replys` = `replys` + 1 where `uid` = new.uid; #增加经验值
  update `post` set `replys` = `replys` + 1,`lastreply` = new.created where `tid` = new.tid;
  if new.isread = 0 AND new.uid <> new.tuid then
    UPDATE `user` SET `messages` = `messages` + 1 where `uid` = new.tuid;
  end if;
END;

#删除评论
CREATE TRIGGER `t_comment_del` After DELETE ON `comment` FOR EACH ROW BEGIN
  update `user` set `exp` = `exp` - 1,`replys` = `replys` - 1 where `uid` = old.uid;
  update `post` set `replys` = `replys` - 1 where `tid` = old.tid;
  if old.isread = 0 then #这条消息未读 则减一
    update `user` set `messages` = `messages` - 1 where `uid` = old.tuid;
  end if;
END;

#编辑评论
CREATE TRIGGER `t_comment_upd` After UPDATE ON `comment` FOR EACH ROW BEGIN
  if old.isread = 0 AND new.isread = 1 then
    UPDATE `user` SET `messages` = `messages` - 1 where `uid` = old.tuid;
  end if;
END;

#增加@消息之前
CREATE TRIGGER `t_atmsg_add_b` Before INSERT ON `atmessage` FOR EACH ROW BEGIN
  if new.uid = new.fuid then
    set new.isread = '1';
  end if;
END;

#增加@消息
CREATE TRIGGER `t_atmsg_add_a` After INSERT ON `atmessage` FOR EACH ROW BEGIN
  if new.isread = 0 AND new.uid <> new.fuid then
    update `user` set `messages` = `messages` + 1 where uid = new.uid;
  end if;
END;

#@消息已读
CREATE TRIGGER `t_atmsg_upd` After UPDATE ON `atmessage` FOR EACH ROW BEGIN
  if old.isread = 0 AND new.isread = 1 then
    UPDATE `user` SET `messages` = `messages` - 1 WHERE `uid` = old.uid;
  end if;
END;


#增加聊天消息
CREATE TRIGGER `t_chat_add` After INSERT ON `chat` FOR EACH ROW BEGIN
  if new.send <> new.recieve then
    update `user` set `messages` = `messages` + 1 where uid = new.recieve;
  end if;
END;

#聊天消息已读
CREATE TRIGGER `t_chat_read` After UPDATE ON `chat` FOR EACH ROW BEGIN
  if old.isread = 0 AND new.isread = 1 then
    UPDATE `user` SET `messages` = `messages` - 1 WHERE `uid` = new.recieve;
  end if;
END;
```
##3.计划任务定义
```sql
#查询手否开启 show variables like '%event%';
#开启 SET PERSIST GLOBAL event_scheduler = ON; mysql 8 可以持久化SET PERSIST
DROP EVENT IF EXISTS `event_aday`;
CREATE EVENT `event_aday`
ON SCHEDULE EVERY 1 DAY STARTS '2016-10-01 00:00:00'
ON COMPLETION PRESERVE
ENABLE
DO update `category` SET `todayposts` = 0 WHERE 1;
```
##4.函数(存储过程)定义
```sql
#新增分类
CREATE PROCEDURE cate_add(
  IN in_title varchar(25),
  IN in_des varchar(150))
BEGIN
  INSERT INTO `category`(`title`,`description`) VALUES (in_title,in_des);
END;

#编辑分类
CREATE PROCEDURE cate_edit(
  IN in_cid int,
  IN in_title varchar(25),
  IN in_des varchar(150))
BEGIN
  UPDATE `category` SET `title` = in_title,`description` = in_des WHERE `cid` = in_cid;
END;

#删除分类 最好不要
CREATE PROCEDURE cate_del(IN in_cid int)
BEGIN
  DELETE FROM `category` WHERE `cid` = in_cid;
END;


#新增用户
CREATE PROCEDURE user_reg(
  IN in_name varchar(25),
  IN in_pass varchar(50),
  IN in_email varchar(25),
  IN sex tinyint)
BEGIN
  INSERT INTO `user`(`username`,`password`,`email`,`sex`) VALUES (in_name,in_pass,in_email,sex);
END;

#删除用户
CREATE PROCEDURE user_del(IN in_uid int)
BEGIN
  delete from `user` where `uid` = in_uid;
END;

#编辑用户
CREATE PROCEDURE user_edit(
  IN in_uid int,
  IN in_sex tinyint,
  IN in_des varchar(200),
  IN in_sites varchar(50),
  IN in_birth date,
  IN in_phone varchar(20))
BEGIN
  UPDATE `user` SET `sex` = in_sex,`description` = in_des,`sites` = in_sites,`birthday` = in_birth,`phone` = in_phone WHERE `uid` = in_uid;
END;

#禁止用户
CREATE PROCEDURE user_bolck(IN in_uid int)
BEGIN
  UPDATE `user` SET `status` = '1' WHERE `uid` = in_uid;
END;

#允许用户
CREATE PROCEDURE user_open(IN in_uid int)
BEGIN
  UPDATE `user` SET `status` = '0' WHERE `uid` = in_uid;
END;

#修改密码
CREATE PROCEDURE user_changepass(
  IN in_uid int,
  IN in_pass varchar(50))
BEGIN
  UPDATE `user` SET `password` = in_pass WHERE `uid` = in_uid;
END;

#发表文章
DROP PROCEDURE IF EXISTS post_add;
CREATE PROCEDURE post_add(
  IN in_cid int,
  IN in_uid int,
  IN in_title varchar(50),
  IN in_content varchar(5000))
BEGIN
  set @author = (select `username` from `user` where `uid` = in_uid);
  INSERT INTO `post`(`cid`,`uid`,`author`,`title`,`content`) VALUES (in_cid,in_uid,@author,in_title,in_content);
END;

#编辑文章
CREATE PROCEDURE post_edit(
  IN in_tid int,
  IN in_title varchar(50),
  IN in_content varchar(5000))
BEGIN
  set @timenow = now();
  update `post` set `title` = in_title,`content` = in_content,`updated` = @timenow,`lastreply` = @timenow where `tid` = in_tid;
END;

#删除文章
CREATE PROCEDURE post_del(IN in_tid int)
BEGIN
  delete from `post` where `tid` = in_tid;
END;

#禁止回复文章
CREATE PROCEDURE post_close_c(IN in_tid int)
BEGIN
  update `post` set `status` = '1' where `tid` = in_tid;
END;

#允许回复文章
CREATE PROCEDURE post_open_c(IN in_tid int)
BEGIN
  update `post` set `status` = '0' where `tid` = in_tid;
END;

#添加评论 回复楼主
DROP PROCEDURE IF EXISTS `comment_add_lz`;
CREATE PROCEDURE comment_add_lz(
  IN in_tid int,
  IN in_uid int,
  IN in_content varchar(2000))
BEGIN
  set @author = (select `username` from `user` where `uid` = in_uid);
  select `uid`,`title` INTO @tuid,@title from `post` where `tid` = in_tid;
  INSERT INTO `comment`(`tid`,`uid`,`tuid`,`author`,`content`,`isread`) VALUES (in_tid,in_uid,@tuid,@author,in_content,@isread);
END;

#楼中楼回复
DROP PROCEDURE IF EXISTS `comment_add_cz`;
CREATE PROCEDURE comment_add_cz(
  IN in_tid int,
  IN in_pid int,
  IN in_uid int,
  IN in_content varchar(2000))
BEGIN
  select `uid` INTO @czuid from `comment` where `id` = in_pid AND `tid` = in_tid;
  if @czuid is not null then
    set @author = (select `username` from `user` where `uid` = in_uid);
    INSERT INTO `comment`(`tid`,`pid`,`uid`,`tuid`,`author`,`content`) VALUES (in_tid,in_pid,in_uid,@czuid,@author,in_content);
    update `comment` set `replys` = `replys` + 1  where `id` = in_pid;
    set @lastid = LAST_INSERT_ID();
    select `title`,`uid` INTO @title,@lzuid from `post` where `tid` = in_tid;
    if @lzuid <> @czuid THEN #发送@给层主
      INSERT INTO `atmessage`(`uid`,`fuid`,`fauthor`,`tid`,`cid`,`title`,`content`) VALUES (@czuid,in_uid,@author,in_tid,@lastid,@title,in_content);
    end if;
  END IF;
END;

#删除评论
DROP PROCEDURE IF EXISTS `comment_del`;
CREATE PROCEDURE comment_del(IN in_id int)
BEGIN
  select `tid`,`pid`,`uid`,`replys` INTO @tid,@pid,@uid,@replys from `comment` where `id` = in_id;
  delete from `comment` where `id` = in_id;
  if @pid >0 then #楼中楼
    update `comment` set `replys` = `replys`-1 where `id` = @pid;
  elseif @replys > 0 then #pid = 0,一般回复 且有子回复删除他们
    delete from `comment` where `pid` = in_id;
  end if;
END;

#编辑评论
CREATE PROCEDURE comment_edit(
  IN in_id int,
  IN in_content varchar(2000))
BEGIN
  update `comment` set `content` = in_content,`updated` = now() where `id` = in_id;
END;

#回复消息已读
CREATE PROCEDURE comment_read_s(IN in_id int)
BEGIN
  UPDATE `comment`SET `isread` = 1 WHERE `id` = in_id;
END;

#某篇回复消息全部已读
CREATE PROCEDURE comment_read_t(
  IN in_uid int,
  IN in_tid int)
BEGIN
  UPDATE `comment`SET `isread` = 1 WHERE `uid` = in_uid AND `tid` = in_tid;
END;

#回复消息全部已读
CREATE PROCEDURE comment_read_a(IN in_uid int)
BEGIN
  UPDATE `comment`SET `isread` = 1 WHERE `uid` = in_uid; 
END;

#收藏文章
DROP PROCEDURE IF EXISTS `star_add`;
CREATE PROCEDURE star_add(
  IN in_uid int,
  IN in_tid int)
BEGIN
  set @title = (select `title` from `post` where `tid` = in_tid);
  if @title is not null then
    INSERT INTO `star`(`uid`,`tid`,`title`) VALUES (in_uid,in_tid,@title);
  end if;
END;

#取消收藏文章
CREATE PROCEDURE star_del_bytid(
  IN in_uid int,
  IN in_tid int)
BEGIN
  DELETE FROM `star` WHERE `uid` = in_uid AND `tid` = in_tid;
END;

#取消收藏文章
CREATE PROCEDURE star_del_byid(IN in_id int)
BEGIN
  DELETE FROM `star` WHERE `id` = in_id;
END;

#添加被@消息
DROP PROCEDURE IF EXISTS `reply_add_at`;
CREATE PROCEDURE reply_add_at(
  IN in_myuid int,
  IN in_tuid int,
  IN in_tid int,#为评论id
  IN in_commentid int,
  IN in_content varchar(200))
BEGIN
  if in_myuid <> in_tuid then
    set @author = (select `username` from `user` where `uid` = in_myuid);
    select `title` INTO @title from `post` where `tid` = in_tid;
    INSERT INTO `atmessage`(`uid`,`fuid`,`fauthor`,`tid`,`cid`,`title`,`content`) VALUES (in_tuid,in_myuid,@author,in_tid,in_commentid,@title,in_content);
  end if;
END;

#@消息已读
CREATE PROCEDURE atmsg_read_s(IN in_id int)
BEGIN
  UPDATE `atmessage`SET `isread` = 1 WHERE `id` = in_id; 
END;

#某篇文章@消息全部已读
CREATE PROCEDURE atmsg_read_t(
  IN in_uid int,
  IN in_tid int)
BEGIN
  UPDATE `atmessage`SET `isread` = 1 WHERE `uid` = in_uid AND `tid` = in_tid; 
END;

#@消息全部已读
CREATE PROCEDURE atmsg_read_a(IN in_uid int)
BEGIN
  UPDATE `atmessage`SET `isread` = 1 WHERE `uid` = in_uid; 
END;


#增加关注
CREATE PROCEDURE follow_add(
  IN in_uid int,
  IN in_tuid int,
  IN in_note varchar(25))
BEGIN
  if in_uid <> in_tuid then
    set @tauthor = (select `username` from `user` where `uid` = in_tuid);
    set @author = (select `username` from `user` where `uid` = in_uid);
  INSERT INTO `follow`(`uid`,`tuid`,`name`,`tname`,`note`) VALUES (in_uid,in_tuid,@author,@tauthor,in_note);
  end if;
END;

#取消关注
CREATE PROCEDURE follow_del(
  IN in_uid int,
  IN in_tuid int)
BEGIN
  DELETE FROM `follow` WHERE `uid` = in_uid AND `tuid` = in_tuid;
END;

#发送聊天消息
CREATE PROCEDURE chat_add(
  IN in_send int,
  IN in_rev int,
  IN in_content varchar(500))
BEGIN
  INSERT INTO `chat`(`sned`,`recieve`,`content`) VALUES (in_send,in_rev,in_content);
END;

#和某人的聊天置为已读
CREATE PROCEDURE chat_read(
  IN in_rev int,
  IN in_send int)
BEGIN
  UPDATE `chat`SET `isread` = 1 WHERE `recieve` = in_rev AND `send` = in_send;
END;
```
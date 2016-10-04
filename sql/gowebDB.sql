---------------------------------------------------
-------论坛系统数据库设计--2351386755@qq.com---------
----------------2016-10-03-------------------------


----------------------------------------------------------------------------------
-------------------------------数据表定义------------------------------------------
----------------------------------------------------------------------------------
CREATE TABLE `user` (
	`uid` int NOT NULL AUTO_INCREMENT COMMENT '用户表用户id',
	`username` varchar(25) NOT NULL COMMENT '用户名',
	`password` varchar(50) NOT NULL COMMENT '密码',
	`email` varchar(25) NOT NULL COMMENT '邮件',
	`status` tinyint NOT NULL DEFAULT 0 COMMENT '0-正常 1-禁止访问',
	`sex` tinyint NOT NULL DEFAULT 0 CHECK IN ('0','1','2') COMMENT '性别0-未知，1-男，2-女',
	`description` varchar(200) NOT NULL DEFAULT '' COMMENT '个人描述',
	`exp` int NOT NULL DEFAULT 0 COMMENT '经验值',
	`sites` varchar(50) NOT NULL DEFAULT '' COMMENT '个人网站',
	`birthday` date NOT NULL DEFAULT '0000-00-00' COMMENT '生日',
	`messages` int NOT NULL DEFAULT 0 COMMENT '新消息数目',
	`posts` int NOT NULL DEFAULT 0 COMMENT '发帖数',
	`replys` int NOT NULL DEFAULT 0 COMMENT '回复数',
	`phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号码',
	`regtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
	PRIMARY KEY (`uid`) ,
	UNIQUE INDEX `user_username` (`username`),
	UNIQUE INDEX `user_email` (`email`)
);
/*
	INSERT INTO `user`(`username`,`password`,`email`) VALUES ('hehe01','just01','2351386755@qq.com');
*/

CREATE TABLE `category` (
	`cid` tinyint NOT NULL AUTO_INCREMENT COMMENT '版块表id',
	`title` varchar(25) NOT NULL COMMENT '版块名字',
	`description` varchar(150) NOT NULL DEFAULT '' COMMENT '描述',
	`posts` int NOT NULL DEFAULT 0 COMMENT '版块帖子计数',
	`todayposts` int NOT NULL DEFAULT 0 COMMENT '今日新帖数目'
	`lastpost` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后发表时间',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
	PRIMARY KEY (`cid`) ,
	UNIQUE INDEX `cat_name` (`title`)
);
/*
INSERT INTO `category`(`title`,`description`) VALUES ('forum01','dghajhgd');
*/

CREATE TABLE `post` (
	`tid` int NOT NULL AUTO_INCREMENT COMMENT '文章表id',
	`cid` tinyint NOT NULL  COMMENT '版块id',
	`uid` int NOT NULL  COMMENT '用户id',
	`author` varchar(25) NOT NULL DEFAULT '' COMMENT '用户名',
	`title` varchar(50) NOT NULL DEFAULT '' COMMENT '标题',
	`content` varchar(5000) NOT NULL DEFAULT '' COMMENT '内容',
	`type` tinyint NOT NULL DEFAULT 0 COMMENT '类型0-一般，1-XX，2-XX',
	`status` tinyint NOT NULL DEFAULT 0 COMMENT '0-正常，1-不可回复',
	`views` int NOT NULL DEFAULT 0 COMMENT '查看数',
	`replys` int NOT NULL DEFAULT 0 COMMENT '回复数',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
	`updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '编辑时间',
	`lastreply` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后回复时间',
	PRIMARY KEY (`tid`) ,
	INDEX `post_new` (`cid`,`lastreply`),
	INDEX `post_who` (`uid`),
	FOREIGN KEY (`cid`) REFERENCES `category` (`cid`) ON DELETE CASCADE,
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

/*
INSERT INTO `post`(`cid`,`uid`,`author`,`title`,`content`) VALUES ('1','1','hehe01','title01','content01');
DELETE FROM `post` WHERE `tid` = 3;
*/


CREATE TABLE `comment` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '评论表id',
	`tid` int NOT NULL COMMENT '帖子id',
	`pid` int NOT NULL DEFAULT 0 COMMENT '父评论id 0-一般回复，0+回复某个回复（楼中楼）',
	`uid` int NOT NULL COMMENT '用户id',
	`author` varchar(25) NOT NULL DEFAULT '' COMMENT '用户名',
	`content` varchar(2000) NOT NULL DEFAULT '' COMMENT '内容',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
	`updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '编辑时间',
	`replys` int NOT NULL DEFAULT 0 COMMENT '回复计数大于0表示有楼中楼回复',
	PRIMARY KEY (`id`) ,
	INDEX `comment_tid` (`tid`),
	INDEX `comment_who` (`uid`),
	INDEX `comment_pid` (`pid`),
	FOREIGN KEY (`tid`) REFERENCES `post` (`tid`) ON DELETE CASCADE,
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);
/*
INSERT INTO `comment`(`tid`,`uid`,`author`,`content`) VALUES ('1','1','hehe01','comment08');
INSERT INTO `comment`(`tid`,`pid`,`uid`,`author`,`content`) VALUES ('1','2','3','hehe03','comment3-3');
*/

CREATE TABLE `star` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '收藏表id',
	`uid` int NOT NULL COMMENT '用户id',
	`tid` int NOT NULL COMMENT '帖子id',
	`title` varchar(100) NOT NULL DEFAULT '' COMMENT '帖子标题',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
	PRIMARY KEY (`id`) ,
	UNIQUE INDEX `star_unique` (`uid`, `tid`),
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `message` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '回复消息表id',
	`uid` int NOT NULL COMMENT '我的uid',
	`fuid` int NOT NULL COMMENT '来自用户uid',
	`fauthor` varchar(25) NOT NULL DEFAULT '' COMMENT '来自用户名'
	`type` tinyint NOT NULL DEFAULT 0 COMMENT '0回复消息 1@消息 2好友请求 3聊天消息 4系统消息',
	`fromid` int NOT NULL COMMENT '可能来自tid也可能来自别的id具体看type对应关系',
	`title` varchar(50) NOT NULL DEFAULT '',
	`content` varchar(200) NOT NULL DEFAULT '' COMMENT '消息内容',
	`isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0未读1已读',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
	PRIMARY KEY (`id`) ,
	INDEX `msg_uid` (`uid`),
	INDEX `msg_from` (`fuid`),
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `follow` (
	`uid` int NOT NULL COMMENT '我的uid',
	`fuid` int NOT NULL COMMENT '对方uid',
	`fname` varchar(25) NOT NULL DEFAULT '' COMMENT '对方用户名'
	`note` varchar(25) NOT NULL DEFAULT '' COMMENT '备注名'
	`relation` tinyint NOT NULL DEFAULT 0 COMMENT '关注表还没想好冗余',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
	PRIMARY KEY (`uid`,`fuid`) ,
	INDEX `friend_rev` (`fuid`),
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
	FOREIGN KEY (`fuid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `chat` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '聊天表id',
	`send` int NOT NULL COMMENT '我的uid',
	`recieve` int NOT NULL COMMENT '来自uid',
	`content` varchar(500) NOT NULL DEFAULT '' COMMENT '内容',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
	`isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0-未读，1-已读，只有recieve读取才置为已读',
	PRIMARY KEY (`id`) ,
	INDEX `chat_my_s` (`send`),
	INDEX `chat_my_r` (`recieve`),
	FOREIGN KEY (`send`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
	FOREIGN KEY (`recieve`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

------------------------------------------------------------------------------------
-------------------------------触发器定义--------------------------------------------
------------------------------------------------------------------------------------
#发表文章 更新category计数 更新用户经验值+3
CREATE TRIGGER `t_post_add` After DELETE ON `post` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` + 3,`posts` = `posts` + 1 where `uid` = new.uid;
	update `category` set `posts` = `posts` + 1, `todayposts` = `todayposts` + 1 where `cid` = new.cid;
END;

#删除文章 更新category计数 更新用户经验值-3
CREATE TRIGGER `t_post_del` After DELETE ON `post` FOR EACH ROW BEGIN
	update `category` set `posts` = `posts` -1 ,`todayposts` = `todayposts` + 1 where `cid` = old.cid;
	update `user` set `exp` = `exp` - 3,`posts` = `posts` -1 where `uid` = old.uid;
END;

#删除评论
CREATE TRIGGER `t_comment_del` After DELETE ON `comment` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` - 1,`replys` = `replys` - 1 where `uid` = old.uid;
	update `post` set `replys` = `replys` - 1 where `tid` = old.tid;
END;

#增加评论
CREATE TRIGGER `t_comment_add` After INSERT ON `comment` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` + 1,`replys` = `replys` + 1 where `uid` = new.uid; #增加经验值
	update `post` set `replys` = `replys` + 1,`lastreply` = new.created where `tid` = new.tid;
END;

#增加消息
CREATE TRIGGER `t_message_add` After INSERT ON `message` FOR EACH ROW BEGIN
	update `user` set `messages` = `messages` + 1 where uid = new.uid;
END;

#增加聊天消息
CREATE TRIGGER `t_chat_add` After INSERT ON `chat` FOR EACH ROW BEGIN
	update `user` set `messages` = `messages` + 1 where uid = new.recieve;
END;

#增加被关注消息
CREATE TRIGGER `t_follow_add` After INSERT ON `follow` FOR EACH ROW BEGIN
	update `user` set `messages` = `messages` + 1 where uid = new.fuid;
END;

------------------------------------------------------------------------------------
-------------------------------计划任务定义------------------------------------------
------------------------------------------------------------------------------------
--查询手否开启 show variables like '%event%';
--开启 SET GLOBAL event_scheduler = ON; mysql 8 可以持久化SET PERSIST
DROP EVENT IF EXISTS `event_aday`;
CREATE EVENT `event_aday`
ON SCHEDULE EVERY 1 DAY STARTS '2016-10-01 00:00:00'
ON COMPLETION PRESERVE
ENABLE
DO
BEGIN
	update `category` SET `todayposts` = 0 WHERE 1;
END;


-----存储过程-----

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
CREATE PROCEDURE user_upd(
	IN in_uid int
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
DROP PROCEDURE IF EXISTS p_add_post;
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

#关闭文章
CREATE PROCEDURE post_del(IN in_tid int)
BEGIN
	delete from `post` where `tid` = in_tid;
END;

#添加评论
CREATE PROCEDURE p_add_comment(
	IN in_tid int,
	IN in_pid int,
	IN in_uid int,
	IN in_content varchar(2000))  
BEGIN
	if in_pid in (select `id`,'0' from `comment` where `tid` = in_tid) then #楼中楼回复
		if in_pid > 0 then #lzl
			update `comment` set `replys` = `replys` + 1  where `id` = in_pid;
		end if;
		INSERT INTO `comment`(`tid`,`pid`,`uid`,`author`,`content`) VALUES (in_tid,in_pid,in_uid,@author,in_content);
	end if;
END;

#删除评论
CREATE PROCEDURE p_del_comment(IN in_id int)
BEGIN
	select @tid = `tid`,@pid = `pid`,@uid = `uid`,@replys = `replys` from `comment` where `id` = in_id;
	if @pid >0 then #楼中楼
		update `comment` set `replys` = `replys`-1 where `id` = @pid;
	elseif @replys > 0 then #pid = 0,一般回复 且有子回复删除他们
		delete from `comment` where `pid` = in_id;
	end if;
	delete from `comment` where `id` = in_id;
END;
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
/*
	INSERT INTO `user`(`username`,`password`,`email`) VALUES ('hehe01','just01','2351386755@qq.com');
*/

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

/*
INSERT INTO `post`(`cid`,`uid`,`author`,`title`,`content`) VALUES ('1','1','hehe01','title01','content01');
DELETE FROM `post` WHERE `tid` = 3;
*/


CREATE TABLE `comment` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '评论表id',
	`tid` int NOT NULL COMMENT '帖子id',
	`pid` int NOT NULL DEFAULT 0 COMMENT '父评论id 0-一般回复0+回复某个回复（楼中楼）',
	`uid` int NOT NULL COMMENT '用户id',
	`touid` int NOT NULL COMMENT '回复对象uid',
	`author` varchar(25) NOT NULL DEFAULT '' COMMENT '用户名',
	`content` varchar(2000) NOT NULL DEFAULT '' COMMENT '内容',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发表时间',
	`updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '编辑时间',
	`replys` int NOT NULL DEFAULT 0 COMMENT '回复计数大于0表示有楼中楼回复',
	PRIMARY KEY (`id`) ,
	INDEX `comment_tid` (`tid`),
	INDEX `comment_who` (`uid`),
	INDEX `comment_pid` (`pid`),
	FOREIGN KEY (`tid`) REFERENCES `post` (`tid`) ON DELETE CASCADE,
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE,
	FOREIGN KEY (`touid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);
/*
INSERT INTO `comment`(`tid`,`uid`,`author`,`content`) VALUES ('1','1','hehe01','comment08');
INSERT INTO `comment`(`tid`,`pid`,`uid`,`author`,`content`) VALUES ('1','2','3','hehe03','comment3-3');
*/

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

CREATE TABLE `message` (
	`id` int NOT NULL AUTO_INCREMENT COMMENT '回复消息表id',
	`uid` int NOT NULL COMMENT '我的uid',
	`fuid` int NOT NULL COMMENT '来自用户uid',
	`fauthor` varchar(25) NOT NULL DEFAULT '' COMMENT '来自用户名',
	`type` tinyint NOT NULL DEFAULT 0 COMMENT '0回复消息 1@消息 2被关注消息 3系统消息',
	`fromid` int NOT NULL COMMENT '可能来自tid也可能来自别的id具体看type对应关系',
	`title` varchar(50) NOT NULL DEFAULT '',
	`content` varchar(200) NOT NULL DEFAULT '' COMMENT '消息内容',
	`isread` tinyint NOT NULL DEFAULT 0 COMMENT '是否已读0未读1已读',
	`created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间',
	PRIMARY KEY (`id`) ,
	INDEX `msg_uid` (`uid`),
	INDEX `msg_from` (`fuid`),
	FOREIGN KEY (`uid`) REFERENCES `user` (`uid`) ON DELETE CASCADE
);

CREATE TABLE `follow` (
	`id` int NOT NULL AUTO_INCREMENT,
	`uid` int NOT NULL COMMENT '我的uid',
	`tuid` int NOT NULL COMMENT '对方uid',
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

------------------------------------------------------------------------------------
-------------------------------触发器定义--------------------------------------------
------------------------------------------------------------------------------------
#发表文章 更新category计数 更新用户经验值+3
CREATE TRIGGER `t_post_add` After INSERT ON `post` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` + 3,`posts` = `posts` + 1 where `uid` = new.uid;
	update `category` set `posts` = `posts` + 1, `todayposts` = `todayposts` + 1 where `cid` = new.cid;
END;

#删除文章 更新category计数 更新用户经验值-3
CREATE TRIGGER `t_post_del` After DELETE ON `post` FOR EACH ROW BEGIN
	update `category` set `posts` = `posts` -1 ,`todayposts` = `todayposts` + 1 where `cid` = old.cid;
	update `user` set `exp` = `exp` - 3,`posts` = `posts` -1 where `uid` = old.uid;
END;

#增加评论
CREATE TRIGGER `t_comment_add` After INSERT ON `comment` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` + 1,`replys` = `replys` + 1 where `uid` = new.uid; #增加经验值
	update `post` set `replys` = `replys` + 1,`lastreply` = new.created where `tid` = new.tid;
END;

#删除评论
CREATE TRIGGER `t_comment_del` After DELETE ON `comment` FOR EACH ROW BEGIN
	update `user` set `exp` = `exp` - 1,`replys` = `replys` - 1 where `uid` = old.uid;
	update `post` set `replys` = `replys` - 1 where `tid` = old.tid;
END;

#增加消息
CREATE TRIGGER `t_message_add` After INSERT ON `message` FOR EACH ROW BEGIN
	update `user` set `messages` = `messages` + 1 where uid = new.uid;
END;

#消息已读
CREATE TRIGGER `t_message_read` After UPDATE ON `message` FOR EACH ROW BEGIN
	if old.isread = 0 AND new.isread = 1 then
		UPDATE `user` SET `messages` = `messages` - 1 WHERE `uid` = old.uid;
 	end if;
END;

#增加聊天消息
CREATE TRIGGER `t_chat_add` After INSERT ON `chat` FOR EACH ROW BEGIN
	update `user` set `messages` = `messages` + 1 where uid = new.recieve;
END;

#聊天消息已读
CREATE TRIGGER `t_chat_read` After UPDATE ON `chat` FOR EACH ROW BEGIN
	if old.isread = 0 AND new.isread = 1 then
		UPDATE `user` SET `messages` = `messages` - 1 WHERE `uid` = new.recieve;
 	end if;
END;

------------------------------------------------------------------------------------
-------------------------------计划任务定义------------------------------------------
------------------------------------------------------------------------------------
--查询手否开启 show variables like '%event%';
--开启 SET PERSIST GLOBAL event_scheduler = ON; mysql 8 可以持久化SET PERSIST
DROP EVENT IF EXISTS `event_aday`;
CREATE EVENT `event_aday`
ON SCHEDULE EVERY 1 DAY STARTS '2016-10-01 00:00:00'
ON COMPLETION PRESERVE
ENABLE
DO update `category` SET `todayposts` = 0 WHERE 1;


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

#添加评论
DROP PROCEDURE IF EXISTS `comment_add`;
CREATE PROCEDURE comment_add(
	IN in_tid int,
	IN in_pid int,
	IN in_uid int,
	IN in_tuid int,
	IN in_content varchar(2000))
BEGIN
	IF in_pid >= 0 THEN 
		IF in_pid > 0 THEN #楼中楼回复
			update `comment` set `replys` = `replys` + 1  where `id` = in_pid;
		END IF;
		
		set @author = (select `username` from `user` where `uid` = in_uid);
		INSERT INTO `comment`(`tid`,`pid`,`uid`,`tuid`,`author`,`content`) VALUES (in_tid,in_pid,in_uid,in_tuid,@author,in_content);
		
		#增加消息
		set @title = (select `title` from `post` where `tid` = in_tid);
		set @lastid = LAST_INSERT_ID();
		INSERT INTO `message`(`uid`,`fuid`,`fauthor`,`type`,`fromid`,`title`,`content`) VALUES (in_tuid,in_uid,@author,'0',@lastid,@title,in_content);
	END IF;
END;

#删除评论
CREATE PROCEDURE comment_del(IN in_id int)
BEGIN
	select @tid = `tid`,@pid = `pid`,@uid = `uid`,@replys = `replys` from `comment` where `id` = in_id;
	if @pid >0 then #楼中楼
		update `comment` set `replys` = `replys`-1 where `id` = @pid;
	elseif @replys > 0 then #pid = 0,一般回复 且有子回复删除他们
		delete from `comment` where `pid` = in_id;
	end if;
	delete from `comment` where `id` = in_id;
END;

#编辑评论
CREATE PROCEDURE comment_edit(
	IN in_id int,
	IN in_content varchar(2000))
BEGIN
	update `comment` set `content` = in_content,`updated` = now() where `id` = in_id;
END;

#收藏文章
CREATE PROCEDURE star_add(
	IN in_uid int,
	IN in_tid int,
	IN in_title varchar(50))
BEGIN
	INSERT INTO `star`(`uid`,`tid`,`title`) VALUES (in_uid,in_tid,in_title);
END;

#取消收藏文章
CREATE PROCEDURE star_del(
	IN in_uid int,
	IN in_tid int)
BEGIN
	DELETE FROM `star` WHERE `uid` = in_uid AND `tid` = in_tid;
END;

#添加消息
CREATE PROCEDURE message_add(
	IN in_uid int,
	IN in_fuid int,
	IN in_type tinyint,
	IN in_fromid int,
	IN in_title varchar(50),
	IN in_content varchar(200))
BEGIN
	set @author = (select `username` from `user` where `uid` = in_uid);
	INSERT INTO `message`(`uid`,`fuid`,`fauthor`,`type`,`fromid`,`title`,`content`)
	VALUES (in_uid,in_fuid,@author,in_type,in_fromid,in_title,in_content);
END;

#消息已读
CREATE PROCEDURE message_read(IN in_id int)
BEGIN
	UPDATE `message`SET `isread` = 1 WHERE `id` = in_id;
END;

#某一类消息全部已读
CREATE PROCEDURE message_read(
	IN in_uid int,
	IN in_type tinyint)
BEGIN
	UPDATE `message`SET `isread` = 1 WHERE `uid` = in_uid AND `type` = in_type; 
END;


#增加关注
CREATE PROCEDURE follow_add(
	IN in_uid int,
	IN in_tuid int,
	IN in_note varchar(25))
BEGIN
	set @tauthor = (select `username` from `user` where `uid` = in_tuid);
	set @fauthor = (select `username` from `user` where `uid` = in_uid);
	INSERT INTO `follow`(`uid`,`tuid`,`tname`,`note`) VALUES (in_uid,in_tuid,@tauthor,in_note);

	#增加消息
	INSERT INTO `message`(`uid`,`fuid`,`fauthor`,`type`,`fromid`) VALUES (in_tuid,in_uid,@fauthor,'2',LAST_INSERT_ID());
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
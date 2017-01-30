---------------------------------------------------
-------论坛系统数据库设计--2351386755@qq.com---------
----------------2016-10-03-------------------------
/*
todo
1.系统提醒表
2.被关注消息表
3.版块置顶表
*/
----------------------------------------------------------------------------------
-------------------------------数据表定义------------------------------------------
----------------------------------------------------------------------------------
DROP TABLE IF EXISTS "user";
CREATE TABLE "user" (
	"uid" serial NOT NULL,
	"username" varchar(15) NOT NULL,
	"password" varchar(40) NOT NULL,
	"email" varchar(30) NOT NULL,
	"sex" int2 NOT NULL DEFAULT 0, --'性别0-未知，1-男，2-女'
	"description" varchar(200) NOT NULL DEFAULT '',
	"sites" varchar(30) NOT NULL DEFAULT '',--'个人主页'
	"birthday" date,
	"messages" int[] NOT NULL DEFAULT '{0,0,0}',--[1]评论消息数目[2]@消息[3]聊天消息
	"exp" int NOT NULL DEFAULT 0 ,--'经验值'
	"posts" int NOT NULL DEFAULT 0,--'发帖数'
	"replys" int NOT NULL DEFAULT 0 ,-- '回复数'
	"phone" varchar(20) NOT NULL DEFAULT '',
	"regtime" timestamp NOT NULL DEFAULT current_timestamp,
	PRIMARY KEY ("uid"),
	CONSTRAINT "user_username" UNIQUE("username"),
	CONSTRAINT "user_email" UNIQUE("email"),
	CONSTRAINT "user_phone" UNIQUE("phone")
)WITH (OIDS=FALSE);


CREATE TABLE "type" (
	"tid" serial NOT NULL,  -- '版块表id',
	"typename" varchar(25) NOT NULL, -- '版块名字',
	"description" varchar(150) NOT NULL DEFAULT '', -- '描述',
	"sticks" int[] ELEMENT REFERENCES "post"."pid",--置顶表
	"posts" int NOT NULL DEFAULT 0, -- '版块帖子计数'
	"replys" int NOT NULL DEFAULT 0, -- 评论计数
	"newposts" int NOT NULL DEFAULT 0, -- '今日新发帖数目',
	"newreplys" int NOT NULL DEFAULT 0, -- '今日新回复数目',
	PRIMARY KEY ("cid") ,
	CONSTRAINT "cat_name" UNIQUE("type")
);


CREATE TABLE "post" (
	"pid" serial NOT NULL,  -- '文章表id',
	"tid" int NOT NULL,  -- '版块id',
	"uid" int NOT NULL,  -- '用户id',
	"author" varchar(15) NOT NULL,--当用户名更改后要修改此
	"title" varchar(50) NOT NULL DEFAULT '', -- '标题'
	"content" text NOT NULL DEFAULT '', -- '内容'
	"tags" varchar(10)[],--标签
	"views" int NOT NULL DEFAULT 0, -- '查看数'
	"replys" int[] NOT NULL DEFAULT '{0,0}' -- '回复数[1]总回复数[2]未读回复数'
	"agrees" int[],--点赞uid列表
	"disagree" int[],--不赞uid同列表
	"status" int2 NOT NULL DEFAULT 0, -- '0-正常，1-不可回复2不可查看'
	"type" int2 NOT NULL DEFAULT 0,-- '类型0-一般，1-管理员加精华'
	"created" timestamp NOT NULL DEFAULT current_timestamp, -- '发表时间'
	"updated" timestamp -- '编辑时间'
	"lastreply" timestamp NOT NULL DEFAULT current_timestamp, -- '最后回复时间'
	PRIMARY KEY ("tid"),
	INDEX "post_new" ("cid","lastreply"),
	INDEX "post_who" ("uid"),
	FOREIGN KEY ("cid") REFERENCES "category" ("cid") ON DELETE CASCADE,
	FOREIGN KEY ("uid") REFERENCES "user" ("uid") ON DELETE CASCADE
);

DROP TABLE IF EXISTS "comment";
CREATE TABLE "comment" (
	"cid" serial NOT NULL,  -- '评论表id'
	"pid" int NOT NULL, -- '帖子id'
	"ppid" int NOT NULL DEFAULT 0, -- '父评论id 0-一般回复0+回复某个回复（楼中楼）'
	"uid" int NOT NULL, -- '用户id'
	"author" varchar(15) NOT NULL,
	"tuid" int NOT NULL, -- '回复对象uid'
	"content" text NOT NULL DEFAULT '', -- '内容'
	"agree" int[],--点赞uid列表
	"disagree" int[],--不赞同uid列表
	"created" timestamp NOT NULL DEFAULT current_timestamp, -- '发表时间'
	"updated" timestamp, -- '编辑时间'
	"replys" int[] NOT NULL DEFAULT 0, -- '回复计数大于0表示有楼中楼回复'
	PRIMARY KEY ("id"),
	INDEX "comment_tid" ("tid"),
	INDEX "comment_who" ("uid"),
	INDEX "comment_pid" ("pid"),
	FOREIGN KEY ("tid") REFERENCES "post" ("tid") ON DELETE CASCADE,
	FOREIGN KEY ("uid") REFERENCES "user" ("uid") ON DELETE CASCADE,
	FOREIGN KEY ("tuid") REFERENCES "user" ("uid") ON DELETE CASCADE
);

DROP TABLE IF EXISTS "at";
CREATE TABLE "at" (
	"aid" serial NOT NULL,  -- '回复和@表 id',
	"uid" int NOT NULL, -- '我的uid',
	"tuid" int NOT NULL -- '对方uid',
	"cid" int NOT NULL -- '来自文章评论id'
	"created" timestamp NOT NULL DEFAULT current_timestamp -- '时间',
	PRIMARY KEY ("id") ,
	INDEX "msg_uid" ("uid"),
	INDEX "msg_from" ("fuid"),
	FOREIGN KEY ("uid") REFERENCES "user" ("uid") ON DELETE CASCADE,
	FOREIGN KEY ("fuid") REFERENCES "user" ("uid") ON DELETE CASCADE
);

DROP TABLE IF EXISTS "star";
CREATE TABLE "star" (
	"sid" serial NOT NULL,  -- '收藏表id'
	"uid" int NOT NULL, -- '用户id'
	"pid" int NOT NULL, -- '帖子id'
	"created" timestamp NOT NULL DEFAULT current_timestamp, -- '收藏时间',
	PRIMARY KEY ("sid"),
	UNIQUE INDEX "star_unique" ("uid", "tid"),
	FOREIGN KEY ("uid") REFERENCES "user" ("uid") ON DELETE CASCADE
);


DROP TABLE IF EXISTS "follow";
CREATE TABLE "follow" (
	"fid" serial NOT NULL,
	"uid" int NOT NULL, -- '我的uid',
	"tuid" int NOT NULL, -- '对方uid',
	"note" varchar(15) NOT NULL DEFAULT '', -- '备注名',
	"created" timestamp NOT NULL DEFAULT current_timestamp,
	PRIMARY KEY ("fid") ,
	UNIQUE INDEX "follow_me_t" ("uid","tuid"),
	INDEX "friend_rev" ("tuid"),
	FOREIGN KEY ("uid") REFERENCES "user" ("uid") ON DELETE CASCADE,
	FOREIGN KEY ("tuid") REFERENCES "user" ("uid") ON DELETE CASCADE
);

DROP TABLE IF EXISTS "chat";
CREATE TABLE "chat" (
	"id" serial NOT NULL -- '聊天表id',
	"uid" int NOT NULL -- '我的uid',
	"tuid" int NOT NULL -- '对方uid',
	"content" text NOT NULL DEFAULT '' -- '内容',
	"created" timestamp NOT NULL DEFAULT current_timestamp -- '时间',
	"isread" bool NOT NULL DEFAULT false -- '是否已读0-未读，1-已读，只有recieve读取才置为已读',
	PRIMARY KEY ("id") ,
	INDEX "chat_my_s" ("send"),
	INDEX "chat_my_r" ("recieve"),
	FOREIGN KEY ("send") REFERENCES "user" ("uid") ON DELETE CASCADE,
	FOREIGN KEY ("recieve") REFERENCES "user" ("uid") ON DELETE CASCADE
);

DROP TABLE IF EXISTS "ban";
CREATE TABLE "ban" (
	"bid" serial NOT NULL, -- '禁止表id'
	"uid" int NOT NULL,
	"reason" varchar(200) NOT NULL DEFAULT '', -- 理由
	"start" timestamp NOT NULL DEFAULT current_timestamp, -- '时间',
	"end" timestamp NOT NULL DEFAULT current_timestamp,
	PRIMARY KEY ("bid")
);

#------------------------------------------------------------------------------------
#-------------------------------触发器定义--------------------------------------------
#------------------------------------------------------------------------------------
#发表文章 更新category计数 更新用户经验值+3
CREATE TRIGGER "t_post_add" After INSERT ON "post" FOR EACH ROW BEGIN
	update "user" set "exp" = "exp" + 3,"posts" = "posts" + 1 where "uid" = new.uid;
	update "category" set "posts" = "posts" + 1,"lastpost" = new.created, "todayposts" = "todayposts" + 1 where "cid" = new.cid;
END;

#删除文章 更新category计数 更新用户经验值-3
DROP TRIGGER IF EXISTS "t_post_del";
CREATE TRIGGER "t_post_del" After DELETE ON "post" FOR EACH ROW BEGIN
	update "category" set "posts" = "posts" -1 ,"todayposts" = "todayposts" - 1 where "cid" = old.cid;
	update "user" set "exp" = "exp" - 3,"posts" = "posts" -1 where "uid" = old.uid;
END;

#增加之前
CREATE TRIGGER "t_--_add_b" Before INSERT ON "--" FOR EACH ROW BEGIN
	if new.uid = new.tuid then
		set new.isread = '1';
	end if;
END;

#增加评论
CREATE TRIGGER "t_--_add_a" After INSERT ON "--" FOR EACH ROW BEGIN
	update "user" set "exp" = "exp" + 1,"replys" = "replys" + 1 where "uid" = new.uid; #增加经验值
	update "post" set "replys" = "replys" + 1,"lastreply" = new.created where "tid" = new.tid;
	if new.isread = 0 AND new.uid <> new.tuid then
		UPDATE "user" SET "messages" = "messages" + 1 where "uid" = new.tuid;
	end if;
END;

#删除评论
CREATE TRIGGER "t_--_del" After DELETE ON "--" FOR EACH ROW BEGIN
	update "user" set "exp" = "exp" - 1,"replys" = "replys" - 1 where "uid" = old.uid;
	update "post" set "replys" = "replys" - 1 where "tid" = old.tid;
	if old.isread = 0 then #这条消息未读 则减一
		update "user" set "messages" = "messages" - 1 where "uid" = old.tuid;
	end if;
END;

#编辑评论
CREATE TRIGGER "t_--_upd" After UPDATE ON "--" FOR EACH ROW BEGIN
	if old.isread = 0 AND new.isread = 1 then
		UPDATE "user" SET "messages" = "messages" - 1 where "uid" = old.tuid;
	end if;
END;

#增加@消息之前
CREATE TRIGGER "t_atmsg_add_b" Before INSERT ON "atmessage" FOR EACH ROW BEGIN
	if new.uid = new.fuid then
		set new.isread = '1';
	end if;
END;

#增加@消息
CREATE TRIGGER "t_atmsg_add_a" After INSERT ON "atmessage" FOR EACH ROW BEGIN
	if new.isread = 0 AND new.uid <> new.fuid then
		update "user" set "messages" = "messages" + 1 where uid = new.uid;
	end if;
END;

#@消息已读
CREATE TRIGGER "t_atmsg_upd" After UPDATE ON "atmessage" FOR EACH ROW BEGIN
	if old.isread = 0 AND new.isread = 1 then
		UPDATE "user" SET "messages" = "messages" - 1 WHERE "uid" = old.uid;
 	end if;
END;


#增加聊天消息
CREATE TRIGGER "t_chat_add" After INSERT ON "chat" FOR EACH ROW BEGIN
	if new.send <> new.recieve then
		update "user" set "messages" = "messages" + 1 where uid = new.recieve;
	end if;
END;

#聊天消息已读
CREATE TRIGGER "t_chat_read" After UPDATE ON "chat" FOR EACH ROW BEGIN
	if old.isread = 0 AND new.isread = 1 then
		UPDATE "user" SET "messages" = "messages" - 1 WHERE "uid" = new.recieve;
 	end if;
END;

#------------------------------------------------------------------------------------
#-------------------------------计划任务定义------------------------------------------
#------------------------------------------------------------------------------------
#查询手否开启 show variables like '%event%';
#开启 SET PERSIST GLOBAL event_scheduler = ON; mysql 8 可以持久化SET PERSIST
DROP EVENT IF EXISTS "event_aday";
CREATE EVENT "event_aday"
ON SCHEDULE EVERY 1 DAY STARTS '2016-10-01 00:00:00'
ON COMPLETION PRESERVE
ENABLE
DO update "category" SET "todayposts" = 0 WHERE 1;


#-----存储过程-----
#新增分类
DROP PROCEDURE IF EXISTS "cate_add";
CREATE PROCEDURE cate_add(
	IN in_title varchar(25),
	IN in_des varchar(150))
BEGIN
	INSERT INTO "category"("title","description") VALUES (in_title,in_des);
END;

#编辑分类
CREATE PROCEDURE cate_edit(
	IN in_cid int,
	IN in_title varchar(25),
	IN in_des varchar(150))
BEGIN
	UPDATE "category" SET "title" = in_title,"description" = in_des WHERE "cid" = in_cid;
END;

#删除分类 最好不要
CREATE PROCEDURE cate_del(IN in_cid int)
BEGIN
	DELETE FROM "category" WHERE "cid" = in_cid;
END;


#新增用户
CREATE PROCEDURE user_reg(
	IN in_name varchar(25),
	IN in_pass varchar(50),
	IN in_email varchar(25),
	IN sex int2)
BEGIN
	INSERT INTO "user"("username","password","email","sex") VALUES (in_name,in_pass,in_email,sex);
END;

#删除用户
CREATE PROCEDURE user_del(IN in_uid int)
BEGIN
	delete from "user" where "uid" = in_uid;
END;

#编辑用户
CREATE PROCEDURE user_edit(
	IN in_uid int,
	IN in_sex int2,
	IN in_des varchar(200),
	IN in_sites varchar(50),
	IN in_birth date,
	IN in_phone varchar(20))
BEGIN
	UPDATE "user" SET "sex" = in_sex,"description" = in_des,"sites" = in_sites,"birthday" = in_birth,"phone" = in_phone WHERE "uid" = in_uid;
END;

#禁止用户
CREATE PROCEDURE user_bolck(IN in_uid int)
BEGIN
	UPDATE "user" SET "status" = '1' WHERE "uid" = in_uid;
END;

#允许用户
CREATE PROCEDURE user_open(IN in_uid int)
BEGIN
	UPDATE "user" SET "status" = '0' WHERE "uid" = in_uid;
END;

#修改密码
CREATE PROCEDURE user_changepass(
	IN in_uid int,
	IN in_pass varchar(50))
BEGIN
	UPDATE "user" SET "password" = in_pass WHERE "uid" = in_uid;
END;

#发表文章
DROP PROCEDURE IF EXISTS post_add;
CREATE PROCEDURE post_add(
	IN in_cid int,
	IN in_uid int,
	IN in_title varchar(50),
	IN in_content varchar(8000))
BEGIN
	set @author = (select "username" from "user" where "uid" = in_uid);
	INSERT INTO "post"("cid","uid","author","title","content") VALUES (in_cid,in_uid,@author,in_title,in_content);
END;

#编辑文章
CREATE PROCEDURE post_edit(
	IN in_tid int,
	IN in_title varchar(50),
	IN in_content varchar(8000))
BEGIN
	set @timenow = now();
	update "post" set "title" = in_title,"content" = in_content,"updated" = @timenow,"lastreply" = @timenow where "tid" = in_tid;
END;

#删除文章
CREATE PROCEDURE post_del(IN in_tid int)
BEGIN
	delete from "post" where "tid" = in_tid;
END;

#禁止回复文章
CREATE PROCEDURE post_close_c(IN in_tid int)
BEGIN
	update "post" set "status" = '1' where "tid" = in_tid;
END;

#允许回复文章
CREATE PROCEDURE post_open_c(IN in_tid int)
BEGIN
	update "post" set "status" = '0' where "tid" = in_tid;
END;

#添加评论 回复楼主
DROP PROCEDURE IF EXISTS "--_add_lz";
CREATE PROCEDURE --_add_lz(
	IN in_tid int,
	IN in_uid int,
	IN in_content varchar(5000))
BEGIN
	set @author = (select "username" from "user" where "uid" = in_uid);
	select "uid","title" INTO @tuid,@title from "post" where "tid" = in_tid;
	INSERT INTO "--"("tid","uid","tuid","author","content","isread") VALUES (in_tid,in_uid,@tuid,@author,in_content,@isread);
END;

#楼中楼回复
DROP PROCEDURE IF EXISTS "--_add_cz";
CREATE PROCEDURE --_add_cz(
	IN in_tid int,
	IN in_pid int,
	IN in_uid int,
	IN in_content varchar(5000))
BEGIN
	select "uid" INTO @czuid from "--" where "id" = in_pid AND "tid" = in_tid;
	if @czuid is not null then
		set @author = (select "username" from "user" where "uid" = in_uid);
		INSERT INTO "--"("tid","pid","uid","tuid","author","content") VALUES (in_tid,in_pid,in_uid,@czuid,@author,in_content);
		update "--" set "replys" = "replys" + 1  where "id" = in_pid;
		set @lastid = LAST_INSERT_ID();
		select "title","uid" INTO @title,@lzuid from "post" where "tid" = in_tid;
		if @lzuid <> @czuid THEN #发送@给层主
			INSERT INTO "atmessage"("uid","fuid","fauthor","tid","cid","title","content") VALUES (@czuid,in_uid,@author,in_tid,@lastid,@title,left(in_content,50));
		end if;
	END IF;
END;

#删除评论
DROP PROCEDURE IF EXISTS "--_del";
CREATE PROCEDURE --_del(IN in_id int)
BEGIN
	select "tid","pid","uid","replys" INTO @tid,@pid,@uid,@replys from "--" where "id" = in_id;
	delete from "--" where "id" = in_id;
	if @pid >0 then #楼中楼
		update "--" set "replys" = "replys"-1 where "id" = @pid;
	elseif @replys > 0 then #pid = 0,一般回复 且有子回复删除他们
		delete from "--" where "pid" = in_id;
	end if;
END;

#编辑评论
CREATE PROCEDURE --_edit(
	IN in_id int,
	IN in_content varchar(5000))
BEGIN
	update "--" set "content" = in_content,"updated" = now() where "id" = in_id;
END;

#回复消息已读
CREATE PROCEDURE --_read_s(IN in_id int)
BEGIN
	UPDATE "--"SET "isread" = 1 WHERE "id" = in_id;
END;

#某篇回复消息全部已读
CREATE PROCEDURE --_read_t(
	IN in_uid int,
	IN in_tid int)
BEGIN
	UPDATE "--"SET "isread" = 1 WHERE "uid" = in_uid AND "tid" = in_tid;
END;

#回复消息全部已读
CREATE PROCEDURE --_read_a(IN in_uid int)
BEGIN
	UPDATE "--"SET "isread" = 1 WHERE "uid" = in_uid; 
END;

#收藏文章
DROP PROCEDURE IF EXISTS "star_add";
CREATE PROCEDURE star_add(
	IN in_uid int,
	IN in_tid int)
BEGIN
	set @title = (select "title" from "post" where "tid" = in_tid);
	if @title is not null then
		INSERT INTO "star"("uid","tid","title") VALUES (in_uid,in_tid,@title);
	end if;
END;

#取消收藏文章
CREATE PROCEDURE star_del_bytid(
	IN in_uid int,
	IN in_tid int)
BEGIN
	DELETE FROM "star" WHERE "uid" = in_uid AND "tid" = in_tid;
END;

#取消收藏文章
CREATE PROCEDURE star_del_byid(IN in_id int)
BEGIN
	DELETE FROM "star" WHERE "id" = in_id;
END;

#添加被@消息
DROP PROCEDURE IF EXISTS "reply_add_at";
CREATE PROCEDURE reply_add_at(
	IN in_myuid int,
	IN in_tuid int,
	IN in_tid int,#为评论id
	IN in_--id int,
	IN in_content varchar(200))
BEGIN
	if in_myuid <> in_tuid then
		set @author = (select "username" from "user" where "uid" = in_myuid);
		select "title" INTO @title from "post" where "tid" = in_tid;
		INSERT INTO "atmessage"("uid","fuid","fauthor","tid","cid","title","content") VALUES (in_tuid,in_myuid,@author,in_tid,in_--id,@title,left(in_content,50));
	end if;
END;

#@消息已读
CREATE PROCEDURE atmsg_read_s(IN in_id int)
BEGIN
	UPDATE "atmessage"SET "isread" = 1 WHERE "id" = in_id; 
END;

#某篇文章@消息全部已读
CREATE PROCEDURE atmsg_read_t(
	IN in_uid int,
	IN in_tid int)
BEGIN
	UPDATE "atmessage"SET "isread" = 1 WHERE "uid" = in_uid AND "tid" = in_tid; 
END;

#@消息全部已读
CREATE PROCEDURE atmsg_read_a(IN in_uid int)
BEGIN
	UPDATE "atmessage"SET "isread" = 1 WHERE "uid" = in_uid; 
END;


#增加关注
CREATE PROCEDURE follow_add(
	IN in_uid int,
	IN in_tuid int,
	IN in_note varchar(25))
BEGIN
	if in_uid <> in_tuid then
		set @tauthor = (select "username" from "user" where "uid" = in_tuid);
		set @author = (select "username" from "user" where "uid" = in_uid);
	INSERT INTO "follow"("uid","tuid","name","tname","note") VALUES (in_uid,in_tuid,@author,@tauthor,in_note);
	end if;
END;

#取消关注
CREATE PROCEDURE follow_del(
	IN in_uid int,
	IN in_tuid int)
BEGIN
	DELETE FROM "follow" WHERE "uid" = in_uid AND "tuid" = in_tuid;
END;

#发送聊天消息
CREATE PROCEDURE chat_add(
	IN in_send int,
	IN in_rev int,
	IN in_content varchar(500))
BEGIN
	INSERT INTO "chat"("sned","recieve","content") VALUES (in_send,in_rev,in_content);
END;

#和某人的聊天置为已读
CREATE PROCEDURE chat_read(
	IN in_rev int,
	IN in_send int)
BEGIN
	UPDATE "chat"SET "isread" = 1 WHERE "recieve" = in_rev AND "send" = in_send;
END;
#---------------------------------------------------
#--------论坛系统数据库设计2351386755@qq.com----------
#------------------2017-3-4-------------------------

#---------------------------------------------------
#----------------数据表定义--------------------------
#---------------------------------------------------

#bigint     -2^63-2^63-1    8字节
#int        -2^31-2^31 - 1  4字节
#smallint   -2^15-2^15 - 1  2字节
#tinyint    0-255           1字节
#FOUND_ROWS() : select
#ROW_COUNT()  : update delete insert.
#LEFT(in_content, 50) 截断字符串

#mysql 判断天数间隔
#1、利用TO_DAYS函数
#select to_days(now()) - to_days('20140831')
#2、利用DATEDIFF函数
#select datediff(now(),'20140831')

#函数TimeStampDiff()是MySQL本身提供的可以计算两个时间间隔的函数，语法为：
#TIMESTAMPDIFF(unit,datetime_expr1,datetime_expr2)，
#其中unit单位有如下几种，分别是：FRAC_SECOND (microseconds), SECOND, MINUTE, HOUR, DAY, WEEK, MONTH, QUARTER, or YEAR

#用户
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `uid`         INT          NOT NULL AUTO_INCREMENT,
  `username`    VARCHAR(25)  NOT NULL,
  `password`    VARCHAR(64)  NOT NULL,
  `email`       VARCHAR(25)  NOT NULL,
  `status`      TINYINT      NOT NULL DEFAULT 0, #0-ok 1-block
  `sex`         TINYINT      NOT NULL DEFAULT 0, #0-unknown 1-man 2-woman
  `exp`         INT          NOT NULL DEFAULT 0, #经验值
  `birthday`    DATE         NOT NULL DEFAULT 0,
  `phone`       VARCHAR(15)  NOT NULL DEFAULT '', #手机号码
  `description` VARCHAR(255) NOT NULL DEFAULT '',
  `site`        VARCHAR(50)  NOT NULL DEFAULT '', #个人网站
  `posts`       INT          NOT NULL DEFAULT 0, #发帖数
  `replys`      INT          NOT NULL DEFAULT 0, #回复数
  `newreplys`   INT          NOT NULL DEFAULT 0, #新回复数
  `newfollows`  INT          NOT NULL DEFAULT 0, #新粉丝数
  `newchats`    INT          NOT NULL DEFAULT 0, #新私信数
  `regtime`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  UNIQUE INDEX `user_username` (`username`),
  UNIQUE INDEX `user_email` (`email`)
);

#分类
DROP TABLE IF EXISTS `cate`;
CREATE TABLE `cate` (
  `cid`         TINYINT      NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(25)  NOT NULL, #版块名字
  `description` VARCHAR(255),
  `sticks`      VARCHAR(255) NOT NULL DEFAULT '', #置顶 tid,tid...
  `posts`       INT          NOT NULL DEFAULT 0, #版块帖子计数
  `todayposts`  INT          NOT NULL DEFAULT 0, #今日新帖数目
  `created`     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cid`),
  UNIQUE INDEX `cate_name` (`name`)
);

#帖子
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
  `tid`       INT            NOT NULL AUTO_INCREMENT,
  `cid`       TINYINT        NOT NULL,
  `uid`       INT            NOT NULL,
  `username`  VARCHAR(25)    NOT NULL, #帖子列表此项非常常见牺牲空间
  `title`     VARCHAR(50)    NOT NULL DEFAULT '',
  `content`   VARCHAR(10000) NOT NULL DEFAULT '',
  `tags`      VARCHAR(255)   NOT NULL DEFAULT '', #标签?#xx?#
  `type`      TINYINT        NOT NULL DEFAULT 0, #0-一般，1-管理员加精华
  `status`    TINYINT        NOT NULL DEFAULT 0, #0-正常，1-不可回复2不可查看
  `views`     INT            NOT NULL DEFAULT 0,
  `replys`    INT            NOT NULL DEFAULT 0,
  `created`   TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated`   TIMESTAMP      NOT NULL DEFAULT 0,
  `lastreply` TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP, #最后回复时间
  PRIMARY KEY (`tid`),
  INDEX `post_last` (`cid`, `lastreply`),
  INDEX `post_user` (`uid`),
  FOREIGN KEY `f_post_cate`(`cid`) REFERENCES `cate` (`cid`)
    ON DELETE CASCADE,
  FOREIGN KEY `f_post_user`(`uid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE
);

#大量数据少用外键
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id`      INT            NOT NULL AUTO_INCREMENT,
  `tid`     INT            NOT NULL, #帖子id
  `pid`     INT            NOT NULL DEFAULT 0, #父评论id 0-顶层回复0+楼中楼且pid为被回复对象
  `uid`     INT            NOT NULL,
  `tuid`    INT            NOT NULL, #回复对象uid(易于查找回复我的)
  `content` VARCHAR(10000) NOT NULL DEFAULT '',
  `replys`  INT            NOT NULL DEFAULT 0, #楼中楼大于0表示有楼中楼回复
  `created` TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` TIMESTAMP      NOT NULL DEFAULT '1970-01-01 00:00:01',
  PRIMARY KEY (`id`),
  INDEX `comment_tid` (`tid`),
  INDEX `comment_user` (`uid`),
  INDEX `comment_pid` (`pid`),
  FOREIGN KEY `f_comment_tid`(`tid`) REFERENCES `post` (`tid`)
    ON DELETE CASCADE,
  FOREIGN KEY `f_comment_user`(`uid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE
);

DROP TABLE IF EXISTS `star`;
CREATE TABLE `star` (
  `id`      INT       NOT NULL AUTO_INCREMENT,
  `uid`     INT       NOT NULL,
  `tid`     INT       NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `star_unique` (`uid`, `tid`),
  FOREIGN KEY `f_star_user`(`uid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE
);

#关注
DROP TABLE IF EXISTS `follow`;
CREATE TABLE `follow` (
  `id`      INT       NOT NULL AUTO_INCREMENT,
  `uid`     INT       NOT NULL,
  `tuid`    INT       NOT NULL, #对方uid
  `note`    VARCHAR(25), #备注名
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `follow_me_t` (`uid`, `tuid`),
  INDEX `follow_rev` (`tuid`),
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE,
  FOREIGN KEY (`tuid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE
);

#聊天表
DROP TABLE IF EXISTS `chat`;
CREATE TABLE `chat` (
  `id`      INT           NOT NULL AUTO_INCREMENT,
  `uid`     INT           NOT NULL,
  `tuid`    INT           NOT NULL, #对方uid
  `content` VARCHAR(2000) NOT NULL DEFAULT '',
  `isread`  TINYINT       NOT NULL DEFAULT '0',
  `created` TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `chat_my_uid` (`uid`),
  INDEX `chat_my_touid` (`tuid`),
  FOREIGN KEY (`uid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE,
  FOREIGN KEY (`tuid`) REFERENCES `user` (`uid`)
    ON DELETE CASCADE
);

#---------------------------------------------------
#------------------触发器定义------------------------
#---------------只更新计数和经验值--------------------
#发表文章
DROP TRIGGER IF EXISTS `t_post_add`;
CREATE TRIGGER `t_post_add`
AFTER INSERT ON `post`
FOR EACH ROW
  BEGIN
    UPDATE `user`
    SET `exp` = `exp` + 3, `posts` = `posts` + 1
    WHERE `uid` = new.uid;
    UPDATE `cate`
    SET `posts` = `posts` + 1, `lastreply` = new.created, `todayposts` = `todayposts` + 1
    WHERE `cid` = new.cid;
  END;

#删除文章
DROP TRIGGER IF EXISTS `t_post_del`;
CREATE TRIGGER `t_post_del`
AFTER DELETE ON `post`
FOR EACH ROW
  BEGIN
    UPDATE `user`
    SET `exp` = `exp` - 3, `posts` = `posts` - 1
    WHERE `uid` = old.uid;

    SET @today = 0;
    IF (date(old.created) = date(now()))
    THEN
      SET @today = 1;
    END IF;

    UPDATE `cate`
    SET `posts` = `posts` - 1, `todayposts` = `todayposts` - @today
    WHERE `cid` = old.cid;
  END;

#增加评论之后
DROP TRIGGER IF EXISTS `t_comment_add`;
CREATE TRIGGER `t_comment_add`
AFTER INSERT ON `comment`
FOR EACH ROW
  BEGIN
    UPDATE `user`
    SET `exp` = `exp` + 1, `replys` = `replys` + 1
    WHERE `uid` = new.uid; #增加经验值

    UPDATE `post`
    SET `replys` = `replys` + 1, `lastreply` = new.created
    WHERE `tid` = new.tid;

    IF new.pid > 0 #楼中楼回复
    THEN
      UPDATE `comment`
      SET `replys` = `replys` + 1
      WHERE id = new.pid;

      #楼中楼 通知楼主
      UPDATE user
      SET newreplys = newreplys + 1
      WHERE uid = (SELECT uid
                   FROM post
                   WHERE tid = new.tid)
            AND uid <> new.tuid;
    END IF;

    #通知
    IF new.uid <> new.tuid
    THEN
      UPDATE user
      SET newreplys = newreplys + 1
      WHERE uid = new.tuid;
    END IF;
  END;

#删除评论
DROP TRIGGER IF EXISTS `t_comment_del`;
CREATE TRIGGER `t_comment_del`
AFTER DELETE ON `comment`
FOR EACH ROW
  BEGIN
    UPDATE `user`
    SET `exp` = `exp` - 1, `replys` = `replys` - 1
    WHERE `uid` = old.uid;

    UPDATE `post`
    SET `replys` = `replys` - 1
    WHERE `tid` = old.tid AND `replys` > 0;

    IF new.pid > 0 #楼中楼回复
    THEN
      UPDATE `comment`
      SET `replys` = `replys` - 1
      WHERE id = new.pid AND `replys` > 0;
    END IF;
  END;

#增加chat
DROP TRIGGER IF EXISTS `t_chat_add`;
CREATE TRIGGER `t_chat_add`
AFTER INSERT ON `chat`
FOR EACH ROW
  IF new.uid <> new.tuid
  THEN
    UPDATE user
    SET newchats = newchats + 1
    WHERE uid = new.tuid;
  END IF;

#删除chat
DROP TRIGGER IF EXISTS `t_chat_del`;
CREATE TRIGGER `t_chat_del`
AFTER DELETE ON `chat`
FOR EACH ROW
  BEGIN
    IF old.isread = 0
    THEN
      UPDATE user
      SET newchats = newchats - 1
      WHERE uid = old.tuid AND newchats > 0;
    END IF;
  END;

#增加关注
DROP TRIGGER IF EXISTS `t_follow_add`;
CREATE TRIGGER `t_follow_add`
AFTER INSERT ON `follow`
FOR EACH ROW
  IF new.uid <> new.tuid
  THEN
    UPDATE user
    SET newfollows = newfollows + 1
    WHERE uid = new.tuid;
  END IF;

#---------------------------------------------------
#--------------计划任务定义--------------------------
#---------------------------------------------------
#查询手否开启 show variables like '%event%';
#开启 SET PERSIST GLOBAL event_scheduler = ON;
DROP EVENT IF EXISTS `event_update_new`;
CREATE EVENT `event_update_new`
  ON SCHEDULE EVERY 1 DAY STARTS '2017-01-01 00:00:00'
  ON COMPLETION PRESERVE
ENABLE
DO UPDATE `cate`
SET `todayposts` = 0
WHERE 1;

#---------------------------------------------------
#---------------------存储过程定义-------------------
#---------------------------------------------------
#添加评论回复楼主
DROP PROCEDURE IF EXISTS `comment_add_lz`;
CREATE PROCEDURE comment_add_lz(
  IN in_tid     INT,
  IN in_uid     INT,
  IN in_content VARCHAR(10000))
  BEGIN
    SELECT uid
    INTO @tuid
    FROM `post`
    WHERE tid = in_tid;
    IF @tuid IS NOT NULL
    THEN
      INSERT INTO `comment` (`tid`, `uid`, `tuid`, `content`)
      VALUES (in_tid, in_uid, @tuid, in_content);
    END IF;
  END;

#楼中楼回复
DROP PROCEDURE IF EXISTS `comment_add_cz`;
CREATE PROCEDURE comment_add_cz(
  IN in_tid     INT,
  IN in_pid     INT,
  IN in_uid     INT,
  IN in_content VARCHAR(10000))
  BEGIN
    SELECT @czuid := `uid`
    FROM `comment`
    WHERE `tid` = in_tid AND `id` = in_pid;

    IF @czuid IS NOT NULL
    THEN
      INSERT INTO `comment` (`tid`, `pid`, `uid`, `tuid`, `content`)
      VALUES (in_tid, in_pid, in_uid, @czuid, in_content);
    END IF;
  END;

#删除评论
DROP PROCEDURE IF EXISTS `comment_del`;
CREATE PROCEDURE comment_del(IN in_uid INT, IN in_id INT)
  BEGIN
    SELECT
      `pid`,
      `replys`
    INTO @pid, @replys
    FROM `comment`
    WHERE `id` = in_id AND `uid` = in_uid;
    IF FOUND_ROWS() > 0
    THEN
      DELETE FROM `comment`
      WHERE `id` = in_id;
      IF @pid = 0 AND @replys > 0 #父评论且有子回复
      THEN
        DELETE FROM comment
        WHERE pid = in_id;
      END IF;
    END IF;
  END;

#删除评论
DROP PROCEDURE IF EXISTS `comment_del_admin`;
CREATE PROCEDURE comment_del_admin(IN in_id INT)
  BEGIN
    SELECT
      `pid`,
      `replys`
    INTO @pid, @replys
    FROM `comment`
    WHERE `id` = in_id;
    IF FOUND_ROWS() > 0
    THEN
      DELETE FROM `comment`
      WHERE `id` = in_id;
      IF @pid = 0 AND @replys > 0 #父评论且有子回复
      THEN
        DELETE FROM comment
        WHERE pid = in_id;
      END IF;
    END IF;
  END;
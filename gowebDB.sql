-- ----------based on postgres sql-------------------
-- ----------------数据表定义--------------------------
-- -------------------v3.0---------------------------

-- LEFT(in_content, 50) 截断字符串

-- 用户表
DROP TABLE IF EXISTS "user";
CREATE TABLE "user" (
  id          SERIAL PRIMARY KEY,
  username    VARCHAR(32) NOT NULL UNIQUE,
  password    VARCHAR(64) NOT NULL,
  email       VARCHAR(32) NOT NULL UNIQUE,
  status      SMALLINT    NOT NULL DEFAULT 0, -- 0-ok 1-block
  sex         SMALLINT    NOT NULL DEFAULT 0, -- 0-unknown 1-man 2-woman
  exp         INT         NOT NULL DEFAULT 0, -- 经验值
  birthday    DATE,
  phone       TEXT,
  description TEXT, -- 个人简介
  site        TEXT, -- 个人网站
  posts       INT         NOT NULL DEFAULT 0, -- 发帖数
  replys      INT         NOT NULL DEFAULT 0, -- 回复数
  regtime     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 用户认证表
DROP TABLE IF EXISTS auth;
CREATE TABLE auth (
  id         SERIAL PRIMARY KEY,
  uid        INT          NOT NULL REFERENCES "user" (id)
    ON DELETE CASCADE,
  type       SMALLINT     NOT NULL, -- 认证类型 1-账号 2-qq 3-微信 4-微博 5-github
  identifier VARCHAR(64)  NOT NULL, -- 自己应用可以不存,第三方应用的唯一标识,微信用户名，GitHub用户名
  credential VARCHAR(128) NOT NULL UNIQUE, -- 密码凭证 token或者第三方token 更改密码过后更改token
  created    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP -- 时间
);

-- 分类表
DROP TABLE IF EXISTS category;
CREATE TABLE category (
  id          SERIAL PRIMARY KEY,
  name        VARCHAR(64) NOT NULL UNIQUE, -- 版块名字
  description TEXT,
  sticks      INT [], -- 置顶 tid...
  posts       INT         NOT NULL DEFAULT 0, -- 版块帖子计数
  created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 帖子表
DROP TABLE IF EXISTS post;
CREATE TABLE post (
  id        SERIAL PRIMARY KEY,
  cid       SMALLINT     NOT NULL REFERENCES category (id)
    ON DELETE CASCADE, -- 板块id
  uid       INT          NOT NULL REFERENCES "user" (id)
    ON DELETE CASCADE, -- 用户id
  username  VARCHAR(32)  NOT NULL, -- 牺牲这个加快查找速度
  title     VARCHAR(128) NOT NULL, -- 标题
  content   TEXT         NOT NULL, -- 内容
  tags      TEXT [], -- 标签?-- xx?--
  type      SMALLINT     NOT NULL DEFAULT 0, -- 0-一般 1-管理员加精华
  status    SMALLINT     NOT NULL DEFAULT 0, -- 0-正常 1-不可回复 2-不可查看
  views     INT          NOT NULL DEFAULT 0, -- 点击数
  replys    INT          NOT NULL DEFAULT 0, -- 回复数
  created   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 发表时间
  updated   TIMESTAMP, -- 编辑时间
  lastreply TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP -- 最后回复时间 用于排序
);


CREATE INDEX i_post_last ON post (lastreply);

-- 评论表
DROP TABLE IF EXISTS comment;
CREATE TABLE comment (
  id      SERIAL PRIMARY KEY,
  tid     INT       NOT NULL, -- 帖子id
  pid     INT       NOT NULL DEFAULT 0, -- 父评论id 0-顶层回复0+楼中楼且pid为被回复对象
  uid     INT       NOT NULL,
  tuid    INT       NOT NULL, -- 回复对象uid(易于查找回复我的)
  content TEXT      NOT NULL DEFAULT '',
  replys  INT       NOT NULL DEFAULT 0, -- 楼中楼大于0表示有楼中楼回复
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP,

  CONSTRAINT f_comment_tid FOREIGN KEY (tid)
  REFERENCES post (tid)
    ON DELETE CASCADE,
  CONSTRAINT f_comment_user FOREIGN KEY (uid)
  REFERENCES user (id)
    ON DELETE CASCADE
);

CREATE INDEX i_comment_tid
  ON comment (tid);
CREATE INDEX i_comment_user
  ON comment (uid);

-- 收藏表
DROP TABLE IF EXISTS star;
CREATE TABLE star (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL,
  tid     INT       NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_star_user FOREIGN KEY (uid)
  REFERENCES user (id)
    ON DELETE CASCADE
);

CREATE INDEX i_star_user
  ON star (uid);
CREATE UNIQUE INDEX i_star_unique
  ON star (uid, tid);

-- 关注表
DROP TABLE IF EXISTS follow;
CREATE TABLE follow (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL,
  tuid    INT       NOT NULL, -- 对方uid
  note    TEXT, -- 备注名
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_follow_user FOREIGN KEY (uid)
  REFERENCES user (id)
    ON DELETE CASCADE,
  CONSTRAINT f_follow_tuser FOREIGN KEY (tuid)
  REFERENCES user (id)
    ON DELETE CASCADE
);

CREATE INDEX i_follow_me
  ON follow (uid);
CREATE INDEX i_follow_other
  ON follow (tuid);
CREATE UNIQUE INDEX i_follow_unique
  ON follow (uid, tuid);

-- 聊天表
DROP TABLE IF EXISTS chat;
CREATE TABLE chat (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL, -- 我的uid 发送方
  tuid    INT       NOT NULL, -- 对方uid
  content TEXT      NOT NULL,
  isread  BOOLEAN   NOT NULL DEFAULT FALSE,
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_chat_user FOREIGN KEY (uid)
  REFERENCES user (id)
    ON DELETE CASCADE,
  CONSTRAINT f_chat_tuser FOREIGN KEY (tuid)
  REFERENCES user (id)
    ON DELETE CASCADE
);

CREATE INDEX i_chat_me
  ON chat (uid);
CREATE INDEX i_chat_other
  ON chat (tuid);

-- ------------------触发器定义------------------------
-- ---------------只更新计数和经验值--------------------

-- 发表帖子触发器
CREATE OR REPLACE FUNCTION func_post_add() RETURNS TRIGGER AS $$
BEGIN
UPDATE user
SET exp = exp + 3, posts = posts + 1
WHERE id = new.uid;

UPDATE category
SET posts = posts + 1
WHERE id = new.cid;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_post_add
  AFTER INSERT
  ON post
  FOR EACH ROW
  EXECUTE PROCEDURE func_post_add();

-- 删除帖子触发器
CREATE OR REPLACE FUNCTION func_post_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE user
SET exp = exp - 3, posts = posts - 1
WHERE id = old.id;

UPDATE category
SET posts = posts - 1
WHERE id = old.id;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_post_del
  AFTER DELETE
  ON post
  FOR EACH ROW
  EXECUTE PROCEDURE func_post_del();

-- 增加评论触发器
CREATE OR REPLACE FUNCTION func_comment_add() RETURNS TRIGGER AS $$
BEGIN
UPDATE user
SET exp = exp + 1, replys = user.replys + 1
WHERE id = new.id;

UPDATE post
SET replys = replys + 1, lastreply = new.created
WHERE tid = new.tid;

IF new.pid > 0 -- 楼中楼回复
THEN
UPDATE comment
SET replys = replys + 1
WHERE id = new.pid;
-- TODO 通知
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_comment_add
  AFTER INSERT
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE func_comment_add();

-- 删除评论触发器
CREATE OR REPLACE FUNCTION func_comment_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE user
SET exp = exp - 1, replys = replys - 1
WHERE id = old.uid;

UPDATE post
SET replys = replys - 1
WHERE tid = old.tid AND replys > 0;

IF new.pid > 0 THEN -- 楼中楼回复
UPDATE comment
SET replys = replys - 1
WHERE id = old.pid AND replys > 0;
ELSEIF old.replys>0 THEN -- 不是楼中楼回复且有子回复删除他们
DELETE FROM comment
WHERE pid = old.id;
END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_comment_del
  AFTER DELETE
  ON comment
  FOR EACH ROW
  EXECUTE PROCEDURE func_comment_del();

-- 增加关注触发器
CREATE OR REPLACE FUNCTION func_follow_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE user
SET newfollows = newfollows + 1
WHERE id = new.tuid;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_follow_add
  AFTER INSERT
  ON follow
  FOR EACH ROW
WHEN (new.uid <> new.tuid)
EXECUTE PROCEDURE func_follow_del();

-- ---------------------常见函数-------------------
-- -----------------------------------------------
-- 添加评论回复楼主
CREATE OR REPLACE FUNCTION comment_add_lz(in_tid INT, in_uid INT, in_content TEXT ) RETURNS BOOLEAN AS $$
DECLARE tuid INT;
BEGIN
SELECT uid
INTO tuid
FROM post
WHERE tid = in_tid;
IF tuid IS NOT NULL THEN
INSERT INTO COMMENT (tid, uid, tuid, content) VALUES (in_tid, in_uid, tuid, in_content) RETURNING id;
RETURN TRUE;
END IF;
RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

-- 添加评论楼中楼回复
CREATE OR REPLACE FUNCTION comment_add_cz(in_tid INT, in_pid INT, in_uid INT, in_content TEXT ) RETURNS BOOLEAN AS $$
DECLARE czuid INT;
DECLARE result BOOLEAN;
BEGIN
SET result = FALSE;
SELECT uid
INTO czuid
FROM comment
WHERE tid = in_tid AND id = in_pid;

IF czuid IS NOT NULL THEN
INSERT INTO COMMENT (tid, pid, uid, tuid, content) VALUES (in_tid, in_pid, in_uid, czuid, in_content) RETURNING id;
SET result = TRUE;
END IF;
RETURN result;
END;
$$ LANGUAGE plpgsql;
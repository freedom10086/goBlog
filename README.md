#GO博客系统
> 一个简单的go语言博客,依赖仅有postgresql数据库
## 配置psql
1. 设置并启动psql
    1. windows
    ```
    .\initdb.exe -D C:\app\pgsql-9.6.2\data -E UTF8
    .\pg_ctl -D "C:\app\pgsql-9.6.2\data" -l logfile start
    ```
    2. linux
    ```
    sudo apt-get install postgresql-client
    sudo apt-get install postgresql
    su postgress //切换到postgress用户
    psql//使用psql命令登录PostgreSQL控制台，以postgress用户
    \password postgres //更改密码
    ```
2. 创建用户 超级用户要谨慎
`createuser [--superuser] yang`

3. 设置密码 用postgres超级管理员登陆
> 为刚刚创建的用户设置密码

`psql postgres`

`\password yang`

`\q`

4. 创建数据库 指定所有权yang
`createdb -O yang -p 5432 goblog`

5. 登陆数据库
`psql -U yang -d goblog -h 127.0.0.1 -p 5432`

> -U 用户
> -d 数据库
> -h host
> -p 端口

## 数据表定义

```sql
-- 用户表
DROP TABLE users;
CREATE TABLE users (
  uid         SERIAL PRIMARY KEY,
  username    VARCHAR(32) NOT NULL UNIQUE,
  password    VARCHAR(64) NOT NULL,
  email       VARCHAR(32) NOT NULL UNIQUE,
  status      SMALLINT    NOT NULL DEFAULT 0, -- 0-ok 1-block
  sex         SMALLINT    NOT NULL DEFAULT 0, -- 0-unknown 1-man 2-woman
  exp         INT         NOT NULL DEFAULT 0, -- 经验值
  birthday    DATE,
  phone       TEXT,
  description TEXT,
  site        TEXT, -- 个人网站
  posts       INT         NOT NULL DEFAULT 0, -- 发帖数
  replys      INT         NOT NULL DEFAULT 0, -- 回复数
  newreplys   INT         NOT NULL DEFAULT 0, -- 新回复数
  newfollows  INT         NOT NULL DEFAULT 0, -- 新粉丝数
  newchats    INT         NOT NULL DEFAULT 0, -- 新私信数
  regtime     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX i_user_name
  ON users (username);

CREATE INDEX i_user_email
  ON users (email);

-- 分类表
DROP TABLE cate;
CREATE TABLE cate (
  cid         SERIAL PRIMARY KEY,
  name        VARCHAR(32) NOT NULL UNIQUE, -- 版块名字
  description TEXT,
  sticks      INT [], -- 置顶 tid...
  posts       INT         NOT NULL DEFAULT 0, -- 版块帖子计数
  created     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX i_cate_name
  ON cate (name);

-- 帖子表
DROP TABLE post;
CREATE TABLE post (
  tid       SERIAL PRIMARY KEY,
  cid       SMALLINT    NOT NULL,
  uid       INT         NOT NULL,
  username  VARCHAR(32) NOT NULL, -- 牺牲这个加快查找速度
  title     VARCHAR(64) NOT NULL,
  content   TEXT        NOT NULL,
  tags      TEXT [], -- 标签?-- xx?--
  type      SMALLINT    NOT NULL DEFAULT 0, -- 0-一般，1-管理员加精华
  status    SMALLINT    NOT NULL DEFAULT 0, -- 0-正常，1-不可回复2不可查看
  views     INT         NOT NULL DEFAULT 0,
  replys    INT         NOT NULL DEFAULT 0,
  created   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated   TIMESTAMP,
  lastreply TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 最后回复时间

  CONSTRAINT f_post_cate FOREIGN KEY (cid)
  REFERENCES cate (cid)
    ON DELETE CASCADE,

  CONSTRAINT f_post_user FOREIGN KEY (uid)
  REFERENCES users (uid)
    ON DELETE CASCADE
);

CREATE INDEX i_post_user
  ON post (uid);

CREATE INDEX i_post_last
  ON post (cid, lastreply);

-- 评论表
DROP TABLE comment;
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
  REFERENCES users (uid)
    ON DELETE CASCADE
);

CREATE INDEX i_comment_tid
  ON comment (tid);
CREATE INDEX i_comment_user
  ON comment (uid);

-- 收藏表
DROP TABLE star;
CREATE TABLE star (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL,
  tid     INT       NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_star_user FOREIGN KEY (uid)
  REFERENCES users (uid)
    ON DELETE CASCADE
);

CREATE INDEX i_star_user
  ON star (uid);
CREATE UNIQUE INDEX i_star_unique
  ON star (uid, tid);

-- 关注表
DROP TABLE follow;
CREATE TABLE follow (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL,
  tuid    INT       NOT NULL, -- 对方uid
  note    TEXT, -- 备注名
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_follow_user FOREIGN KEY (uid)
  REFERENCES users (uid)
    ON DELETE CASCADE,
  CONSTRAINT f_follow_tuser FOREIGN KEY (tuid)
  REFERENCES users (uid)
    ON DELETE CASCADE
);

CREATE INDEX i_follow_me
  ON follow (uid);
CREATE INDEX i_follow_other
  ON follow (tuid);
CREATE UNIQUE INDEX i_follow_unique
  ON follow (uid, tuid);

-- 聊天表
DROP TABLE chat;
CREATE TABLE chat (
  id      SERIAL PRIMARY KEY,
  uid     INT       NOT NULL, -- 我的uid 发送方
  tuid    INT       NOT NULL, -- 对方uid
  content TEXT      NOT NULL,
  isread  BOOLEAN   NOT NULL DEFAULT FALSE,
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT f_chat_user FOREIGN KEY (uid)
  REFERENCES users (uid)
    ON DELETE CASCADE,
  CONSTRAINT f_chat_tuser FOREIGN KEY (tuid)
  REFERENCES users (uid)
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
UPDATE users
SET exp = exp + 3, posts = posts + 1
WHERE uid = new.uid;

UPDATE cate
SET posts = posts + 1
WHERE cid = new.cid;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_post_add
AFTER INSERT ON post
FOR EACH ROW
  EXECUTE PROCEDURE func_post_add();

-- 删除帖子触发器
CREATE OR REPLACE FUNCTION func_post_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE users
SET exp = exp - 3, posts = posts - 1
WHERE uid = old.uid;

UPDATE cate
SET posts = posts - 1
WHERE cid = old.cid;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_post_del
AFTER DELETE ON post
FOR EACH ROW
  EXECUTE PROCEDURE func_post_del();

-- 增加评论触发器
CREATE OR REPLACE FUNCTION func_comment_add() RETURNS TRIGGER AS $$
BEGIN
UPDATE users
SET exp = exp + 1, replys = users.replys + 1
WHERE uid = new.uid;

UPDATE post
SET replys = replys + 1, lastreply = new.created
WHERE tid = new.tid;

IF new.pid > 0 -- 楼中楼回复
THEN
UPDATE comment
SET replys = replys + 1
WHERE id = new.pid;

-- 楼中楼 通知楼主
UPDATE users
SET newreplys = newreplys + 1
WHERE uid <> new.tuid AND uid = (SELECT uid
                                 FROM post
                                 WHERE tid = new.tid);
END IF;
-- 通知
IF new.uid <> new.tuid
THEN
UPDATE users
SET newreplys = newreplys + 1
WHERE uid = new.tuid;
END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_comment_add
AFTER INSERT ON comment
FOR EACH ROW
  EXECUTE PROCEDURE func_comment_add();

-- 删除评论触发器
CREATE OR REPLACE FUNCTION func_comment_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE users
SET exp = exp - 1, replys = replys - 1
WHERE uid = old.uid;

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
AFTER DELETE ON comment
FOR EACH ROW
  EXECUTE PROCEDURE func_comment_del();

-- 增加聊天触发器
CREATE OR REPLACE FUNCTION func_chat_add() RETURNS TRIGGER AS $$
BEGIN
IF new.uid <> new.tuid
THEN
UPDATE users
SET newchats = newchats + 1
WHERE uid = new.tuid;
END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER t_chat_add ON chat;
CREATE TRIGGER t_chat_add
AFTER INSERT ON chat
FOR EACH ROW
  EXECUTE PROCEDURE func_chat_add();

-- 删除聊天触发器
CREATE OR REPLACE FUNCTION func_chat_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE users
SET newchats = newchats - 1
WHERE uid = old.tuid AND newchats > 0;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER t_chat_del
AFTER DELETE ON chat
FOR EACH ROW
WHEN (old.isread = FALSE )
EXECUTE PROCEDURE func_chat_del();

-- 增加关注触发器
CREATE OR REPLACE FUNCTION func_follow_del() RETURNS TRIGGER AS $$
BEGIN
UPDATE users
SET newfollows = newfollows + 1
WHERE uid = new.tuid;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS t_follow_add;
CREATE TRIGGER t_follow_add
AFTER INSERT ON follow
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
```

## 依赖及参考
- [Chart.js](http://chartjs.org/)
- [highlight.js](http://git.io/hljslicense)
- [bootstrap](https://getbootstrap.com)
- [marked.js](https://github.com/chjj/marked)

## bug 反馈
- [2351386755@qq.com](mailto://2351386755@qq.com)    
- [yangluo.chn@gmail.com](mailto://yangluo.chn@gmail.com)

## License
```
免费使用随便传播
```
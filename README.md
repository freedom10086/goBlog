# 我的博客系统文档


## before
>sudo apt install mysql-server
>sudo pip install flask
>sudo pip install pymysql
>CREATE DATABASE MyBlog DEFAULT CHARACTER SET utf8mb4;

------

## 1. 数据库定义

>评论表 `comments`

- `ID` 唯一id 自增
- `post_ID` 哪一篇文章的id
- `author_ID` 发表者的id
- `author_IP` 发表者的ip
- `time` 发表的事件
- `content` 发表的内容
- `agent` 发表者的user-agent
- `type` 发表评论的类型
- `parent` 发表评论的父级评论（楼中楼）或者为 0
- 主码`ID`

---

>网站设置表 `options`

- `ID` 唯一id 自增
- `name` 设置名
- `value` 设置值
- `des` 设置描述
- 主码 `ID`

---
>所有文章表  `posts`

- `ID` 唯一id 自增
- `fid` 论坛所属板块
- `title` 文章标题
- `content` 文章内容
- `user_ID` 发表者id
- `user_name` 发表者用户名
- `time` 发表时间
- `modified_time` 最后修改时间
- `comment_count` 评论计数
- `like_count` 喜欢计数
- `status` 目前状态 0-开放 1-关闭
- `comment_status` 评论状态 0-开启 1-关闭
- `tages` 文章标签 多个用`,`分割
- 主码 `ID`

---
>所有用户表  `users`

- `ID` 唯一id 自增
- `username` 用户名
- `password` 密码
- `email` 邮件
- `email_status` 邮件验证状态 0-未 1-验证
- `registered` 注册时间
- `lastlogin` 最后登录时间
- `sites` 用户主页
- `sex` 性别 0-保密 1-男 2-女
- `avatar` 头像地址
- `description` 用户签名
- `exp` 用户积分
- `status` 用户状态 0-正常 1-限制登陆
- `group` 用户组 0-普通 1-管理员 *todo more*
- `birthday` 生日
- `number` 电话
- `position` 住址
- 主码 `ID`
- 唯一码 `username`

---
>用户关系表 `user_relation`
>用于表示拥护之间的关系 如好友/关注

- `ID` 唯一id 自增
- `uid` 用户id
- `touid` 对方用户id
- `relation` 关系类型 0-好友 1-关注
- `description` 关系描述
- 主码 (`uid`,`touid`)

---
>用户消息表 `message`

- `ID` 唯一id 自增
- `touid` 用户id
- `from_uid` 消息发送方id
- `from_username` 消息发送方用户名
- `time` 发送时间
- `message_title` 发送消息标题
- `message` 消息正文
- `type` 消息种类 0-评论提醒 1-好友消息 2-系统消息
- `status` 消息状态 0-未读 1-已读
- `fromID` 消息来源id 评论则评论id 以此类推
- 主码 `ID`

---
>用户收藏表 `star`

- `ID` 唯一id 自增
- `user` 用户id
- `time` 收藏时间
- `type` 收藏类型 *todo default 0*
- `star_id` 收藏的id 文章就是文章的id
- `description` 收藏的描述 可以用作备注
- 主码 `ID`

---
>用户签到表 `sign`

- `ID` 唯一id 自增
- `user` 用户id
- `time` 签到时间
- `sign_prize` 签到获得的奖励积分
- `sign_des` 签到的描述 想说的话
- `ip` 签到ip
- 主码 `ID`

---
>首页表 `portal`

- `ID` 唯一id 自增
- `type` 类型 0-一般文章 *todo*
- `from_id` 首页文章的id
- `img` 图片 没有就空
- `title` 标题
- `description` 描述
- 主码 `ID`

---
## 2. 接口定义
*todo*

---
##附录数据库创建
```
CREATE TABLE `comments` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `post_ID` int(11) unsigned NOT NULL DEFAULT '0',
  `author_ID` int(11) unsigned NOT NULL DEFAULT '0',
  `author_IP` varchar(50) NOT NULL DEFAULT '',
  `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `content` text NOT NULL,
  `agent` varchar(200) NOT NULL DEFAULT '',
  `type` varchar(20) NOT NULL DEFAULT '',
  `parent` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `options` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `value` varchar(50) NOT NULL,
  `des` varchar(50) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `posts` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `fid` mediumint(8) unsigned NOT NULL DEFAULT '0',
  `title` varchar(100) NOT NULL DEFAULT '',
  `content` longtext NOT NULL,
  `user_ID` int(11) unsigned NOT NULL DEFAULT '0',
  `user_name` varchar(20) NOT NULL DEFAULT '',
  `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `modified_time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `comment_count` int(11) NOT NULL DEFAULT '0',
  `like_count` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `comment_status` tinyint(1) NOT NULL DEFAULT '0',
  `tages` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `users` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL DEFAULT '',
  `password` varchar(20) NOT NULL DEFAULT '',
  `email` varchar(50) NOT NULL DEFAULT '',
  `email_status` tinyint(1) NOT NULL DEFAULT '0',
  `registered` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `lastlogin` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `sites` varchar(50) NOT NULL DEFAULT '',
  `sex` tinyint(1) DEFAULT '0',
  `avatar` varchar(50) NOT NULL DEFAULT '',
  `description` varchar(100) DEFAULT NULL,
  `exp` int(11) unsigned NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `group` tinyint(1) NOT NULL DEFAULT '0',
  `birthday` date DEFAULT NULL DEFAULT '0000-00-00 00:00:00',
  `number` varchar(20) NOT NULL DEFAULT '0',
  `position` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_relation` (
  `ID` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `touid` int(11) UNSIGNED NOT NULL DEFAULT '0',
  `relation` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`uid`,`touid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `message` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `touid` int(11) unsigned NOT NULL DEFAULT '0',
  `from_uid` int(11) unsigned NOT NULL DEFAULT '0',
  `from_username` varchar(20) NOT NULL DEFAULT '',
  `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `message_title` varchar(100) NOT NULL DEFAULT '',
  `message` text NOT NULL,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `fromID` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



CREATE TABLE `star` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(11) unsigned NOT NULL,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `star_id` int(11) unsigned NOT NULL DEFAULT '0',
  `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `sign` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(11) unsigned NOT NULL DEFAULT '0',
  `time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `sign_prize` varchar(50) NOT NULL DEFAULT '',
  `sign_des` varchar(50) NOT NULL DEFAULT '',
  `ip` varchar(20) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `portal` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL DEFAULT '0',
  `from_id` int(11) NOT NULL DEFAULT '0',
  `img` varchar(50) NOT NULL DEFAULT '',
  `title` varchar(50) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```
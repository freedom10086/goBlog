package models

import (
	"time"
)

/*
CREATE TABLE `comments` (
  `cid` integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tid` integer NOT NULL DEFAULT 0,
  `uid` integer NOT NULL DEFAULT 0,
  `pid` integer NOT NULL DEFAULT 0,
  `ip` vachar(30) NOT NULL DEFAULT '',
  `agent` varchar(200) NOT NULL DEFAULT '',
  `created` datetime NOT NULL,
  `updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `content` varchar(5000) NOT NULL DEFAULT '' ,
  `author` varchar(30) NOT NULL DEFAULT ''
) ENGINE=InnoDB;
CREATE INDEX `comments_tid` ON `comments` (`tid`);
CREATE INDEX `comments_pid` ON `comments` (`pid`);
CREATE INDEX `comments_uid` ON `comments` (`uid`);
*/

type Comment struct {
	Cid     int
	Tid     int
	Uid     int
	Pid     int //父评论 0-一般 >0 楼中楼
	Ip      string
	Agent   string
	Created time.Time
	Updated time.Time
	Content string
	Author  string
}

package models

import (
	"time"
)

/*
 */

type Star struct {
	Sid     int       //id
	Uid     int       //用户id
	Tid     int       //收藏帖子id
	Title   string    //收藏帖子标题
	Created time.Time //时间
}

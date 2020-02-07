package repository

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"goBlog/conf"
	"net/smtp"
	"strings"
	"time"
)

//发送邮箱
//可以群发 to;to1;to2
//http://www.cnblogs.com/linecheng/p/5861468.html
/**
From:sender_user@demo.net
To:to_user@demo.net
Subject:这是主题
Mime-Version：1.0 //通常是1.0
Content-Type：Multipart/mixed;Boundary="THIS_IS_BOUNDARY_JUST_MAKE_YOURS" //boundary为分界字符，跟http传文件时类似
Date:当前时间


--THIS_IS_BOUNDARY_JUST_MAKE_YOURS         //boundary前边需要加上连接符 -- ， 首部和第一个boundary之间有两个空行
Content-Type:text/plain;chart-set=utf-8
                                            //单个部分的首部和正文间有个空行
这是正文1
这是正文2

--THIS_IS_BOUNDARY_JUST_MAKE_YOURS                  //每个部分的与上一部分之间一个空行
Content-Type：image/jpg;name="test.jpg"
Content-Transfer-Encoding:base64
Content-Description:这个是描述
                                            //单个部分的首部和正文间有个空行
base64编码的文件                              //文件内容使用base64 编码，单行不超过80字节，需要插入\r\n进行换行
--THIS_IS_BOUNDARY_JUST_MAKE_YOURS--        //最后结束的标识--boundary--
*/
func SendMail(to, subject, content string) error {
	user := conf.Conf.MailUsername
	password := conf.Conf.MailPassword
	host := conf.Conf.MailHost

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	buffer := bytes.NewBuffer(nil)
	boudary := "#####"
	header := fmt.Sprintf("To:%s\r\n"+
		"From:%s\r\n"+
		"Subject:%s\r\n"+
		"Content-Type:Multipart/mixed;Boundary=\"%s\"\r\n"+
		"Mime-Version:1.0\r\n"+
		"Date:%s\r\n\r\n", to, user, subject, boudary, time.Now().String())
	body := fmt.Sprintf("--%s\r\nContent-Type:text/html;chart-set=utf-8\r\n%s\r\n",
		boudary, content) // body 可以循环写入多个
	tail := fmt.Sprintf("--%s--", boudary)

	buffer.WriteString(header)
	buffer.WriteString(body)
	buffer.WriteString(tail)

	return smtp.SendMail(host, auth, user, strings.Split(to, ";"), buffer.Bytes())
}

//存入数据库 md5(password)
func Md5_encode(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return md5pass
}

package repository

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
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

type MailAttach struct {
	name        string // eg a.jpg
	data        []byte
	contentType string // eg image/jpg
}

func SendPlainMail(to, subject, content string) error {
	user := conf.Conf.MailUsername
	password := conf.Conf.MailPassword
	host := conf.Conf.MailHost

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	contentType := "Content-Type: text/plain" + "; charset=utf-8"

	msg := []byte("To: " + to + "\r\n" +
		"From: " + user + "\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n" +
		"Mime-Version:1.0\r\n" +
		fmt.Sprintf("Date:%s\r\n", time.Now().String()) +
		"\r\n" + content)
	tos := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, tos, msg)
	return err
}

func SendHtmlMail(to, subject, content string) error {
	user := conf.Conf.MailUsername
	password := conf.Conf.MailPassword
	host := conf.Conf.MailHost

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	contentType := "Content-Type: text/html" + "; charset=utf-8"

	msg := []byte("To: " + to + "\r\n" +
		"From: " + user + "\r\n" +
		"Subject: " + subject + "\r\n" +
		contentType + "\r\n" +
		"Mime-Version:1.0\r\n" +
		fmt.Sprintf("Date:%s\r\n", time.Now().String()) +
		"\r\n" + content)
	tos := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, tos, msg)
	return err
}

func sendMailWithAttach(to, subject, content string, attachs []MailAttach) error {
	user := conf.Conf.MailUsername
	password := conf.Conf.MailPassword
	host := conf.Conf.MailHost

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	buffer := bytes.NewBuffer(nil)

	boudary := "####"
	header := fmt.Sprintf("To:%s\r\n"+
		"From:%s\r\n"+
		"Subject:%s\r\n"+
		"Content-Type:multipart/mixed;Boundary=\"%s\"\r\n"+
		"Mime-Version:1.0\r\n"+
		"Date:%s\r\n", to, user, subject, boudary, time.Now().String())
	buffer.WriteString(header)

	msg := fmt.Sprintf("\r\n\r\n--"+boudary+"\r\n"+
		"Content-Type:text/html;charset=utf-8\r\n\r\n%s\r\n", content)
	buffer.WriteString(msg)

	for _, attach := range attachs {
		gap := fmt.Sprintf(
			"\r\n--%s\r\n"+
				"Content-Transfer-Encoding: base64\r\n"+
				"Content-Disposition: attachment;\r\n"+
				"Content-Type:%s;name=\"%s\"\r\n", boudary, attach.contentType, attach.name)
		buffer.WriteString(gap)

		base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(attach.data)))
		base64.StdEncoding.Encode(base64Bytes, attach.data)
		buffer.WriteString("\r\n")
		for i, l := 0, len(base64Bytes); i < l; i++ {
			buffer.WriteByte(base64Bytes[i])
			if (i+1)%76 == 0 {
				buffer.WriteString("\r\n")
			}
		}
	}

	buffer.WriteString("\r\n--" + boudary + "--")
	sendto := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendto, buffer.Bytes())
	return err
}

//存入数据库 md5(password)
func Md5_encode(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return md5pass
}

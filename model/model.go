package model

import (
	"strings"
	"net/smtp"
	"fmt"
	"crypto/md5"
	"bytes"
	"time"
)


//发送邮箱
//可以群发 to;to1;to2
//http://www.cnblogs.com/linecheng/p/5861468.html
func SendMail(to, subject, content string) error {
	user := "2351386755@qq.com"
	password := "qjadozmyfdgpdhhj"
	host := "smtp.qq.com:587"

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	buffer := bytes.NewBuffer(nil)
	boudary := "#####"
	header := fmt.Sprintf("To:%s\r\n" +
		"From:%s\r\n" +
		"Subject:%s\r\n" +
		"Content-Type:multipart/mixed;Boundary=\"%s\"\r\n" +
		"Mime-Version:1.0\r\n" +
		"Date:%s\r\n", to, user, subject, boudary, time.Now().String())
	buffer.WriteString(header)
	c := fmt.Sprintf("\r\n" +
		"\r\n" +
		"--%s\r\n" +
		"Content-Type:text/plain;charset=utf-8\r\n" +
		"\r\n" +
		"%s\r\n", boudary, content)
	buffer.WriteString(c)
	return smtp.SendMail(host, auth, user, strings.Split(to, ";"), buffer.Bytes())
}

//存入数据库 md5(password)
func Md5_encode(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return md5pass
}

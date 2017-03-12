package model

import (
	"strings"
	"net/smtp"
	"fmt"
	"crypto/md5"
)


//发送邮箱
func SendMail(to, subject, content string) error {
	user := "2351386755@qq.com"
	password := "StrikeFreedom"
	host := "smtp.qq.com:587"

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	msg := []byte("To:" + to +
		"\r\n" +
		"Subject: " + subject +
		"\r\n\r\n" +
		content)

	err := smtp.SendMail(host, auth, user, strings.Split(to, ";"), msg)
	return err
}

//存入数据库 md5(password)
func Md5_encode(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return md5pass
}

package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/base64"
	"io"
	"crypto/rand"
)

//生成32位md5字串
func GenMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func GenGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GenMd5String(base64.URLEncoding.EncodeToString(b))
}

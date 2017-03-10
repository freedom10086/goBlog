package code

import "errors"

var (
	ERR_PARAMETER = errors.New("less parameter or parameter is invalid!")

	ERR_TOKEN_TIMEOUT = errors.New("token is time out")
	ERR_TOKEN_INVALID = errors.New("token is invalid")

	ERR_LOGIN = errors.New("username or password is error")

	ErrInvalidPara = errors.New("参数错误!")
	ErrNoInsert    = errors.New("操作失败!")
	ErrNoDelete    = errors.New("删除失败!")
	ErrNoUpdate    = errors.New("没有更改!")
	ErrReply       = errors.New("此文章无法回复!")
	ErrLogin       = errors.New("账号异常,你没有权限登陆!")

	CODE_OK    = 200
	CODE_ERROR = 500
)

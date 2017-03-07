package code

import "errors"

var (
	ERR_PARAMETER  = errors.New("less parameter or parameter is invalid!")

	ERR_TOKEN_TIMEOUT = errors.New("token is time out")
	ERR_TOKEN_INVALID = errors.New("token is invalid")

	ERR_LOGIN = errors.New("username or password is error")


	ErrInvalidPara error = errors.New("参数错误!")
	ErrNoInsert    error = errors.New("操作失败!")
	ErrNoDelete    error = errors.New("删除失败!")
	ErrNoUpdate    error = errors.New("没有更改!")
	ErrReply       error = errors.New("此文章无法回复!")
	ErrLogin             = errors.New("账号异常,你没有权限登陆!")
)

package repository

import "errors"

var (
	ErrNoRecord error = errors.New("没有查询到记录")
	ErrNoAffect error = errors.New("没有改变")
	ErrNoAuth   error = errors.New("没有权限")
	ErrParama   error = errors.New("参数错误")
	ErrExcution error = errors.New("执行错误")
)

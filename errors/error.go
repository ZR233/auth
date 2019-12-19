package errors

import "errors"

var (
	ErrRecordNotExist   = errors.New("记录不存在")
	ErrRecordExist      = errors.New("记录已存在")
	ErrPermissionDenied = errors.New("权限不足")
)

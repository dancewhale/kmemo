package repository

import "errors"

var (
	// ErrNotFound 记录不存在
	ErrNotFound = errors.New("record not found")

	// ErrDuplicateKey 唯一键冲突
	ErrDuplicateKey = errors.New("duplicate key")

	// ErrInvalidInput 输入参数无效
	ErrInvalidInput = errors.New("invalid input")

	// ErrConcurrentUpdate 并发更新冲突
	ErrConcurrentUpdate = errors.New("concurrent update conflict")
)

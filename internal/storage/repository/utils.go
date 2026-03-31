package repository

import (
	"gorm.io/gorm"
)

// convertError 统一错误转换
func convertError(err error) error {
	if err == nil {
		return nil
	}
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	if err == gorm.ErrDuplicatedKey {
		return ErrDuplicateKey
	}
	return err
}

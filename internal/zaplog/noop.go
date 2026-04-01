package zaplog

import "go.uber.org/zap"

func Nop() *zap.Logger {
	return zap.NewNop()
}

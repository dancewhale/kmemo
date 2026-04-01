package zaplog

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor(base *zap.Logger) grpc.UnaryClientInterceptor {
	if base == nil {
		base = zap.NewNop()
	}

	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx, requestID := EnsureRequestID(ctx)

		md, _ := metadata.FromOutgoingContext(ctx)
		md = md.Copy()
		md.Set(MetadataRequestIDKey, requestID)
		ctx = metadata.NewOutgoingContext(ctx, md)

		logger := base.With(
			zap.String("request_id", requestID),
			zap.String("grpc_method", method),
			zap.String("target", cc.Target()),
		)

		startedAt := time.Now()
		logger.Debug("grpc call started")
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(startedAt)
		if err != nil {
			logger.Error("grpc call failed", zap.Duration("duration", duration), zap.Error(err))
			return err
		}
		logger.Info("grpc call finished", zap.Duration("duration", duration))
		return nil
	}
}

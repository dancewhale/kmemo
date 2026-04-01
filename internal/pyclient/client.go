package pyclient

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"
	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/config"
	"kmemo/internal/zaplog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the generated KmemoProcessor gRPC client. Business calls stay in services/.
type Client struct {
	conn   *grpc.ClientConn
	api    kmemov1.KmemoProcessorClient
	logger *zap.Logger
	target string
}

// New dials the Python worker. Close when shutting down the host.
func New(ctx context.Context, cfg config.Config) (*Client, error) {
	ctx = zaplog.WithLogger(ctx, zaplog.FromContext(ctx).Named("pyclient"))
	ctx, _ = zaplog.EnsureRequestID(ctx)
	logger := zaplog.FromContext(ctx)

	logger.Info("python client dial started", zap.String("target", cfg.PythonGRPCAddr))

	ctx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.PythonGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(zaplog.UnaryClientInterceptor(logger)),
	)
	if err != nil {
		logger.Error("python client dial failed", zap.Error(err), zap.String("target", cfg.PythonGRPCAddr))
		return nil, fmt.Errorf("dial python grpc: %w", err)
	}
	logger.Info("python client connected", zap.String("target", cfg.PythonGRPCAddr))
	return &Client{
		conn:   conn,
		api:    kmemov1.NewKmemoProcessorClient(conn),
		logger: logger,
		target: cfg.PythonGRPCAddr,
	}, nil
}

// API exposes the raw generated client for future service layers.
// TODO: remove direct exposure once facades in internal/services exist.
func (c *Client) API() kmemov1.KmemoProcessorClient {
	return c.api
}

// Close releases the underlying connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	if c.logger != nil {
		c.logger.Info("python client closing", zap.String("target", c.target))
	}
	return c.conn.Close()
}

var _ io.Closer = (*Client)(nil)

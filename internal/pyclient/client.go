package pyclient

import (
	"context"
	"fmt"
	"io"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the generated KmemoProcessor gRPC client. Business calls stay in services/.
type Client struct {
	conn *grpc.ClientConn
	api  kmemov1.KmemoProcessorClient
}

// New dials the Python worker. Close when shutting down the host.
func New(ctx context.Context, cfg config.Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.PythonGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial python grpc: %w", err)
	}
	return &Client{
		conn: conn,
		api:  kmemov1.NewKmemoProcessorClient(conn),
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
	return c.conn.Close()
}

var _ io.Closer = (*Client)(nil)

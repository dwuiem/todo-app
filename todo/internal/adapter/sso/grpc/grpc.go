package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"time"
	"todo/gen/sso"
)

type Client struct {
	api sso.AuthClient
}

func New(
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: sso.NewAuthClient(cc),
	}, nil
}

func (c *Client) Register(ctx context.Context, request *sso.RegisterRequest) (uuid.UUID, error) {
	const op = "grpc.Register"

	response, err := c.api.Register(ctx, request)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := uuid.Parse(response.GetUserId())
	if err != nil {
		return uuid.Nil, status.Error(codes.Internal, err.Error())
	}

	return id, nil
}

func (c *Client) Login(ctx context.Context, request *sso.LoginRequest) (string, error) {
	const op = "grpc.Login"

	response, err := c.api.Login(ctx, request)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return response.Token, nil
}

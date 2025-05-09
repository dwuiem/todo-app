package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/gen/sso"
	"sso/internal/service/auth"
	"sso/internal/storage"
)

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Service
}

type Service interface {
	Login(
		ctx context.Context,
		username string,
		password string,
		appId int,
	) (token string, err error)
	Register(
		ctx context.Context,
		username string,
		password string,
	) (userId int64, err error)
	IsAdmin(
		ctx context.Context,
		userId int64,
	) (isAdmin bool, err error)
}

func RegisterGRPC(gRPC *grpc.Server, auth Service) {
	sso.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *sso.LoginRequest,
) (*sso.LoginResponse, error) {

	if err := validateLogin(in); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, in.Username, in.Password, int(in.AppId))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {

	if err := validateRegister(in); err != nil {
		return nil, err
	}

	userId, err := s.auth.Register(ctx, in.Username, in.Password)

	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sso.RegisterResponse{
		UserId: userId,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	in *sso.IsAdminRequest,
) (*sso.IsAdminResponse, error) {
	if err := validateIsAdmin(in); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sso.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateRegister(r *sso.RegisterRequest) error {
	if r.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is empty")
	}

	if r.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}

func validateLogin(r *sso.LoginRequest) error {
	if r.GetUsername() == "" {
		return status.Error(codes.InvalidArgument, "username is empty")
	}

	if r.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if r.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "app_id is empty")
	}

	return nil
}

func validateIsAdmin(r *sso.IsAdminRequest) error {
	if r.GetUserId() == 0 {
		return status.Error(codes.InvalidArgument, "userId is empty")
	}
	return nil
}

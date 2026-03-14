package main

import (
	"context"
	"log/slog"

	"github.com/sjingzhao/gong/pb/example"
	"github.com/sjingzhao/gong/pkg/logx"
)

// UserServiceImpl 实现 example.UserServiceServer 接口
type UserServiceImpl struct {
	example.UnimplementedUserServiceServer
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *example.UserRequest) (*example.User, error) {
	logx.FromContext(ctx).InfoContext(ctx, "grpc request", slog.Any("GetUser", req))

	return &example.User{
		Id:    req.UserId,
		Name:  "User_" + string(rune('A'+req.UserId%26)),
		Email: "user" + string(req.UserId) + "@example.com",
	}, nil
}

func (s *UserServiceImpl) AddUser(ctx context.Context, req *example.User) (*example.UserRequest, error) {
	logx.FromContext(ctx).InfoContext(ctx, "grpc request", slog.Any("AddUser", req))
	return &example.UserRequest{UserId: req.Id + 1}, nil
}

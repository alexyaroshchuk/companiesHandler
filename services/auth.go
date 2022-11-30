package services

import (
	"context"

	"companiesHandler/auth"
	"companiesHandler/models"
	"companiesHandler/proto/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	userService *UserService
	jwtManager  *auth.JWTManager
}

func NewAuthServer(userService *UserService, jwtManager *auth.JWTManager) pb.AuthServiceServer {
	return &AuthServer{userService: userService, jwtManager: jwtManager}
}

func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.userService.Get(ctx, &pb.UserRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	token, err := server.jwtManager.Generate(&models.User{
		Username: user.User.Name,
		Role:     user.User.Role,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	return &pb.LoginResponse{AccessToken: token}, nil
}

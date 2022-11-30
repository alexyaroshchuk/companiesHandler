package services

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"companiesHandler/models"
	"companiesHandler/proto/pb"
	"companiesHandler/repos"
)

type UserServicer interface {
	Create(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CompanyResponse, error)
	Get(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyResponse, error)
}

type UserService struct {
	Repo *repos.RepoUser
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		Repo: repos.NewUserRepo(db),
	}
}

func (us *UserService) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	err := us.Repo.Save(ctx, &models.User{
		Username:       req.GetName(),
		HashedPassword: req.GetPassword(),
		Role:           req.GetRole(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save user: %v", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Name: req.GetName(),
			Role: req.GetRole(),
		},
	}, nil
}

func (us *UserService) Get(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := us.Repo.Find(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Name: user.Username,
			Role: user.Role,
		},
	}, nil
}

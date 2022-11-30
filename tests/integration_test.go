package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"companiesHandler/mappers"
	proto "companiesHandler/proto/pb"
)

//nolint
type config struct {
	GRPCAddr   string `env:"GRPC_ADDR,default=localhost:8080"`
	DBHost     string `env:"DB_DRIVER,default=localhost"`
	DBDriver   string `env:"DB_DRIVER,default=postgres"`
	DBUser     string `env:"DB_USER,default=test"`
	DBPassword string `env:"DB_PASSWORD,default=password"`
	DBName     string `env:"DB_NAME,default=companies"`
	DBPort     string `env:"DB_PORT,default=5432"`
}

//nolint
type dbCompanies struct {
	ID                string
	Name              string
	Description       string
	AmountOfEmployees string
	Type              string
	Registered        string
}

//nolint
type FunctionalTestSuite struct {
	suite.Suite
	config config

	dbManager *dbCompanies
}

//nolint
func TestFunctionalCases(t *testing.T) {
	suite.Run(t, &FunctionalTestSuite{})
}

//nolint
func (s *FunctionalTestSuite) SetupSuite() {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	comp := proto.NewCompanyServiceClient(conn)
	usr := proto.NewUserServiceClient(conn)
	auth := proto.NewAuthServiceClient(conn)
	usrReq := &proto.CreateUserRequest{
		Name:     "Alex",
		Password: "12345678",
		Role:     "admin",
	}
	usrResp, err := usr.Create(context.Background(), usrReq)
	fmt.Println(usrResp)

	compReq := &proto.CreateCompanyRequest{
		Name:             "Alex",
		Description:      nil,
		AmountOfEmployee: 3,
		Registered:       true,
		Type:             mappers.StringToTypes[mappers.Cooperative],
	}
	compResp, err := comp.Create(context.Background(), compReq)
	fmt.Println(compResp)

	authReq := &proto.LoginRequest{
		Username: "Alex",
		Password: "12345678",
	}
	authResp, err := auth.Login(context.Background(), authReq)
	fmt.Println(authResp)
}

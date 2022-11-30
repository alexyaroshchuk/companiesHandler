package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"companiesHandler/mappers"
	"companiesHandler/models"
	"companiesHandler/producers"
	"companiesHandler/proto/pb"
	"companiesHandler/repos"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
	"github.com/satori/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompanyServicer interface {
	Create(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CompanyResponse, error)
	Patch(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.CompanyResponse, error)
	Get(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyResponse, error)
	Delete(ctx context.Context, req *pb.CompanyRequest) (*pb.DeleteCompanyResponse, error)
}

type CompanyService struct {
	Repo     *repos.RepoCompany
	Producer *kafka.Producer
}

func NewCompanyService(db *sql.DB, p *kafka.Producer) *CompanyService {
	return &CompanyService{
		Repo:     repos.NewCompanyRepo(db),
		Producer: p,
	}
}

func (cs *CompanyService) Create(ctx context.Context, req *pb.CreateCompanyRequest) (*pb.CompanyResponse, error) {
	id := uuid.NewV4().String()
	company := models.Company{
		UUID:              id,
		Name:              req.GetName(),
		Description:       req.GetDescription(),
		AmountOfEmployees: req.GetAmountOfEmployee(),
		Registered:        req.GetRegistered(),
		Type:              mappers.TypesToString[req.GetType()],
	}
	err := cs.Repo.CreateCompany(ctx, company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create company: %v", err)
	}

	err = cs.sendDataToKafka(company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot send data to kafka: %v", err)
	}

	return &pb.CompanyResponse{
		Company: &pb.Company{
			UUID:              id,
			Name:              req.GetName(),
			Description:       req.Description,
			AmountOfEmployees: req.GetAmountOfEmployee(),
			Registered:        req.GetRegistered(),
			Type:              req.GetType(),
		},
	}, nil
}

func (cs *CompanyService) Get(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyResponse, error) {
	company, err := cs.Repo.GetCompany(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get company: %v", err)
	}

	return &pb.CompanyResponse{
		Company: &pb.Company{
			UUID:              company.UUID,
			Name:              company.Name,
			Description:       &company.Description,
			AmountOfEmployees: company.AmountOfEmployees,
			Registered:        company.Registered,
			Type:              mappers.StringToTypes[company.Type],
		},
	}, nil
}

func (cs *CompanyService) Patch(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.CompanyResponse, error) {
	var cType *string
	if req.Type != nil {
		sType := mappers.TypesToString[req.GetType()]
		cType = &sType
	}
	err := cs.Repo.UpdateCompany(ctx, models.CompanyForUpdate{
		UUID:              req.GetUUID(),
		Name:              req.Name,
		Description:       req.Description,
		AmountOfEmployees: req.AmountOfEmployee,
		Registered:        req.Registered,
		Type:              cType,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot update company: %v", err)
	}

	company, err := cs.Repo.GetCompany(ctx, req.GetUUID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot fetch company: %v", err)
	}

	err = cs.sendDataToKafka(company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot send data to kafka: %v", err)
	}

	return &pb.CompanyResponse{
		Company: &pb.Company{
			UUID:              company.UUID,
			Name:              company.Name,
			Description:       &company.Description,
			AmountOfEmployees: company.AmountOfEmployees,
			Registered:        company.Registered,
			Type:              mappers.StringToTypes[company.Type],
		},
	}, nil
}

func (cs *CompanyService) Delete(ctx context.Context, req *pb.CompanyRequest) (*pb.DeleteCompanyResponse, error) {
	err := cs.Repo.DeleteCompany(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete company: %v", err)
	}

	return &pb.DeleteCompanyResponse{
		Success: true,
	}, nil
}

func (cs *CompanyService) sendDataToKafka(company models.Company) error {
	stringJson, _ := json.Marshal(company)
	fmt.Println(string(stringJson))
	err := producers.Produce(cs.Producer, string(stringJson), os.Getenv("KAFKA_TOPIC"))
	if err != nil {
		return err
	}

	return nil
}

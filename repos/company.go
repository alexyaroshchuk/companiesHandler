package repos

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"companiesHandler/models"
)

type CompanyRepo interface {
	GetCompany(ctx context.Context, UUID string) (models.Company, error)
	UpdateCompany(ctx context.Context, c models.CompanyForUpdate) error
	CreateCompany(ctx context.Context, c models.Company) error
	DeleteCompany(ctx context.Context, UUID string) error
}

type RepoCompany struct {
	Db *sql.DB
}

func NewCompanyRepo(postCollection *sql.DB) *RepoCompany {
	return &RepoCompany{
		Db: postCollection,
	}
}

func (rc *RepoCompany) GetCompany(ctx context.Context, UUID string) (models.Company, error) {
	sqlStatement := `select * from "companies" where id=$1`
	var company models.Company
	err := rc.Db.QueryRow(sqlStatement, UUID).Scan(
		&company.UUID,
		&company.Name,
		&company.Description,
		&company.AmountOfEmployees,
		&company.Registered,
		&company.Type,
	)
	if err != nil {
		log.Printf("Company couldn't be fetched, %v\n", err)
		return company, err
	}

	return company, nil
}

func (rc *RepoCompany) UpdateCompany(ctx context.Context, c models.CompanyForUpdate) error {
	sqlStatement := `UPDATE "companies" SET `
	qParts := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)
	i := 1
	if c.Name != nil {
		qParts = append(qParts, fmt.Sprintf(`name = $%d`, i))
		i += 1
		args = append(args, c.Name)
	}
	if c.Description != nil {
		qParts = append(qParts, fmt.Sprintf(`description = $%d`, i))
		i += 1
		args = append(args, c.Description)
	}
	if c.AmountOfEmployees != nil {
		qParts = append(qParts, fmt.Sprintf(`amount_of_employees = $%d`, i))
		i += 1
		args = append(args, c.AmountOfEmployees)
	}
	if c.Registered != nil {
		qParts = append(qParts, fmt.Sprintf(`reqistered = $%d`, i))
		i += 1
		args = append(args, c.Registered)
	}
	if c.Type != nil {
		qParts = append(qParts, fmt.Sprintf(`type = $%d`, i))
		i += 1
		args = append(args, c.Type)
	}

	sqlStatement += strings.Join(qParts, ",") + fmt.Sprintf(` WHERE id =  $%d`, i)
	args = append(args, c.UUID)

	_, err := rc.Db.Exec(sqlStatement, args...)
	if err != nil {
		log.Printf("Company couldn't be updated, %v\n", err)
		return err
	}

	return nil
}

func (rc *RepoCompany) CreateCompany(ctx context.Context, c models.Company) error {
	sqlStatement := `INSERT INTO "companies" (
                     id, name, description, amount_of_employees, registered, type) 
					 VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := rc.Db.Exec(sqlStatement, c.UUID, c.Name, c.Description, c.AmountOfEmployees, c.Registered, c.Type); err != nil {
		log.Printf("Company couldn't be inserted, %v\n", err)
		return err
	}

	return nil
}

func (rc *RepoCompany) DeleteCompany(ctx context.Context, UUID string) error {
	sqlStatement := `delete from "companies" where id=$1`
	if _, err := rc.Db.Exec(sqlStatement, UUID); err != nil {
		if err != nil {
			log.Printf("Company couldn't be deleted, %v\n", err)
			return err
		}
	}

	return nil
}

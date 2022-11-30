package repos

import (
	"context"
	"database/sql"
	"log"

	"companiesHandler/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserStore interface {
	Find(ctx context.Context, username, password string) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
}

func NewUserRepo(postCollection *sql.DB) *RepoUser {
	return &RepoUser{
		Db: postCollection,
	}
}

type RepoUser struct {
	Db *sql.DB
}

func (ur *RepoUser) Find(ctx context.Context, username, password string) (*models.User, error) {
	sqlStatement := `select * from "users" where username=$1`
	var user models.User
	err := ur.Db.QueryRow(sqlStatement, username).Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Role,
	)

	if !user.IsCorrectPassword(password) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}

	if err != nil {
		log.Printf("User couldn't be fetched, %v\n", err)
		return &user, err
	}

	return &user, nil
}

func (ur *RepoUser) Save(ctx context.Context, user *models.User) error {
	user, err := models.NewUser(user.Username, user.HashedPassword, user.Role)
	if err != nil {
		log.Printf("User couldn't be created, %v\n", err)
	}
	sqlStatement := `INSERT INTO "users" (username, password, role) 
					 VALUES ($1, $2, $3)`

	if _, err := ur.Db.Exec(sqlStatement, user.Username, user.HashedPassword, user.Role); err != nil {
		log.Printf("User couldn't be inserted, %v\n", err)
		return err
	}

	return nil
}

package repo

import (
	"context"
	"errors"

	"back/lk/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserInterface interface {
	GetUser(ctx context.Context, email string, id uint64) (*entity.User, error)
	CreateRest(ctx context.Context, name, email, password string) (entity.User, error)
	GetPassword(ctx context.Context, id uint64) (entity.Password, error)
}

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) UserInterface {
	return &User{db: db}
}

func (u *User) GetUser(ctx context.Context, email string, id uint64) (*entity.User, error) {
	if email == "" && id == 0 {
		return nil, errors.New("")
	}
	var res entity.User
	query := `select id, name, email from restaurant where email=$1`
	err := u.db.QueryRow(ctx, query, email).Scan(&res.ID, &res.Name, &res.Email)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (u *User) CreateRest(ctx context.Context, name, email, password string) (entity.User, error) {
	query := `insert into restaurant(name, email, password) values ($1, $2, $3) returning id, name, email;`
	var res entity.User
	err := u.db.QueryRow(ctx, query, name, email, password).Scan(&res.ID, &res.Name, &res.Email)
	if err != nil {
		return entity.User{}, err
	}
	return res, nil
}

func (u *User) GetPassword(ctx context.Context, id uint64) (entity.Password, error) {
	query := `select password from restaurant where id=$1;`
	var res string
	err := u.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return "", err
	}
	return entity.Password(res), nil
}

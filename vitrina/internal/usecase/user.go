package usecase

import (
	"back/vitrina/internal/entity"
	"back/vitrina/internal/repo"

	"context"
)

type UserInterface interface {
	Signup(ctx context.Context, data entity.AuthRequest) (entity.User, error)
	SignIn(ctx context.Context, data entity.AuthRequest) (entity.User, error)
}

type User struct {
	repo repo.UserInterface
}

func NewUser(r repo.UserInterface) UserInterface {
	return &User{repo: r}
}

func (u *User) Signup(ctx context.Context, data entity.AuthRequest) (entity.User, error) {
	user, err := u.repo.GetUser(ctx, data.Email, 0)
	if err != nil {
		return entity.User{}, err
	}
	if user != nil {
		return entity.User{}, repo.ErrEqualUser
	}
	var pass entity.Password
	err = pass.Hash(data.Password)
	if err != nil {
		return entity.User{}, err
	}
	res, err := u.repo.CreateUser(ctx, data.Name, data.Email, string(pass))
	if err != nil {
		return entity.User{}, err
	}
	return res, nil
}

func (u *User) SignIn(ctx context.Context, data entity.AuthRequest) (entity.User, error) {
	user, err := u.repo.GetUser(ctx, data.Email, 0)
	if err != nil {
		return entity.User{}, err
	}
	if user == nil {
		return entity.User{}, repo.ErrWrongLoginOrPassword
	}
	password, err := u.repo.GetPassword(ctx, user.ID)
	if err != nil {
		return entity.User{}, err
	}
	if !password.IsEqual(data.Password) {
		return entity.User{}, repo.ErrWrongLoginOrPassword
	}
	return *user, nil
}
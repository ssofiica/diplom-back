package repo

import "errors"

var (
	ErrWrongLoginOrPassword = errors.New("Неверный логин или пароль")
	ErrEqualUser = errors.New("Такой пользователь уже существует")
)

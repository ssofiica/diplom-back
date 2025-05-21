package delivery

import "errors"

var (
	ErrDefault400 = errors.New("Неверный запрос")
	ErrDefault401 = errors.New("Неавторизован")
	ErrDefault500 = errors.New("Ошибка сервера")

	ErrTokenGenerate = errors.New("Ошибка генерации токена")
	ErrNoRequestVars = errors.New("Остуствуют параметры запроса")
	ErrNotValidBody  = errors.New("Неверные данные тела запроса")
	ErrWrongStatus = errors.New("Неверный статус")
)

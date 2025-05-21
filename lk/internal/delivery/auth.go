package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"back/lk/internal/entity"
	"back/lk/internal/usecase"
	"back/lk/internal/utils/request"
	"back/lk/internal/utils/response"
)

var userKey string = "user"

type AuthHandler struct {
	usecase usecase.UserInterface
	jwt     entity.JWT
}

func NewAuthHandler(u usecase.UserInterface, t entity.JWT) *AuthHandler {
	return &AuthHandler{usecase: u, jwt: t}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	handler := "SignUp"
	payload := entity.AuthRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, handler, ErrDefault400)
		return
	}
	if !payload.Valid() {
		response.WithError(w, 400, handler, ErrDefault400)
		return
	}
	userData, err := h.usecase.Signup(context.Background(), payload)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, usecase.ErrWrongLoginOrPassword) {
			response.WithError(w, 500, handler, usecase.ErrWrongLoginOrPassword)
			return
		}
		if errors.Is(err, usecase.ErrEqualUser) {
			response.WithError(w, 400, handler, usecase.ErrEqualUser)
			return
		}
		response.WithError(w, 500, handler, ErrDefault500)
		return
	}
	jwtToken, err := h.jwt.GenerateToken(userData.ID, userData.Name, userData.Email)
	if err != nil {
		response.WithError(w, 500, handler, ErrTokenGenerate)
		return
	}
	w.Header().Set("Authorization", "Bearer "+jwtToken)
	res := entity.AuthResponse{Token: jwtToken}
	response.WriteData(w, res, 200)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	handler := "SignIn"
	payload := entity.AuthRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, handler, ErrDefault400)
		return
	}
	payload.Name = "def"
	if !payload.Valid() {
		response.WithError(w, 400, handler, ErrDefault400)
		return
	}
	userData, err := h.usecase.SignIn(context.Background(), payload)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, usecase.ErrWrongLoginOrPassword) {
			response.WithError(w, 500, handler, usecase.ErrWrongLoginOrPassword)
			return
		}
		response.WithError(w, 500, handler, ErrDefault500)
		return
	}
	jwtToken, err := h.jwt.GenerateToken(userData.ID, userData.Name, userData.Email)
	if err != nil {
		response.WithError(w, 500, handler, ErrTokenGenerate)
		return
	}
	w.Header().Set("Authorization", "Bearer "+jwtToken)
	res := entity.AuthResponse{Token: jwtToken}
	response.WriteData(w, res, 200)
}

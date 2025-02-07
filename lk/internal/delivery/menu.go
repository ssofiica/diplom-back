package delivery

import (
	"back/lk/internal/entity"
	"back/lk/internal/usecase"
	"back/lk/internal/utils/request"
	"back/lk/internal/utils/response"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MenuHandler struct {
	usecase usecase.MenuInterface
}

func NewMenuHandler(u usecase.MenuInterface) *MenuHandler {
	return &MenuHandler{usecase: u}
}

func (h *MenuHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	payload := entity.FoodDTO{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "AddFood", err)
		return
	}
	res, err := h.usecase.AddFood(context.Background(), payload.ToFood())
	if err != nil {
		response.WithError(w, 500, "AddFood", err)
		return
	}
	response.WriteData(w, res, 200)
}

func (h *MenuHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	payload := entity.CategoryDTO{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "AddCategory", err)
		return
	}
	res, err := h.usecase.AddCategory(context.Background(), payload.ToCategory())
	if err != nil {
		response.WithError(w, 500, "AddCategory", err)
		return
	}
	response.WriteData(w, res, 200)
}

func (h *MenuHandler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	if value == "" {
		fmt.Println("no id")
		response.WithError(w, 400, "DeleteFood", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "DeleteFood", err)
		return
	}
	err = h.usecase.DeleteFood(context.Background(), uint64(id))
	if err != nil {
		response.WithError(w, 500, "DeleteFood", err)
		return
	}
	response.WriteData(w, "Блюдо удалено", 200)
}

func (h *MenuHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: чтото сделать с блюдами, которые находятся в этой категории
	vars := mux.Vars(r)
	value := vars["id"]
	if value == "" {
		fmt.Println("no id")
		response.WithError(w, 400, "DeleteCategory", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "DeleteCategory", err)
		return
	}
	err = h.usecase.DeleteCategory(context.Background(), uint64(id))
	if err != nil {
		response.WithError(w, 500, "DeleteCategory", err)
		return
	}
	response.WriteData(w, "Категория удалена", 200)
}

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

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	restId := uint64(1)
	res, err := h.usecase.GetMenu(context.Background(), restId)
	if err != nil {
		response.WithError(w, 500, "GetMenu", err)
		return
	}
	resDTO := make([]entity.CategoryDTO, len(res))
	for i, c := range res {
		resDTO[i] = c.ToDTO()
	}
	response.WriteData(w, resDTO, 200)
}

func (h *MenuHandler) GetCategoryList(w http.ResponseWriter, r *http.Request) {
	restId := uint64(1)
	res, err := h.usecase.GetCategoryList(context.Background(), restId)
	if err != nil {
		response.WithError(w, 500, "GetCategoryList", err)
		return
	}
	resDTO := make([]entity.CategoryDTO, len(res))
	for i, c := range res {
		resDTO[i] = c.ToDTO()
	}
	response.WriteData(w, resDTO, 200)
}

// еда по статусу для определенной категории
func (h *MenuHandler) GetFoodByStatus(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var val1 string = params.Get("status")
	var status entity.FoodStatus
	status.Scan(val1)

	vars := mux.Vars(r)
	val2 := vars["category_id"]
	if val2 == "" {
		fmt.Println("no id")
		response.WithError(w, 400, "GetFoodByStatus", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(val2)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "GetFoodByStatus", err)
		return
	}

	res, err := h.usecase.GetFoodByStatus(context.Background(), status, uint64(id))
	if err != nil {
		response.WithError(w, 500, "GetFoodByStatus", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
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
	response.WriteData(w, res.ToDTO(), 200)
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

func (h *MenuHandler) EditFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	if value == "" {
		fmt.Println("no id")
		response.WithError(w, 400, "EditFood", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "EditFood", err)
		return
	}
	payload := entity.EditFood{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "EditFood", err)
		return
	}
	if payload.Status != "" && !entity.IsFoodStatus(payload.Status) {
		response.WithError(w, 400, "EditFood", ErrWrongStatus)
		return
	}
	res, err := h.usecase.EditFood(context.Background(), uint32(id), payload)
	if err != nil {
		response.WithError(w, 500, "EditFood", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *MenuHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	value := vars["id"]
	if value == "" {
		fmt.Println("no id")
		response.WithError(w, 400, "ChangeStatus", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "ChangeStatus", err)
		return
	}
	payload := entity.ChangeStatusRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "ChangeStatus", err)
		return
	}
	if !entity.IsFoodStatus(payload.Status) {
		response.WithError(w, 400, "ChangeStatus", ErrWrongStatus)
		return
	}
	err = h.usecase.ChangeStatus(context.Background(), uint32(id), payload.Status)
	if err != nil {
		response.WithError(w, 500, "ChangeStatus", err)
		return
	}
	response.WriteData(w, nil, 200)
}

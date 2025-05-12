package delivery

import (
	"back/vitrina/internal/entity"
	"back/vitrina/internal/usecase"
	"back/vitrina/utils/request"
	"back/vitrina/utils/response"
	"errors"
	"fmt"
	"strconv"

	"context"
	"net/http"

	"github.com/gorilla/mux"
)

var userKey string = "user"

type OrderHandler struct {
	usecase usecase.OrderInterface
}

func NewOrderHandler(u usecase.OrderInterface) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func (h *OrderHandler) ChangeFoodCountInBasket(w http.ResponseWriter, r *http.Request) {
	restId := uint32(1)
	//userId := uint32(1)
	user, ok := r.Context().Value(userKey).(entity.User)
	if !ok {
		response.WithError(w, 401, "GetCurrent", ErrDefault401)
		return
	}
	payload := entity.RequestAddFood{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "AddFood", err)
		return
	}
	if !payload.Valid() {
		response.WithError(w, 400, "AddFood", ErrNotValidBody)
		return
	}
	err := h.usecase.AddFoodToOrder(context.Background(), user.ID, restId, payload)
	if err != nil {
		if errors.Is(usecase.ErrFoodStoped, err) {
			response.WithError(w, 409, "AddFood", err)
			return
		}
		response.WithError(w, 500, "AddFood", err)
		return
	}
	basket, err := h.usecase.GetBasket(context.Background(), user.ID, 0)
	if err != nil {
		if errors.Is(usecase.ErrFoodStoped, err) {
			response.WithError(w, 409, "AddFood", err)
			return
		}
		response.WithError(w, 500, "AddFood", err)
		return
	}
	response.WriteData(w, basket.ToDTO(), 200)
}

func (h *OrderHandler) GetUserBasket(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	userId := uint32(1)
	res, err := h.usecase.GetBasket(context.Background(), userId, 0)
	if err != nil {
		if errors.Is(usecase.ErrFoodStoped, err) {
			response.WithError(w, 409, "GetBasket", err)
			return
		}
		response.WithError(w, 500, "GetBasket", err)
		return
	}
	if res.Id == 0 {
		response.WriteData(w, "У вас нет корзины", 200)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	//userId := uint32(1)
	vars := mux.Vars(r)
	value := vars["id"]
	if value == "" {
		response.WithError(w, 400, "GetOrderById", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "GetOrderById", err)
		return
	}
	res, err := h.usecase.GetOrderById(context.Background(), uint32(id))
	if err != nil {
		response.WithError(w, 500, "GetOrderById", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) UpdateBasketInfo(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	userId := uint32(1)
	payload := entity.RequestBasketInfo{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "UpdateBasketInfo", err)
		return
	}
	if !payload.Valid() {
		response.WithError(w, 400, "UpdateBasketInfo", ErrNotValidBody)
		return
	}
	res, err := h.usecase.UpdateBasketInfo(context.Background(), userId, payload)
	if err != nil {
		response.WithError(w, 500, "UpdateBasketInfo", err)
		return
	}
	if res.Id == 0 {
		response.WriteData(w, "У вас нет корзины", 200)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) Pay(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	userId := uint32(1)
	id, err := h.usecase.Pay(context.Background(), userId)
	if err != nil {
		if errors.Is(err, usecase.ErrNeedAddress) {
			response.WriteData(w, err.Error(), 200)
			return
		}
		response.WithError(w, 500, "Pay", err)
		return
	}
	res, err := h.usecase.GetOrderById(context.Background(), id)
	if err != nil {
		response.WithError(w, 500, "Pay", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	user, ok := r.Context().Value(userKey).(entity.User)
	if !ok {
		response.WithError(w, 401, "GetCurrent", ErrDefault401)
		return
	}
	// userId := uint32(1)
	res, err := h.usecase.Current(context.Background(), user.ID)
	if err != nil {
		response.WithError(w, 500, "GetCurrent", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) GetArchive(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	userId := uint32(1)
	res, err := h.usecase.Archive(context.Background(), userId)
	if err != nil {
		response.WithError(w, 500, "GetCurrent", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

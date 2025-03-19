package delivery

import (
	"back/vitrina/internal/entity"
	"back/vitrina/internal/usecase"
	"back/vitrina/utils/request"
	"back/vitrina/utils/response"
	"errors"

	"context"
	"net/http"
)

type OrderHandler struct {
	usecase usecase.OrderInterface
}

func NewOrderHandler(u usecase.OrderInterface) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func (h *OrderHandler) AddFoodToOrder(w http.ResponseWriter, r *http.Request) {
	restId := uint32(1)
	userId := uint32(1)
	payload := entity.RequestAddFood{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "AddFood", err)
		return
	}
	if !payload.Valid() {
		response.WithError(w, 400, "AddFood", ErrNotValidBody)
		return
	}
	err := h.usecase.AddFoodToOrder(context.Background(), userId, restId, payload)
	if err != nil {
		if errors.Is(usecase.ErrFoodStoped, err) {
			response.WithError(w, 409, "AddFood", err)
			return
		}
		response.WithError(w, 500, "AddFood", err)
		return
	}
	response.WriteData(w, "ok", 200)
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

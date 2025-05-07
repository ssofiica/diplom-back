package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"back/lk/internal/usecase"
	"back/lk/internal/utils/response"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	usecase usecase.OrderInterface
}

func NewOrder(u usecase.OrderInterface) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func (h *OrderHandler) GetMiniOrders(w http.ResponseWriter, r *http.Request) {
	restId := uint32(1)
	params := r.URL.Query()
	value := params.Get("status")
	if value == "" {
		response.WithError(w, 400, "GetOrdersNew", errors.New("missing request var"))
		return
	}
	res, err := h.usecase.GetMiniOrders(context.Background(), restId, value)
	if err != nil {
		response.WithError(w, 500, "GetOrdersNew", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
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

func (h *OrderHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	//restId := uint32(1)
	params := r.URL.Query()
	status := params.Get("status")
	if status == "" {
		response.WithError(w, 400, "ChangeStatus", errors.New("missing request var"))
		return
	}
	vars := mux.Vars(r)
	value1 := vars["id"]
	if value1 == "" {
		response.WithError(w, 400, "ChangeStatus", errors.New("missing request var"))
		return
	}
	id, err := strconv.Atoi(value1)
	if err != nil {
		fmt.Println("err in converting str to int")
		response.WithError(w, 400, "ChangeStatus", err)
		return
	}
	err = h.usecase.UpdateStatus(context.Background(), uint32(id), status)
	if err != nil {
		response.WithError(w, 500, "ChangeStatus", err)
		return
	}
	response.WriteData(w, "ok", 200)
}

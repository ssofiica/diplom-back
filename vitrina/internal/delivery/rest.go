package delivery

import (
	"back/vitrina/entity"
	"back/vitrina/internal/usecase"
	"back/vitrina/utils/response"

	"context"
	"net/http"
)

type RestHandler struct {
	usecase usecase.RestInterface
}

func NewRestHandler(u usecase.RestInterface) *RestHandler {
	return &RestHandler{usecase: u}
}

func (h *RestHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	restId := uint64(1)
	res, err := h.usecase.GetInfo(context.Background(), restId)
	if err != nil {
		response.WithError(w, 500, "GetInfo", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *RestHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
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

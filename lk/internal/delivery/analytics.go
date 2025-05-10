package delivery

import (
	"context"
	"net/http"

	"back/lk/internal/entity"
	"back/lk/internal/usecase"
	"back/lk/internal/utils/request"
	"back/lk/internal/utils/response"
)

type Analytics struct {
	usecase usecase.AnalyticsInterface
}

func NewAnalytics(u usecase.AnalyticsInterface) *Analytics {
	return &Analytics{usecase: u}
}

func (h *Analytics) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	payload := entity.DateIntervalRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "GetAnalytics", err)
		return
	}
	res, err := h.usecase.GetAnalytics(context.Background(), restId, payload.Start, payload.End)
	if err != nil {
		response.WithError(w, 500, "GetAnalytics", err)
		return
	}
	response.WriteData(w, res, 200)
}

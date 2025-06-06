package delivery

import (
	"context"
	"errors"
	"net/http"

	"back/lk/internal/entity"
	"back/lk/internal/usecase"
	"back/lk/internal/utils/response"
)

type Analytics struct {
	usecase usecase.AnalyticsInterface
}

func NewAnalytics(u usecase.AnalyticsInterface) *Analytics {
	return &Analytics{usecase: u}
}

func (h *Analytics) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	rest, ok := r.Context().Value(userKey).(entity.User)
	if !ok {
		response.WithError(w, 401, "GetCategoryList", ErrDefault401)
		return
	}
	params := r.URL.Query()
	start := params.Get("start")
	end := params.Get("end")
	if start == "" || end == "" {
		response.WithError(w, 400, "GetAnalytics", errors.New("missing request var"))
		return
	}
	res, err := h.usecase.GetAnalytics(context.Background(), rest.ID, start, end)
	if err != nil {
		response.WithError(w, 500, "GetAnalytics", err)
		return
	}
	response.WriteData(w, res, 200)
}

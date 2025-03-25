package delivery

import (
	"back/lk/internal/usecase"
	"back/lk/internal/utils/response"
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
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

func (h *RestHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		response.WithError(w, 400, "UploadImage", err)
		return
	}

	id := r.FormValue("restaurant_id")
	restId, err := strconv.Atoi(id)
	if err != nil {
		response.WithError(w, 400, "UploadImage", err)
		return
	}

	file, fileHeader, err := r.FormFile("logo_url")
	if err != nil {
		response.WithError(w, 400, "UploadImage", err)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.WithError(w, 500, "UploadImage", err)
		return
	}
	fileExtension := filepath.Ext(fileHeader.Filename)
	mimeType := GetMimeType(fileExtension)
	if mimeType == "" {
		response.WithError(w, 400, "UploadImage", err)
		return
	}
	err = h.usecase.UploadLogo(context.Background(), fileBytes, fileExtension, mimeType, uint64(restId))
	if err != nil {
		response.WithError(w, 500, "UploadImage", err)
		return
	}
	response.WriteData(w, nil, 200)
}

func GetMimeType(ext string) string {
	switch ext {
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}

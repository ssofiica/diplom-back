package delivery

import (
	"back/lk/internal/entity"
	"back/lk/internal/usecase"
	"back/lk/internal/utils/request"
	"back/lk/internal/utils/response"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path/filepath"
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

func (h *RestHandler) UploadBaseInfo(w http.ResponseWriter, r *http.Request) {
	restId := uint64(1)
	payload := entity.BaseInfoRequest{}
	if err := request.GetRequestData(r, &payload); err != nil {
		response.WithError(w, 400, "UploadBaseInfo", err)
		return
	}
	res, err := h.usecase.UploadBaseInfo(context.Background(), payload.FromDTO(), restId)
	if err != nil {
		response.WithError(w, 500, "UploadBaseInfo", err)
		return
	}
	response.WriteData(w, res.ToDTO(), 200)
}

func (h *RestHandler) UploadDescriptionsAndImages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		response.WithError(w, 400, "UploadDescriptionsAndImages", err)
		return
	}

	content := entity.DescripAndImgs{}
	if descriptions := r.FormValue("descriptions"); descriptions != "" {
		if err = json.Unmarshal([]byte(descriptions), &content.Description); err != nil {
			response.WithError(w, 400, "UploadDescriptionsAndImages", err)
			return
		}
		if dIndexes := r.FormValue("descrip_indexes"); dIndexes != "" {
			if err = json.Unmarshal([]byte(dIndexes), &content.DescripIndexes); err != nil {
				response.WithError(w, 400, "UploadDescriptionsAndImages", err)
				return
			}
		}
		if len(content.DescripIndexes) != len(content.Description) {
			response.WithError(w, 400, "UploadDescriptionsAndImages", errors.New("Нужны индексы описаний"))
			return
		}
	}

	files := r.MultipartForm.File["images"]
	var imgIndexArray []uint8
	if files != nil {
		if imgIndexes := r.FormValue("img_indexes"); imgIndexes != "" {
			if err = json.Unmarshal([]byte(imgIndexes), &imgIndexArray); err != nil {
				response.WithError(w, 400, "UploadDescriptionsAndImages", err)
				return
			}
		}
		if len(files) != len(imgIndexArray) {
			response.WithError(w, 400, "UploadDescriptionsAndImages", errors.New("Нужны индексы картинок"))
			return
		}
	}
	var img entity.Img
	for i, fileHeader := range files {
		// Открываем файл
		file, err := fileHeader.Open()
		if err != nil {
			response.WithError(w, 500, "UploadDescriptionsAndImages", err)
			return
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			response.WithError(w, 500, "UploadDescriptionsAndImages", err)
			return
		}
		img.Data = fileBytes
		img.Ext = filepath.Ext(fileHeader.Filename)
		mime := GetMimeType(img.Ext)
		if mime == "" {
			response.WithError(w, 400, "UploadImage", errors.New("cant take mime type"))
			return
		}
		img.Index = imgIndexArray[i]
		img.MimeType = mime
		content.Img = append(content.Img, img)
	}
	err = h.usecase.UploadDescriptionAndImages(context.Background(), &content, restId)
	if err != nil {
		response.WithError(w, 500, "UploadDescriptionsAndImages", err)
		return
	}
	response.WriteData(w, nil, 200)
}

func (h *RestHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
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
	case ".jpg", ".JPG":
		return "image/jpeg"
	case ".jpeg", ".JPEG":
		return "image/jpeg"
	case ".png", ".PNG":
		return "image/png"
	default:
		return ""
	}
}

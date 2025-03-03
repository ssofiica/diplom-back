package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WithError(w http.ResponseWriter, statusCode int, handlerName string, err error) {
	fmt.Println(err.Error(), "handler", handlerName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(`{"error":"Error occurred"}`))
}

func WriteData(w http.ResponseWriter, data interface{}, statusCode int) {
	if data == nil {
		data = "Успешный ответ"
	}
	body, err := json.Marshal(data)
	if err != nil {
		WithError(w, 500, "", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

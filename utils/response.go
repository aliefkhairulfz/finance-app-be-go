package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Send(w http.ResponseWriter, status string, code int, message string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ApiResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	})

	return nil
}

package utils

import (
	"encoding/json"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
)

type ServerResponse struct {
	Success  bool        `json:"success"`
	Code     string      `json:"code,omitempty"`
	Messages interface{} `json:"messages,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

func WriteErrorJsonResponse(w http.ResponseWriter, errorCode string, statusCode int, messages interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &ServerResponse{
		Success:  false,
		Code:     errorCode,
		Messages: messages,
	}

	json.NewEncoder(w).Encode(res)
}

func WriteSuccessJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := &ServerResponse{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(res)
}

func WriteInternalErrorJsonResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	res := &ServerResponse{
		Success:  false,
		Code:     constants.ErrorCodes.InternalServerError,
		Messages: []string{"Something went wrong"},
	}

	json.NewEncoder(w).Encode(res)
}

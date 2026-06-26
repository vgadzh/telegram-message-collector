package httpx

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, status int, message string) {
	_ = JSON(w, status, ErrorResponse{Error: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "unauthorized")
}

func InternalServerError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}

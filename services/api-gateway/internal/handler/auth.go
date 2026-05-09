package handler

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type AuthHandler struct {
	conn *grpc.ClientConn
}

func NewAuthHandler(conn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{conn: conn}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}
	
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "auth-service not yet connected"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "auth-service not yet connected"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "auth-service not yet connected"})
}

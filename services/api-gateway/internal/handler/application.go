package handler

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ApplicationHandler struct {
	conn *grpc.ClientConn
}

func NewApplicationHandler(conn *grpc.ClientConn) *ApplicationHandler {
	return &ApplicationHandler{conn: conn}
}

func (h *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "application-service not yet connected"})
}
func (h *ApplicationHandler) GetMyApplications(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "application-service not yet connected"})
}
func (h *ApplicationHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusNoContent, nil)
}
func (h *ApplicationHandler) GetListingApplications(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "application-service not yet connected"})
}
func (h *ApplicationHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "application-service not yet connected"})
}
func (h *ApplicationHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "application-service not yet connected"})
}

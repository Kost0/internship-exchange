package handler

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ListingHandler struct {
	conn *grpc.ClientConn
}

func NewListingHandler(conn *grpc.ClientConn) *ListingHandler {
	return &ListingHandler{conn: conn}
}

func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) GetMyListings(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusNoContent, nil)
}
func (h *ListingHandler) PublishListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}
func (h *ListingHandler) CloseListing(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "listing-service not yet connected"})
}

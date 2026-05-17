package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	listingpb "github.com/Kost0/internship-exchange/proto/listing"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/middleware"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ListingHandler struct {
	client listingpb.ListingServiceClient
}

func NewListingHandler(conn *grpc.ClientConn) *ListingHandler {
	return &ListingHandler{client: listingpb.NewListingServiceClient(conn)}
}

func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	res, err := h.client.GetListings(r.Context(), &listingpb.GetListingsRequest{
		Query:          q.Get("query"),
		Format:         q.Get("format"),
		EmploymentType: q.Get("employmentType"),
		City:           q.Get("city"),
		Skill:          q.Get("skill"),
		Page:           int32QueryParam(q.Get("page"), 1),
		Limit:          int32QueryParam(q.Get("limit"), 12),
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ListingHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.client.GetListing(r.Context(), &listingpb.GetListingRequest{Id: id})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ListingHandler) GetMyListings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.client.GetMyListings(r.Context(), &listingpb.GetMyListingsRequest{
		CompanyId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		Requirements   string `json:"requirements"`
		WhatWeOffer    string `json:"whatWeOffer"`
		City           string `json:"city"`
		Format         string `json:"format"`
		EmploymentType string `json:"employmentType"`
		SalaryFrom     int64  `json:"salaryFrom"`
		SalaryTo       int64  `json:"salaryTo"`
		SalaryCurrency string `json:"salaryCurrency"`
		Deadline       string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.CreateListing(r.Context(), &listingpb.CreateListingRequest{
		CompanyId:      userID,
		Title:          body.Title,
		Description:    body.Description,
		Requirements:   body.Requirements,
		WhatWeOffer:    body.WhatWeOffer,
		City:           body.City,
		Format:         body.Format,
		EmploymentType: body.EmploymentType,
		SalaryFrom:     body.SalaryFrom,
		SalaryTo:       body.SalaryTo,
		SalaryCurrency: body.SalaryCurrency,
		Deadline:       body.Deadline,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, res)
}

func (h *ListingHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	var body struct {
		Title          string `json:"title"`
		Description    string `json:"description"`
		Requirements   string `json:"requirements"`
		WhatWeOffer    string `json:"whatWeOffer"`
		City           string `json:"city"`
		Format         string `json:"format"`
		EmploymentType string `json:"employmentType"`
		SalaryFrom     int64  `json:"salaryFrom"`
		SalaryTo       int64  `json:"salaryTo"`
		SalaryCurrency string `json:"salaryCurrency"`
		Deadline       string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateListing(r.Context(), &listingpb.UpdateListingRequest{
		Id:             id,
		CompanyId:      userID,
		Title:          body.Title,
		Description:    body.Description,
		Requirements:   body.Requirements,
		WhatWeOffer:    body.WhatWeOffer,
		City:           body.City,
		Format:         body.Format,
		EmploymentType: body.EmploymentType,
		SalaryFrom:     body.SalaryFrom,
		SalaryTo:       body.SalaryTo,
		SalaryCurrency: body.SalaryCurrency,
		Deadline:       body.Deadline,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	_, err := h.client.DeleteListing(r.Context(), &listingpb.DeleteListingRequest{
		Id:        id,
		CompanyId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ListingHandler) PublishListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	res, err := h.client.PublishListing(r.Context(), &listingpb.PublishListingRequest{
		Id:        id,
		CompanyId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)
		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ListingHandler) CloseListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	res, err := h.client.CloseListing(r.Context(), &listingpb.CloseListingRequest{
		Id:        id,
		CompanyId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)
		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func int32QueryParam(s string, fallback int32) int32 {
	if s == "" {
		return fallback
	}

	var n int32

	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return fallback
	}

	return n
}
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	profilepb "github.com/Kost0/internship-exchange/proto/profile"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/dto"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	listingpb "github.com/Kost0/internship-exchange/proto/listing"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/middleware"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ListingHandler struct {
	client  listingpb.ListingServiceClient
	profile profilepb.ProfileServiceClient
}

func NewListingHandler(listingConn, profileConn *grpc.ClientConn) *ListingHandler {
	return &ListingHandler{
		client:  listingpb.NewListingServiceClient(listingConn),
		profile: profilepb.NewProfileServiceClient(profileConn),
	}
}

func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	res, err := h.client.GetListings(r.Context(), &listingpb.GetListingsRequest{
		Query: q.Get("query"), Format: q.Get("format"),
		EmploymentType: q.Get("employmentType"), City: q.Get("city"),
		Skill: q.Get("skill"),
		Page:  int32QueryParam(q.Get("page"), 1),
		Limit: int32QueryParam(q.Get("limit"), 12),
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	result := dto.GetListingsResponse{
		Total: res.Total, Page: res.Page, Limit: res.Limit,
		Items: []dto.ListingResponse{},
	}
	for _, l := range res.Items {
		result.Items = append(result.Items, protoToListingDTO(l))

	}

	proxy.WriteJSON(w, http.StatusOK, result)
}

func (h *ListingHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.client.GetListing(r.Context(), &listingpb.GetListingRequest{Id: id})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	l := protoToListingDTO(res)

	proxy.WriteJSON(w, http.StatusOK, l)
}

func (h *ListingHandler) GetMyListings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.client.GetMyListings(r.Context(), &listingpb.GetMyListingsRequest{CompanyId: userID})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	items := []dto.ListingResponse{}

	for _, l := range res.Items {
		items = append(items, protoToListingDTO(l))
	}

	proxy.WriteJSON(w, http.StatusOK, items)
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
		CompanyId: userID, Title: body.Title, Description: body.Description,
		Requirements: body.Requirements, WhatWeOffer: body.WhatWeOffer,
		City: body.City, Format: body.Format, EmploymentType: body.EmploymentType,
		SalaryFrom: body.SalaryFrom, SalaryTo: body.SalaryTo,
		SalaryCurrency: body.SalaryCurrency, Deadline: body.Deadline,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	l := protoToListingDTO(res)

	proxy.WriteJSON(w, http.StatusCreated, l)
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
		Id: id, CompanyId: userID, Title: body.Title, Description: body.Description,
		Requirements: body.Requirements, WhatWeOffer: body.WhatWeOffer,
		City: body.City, Format: body.Format, EmploymentType: body.EmploymentType,
		SalaryFrom: body.SalaryFrom, SalaryTo: body.SalaryTo,
		SalaryCurrency: body.SalaryCurrency, Deadline: body.Deadline,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	l := protoToListingDTO(res)

	proxy.WriteJSON(w, http.StatusOK, l)
}

func (h *ListingHandler) PublishListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	res, err := h.client.PublishListing(r.Context(), &listingpb.PublishListingRequest{Id: id, CompanyId: userID})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	l := protoToListingDTO(res)

	proxy.WriteJSON(w, http.StatusOK, l)
}

func (h *ListingHandler) CloseListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	res, err := h.client.CloseListing(r.Context(), &listingpb.CloseListingRequest{Id: id, CompanyId: userID})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	l := protoToListingDTO(res)

	proxy.WriteJSON(w, http.StatusOK, l)
}

func (h *ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	h.syncCompany(r.Context(), userID)

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

func (h *ListingHandler) syncCompany(ctx context.Context, userID string) {
	profile, err := h.profile.GetMyCompanyProfile(ctx, &profilepb.GetMyCompanyProfileRequest{
		UserId: userID,
	})
	if err != nil {
		return
	}

	_, _ = h.client.SyncCompany(ctx, &listingpb.SyncCompanyRequest{
		UserId:   userID,
		Name:     profile.Name,
		LogoUrl:  profile.LogoUrl,
		Industry: profile.Industry,
		City:     profile.City,
	})
}

func protoToListingDTO(l *listingpb.ListingResponse) dto.ListingResponse {
	res := dto.ListingResponse{
		ID: l.Id, CompanyID: l.CompanyId, Title: l.Title,
		Description: l.Description, Requirements: l.Requirements,
		WhatWeOffer: l.WhatWeOffer, City: l.City, Format: l.Format,
		EmploymentType: l.EmploymentType, SalaryFrom: l.SalaryFrom,
		SalaryTo: l.SalaryTo, SalaryCurrency: l.SalaryCurrency,
		Deadline: l.Deadline, Status: l.Status,
		CreatedAt: l.CreatedAt, UpdatedAt: l.UpdatedAt,
		Skills: []dto.ListingSkillResponse{},
	}

	if l.Company != nil {
		res.Company = &dto.CompanyInfoResponse{
			ID: l.Company.Id, Name: l.Company.Name,
			LogoURL: l.Company.LogoUrl, Industry: l.Company.Industry,
			City: l.Company.City,
		}
	}

	for _, s := range l.Skills {
		res.Skills = append(res.Skills, dto.ListingSkillResponse{
			ID: s.Id, ListingID: s.ListingId,
			Skill: s.Skill, IsRequired: s.IsRequired,
		})
	}

	return res
}

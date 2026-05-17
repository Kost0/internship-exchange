package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/dto"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	applicationpb "github.com/Kost0/internship-exchange/proto/application"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/middleware"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ApplicationHandler struct {
	client applicationpb.ApplicationServiceClient
}

func NewApplicationHandler(conn *grpc.ClientConn) *ApplicationHandler {
	return &ApplicationHandler{client: applicationpb.NewApplicationServiceClient(conn)}
}

func (h *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())

	var body struct {
		ListingID   string `json:"listingId"`
		CoverLetter string `json:"coverLetter"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.Apply(r.Context(), &applicationpb.ApplyRequest{
		StudentId: studentID, ListingId: body.ListingID, CoverLetter: body.CoverLetter,
	})

	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, protoToAppDTO(res))
}

func (h *ApplicationHandler) GetMyApplications(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	res, err := h.client.GetMyApplications(r.Context(), &applicationpb.GetMyApplicationsRequest{StudentId: studentID})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	items := []dto.ApplicationResponse{}

	for _, a := range res.Items {
		items = append(items, protoToAppDTO(a))
	}

	proxy.WriteJSON(w, http.StatusOK, items)
}

func (h *ApplicationHandler) GetListingApplications(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	listingID := chi.URLParam(r, "id")

	res, err := h.client.GetListingApplications(r.Context(), &applicationpb.GetListingApplicationsRequest{
		ListingId: listingID, CompanyId: companyID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	items := []dto.ApplicationResponse{}

	for _, a := range res.Items {
		items = append(items, protoToAppDTO(a))
	}

	proxy.WriteJSON(w, http.StatusOK, items)
}

func (h *ApplicationHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	companyID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")
	var body struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.ChangeStatus(r.Context(), &applicationpb.ChangeStatusRequest{
		Id: id, CompanyId: companyID, Status: body.Status, Comment: body.Comment,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, protoToAppDTO(res))
}

func (h *ApplicationHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	companyID := middleware.GetUserID(r.Context())
	res, err := h.client.GetHistory(r.Context(), &applicationpb.GetHistoryRequest{
		Id: id, CompanyId: companyID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	events := []dto.ApplicationEventResponse{}

	for _, e := range res.Events {
		events = append(events, dto.ApplicationEventResponse{
			ID: e.Id, ApplicationID: e.ApplicationId,
			OldStatus: e.OldStatus, NewStatus: e.NewStatus,
			Comment: e.Comment, ChangedAt: e.ChangedAt,
		})
	}

	proxy.WriteJSON(w, http.StatusOK, events)
}

func (h *ApplicationHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	_, err := h.client.Withdraw(r.Context(), &applicationpb.WithdrawRequest{
		Id:        id,
		StudentId: studentID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func protoToAppDTO(a *applicationpb.ApplicationResponse) dto.ApplicationResponse {
	res := dto.ApplicationResponse{
		ID: a.Id, StudentID: a.StudentId, ListingID: a.ListingId,
		CoverLetter: a.CoverLetter, Status: a.Status,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
		Events: []dto.ApplicationEventResponse{},
	}
	for _, e := range a.Events {
		res.Events = append(res.Events, dto.ApplicationEventResponse{
			ID: e.Id, ApplicationID: e.ApplicationId,
			OldStatus: e.OldStatus, NewStatus: e.NewStatus,
			Comment: e.Comment, ChangedAt: e.ChangedAt,
		})
	}
	return res
}

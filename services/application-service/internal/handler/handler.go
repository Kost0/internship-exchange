package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	applicationpb "github.com/Kost0/internship-exchange/proto/application"
	"github.com/Kost0/internship-exchange/services/application-service/internal/model"
	"github.com/Kost0/internship-exchange/services/application-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/application-service/internal/service"
)

type ApplicationHandler struct {
	applicationpb.UnimplementedApplicationServiceServer
	svc *service.ApplicationService
}

func NewApplicationHandler(svc *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{svc: svc}
}

func (h *ApplicationHandler) Apply(ctx context.Context, req *applicationpb.ApplyRequest) (*applicationpb.ApplicationResponse, error) {
	if req.StudentId == "" || req.ListingId == "" {
		return nil, status.Error(codes.InvalidArgument, "student_id and listing_id are required")
	}

	app, err := h.svc.Apply(ctx, req.StudentId, req.ListingId, req.CoverLetter)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return nil, status.Error(codes.AlreadyExists, "already applied to this listing")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return appToProto(app), nil
}

func (h *ApplicationHandler) GetMyApplications(ctx context.Context, req *applicationpb.GetMyApplicationsRequest) (*applicationpb.GetApplicationsResponse, error) {
	apps, err := h.svc.GetMyApplications(ctx, req.StudentId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &applicationpb.GetApplicationsResponse{}
	for i := range apps {
		resp.Items = append(resp.Items, appToProto(&apps[i]))
	}

	return resp, nil
}

func (h *ApplicationHandler) Withdraw(ctx context.Context, req *applicationpb.WithdrawRequest) (*applicationpb.WithdrawResponse, error) {
	if err := h.svc.Withdraw(ctx, req.Id, req.StudentId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "application not found or cannot be withdrawn")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &applicationpb.WithdrawResponse{Success: true}, nil
}

func (h *ApplicationHandler) GetListingApplications(ctx context.Context, req *applicationpb.GetListingApplicationsRequest) (*applicationpb.GetApplicationsResponse, error) {
	apps, err := h.svc.GetListingApplications(ctx, req.ListingId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &applicationpb.GetApplicationsResponse{}
	for i := range apps {
		resp.Items = append(resp.Items, appToProto(&apps[i]))
	}

	return resp, nil
}

func (h *ApplicationHandler) ChangeStatus(ctx context.Context, req *applicationpb.ChangeStatusRequest) (*applicationpb.ApplicationResponse, error) {
	newStatus := model.ApplicationStatus(req.Status)

	app, err := h.svc.ChangeStatus(ctx, req.Id, req.CompanyId, newStatus, req.Comment)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "application not found")
	}
	if errors.Is(err, repository.ErrInvalidTransition) {
		return nil, status.Error(codes.InvalidArgument, "invalid status transition")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return appToProto(app), nil
}

func (h *ApplicationHandler) GetHistory(ctx context.Context, req *applicationpb.GetHistoryRequest) (*applicationpb.GetHistoryResponse, error) {
	events, err := h.svc.GetHistory(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "application not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &applicationpb.GetHistoryResponse{}
	for i := range events {
		resp.Events = append(resp.Events, eventToProto(&events[i]))
	}

	return resp, nil
}

func appToProto(a *model.Application) *applicationpb.ApplicationResponse {
	resp := &applicationpb.ApplicationResponse{
		Id:          a.ID,
		StudentId:   a.StudentID,
		ListingId:   a.ListingID,
		CoverLetter: a.CoverLetter,
		Status:      string(a.Status),
		CreatedAt:   a.CreatedAt.String(),
		UpdatedAt:   a.UpdatedAt.String(),
	}

	for i := range a.Events {
		resp.Events = append(resp.Events, eventToProto(&a.Events[i]))
	}

	return resp
}

func eventToProto(e *model.ApplicationEvent) *applicationpb.ApplicationEvent {
	return &applicationpb.ApplicationEvent{
		Id:            e.ID,
		ApplicationId: e.ApplicationID,
		OldStatus:     string(e.OldStatus),
		NewStatus:     string(e.NewStatus),
		Comment:       e.Comment,
		ChangedAt:     e.ChangedAt.String(),
	}
}

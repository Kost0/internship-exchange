package handler

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingpb "github.com/Kost0/internship-exchange/proto/listing"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/model"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/service"
)

type ListingHandler struct {
	listingpb.UnimplementedListingServiceServer
	svc *service.ListingService
}

func NewListingHandler(svc *service.ListingService) *ListingHandler {
	return &ListingHandler{svc: svc}
}

func (h *ListingHandler) GetListings(ctx context.Context, req *listingpb.GetListingsRequest) (*listingpb.GetListingsResponse, error) {
	filter := model.ListingsFilter{
		Query:          req.Query,
		Format:         req.Format,
		EmploymentType: req.EmploymentType,
		City:           req.City,
		Skill:          req.Skill,
		Page:           req.Page,
		Limit:          req.Limit,
	}

	listings, total, err := h.svc.GetListings(ctx, filter)
	if err != nil {
		log.Printf("GetListings error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &listingpb.GetListingsResponse{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}
	for i := range listings {
		resp.Items = append(resp.Items, listingToProto(&listings[i]))
	}

	return resp, nil
}

func (h *ListingHandler) GetListing(ctx context.Context, req *listingpb.GetListingRequest) (*listingpb.ListingResponse, error) {
	listing, err := h.svc.GetListing(ctx, req.Id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "listing not found")
	}
	if err != nil {
		log.Printf("GetListing error id=%s: %v", req.Id, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return listingToProto(listing), nil
}

func (h *ListingHandler) GetMyListings(ctx context.Context, req *listingpb.GetMyListingsRequest) (*listingpb.GetMyListingsResponse, error) {
	listings, err := h.svc.GetMyListings(ctx, req.CompanyId)
	if err != nil {
		log.Printf("GetMyListings error companyUserId=%s: %v", req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	resp := &listingpb.GetMyListingsResponse{}
	for i := range listings {
		resp.Items = append(resp.Items, listingToProto(&listings[i]))
	}

	return resp, nil
}

func (h *ListingHandler) CreateListing(ctx context.Context, req *listingpb.CreateListingRequest) (*listingpb.ListingResponse, error) {
	l := model.Listing{
		Title:          req.Title,
		Description:    req.Description,
		Requirements:   req.Requirements,
		WhatWeOffer:    req.WhatWeOffer,
		City:           req.City,
		Format:         model.ListingFormat(req.Format),
		EmploymentType: model.EmploymentType(req.EmploymentType),
		SalaryFrom:     req.SalaryFrom,
		SalaryTo:       req.SalaryTo,
		SalaryCurrency: req.SalaryCurrency,
		Deadline:       req.Deadline,
	}

	listing, err := h.svc.CreateListing(ctx, req.CompanyId, l)
	if err != nil {
		log.Printf("CreateListing error companyUserId=%s: %v", req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return listingToProto(listing), nil
}

func (h *ListingHandler) UpdateListing(ctx context.Context, req *listingpb.UpdateListingRequest) (*listingpb.ListingResponse, error) {
	l := model.Listing{
		Title:          req.Title,
		Description:    req.Description,
		Requirements:   req.Requirements,
		WhatWeOffer:    req.WhatWeOffer,
		City:           req.City,
		Format:         model.ListingFormat(req.Format),
		EmploymentType: model.EmploymentType(req.EmploymentType),
		SalaryFrom:     req.SalaryFrom,
		SalaryTo:       req.SalaryTo,
		SalaryCurrency: req.SalaryCurrency,
		Deadline:       req.Deadline,
	}

	listing, err := h.svc.UpdateListing(ctx, req.Id, req.CompanyId, l)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "listing not found or not a draft")
	}
	if errors.Is(err, repository.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}
	if err != nil {
		log.Printf("UpdateListing error id=%s companyUserId=%s: %v", req.Id, req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return listingToProto(listing), nil
}

func (h *ListingHandler) DeleteListing(ctx context.Context, req *listingpb.DeleteListingRequest) (*listingpb.DeleteListingResponse, error) {
	if err := h.svc.DeleteListing(ctx, req.Id, req.CompanyId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "listing not found or not a draft")
		}
		log.Printf("DeleteListing error id=%s companyUserId=%s: %v", req.Id, req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &listingpb.DeleteListingResponse{Success: true}, nil
}

func (h *ListingHandler) PublishListing(ctx context.Context, req *listingpb.PublishListingRequest) (*listingpb.ListingResponse, error) {
	listing, err := h.svc.PublishListing(ctx, req.Id, req.CompanyId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "listing not found")
	}
	if errors.Is(err, repository.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}
	if errors.Is(err, repository.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, "listing already active")
	}
	if err != nil {
		log.Printf("PublishListing error id=%s companyUserId=%s: %v", req.Id, req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return listingToProto(listing), nil
}

func (h *ListingHandler) CloseListing(ctx context.Context, req *listingpb.CloseListingRequest) (*listingpb.ListingResponse, error) {
	listing, err := h.svc.CloseListing(ctx, req.Id, req.CompanyId)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "listing not found")
	}
	if errors.Is(err, repository.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}
	if errors.Is(err, repository.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, "listing already closed")
	}
	if err != nil {
		log.Printf("CloseListing error id=%s companyUserId=%s: %v", req.Id, req.CompanyId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return listingToProto(listing), nil
}

func listingToProto(l *model.Listing) *listingpb.ListingResponse {
	resp := &listingpb.ListingResponse{
		Id:             l.ID,
		CompanyId:      l.CompanyID,
		Title:          l.Title,
		Description:    l.Description,
		Requirements:   l.Requirements,
		WhatWeOffer:    l.WhatWeOffer,
		City:           l.City,
		Format:         string(l.Format),
		EmploymentType: string(l.EmploymentType),
		SalaryFrom:     l.SalaryFrom,
		SalaryTo:       l.SalaryTo,
		SalaryCurrency: l.SalaryCurrency,
		Deadline:       l.Deadline,
		Status:         string(l.Status),
		CreatedAt:      l.CreatedAt.String(),
		UpdatedAt:      l.UpdatedAt.String(),
	}

	if l.Company != nil {
		resp.Company = &listingpb.CompanyInfo{
			Id:       l.Company.ID,
			UserId:   l.Company.UserID,
			Name:     l.Company.Name,
			LogoUrl:  l.Company.LogoURL,
			Industry: l.Company.Industry,
			City:     l.Company.City,
		}
	}

	for i := range l.Skills {
		resp.Skills = append(resp.Skills, &listingpb.ListingSkill{
			Id:         l.Skills[i].ID,
			ListingId:  l.Skills[i].ListingID,
			Skill:      l.Skills[i].Skill,
			IsRequired: l.Skills[i].IsRequired,
		})
	}

	return resp
}

func (h *ListingHandler) SyncCompany(ctx context.Context, req *listingpb.SyncCompanyRequest) (*listingpb.SyncCompanyResponse, error) {
	err := h.svc.SyncCompany(ctx, req.UserId, req.Name, req.LogoUrl, req.Industry, req.City)
	if err != nil {
		log.Printf("SyncCompany error userId=%s: %v", req.UserId, err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &listingpb.SyncCompanyResponse{Success: true}, nil
}

func (h *ListingHandler) GetCompany(ctx context.Context, req *listingpb.GetCompanyRequest) (*listingpb.CompanyResponse, error) {
	company, err := h.svc.GetCompany(ctx, req.Id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "company not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &listingpb.CompanyResponse{
		Id: company.ID, UserId: company.UserID, Name: company.Name,
		LogoUrl: company.LogoURL, Industry: company.Industry, City: company.City,
	}, nil
}

package service

import (
	"context"
	"errors"

	"github.com/Kost0/internship-exchange/services/listing-service/internal/model"
	"github.com/Kost0/internship-exchange/services/listing-service/internal/repository"
)

type ListingService struct {
	listings  *repository.ListingRepository
	companies *repository.CompanyRepository
}

func NewListingService(
	listings *repository.ListingRepository,
	companies *repository.CompanyRepository,
) *ListingService {
	return &ListingService{listings: listings, companies: companies}
}

func (s *ListingService) GetListings(ctx context.Context, f model.ListingsFilter) ([]model.Listing, int64, error) {
	return s.listings.GetAll(ctx, f)
}

func (s *ListingService) GetListing(ctx context.Context, id string) (*model.Listing, error) {
	return s.listings.GetByID(ctx, id)
}

func (s *ListingService) GetMyListings(ctx context.Context, userID string) ([]model.Listing, error) {
	company, err := s.companies.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.listings.GetByCompanyID(ctx, company.ID)
}

func (s *ListingService) CreateListing(ctx context.Context, userID string, l model.Listing) (*model.Listing, error) {
	company, err := s.companies.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	l.CompanyID = company.ID

	return s.listings.Create(ctx, l)
}

func (s *ListingService) UpdateListing(ctx context.Context, id, userID string, l model.Listing) (*model.Listing, error) {
	company, err := s.companies.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result, err := s.listings.Update(ctx, id, company.ID, l)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, repository.ErrNotFound
	}

	return result, err
}

func (s *ListingService) DeleteListing(ctx context.Context, id, userID string) error {
	company, err := s.companies.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.listings.Delete(ctx, id, company.ID)
}

func (s *ListingService) PublishListing(ctx context.Context, id, userID string) (*model.Listing, error) {
	company, err := s.companies.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	listing, err := s.listings.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if listing.CompanyID != company.ID {
		return nil, repository.ErrForbidden
	}

	if listing.Status == model.StatusActive {
		return nil, repository.ErrConflict
	}

	return s.listings.SetStatus(ctx, id, company.ID, model.StatusActive)
}

func (s *ListingService) CloseListing(ctx context.Context, id, userID string) (*model.Listing, error) {
	company, err := s.companies.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	listing, err := s.listings.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if listing.CompanyID != company.ID {
		return nil, repository.ErrForbidden
	}

	if listing.Status == model.StatusClosed {
		return nil, repository.ErrConflict
	}

	return s.listings.SetStatus(ctx, id, company.ID, model.StatusClosed)
}

func (s *ListingService) SyncCompany(ctx context.Context, userID, name, logo, industry, city string) error {
	return s.companies.SyncFromProfile(ctx, userID, name, logo, industry, city)
}

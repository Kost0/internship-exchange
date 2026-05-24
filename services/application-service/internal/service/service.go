package service

import (
	"context"
	"errors"

	"github.com/Kost0/internship-exchange/services/application-service/internal/model"
	"github.com/Kost0/internship-exchange/services/application-service/internal/publisher"
	"github.com/Kost0/internship-exchange/services/application-service/internal/repository"
)

type ApplicationService struct {
	repo      *repository.ApplicationRepository
	publisher *publisher.Publisher
}

func NewApplicationService(repo *repository.ApplicationRepository, pub *publisher.Publisher) *ApplicationService {
	return &ApplicationService{repo: repo, publisher: pub}
}

func (s *ApplicationService) Apply(ctx context.Context, studentID, listingID, coverLetter, studentEmail, companyEmail string) (*model.Application, error) {
	app, err := s.repo.Create(ctx, studentID, listingID, coverLetter, studentEmail, companyEmail)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishApplicationCreated(ctx, app.ID, app.StudentID, app.ListingID, app.CompanyEmail)

	return app, nil
}

func (s *ApplicationService) GetMyApplications(ctx context.Context, studentID string) ([]model.Application, error) {
	apps, err := s.repo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.LoadEvents(ctx, apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func (s *ApplicationService) Withdraw(ctx context.Context, id, studentID string) error {
	return s.repo.Delete(ctx, id, studentID)
}

func (s *ApplicationService) GetListingApplications(ctx context.Context, listingID string) ([]model.Application, error) {
	apps, err := s.repo.GetByListingID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.LoadEvents(ctx, apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func (s *ApplicationService) ChangeStatus(ctx context.Context, id, companyID string, newStatus model.ApplicationStatus, comment string) (*model.Application, error) {
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !app.Status.CanTransitionTo(newStatus) {
		return nil, repository.ErrInvalidTransition
	}

	oldStatus := app.Status

	app, err = s.repo.UpdateStatus(ctx, id, newStatus, comment)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.AddEvent(ctx, id, oldStatus, newStatus, comment)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishStatusChanged(ctx,
		app.ID, app.StudentID, app.ListingID,
		string(oldStatus), string(newStatus), comment, app.StudentEmail,
	)

	events, err := s.repo.GetEvents(ctx, id)
	if err != nil {
		return nil, err
	}
	app.Events = events

	return app, nil
}

func (s *ApplicationService) GetHistory(ctx context.Context, id string) ([]model.ApplicationEvent, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, repository.ErrNotFound
		}

		return nil, err
	}

	return s.repo.GetEvents(ctx, id)
}

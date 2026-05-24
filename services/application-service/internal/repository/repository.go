package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kost0/internship-exchange/services/application-service/internal/model"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("already applied")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type ApplicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(ctx context.Context, studentID, listingID, coverLetter string) (*model.Application, error) {
	query := `
		INSERT INTO applications (student_id, listing_id, cover_letter)
		VALUES ($1, $2, $3)
		RETURNING id, student_id, listing_id, cover_letter, status, created_at, updated_at
	`

	app := &model.Application{}
	err := r.db.QueryRow(ctx, query, studentID, listingID, coverLetter).Scan(
		&app.ID, &app.StudentID, &app.ListingID,
		&app.CoverLetter, &app.Status,
		&app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrConflict
		}

		return nil, err
	}

	return app, nil
}

func (r *ApplicationRepository) GetByStudentID(ctx context.Context, studentID string) ([]model.Application, error) {
	query := `
		SELECT id, student_id, listing_id, cover_letter, status, created_at, updated_at
		FROM applications
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Application
	for rows.Next() {
		app := model.Application{}

		if err := rows.Scan(
			&app.ID, &app.StudentID, &app.ListingID,
			&app.CoverLetter, &app.Status,
			&app.CreatedAt, &app.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, app)
	}

	return result, nil
}

func (r *ApplicationRepository) GetByListingID(ctx context.Context, listingID string) ([]model.Application, error) {
	query := `
		SELECT id, student_id, listing_id, cover_letter, status, created_at, updated_at
		FROM applications
		WHERE listing_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Application
	for rows.Next() {
		app := model.Application{}

		if err := rows.Scan(
			&app.ID, &app.StudentID, &app.ListingID,
			&app.CoverLetter, &app.Status,
			&app.CreatedAt, &app.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, app)
	}

	return result, nil
}

func (r *ApplicationRepository) GetByID(ctx context.Context, id string) (*model.Application, error) {
	query := `
		SELECT id, student_id, listing_id, cover_letter, status, created_at, updated_at
		FROM applications
		WHERE id = $1
	`

	app := &model.Application{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&app.ID, &app.StudentID, &app.ListingID,
		&app.CoverLetter, &app.Status,
		&app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return app, nil
}

func (r *ApplicationRepository) Delete(ctx context.Context, id, studentID string) error {
	query := `
		DELETE FROM applications
		WHERE id = $1 AND student_id = $2 AND status = 'applied'
	`

	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id string, newStatus model.ApplicationStatus, comment string) (*model.Application, error) {
	query := `
		UPDATE applications
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, student_id, listing_id, cover_letter, status, created_at, updated_at
	`

	app := &model.Application{}

	err := r.db.QueryRow(ctx, query, newStatus, id).Scan(
		&app.ID, &app.StudentID, &app.ListingID,
		&app.CoverLetter, &app.Status,
		&app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return app, nil
}

func (r *ApplicationRepository) AddEvent(ctx context.Context, applicationID string, oldStatus, newStatus model.ApplicationStatus, comment string) (*model.ApplicationEvent, error) {
	query := `
		INSERT INTO application_events (application_id, old_status, new_status, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, application_id, old_status, new_status, comment, changed_at
	`

	event := &model.ApplicationEvent{}
	var dbComment *string

	err := r.db.QueryRow(ctx, query, applicationID, oldStatus, newStatus, nilIfEmpty(comment)).Scan(
		&event.ID, &event.ApplicationID,
		&event.OldStatus, &event.NewStatus,
		&dbComment, &event.ChangedAt,
	)
	if err != nil {
		return nil, err
	}

	event.Comment = derefStr(dbComment)
	return event, nil
}

func (r *ApplicationRepository) GetEvents(ctx context.Context, applicationID string) ([]model.ApplicationEvent, error) {
	query := `
		SELECT id, application_id, old_status, new_status, comment, changed_at
		FROM application_events
		WHERE application_id = $1
		ORDER BY changed_at ASC
	`

	rows, err := r.db.Query(ctx, query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.ApplicationEvent
	for rows.Next() {
		e := model.ApplicationEvent{}
		var dbComment *string

		if err := rows.Scan(
			&e.ID, &e.ApplicationID,
			&e.OldStatus, &e.NewStatus,
			&dbComment, &e.ChangedAt,
		); err != nil {
			return nil, err
		}

		e.Comment = derefStr(dbComment)
		events = append(events, e)
	}

	return events, nil
}

func (r *ApplicationRepository) LoadEvents(ctx context.Context, apps []model.Application) error {
	if len(apps) == 0 {
		return nil
	}

	ids := make([]string, len(apps))
	for i, a := range apps {
		ids[i] = a.ID
	}

	query := `
		SELECT id, application_id, old_status, new_status, comment, changed_at
		FROM application_events
		WHERE application_id = ANY($1)
		ORDER BY changed_at ASC
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return err
	}
	defer rows.Close()

	eventsMap := make(map[string][]model.ApplicationEvent)
	for rows.Next() {
		e := model.ApplicationEvent{}
		var dbComment *string

		if err := rows.Scan(
			&e.ID, &e.ApplicationID,
			&e.OldStatus, &e.NewStatus,
			&dbComment, &e.ChangedAt,
		); err != nil {
			return err
		}

		e.Comment = derefStr(dbComment)
		eventsMap[e.ApplicationID] = append(eventsMap[e.ApplicationID], e)
	}

	for i := range apps {
		apps[i].Events = eventsMap[apps[i].ID]
	}

	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "unique") || contains(err.Error(), "duplicate"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}

		return false
	}()
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

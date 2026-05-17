package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kost0/internship-exchange/services/listing-service/internal/model"
)

var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")
var ErrConflict = errors.New("conflict")

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) GetOrCreate(ctx context.Context, userID string) (*model.Company, error) {
	query := `
		INSERT INTO companies (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, name, logo_url, industry, city, created_at, updated_at
	`

	c := &model.Company{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&c.ID, &c.UserID, &c.Name, &c.LogoURL,
		&c.Industry, &c.City, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *CompanyRepository) GetByUserID(ctx context.Context, userID string) (*model.Company, error) {
	query := `
		SELECT id, user_id, name, logo_url, industry, city, created_at, updated_at
		FROM companies
		WHERE user_id = $1
	`

	c := &model.Company{}

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&c.ID, &c.UserID, &c.Name, &c.LogoURL,
		&c.Industry, &c.City, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return c, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id string) (*model.Company, error) {
	query := `
		SELECT id, user_id, name, logo_url, industry, city, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	c := &model.Company{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.UserID, &c.Name, &c.LogoURL,
		&c.Industry, &c.City, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return c, nil
}

func (r *CompanyRepository) SyncFromProfile(ctx context.Context, userID, name, logoURL, industry, city string) error {
	query := `
		INSERT INTO companies (user_id, name, logo_url, industry, city)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
		SET name = $2, logo_url = $3, industry = $4, city = $5, updated_at = NOW()
	`

	_, err := r.db.Exec(ctx, query, userID, name, logoURL, industry, city)
	return err
}

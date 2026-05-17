package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kost0/internship-exchange/services/profile-service/internal/model"
)

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
		RETURNING id, user_id, name, tagline, description, industry, size,
		          founded_year, website, contact_email, city, country,
		          is_remote_friendly, logo_url, created_at, updated_at
	`

	c := &model.Company{}
	var tagline, description, industry, size, website, contactEmail, city, country, logoURL *string
	var foundedYear *int32

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&c.ID, &c.UserID, &c.Name,
		&tagline, &description, &industry, &size,
		&foundedYear, &website, &contactEmail, &city, &country,
		&c.IsRemoteFriendly, &logoURL,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.Tagline = derefStr(tagline)
	c.Description = derefStr(description)
	c.Industry = derefStr(industry)
	c.Size = derefStr(size)
	c.Website = derefStr(website)
	c.ContactEmail = derefStr(contactEmail)
	c.City = derefStr(city)
	c.Country = derefStr(country)
	c.LogoURL = derefStr(logoURL)
	if foundedYear != nil {
		c.FoundedYear = *foundedYear
	}

	return c, nil
}

func (r *CompanyRepository) GetByUserID(ctx context.Context, userID string) (*model.Company, error) {
	query := `
		SELECT id, user_id, name, tagline, description, industry, size,
		       founded_year, website, contact_email, city, country,
		       is_remote_friendly, logo_url, created_at, updated_at
		FROM companies
		WHERE user_id = $1
	`

	c := &model.Company{}
	var tagline, description, industry, size, website, contactEmail, city, country, logoURL *string
	var foundedYear *int32

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&c.ID, &c.UserID, &c.Name,
		&tagline, &description, &industry, &size,
		&foundedYear, &website, &contactEmail, &city, &country,
		&c.IsRemoteFriendly, &logoURL,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	c.Tagline = derefStr(tagline)
	c.Description = derefStr(description)
	c.Industry = derefStr(industry)
	c.Size = derefStr(size)
	c.Website = derefStr(website)
	c.ContactEmail = derefStr(contactEmail)
	c.City = derefStr(city)
	c.Country = derefStr(country)
	c.LogoURL = derefStr(logoURL)
	if foundedYear != nil {
		c.FoundedYear = *foundedYear
	}

	return c, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id string) (*model.Company, error) {
	query := `
		SELECT id, user_id, name, tagline, description, industry, size,
		       founded_year, website, contact_email, city, country,
		       is_remote_friendly, logo_url, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	c := &model.Company{}
	var tagline, description, industry, size, website, contactEmail, city, country, logoURL *string
	var foundedYear *int32

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.UserID, &c.Name,
		&tagline, &description, &industry, &size,
		&foundedYear, &website, &contactEmail, &city, &country,
		&c.IsRemoteFriendly, &logoURL,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	c.Tagline = derefStr(tagline)
	c.Description = derefStr(description)
	c.Industry = derefStr(industry)
	c.Size = derefStr(size)
	c.Website = derefStr(website)
	c.ContactEmail = derefStr(contactEmail)
	c.City = derefStr(city)
	c.Country = derefStr(country)
	c.LogoURL = derefStr(logoURL)
	if foundedYear != nil {
		c.FoundedYear = *foundedYear
	}

	return c, nil
}

func (r *CompanyRepository) Update(ctx context.Context, userID string, c model.Company) (*model.Company, error) {
	query := `
		UPDATE companies
		SET name               = $1,
		    tagline            = $2,
		    description        = $3,
		    industry           = $4,
		    size               = $5,
		    founded_year       = $6,
		    website            = $7,
		    contact_email      = $8,
		    city               = $9,
		    country            = $10,
		    is_remote_friendly = $11,
		    updated_at         = NOW()
		WHERE user_id = $12
		RETURNING id, user_id, name, tagline, description, industry, size,
		          founded_year, website, contact_email, city, country,
		          is_remote_friendly, logo_url, created_at, updated_at
	`

	result := &model.Company{}
	var tagline, description, industry, size, website, contactEmail, city, country, logoURL *string
	var foundedYear *int32

	err := r.db.QueryRow(ctx, query,
		c.Name, nilIfEmpty(c.Tagline), nilIfEmpty(c.Description),
		nilIfEmpty(c.Industry), nilIfEmpty(c.Size),
		nullIfZero(int(c.FoundedYear)),
		nilIfEmpty(c.Website), nilIfEmpty(c.ContactEmail),
		nilIfEmpty(c.City), nilIfEmpty(c.Country),
		c.IsRemoteFriendly,
		userID,
	).Scan(
		&result.ID, &result.UserID, &result.Name,
		&tagline, &description, &industry, &size,
		&foundedYear, &website, &contactEmail, &city, &country,
		&result.IsRemoteFriendly, &logoURL,
		&result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	result.Tagline = derefStr(tagline)
	result.Description = derefStr(description)
	result.Industry = derefStr(industry)
	result.Size = derefStr(size)
	result.Website = derefStr(website)
	result.ContactEmail = derefStr(contactEmail)
	result.City = derefStr(city)
	result.Country = derefStr(country)
	result.LogoURL = derefStr(logoURL)
	if foundedYear != nil {
		result.FoundedYear = *foundedYear
	}

	return result, nil
}

func (r *CompanyRepository) SetLogoURL(ctx context.Context, userID, url string) error {
	query := `UPDATE companies SET logo_url = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, url, userID)
	return err
}

func nullIfZero(n int) any {
	if n == 0 {
		return nil
	}
	return n
}

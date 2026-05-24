package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kost0/internship-exchange/services/listing-service/internal/model"
)

type ListingRepository struct {
	db *pgxpool.Pool
}

func NewListingRepository(db *pgxpool.Pool) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) GetAll(ctx context.Context, f model.ListingsFilter) ([]model.Listing, int64, error) {
	conditions := []string{"l.status = 'active'"}
	args := make([]any, 0)
	idx := 1

	if f.Query != "" {
		conditions = append(conditions, fmt.Sprintf("(l.title ILIKE $%d OR l.description ILIKE $%d)", idx, idx+1))
		like := "%" + f.Query + "%"
		args = append(args, like, like)
		idx += 2
	}

	if f.Format != "" {
		conditions = append(conditions, fmt.Sprintf("l.format = $%d", idx))
		args = append(args, f.Format)
		idx++
	}

	if f.EmploymentType != "" {
		conditions = append(conditions, fmt.Sprintf("l.employment_type = $%d", idx))
		args = append(args, f.EmploymentType)
		idx++
	}

	if f.City != "" {
		conditions = append(conditions, fmt.Sprintf("l.city ILIKE $%d", idx))
		args = append(args, "%"+f.City+"%")
		idx++
	}

	if f.Skill != "" {
		conditions = append(conditions, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM listing_skills ls WHERE ls.listing_id = l.id AND ls.skill ILIKE $%d)", idx,
		))
		args = append(args, "%"+f.Skill+"%")
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM listings l %s`, where)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.Limit == 0 {
		f.Limit = 12
	}
	if f.Page == 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	listingsQuery := fmt.Sprintf(`
	SELECT l.id, l.company_id, l.title, l.description, l.requirements, l.what_we_offer,
	       l.city, l.format, l.employment_type, l.salary_from, l.salary_to, l.salary_currency,
	       l.deadline::text, l.status, l.created_at, l.updated_at,
	       c.id,
			c.user_id,
			COALESCE(c.name, ''),
			COALESCE(c.logo_url, ''),
			COALESCE(c.industry, ''),
			COALESCE(c.city, '')
	FROM listings l
	JOIN companies c ON c.id = l.company_id
	%s
	ORDER BY l.created_at DESC
	LIMIT $%d OFFSET $%d
`, where, idx, idx+1)

	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(ctx, listingsQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var listings []model.Listing
	for rows.Next() {
		l, err := scanListing(rows)

		if err != nil {
			return nil, 0, err
		}

		listings = append(listings, *l)
	}

	if err := r.loadSkillsForListings(ctx, listings); err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

func (r *ListingRepository) GetByID(ctx context.Context, id string) (*model.Listing, error) {
	query := `
	SELECT l.id, l.company_id, l.title, l.description, l.requirements, l.what_we_offer,
	       l.city, l.format, l.employment_type, l.salary_from, l.salary_to, l.salary_currency,
	       l.deadline::text, l.status, l.created_at, l.updated_at,
	       c.id, c.user_id,
	       COALESCE(c.name, ''),
	       COALESCE(c.logo_url, ''),
	       COALESCE(c.industry, ''),
	       COALESCE(c.city, '')
	FROM listings l
	JOIN companies c ON c.id = l.company_id
	WHERE l.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	l, err := scanListing(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	skills, err := r.getSkills(ctx, l.ID)
	if err != nil {
		return nil, err
	}

	l.Skills = skills

	return l, nil
}

func (r *ListingRepository) GetByCompanyID(ctx context.Context, companyID string) ([]model.Listing, error) {
	query := `
	SELECT l.id, l.company_id, l.title, l.description, l.requirements, l.what_we_offer,
	       l.city, l.format, l.employment_type, l.salary_from, l.salary_to, l.salary_currency,
	       l.deadline::text, l.status, l.created_at, l.updated_at,
	       c.id, c.user_id,
	       COALESCE(c.name, ''),
	       COALESCE(c.logo_url, ''),
	       COALESCE(c.industry, ''),
	       COALESCE(c.city, '')
	FROM listings l
	JOIN companies c ON c.id = l.company_id
	WHERE l.company_id = $1
	ORDER BY l.created_at DESC
`

	rows, err := r.db.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listings []model.Listing
	for rows.Next() {
		l, err := scanListing(rows)

		if err != nil {
			return nil, err
		}

		listings = append(listings, *l)
	}

	if err := r.loadSkillsForListings(ctx, listings); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *ListingRepository) Create(ctx context.Context, l model.Listing) (*model.Listing, error) {
	query := `
		INSERT INTO listings (company_id, title, description, requirements, what_we_offer,
		                      city, format, employment_type, salary_from, salary_to,
		                      salary_currency, deadline)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, company_id, title, description, requirements, what_we_offer,
		          city, format, employment_type, salary_from, salary_to, salary_currency,
		          deadline::text, status, created_at, updated_at
	`

	result := &model.Listing{}

	var city *string
	var salaryFrom *int64
	var salaryTo *int64
	var salaryCurrency *string
	var deadline *string

	err := r.db.QueryRow(ctx, query,
		l.CompanyID, l.Title, l.Description, l.Requirements, l.WhatWeOffer,
		nilIfEmpty(l.City), l.Format, l.EmploymentType,
		nullIfZeroInt(l.SalaryFrom), nullIfZeroInt(l.SalaryTo),
		nilIfEmpty(l.SalaryCurrency), nilIfEmpty(l.Deadline),
	).Scan(
		&result.ID, &result.CompanyID, &result.Title, &result.Description,
		&result.Requirements, &result.WhatWeOffer,
		&city, &result.Format, &result.EmploymentType,
		&salaryFrom, &salaryTo, &salaryCurrency,
		&deadline, &result.Status,
		&result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	result.City = derefStr(city)
	result.SalaryFrom = derefInt64(salaryFrom)
	result.SalaryTo = derefInt64(salaryTo)
	result.SalaryCurrency = derefStr(salaryCurrency)
	result.Deadline = derefStr(deadline)

	return result, nil
}

func (r *ListingRepository) Update(ctx context.Context, id, companyID string, l model.Listing) (*model.Listing, error) {
	query := `
		UPDATE listings
		SET title           = $1,
		    description     = $2,
		    requirements    = $3,
		    what_we_offer   = $4,
		    city            = $5,
		    format          = $6,
		    employment_type = $7,
		    salary_from     = $8,
		    salary_to       = $9,
		    salary_currency = $10,
		    deadline        = $11,
		    updated_at      = NOW()
		WHERE id = $12 AND company_id = $13 AND status = 'draft'
		RETURNING id, company_id, title, description, requirements, what_we_offer,
		          city, format, employment_type, salary_from, salary_to, salary_currency,
		          deadline::text, status, created_at, updated_at
	`

	result := &model.Listing{}

	var city *string
	var salaryFrom *int64
	var salaryTo *int64
	var salaryCurrency *string
	var deadline *string

	err := r.db.QueryRow(ctx, query,
		l.Title, l.Description, l.Requirements, l.WhatWeOffer,
		nilIfEmpty(l.City), l.Format, l.EmploymentType,
		nullIfZeroInt(l.SalaryFrom), nullIfZeroInt(l.SalaryTo),
		nilIfEmpty(l.SalaryCurrency), nilIfEmpty(l.Deadline),
		id, companyID,
	).Scan(
		&result.ID, &result.CompanyID, &result.Title, &result.Description,
		&result.Requirements, &result.WhatWeOffer,
		&city, &result.Format, &result.EmploymentType,
		&salaryFrom, &salaryTo, &salaryCurrency,
		&deadline, &result.Status,
		&result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result.City = derefStr(city)
	result.SalaryFrom = derefInt64(salaryFrom)
	result.SalaryTo = derefInt64(salaryTo)
	result.SalaryCurrency = derefStr(salaryCurrency)
	result.Deadline = derefStr(deadline)

	return result, nil
}
func (r *ListingRepository) Delete(ctx context.Context, id, companyID string) error {
	query := `DELETE FROM listings WHERE id = $1 AND company_id = $2 AND status = 'draft'`
	tag, err := r.db.Exec(ctx, query, id, companyID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *ListingRepository) SetStatus(ctx context.Context, id, companyID string, status model.ListingStatus) (*model.Listing, error) {
	query := `
		UPDATE listings
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND company_id = $3
		RETURNING id, company_id, title, description, requirements, what_we_offer,
		          city, format, employment_type, salary_from, salary_to, salary_currency,
		          deadline::text, status, created_at, updated_at
	`

	result := &model.Listing{}

	var city *string
	var salaryFrom *int64
	var salaryTo *int64
	var salaryCurrency *string
	var deadline *string

	err := r.db.QueryRow(ctx, query, status, id, companyID).Scan(
		&result.ID, &result.CompanyID, &result.Title, &result.Description,
		&result.Requirements, &result.WhatWeOffer,
		&city, &result.Format, &result.EmploymentType,
		&salaryFrom, &salaryTo, &salaryCurrency,
		&deadline, &result.Status,
		&result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result.City = derefStr(city)
	result.SalaryFrom = derefInt64(salaryFrom)
	result.SalaryTo = derefInt64(salaryTo)
	result.SalaryCurrency = derefStr(salaryCurrency)
	result.Deadline = derefStr(deadline)

	return result, nil
}

func (r *ListingRepository) getSkills(ctx context.Context, listingID string) ([]model.ListingSkill, error) {
	query := `
		SELECT id, listing_id, skill, is_required
		FROM listing_skills
		WHERE listing_id = $1
	`

	rows, err := r.db.Query(ctx, query, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []model.ListingSkill

	for rows.Next() {
		s := model.ListingSkill{}
		if err := rows.Scan(&s.ID, &s.ListingID, &s.Skill, &s.IsRequired); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, nil
}

func (r *ListingRepository) loadSkillsForListings(ctx context.Context, listings []model.Listing) error {
	if len(listings) == 0 {
		return nil
	}

	ids := make([]string, len(listings))
	for i, l := range listings {
		ids[i] = l.ID
	}

	query := `
		SELECT id, listing_id, skill, is_required
		FROM listing_skills
		WHERE listing_id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return err
	}
	defer rows.Close()

	skillsMap := make(map[string][]model.ListingSkill)
	for rows.Next() {
		s := model.ListingSkill{}
		if err := rows.Scan(&s.ID, &s.ListingID, &s.Skill, &s.IsRequired); err != nil {
			return err
		}
		skillsMap[s.ListingID] = append(skillsMap[s.ListingID], s)
	}

	for i := range listings {
		listings[i].Skills = skillsMap[listings[i].ID]
	}

	return nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanListing(row scannable) (*model.Listing, error) {
	l := &model.Listing{Company: &model.Company{}}

	var city *string
	var salaryFrom *int64
	var salaryTo *int64
	var salaryCurrency *string
	var deadline *string

	var cUserID *string
	var cName *string
	var cLogo *string
	var cIndustry *string
	var cCity *string

	err := row.Scan(
		&l.ID, &l.CompanyID, &l.Title, &l.Description, &l.Requirements, &l.WhatWeOffer,
		&city, &l.Format, &l.EmploymentType, &salaryFrom, &salaryTo,
		&salaryCurrency, &deadline, &l.Status, &l.CreatedAt, &l.UpdatedAt,
		&l.Company.ID, &cUserID, &cName, &cLogo, &cIndustry, &cCity,
	)
	if err != nil {
		return nil, err
	}

	l.City = derefStr(city)
	l.SalaryFrom = derefInt64(salaryFrom)
	l.SalaryTo = derefInt64(salaryTo)
	l.SalaryCurrency = derefStr(salaryCurrency)
	l.Deadline = derefStr(deadline)

	l.Company.UserID = derefStr(cUserID)
	l.Company.Name = derefStr(cName)
	l.Company.LogoURL = derefStr(cLogo)
	l.Company.Industry = derefStr(cIndustry)
	l.Company.City = derefStr(cCity)

	return l, nil
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}

func nullIfZeroInt(n int64) any {
	if n == 0 {
		return nil
	}

	return n
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt64(n *int64) int64 {
	if n == nil {
		return 0
	}
	return *n
}

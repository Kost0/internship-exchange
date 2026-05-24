package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kost0/internship-exchange/services/profile-service/internal/model"
)

var ErrNotFound = errors.New("not found")

type StudentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetOrCreate(ctx context.Context, userID string) (*model.Student, error) {
	query := `
		INSERT INTO students (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO UPDATE SET updated_at = NOW()
		RETURNING id, user_id, first_name, last_name, phone, city, bio,
		          avatar_url, resume_url, github_url, linkedin_url, portfolio_url,
		          created_at, updated_at
	`

	s := &model.Student{}
	var phone, city, bio, avatarURL, resumeURL, githubURL, linkedinURL, portfolioURL *string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.FirstName, &s.LastName,
		&phone, &city, &bio,
		&avatarURL, &resumeURL, &githubURL, &linkedinURL, &portfolioURL,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	s.Phone = derefStr(phone)
	s.City = derefStr(city)
	s.Bio = derefStr(bio)
	s.AvatarURL = derefStr(avatarURL)
	s.ResumeURL = derefStr(resumeURL)
	s.GithubURL = derefStr(githubURL)
	s.LinkedinURL = derefStr(linkedinURL)
	s.PortfolioURL = derefStr(portfolioURL)

	return s, nil
}

func (r *StudentRepository) GetByUserID(ctx context.Context, userID string) (*model.Student, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, city, bio,
		       avatar_url, resume_url, github_url, linkedin_url, portfolio_url,
		       created_at, updated_at
		FROM students
		WHERE user_id = $1
	`

	s := &model.Student{}
	var phone, city, bio, avatarURL, resumeURL, githubURL, linkedinURL, portfolioURL *string

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.FirstName, &s.LastName,
		&phone, &city, &bio,
		&avatarURL, &resumeURL, &githubURL, &linkedinURL, &portfolioURL,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	s.Phone = derefStr(phone)
	s.City = derefStr(city)
	s.Bio = derefStr(bio)
	s.AvatarURL = derefStr(avatarURL)
	s.ResumeURL = derefStr(resumeURL)
	s.GithubURL = derefStr(githubURL)
	s.LinkedinURL = derefStr(linkedinURL)
	s.PortfolioURL = derefStr(portfolioURL)

	return s, nil
}

func (r *StudentRepository) GetByID(ctx context.Context, id string) (*model.Student, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, city, bio,
		       avatar_url, resume_url, github_url, linkedin_url, portfolio_url,
		       created_at, updated_at
		FROM students
		WHERE id = $1
	`

	s := &model.Student{}
	var phone, city, bio, avatarURL, resumeURL, githubURL, linkedinURL, portfolioURL *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.UserID, &s.FirstName, &s.LastName,
		&phone, &city, &bio,
		&avatarURL, &resumeURL, &githubURL, &linkedinURL, &portfolioURL,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	s.Phone = derefStr(phone)
	s.City = derefStr(city)
	s.Bio = derefStr(bio)
	s.AvatarURL = derefStr(avatarURL)
	s.ResumeURL = derefStr(resumeURL)
	s.GithubURL = derefStr(githubURL)
	s.LinkedinURL = derefStr(linkedinURL)
	s.PortfolioURL = derefStr(portfolioURL)

	return s, nil
}

func (r *StudentRepository) Update(ctx context.Context, userID string, fields map[string]any) (*model.Student, error) {
	query := `
		UPDATE students
		SET first_name    = $1,
		    last_name     = $2,
		    phone         = $3,
		    city          = $4,
		    bio           = $5,
		    github_url    = $6,
		    linkedin_url  = $7,
		    portfolio_url = $8,
		    updated_at    = NOW()
		WHERE user_id = $9
		RETURNING id, user_id, first_name, last_name, phone, city, bio,
		          avatar_url, resume_url, github_url, linkedin_url, portfolio_url,
		          created_at, updated_at
	`

	s := &model.Student{}
	var phone, city, bio, avatarURL, resumeURL, githubURL, linkedinURL, portfolioURL *string

	err := r.db.QueryRow(ctx, query,
		fields["first_name"], fields["last_name"],
		fields["phone"], fields["city"], fields["bio"],
		fields["github_url"], fields["linkedin_url"], fields["portfolio_url"],
		userID,
	).Scan(
		&s.ID, &s.UserID, &s.FirstName, &s.LastName,
		&phone, &city, &bio,
		&avatarURL, &resumeURL, &githubURL, &linkedinURL, &portfolioURL,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	s.Phone = derefStr(phone)
	s.City = derefStr(city)
	s.Bio = derefStr(bio)
	s.AvatarURL = derefStr(avatarURL)
	s.ResumeURL = derefStr(resumeURL)
	s.GithubURL = derefStr(githubURL)
	s.LinkedinURL = derefStr(linkedinURL)
	s.PortfolioURL = derefStr(portfolioURL)

	return s, nil
}

func (r *StudentRepository) SetAvatarURL(ctx context.Context, userID, url string) error {
	query := `UPDATE students SET avatar_url = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, url, userID)
	return err
}

func (r *StudentRepository) SetResumeURL(ctx context.Context, userID, url string) error {
	query := `UPDATE students SET resume_url = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, url, userID)
	return err
}

func (r *StudentRepository) GetEducations(ctx context.Context, studentID string) ([]model.Education, error) {
	query := `
		SELECT id, student_id, university, faculty, specialization,
		       degree, start_year, end_year, gpa, is_current
		FROM educations
		WHERE student_id = $1
		ORDER BY start_year DESC
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Education
	for rows.Next() {
		e := model.Education{}
		var faculty, specialization, degree *string
		var endYear *int32
		var gpa *float64

		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.University,
			&faculty, &specialization, &degree,
			&e.StartYear, &endYear, &gpa, &e.IsCurrent,
		); err != nil {
			return nil, err
		}

		e.Faculty = derefStr(faculty)
		e.Specialization = derefStr(specialization)
		e.Degree = derefStr(degree)
		if endYear != nil {
			e.EndYear = *endYear
		}
		if gpa != nil {
			e.GPA = *gpa
		}

		result = append(result, e)
	}

	return result, nil
}

func (r *StudentRepository) AddEducation(ctx context.Context, studentID string, e model.Education) (*model.Education, error) {
	query := `
		INSERT INTO educations (student_id, university, faculty, specialization, degree, start_year, end_year, gpa, is_current)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, student_id, university, faculty, specialization, degree, start_year, end_year, gpa, is_current
	`

	result := &model.Education{}
	var faculty, specialization, degree *string
	var endYear *int32
	var gpa *float64

	err := r.db.QueryRow(ctx, query,
		studentID, e.University,
		nilIfEmpty(e.Faculty), nilIfEmpty(e.Specialization), nilIfEmpty(e.Degree),
		e.StartYear, nullIfZeroInt32(e.EndYear), nullIfZeroFloat(e.GPA), e.IsCurrent,
	).Scan(
		&result.ID, &result.StudentID, &result.University,
		&faculty, &specialization, &degree,
		&result.StartYear, &endYear, &gpa, &result.IsCurrent,
	)
	if err != nil {
		return nil, err
	}

	result.Faculty = derefStr(faculty)
	result.Specialization = derefStr(specialization)
	result.Degree = derefStr(degree)
	if endYear != nil {
		result.EndYear = *endYear
	}
	if gpa != nil {
		result.GPA = *gpa
	}

	return result, nil
}

func (r *StudentRepository) UpdateEducation(ctx context.Context, id, studentID string, e model.Education) (*model.Education, error) {
	query := `
		UPDATE educations
		SET university = $1, faculty = $2, specialization = $3,
		    degree = $4, start_year = $5, end_year = $6, gpa = $7, is_current = $8
		WHERE id = $9 AND student_id = $10
		RETURNING id, student_id, university, faculty, specialization, degree, start_year, end_year, gpa, is_current
	`

	result := &model.Education{}
	var faculty, specialization, degree *string
	var endYear *int32
	var gpa *float64

	err := r.db.QueryRow(ctx, query,
		e.University, nilIfEmpty(e.Faculty), nilIfEmpty(e.Specialization),
		nilIfEmpty(e.Degree), e.StartYear, nullIfZeroInt32(e.EndYear),
		nullIfZeroFloat(e.GPA), e.IsCurrent,
		id, studentID,
	).Scan(
		&result.ID, &result.StudentID, &result.University,
		&faculty, &specialization, &degree,
		&result.StartYear, &endYear, &gpa, &result.IsCurrent,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	result.Faculty = derefStr(faculty)
	result.Specialization = derefStr(specialization)
	result.Degree = derefStr(degree)
	if endYear != nil {
		result.EndYear = *endYear
	}
	if gpa != nil {
		result.GPA = *gpa
	}

	return result, nil
}

func (r *StudentRepository) DeleteEducation(ctx context.Context, id, studentID string) error {
	query := `DELETE FROM educations WHERE id = $1 AND student_id = $2`
	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *StudentRepository) GetExperiences(ctx context.Context, studentID string) ([]model.Experience, error) {
	query := `
		SELECT id, student_id, company_name, position, description,
		       start_date::text, end_date::text, is_current, format
		FROM experiences
		WHERE student_id = $1
		ORDER BY start_date DESC NULLS LAST
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Experience
	for rows.Next() {
		e := model.Experience{}
		var description, startDate, endDate, format *string

		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.CompanyName, &e.Position,
			&description, &startDate, &endDate, &e.IsCurrent, &format,
		); err != nil {
			return nil, err
		}

		e.Description = derefStr(description)
		e.StartDate = derefStr(startDate)
		e.EndDate = derefStr(endDate)
		e.Format = derefStr(format)

		result = append(result, e)
	}

	return result, nil
}

func (r *StudentRepository) AddExperience(ctx context.Context, studentID string, e model.Experience) (*model.Experience, error) {
	query := `
		INSERT INTO experiences (student_id, company_name, position, description, start_date, end_date, is_current, format)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, student_id, company_name, position, description, start_date::text, end_date::text, is_current, format
	`

	result := &model.Experience{}
	var description, startDate, endDate, format *string

	err := r.db.QueryRow(ctx, query,
		studentID, e.CompanyName, e.Position,
		nilIfEmpty(e.Description), nilIfEmpty(e.StartDate), nilIfEmpty(e.EndDate),
		e.IsCurrent, nilIfEmpty(e.Format),
	).Scan(
		&result.ID, &result.StudentID, &result.CompanyName, &result.Position,
		&description, &startDate, &endDate, &result.IsCurrent, &format,
	)
	if err != nil {
		return nil, err
	}

	result.Description = derefStr(description)
	result.StartDate = derefStr(startDate)
	result.EndDate = derefStr(endDate)
	result.Format = derefStr(format)

	return result, nil
}

func (r *StudentRepository) UpdateExperience(ctx context.Context, id, studentID string, e model.Experience) (*model.Experience, error) {
	query := `
		UPDATE experiences
		SET company_name = $1, position = $2, description = $3,
		    start_date = $4, end_date = $5, is_current = $6, format = $7
		WHERE id = $8 AND student_id = $9
		RETURNING id, student_id, company_name, position, description, start_date::text, end_date::text, is_current, format
	`

	result := &model.Experience{}
	var description, startDate, endDate, format *string

	err := r.db.QueryRow(ctx, query,
		e.CompanyName, e.Position,
		nilIfEmpty(e.Description), nilIfEmpty(e.StartDate), nilIfEmpty(e.EndDate),
		e.IsCurrent, nilIfEmpty(e.Format),
		id, studentID,
	).Scan(
		&result.ID, &result.StudentID, &result.CompanyName, &result.Position,
		&description, &startDate, &endDate, &result.IsCurrent, &format,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	result.Description = derefStr(description)
	result.StartDate = derefStr(startDate)
	result.EndDate = derefStr(endDate)
	result.Format = derefStr(format)

	return result, nil
}

func (r *StudentRepository) DeleteExperience(ctx context.Context, id, studentID string) error {
	query := `DELETE FROM experiences WHERE id = $1 AND student_id = $2`
	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *StudentRepository) GetProjects(ctx context.Context, studentID string) ([]model.Project, error) {
	query := `
		SELECT id, student_id, title, description, url, techs, start_date::text, end_date::text
		FROM projects
		WHERE student_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Project
	for rows.Next() {
		p := model.Project{}
		var description, url, startDate, endDate *string

		if err := rows.Scan(
			&p.ID, &p.StudentID, &p.Title,
			&description, &url, &p.Techs,
			&startDate, &endDate,
		); err != nil {
			return nil, err
		}

		p.Description = derefStr(description)
		p.URL = derefStr(url)
		p.StartDate = derefStr(startDate)
		p.EndDate = derefStr(endDate)

		result = append(result, p)
	}

	return result, nil
}

func (r *StudentRepository) AddProject(ctx context.Context, studentID string, p model.Project) (*model.Project, error) {
	query := `
		INSERT INTO projects (student_id, title, description, url, techs, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, student_id, title, description, url, techs, start_date::text, end_date::text
	`

	result := &model.Project{}
	var description, url, startDate, endDate *string

	err := r.db.QueryRow(ctx, query,
		studentID, p.Title,
		nilIfEmpty(p.Description), nilIfEmpty(p.URL),
		p.Techs,
		nilIfEmpty(p.StartDate), nilIfEmpty(p.EndDate),
	).Scan(
		&result.ID, &result.StudentID, &result.Title,
		&description, &url, &result.Techs,
		&startDate, &endDate,
	)
	if err != nil {
		return nil, err
	}

	result.Description = derefStr(description)
	result.URL = derefStr(url)
	result.StartDate = derefStr(startDate)
	result.EndDate = derefStr(endDate)

	return result, nil
}

func (r *StudentRepository) UpdateProject(ctx context.Context, id, studentID string, p model.Project) (*model.Project, error) {
	query := `
		UPDATE projects
		SET title = $1, description = $2, url = $3, techs = $4, start_date = $5, end_date = $6
		WHERE id = $7 AND student_id = $8
		RETURNING id, student_id, title, description, url, techs, start_date::text, end_date::text
	`

	result := &model.Project{}
	var description, url, startDate, endDate *string

	err := r.db.QueryRow(ctx, query,
		p.Title, nilIfEmpty(p.Description), nilIfEmpty(p.URL),
		p.Techs, nilIfEmpty(p.StartDate), nilIfEmpty(p.EndDate),
		id, studentID,
	).Scan(
		&result.ID, &result.StudentID, &result.Title,
		&description, &url, &result.Techs,
		&startDate, &endDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	result.Description = derefStr(description)
	result.URL = derefStr(url)
	result.StartDate = derefStr(startDate)
	result.EndDate = derefStr(endDate)

	return result, nil
}

func (r *StudentRepository) DeleteProject(ctx context.Context, id, studentID string) error {
	query := `DELETE FROM projects WHERE id = $1 AND student_id = $2`
	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *StudentRepository) GetSkills(ctx context.Context, studentID string) ([]model.StudentSkill, error) {
	query := `SELECT id, student_id, skill, level FROM student_skills WHERE student_id = $1`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.StudentSkill
	for rows.Next() {
		s := model.StudentSkill{}
		if err := rows.Scan(&s.ID, &s.StudentID, &s.Skill, &s.Level); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *StudentRepository) GetLanguages(ctx context.Context, studentID string) ([]model.StudentLanguage, error) {
	query := `SELECT id, student_id, language, level FROM student_languages WHERE student_id = $1`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.StudentLanguage
	for rows.Next() {
		l := model.StudentLanguage{}
		if err := rows.Scan(&l.ID, &l.StudentID, &l.Language, &l.Level); err != nil {
			return nil, err
		}
		result = append(result, l)
	}

	return result, nil
}

func (r *StudentRepository) AddSkill(ctx context.Context, studentID string, skill, level string) (*model.StudentSkill, error) {
	query := `
        INSERT INTO student_skills (student_id, skill, level)
        VALUES ($1, $2, $3)
        RETURNING id, student_id, skill, level
    `
	s := &model.StudentSkill{}
	err := r.db.QueryRow(ctx, query, studentID, skill, level).Scan(&s.ID, &s.StudentID, &s.Skill, &s.Level)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *StudentRepository) DeleteSkill(ctx context.Context, id, studentID string) error {
	query := `DELETE FROM student_skills WHERE id = $1 AND student_id = $2`

	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *StudentRepository) AddLanguage(ctx context.Context, studentID string, language, level string) (*model.StudentLanguage, error) {
	query := `
        INSERT INTO student_languages (student_id, language, level)
        VALUES ($1, $2, $3)
        RETURNING id, student_id, language, level
    `

	l := &model.StudentLanguage{}
	err := r.db.QueryRow(ctx, query, studentID, language, level).Scan(&l.ID, &l.StudentID, &l.Language, &l.Level)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (r *StudentRepository) DeleteLanguage(ctx context.Context, id, studentID string) error {
	query := `DELETE FROM student_languages WHERE id = $1 AND student_id = $2`

	tag, err := r.db.Exec(ctx, query, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullIfZeroInt32(n int32) any {
	if n == 0 {
		return nil
	}
	return n
}

func nullIfZeroFloat(f float64) any {
	if f == 0 {
		return nil
	}
	return f
}

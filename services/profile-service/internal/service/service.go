package service

import (
	"context"
	"errors"

	"github.com/Kost0/internship-exchange/services/profile-service/internal/model"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/storage"
)

type ProfileService struct {
	students  *repository.StudentRepository
	companies *repository.CompanyRepository
	storage   *storage.MinioStorage
}

func NewProfileService(
	students *repository.StudentRepository,
	companies *repository.CompanyRepository,
	storage *storage.MinioStorage,
) *ProfileService {
	return &ProfileService{
		students:  students,
		companies: companies,
		storage:   storage,
	}
}

func (s *ProfileService) GetMyStudentProfile(ctx context.Context, userID string) (*model.Student, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.loadStudentRelations(ctx, student)
}

func (s *ProfileService) GetStudentProfile(ctx context.Context, id string) (*model.Student, error) {
	student, err := s.students.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}

		return nil, err
	}

	return s.loadStudentRelations(ctx, student)
}

func (s *ProfileService) UpdateStudentProfile(ctx context.Context, userID string, fields map[string]any) (*model.Student, error) {
	student, err := s.students.Update(ctx, userID, fields)
	if err != nil {
		return nil, err
	}

	return s.loadStudentRelations(ctx, student)
}

func (s *ProfileService) UploadAvatar(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	url, err := s.storage.UploadAvatar(ctx, userID, data, contentType)
	if err != nil {
		return "", err
	}

	return url, s.students.SetAvatarURL(ctx, userID, url)
}

func (s *ProfileService) UploadResume(ctx context.Context, userID string, data []byte) (string, error) {
	url, err := s.storage.UploadResume(ctx, userID, data)
	if err != nil {
		return "", err
	}

	return url, s.students.SetResumeURL(ctx, userID, url)
}

func (s *ProfileService) GetResumeURL(ctx context.Context, studentID string) (string, error) {
	student, err := s.students.GetByID(ctx, studentID)
	if err != nil {
		return "", err
	}

	if student.ResumeURL == "" {
		return "", errors.New("resume not uploaded")
	}

	return s.storage.GetResumePresignedURL(ctx, student.UserID)
}

func (s *ProfileService) AddEducation(ctx context.Context, userID string, e model.Education) (*model.Education, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.AddEducation(ctx, student.ID, e)
}

func (s *ProfileService) UpdateEducation(ctx context.Context, id, userID string, e model.Education) (*model.Education, error) {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.UpdateEducation(ctx, id, student.ID, e)
}

func (s *ProfileService) DeleteEducation(ctx context.Context, id, userID string) error {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.students.DeleteEducation(ctx, id, student.ID)
}

func (s *ProfileService) AddExperience(ctx context.Context, userID string, e model.Experience) (*model.Experience, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.AddExperience(ctx, student.ID, e)
}

func (s *ProfileService) UpdateExperience(ctx context.Context, id, userID string, e model.Experience) (*model.Experience, error) {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.UpdateExperience(ctx, id, student.ID, e)
}

func (s *ProfileService) DeleteExperience(ctx context.Context, id, userID string) error {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.students.DeleteExperience(ctx, id, student.ID)
}

func (s *ProfileService) AddProject(ctx context.Context, userID string, p model.Project) (*model.Project, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.AddProject(ctx, student.ID, p)
}

func (s *ProfileService) UpdateProject(ctx context.Context, id, userID string, p model.Project) (*model.Project, error) {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.UpdateProject(ctx, id, student.ID, p)
}

func (s *ProfileService) DeleteProject(ctx context.Context, id, userID string) error {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.students.DeleteProject(ctx, id, student.ID)
}

func (s *ProfileService) GetMyCompanyProfile(ctx context.Context, userID string) (*model.Company, error) {
	return s.companies.GetOrCreate(ctx, userID)
}

func (s *ProfileService) GetCompanyProfile(ctx context.Context, id string) (*model.Company, error) {
	return s.companies.GetByID(ctx, id)
}

func (s *ProfileService) UpdateCompanyProfile(ctx context.Context, userID string, c model.Company) (*model.Company, error) {
	return s.companies.Update(ctx, userID, c)
}

func (s *ProfileService) UploadLogo(ctx context.Context, userID string, data []byte, contentType string) (string, error) {
	url, err := s.storage.UploadLogo(ctx, userID, data, contentType)
	if err != nil {
		return "", err
	}

	return url, s.companies.SetLogoURL(ctx, userID, url)
}

func (s *ProfileService) loadStudentRelations(ctx context.Context, student *model.Student) (*model.Student, error) {
	var err error

	student.Educations, err = s.students.GetEducations(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	student.Experiences, err = s.students.GetExperiences(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	student.Projects, err = s.students.GetProjects(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	student.Skills, err = s.students.GetSkills(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	student.Languages, err = s.students.GetLanguages(ctx, student.ID)
	if err != nil {
		return nil, err
	}

	return student, nil
}

func (s *ProfileService) AddSkill(ctx context.Context, userID, skill, level string) (*model.StudentSkill, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return s.students.AddSkill(ctx, student.ID, skill, level)
}

func (s *ProfileService) DeleteSkill(ctx context.Context, id, userID string) error {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.students.DeleteSkill(ctx, id, student.ID)
}

func (s *ProfileService) AddLanguage(ctx context.Context, userID, language, level string) (*model.StudentLanguage, error) {
	student, err := s.students.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.students.AddLanguage(ctx, student.ID, language, level)
}

func (s *ProfileService) DeleteLanguage(ctx context.Context, id, userID string) error {
	student, err := s.students.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return s.students.DeleteLanguage(ctx, id, student.ID)
}

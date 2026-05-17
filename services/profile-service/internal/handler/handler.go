package handler

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	profilepb "github.com/Kost0/internship-exchange/proto/profile"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/model"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/repository"
	"github.com/Kost0/internship-exchange/services/profile-service/internal/service"
)

type ProfileHandler struct {
	profilepb.UnimplementedProfileServiceServer
	svc *service.ProfileService
}

func NewProfileHandler(svc *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) GetMyStudentProfile(ctx context.Context, req *profilepb.GetMyStudentProfileRequest) (*profilepb.StudentProfileResponse, error) {
	student, err := h.svc.GetMyStudentProfile(ctx, req.UserId)
	if err != nil {
		log.Printf("GetMyStudentProfile error userID=%s: %v", req.UserId, err)

		return nil, status.Error(codes.Internal, "internal error")
	}

	return studentToProto(student), nil
}

func (h *ProfileHandler) GetStudentProfile(ctx context.Context, req *profilepb.GetStudentProfileRequest) (*profilepb.StudentProfileResponse, error) {
	student, err := h.svc.GetStudentProfile(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "student not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return studentToProto(student), nil
}

func (h *ProfileHandler) UpdateStudentProfile(ctx context.Context, req *profilepb.UpdateStudentProfileRequest) (*profilepb.StudentProfileResponse, error) {
	fields := map[string]any{
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"phone":         req.Phone,
		"city":          req.City,
		"bio":           req.Bio,
		"github_url":    req.GithubUrl,
		"linkedin_url":  req.LinkedinUrl,
		"portfolio_url": req.PortfolioUrl,
	}

	student, err := h.svc.UpdateStudentProfile(ctx, req.UserId, fields)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return studentToProto(student), nil
}

func (h *ProfileHandler) UploadAvatar(ctx context.Context, req *profilepb.UploadAvatarRequest) (*profilepb.UploadAvatarResponse, error) {
	url, err := h.svc.UploadAvatar(ctx, req.UserId, req.Data, req.ContentType)
	if err != nil {
		return nil, status.Error(codes.Internal, "upload failed")
	}

	return &profilepb.UploadAvatarResponse{AvatarUrl: url}, nil
}

func (h *ProfileHandler) UploadResume(ctx context.Context, req *profilepb.UploadResumeRequest) (*profilepb.UploadResumeResponse, error) {
	url, err := h.svc.UploadResume(ctx, req.UserId, req.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, "upload failed")
	}

	return &profilepb.UploadResumeResponse{ResumeUrl: url}, nil
}

func (h *ProfileHandler) GetResumeURL(ctx context.Context, req *profilepb.GetResumeURLRequest) (*profilepb.GetResumeURLResponse, error) {
	url, err := h.svc.GetResumeURL(ctx, req.StudentId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "student not found")
		}

		return nil, status.Error(codes.NotFound, "resume not uploaded")
	}

	return &profilepb.GetResumeURLResponse{Url: url}, nil
}

func (h *ProfileHandler) AddEducation(ctx context.Context, req *profilepb.AddEducationRequest) (*profilepb.EducationResponse, error) {
	edu := model.Education{
		University:     req.University,
		Faculty:        req.Faculty,
		Specialization: req.Specialization,
		Degree:         req.Degree,
		StartYear:      req.StartYear,
		EndYear:        req.EndYear,
		GPA:            req.Gpa,
		IsCurrent:      req.IsCurrent,
	}

	result, err := h.svc.AddEducation(ctx, req.UserId, edu)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return educationToProto(result), nil
}

func (h *ProfileHandler) UpdateEducation(ctx context.Context, req *profilepb.UpdateEducationRequest) (*profilepb.EducationResponse, error) {
	edu := model.Education{
		University:     req.University,
		Faculty:        req.Faculty,
		Specialization: req.Specialization,
		Degree:         req.Degree,
		StartYear:      req.StartYear,
		EndYear:        req.EndYear,
		GPA:            req.Gpa,
		IsCurrent:      req.IsCurrent,
	}

	result, err := h.svc.UpdateEducation(ctx, req.Id, req.UserId, edu)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "education not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return educationToProto(result), nil
}

func (h *ProfileHandler) DeleteEducation(ctx context.Context, req *profilepb.DeleteEducationRequest) (*profilepb.DeleteResponse, error) {
	if err := h.svc.DeleteEducation(ctx, req.Id, req.UserId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "education not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &profilepb.DeleteResponse{Success: true}, nil
}

func (h *ProfileHandler) AddExperience(ctx context.Context, req *profilepb.AddExperienceRequest) (*profilepb.ExperienceResponse, error) {
	exp := model.Experience{
		CompanyName: req.CompanyName,
		Position:    req.Position,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsCurrent:   req.IsCurrent,
		Format:      req.Format,
	}

	result, err := h.svc.AddExperience(ctx, req.UserId, exp)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return experienceToProto(result), nil
}

func (h *ProfileHandler) UpdateExperience(ctx context.Context, req *profilepb.UpdateExperienceRequest) (*profilepb.ExperienceResponse, error) {
	exp := model.Experience{
		CompanyName: req.CompanyName,
		Position:    req.Position,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsCurrent:   req.IsCurrent,
		Format:      req.Format,
	}

	result, err := h.svc.UpdateExperience(ctx, req.Id, req.UserId, exp)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "experience not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return experienceToProto(result), nil
}

func (h *ProfileHandler) DeleteExperience(ctx context.Context, req *profilepb.DeleteExperienceRequest) (*profilepb.DeleteResponse, error) {
	if err := h.svc.DeleteExperience(ctx, req.Id, req.UserId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "experience not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &profilepb.DeleteResponse{Success: true}, nil
}

func (h *ProfileHandler) AddProject(ctx context.Context, req *profilepb.AddProjectRequest) (*profilepb.ProjectResponse, error) {
	techs := req.Techs
	if techs == nil {
		techs = []string{}
	}

	proj := model.Project{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.Url,
		Techs:       techs,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	result, err := h.svc.AddProject(ctx, req.UserId, proj)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return projectToProto(result), nil
}

func (h *ProfileHandler) UpdateProject(ctx context.Context, req *profilepb.UpdateProjectRequest) (*profilepb.ProjectResponse, error) {
	techs := req.Techs
	if techs == nil {
		techs = []string{}
	}

	proj := model.Project{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.Url,
		Techs:       techs,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	result, err := h.svc.UpdateProject(ctx, req.Id, req.UserId, proj)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return projectToProto(result), nil
}

func (h *ProfileHandler) DeleteProject(ctx context.Context, req *profilepb.DeleteProjectRequest) (*profilepb.DeleteResponse, error) {
	if err := h.svc.DeleteProject(ctx, req.Id, req.UserId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "project not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &profilepb.DeleteResponse{Success: true}, nil
}

func (h *ProfileHandler) GetMyCompanyProfile(ctx context.Context, req *profilepb.GetMyCompanyProfileRequest) (*profilepb.CompanyProfileResponse, error) {
	company, err := h.svc.GetMyCompanyProfile(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return companyToProto(company), nil
}

func (h *ProfileHandler) GetCompanyProfile(ctx context.Context, req *profilepb.GetCompanyProfileRequest) (*profilepb.CompanyProfileResponse, error) {
	company, err := h.svc.GetCompanyProfile(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "company not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return companyToProto(company), nil
}

func (h *ProfileHandler) UpdateCompanyProfile(ctx context.Context, req *profilepb.UpdateCompanyProfileRequest) (*profilepb.CompanyProfileResponse, error) {
	c := model.Company{
		Name:             req.Name,
		Tagline:          req.Tagline,
		Description:      req.Description,
		Industry:         req.Industry,
		Size:             req.Size,
		FoundedYear:      req.FoundedYear,
		Website:          req.Website,
		ContactEmail:     req.ContactEmail,
		City:             req.City,
		Country:          req.Country,
		IsRemoteFriendly: req.IsRemoteFriendly,
	}

	company, err := h.svc.UpdateCompanyProfile(ctx, req.UserId, c)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return companyToProto(company), nil
}

func (h *ProfileHandler) UploadLogo(ctx context.Context, req *profilepb.UploadLogoRequest) (*profilepb.UploadLogoResponse, error) {
	url, err := h.svc.UploadLogo(ctx, req.UserId, req.Data, req.ContentType)
	if err != nil {
		return nil, status.Error(codes.Internal, "upload failed")
	}

	return &profilepb.UploadLogoResponse{LogoUrl: url}, nil
}

func (h *ProfileHandler) AddSkill(ctx context.Context, req *profilepb.AddSkillRequest) (*profilepb.SkillResponse, error) {
	skill, err := h.svc.AddSkill(ctx, req.UserId, req.Skill, req.Level)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilepb.SkillResponse{
		Id: skill.ID, StudentId: skill.StudentID, Skill: skill.Skill, Level: skill.Level,
	}, nil
}

func (h *ProfileHandler) DeleteSkill(ctx context.Context, req *profilepb.DeleteSkillRequest) (*profilepb.DeleteResponse, error) {
	if err := h.svc.DeleteSkill(ctx, req.Id, req.UserId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "skill not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilepb.DeleteResponse{Success: true}, nil
}

func (h *ProfileHandler) AddLanguage(ctx context.Context, req *profilepb.AddLanguageRequest) (*profilepb.LanguageResponse, error) {
	lang, err := h.svc.AddLanguage(ctx, req.UserId, req.Language, req.Level)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilepb.LanguageResponse{
		Id: lang.ID, StudentId: lang.StudentID, Language: lang.Language, Level: lang.Level,
	}, nil
}

func (h *ProfileHandler) DeleteLanguage(ctx context.Context, req *profilepb.DeleteLanguageRequest) (*profilepb.DeleteResponse, error) {
	if err := h.svc.DeleteLanguage(ctx, req.Id, req.UserId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "language not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &profilepb.DeleteResponse{Success: true}, nil
}

func studentToProto(s *model.Student) *profilepb.StudentProfileResponse {
	resp := &profilepb.StudentProfileResponse{
		Id:           s.ID,
		UserId:       s.UserID,
		FirstName:    s.FirstName,
		LastName:     s.LastName,
		Phone:        s.Phone,
		City:         s.City,
		Bio:          s.Bio,
		AvatarUrl:    s.AvatarURL,
		ResumeUrl:    s.ResumeURL,
		GithubUrl:    s.GithubURL,
		LinkedinUrl:  s.LinkedinURL,
		PortfolioUrl: s.PortfolioURL,
	}

	for i := range s.Educations {
		resp.Educations = append(resp.Educations, educationToProto(&s.Educations[i]))
	}

	for i := range s.Experiences {
		resp.Experiences = append(resp.Experiences, experienceToProto(&s.Experiences[i]))
	}

	for i := range s.Projects {
		resp.Projects = append(resp.Projects, projectToProto(&s.Projects[i]))
	}

	for i := range s.Skills {
		resp.Skills = append(resp.Skills, &profilepb.SkillResponse{
			Id: s.Skills[i].ID, StudentId: s.Skills[i].StudentID,
			Skill: s.Skills[i].Skill, Level: s.Skills[i].Level,
		})
	}

	for i := range s.Languages {
		resp.Languages = append(resp.Languages, &profilepb.LanguageResponse{
			Id: s.Languages[i].ID, StudentId: s.Languages[i].StudentID,
			Language: s.Languages[i].Language, Level: s.Languages[i].Level,
		})
	}

	return resp
}

func educationToProto(e *model.Education) *profilepb.EducationResponse {
	return &profilepb.EducationResponse{
		Id: e.ID, StudentId: e.StudentID, University: e.University,
		Faculty: e.Faculty, Specialization: e.Specialization,
		Degree: e.Degree, StartYear: e.StartYear, EndYear: e.EndYear,
		Gpa: e.GPA, IsCurrent: e.IsCurrent,
	}
}

func experienceToProto(e *model.Experience) *profilepb.ExperienceResponse {
	return &profilepb.ExperienceResponse{
		Id: e.ID, StudentId: e.StudentID, CompanyName: e.CompanyName,
		Position: e.Position, Description: e.Description,
		StartDate: e.StartDate, EndDate: e.EndDate,
		IsCurrent: e.IsCurrent, Format: e.Format,
	}
}

func projectToProto(p *model.Project) *profilepb.ProjectResponse {
	return &profilepb.ProjectResponse{
		Id: p.ID, StudentId: p.StudentID, Title: p.Title,
		Description: p.Description, Url: p.URL, Techs: p.Techs,
		StartDate: p.StartDate, EndDate: p.EndDate,
	}
}

func companyToProto(c *model.Company) *profilepb.CompanyProfileResponse {
	return &profilepb.CompanyProfileResponse{
		Id: c.ID, UserId: c.UserID, Name: c.Name, Tagline: c.Tagline,
		Description: c.Description, Industry: c.Industry, Size: c.Size,
		FoundedYear: c.FoundedYear, Website: c.Website, ContactEmail: c.ContactEmail,
		City: c.City, Country: c.Country, IsRemoteFriendly: c.IsRemoteFriendly,
		LogoUrl: c.LogoURL,
	}
}

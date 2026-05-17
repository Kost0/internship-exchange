package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"

	profilepb "github.com/Kost0/internship-exchange/proto/profile"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/middleware"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
	"github.com/go-chi/chi/v5"
)

type ProfileHandler struct {
	client profilepb.ProfileServiceClient
}

func NewProfileHandler(conn *grpc.ClientConn) *ProfileHandler {
	return &ProfileHandler{client: profilepb.NewProfileServiceClient(conn)}
}

func (h *ProfileHandler) GetMyStudentProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.client.GetMyStudentProfile(r.Context(), &profilepb.GetMyStudentProfileRequest{
		UserId: userID,
	})
	if err != nil {
		log.Printf("GetMyStudentProfile error: %v", err)
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) GetStudentProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.client.GetStudentProfile(r.Context(), &profilepb.GetStudentProfileRequest{
		Id: id,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) UpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Phone        string `json:"phone"`
		City         string `json:"city"`
		Bio          string `json:"bio"`
		GithubURL    string `json:"githubUrl"`
		LinkedinURL  string `json:"linkedinUrl"`
		PortfolioURL string `json:"portfolioUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateStudentProfile(r.Context(), &profilepb.UpdateStudentProfileRequest{
		UserId:       userID,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Phone:        body.Phone,
		City:         body.City,
		Bio:          body.Bio,
		GithubUrl:    body.GithubURL,
		LinkedinUrl:  body.LinkedinURL,
		PortfolioUrl: body.PortfolioURL,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "file too large")

		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "avatar field required")

		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		proxy.WriteError(w, http.StatusInternalServerError, "failed to read file")

		return
	}

	res, err := h.client.UploadAvatar(r.Context(), &profilepb.UploadAvatarRequest{
		UserId:      userID,
		Data:        data,
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"avatarUrl": res.AvatarUrl})
}

func (h *ProfileHandler) UploadResume(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "file too large")

		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "resume field required")

		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		proxy.WriteError(w, http.StatusInternalServerError, "failed to read file")

		return
	}

	res, err := h.client.UploadResume(r.Context(), &profilepb.UploadResumeRequest{
		UserId: userID,
		Data:   data,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"resumeUrl": res.ResumeUrl})
}

func (h *ProfileHandler) GetResumeURL(w http.ResponseWriter, r *http.Request) {
	studentID := chi.URLParam(r, "id")

	res, err := h.client.GetResumeURL(r.Context(), &profilepb.GetResumeURLRequest{
		StudentId: studentID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"url": res.Url})
}

func (h *ProfileHandler) AddEducation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		University     string  `json:"university"`
		Faculty        string  `json:"faculty"`
		Specialization string  `json:"specialization"`
		Degree         string  `json:"degree"`
		StartYear      int32   `json:"startYear"`
		EndYear        int32   `json:"endYear"`
		GPA            float64 `json:"gpa"`
		IsCurrent      bool    `json:"isCurrent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.AddEducation(r.Context(), &profilepb.AddEducationRequest{
		UserId:         userID,
		University:     body.University,
		Faculty:        body.Faculty,
		Specialization: body.Specialization,
		Degree:         body.Degree,
		StartYear:      body.StartYear,
		EndYear:        body.EndYear,
		Gpa:            body.GPA,
		IsCurrent:      body.IsCurrent,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, res)
}

func (h *ProfileHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	var body struct {
		University     string  `json:"university"`
		Faculty        string  `json:"faculty"`
		Specialization string  `json:"specialization"`
		Degree         string  `json:"degree"`
		StartYear      int32   `json:"startYear"`
		EndYear        int32   `json:"endYear"`
		GPA            float64 `json:"gpa"`
		IsCurrent      bool    `json:"isCurrent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateEducation(r.Context(), &profilepb.UpdateEducationRequest{
		Id:             id,
		UserId:         userID,
		University:     body.University,
		Faculty:        body.Faculty,
		Specialization: body.Specialization,
		Degree:         body.Degree,
		StartYear:      body.StartYear,
		EndYear:        body.EndYear,
		Gpa:            body.GPA,
		IsCurrent:      body.IsCurrent,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	_, err := h.client.DeleteEducation(r.Context(), &profilepb.DeleteEducationRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfileHandler) AddExperience(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		CompanyName string `json:"companyName"`
		Position    string `json:"position"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		IsCurrent   bool   `json:"isCurrent"`
		Format      string `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.AddExperience(r.Context(), &profilepb.AddExperienceRequest{
		UserId:      userID,
		CompanyName: body.CompanyName,
		Position:    body.Position,
		Description: body.Description,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		IsCurrent:   body.IsCurrent,
		Format:      body.Format,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, res)
}

func (h *ProfileHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	var body struct {
		CompanyName string `json:"companyName"`
		Position    string `json:"position"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		IsCurrent   bool   `json:"isCurrent"`
		Format      string `json:"format"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateExperience(r.Context(), &profilepb.UpdateExperienceRequest{
		Id:          id,
		UserId:      userID,
		CompanyName: body.CompanyName,
		Position:    body.Position,
		Description: body.Description,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		IsCurrent:   body.IsCurrent,
		Format:      body.Format,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	_, err := h.client.DeleteExperience(r.Context(), &profilepb.DeleteExperienceRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfileHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Techs       []string `json:"techs"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.AddProject(r.Context(), &profilepb.AddProjectRequest{
		UserId:      userID,
		Title:       body.Title,
		Description: body.Description,
		Url:         body.URL,
		Techs:       body.Techs,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, res)
}

func (h *ProfileHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	var body struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Techs       []string `json:"techs"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateProject(r.Context(), &profilepb.UpdateProjectRequest{
		Id:          id,
		UserId:      userID,
		Title:       body.Title,
		Description: body.Description,
		Url:         body.URL,
		Techs:       body.Techs,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	id := chi.URLParam(r, "id")

	_, err := h.client.DeleteProject(r.Context(), &profilepb.DeleteProjectRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProfileHandler) GetMyCompanyProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	res, err := h.client.GetMyCompanyProfile(r.Context(), &profilepb.GetMyCompanyProfileRequest{
		UserId: userID,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.client.GetCompanyProfile(r.Context(), &profilepb.GetCompanyProfileRequest{
		Id: id,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var body struct {
		Name             string `json:"name"`
		Tagline          string `json:"tagline"`
		Description      string `json:"description"`
		Industry         string `json:"industry"`
		Size             string `json:"size"`
		FoundedYear      int32  `json:"foundedYear"`
		Website          string `json:"website"`
		ContactEmail     string `json:"contactEmail"`
		City             string `json:"city"`
		Country          string `json:"country"`
		IsRemoteFriendly bool   `json:"isRemoteFriendly"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.UpdateCompanyProfile(r.Context(), &profilepb.UpdateCompanyProfileRequest{
		UserId:           userID,
		Name:             body.Name,
		Tagline:          body.Tagline,
		Description:      body.Description,
		Industry:         body.Industry,
		Size:             body.Size,
		FoundedYear:      body.FoundedYear,
		Website:          body.Website,
		ContactEmail:     body.ContactEmail,
		City:             body.City,
		Country:          body.Country,
		IsRemoteFriendly: body.IsRemoteFriendly,
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, res)
}

func (h *ProfileHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "file too large")

		return
	}

	file, header, err := r.FormFile("logo")
	if err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "logo field required")

		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		proxy.WriteError(w, http.StatusInternalServerError, "failed to read file")

		return
	}

	res, err := h.client.UploadLogo(r.Context(), &profilepb.UploadLogoRequest{
		UserId:      userID,
		Data:        data,
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		proxy.WriteGRPCError(w, err)

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{"logoUrl": res.LogoUrl})
}

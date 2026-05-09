package handler

import (
	"net/http"

	"google.golang.org/grpc"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type ProfileHandler struct {
	conn *grpc.ClientConn
}

func NewProfileHandler(conn *grpc.ClientConn) *ProfileHandler {
	return &ProfileHandler{conn: conn}
}

func (h *ProfileHandler) GetMyStudentProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) GetStudentProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UploadResume(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) GetResumeURL(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) AddEducation(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UpdateEducation(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) DeleteEducation(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusNoContent, nil)
}
func (h *ProfileHandler) AddExperience(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusNoContent, nil)
}
func (h *ProfileHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusCreated, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusNoContent, nil)
}
func (h *ProfileHandler) GetMyCompanyProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) GetCompanyProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UpdateCompanyProfile(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}
func (h *ProfileHandler) UploadLogo(w http.ResponseWriter, r *http.Request) {
	proxy.WriteJSON(w, http.StatusOK, map[string]string{"message": "profile-service not yet connected"})
}

package handler

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/Kost0/internship-exchange/proto/auth"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type AuthHandler struct {
	client authpb.AuthServiceClient
}

func NewAuthHandler(conn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{client: authpb.NewAuthServiceClient(conn)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.Register(r.Context(), &authpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			proxy.WriteError(w, http.StatusConflict, "email already taken")

			return
		}
		proxy.WriteError(w, http.StatusInternalServerError, "internal error")

		return
	}

	proxy.WriteJSON(w, http.StatusCreated, map[string]any{
		"user": map[string]string{
			"id":    res.UserId,
			"email": res.Email,
			"role":  res.Role,
		},
		"tokens": map[string]string{
			"accessToken":  res.AccessToken,
			"refreshToken": res.RefreshToken,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.Login(r.Context(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			proxy.WriteError(w, http.StatusUnauthorized, "invalid credentials")

			return
		}

		proxy.WriteError(w, http.StatusInternalServerError, "internal error")

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]any{
		"user": map[string]string{
			"id":    res.UserId,
			"email": res.Email,
			"role":  res.Role,
		},
		"tokens": map[string]string{
			"accessToken":  res.AccessToken,
			"refreshToken": res.RefreshToken,
		},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		proxy.WriteError(w, http.StatusBadRequest, "invalid request body")

		return
	}

	res, err := h.client.RefreshToken(r.Context(), &authpb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			proxy.WriteError(w, http.StatusUnauthorized, "invalid refresh token")

			return
		}
		proxy.WriteError(w, http.StatusInternalServerError, "internal error")

		return
	}

	proxy.WriteJSON(w, http.StatusOK, map[string]string{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

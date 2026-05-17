package proxy

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

func WriteGRPCError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.NotFound:
		WriteError(w, http.StatusNotFound, status.Convert(err).Message())
	case codes.AlreadyExists:
		WriteError(w, http.StatusConflict, status.Convert(err).Message())
	case codes.InvalidArgument:
		WriteError(w, http.StatusBadRequest, status.Convert(err).Message())
	case codes.Unauthenticated:
		WriteError(w, http.StatusUnauthorized, status.Convert(err).Message())
	case codes.PermissionDenied:
		WriteError(w, http.StatusForbidden, status.Convert(err).Message())
	default:
		WriteError(w, http.StatusInternalServerError, "internal error")
	}
}

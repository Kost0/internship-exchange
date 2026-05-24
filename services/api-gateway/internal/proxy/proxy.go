package proxy

import (
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if msg, ok := v.(proto.Message); ok {
		b, err := protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   false,
		}.Marshal(msg)

		if err == nil {
			w.Write(b)

			return
		}
	}

	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, statusCode int, msg string) {
	WriteJSON(w, statusCode, map[string]string{"error": msg})
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

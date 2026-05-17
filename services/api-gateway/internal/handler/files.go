package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/proxy"
)

type FileProxyHandler struct {
	base          *url.URL
	reverseProxy  *httputil.ReverseProxy
	publicBuckets map[string]bool
}

func NewFileProxyHandler(minioAddr string) *FileProxyHandler {
	base := &url.URL{
		Scheme: "http",
		Host:   minioAddr,
	}

	return &FileProxyHandler{
		base:         base,
		reverseProxy: httputil.NewSingleHostReverseProxy(base),
		publicBuckets: map[string]bool{
			"avatars": true,
			"logos":   true,
		},
	}
}

func (h *FileProxyHandler) ServePublicFile(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	if !h.publicBuckets[bucket] {
		proxy.WriteError(w, http.StatusNotFound, "file not found")
		
		return
	}

	path := chi.URLParam(r, "path")
	if path == "" {
		proxy.WriteError(w, http.StatusNotFound, "file not found")

		return
	}

	r.URL.Scheme = h.base.Scheme
	r.URL.Host = h.base.Host
	r.Host = h.base.Host
	r.URL.Path = fmt.Sprintf("/%s/%s", bucket, path)

	h.reverseProxy.ServeHTTP(w, r)
}

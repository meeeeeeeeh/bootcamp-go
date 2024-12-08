package handler

import (
	"day_03/internal/middleware"
	"net/http"
)

func (h *Handler) InitRoutes() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", h.HandleUploadPageHtml)
	r.HandleFunc("/api/places", h.HandleUploadPageJson)
	r.HandleFunc("/api/recommend", middleware.AuthMiddleware(h.cfg, http.HandlerFunc(h.HandleRecommendation)))
	r.HandleFunc("/api/get_token", h.HandleGetToken)

	return r
}

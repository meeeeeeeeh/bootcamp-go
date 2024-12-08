package handler

import (
	"day_03/internal/config"
	"day_03/internal/model"
	"day_03/internal/service"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	pageSize             = 10
	recommendationAmount = 3
)

type Handler struct {
	svc *service.Service
	cfg *config.Config
}

func NewHandler(config *config.Config, s *service.Service) *Handler {
	return &Handler{
		svc: s,
		cfg: config,
	}
}

func (h *Handler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	var request model.RequestData

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.UserName == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	claims := &model.Claims{
		UserName: request.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JwtSecret))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

func (h *Handler) HandleRecommendation(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	latNumber, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lon := r.URL.Query().Get("lon")
	lonNumber, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recommend, err := h.svc.GetRecommendations(latNumber, lonNumber, recommendationAmount)
	if err != nil {
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	response := model.ResponseRecommendation{
		Name:   "Recommendation",
		Places: recommend,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) HandleUploadPageJson(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'page' value: '%d'", pageNumber),
		})
		return
	}

	offset := (pageNumber - 1) * pageSize
	places, total, err := h.svc.GetPlaces(pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to get places", http.StatusInternalServerError)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	if pageNumber > totalPages-1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'page' value: '%d'", pageNumber),
		})
		return
	}

	response := model.ResponsePlaces{
		Name:     "Places",
		Total:    total,
		Places:   places,
		PrevPage: pageNumber - 1,
		NextPage: pageNumber + 1,
		LastPage: totalPages,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) HandleUploadPageHtml(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", page), http.StatusBadRequest)
		return
	}

	offset := (pageNumber - 1) * pageSize
	places, total, err := h.svc.GetPlaces(pageSize, offset)
	if err != nil {
		http.Error(w, "Failed to get places", http.StatusInternalServerError)
		log.Fatalf("Failed to get places, err: %v", err)
	}

	totalPages := (total + pageSize - 1) / pageSize
	if pageNumber > totalPages-1 {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", page), http.StatusBadRequest)
		return
	}

	data := model.PageData{
		Places:       places,
		Total:        total,
		PageNumber:   pageNumber,
		TotalPages:   totalPages - 1,
		PreviousPage: pageNumber - 1,
		NextPage:     pageNumber + 1,
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalf("Error loading template: %s", err)
	}
	tmpl.Execute(w, data)
}

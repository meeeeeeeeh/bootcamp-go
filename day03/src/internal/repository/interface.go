package repository

import "day_03/internal/model"

type ESRepo interface {
	GetPlaces(limit int, offset int) ([]model.Place, int, error)
	GetRecommendations(lat, lon float64, amount int) ([]model.Place, error)
	GetTotalPlacesCount() (int, error)
}

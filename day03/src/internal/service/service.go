package service

import (
	"day_03/internal/model"
	"day_03/internal/repository"
	"fmt"
)

type Service struct {
	esRepo repository.ESRepo
}

func NewService(repo repository.ESRepo) *Service {
	return &Service{esRepo: repo}
}

func (s *Service) GetRecommendations(lat, lon float64, amount int) ([]model.Place, error) {
	recommend, err := s.esRepo.GetRecommendations(lat, lon, amount)
	if err != nil {
		return nil, fmt.Errorf("cannot get recommendations, err: %v", err)
	}
	return recommend, nil
}

func (s *Service) GetPlaces(limit, offset int) ([]model.Place, int, error) {
	places, amount, err := s.esRepo.GetPlaces(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get places, err: %v", err)
	}
	return places, amount, nil
}

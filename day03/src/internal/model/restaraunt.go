package model

import "github.com/dgrijalva/jwt-go"

type Place struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type PageData struct {
	Places       []Place
	Total        int
	PageNumber   int
	TotalPages   int
	PreviousPage int
	NextPage     int
}

type ResponsePlaces struct {
	Name     string  `json:"name"`
	Total    int     `json:"total"`
	Places   []Place `json:"places"`
	PrevPage int     `json:"prev_page"`
	NextPage int     `json:"next_page"`
	LastPage int     `json:"last_page"`
}

type ResponseRecommendation struct {
	Name   string  `json:"name"`
	Places []Place `json:"places"`
}

type Claims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

type RequestData struct {
	UserName string `json:"username"`
}

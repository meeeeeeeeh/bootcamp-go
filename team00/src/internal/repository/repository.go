package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"team-00/internal/config"
	"team-00/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) SaveAnomaly(anomaly model.Anomaly) error {
	res := r.db.Create(&anomaly)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Name, cfg.Postgres.Port)

	conn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot connect to db, err: %v", err)
	}
	return &Repository{db: conn}, nil
}

func (r *Repository) Close() {
	sqlDB, _ := r.db.DB()
	sqlDB.Close()
}

package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math"
	proto "team-00/generated"
	"team-00/internal/config"
	"team-00/internal/model"
)

const maxFrequencyCountForUpdate = 50

type Service struct {
	client         proto.RandomaliensServiceClient
	repo           DbRepo
	k              float64
	id             int64
	mean           float64
	std            float64
	frequencyCount int
	done           chan struct{}
}

func (s *Service) UpdateStats(frequency float64) {
	prevMean := s.mean
	s.mean = prevMean + (frequency-prevMean)/float64(s.frequencyCount)
	s.std = math.Sqrt(((float64(s.frequencyCount)-1)*s.std*s.std + (frequency-prevMean)*(frequency-s.mean)) / float64(s.frequencyCount))
	s.frequencyCount++
}

func (s *Service) IsAnomaly(frequency float64) bool {
	return math.Abs(frequency-s.mean) > s.k*s.mean
}

func (s *Service) ProcessAnomalies() {
	ctx := context.Background()
	stream, err := s.client.StreamFrequency(ctx, &proto.FrequencyIn{ClientId: s.id})
	if err != nil {
		log.Fatalf("cannot handle frequency stream, err: %v", err)
	}

	for {
		select {
		case <-s.done:
			log.Println("process anomalies finishes")
			return

		default:
			frequency, err := stream.Recv()
			if err != nil {
				log.Fatalf("cannot get frequency, err: %v", err)
			}
			s.frequencyCount++

			if s.frequencyCount <= maxFrequencyCountForUpdate {
				s.UpdateStats(frequency.Frequency)
				if s.frequencyCount%10 == 0 {
					log.Printf("frequency count: %d, mean: %f, std: %f, frequency: %f, k: %f\n", s.frequencyCount, s.mean, s.std, frequency.Frequency, s.k)
				}
			} else {
				if s.IsAnomaly(frequency.Frequency) {
					anomaly := model.Anomaly{
						Dt:        frequency.Dt.AsTime(),
						SessionId: frequency.SessionId,
						Frequency: frequency.Frequency,
					}
					err := s.repo.SaveAnomaly(anomaly)
					if err != nil {
						log.Fatalf("cannot save anomaly, err: %v", err)
					}
				}
			}
		}
	}
}

func (s *Service) StopProcessingAnomalies() {
	s.done <- struct{}{}
}

func NewClient(cfg *config.Config, db DbRepo, k float64, id int64) (*Service, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.Service.Host, cfg.Service.Port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to randomalies service, err: %v", err)
	}
	client := proto.NewRandomaliensServiceClient(conn)

	s := &Service{
		client: client,
		repo:   db,
		k:      k,
		id:     id,
		done:   make(chan struct{}),
	}
	go s.ProcessAnomalies()
	return s, nil
}

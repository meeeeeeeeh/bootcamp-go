package sevice

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	proto "team-00/generated"
	"time"
)

type Service struct {
	proto.UnimplementedRandomaliensServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) StreamFrequency(in *proto.FrequencyIn, stream proto.RandomaliensService_StreamFrequencyServer) error {
	sessionId := uuid.New().String()
	mean := -10 + rand.Float64()*(10-(-10))
	std := 0.3 + rand.Float64()*(1.5-0.3)

	log.Printf("client_id: %d,session_id: %s, mean: %f, std: %f\n", in.ClientId, sessionId, mean, std)

	for {
		frequency := rand.NormFloat64() + mean
		currentTime := time.Now().UTC()
		timestamp := timestamppb.New(currentTime)

		response := &proto.FrequencyOut{
			SessionId: sessionId,
			Frequency: frequency,
			Dt:        timestamp,
		}

		if err := stream.SendMsg(response); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
}

package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "team-00/generated"
	"team-00/internal/config"
	"team-00/internal/sevice"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)

	}

	grpcServer := grpc.NewServer()
	s := sevice.NewService()
	proto.RegisterRandomaliensServiceServer(grpcServer, s)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Service.Host, cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
}

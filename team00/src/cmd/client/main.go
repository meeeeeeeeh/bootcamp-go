package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"team-00/internal/client"
	"team-00/internal/config"
	"team-00/internal/repository"
)

func main() {
	var clientId int64 = 1
	k := flag.Float64("k", 3.0, "STD anomaly coefficient")
	flag.Parse()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)

	}

	dbRepo, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("cannot initialize db repository: %v", err)
	}
	defer dbRepo.Close()

	c, err := client.NewClient(cfg, dbRepo, *k, clientId)
	if err != nil {
		log.Fatalf("cannot initialize client service: %v", err)
	}
	defer c.StopProcessingAnomalies()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("the program was finished")
}

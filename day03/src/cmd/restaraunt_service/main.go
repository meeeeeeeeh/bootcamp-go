package main

import (
	"day_03/internal/config"
	handler "day_03/internal/handlers"
	"day_03/internal/repository"
	"day_03/internal/server"
	"day_03/internal/service"
	"log"
	"os"
	"os/signal"
)

//запуск проекта в докере
// docker-compose up --build

//docker-compose down --volumes
// флаг --volumes удаляет все тома из контейнера - например данные которые были загружены в эластиксерч

// env файл для локального запуска - ELASTICSEARCH_ADDRESS="https://localhost:9200"
// env файл для запуска  в докере - ELASTICSEARCH_ADDRESS="http://elasticsearch:9200"

//запустить elasticsearch
//cd /usr/local/elasticsearch
//./bin/elasticsearch
// ./usr/local/elasticsearch/bin/elasticsearch

//проверка соединения
//curl -k -u "endeharh:123456" -X GET "https://localhost:9200/"
//удалить индекс
//curl -k -u "endeharh:123456" -X DELETE "https://localhost:9200/index_name"
//вывести все индексы
//curl -k -u "endeharh:123456" -X GET "https://localhost:9200/_cat/indices?v"

//вывести данные по индексу
//curl -X GET "https://localhost:9200/places/_search?pretty" -u "endeharh:123456" -k

//curl -X GET "localhost:9200/"
//password yl=-1H-8GW4rutOy0=FN

//user endeharh
//password 123456

//увеличить лимит значений которй возвращает эластисерч
//curl -X PUT "https://localhost:9200/places/_settings" \
//-u "endeharh:123456" \
//-H "Content-Type: application/json" \
//-k \
//-d '{
//"index": {
//"max_result_window": 50000
//}
//}'

//проверить лимит значений
// curl -k -u "endeharh:123456" -X GET "https://localhost:9200/places/_settings?pretty" -H 'Content-Type: application/json'

//чтобы эластисерч потреблял меньше памяти нужно установить лимит

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot initialize config, err: %v", err)
	}

	elasticSearch, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatalf("cannot create new repository, err: %v", err)
	}

	restaurantService := service.NewService(elasticSearch)
	handle := handler.NewHandler(cfg, restaurantService)
	svr := server.New(cfg.AddressSvr, handle.InitRoutes())

	go func() {
		err := svr.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server start...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	err = svr.Stop()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server stop...")
}

package main

import (
	"bytes"
	"context"
	"day_03/internal/config"
	"day_03/internal/model"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot initialize config, err: %v", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			conf.AddressES,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("error creating the client: %v", err)
	}

	//wait for the connection to es
	for {
		res, err := es.Info()
		if err == nil {
			defer res.Body.Close()
			// Проверяем статус ответа
			if res.IsError() {
				log.Println("Elasticsearch is not ready, retrying...")
			} else {
				log.Println("Elasticsearch is ready!")
				break
			}
		} else {
			log.Println("Waiting for Elasticsearch...")
		}
		time.Sleep(2 * time.Second)
	}

	//file path supposed to run project from route
	err = initIndex(es, "places", "internal/config/elasticsearch/mapping.json")
	if err != nil {
		log.Fatalf("cannot create new index, %v", err)
	}

	err = uploadPlaceData(es, "internal/config/elasticsearch/data.csv")
	if err != nil {
		log.Fatalf("cannot upload the index data to elasticsearch, %v", err)
	}
}

func uploadPlaceData(es *elasticsearch.Client, dataFilePath string) error {
	file, err := os.Open(dataFilePath)
	if err != nil {
		return fmt.Errorf("error opening the file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %v", err)
	}

	var buf bytes.Buffer
	for _, record := range records {
		id := record[0]
		name := record[1]
		address := record[2]
		phone := record[3]
		lon, _ := strconv.ParseFloat(record[4], 64)
		lat, _ := strconv.ParseFloat(record[5], 64)

		place := model.Place{
			Name:    name,
			Address: address,
			Phone:   phone,
			Location: model.GeoPoint{
				Lat: lat,
				Lon: lon,
			},
		}

		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "places", "_id" : "%s" } }%s`, id, "\n"))
		data, err := json.Marshal(place)
		if err != nil {
			return fmt.Errorf("error marshaling document: %v", err)
		}

		buf.Grow(len(meta) + len(data) + 1)
		buf.Write(meta)
		buf.Write(data)
		buf.WriteByte('\n')
	}

	bulkReq := esapi.BulkRequest{
		Body: bytes.NewReader(buf.Bytes()),
	}

	bulkRes, err := bulkReq.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("error performing bulk request: %v", err)
	}
	defer bulkRes.Body.Close()

	if bulkRes.IsError() {
		return fmt.Errorf("bulk request error: %s", bulkRes.String())
	}

	log.Println("Data indexed successfully.")
	return nil
}

func initIndex(es *elasticsearch.Client, indexName, mappingFilePath string) error {
	byteValue, err := os.ReadFile(mappingFilePath)
	if err != nil {
		return fmt.Errorf("error oppening the file: %v", err)
	}

	//increases response limit til 20k instead of default 10k
	mapping := string(byteValue)
	requestBody := fmt.Sprintf(`{
		"settings": {
			"index": {
				"max_result_window": 20000
			}
		},
		"mappings": %s
	}`, mapping)

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader([]byte(requestBody)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return fmt.Errorf("error creating the index: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating the index: %s", res.String())
	}

	log.Println("Index created.")
	return nil
}

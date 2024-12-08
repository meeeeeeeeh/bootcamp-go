package repository

import (
	"context"
	"day_03/internal/config"
	"day_03/internal/model"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"strconv"
	"strings"
)

type repository struct {
	client *elasticsearch.Client
}

func NewRepository(config *config.Config) (*repository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.AddressES,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %v", err)
	}

	return &repository{
		client: es,
	}, nil
}

func (r *repository) GetRecommendations(lat, lon float64, amount int) ([]model.Place, error) {
	var places []model.Place

	query := fmt.Sprintf(`{
	"size": %d,
	"query": {
		"match_all": {}
	},
	"sort": [
		{
			"_geo_distance": {
				"location": {
					"lat": %f,
					"lon": %f
				},
				"order": "asc",
				"unit": "km",
				"mode": "min",
				"distance_type": "arc",
				"ignore_unmapped": true
			}
		}
	]
}`, amount, lat, lon)

	req := esapi.SearchRequest{
		Index: []string{"places"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		var place model.Place
		source := hit.(map[string]interface{})["_source"]
		data, _ := json.Marshal(source)
		json.Unmarshal(data, &place)

		// Elasticsearch automatically assigns a unique identifier (_id) to each document when it is indexed.
		id, _ := hit.(map[string]interface{})["_id"].(string)
		place.ID, _ = strconv.Atoi(id)
		places = append(places, place)
	}

	return places, nil
}

// returns a list of items, a total number of hits and  an error in case of one
func (r *repository) GetPlaces(limit int, offset int) ([]model.Place, int, error) {
	var places []model.Place

	query := fmt.Sprintf(`{
        "from": %d,
        "size": %d,
        "query": {
            "match_all": {}
        }
    }`, offset, limit)

	req := esapi.SearchRequest{
		Index: []string{"places"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, err
	}

	hits := response["hits"].(map[string]interface{})["hits"].([]interface{})
	//total := int(response["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	for _, hit := range hits {
		var place model.Place
		source := hit.(map[string]interface{})["_source"]
		data, _ := json.Marshal(source)
		json.Unmarshal(data, &place)

		// Elasticsearch automatically assigns a unique identifier (_id) to each document when it is indexed.
		id, _ := hit.(map[string]interface{})["_id"].(string)
		place.ID, _ = strconv.Atoi(id)

		places = append(places, place)
	}

	total, err := r.GetTotalPlacesCount()
	if err != nil {
		return nil, 0, fmt.Errorf("cannot get total index amount, err: %v", err)
	}

	return places, total, nil
}

func (r *repository) GetTotalPlacesCount() (int, error) {
	req := esapi.CountRequest{
		Index: []string{"places"},
		Body:  strings.NewReader(`{"query": {"match_all": {}}}`),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, err
	}

	total := int(response["count"].(float64))

	return total, nil
}

package dbreader

import (
	"day01/internal/model"
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

type readerXML struct {
	file string
}

type readerJSON struct {
	file string
}

type DBReader interface {
	Read() (*model.Storage, error)
}

func NewReader(filename string) (DBReader, error) {
	if strings.Contains(filename, ".xml") {
		return &readerXML{file: filename}, nil
	}
	if strings.Contains(filename, ".json") {
		return &readerJSON{file: filename}, nil
	}
	err := errors.New("file extension is not supported")
	return nil, err
}

func (r readerXML) Read() (*model.Storage, error) {
	var db model.Storage
	file, err := os.ReadFile(r.file)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(file, &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (r readerJSON) Read() (*model.Storage, error) {
	var db model.Storage
	file, err := os.ReadFile(r.file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

package main

import (
	"day01/internal/dbreader"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
)

func printOppositeFile(file string) error {
	r, err := dbreader.NewReader(file)
	if err != nil {
		return err
	}

	data, err := r.Read()
	if err != nil {
		return err
	}

	var b []byte
	if strings.Contains(file, ".xml") {
		b, err = json.MarshalIndent(data, "", "    ")
		if err != nil {
			return err
		}
	} else if strings.Contains(file, ".json") {
		b, err = xml.MarshalIndent(data, "", "    ")
		if err != nil {
			return err
		}
	} else {
		err = errors.New("file extention is not supported")
		return err
	}

	fmt.Println(string(b))
	return nil
}

func main() {
	file := flag.String("f", "", "file name")
	flag.Parse()

	err := printOppositeFile(*file)
	if err != nil {
		log.Fatalln(err)
	}

}

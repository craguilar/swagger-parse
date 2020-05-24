package main

import (
	"encoding/csv"

	"log"
	"os"
)

func main() {
	fileName := "swagger.yaml"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	log.Printf("%s", fileName)
	yaml, err := ReadYaml(fileName)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	definitions, err := ParseGetPaths(yaml)
	if err != nil {
		log.Fatalf("Parse get paths: %v", err)
	}
	rows, err := ParseDefinitions(definitions, yaml)
	if err != nil {
		log.Fatalf("Parse definitions: %v", err)
	}
	err = saveToCsv(rows)
	if err != nil {
		log.Fatalf("Save to CSV: %v", err)
	}
}

func saveToCsv(rows [][]string) error {
	file, err := os.Create("data-matrix.csv")
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range rows {
		err := writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}

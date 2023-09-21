package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type configuration struct {
	Interval int
	Websites []website
}

func readConfiguration() configuration {
	var cFlag = flag.String("c", "configuration.json", "path to configuration file")
	flag.Parse()

	jsonFile, err := os.Open(*cFlag)
	if err != nil {
		panic(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var configuration configuration
	json.Unmarshal(byteValue, &configuration)

	defer jsonFile.Close()

	return configuration
}

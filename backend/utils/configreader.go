package utils

import (
	"io"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func getIntAttributeFromConfigFile(attribute string) int {
	file, err := os.Open("../config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]int
	yaml.Unmarshal([]byte(byteValue), &result)

	return result[attribute]
}

func getStringAttributeFromConfigFile(attribute string) string {
	file, err := os.Open("../config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]string
	yaml.Unmarshal([]byte(byteValue), &result)

	return result[attribute]
}

func GetMapDimension() int {
	return getIntAttributeFromConfigFile("MapDimension")
}

func GetExAgentVisionRange() int {
	return getIntAttributeFromConfigFile("ExAgentVisionRange")
}

func GetExName() string {
	return getStringAttributeFromConfigFile("ExName")
}

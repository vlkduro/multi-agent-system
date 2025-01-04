package utils

import (
	"io"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func getIntAttributeFromConfigFile(attribute string) int {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]int
	yaml.Unmarshal([]byte(byteValue), &result)

	return result[attribute]
}

func getFloat64AttributeFromConfigFile(attribute string) float64 {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]float64
	yaml.Unmarshal([]byte(byteValue), &result)

	return result[attribute]
}

func getStringAttributeFromConfigFile(attribute string) string {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var result map[string]string
	yaml.Unmarshal([]byte(byteValue), &result)

	return result[attribute]
}

func GetNumberBees() int {
	return getIntAttributeFromConfigFile("NumberBees")
}

func GetMaxNectar() int {
	return getIntAttributeFromConfigFile("MaxNectar")
}

func GetNumberFlowers() int {
	return getIntAttributeFromConfigFile("NumberFlowers")
}

func GetNumberFlowerPatches() int {
	return getIntAttributeFromConfigFile("NumberFlowerPatches")
}

func GetNumberObjects() int {
	return getIntAttributeFromConfigFile("NumberFlowers") + 1
}

func GetMapDimension() int {
	return getIntAttributeFromConfigFile("MapDimension")
}

func GetBeeAgentVisionRange() float64 {
	return getFloat64AttributeFromConfigFile("BeeAgentVisionRange")
}

func GetExName() string {
	return getStringAttributeFromConfigFile("ExName")
}

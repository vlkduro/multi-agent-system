package utils

import (
	"io"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func getIntAttributeFromConfigFile(attribute string) int {
	file, err := os.Open("/Users/quentin.v/ai30_valakou_martins_chartier_bidaux/config.yaml")
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
	file, err := os.Open("/Users/quentin.v/ai30_valakou_martins_chartier_bidaux/config.yaml")
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
	file, err := os.Open("/Users/quentin.v/ai30_valakou_martins_chartier_bidaux/config.yaml")
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

func GetNumberHornets() int {
	return getIntAttributeFromConfigFile("NumberHornets")
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

func GetProducedNectarPerTurn() int {
	return getIntAttributeFromConfigFile("ProducedNectarPerTurn")
}

func GetMaxNectarHeld() int {
	return getIntAttributeFromConfigFile("MaxNectarHeld")
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

func GetHornetAgentVisionRange() float64 {
	return getFloat64AttributeFromConfigFile("HornetAgentVisionRange")
}

func GetExName() string {
	return getStringAttributeFromConfigFile("ExName")
}

func GetBeeCreationCost() int {
	return getIntAttributeFromConfigFile("BeeCreationCost")
}

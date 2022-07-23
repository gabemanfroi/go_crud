package validations

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
)

var necessaryDtos = []string{"create", "read", "update"}

func getMissingDtos(existing, necessary []string) []string {
	missingMap := CreateMissingMap(existing)

	var missingDtos []string
	for _, element := range necessary {
		if !missingMap[element] {
			missingDtos = append(missingDtos, element)
		}
	}
	return missingDtos
}

func CheckIfModelHasAllDtos(modelName string, directory string) {

	var missingDtos []string

	existingDtosFiles, _ := ioutil.ReadDir(directory + "/domain/DTO/" + modelName)

	if len(existingDtosFiles) == 0 {
		log.Fatalf("No DTOs found for [%s]", modelName)
	}

	existingDtos := getExistingDtos(existingDtosFiles)

	if len(existingDtosFiles) < 3 {
		missingDtos = getMissingDtos(existingDtos, necessaryDtos)
	}

	if len(missingDtos) > 0 {
		log.Fatalf("Missing DTOs for [%s] %s", modelName, missingDtos)
	}
}

func getExistingDtos(existingDtosFiles []fs.FileInfo) []string {
	var existingDtos []string
	for _, f := range existingDtosFiles {
		existingDtos = append(existingDtos, strings.Split(f.Name(), ".")[0])
	}
	return existingDtos
}

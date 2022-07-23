package validations

import (
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
)

var necessaryValidators = []string{"create", "update"}

func getMissingValidators(existing, necessary []string) []string {
	missingMap := CreateMissingMap(existing)

	var missingValidators []string
	for _, element := range necessary {
		if !missingMap[element] {
			missingValidators = append(missingValidators, element)
		}
	}
	return missingValidators
}

func CheckIfModelHasAllValidators(modelName string, directory string) {

	var missingValidators []string

	existingValidatorsFiles, _ := ioutil.ReadDir(directory + "/infra/validators/" + modelName)

	if len(existingValidatorsFiles) == 0 {
		log.Fatalf("No Validators found for [%s]", modelName)
	}

	existingValidators := getExistingValidators(existingValidatorsFiles)

	if len(existingValidatorsFiles) < 3 {
		missingValidators = getMissingValidators(existingValidators, necessaryValidators)
	}

	if len(missingValidators) > 0 {
		log.Fatalf("Missing DTOs for [%s] %s", modelName, missingValidators)
	}
}

func getExistingValidators(existingValidatorsFiles []fs.FileInfo) []string {
	var existingValidators []string
	for _, f := range existingValidatorsFiles {
		existingValidators = append(existingValidators, strings.Split(f.Name(), ".")[0])
	}
	return existingValidators
}

package validations

import (
	"io/ioutil"
	"log"
)

func CheckIfModelExists(modelName string, directory string) {

	_, err := ioutil.ReadFile(directory + "/domain/models/" + modelName + ".go")

	if err != nil {
		log.Fatalf("model [%s] does not exist", modelName)
	}
}

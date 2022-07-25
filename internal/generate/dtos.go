package generate

import (
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/prompt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"io/ioutil"
)

func getReadDTOProperties(modelName string, modelProperties []ModelProperty) []DTOProperty {
	return getDTOProperties(modelName, modelProperties, "read")
}

func getUpdateDTOProperties(modelName string, modelProperties []ModelProperty) []DTOProperty {
	return getDTOProperties(modelName, modelProperties, "update")
}

func getCreateDTOProperties(modelName string, modelProperties []ModelProperty) []DTOProperty {
	return getDTOProperties(modelName, modelProperties, "create")
}

func getDTOProperties(modelName string, modelProperties []ModelProperty, name string) []DTOProperty {
	file, err := getExistingModelFile(modelName, name)
	if err != nil {
		return initDTOCreationPrompt(modelName, modelProperties, name)
	}
	return getDTOPropertiesFromFile(file)
}

func getExistingModelFile(modelName string, name string) ([]byte, error) {
	dtosDir := fmt.Sprintf("%s/domain/DTO/%s", utils.GetWorkingDirectory(), modelName)
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.go", dtosDir, name))
	return file, err
}

func initDTOCreationPrompt(modelName string, modelProperties []ModelProperty, name string) []DTOProperty {
	if prompt.GetBoolean(prompt.GetDtoNotFoundPromptContent(modelName, name)) {
		return createDTO(modelName, modelProperties, name)
	}
	return nil
}

func createDTO(modelName string, modelProperties []ModelProperty, name string) []DTOProperty {

	var properties []DTOProperty
	var propertyNames []string
	propertyMap := make(map[string]ModelProperty)
	chooseAnotherProperty := true

	for _, p := range modelProperties {
		propertyNames = append(propertyNames, p.Name)
		propertyMap[p.Name] = p
	}

	for chooseAnotherProperty {

		selectedProperty := prompt.GetSelection(prompt.GetSelectPropertyPromptContent(modelName), propertyNames)
		createdProperty := createDtoProperty(propertyMap, selectedProperty)

		if isUpdateDTO(name) {
			createdProperty.DataType = "*" + createdProperty.DataType
		}

		properties = append(properties, createdProperty)
		propertyNames = utils.RemoveStringFromArray(propertyNames, createdProperty.Name)

		if noPropertiesRemaining(propertyNames) {
			return properties
		}

		chooseAnotherProperty = prompt.GetBoolean(addAnotherPropertyPromptContent)
	}

	return properties
}

func isUpdateDTO(name string) bool {
	return name == "update"
}

func noPropertiesRemaining(propertyNames []string) bool {
	return len(propertyNames) == 0
}

func createDtoProperty(propertyMap map[string]ModelProperty, selectedProperty string) DTOProperty {
	createdProperty := DTOProperty{
		Property: Property{
			Name:     propertyMap[selectedProperty].Name,
			DataType: propertyMap[selectedProperty].DataType,
		},
		JsonString: getJsonString(propertyMap[selectedProperty].Name),
	}
	return createdProperty
}

func getJsonString(propertyName string) string {
	return fmt.Sprintf("`json:\"%s\"`", utils.PascalToCamelCase(propertyName))
}

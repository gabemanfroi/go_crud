package generate

import (
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/prompt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"io/ioutil"
)

type DTOProperty struct {
	Name       string
	DataType   string
	JsonString string
}

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
	dtosDir := fmt.Sprintf("%s/domain/DTO/%s", utils.GetWorkingDirectory(), modelName)

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.go", dtosDir, name))

	if err != nil {
		dtosNotFoundPromptContent := prompt.Content{
			Label: fmt.Sprintf("No %s DTO was found for the model %s would you like to create it?", name, modelName),
		}
		if prompt.GetBoolean(dtosNotFoundPromptContent) {
			return createDTO(modelName, modelProperties, name)
		}
	}

	return getDTOPropertiesFromFile(content)
}

func createDTO(modelName string, modelProperties []ModelProperty, name string) []DTOProperty {
	dtosPropertiesPromptContent := prompt.Content{
		Label: "Choose which property of " + modelName + " you want to add to the DTO ",
	}

	var properties []DTOProperty
	var propertyNames []string
	propertyMap := make(map[string]ModelProperty)
	chooseAnotherProperty := true

	for _, p := range modelProperties {
		propertyNames = append(propertyNames, p.Name)
		propertyMap[p.Name] = p
	}

	for chooseAnotherProperty {

		selectedProperty := prompt.GetSelection(dtosPropertiesPromptContent, propertyNames)
		createdProperty := DTOProperty{
			Name:       propertyMap[selectedProperty].Name,
			DataType:   propertyMap[selectedProperty].DataType,
			JsonString: fmt.Sprintf("`json:\"%s\"`", utils.PascalToCamelCase(propertyMap[selectedProperty].Name)),
		}

		if name == "update" {
			createdProperty.DataType = "*" + createdProperty.DataType
		}

		properties = append(properties, createdProperty)
		propertyNames = utils.RemoveStringFromArray(propertyNames, createdProperty.Name)

		if len(propertyNames) == 0 {
			return properties
		}

		chooseAnotherProperty = prompt.GetBoolean(addAnotherPropertyPromptContent)
	}

	return properties
}

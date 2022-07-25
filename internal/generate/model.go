package generate

import (
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/prompt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"io/ioutil"
	"strings"
)

func getModelProperties(modelName string) []ModelProperty {
	content, err := getModelFile(modelName, utils.GetWorkingDirectory())
	if err != nil {
		return initModelCreationPrompt(modelName)
	}
	return getExistingModelProperties(content)
}

func initModelCreationPrompt(modelName string) []ModelProperty {
	modelNotFoundPromptContent := prompt.Content{
		Label: "Could not find a model named " + modelName + ", would you like to create one?",
	}

	if prompt.GetBoolean(modelNotFoundPromptContent) {
		return CreateModel(modelName)
	}
	return nil
}

func getModelFile(modelName string, directory string) ([]byte, error) {
	content, err := ioutil.ReadFile(directory + "/domain/models/" + modelName + ".go")
	if err != nil {
		return nil, err
	}
	return content, err
}

func CreateModel(modelName string) []ModelProperty {
	fmt.Printf("Creating model %s\n", modelName)
	return getModelWithProperties()
}

func getModelWithProperties() []ModelProperty {
	addOneMoreProperty := true

	var properties []ModelProperty

	for addOneMoreProperty {
		createdProperty := createModelProperty()
		addCreatedPropertyToArray(&properties, createdProperty)
		addOneMoreProperty = prompt.GetBoolean(addAnotherPropertyPromptContent)
	}
	return properties
}

func propertyIsRelational(property ModelProperty) bool {
	return utils.ArrayContainsString(getExistingModels(), property.DataType)
}

func addCreatedPropertyToArray(properties *[]ModelProperty, createdProperty ModelProperty) {
	*properties = append(*properties, createdProperty)
	handleRelationalPropertyCreation(properties, createdProperty)
}

func handleRelationalPropertyCreation(properties *[]ModelProperty, createdProperty ModelProperty) {
	if propertyIsRelational(createdProperty) {
		*properties = append(*properties, ModelProperty{
			Property: Property{
				Name:     createdProperty.DataType + "Id",
				DataType: "uint",
			},
			GormString: "`gorm:\"not null\"`",
		})
	}
}

func createModelProperty() ModelProperty {
	var property ModelProperty
	property.Name = prompt.GetInput(propertyNamePromptContent)

	property.DataType = prompt.GetSelection(propertyDataTypePromptContent, append(dataTypesOptions, getExistingModels()...))

	handlePropertyNullability(&property)
	return property
}

func handlePropertyNullability(property *ModelProperty) {
	propertyNullablePromptContent := prompt.Content{
		Label: fmt.Sprintf("Is %s nullable?", property.Name),
	}
	property.Nullable = prompt.GetBoolean(propertyNullablePromptContent)
	if !property.Nullable {
		property.GormString = "`gorm:\"not null\"`"
		getMinimums(property)
	}
}

func getMinimums(property *ModelProperty) {
	if property.DataType == "string" {
		getMinimumLength(property)
	}
	if property.DataType == "uint" {
		getMinimumValue(property)
	}
}

func getMinimumValue(property *ModelProperty) {
	property.MinimumValue = prompt.GetNumber(prompt.GetMinimumValuePromptContent(property.Name))
}

func getMinimumLength(property *ModelProperty) {
	property.MinimumLength = prompt.GetNumber(prompt.GetMinimumLengthValuePromptContent(property.Name))
}

func getExistingModels() []string {
	var models []string
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/domain/models", utils.GetWorkingDirectory()))

	utils.Check(err, "could not read from dir")

	for _, f := range files {
		modelName := strings.Split(utils.SnakeToPascalCase(f.Name()), ".")[0]
		if modelName != "README" {
			models = append(models, modelName)
		}

	}

	return models
}

func getExistingModelProperties(file []byte) []ModelProperty {
	return getModelPropertiesFromFile(file)
}

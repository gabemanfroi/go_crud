package generate

import (
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/prompt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"io/ioutil"
)

var (
	propertyNamePromptContent = prompt.Content{
		ErrorMsg: "Please provide a property name.",
		Label:    "Enter a name for the property (Pascal Case, please)",
	}
	propertyDataTypePromptContent = prompt.Content{
		Label: "Select the property Data type",
	}
	addAnotherPropertyPromptContent = prompt.Content{
		Label: "Property added, would you like to add more properties?",
	}
	dataTypesOptions = []string{"string", "uint", "bool"}
)

func getModelProperties(modelName string) []ModelProperty {

	err := GetModel(modelName, utils.GetWorkingDirectory())

	if err != nil {

		modelNotFoundPromptContent := prompt.Content{
			Label: "Could not find a model named " + modelName + ", would you like to create one?",
		}

		if prompt.GetBoolean(modelNotFoundPromptContent) {
			return CreateModel(modelName)
		}
	}

	return nil
}

func GetModel(modelName string, directory string) error {
	_, err := ioutil.ReadFile(directory + "/domain/models/" + modelName + ".go")
	return err
}

type ModelProperty struct {
	Name          string
	DataType      string
	Nullable      bool
	MinimumValue  uint
	MinimumLength uint
	GormString    string
}

func CreateModel(modelName string) []ModelProperty {
	fmt.Printf("Creating model %s\n", modelName)
	return PopulateModelWithProperties()
}

func PopulateModelWithProperties() []ModelProperty {
	addOneMoreProperty := true

	var properties []ModelProperty

	for addOneMoreProperty {
		properties = append(properties, CreateProperty())
		addOneMoreProperty = prompt.GetBoolean(addAnotherPropertyPromptContent)
	}
	return properties
}

func CreateProperty() ModelProperty {
	var property ModelProperty
	property.Name = prompt.GetInput(propertyNamePromptContent)
	property.DataType = prompt.GetSelection(propertyDataTypePromptContent, dataTypesOptions)
	propertyNullablePromptContent := prompt.Content{
		Label: fmt.Sprintf("Is %s nullable?", property.Name),
	}
	property.Nullable = prompt.GetBoolean(propertyNullablePromptContent)
	if !property.Nullable {
		property.GormString = ("`gorm:\"not null\"`")
		getMinimums(property)
	}
	return property
}

func getMinimums(property ModelProperty) {
	if property.DataType == "string" {
		getMinimumLength(property)
	}
	if property.DataType == "uint" {
		getMinimumValue(property)
	}
}

func getMinimumValue(property ModelProperty) {
	propertyMinimumLengthPromptContent := prompt.Content{
		ErrorMsg: "Enter a valid number.",
		Label:    fmt.Sprintf("Please enter the minimum value for %s", property.Name),
	}
	property.MinimumValue = prompt.GetNumber(propertyMinimumLengthPromptContent)
}

func getMinimumLength(property ModelProperty) {
	propertyMinimumLengthPromptContent := prompt.Content{
		ErrorMsg: "Enter a valid number.",
		Label:    fmt.Sprintf("Please enter the minimum length for %s", property.Name),
	}
	property.MinimumLength = prompt.GetNumber(propertyMinimumLengthPromptContent)
}

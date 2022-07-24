package generate

import (
	"fmt"
	"github.com/gabemanfroi/go_crud/internal/prompt"
	"github.com/gabemanfroi/go_crud/internal/utils"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	addAnotherRulePromptContent = prompt.Content{
		Label: "Rule added, would you like to add more rules?",
	}
	rules          = []string{"required", "omitempty", "min", "max", "gte", "lte", "lt", "gt"}
	numericOptions = []string{"min", "max", "gte", "lte", "lt", "gt"}
	propertyMap    = make(map[string]ModelProperty)
	propertyNames  []string
)

type ValidatorProperty struct {
	Name           string
	DataType       string
	ValidateString string
}

func getCreateValidatorProperties(modelName string, modelProperties []ModelProperty) []ValidatorProperty {
	return getValidatorProperties(modelName, modelProperties, "create")
}

func getUpdateValidatorProperties(modelName string, modelProperties []ModelProperty) []ValidatorProperty {
	return getValidatorProperties(modelName, modelProperties, "update")
}

func getValidatorProperties(modelName string, modelProperties []ModelProperty, name string) []ValidatorProperty {
	validatorsDir := fmt.Sprintf("%s/infra/validators/%s", utils.GetWorkingDirectory(), modelName)

	utils.CreateDirectoryIfNotExists(validatorsDir)

	content, err := ioutil.ReadFile(fmt.Sprintf("validatorsDir/%s.go", name))

	if err != nil {
		return initValidatorCreationPrompt(modelName, modelProperties, name)
	}

	return getValidatorPropertiesFromFile(content)
}

func initValidatorCreationPrompt(modelName string, modelProperties []ModelProperty, name string) []ValidatorProperty {
	validatorsNotFoundPromptContent := prompt.Content{
		Label: fmt.Sprintf("No %s Validator was found for the model %s would you like to create it?", name, modelName),
	}
	if prompt.GetBoolean(validatorsNotFoundPromptContent) {
		return createValidator(modelName, modelProperties, name)
	}
	return nil
}

func createValidator(modelName string, modelProperties []ModelProperty, name string) []ValidatorProperty {

	var properties []ValidatorProperty

	chooseAnotherProperty := true

	initializePropertiesHandlers(modelProperties)

	for chooseAnotherProperty {
		availableRules := []string{"required", "omitempty", "min", "max", "gte", "lte", "lt", "gt"}
		createdProperty := createProperty(modelName, name, availableRules)
		properties = append(properties, createdProperty)

		if len(propertyNames) == 0 {
			return properties
		}
		chooseAnotherProperty = prompt.GetBoolean(addAnotherPropertyPromptContent)
	}

	return properties
}

func createProperty(modelName string, name string, availableRules []string) ValidatorProperty {

	var selectedRules []string
	chooseAnotherRule := true

	if name == "update" {
		availableRules = utils.RemoveStringFromArray(availableRules, "omitempty")
		selectedRules = append(selectedRules, "omitempty")
	}

	createdProperty := initializeProperty(modelName)

	for chooseAnotherRule {
		validateOptionsPromptContent := prompt.Content{
			Label: "Choose the rule you want to apply for " + createdProperty.Name,
		}
		selectedRule := prompt.GetSelection(validateOptionsPromptContent, rules)
		availableRules = utils.RemoveStringFromArray(availableRules, selectedRule)
		for _, o := range numericOptions {
			if selectedRule == o {
				selectedRule = getSelectedRuleWithValue(selectedRule)
				break
			}
		}
		selectedRules = append(selectedRules, selectedRule)
		availableRules = utils.RemoveStringFromArray(availableRules, selectedRule)
		chooseAnotherRule = prompt.GetBoolean(addAnotherRulePromptContent)
	}

	createdProperty.ValidateString = getFormatedString(selectedRules)
	propertyNames = utils.RemoveStringFromArray(propertyNames, createdProperty.Name)
	return createdProperty
}

func getFormatedString(selectedRules []string) string {
	return fmt.Sprintf("`validate:\"%s\"`", strings.Join(selectedRules, ","))
}

func initializePropertiesHandlers(modelProperties []ModelProperty) {
	if len(propertyNames) > 0 {
		propertyNames = nil
	}

	for _, p := range modelProperties {
		propertyNames = append(propertyNames, p.Name)
		propertyMap[p.Name] = p
	}
}

func initializeProperty(modelName string) ValidatorProperty {
	validatorPropertiesPromptContent := prompt.Content{
		Label: "Choose which property of " + modelName + " you want to add to the Validator ",
	}
	selectedProperty := prompt.GetSelection(validatorPropertiesPromptContent, propertyNames)
	createdProperty := ValidatorProperty{
		Name:     propertyMap[selectedProperty].Name,
		DataType: propertyMap[selectedProperty].DataType,
	}
	return createdProperty
}

func getSelectedRuleWithValue(selectedRule string) string {

	selectedRuleValuePromptContent := prompt.Content{
		ErrorMsg: "Enter a valid number.",
		Label:    fmt.Sprintf("Please enter the minimum value for %s", selectedRule),
	}
	selectedRule = selectedRule + "=" + strconv.Itoa(int(prompt.GetNumber(selectedRuleValuePromptContent)))
	return selectedRule
}

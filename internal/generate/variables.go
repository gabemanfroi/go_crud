package generate

import "github.com/gabemanfroi/go_crud/internal/prompt"

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
	dataTypesOptions            = []string{"string", "uint", "bool"}
	addAnotherRulePromptContent = prompt.Content{
		Label: "Rule added, would you like to add more rules?",
	}
	numericOptions = []string{"min", "max", "gte", "lte", "lt", "gt"}
	propertyMap    = make(map[string]ModelProperty)
	propertyNames  []string
)

package generate

import (
	"regexp"
	"strings"
)

func getDTOPropertiesFromFile(file []byte) []DTOProperty {
	var properties []DTOProperty
	content := normalizeFileContent(file)
	for _, p := range extractPropertiesFromLines(getLines(content)) {
		properties = append(properties, DTOProperty{Property: p})
	}
	return properties
}

func getValidatorPropertiesFromFile(file []byte) []ValidatorProperty {
	var properties []ValidatorProperty
	content := normalizeFileContent(file)
	for _, p := range extractPropertiesFromLines(getLines(content)) {
		properties = append(properties, ValidatorProperty{Property: p})
	}
	return properties
}

func getModelPropertiesFromFile(file []byte) []ModelProperty {
	var properties []ModelProperty
	content := normalizeFileContent(file)
	for _, p := range extractPropertiesFromLines(getLines(content)) {
		properties = append(properties, ModelProperty{Property: p})
	}
	return properties
}

func extractPropertiesFromLines(lines []string) []Property {
	var properties []Property
	for _, line := range lines {
		properties = populatePropertiesArray(line, properties)
	}
	return properties
}

func populatePropertiesArray(line string, properties []Property) []Property {
	property := extractPropertyFromLine(line)
	if property.Name != "" {
		properties = append(properties, Property{Name: property.Name, DataType: property.DataType})
	}
	return properties
}

func extractPropertyFromLine(line string) Property {
	re := regexp.MustCompile("\\t\\S+\\s*\\S+")

	match := re.FindStringSubmatch(line)
	return getPropertyFromRegex(match)
}

func getPropertyFromRegex(match []string) Property {
	if match != nil {
		lineWithoutTab := strings.Replace(match[0], "\t", "", -1)
		if propertyIsNotNeeded(lineWithoutTab) {
			return createPropertyFromLine(lineWithoutTab)
		}
		return Property{}
	}
	return Property{}
}

func propertyIsNotNeeded(lineWithoutTab string) bool {
	return lineWithoutTab != "gorm.Model"
}

func createPropertyFromLine(lineWithoutTab string) Property {
	propertyLine := strings.Fields(lineWithoutTab)
	return Property{Name: propertyLine[0], DataType: propertyLine[1]}
}

func getLines(newContent string) []string {
	lines := strings.Split(newContent, "\n")
	return lines
}

func normalizeFileContent(content []byte) string {
	newContent := strings.Split(string(content), "{")[1]
	newContent = strings.Split(newContent, "}")[0]
	return newContent
}

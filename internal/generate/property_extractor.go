package generate

import (
	"regexp"
	"strings"
)

func getDTOPropertiesFromFile(content []byte) []DTOProperty {
	newContent := removeBracketsFromContent(content)
	return extractDTOPropertiesFromLines(getLines(newContent))
}

func extractDTOPropertiesFromLines(lines []string) []DTOProperty {
	var properties []DTOProperty
	for _, line := range lines {
		property := extractPropertyFromLine(line)
		if property != "" {
			properties = append(properties, DTOProperty{Name: property})
		}

	}
	return properties
}

func getValidatorPropertiesFromFile(content []byte) []ValidatorProperty {
	newContent := removeBracketsFromContent(content)
	return extractValidatorPropertiesFromLines(getLines(newContent))
}

func extractValidatorPropertiesFromLines(lines []string) []ValidatorProperty {
	var properties []ValidatorProperty
	for _, line := range lines {
		property := extractPropertyFromLine(line)
		if property != "" {
			properties = append(properties, ValidatorProperty{Name: property})
		}

	}
	return properties
}

func extractPropertyFromLine(line string) string {
	re := regexp.MustCompile("\t[^\\s]+")

	match := re.FindStringSubmatch(line)
	if match != nil {
		return strings.Replace(match[0], "\t", "", -1)
	}
	return ""
}

func getLines(newContent string) []string {
	lines := strings.Split(newContent, "\n")
	return lines
}

func removeBracketsFromContent(content []byte) string {
	newContent := strings.Split(string(content), "{")[1]
	newContent = strings.Split(newContent, "}")[0]
	return newContent
}

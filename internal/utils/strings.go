package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func PascalToCamelCase(s string) string {
	return strings.ToLower(s[0:1]) + s[1:len(s)]
}

func RemoveStringFromArray(array []string, s string) []string {
	for i, v := range array {
		if v == s {
			return append(array[:i], array[i+1:]...)
		}
	}
	return array
}

func SnakeToCamelCase(modelName string) string {
	var camelCasedName string
	if !strings.Contains(modelName, "_") {
		return modelName
	}
	words := strings.Split(modelName, "_")
	for i, w := range words {
		if i == 0 {
			camelCasedName += w
		} else {
			camelCasedName += cases.Title(language.BrazilianPortuguese, cases.NoLower).String(w)
		}
	}
	return camelCasedName
}

func SnakeToPascalCase(modelName string) string {
	var pascalCasedName string
	if !strings.Contains(modelName, "_") {
		return cases.Title(language.BrazilianPortuguese, cases.NoLower).String(modelName)
	}
	words := strings.Split(modelName, "_")
	for _, w := range words {
		pascalCasedName += cases.Title(language.BrazilianPortuguese, cases.NoLower).String(w)
	}
	return pascalCasedName
}

package generate

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func GetCamelCasedModelName(modelName string) string {
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

func GetPascalCasedModelName(modelName string) string {
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

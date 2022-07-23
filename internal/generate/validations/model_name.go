package validations

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func CamelizeModelName(modelName string) string {
	var camelized string
	if !strings.Contains(modelName, "_") {
		return modelName
	}
	words := strings.Split(modelName, "_")
	for i, w := range words {
		if i == 0 {
			camelized += w
		} else {
			camelized += cases.Title(language.BrazilianPortuguese, cases.NoLower).String(w)
		}
	}
	return camelized
}

func PascalizeModelName(modelName string) string {
	var pascalized string
	if !strings.Contains(modelName, "_") {
		return modelName
	}
	words := strings.Split(modelName, "_")
	for _, w := range words {
		pascalized += cases.Title(language.BrazilianPortuguese, cases.NoLower).String(w)
	}
	return pascalized
}

package utils

import "strings"

func PascalToCamelCase(s string) string {
	return strings.ToLower(s[0:1]) + s[1:len(s)]
}

func RemoveStringFromArray(a []string, s string) []string {
	for i, v := range a {
		if v == s {
			return append(a[:i], a[i+1:]...)
		}
	}
	return a
}

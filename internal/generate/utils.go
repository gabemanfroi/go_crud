package generate

func CreateMissingMap(firstArray []string) map[string]bool {
	missingMap := make(map[string]bool)
	for _, element := range firstArray {
		missingMap[element] = true
	}
	return missingMap
}

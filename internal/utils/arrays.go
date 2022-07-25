package utils

func ArrayContainsString(a []string, s string) bool {
	for _, w := range a {
		if w == s {
			return true
		}
	}
	return false
}

package analyze

import (
	"strings"
)

func IsValidEpStr(i1 string) bool {
	return epstrExp.MatchString(i1)
}

func CleanSymbols(i1 string) string {
	i1 = strings.ReplaceAll(i1, "[", "")
	i1 = strings.ReplaceAll(i1, "]", "")
	i1 = strings.ReplaceAll(i1, "(", "")
	i1 = strings.ReplaceAll(i1, ")", "")

	return i1
}

func CleanExtension(i1 string) string {
	temp := strings.ToLower(i1)
	var x int
	for _, v1 := range extensions {
		if x = strings.Index(temp, v1); x > 0 {
			i1 = strings.TrimSpace(i1[:x])
			break
		}
	}

	return i1
}

func HasExtension(i1 string) bool {
	temp := strings.ToLower(i1)
	for _, v1 := range extensions {
		if strings.Contains(temp, v1) {
			return true
		}
	}

	return false
}

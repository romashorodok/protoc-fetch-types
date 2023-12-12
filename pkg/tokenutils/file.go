package tokenutils

import "strings"

func countRunes(str string, char rune) int {
	count := 0
	for _, c := range str {
		if c == char {
			count++
		}
	}
	return count
}

func IsRoot(filePath string) bool {
	path := strings.Split(filePath, "/")
	return strings.HasSuffix(path[0], ".proto")
}

func BackwardPath(filePath string) string {
	// TODO: Only unix paths
	return AppendBackwards(filePath, countRunes(filePath, '/'))
}

func GetBackwardCount(filePath string) int {
	return countRunes(filePath, '/')
}

func AppendBackwards(filePath string, backwards int) string {
	var result string
	for i := 0; i < backwards; i++ {
		result += "../"
	}
	result += filePath
	return result
}

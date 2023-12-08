package requestfunc

import (
	"fmt"
	"regexp"
	"strings"
)

func extractParams(target string) ([]string, error) {
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	matches := re.FindAllStringSubmatch(target, -1)

	if matches == nil {
		return nil, fmt.Errorf("no matches found")
	}

	var values []string
	for _, match := range matches {
		if len(match) == 2 {
			values = append(values, match[1])
		}
	}

	return values, nil
}

func replaceGoogleTypeParams(pattern string, values []string, sep string) (string, []string) {
	result := pattern
	// Google has different tpyeof patterns
	// https://github.com/googleapis/googleapis/blob/2aa3b1d5a90d05e0606d11787de475b0df068d49/google/api/http.proto#L72
	for idx, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) == 2 {
			paramName := parts[0]
			values[idx] = paramName

			// Handle the format {key=value}
			result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", value), "${"+sep+paramName+"}")
		} else {
			// Handle the format {key}
			result = strings.ReplaceAll(result, fmt.Sprintf("{%s}", value), "${"+sep+value+"}")
		}
	}

	return result, values
}

// returns - (uriPath string, paramsInUrl []string, err error)
func parsePattern(pattern string) (string, []string, error) {
	params, err := extractParams(pattern)
	if err != nil {
		return "", nil, err
	}
	uriPath, params := replaceGoogleTypeParams(pattern, params, "params.")
	return uriPath, params, err
}

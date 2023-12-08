package requestfunc

import (
	"errors"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

var InconsistentHttpAnnotationError = errors.New("inconsistent http annotation")

func getHttpField(httpRule *annotations.HttpRule) (pattern, method string) {
	switch {
	case httpRule.GetGet() != "":
		return httpRule.GetGet(), "GET"
	case httpRule.GetPut() != "":
		return httpRule.GetPut(), "PUT"
	case httpRule.GetPost() != "":
		return httpRule.GetPost(), "POST"
	case httpRule.GetDelete() != "":
		return httpRule.GetDelete(), "DELETE"
	case httpRule.GetPatch() != "":
		return httpRule.GetPatch(), "PATCH"
	default:
		return "", ""
	}
}

// return - (pattern string, requestMethod string, err error)
func googleHttpAnnotation(method *proxy.MethodProxy) (string, string, error) {
	httpRule_Any := proto.GetExtension(method.GetOptions(), annotations.E_Http)

	if httpRule, ok := httpRule_Any.(*annotations.HttpRule); ok {
		pattern, method := getHttpField(httpRule)
		if pattern == "" || method == "" {
			return "", "", InconsistentHttpAnnotationError
		}
		return pattern, method, nil
	}

	return "", "", InconsistentHttpAnnotationError
}

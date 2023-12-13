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
    // https://github.com/googleapis/googleapis/blob/2aa3b1d5a90d05e0606d11787de475b0df068d49/google/api/annotations.proto#L30C23-L30C23
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

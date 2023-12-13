package proxy

import (
	"strings"
	"unicode"

	"google.golang.org/protobuf/types/descriptorpb"
)

const unixPathSeparator = "/"

func uppercase(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := unicode.ToUpper(rune(s[0]))
	return string(firstChar) + s[1:]
}

func GetNamespaceTokens(file *descriptorpb.FileDescriptorProto) []string {
	namespace := strings.Split(file.GetName(), unixPathSeparator)
	namespace = namespace[:len(namespace)-1]
	if len(namespace) == 0 || namespace[len(namespace)-1] != file.GetPackage() {
		namespace = append(namespace, file.GetPackage())
	}
	return namespace
}

// Determine name for aliases in imports and types
func ImportAliasFromFilePath(file *descriptorpb.FileDescriptorProto) string {
	namespace := strings.Split(file.GetName(), unixPathSeparator)
	namespace = append(namespace[:len(namespace)-1],
		strings.TrimSuffix(namespace[len(namespace)-1], ".proto"),
	)

	var result string
	for idx, token := range namespace {
		if idx == 0 {
			result += token
			continue
		}
		result += uppercase(token)
	}
	return result
}

func GetNamespace(file *descriptorpb.FileDescriptorProto) string {
	// NOTE: nested namespace format as data.product.metadata. Decided use alias and single namespace from file.Package().
	// TODO: `fetch_types.ts` may contain full names to import it as single module.
	// return strings.Join(GetNamespaceTokens(file), ".")

	return file.GetPackage()
}

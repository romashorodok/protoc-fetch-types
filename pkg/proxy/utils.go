package proxy

import (
	"strings"
	"unicode"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	unixPathSeparator = "/"
	packageSeparator  = "."
)

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

// NOTE: It's actually support nested namespaces but like that `namespace room.models { ... }`
// To do that it should be in one folder and has package name like that `room.models`.
// EXAMPLE:
// import { room as modelsRoom_models } from "./models/room_models";
// ...
//
//	↓
//
// export type RoomCreateResponse = { room: modelsRoom_models.models.Room; };
func PackageNamespacePrefix(file *descriptorpb.FileDescriptorProto) string {
	var result string
	namespacePackage := strings.Split(file.GetPackage(), packageSeparator)
	if len(namespacePackage) > 1 {
		namespacePackage = namespacePackage[1:]
		for _, packageNamespace := range namespacePackage {
			result += packageNamespace + packageSeparator
		}
	}
    return strings.TrimRight(result, packageSeparator)
}

func PackageNamespaceSuffix(file *descriptorpb.FileDescriptorProto) string {
    return strings.SplitN(file.GetPackage(), packageSeparator, 2)[0]
}

func GetNamespace(file *descriptorpb.FileDescriptorProto) string {
	// NOTE: nested namespace format as data.product.metadata. Decided use alias and single namespace from file.Package().
	// TODO: `fetch_types.ts` may contain full names to import it as single module.
	// return strings.Join(GetNamespaceTokens(file), ".")

	return file.GetPackage()
}

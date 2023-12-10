package proxy

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

const unixPathSeparator = "/"

func GetNamespaceTokens(file *descriptorpb.FileDescriptorProto) []string {
	namespace := strings.Split(file.GetName(), unixPathSeparator)
	namespace = namespace[:len(namespace)-1]
	if len(namespace) == 0 || namespace[len(namespace)-1] != file.GetPackage() {
		namespace = append(namespace, file.GetPackage())
	}
	return namespace
}

func GetNamespace(file *descriptorpb.FileDescriptorProto) string {
	return strings.Join(GetNamespaceTokens(file), ".")
}

package tokenutils

import (
	"fmt"
	"unicode"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Uppercase(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := unicode.ToUpper(rune(s[0]))
	return string(firstChar) + s[1:]
}

func Lowercase(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := unicode.ToLower(rune(s[0]))
	return string(firstChar) + s[1:]
}

func TsType(t *descriptorpb.FieldDescriptorProto_Type) string {
	switch *t {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "number"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "boolean"
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		return "Object"
	default:
		return "any"
	}
}

func TsArray(field *descriptorpb.FieldDescriptorProto) bool {
	return field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

func TypeAliasName(message *proxy.MessageProxy) string {
	return fmt.Sprintf("%sPackage%s", message.GetName(), Uppercase(message.GetPackageName()))
}

func TypeAliasParamName(message *proxy.MessageProxy) string {
	return fmt.Sprintf("%s_%s", Lowercase(message.GetName()), message.GetPackageName())
}

package typealias

import (
	"embed"
	"fmt"
	"io"
	"log"
	"unicode"

	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

func uppercase(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := unicode.ToUpper(rune(s[0]))
	return string(firstChar) + s[1:]
}

func lowercase(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := unicode.ToLower(rune(s[0]))
	return string(firstChar) + s[1:]
}

func tsType(t *descriptorpb.FieldDescriptorProto_Type) string {
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

func tsArray(field *descriptorpb.FieldDescriptorProto) bool {
	return field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

func typeAliasName(message *proxy.MessageProxy) string {
	return fmt.Sprintf("%sPackage%s", message.GetName(), uppercase(message.GetPackageName()))
}

type TypeAlias struct {
	ref  *proxy.MessageProxy
	tmpl *templatebuilder.TemplateBuilder
}

type template_TypeAliasParams struct {
	Type  string
	Name  string
	Array bool
}

type template_TypeAlias struct {
	Name       string
	Parameters []*template_TypeAliasParams
}

func (s *TypeAlias) tmplStruct() template_TypeAlias {
	var params []*template_TypeAliasParams

	for _, message := range s.ref.GetFieldsMessages() {
		var isArray bool
		if field, err := message.GetOriginField(); err == nil {
			isArray = tsArray(field)
		}

		params = append(params,
			&template_TypeAliasParams{
				Name:  fmt.Sprintf("%s_%s", lowercase(message.GetName()), message.GetPackageName()),
				Type:  typeAliasName(message),
				Array: isArray,
			})
	}

	for _, field := range s.ref.GetPrimitiveFields() {
		params = append(params,
			&template_TypeAliasParams{
				Name:  field.GetName(),
				Type:  tsType(field.Type),
				Array: tsArray(field),
			},
		)
	}

	return template_TypeAlias{
		Name:       typeAliasName(s.ref),
		Parameters: params,
	}
}

func (s *TypeAlias) WriteInto(in io.Writer) error {
	log.Println(s.tmplStruct())
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

const TYPE_ALIAS_TEMPLATE_FILE = "templates/type_alias.tmpl"

func New(storage embed.FS, ref *proxy.MessageProxy) *TypeAlias {
	templateFile, err := storage.ReadFile(TYPE_ALIAS_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", TYPE_ALIAS_TEMPLATE_FILE, storage)
	}

	log.Println(templateFile)

	tmpl := templatebuilder.New(templateFile, TYPE_ALIAS_TEMPLATE_FILE)

	return &TypeAlias{
		ref:  ref,
		tmpl: tmpl,
	}
}

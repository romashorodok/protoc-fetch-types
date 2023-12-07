package typealias

import (
	"embed"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
	"google.golang.org/protobuf/types/descriptorpb"
)

type typeAliasParameter struct {
	Type string
	Name string

	ref *descriptorpb.FieldDescriptorProto
}

type TypeAliasObject struct {
	TypeName string

	ref  *descriptorpb.DescriptorProto
	tmpl *templatebuilder.TemplateBuilder
}

func (s *TypeAliasObject) WriteInto(in io.Writer) error {
	return nil
}

type NewParams struct {
	TypeAliasObject
}

func New(storage embed.FS, ref *descriptorpb.DescriptorProto, params *NewParams) *TypeAliasObject {
	templateFile, err := storage.ReadFile(resources.TYPE_ALIAS_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", resources.TYPE_ALIAS_TEMPLATE_FILE, storage)
	}

	tmpl := templatebuilder.New(templateFile, resources.TYPE_ALIAS_TEMPLATE_FILE)

	return &TypeAliasObject{
		TypeName: params.TypeName,

		tmpl: tmpl,
		ref:  ref,
	}
}

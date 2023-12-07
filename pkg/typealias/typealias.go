package typealias

import (
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

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
		name := message.GetName()
		packageName := message.GetPackageName()

		params = append(params, &template_TypeAliasParams{
			Name: fmt.Sprintf("%sPackage_%s", packageName, name),
			Type: "string",
		})
	}

	return template_TypeAlias{
		Name:       s.ref.GetName(),
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

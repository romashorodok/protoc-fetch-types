package typealias

import (
	"embed"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/tokenutils"
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
		var isArray bool
		if field, err := message.GetOriginField(); err == nil {
			isArray = tokenutils.TsArray(field)
		}

		params = append(params,
			&template_TypeAliasParams{
				Name:  tokenutils.TypeAliasParamName(message),
				Type:  tokenutils.TypeAliasName(message),
				Array: isArray,
			})
	}

	for _, field := range s.ref.GetPrimitiveFields() {
		params = append(params,
			&template_TypeAliasParams{
				Name:  field.GetName(),
				Type:  tokenutils.TsType(field.Type),
				Array: tokenutils.TsArray(field),
			},
		)
	}

	return template_TypeAlias{
		Name:       tokenutils.TypeAliasName(s.ref),
		Parameters: params,
	}
}

func (s *TypeAlias) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

const TYPE_ALIAS_TEMPLATE_FILE = "templates/type_alias.tmpl"

func New(storage embed.FS, ref *proxy.MessageProxy) *TypeAlias {
	templateFile, err := storage.ReadFile(TYPE_ALIAS_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", TYPE_ALIAS_TEMPLATE_FILE, storage)
	}
	tmpl := templatebuilder.New(templateFile, TYPE_ALIAS_TEMPLATE_FILE)

	return &TypeAlias{
		ref:  ref,
		tmpl: tmpl,
	}
}

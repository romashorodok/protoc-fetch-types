package importfrom

import (
	"embed"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

type ImportFrom struct {
	tmpl     *templatebuilder.TemplateBuilder
	template *template_ImportFrom
}

type template_ImportFrom struct {
	Namespace string
	AliasName string
	FilePath  string
}

func (s *ImportFrom) tmplStruct() *template_ImportFrom {
	return s.template
}

func (s *ImportFrom) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

type NewParams struct {
	Storage   embed.FS
	Namespace string
	AliasName string
	FilePath  string
}

const IMPORT_FROM_TEMPLATE_FILE = "templates/importfrom.tmpl"

func New(params *NewParams) *ImportFrom {
	templateFile, err := params.Storage.ReadFile(IMPORT_FROM_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", IMPORT_FROM_TEMPLATE_FILE, params.Storage)
	}
	tmpl := templatebuilder.New(templateFile, IMPORT_FROM_TEMPLATE_FILE)

	return &ImportFrom{
		tmpl: tmpl,
		template: &template_ImportFrom{
			Namespace: params.Namespace,
			AliasName: params.AliasName,
			FilePath:  params.FilePath,
		},
	}
}

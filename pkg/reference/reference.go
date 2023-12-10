package reference

import (
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

type Reference struct {
	tmpl *templatebuilder.TemplateBuilder
	name string
	path string
}

type template_Reference struct {
	Reference string
}

func (s *Reference) tmplStruct() *template_Reference {
	return &template_Reference{
		Reference: fmt.Sprintf("reference path=\"%s\"", s.path),
	}
}

func (s *Reference) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

type NewParams struct {
	Storage embed.FS
	Path    string
	Name    string
}

const IMPORT_TEMPLATE_FILE = "templates/reference.tmpl"

func New(params *NewParams) *Reference {
	templateFile, err := params.Storage.ReadFile(IMPORT_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", IMPORT_TEMPLATE_FILE, params.Storage)
	}
	tmpl := templatebuilder.New(templateFile, IMPORT_TEMPLATE_FILE)

	return &Reference{
		tmpl: tmpl,
		name: params.Name,
		path: params.Path,
	}
}

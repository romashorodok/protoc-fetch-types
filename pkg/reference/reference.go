package reference

import (
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

type Reference struct {
	tmpl     *templatebuilder.TemplateBuilder
	FilePath string
}

type template_Reference struct {
	Reference string
}

func (s *Reference) tmplStruct() *template_Reference {
	return &template_Reference{
		Reference: fmt.Sprintf("reference path=\"%s\"", s.FilePath),
	}
}

func (s *Reference) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

type NewParams struct {
	Storage  embed.FS
	FilePath string
}

const REFERENCE_TEMPLATE_FILE = "templates/reference.tmpl"

func New(params *NewParams) *Reference {
	templateFile, err := params.Storage.ReadFile(REFERENCE_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", REFERENCE_TEMPLATE_FILE, params.Storage)
	}
	tmpl := templatebuilder.New(templateFile, REFERENCE_TEMPLATE_FILE)

	return &Reference{
		tmpl:     tmpl,
		FilePath: params.FilePath,
	}
}

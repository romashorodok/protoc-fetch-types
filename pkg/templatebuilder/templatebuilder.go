package templatebuilder

import (
	"fmt"
	"html/template"
	"io"
	"log"
)

type TemplateBuilder struct {
	tmpl *template.Template
}

func (t *TemplateBuilder) WriteInto(in io.Writer, target any) error {
	return t.tmpl.Execute(in, target)
}

func generic(s string) template.HTML {
	return template.HTML(fmt.Sprintf("<%s>", s))
}

func closing(s string) template.HTML {
	return template.HTML(fmt.Sprintf("<%s />", s))
}

func unsafe(s string) template.HTML {
	return template.HTML(s)
}

func New(templateFile []byte, name string) *TemplateBuilder {
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"generic": generic,
			"unsafe":  unsafe,
            "closing": closing,
		}).
		Parse(string(templateFile))
	if err != nil {
		log.Panicf("Unable create tempalte %s. %s", name, err)
	}

	return &TemplateBuilder{
		tmpl: tmpl,
	}
}

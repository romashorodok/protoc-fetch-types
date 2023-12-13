package namespace

import (
	"embed"
	"io"
	"log"
	"strings"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
)

type _namespace struct {
	tmpl *templatebuilder.TemplateBuilder
	Name string
	Body string
}

type template_Namespace struct {
	Name string
	Body string
}

func (s *_namespace) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, &template_Namespace{
		Name: s.Name,
		Body: strings.ReplaceAll(s.Body, "\n", "\n"+"    "),
	})
}

type NewParams struct {
	Storage embed.FS
	Name    string
	Body    string
}

const NAMESPACE_TEMPLATE_FILE = "templates/namespace.tmpl"

func New(params *NewParams) *_namespace {
	templateFile, err := params.Storage.ReadFile(NAMESPACE_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", NAMESPACE_TEMPLATE_FILE, params.Storage)
	}
	tmpl := templatebuilder.New(templateFile, NAMESPACE_TEMPLATE_FILE)

	return &_namespace{
		tmpl: tmpl,
		Name: params.Name,
		Body: params.Body,
	}
}

type NamespaceTree interface {
	io.WriteCloser
	GetResult() *strings.Builder
}

type namespaceTree struct {
	namesapces       strings.Builder
	body             strings.Builder
	namespacesTokens []string
	storage          embed.FS
}

func (nt *namespaceTree) generateTree(namespaceTokens []string) strings.Builder {
	if len(namespaceTokens) == 0 {
		return nt.body
	}

	var result strings.Builder
	namespace := nt.generateTree(namespaceTokens[1:])

	nresult := New(&NewParams{
		Storage: nt.storage,
		Name:    namespaceTokens[0],
		Body:    namespace.String(),
	})

	_ = nresult.WriteInto(&result)

	return result
}

func (nt *namespaceTree) Close() error {
	nt.namesapces = nt.generateTree(nt.namespacesTokens)
	return nil
}

func (nt *namespaceTree) Write(p []byte) (n int, err error) {
	return nt.body.Write(p)
}

func (nt *namespaceTree) GetResult() *strings.Builder {
	return &nt.namesapces
}

var _ NamespaceTree = (*namespaceTree)(nil)

type NewNestedNamespaceParams struct {
	Storage         embed.FS
	NamespaceTokens []string
}

func NewNamespaceTree(params *NewNestedNamespaceParams) NamespaceTree {
	return &namespaceTree{
		namesapces:       strings.Builder{},
		body:             strings.Builder{},
		namespacesTokens: params.NamespaceTokens,
		storage:          params.Storage,
	}
}

package requestfunc

import (
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/tokenutils"
)

type RequestFunc struct {
	tmpl                    *templatebuilder.TemplateBuilder
	ref                     *proxy.MethodProxy
	messageFilenameRegistrt proxy.T_MessageFilenameRegistry
}

type template_RequestParams struct {
	Name string
	Type string
}

type template_RequestFunc struct {
	Name          string
	RequestMethod string
	BodyTypeName  string
	UriPath       string

	RequestParamsTypeName string
	RequestParams         []*template_RequestParams
}

func (s *RequestFunc) GetInputMessage() *proxy.MessageProxy {
	return s.ref.GetInputMessage()
}

func (s *RequestFunc) tmplStruct() *template_RequestFunc {
	inputMessage := s.GetInputMessage()

	pattern, requestMethod, err := googleHttpAnnotation(s.ref)
	if err != nil {
		log.Printf("%s. For %s requestfunc.googleHtttpAnnotation", err, s.ref.GetName())
		return nil
	}

	uriPath, params, err := parsePattern(pattern)
	if err != nil {
		log.Printf("%s. For %s requestfunc.parsePattern", err, s.ref.GetName())
		return nil
	}

	var requestParams []*template_RequestParams

	for _, param := range params {
		requestParams = append(requestParams, &template_RequestParams{
			Name: param,
			Type: "string",
		})
	}

	return &template_RequestFunc{
		Name:                  s.ref.GetName(),
		RequestMethod:         requestMethod,
		BodyTypeName:          tokenutils.TypeAliasName(inputMessage),
		UriPath:               uriPath,
		RequestParamsTypeName: fmt.Sprintf("%sParams", s.ref.GetName()),
		RequestParams:         requestParams,
	}
}

func (s *RequestFunc) WriteInto(in io.Writer) error {
	return s.tmpl.WriteInto(in, s.tmplStruct())
}

const TYPE_ALIAS_TEMPLATE_FILE = "templates/request_func.tmpl"

type NewParamsRequest struct {
	Storage                 embed.FS
	MessageFilenameRegistry proxy.T_MessageFilenameRegistry
	Ref                     *proxy.MethodProxy
}

func New(params *NewParamsRequest) *RequestFunc {
	templateFile, err := params.Storage.ReadFile(TYPE_ALIAS_TEMPLATE_FILE)
	if err != nil {
		log.Panicf("Unable read %s at storage %+v", TYPE_ALIAS_TEMPLATE_FILE, params.Storage)
	}
	tmpl := templatebuilder.New(templateFile, TYPE_ALIAS_TEMPLATE_FILE)

	return &RequestFunc{
		ref:                     params.Ref,
		tmpl:                    tmpl,
		messageFilenameRegistrt: params.MessageFilenameRegistry,
	}
}

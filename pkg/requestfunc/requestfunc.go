package requestfunc

import (
	"embed"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/templatebuilder"
	"google.golang.org/protobuf/types/descriptorpb"
)

type RequestFuncObject struct {
	Method  string
	Pattern string

	tmpl *templatebuilder.TemplateBuilder
	ref  *descriptorpb.MethodDescriptorProto
}

type NewParams struct {
	RequestFuncObject
}

func New(storage embed.FS, ref *descriptorpb.MethodDescriptorProto, params *NewParams) *RequestFuncObject {
	// templateFile, err := storage.ReadFile(resources.FUNC_REQUEST_TEMPLATE_FILE)
	// if err != nil {
	// 	log.Panicf("Unable read %s at storage %+v", resources.FUNC_REQUEST_TEMPLATE_FILE, storage)
	// }
	//
	// tmpl := templatebuilder.New(templateFile, resources.FUNC_REQUEST_TEMPLATE_FILE)

	return &RequestFuncObject{
		Method:  params.Method,
		Pattern: params.Pattern,

		// tmpl: tmpl,
		ref: ref,
	}
}

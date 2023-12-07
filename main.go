package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/requestfunc"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/typealias"
	"google.golang.org/protobuf/proto"
)

//go:embed templates/func_request.tmpl templates/type_alias.tmpl
var storage embed.FS

type (
	RequestFuncTree = map[resources.ProtoID]*requestfunc.RequestFuncObject
	TypeAliasTree   = map[resources.ProtoID]*typealias.TypeAliasObject

	MessageRegistry = map[resources.ProtoID]*proxy.MessageProxy
	MethodRegistry  = map[resources.ProtoID]*proxy.MethodProxy
)

var (
	methodRegistry  MethodRegistry  = make(MethodRegistry)
	messageRegistry MessageRegistry = make(MessageRegistry)
)

func fillRegistry(request *plugin.CodeGeneratorRequest) {
	for _, protoFile := range request.ProtoFile {
		packageName := protoFile.GetPackage()

		// TODO: i know dependency of each file, i need get messages from it by aliasing it
		// each package may be with same name but exist in different directory
		// user/meta/fields.proto - meta
		// product/meta/fields.proto - meta
		log.Println(protoFile.GetDependency())

		for _, protoService := range protoFile.Service {
			serviceName := protoService.GetName()

			for _, method := range protoService.Method {
				methodProxy := proxy.NewMethodProxy(
					fmt.Sprintf(".%s.%s", packageName, serviceName),
					method,
				)
				methodRegistry[methodProxy.GetProtoID()] = methodProxy
			}
		}

		for _, protoMessage := range protoFile.MessageType {
			messageProxy := proxy.NewMessageProxy(
				fmt.Sprintf(".%s", packageName),
				protoMessage,
			)
			messageRegistry[messageProxy.GetProtoID()] = messageProxy
		}
	}
}

func generate(req *plugin.CodeGeneratorRequest) string {
	var builder strings.Builder

	fillRegistry(req)

	for protoID, method := range methodRegistry {
		methodInputMessageID := method.GetInputType()
		_ = protoID

		message, exist := messageRegistry[methodInputMessageID]
		if !exist {
			continue
		}

		message.GetFieldsMessages()

		// log.Println(messageFields)

	}

	// typeAliasTree := make(TypeAliasTree)
	// requestFuncTree := make(RequestFuncTree)

	_ = req

	_ = storage

	return builder.String()
}

func main() {
	request, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Panic("Unable read stdin", err)
	}

	req := &plugin.CodeGeneratorRequest{}
	if err = proto.Unmarshal(request, req); err != nil {
		log.Panic("Unable deserialize request", err)
	}

	_ = generate(req)

	resp := &plugin.CodeGeneratorResponse{
		File: []*plugin.CodeGeneratorResponse_File{
			{
				Name:    proto.String("fetch-types.ts"),
				Content: proto.String("test"),
			},
		},
	}

	res, err := proto.Marshal(resp)
	if err != nil {
		log.Panic("Unable serialize response", err)
	}

	_, err = os.Stdout.Write(res)
	if err != nil {
		log.Panic("Unable send response", err)
	}
}

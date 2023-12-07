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
	T_requestFuncTree = map[resources.ProtoID]*requestfunc.RequestFuncObject
	T_typeAliasTree   = map[resources.ProtoID]*typealias.TypeAliasObject
)

var (
	methodRegistry = make(proxy.T_methodRegistry)
	// TODO: Use registry with file name instead
	messageRegistry = make(proxy.T_messageRegistry)

	messageFilenameRegistry = make(proxy.T_messageFilenameRegistry)
)

func fillRegistry(request *plugin.CodeGeneratorRequest) {
	for _, protoFile := range request.ProtoFile {
		packageName := protoFile.GetPackage()

		// TODO: i know dependency of each file, i need get messages from it by aliasing it
		// each package may be with same name but exist in different directory
		// user/meta/fields.proto - meta
		// product/meta/fields.proto - meta
		// log.Println(protoFile.GetDependency())

		for _, protoService := range protoFile.Service {
			serviceName := protoService.GetName()

			for _, method := range protoService.Method {
				methodProxy := proxy.NewMethodProxy(
					fmt.Sprintf(".%s.%s", packageName, serviceName),
					protoFile,
					method,
				)
				methodRegistry[methodProxy.GetProtoID()] = methodProxy
			}
		}

		for _, protoMessage := range protoFile.MessageType {
			messageProxy := proxy.NewMessageProxy(
				&proxy.NewMessageProxyParams{
					PackageID:               fmt.Sprintf(".%s", packageName),
					File:                    protoFile,
					DescriptorProto:         protoMessage,
					MessageFilenameRegistry: messageFilenameRegistry,
				},
			)
			messageRegistry[messageProxy.GetProtoID()] = messageProxy
			messageFilenameRegistry[messageProxy.GetFilenameProtoID()] = messageProxy
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

		nestedMessages := message.GetFieldsMessages()

        log.Println("For message", message.GetName())
		log.Println(nestedMessages)


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

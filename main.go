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
	T_typeAliasTree   = map[resources.FilenameProtoID]*typealias.TypeAlias
)

// requestFuncTree = make(T_requestFuncTree)
var typeAliasTree = make(T_typeAliasTree)

var (
	methodRegistry = make(proxy.T_methodRegistry)
	// TODO: Use registry with file name instead
	messageRegistry = make(proxy.T_messageRegistry)

	messageFilenameRegistry = make(proxy.T_messageFilenameRegistry)
)

func fillRegistry(request *plugin.CodeGeneratorRequest) {
	for _, protoFile := range request.ProtoFile {

		packageName := protoFile.GetPackage()

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

		for _, msg := range nestedMessages {
			typeAliasTree[msg.GetFilenameProtoID()] = typealias.New(storage, msg)
		}
		typeAliasTree[message.GetFilenameProtoID()] = typealias.New(storage, message)
	}

	for _, typeAliases := range typeAliasTree {
		typeAliases.WriteInto(&builder)
	}

	log.Println(typeAliasTree)

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

	targetFile := "fetch_types.proto"

	var generateTargetFile bool
	for _, file := range req.FileToGenerate {
		if file == targetFile {
			generateTargetFile = true
			break
		}
	}

	if generateTargetFile {
		result := generate(req)

		respFiles := []*plugin.CodeGeneratorResponse_File{
			{
				Name:    proto.String(fmt.Sprintf("%s%s", strings.TrimSuffix(targetFile, ".proto"), ".ts")),
				Content: proto.String(result),
			},
		}

		resp := &plugin.CodeGeneratorResponse{
			File: respFiles,
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
}

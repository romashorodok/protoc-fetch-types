package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/namespace"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/requestfunc"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/typealias"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

//go:embed templates/request_func.tmpl templates/type_alias.tmpl templates/namespace.tmpl
var storage embed.FS

const UnixPathSeparator = "/"

type (
	T_requestFuncTree = map[resources.ProtoID]*requestfunc.RequestFunc
	T_typeAliasTree   = map[resources.FilenameProtoID]*typealias.TypeAlias

	T_file      = map[string]*descriptorpb.FileDescriptorProto
	T_namespace = map[string]*descriptorpb.FileDescriptorProto
)

var (
	requestFuncTree = make(T_requestFuncTree)
	typeAliasTree   = make(T_typeAliasTree)
)

var (
	methodRegistry    = make(proxy.T_MethodRegistry)
	fileRegistry      = make(T_file)
	namespaceRegistry = make(T_namespace)

	messageFilenameRegistry = make(proxy.T_MessageFilenameRegistry)
	methodFilenameRegistry  = make(proxy.T_MethodFilenameRegistry)
)

func GetNamespaceTokens(file *descriptorpb.FileDescriptorProto) []string {
	namespace := strings.Split(file.GetName(), UnixPathSeparator)
	namespace = namespace[:len(namespace)-1]
	if len(namespace) == 0 || namespace[len(namespace)-1] != file.GetPackage() {
		namespace = append(namespace, file.GetPackage())
	}
	return namespace
}

func GetNamespace(file *descriptorpb.FileDescriptorProto) string {
	return strings.Join(GetNamespaceTokens(file), ".")
}

func fillRegistry(request *plugin.CodeGeneratorRequest) {
	var names []string
	// Protoc start requesting top level files in tree from bottom to up
	// I cannot access from bottom files upper files. They must be
	for _, file := range request.ProtoFile {
		names = append(names, file.GetName())

		fileRegistry[file.GetName()] = file
		namespaceRegistry[GetNamespace(file)] = file

		packageName := file.GetPackage()

		for _, protoService := range file.Service {
			serviceName := protoService.GetName()

			for _, method := range protoService.Method {
				methodProxy := proxy.NewMethodProxy(
					&proxy.NewMethodProxyParams{
						ServiceID:               fmt.Sprintf(".%s.%s", packageName, serviceName),
						File:                    file,
						MethodDescriptor:        method,
						MessageFilenameRegistry: messageFilenameRegistry,
					},
				)
				methodRegistry[methodProxy.GetProtoID()] = methodProxy
				methodFilenameRegistry[methodProxy.GetFilenameProtoID()] = methodProxy
			}
		}

		for _, protoMessage := range file.MessageType {
			messageProxy := proxy.NewMessageProxy(
				&proxy.NewMessageProxyParams{
					PackageID:               fmt.Sprintf(".%s", packageName),
					File:                    file,
					DescriptorProto:         protoMessage,
					MessageFilenameRegistry: messageFilenameRegistry,
				},
			)
			messageFilenameRegistry[messageProxy.GetFilenameProtoID()] = messageProxy
		}
	}

	log.Printf("%s files has access to %s", request.FileToGenerate, names)
}

func tsFilename(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name)) + ".ts"
}

func generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	var response plugin.CodeGeneratorResponse

	fillRegistry(req)

	for _, requestedFile := range req.FileToGenerate {
		file, exist := fileRegistry[requestedFile]
		if !exist {
			log.Printf("[Warning] Requestd file not found in registry. Ignoring %s", requestedFile)
			continue
		}
		var builder strings.Builder

		namespaceWriter := namespace.NewNamespaceTree(
			&namespace.NewNestedNamespaceParams{
				Storage:         storage,
				NamespaceTokens: GetNamespaceTokens(file),
			},
		)

		for _, method := range methodFilenameRegistry {
			if requestedFile != method.GetFileName() {
				continue
			}

			requestFunc := requestfunc.New(
				&requestfunc.NewParamsRequest{
					Storage:                 storage,
					MessageFilenameRegistry: messageFilenameRegistry,
					Ref:                     method,
				},
			)

			_ = requestFunc.WriteInto(namespaceWriter)
		}

		// for _, message := range messageFilenameRegistry {
		// 	typeAlias := typealias.New(storage, message)
		//
		// 	_ = typeAlias.WriteInto(namespaceWriter)
		// }

		// log.Println("Current file", requestedFile)
		// _ = requestFuncs
		// _ = typeAliases

		namespaceWriter.Close()
		builder.WriteString(namespaceWriter.GetResult().String())

		response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(tsFilename(requestedFile)),
			Content: proto.String(builder.String()),
		})
	}

	// fillRegistry(req)

	// for _, method := range methodRegistry {
	// 	inputMessage := method.GetInputMessage()
	// 	outputMessage := method.GetOutputMessage()
	//
	// 	_, exist := messageFilenameRegistry[inputMessage.GetFilenameProtoID()]
	// 	if !exist {
	// 		continue
	// 	}
	//
	// 	for _, msg := range inputMessage.GetFieldsMessages() {
	// 		typeAliasTree[msg.GetFilenameProtoID()] = typealias.New(storage, msg)
	// 	}
	// 	for _, msg := range outputMessage.GetFieldsMessages() {
	// 		typeAliasTree[msg.GetFilenameProtoID()] = typealias.New(storage, msg)
	// 	}
	//
	// 	typeAliasTree[inputMessage.GetFilenameProtoID()] = typealias.New(storage, inputMessage)
	// 	typeAliasTree[outputMessage.GetFilenameProtoID()] = typealias.New(storage, outputMessage)
	//
	// 	requestFuncTree[method.GetFilenameProtoID()] = requestfunc.New(
	// 		&requestfunc.NewParamsRequest{
	// 			Storage:                 storage,
	// 			MessageFilenameRegistry: messageFilenameRegistry,
	// 			Ref:                     method,
	// 		},
	// 	)
	// }
	//
	// for _, typeAliases := range typeAliasTree {
	// 	_ = typeAliases.WriteInto(&builder)
	// }
	//
	// for _, requestFunc := range requestFuncTree {
	// 	_ = requestFunc.WriteInto(&builder)
	// }

	// return builder.String()

	return &response
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

	// targetFile := "fetch_types.proto"

	// var generateTargetFile bool
	// for _, file := range req.FileToGenerate {
	// 	if file == targetFile {
	// 		generateTargetFile = true
	// 		break
	// 	}
	// }

	resp := generate(req)

	// respFiles := []*plugin.CodeGeneratorResponse_File{
	// 	{
	// 		Name:    proto.String(fmt.Sprintf("%s%s", strings.TrimSuffix(targetFile, ".proto"), ".ts")),
	// 		Content: proto.String(result),
	// 	},
	// }

	// resp := &plugin.CodeGeneratorResponse{
	// 	File: respFiles,
	// }

	res, err := proto.Marshal(resp)
	if err != nil {
		log.Panic("Unable serialize response", err)
	}

	_, err = os.Stdout.Write(res)
	if err != nil {
		log.Panic("Unable send response", err)
	}
}

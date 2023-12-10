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
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/reference"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/requestfunc"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/typealias"
	"google.golang.org/protobuf/proto"
)

//go:embed templates/request_func.tmpl templates/type_alias.tmpl templates/namespace.tmpl templates/reference.tmpl
var storage embed.FS

func fillRegistry(request *plugin.CodeGeneratorRequest) *proxy.Registry {
	registry := proxy.NewRegistry()

	var fileNames []string
	var packageNames []string

	// Protoc start requesting top level files in tree from bottom to up
	// I cannot access from bottom files upper files. They must be imported
	for _, file := range request.ProtoFile {
		fileNames = append(fileNames, file.GetName())
		packageNames = append(packageNames, file.GetPackage())
		registry.File[file.GetName()] = file

		for _, protoService := range file.Service {
			serviceName := protoService.GetName()

			for _, method := range protoService.Method {
				methodProxy := proxy.NewMethodProxy(
					&proxy.NewMethodProxyParams{
						MethodDescriptor: method,
						ServiceID:        fmt.Sprintf(".%s.%s", file.GetPackage(), serviceName),
						Registy:          registry,
						File:             file,
					},
				)
				registry.Method[methodProxy.GetFilenameProtoID()] = methodProxy
			}
		}

		for _, protoMessage := range file.MessageType {
			messageProxy := proxy.NewMessageProxy(
				&proxy.NewMessageProxyParams{
					DescriptorProto: protoMessage,
					Registry:        registry,
					File:            file,
				},
			)
			registry.Message[messageProxy.GetFilenameProtoID()] = messageProxy
		}
	}
	log.Printf("%s files has access to %s files contains %s", request.FileToGenerate, fileNames, packageNames)
	return registry
}

func tsFilename(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name)) + ".ts"
}

func generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	registry := fillRegistry(req)
	var response plugin.CodeGeneratorResponse

	for _, requestedFile := range req.FileToGenerate {
		file, exist := registry.File[requestedFile]
		if !exist {
			log.Printf("[Warning] Requestd file not found in registry. Ignoring %s", requestedFile)
			continue
		}
		var builder strings.Builder
		log.Printf("[%s] file has dependency %s", file.GetName(), file.GetDependency())

		for _, filePath := range file.GetDependency() {
			dependencyFile, exist := registry.File[filePath]
			if !exist {
				continue
			}
			importPath := tsFilename(dependencyFile.GetName())
			_ = reference.New(&reference.NewParams{
				Storage: storage,
				Path:    importPath,
				Name:    dependencyFile.GetPackage(),
			}).WriteInto(&builder)
		}

		namespaceWriter := namespace.NewNamespaceTree(
			&namespace.NewNestedNamespaceParams{
				Storage:         storage,
				NamespaceTokens: proxy.GetNamespaceTokens(file),
			},
		)

		for _, method := range registry.Method {
			if requestedFile != method.GetFileName() {
				continue
			}

			// TODO: rename input to request
			inputMessage := method.GetInputMessage()

			for _, localMessage := range inputMessage.GetLocalFieldMessages() {
				_, exist := registry.AlredyExisted[localMessage.GetFilenameProtoID()]
				if exist {
					continue
				}
				_ = typealias.New(storage, localMessage).
					WriteInto(namespaceWriter)
				registry.AlredyExisted[localMessage.GetFilenameProtoID()] = struct{}{}
			}

			_, exist := registry.AlredyExisted[inputMessage.GetFilenameProtoID()]
			if !exist {
				_ = typealias.New(storage, inputMessage).
					WriteInto(namespaceWriter)
				registry.AlredyExisted[inputMessage.GetFilenameProtoID()] = struct{}{}
			}

			_, exist = registry.AlredyExisted[method.GetFilenameProtoID()]
			if !exist {
				_ = requestfunc.New(
					&requestfunc.NewParamsRequest{
						Storage: storage,
						Ref:     method,
					},
				).WriteInto(namespaceWriter)
				registry.AlredyExisted[method.GetFilenameProtoID()] = struct{}{}
			}
		}

		namespaceWriter.Close()
		builder.WriteString(namespaceWriter.GetResult().String())

		response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(tsFilename(requestedFile)),
			Content: proto.String(builder.String()),
		})
	}
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

package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/importfrom"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/namespace"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/proxy"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/reference"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/requestfunc"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/tokenutils"
	"github.com/romashorodok/protoc-gen-fetch-types/pkg/typealias"
	"google.golang.org/protobuf/proto"
)

const FETCH_TYPES_FILENAME = "fetch_types.proto"

//go:embed templates/request_func.tmpl templates/type_alias.tmpl templates/namespace.tmpl templates/reference.tmpl templates/importfrom.tmpl
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

type GeneratorOptions struct {
	ImportIgnorePrefixes []string
}

func options(compilerOptions string) *GeneratorOptions {
	var result GeneratorOptions

	options := strings.Split(compilerOptions, ",")

	for _, opt := range options {
		slice := strings.Split(opt, "=")
		name, value := slice[0], slice[1]
		switch name {
		case "import_ignore":
			re := regexp.MustCompile(`"([^"]+)"`)
			matches := re.FindAllString(value, -1)
			for idx, match := range matches {
				match = strings.Trim(match, `"`)
				matches[idx] = match
			}
			result.ImportIgnorePrefixes = matches
		}
	}
	return &result
}

func generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	registry := fillRegistry(req)
	var response plugin.CodeGeneratorResponse

	generatorOptions := options(req.GetParameter())

	for _, requestedFile := range req.FileToGenerate {
		file, exist := registry.File[requestedFile]
		if !exist {
			log.Printf("[Warning] Requestd file not found in registry. Ignoring %s", requestedFile)
			continue
		}
		var builder strings.Builder
		log.Printf("[%s] file has dependency %s", file.GetName(), file.GetDependency())
		requestedNamespace := proxy.GetNamespaceTokens(file)
		if tokenutils.HasNamespaceToken(requestedNamespace, generatorOptions.ImportIgnorePrefixes) {
			log.Printf("[IGNORE] Ignored generating %s", file.GetName())
			continue
		}

		// TODO: Make file proxy
		for _, filePath := range file.GetDependency() {
			dependencyFile, exist := registry.File[filePath]
			if !exist {
				continue
			}

			// If dependency file in ignore just skip it
			namespaceTokens := proxy.GetNamespaceTokens(dependencyFile)
			if tokenutils.HasNamespaceToken(namespaceTokens, generatorOptions.ImportIgnorePrefixes) {
				log.Printf("[IGNORE] Ignored import for %s", dependencyFile.GetName())
				continue
			}

			// Check if current file is not the root of tree.
			if !tokenutils.IsRoot(requestedFile) {
				// If dependency file is on the root. I need make backward path
				if tokenutils.IsRoot(filePath) {
					backwards := tokenutils.GetBackwardCount(requestedFile)
					filePath = tokenutils.AppendBackwards(filePath, backwards)
				} else {
					filePath = tokenutils.BackwardPath(filePath)
				}
			} else {
				filePath = dependencyFile.GetName()
			}
			filePath = tsFilename(filePath)

			_ = reference.New(&reference.NewParams{
				Storage:  storage,
				FilePath: filePath,
			}).WriteInto(&builder)

			_ = importfrom.New(&importfrom.NewParams{
				Storage:   storage,
				Namespace: proxy.PackageNamespaceSuffix(dependencyFile),
				AliasName: proxy.ImportAliasFromFilePath(dependencyFile),
				FilePath:  "./" + strings.TrimSuffix(filePath, ".ts"),
			}).WriteInto(&builder)
		}

		namespaceWriter := namespace.NewNamespaceTree(
			&namespace.NewNestedNamespaceParams{
				Storage:         storage,
				NamespaceTokens: []string{proxy.GetNamespace(file)},

				// NOTE: When will need nested namespaces
				// NamespaceTokens: proxy.GetNamespaceTokens(file),
			},
		)

		for filenameProtoID, message := range registry.Message {
			if !strings.HasPrefix(filenameProtoID, requestedFile) {
				continue
			}
			_, exist := registry.AlredyExisted[message.GetFilenameProtoID()]
			if !exist {
				_ = typealias.New(storage, message).
					WriteInto(namespaceWriter)
				registry.AlredyExisted[message.GetFilenameProtoID()] = struct{}{}
			}
		}

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

	resp := generate(req)

	res, err := proto.Marshal(resp)
	if err != nil {
		log.Panic("Unable serialize response", err)
	}

	_, err = os.Stdout.Write(res)
	if err != nil {
		log.Panic("Unable send response", err)
	}
}

package proxy

import (
	"fmt"
	"strings"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MethodProxy struct {
	*descriptorpb.MethodDescriptorProto
	serviceID               string
	file                    *descriptorpb.FileDescriptorProto
	messageFilenameRegistry T_messageFilenameRegistry
}

func (p *MethodProxy) GetPackageName() string {
	return p.file.GetPackage()
}

func (p *MethodProxy) GetFilenameProtoID() resources.FilenameProtoID {
	return fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, p.GetName())
}

// TOOD: refactor it and make reusable
func (p *MethodProxy) searchMessageByTypeName(typeName string) *MessageProxy {
	typeNameParts := strings.Split(typeName, ".")
	typeNameSuffix := typeNameParts[len(typeNameParts)-1]

	localMessage := fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, typeNameSuffix)

	// If message in same package as parent message
	{
		message, exist := p.messageFilenameRegistry[localMessage]
		if exist {
			return message
		}
	}

	// If in other packages
	for _, dependencyFile := range p.file.GetDependency() {
		filenameProtoID := fmt.Sprintf("%s:%s", dependencyFile, typeNameSuffix)

		message, exist := p.messageFilenameRegistry[filenameProtoID]
		if !exist {
			continue
		}

		targetTypeName := message.GetProtoID()

		if targetTypeName == typeName {
			return message
		}
	}

	return nil
}

func (p *MethodProxy) GetInputMessage() *MessageProxy {
	return p.searchMessageByTypeName(p.GetInputType())
}

func (p *MethodProxy) GetOutputMessage() *MessageProxy {
	return p.searchMessageByTypeName(p.GetOutputType())
}

func (p *MethodProxy) GetProtoID() string {
	return fmt.Sprintf("%s.%s", p.serviceID, p.GetName())
}

var _ ProtoProxy = (*MethodProxy)(nil)

type NewMethodProxyParams struct {
	ServiceID               string
	File                    *descriptorpb.FileDescriptorProto
	MethodDescriptor        *descriptorpb.MethodDescriptorProto
	MessageFilenameRegistry map[resources.FilenameProtoID]*MessageProxy
}

func NewMethodProxy(prarams *NewMethodProxyParams) *MethodProxy {
	return &MethodProxy{
		MethodDescriptorProto:   prarams.MethodDescriptor,
		serviceID:               prarams.ServiceID,
		file:                    prarams.File,
		messageFilenameRegistry: prarams.MessageFilenameRegistry,
	}
}

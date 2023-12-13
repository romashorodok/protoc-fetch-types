package proxy

import (
	"fmt"
	"strings"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MethodProxy struct {
	*descriptorpb.MethodDescriptorProto

	serviceID string
	registry  *Registry
	file      *descriptorpb.FileDescriptorProto
}

func (p *MethodProxy) GetPackageName() string {
	return p.file.GetPackage()
}

func (p *MethodProxy) GetFileName() string {
	return p.file.GetName()
}

func (p *MethodProxy) GetFilenameProtoID() resources.FilenameProtoID {
	return fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, p.GetName())
}

// Search message in current package
func (p *MethodProxy) searcLocalMessageByTypeName(typeName string) *MessageProxy {
	typeNameParts := strings.Split(typeName, ".")
	typeNameSuffix := typeNameParts[len(typeNameParts)-1]
	localMessage := fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, typeNameSuffix)
	{
		message, exist := p.registry.Message[localMessage]
		if exist {
			return message
		}
	}
	return nil
}

func (p *MethodProxy) searchMessageByTypeName(typeName string) *MessageProxy {
	if localMessage := p.searcLocalMessageByTypeName(typeName); localMessage != nil {
		return localMessage
	}

	typeNameParts := strings.Split(typeName, ".")
	typeNameSuffix := typeNameParts[len(typeNameParts)-1]
	for _, dependencyFile := range p.file.GetDependency() {
		filenameProtoID := fmt.Sprintf("%s:%s", dependencyFile, typeNameSuffix)
		message, exist := p.registry.Message[filenameProtoID]
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
	MethodDescriptor *descriptorpb.MethodDescriptorProto
	ServiceID        string
	Registy          *Registry
	File             *descriptorpb.FileDescriptorProto
}

func NewMethodProxy(prarams *NewMethodProxyParams) *MethodProxy {
	return &MethodProxy{
		MethodDescriptorProto: prarams.MethodDescriptor,
		serviceID:             prarams.ServiceID,
		registry:              prarams.Registy,
		file:                  prarams.File,
	}
}

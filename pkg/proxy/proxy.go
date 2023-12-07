package proxy

import (
	"fmt"
	"strings"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"google.golang.org/protobuf/types/descriptorpb"
)

const filenameSeparator = ":"

type (
	T_messageRegistry         = map[resources.ProtoID]*MessageProxy
	T_methodRegistry          = map[resources.ProtoID]*MethodProxy
	T_messageFilenameRegistry = map[resources.FilenameProtoID]*MessageProxy
)

type ProtoProxy interface {
	GetProtoID() string
	GetFilenameProtoID() resources.FilenameProtoID
	GetPackageName() string
}

type MethodProxy struct {
	*descriptorpb.MethodDescriptorProto
	serviceID string
	file      *descriptorpb.FileDescriptorProto
}

func (p *MethodProxy) GetPackageName() string {
	return p.file.GetPackage()
}

func (p *MethodProxy) GetFilenameProtoID() resources.FilenameProtoID {
	return fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, p.GetName())
}

func (p *MethodProxy) GetProtoID() string {
	return fmt.Sprintf("%s.%s", p.serviceID, p.GetName())
}

var _ ProtoProxy = (*MethodProxy)(nil)

func NewMethodProxy(serviceID string, file *descriptorpb.FileDescriptorProto, pb *descriptorpb.MethodDescriptorProto) *MethodProxy {
	return &MethodProxy{
		MethodDescriptorProto: pb,
		serviceID:             serviceID,
		file:                  file,
	}
}

type MessageProxy struct {
	*descriptorpb.DescriptorProto
	packageID               string
	file                    *descriptorpb.FileDescriptorProto
	messageFilenameRegistry T_messageFilenameRegistry
}

// GetPackageName implements ProtoProxy.
func (p *MessageProxy) GetPackageName() string {
	return p.file.GetPackage()
}

func (p *MessageProxy) GetFilenameProtoID() resources.FilenameProtoID {
	return fmt.Sprintf("%s%s%s", p.file.GetName(), filenameSeparator, p.GetName())
}

func (p *MessageProxy) searchMessageByTypeName(typeName string) *MessageProxy {
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

func (p *MessageProxy) GetFieldsMessages() []*MessageProxy {
	var result []*MessageProxy

	for _, field := range p.Field {
		if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}

		result = append(result,
			p.searchMessageByTypeName(field.GetTypeName()),
		)
	}

	return result
}

func (p *MessageProxy) GetProtoID() string {
	return fmt.Sprintf("%s.%s", p.packageID, p.GetName())
}

var _ ProtoProxy = (*MessageProxy)(nil)

type NewMessageProxyParams struct {
	PackageID               string
	File                    *descriptorpb.FileDescriptorProto
	DescriptorProto         *descriptorpb.DescriptorProto
	MessageFilenameRegistry T_messageFilenameRegistry
}

func NewMessageProxy(params *NewMessageProxyParams) *MessageProxy {
	return &MessageProxy{
		DescriptorProto:         params.DescriptorProto,
		packageID:               params.PackageID,
		file:                    params.File,
		messageFilenameRegistry: params.MessageFilenameRegistry,
	}
}

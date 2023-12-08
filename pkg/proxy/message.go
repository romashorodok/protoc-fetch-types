package proxy

import (
	"fmt"
	"strings"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MessageProxy struct {
	*descriptorpb.DescriptorProto
	packageID               string
	file                    *descriptorpb.FileDescriptorProto
	messageFilenameRegistry T_messageFilenameRegistry

	// Indicates the origin field when the proxy is generated by another proxy field.
	origin *descriptorpb.FieldDescriptorProto
}

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
		message := p.searchMessageByTypeName(field.GetTypeName())
		if message != nil {
			message.origin = field
		}
		result = append(result, message)
	}
	return result
}

func (p *MessageProxy) GetPrimitiveFields() []*descriptorpb.FieldDescriptorProto {
	var result []*descriptorpb.FieldDescriptorProto
	for _, field := range p.Field {
		if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}
		result = append(result, field)
	}
	return result
}

func (p *MessageProxy) GetOriginField() (*descriptorpb.FieldDescriptorProto, error) {
	if p.origin == nil {
		return nil, NotFoundOriginError
	}
	return p.origin, nil
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
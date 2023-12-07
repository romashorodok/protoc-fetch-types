package proxy

import (
	"fmt"
	"log"

	"google.golang.org/protobuf/types/descriptorpb"
	// "google.golang.org/protobuf/reflect/protodesc"
)

type ProtoProxy interface {
	GetProtoID() string
}

type MethodProxy struct {
	*descriptorpb.MethodDescriptorProto
	serviceID string
}

func (p *MethodProxy) GetProtoID() string {
	return fmt.Sprintf("%s.%s", p.serviceID, p.GetName())
}

var _ ProtoProxy = (*MethodProxy)(nil)

func NewMethodProxy(serviceID string, pb *descriptorpb.MethodDescriptorProto) *MethodProxy {
	return &MethodProxy{
		MethodDescriptorProto: pb,
		serviceID:             serviceID,
	}
}

type MessageProxy struct {
	*descriptorpb.DescriptorProto
	packageID string
}

func (p *MessageProxy) GetFieldsMessages() []*MessageProxy {
	log.Println(p.GetName())

	var result []*MessageProxy

	for _, field := range p.Field {
		if field.GetType() != descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}

		// TODO: get nested fields
	}

	return result
}

func (p *MessageProxy) GetProtoID() string {
	return fmt.Sprintf("%s.%s", p.packageID, p.GetName())
}

var _ ProtoProxy = (*MessageProxy)(nil)

func NewMessageProxy(packageID string, pb *descriptorpb.DescriptorProto) *MessageProxy {
	return &MessageProxy{
		DescriptorProto: pb,
		packageID:       packageID,
	}
}

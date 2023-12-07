package resources

import (
	"log"

	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	TYPE_ALIAS_TEMPLATE_FILE   = "templates/type_alias.tmpl"
	FUNC_REQUEST_TEMPLATE_FILE = "templates/func_request.tmpl"
)

type ProtoID = string

var (
	messageRegistry = make(map[ProtoID]*descriptorpb.DescriptorProto)
	methodRegistry  = make(map[ProtoID]*descriptorpb.MethodDescriptorProto)
)

//***

func SetMessage(id ProtoID, pb *descriptorpb.DescriptorProto) {
	messageRegistry[id] = pb
}

func GetMessage(id ProtoID) *descriptorpb.DescriptorProto {
	message, exist := messageRegistry[id]
	if !exist {
		log.Panicf("Unable find `%s` message in registry.", id)
	}
	return message
}

//***

func SetMethod(id ProtoID, pb *descriptorpb.MethodDescriptorProto) {
	methodRegistry[id] = pb
}

func GetMethod(id ProtoID) *descriptorpb.MethodDescriptorProto {
	message, exist := methodRegistry[id]
	if !exist {
		log.Panicf("Unable find `%s` method in registry.", id)
	}
	return message
}

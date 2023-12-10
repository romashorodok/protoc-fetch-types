package proxy

import (
	"errors"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
	"google.golang.org/protobuf/types/descriptorpb"
)

var NotFoundOriginError error = errors.New("not found parent fiel")

const filenameSeparator = ":"

type name = string

type (
	R_AlredyExisted map[resources.FilenameProtoID]struct{}
	R_Method        map[resources.FilenameProtoID]*MethodProxy
	R_Message       map[resources.FilenameProtoID]*MessageProxy
	R_File          map[name]*descriptorpb.FileDescriptorProto
)

type ProtoProxy interface {
	GetProtoID() string
	GetFilenameProtoID() resources.FilenameProtoID
	GetPackageName() string
}

type Registry struct {
	// Contains only local alredy existed keys
	AlredyExisted R_AlredyExisted
	Message       R_Message
	Method        R_Method
	File          R_File
}

func NewRegistry() *Registry {
	return &Registry{
		AlredyExisted: make(R_AlredyExisted),
		Message:       make(R_Message),
		Method:        make(R_Method),
		File:          make(R_File),
	}
}

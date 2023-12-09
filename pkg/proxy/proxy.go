package proxy

import (
	"errors"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
)

var NotFoundOriginError error = errors.New("not found parent fiel")

const filenameSeparator = ":"

type (
	T_MessageRegistry = map[resources.ProtoID]*MessageProxy
	T_MethodRegistry  = map[resources.ProtoID]*MethodProxy

	T_MessageFilenameRegistry = map[resources.FilenameProtoID]*MessageProxy
	T_MethodFilenameRegistry  = map[resources.FilenameProtoID]*MethodProxy
)

type ProtoProxy interface {
	GetProtoID() string
	GetFilenameProtoID() resources.FilenameProtoID
	GetPackageName() string
}

package proxy

import (
	"errors"

	"github.com/romashorodok/protoc-gen-fetch-types/pkg/resources"
)

var NotFoundOriginError error = errors.New("not found parent fiel")

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

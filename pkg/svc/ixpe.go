package svc

//go:generate mkdir -p ../proto/generated
//go:generate sh -c "protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

import (
	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
)

// IPXEBuilder manages iPXEs
type IPXEBuilder struct {
	iPXEBuilder.UnimplementedIPXEBuilderServer
}

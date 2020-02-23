package svc

//go:generate mkdir -p ../proto/generated
//go:generate sh -c "protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

import (
	"errors"
	"strings"

	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	"gitlab.com/bloom42/libs/rz-go/log"
)

const (
	totalCompileSteps = 2247
)

// IPXEBuilder manages iPXEs.
type IPXEBuilder struct {
	iPXEBuilder.UnimplementedIPXEBuilderServer
	Builder *workers.Builder
}

// Extract extracts the iPXE source code.
func (i *IPXEBuilder) Extract() error {
	return i.Builder.Extract()
}

// Create creates an iPXE.
func (i *IPXEBuilder) Create(req *iPXEBuilder.IPXE, srv iPXEBuilder.IPXEBuilder_CreateServer) error {
	delta := int64(totalCompileSteps)
	stdoutChan, stderrChan, doneChan, errChan := make(chan string), make(chan string), make(chan string), make(chan error)

	go i.Builder.Build(req.GetScript(), req.GetPlatform(), req.GetDriver(), req.GetExtension(), stdoutChan, stderrChan, doneChan, errChan)

	log.Info("Building iPXE")

	for {
		select {
		case <-stdoutChan:
			delta--
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: delta,
				Path:  "",
			}); err != nil {
				return err
			}
		case stderrErr := <-stderrChan:
			if !strings.Contains(stderrErr, "ar:") { // `ar` prints to `stderr` for some reason
				return errors.New(stderrErr)
			}

			delta--
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: delta,
				Path:  "",
			}); err != nil {
				return err
			}
		case outPath := <-doneChan:
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: 0,
				Path:  outPath,
			}); err != nil {
				return err
			}

			return nil
		case err := <-errChan:
			return err
		}
	}
}

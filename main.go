package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/pojntfx/ipxebuilderd/pkg/svc"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	builder := workers.Builder{
		BasePath: filepath.Join(os.TempDir(), "ipxebuilderd", uuid.NewV4().String()),
	}

	listener, err := net.Listen("tcp", ":3009")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	iPXEService := svc.IPXEBuilder{
		Builder: &builder,
	}

	iPXEBuilder.RegisterIPXEBuilderServer(server, &iPXEService)

	if err := iPXEService.Extract(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

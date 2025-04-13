package grpcserver

import (
	"google.golang.org/grpc"

	v1 "github.com/sudeeya/avito-assignment/internal/controller/grpc/v1"
)

func NewServer(serviceServer v1.PVZServiceServer) *grpc.Server {
	server := grpc.NewServer()

	v1.RegisterPVZServiceServer(server, serviceServer)

	return server
}

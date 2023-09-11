package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/edmarfelipe/grpc-load-balancing/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func main() {
	slog.Info("Starting server...")
	srv := grpc.NewServer()
	shared.RegisterUserServer(srv, &UserServer{})

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(fmt.Errorf("could not connect: %w", err))
	}
	slog.Info("Server is running on port 8001")
	srv.Serve(listener)
}

type UserServer struct {
	shared.UnimplementedUserServer
}

func (us *UserServer) Hello(ctx context.Context, req *shared.Request) (*shared.Reply, error) {
	slog.Info("Recive Hello with ID: " + req.Id)
	p, _ := peer.FromContext(ctx)

	return &shared.Reply{
		Message: "Hello from: " + p.Addr.String(),
	}, nil
}

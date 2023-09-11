package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/edmarfelipe/grpc-load-balancing/shared"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverPord := os.Getenv("SERVER_HOST")
	slog.Info("Starting client", "serverPort", serverPord)

	cc, err := startClient(serverPord)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling /hello")

		usrService := shared.NewUserClient(cc)
		result, err := usrService.Hello(context.Background(), &shared.Request{Id: uuid.NewString()})
		if err != nil {
			slog.Error("Error while calling Hello RPC", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Reply from the server: "+result.Message)
	})
	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		slog.Error("Error while starting server", "err", err)
	}
}

func startClient(serverPort string) (*grpc.ClientConn, error) {
	cc, err := grpc.Dial(
		serverPort,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	return cc, nil
}

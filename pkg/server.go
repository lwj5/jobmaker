package main

import (
	"log"
	"net"

	"github.com/lwj5/jobmaker/pkg/jobmaker"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement the service.
type server struct {
	jobmaker.UnimplementedJobmakerServer
}

func main() {
	log.Printf("Starting server...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Started server on port %v.", port)

	s := grpc.NewServer()
	jobmaker.RegisterJobmakerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

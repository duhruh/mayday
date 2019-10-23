package main

import (
	"log"
	"net"

	"github.com/docker/mayday/pkg/mayday"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	mayday.NewServer(s)
	println("server listening")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

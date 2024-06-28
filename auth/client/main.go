package main

import (
	"context"
	"log"
	"time"

	pb "github.com/iamstep4ik/quick-meet/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRegisterUserClient(conn)

	username := "testuser"
	email := "testuser@example.com"
	password := "P@ssw0rd!"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Register(ctx, &pb.RegisterRequest{Username: username, Email: email, Password: password})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Register Response: %s", r.GetMessage())
}

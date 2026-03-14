package main

import (
	"context"
	"log"

	"github.com/sjingzhao/gong/pb/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, err := grpc.NewClient(":8091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	// 获取 UserServiceClient
	userClient := example.NewUserServiceClient(client)

	// 调用 unary RPC
	resp, err := userClient.GetUser(context.Background(), &example.UserRequest{UserId: 42})
	if err != nil {
		log.Fatal("GetUser error:", err)
	}
	log.Printf("Unary response: %+v", resp)

	// 调用 streaming RPC
	stream, err := userClient.AddUser(context.Background(), &example.User{Id: 1, Name: "x", Email: "e"})
	if err != nil {
		log.Fatal("StreamUsers error:", err)
	}

	log.Printf("Stream message: %+v", stream)
}

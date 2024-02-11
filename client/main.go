package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jittakal/todo-api-hexagonal/api/proto"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create a TaskService client
	client := proto.NewTaskServiceClient(conn)

	// Define a context with timeout
	ctx := context.Background()

	// Call the GetTask RPC
	getTaskResponse, err := client.GetTask(ctx, &proto.GetTaskRequest{Id: "77c7a363-1893-4cd7-beb1-5b5bdca53723"})
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	// Print the task details
	fmt.Printf("Task ID: %s\n", getTaskResponse.Id)
	fmt.Printf("Task Title: %s\n", getTaskResponse.Title)
	fmt.Printf("Task Done: %t\n", getTaskResponse.Done)
}

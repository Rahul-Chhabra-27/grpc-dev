package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	calculatorProto "rahulchhabra.io/proto/calculator"
)

func unaryRPC(grpcClient calculatorProto.CalculatorServiceClient) {
	fmt.Println("Calling the Calculator RPC")
	request := &calculatorProto.SumRequest{
		First:  10,
		Second: 12,
	}
	response, err := grpcClient.Calculator(context.Background(), request)

	if err != nil {
		log.Fatalf("Error wile calling calculator RPC %v \n", err)
	}

	fmt.Printf("Addition of two numbers is %d", response.SumResult)
}
func main() {

	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't not connect: %v", err)
	}
	defer connection.Close()

	grpcClient := calculatorProto.NewCalculatorServiceClient(connection)
	// unary function ...
	unaryRPC(grpcClient)
}

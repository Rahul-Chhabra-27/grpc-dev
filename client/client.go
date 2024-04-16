package main

import (
	"context"
	"fmt"
	"io"
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
func ServerStreamingRPC(client calculatorProto.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming..")

	// request body..
	request := &calculatorProto.PrimeDecompositionRequest{
		PrimeNumber: int64(36),
	}

	// calling the PrimeNumberDecomposition RPC...
	stream, err := client.PrimeNumberDecomposition(context.Background(), request)

	if err != nil {
		log.Fatalf("Something happend %v\n", err)
	}
	// Recieving the server response in the form of stream..
	for {
		// recieving the chunk of stream.
		res, err := stream.Recv()

		// stream completed.
		if err == io.EOF {
			fmt.Println("Steam Completed")
			break
		} else if err != nil {
			log.Fatalf("Something happend %v \n", err)
		} else {
			fmt.Printf("Factors are : %d\n", res.PrimeFactor)
		}
	}
}
func main() {

	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't not connect: %v", err)
	}
	defer connection.Close()

	grpcClient := calculatorProto.NewCalculatorServiceClient(connection)
	// unary function ...
//	unaryRPC(grpcClient)

	ServerStreamingRPC(grpcClient);
}

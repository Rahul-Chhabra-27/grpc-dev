package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	calculatorProto "rahulchhabra.io/proto/calculator"
)

// Unary RPC function for the client..
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
// Client Streaming RPC function for the client..
func ClientStreamingRPC(client calculatorProto.CalculatorServiceClient) {

	// request slice
	requests := []*calculatorProto.SumOfTheArrayElementsRequest{
		&calculatorProto.SumOfTheArrayElementsRequest{
			Element: 1,
		},
		&calculatorProto.SumOfTheArrayElementsRequest{
			Element: 2,
		},
		&calculatorProto.SumOfTheArrayElementsRequest{
			Element: 3,
		},
		&calculatorProto.SumOfTheArrayElementsRequest{
			Element: 4,
		},
		&calculatorProto.SumOfTheArrayElementsRequest{
			Element: 5,
		},
	}

	// calling the sumofTheArrayElements RPC.
	stream, err := client.SumOfTheArrayElements(context.Background())
	if err != nil {
		log.Fatalf("Error while calling SumOfTheArrayElements RPC.. %v", err)
	}

	for _, request := range requests {
		fmt.Printf("Sending request %v \n", request)
		stream.Send(request)
		// waiting for one second to see the stream in the o/p console one by one..
		time.Sleep(time.Second)
	}

	response, err := stream.CloseAndRecv();
	if err != nil {
		log.Fatalf("Error while recieving the response .. %v", err)
	}
	fmt.Printf("Response : %v\n", response.Sumofallelements)
}
func BiDirectionalStreamingRPC(client calculatorProto.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDirectional streaming..")
	// request slice
	requests := []*calculatorProto.FindMaximumRequest{
		&calculatorProto.FindMaximumRequest{
			Number: 1,
		},
		&calculatorProto.FindMaximumRequest{
			Number: 5,
		},
		&calculatorProto.FindMaximumRequest{
			Number: 3,
		},
		&calculatorProto.FindMaximumRequest{
			Number: 6,
		},
		&calculatorProto.FindMaximumRequest{
			Number: 2,
		},
	}
	stream, err := client.FindMaximum(context.Background());
	if err != nil {
		log.Fatalf("Error while calling FindMaximum RPC %v", err)
	}

	// channel to wait for the goroutine to complete..
	watch := make(chan struct{})
	// goroutine to send the stream..
	go func () {
		for _, request := range requests {
			fmt.Printf("Sending request %v\n", request)
			stream.Send(request)
			time.Sleep(time.Millisecond * 1000)
		}
		// closing the stream after sending all the requests..
		stream.CloseSend()
	}()
	// goroutine to recieve the stream..
	go func () {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving the stream %v\n", err)
				break
			}
			fmt.Printf("Recieved : %v\n", response.Maximum)
		}
		close(watch)
	}()
		<-watch
}
func main() {

	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't not connect: %v", err)
	}
	defer connection.Close()

	grpcClient := calculatorProto.NewCalculatorServiceClient(connection)
	//1. unary function ...
	//	unaryRPC(grpcClient)
	//2. server streaming function..
	//ServerStreamingRPC(grpcClient)

	// 3. client streaing function.
	//ClientStreamingRPC(grpcClient)

	// 4. Bi-Directional streaming function.
	BiDirectionalStreamingRPC(grpcClient)
}

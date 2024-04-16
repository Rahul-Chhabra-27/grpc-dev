package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	calculatorProto "rahulchhabra.io/proto/calculator"
)

type Config struct {
	calculatorProto.UnimplementedCalculatorServiceServer
}

// Unary RPC function
func (*Config) Calculator(ctx context.Context, request *calculatorProto.SumRequest) (response *calculatorProto.SumResponse, err error) {
	// getting the first number from the request body.
	firstNumber := request.First
	// getting the second number from the request body.
	secondNumber := request.Second

	// returning the SumResponse ;
	return &calculatorProto.SumResponse{
		SumResult: firstNumber + secondNumber,
	}, nil
}
func (*Config) PrimeNumberDecomposition(request *calculatorProto.PrimeDecompositionRequest, stream calculatorProto.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Println("Server Streaming Started...")
	var factor int64 = 2
	var number int64 = request.PrimeNumber
	for {
		if number <= 1 {
			break
		} else if number%factor == 0 {
			stream.Send(&calculatorProto.PrimeDecompositionResponse{
				PrimeFactor: factor,
			})
			number /= factor
		} else {
			fmt.Println("Factor is incremented by one++ ")
			factor++
		}
	}
	return nil
}
func main() {
	// listen on the port
	listner, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to start the server %v\n", err)
	}
	// create a new gRPC server...
	grpcServer := grpc.NewServer()
	// register the calculator service.
	calculatorProto.RegisterCalculatorServiceServer(grpcServer, &Config{})
	log.Printf("Starting gRPC listener on port " + "localhost:50051")
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	calculatorProto "rahulchhabra.io/proto/calculator"
)

type Config struct {
	calculatorProto.UnimplementedCalculatorServiceServer
}

//  Unary RPC function
func (*Config) Calculator(ctx context.Context, request *calculatorProto.SumRequest) (response *calculatorProto.SumResponse, err error) {
	// getting the first number from the request body.
	firstNumber := request.First;
	// getting the second number from the request body.
	secondNumber := request.Second

	// returning the SumResponse ;
	return &calculatorProto.SumResponse{
		SumResult: firstNumber + secondNumber,
	}, nil;
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

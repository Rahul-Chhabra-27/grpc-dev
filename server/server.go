package main

import (
	"context"
	"fmt"
	"io"
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
// Server Streaming..
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
// Client Streaming..
func (*Config) SumOfTheArrayElements(stream calculatorProto.CalculatorService_SumOfTheArrayElementsServer) error {
	sumResult := int64(0)
	// recieving the stream...
	for {
		// recieving the chunk of stream..
		chunk, err := stream.Recv()

		// stream completd...
		if err == io.EOF {
			response := &calculatorProto.SumOfTheArrayElementsResponse{
				Sumofallelements: int64(sumResult),
			}
			// sending the response and closing the stream..
			return stream.SendAndClose(response);

		} else if err != nil {
			//some error occured...
			log.Fatalf("Error while Reading client streaming %v\n",err)

		} else {
			// add the element coming from stream to the sumResult variable.
			sumResult += int64(chunk.Element);
		} 
	}
}
// Bi-Directional Streaming..

// FindMaximum function will find the maximum number from the stream of numbers.
// and send the maximum number to the client.
func (*Config) FindMaximum(stream calculatorProto.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request..")
	maximum := int64(0)
	for {
		// recieving the chunk of stream.
		chunk, err := stream.Recv()
		// stream completed...
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading the client stream %v\n",err)
			return err
		} else {
			number := chunk.Number
			if number > maximum {
				maximum = number
				// sending the maximum number to the client.
				sendingError := stream.Send(&calculatorProto.FindMaximumResponse{
					Maximum: maximum,
				})

				if sendingError != nil {
					log.Fatalf("Error while sending the maximum number to the client %v\n",sendingError)
					return sendingError
				}
			}
		}
	}
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

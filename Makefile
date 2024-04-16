all:
	protoc -I ./proto \
	--go_out ./proto --go_opt paths=source_relative \
	--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
	./proto/calculator/calculator.proto    

grpc-server:
		go run server/server.go


grpc-client:
		go run client/client.go

path-variable:
	 export PATH="$PATH:$(go env GOPATH)/bin";
package grpc_application

import (
	"context"
	example_proto_messages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
)

func (ls *GrpcServer) HelloWorld(ctx context.Context, in *example_proto_messages.HelloRequest) (*example_proto_messages.HelloResponse, error) {
	return &example_proto_messages.HelloResponse{User: &example_proto_messages.UserData{Name: "Robson"}}, nil
}

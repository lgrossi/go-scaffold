package grpc_application

import (
	"database/sql"
	"github.com/lgrossi/go-scaffold/src/configs"
	"github.com/lgrossi/go-scaffold/src/database"
	example_proto_messages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"github.com/lgrossi/go-scaffold/src/network"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer struct {
	DB *sql.DB
	example_proto_messages.ExampleServiceServer
	network.ServerInterface
}

func (ls *GrpcServer) Initialize(gConfigs configs.GlobalConfigs) error {
	ls.DB = database.PullConnection(gConfigs)
	return nil
}

func (ls *GrpcServer) Run(gConfigs configs.GlobalConfigs) error {
	c, err := net.Listen("tcp", gConfigs.ServerConfigs.Grpc.Format())

	if err != nil {
		return err
	}

	server := grpc.NewServer()
	example_proto_messages.RegisterExampleServiceServer(server, ls)

	return server.Serve(c)
}

func (ls *GrpcServer) GetName() string {
	return "gRPC"
}

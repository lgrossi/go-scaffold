package grpc_application

import (
	"database/sql"
	example_proto_messages "github.com/lgrossi/go-scaffold/src/grpc/example_proto_defs"
	"github.com/lgrossi/go-scaffold/src/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrpcServer_GetName(t *testing.T) {
	type fields struct {
		DB                 *sql.DB
		LoginServiceServer example_proto_messages.ExampleServiceServer
		ServerInterface    network.ServerInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		"",
		fields{},
		"gRPC",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &GrpcServer{
				DB:                 tt.fields.DB,
				ExampleServiceServer: tt.fields.LoginServiceServer,
				ServerInterface:    tt.fields.ServerInterface,
			}
			assert.Equal(t, tt.want, ls.GetName())
		})
	}
}

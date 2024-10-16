package grpc_application

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	example_proto_messages "github.com/tibia-oce/discord-bot/src/grpc/example_proto_defs"
	"github.com/tibia-oce/discord-bot/src/network"
	"testing"
)

func TestGrpcServer_GetName(t *testing.T) {
	type fields struct {
		DB                   *sql.DB
		exampleServiceServer example_proto_messages.ExampleServiceServer
		ServerInterface      network.ServerInterface
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
				DB:                   tt.fields.DB,
				ExampleServiceServer: tt.fields.exampleServiceServer,
				ServerInterface:      tt.fields.ServerInterface,
			}
			assert.Equal(t, tt.want, ls.GetName())
		})
	}
}

package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tibia-oce/discord-bot/src/network"
	"google.golang.org/grpc"
	"testing"
)

func TestApi_GetName(t *testing.T) {
	type fields struct {
		Router          *gin.Engine
		DB              *sql.DB
		GrpcConnection  *grpc.ClientConn
		ServerInterface network.ServerInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{{
		"",
		fields{},
		"api",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_api := &Api{
				Router:          tt.fields.Router,
				DB:              tt.fields.DB,
				GrpcConnection:  tt.fields.GrpcConnection,
				ServerInterface: tt.fields.ServerInterface,
			}
			assert.Equal(t, tt.want, _api.GetName())
		})
	}
}

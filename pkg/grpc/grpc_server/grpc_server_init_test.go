package grpc_server

import (
	"context"
	"testing"
)

func TestNewGrpcServer(t *testing.T) {
	stopCh := make(chan struct{})
	server := NewGrpcServer(9999, context.TODO(), stopCh)
	server.StartGrpcServerTcp()
}

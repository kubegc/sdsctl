package grpc_server

import (
	"context"
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/grpc/pb_gen"
	"github.com/WANNA959/sdsctl/pkg/internal"
	"github.com/WANNA959/sdsctl/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type GrpcServer struct {
	*pb_gen.UnimplementedSdsCtlServiceServer
	ctx     context.Context
	stopCh  chan struct{}
	port    int
	service *internal.NetworkControllerService
}

var logger = utils.GetLogger()
var gServer *GrpcServer

func GetGServer() *GrpcServer {
	return gServer
}

func NewGrpcServer(port int, ctx context.Context, stopCh chan struct{}) *GrpcServer {
	s := &GrpcServer{
		ctx:    ctx,
		stopCh: stopCh,
		port:   port,
	}

	s.service = internal.NewLiteNCService()
	return s
}

func (s *GrpcServer) StartGrpcServerTcp() error {
	defer logger.Debug("StartGrpcServerTcp done")

	tcpAddr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", tcpAddr)
	defer lis.Close()
	if err != nil {
		logger.Errorf("tcp failed to listen: %v", err)
		return err
	}

	gopts := []grpc.ServerOption{}
	server := grpc.NewServer(gopts...)
	// register reflection for grpcurl service
	reflection.Register(server)
	// register service
	pb_gen.RegisterSdsCtlServiceServer(server, s)
	logger.Infof("grpc server ready to serve at %+v", tcpAddr)

	go func() {
		for {
			select {
			case <-s.stopCh:
				server.GracefulStop()
				return
			}
		}
	}()

	if err := server.Serve(lis); err != nil {
		logger.Errorf("grpc server failed to serve: %v", err)
		return err
	}
	return nil
}

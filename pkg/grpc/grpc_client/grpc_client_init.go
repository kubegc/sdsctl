package grpc_client

import (
	"errors"
	"fmt"
	"github.com/WANNA959/sdsctl/pkg/grpc/pb_gen"
	"github.com/WANNA959/sdsctl/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	C    pb_gen.SdsCtlServiceClient
	Ip   string
	Port string
}

var logger = utils.GetLogger()

func NewGrpcClient(ip, port string) (*GrpcClient, error) {
	client := &GrpcClient{
		C:    nil,
		Ip:   ip,
		Port: port,
	}
	err := client.InitGrpcClientConn()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *GrpcClient) InitGrpcClientConn() error {
	// Set up a connection to the server.
	var address string
	if len(c.Ip) == 0 || len(c.Port) == 0 {
		logger.Error("ip and port can't be empty")
		return errors.New("ip and port can't be empty")
	}
	address = fmt.Sprintf("%s:%s", c.Ip, c.Port)

	var dialOpt []grpc.DialOption
	insecure.NewCredentials()
	dialOpt = append(dialOpt, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// init grpc client
	conn, err := grpc.Dial(address, dialOpt...)
	if err != nil {
		logger.Errorf("can't connect: %v", err)
		return err
	}
	c.C = pb_gen.NewSdsCtlServiceClient(conn)
	return nil
}

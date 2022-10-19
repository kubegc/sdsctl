package internal

import "github.com/op/go-logging"

type NetworkControllerService struct {
}

var logger *logging.Logger

func NewLiteNCService() *NetworkControllerService {
	return &NetworkControllerService{}
}

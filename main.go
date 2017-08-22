package gemstone

import (
	"github.com/cryptopay-dev/gemstone/logger"
	"google.golang.org/grpc"
)

type Service interface {
	Run() error
	Server() *grpc.Server
	Client(string) (*grpc.ClientConn, error)
	Logger() logger.Logger
}

func NewService(opts ...Option) (Service, error) {
	return newService(opts...)
}

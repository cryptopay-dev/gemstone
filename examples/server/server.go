package main

import (
	"fmt"
	"os"

	"github.com/cryptopay-dev/gemstone"
	proto "github.com/cryptopay-dev/gemstone/examples/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
)

type Greeter struct {
	service gemstone.Service
}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest) (res *proto.HelloResponse, err error) {
	g.service.Logger().Infof("We have new greeting request from %s", req.Name)

	return &proto.HelloResponse{Greeting: "Hello, " + req.Name}, nil
}

func main() {
	service, err := gemstone.NewService(
		gemstone.Version("1.0.0"),
		gemstone.Name("greeter"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	greeter := &Greeter{
		service: service,
	}
	proto.RegisterGreeterServer(service.Server(), greeter)
	reflection.Register(service.Server())

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

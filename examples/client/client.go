package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cryptopay-dev/gemstone"
	proto "github.com/cryptopay-dev/gemstone/examples/proto"
	"golang.org/x/net/context"
)

func main() {
	// Create a new service. Optionally include some options here.
	service, err := gemstone.NewService()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create new greeter client
	client, err := service.Client("greeter")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	greeter := proto.NewGreeterClient(client)

	// Call the greeter
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*2)
	rsp, err := greeter.Hello(ctx, &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}

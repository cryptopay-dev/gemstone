package gemstone

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptopay-dev/gemstone/internal"
	"github.com/cryptopay-dev/gemstone/logger"
	"github.com/cryptopay-dev/gemstone/registry"
	"github.com/cryptopay-dev/gemstone/registry/consul"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type DefaultService struct {
	server  *grpc.Server
	service registry.Service
	options Options
}

func newService(opts ...Option) (*DefaultService, error) {
	options := newOptions(opts...)

	// Obtaining current registry
	if options.Registry == nil {
		reg, err := consul.New()
		if err != nil {
			return nil, err
		}

		options.Registry = reg
	}

	return &DefaultService{
		server:  grpc.NewServer(),
		options: options,
	}, nil
}

func (s *DefaultService) Logger() logger.Logger {
	return s.options.Logger
}

func (s *DefaultService) Use() error {
	return nil
}

func (s *DefaultService) Run() error {
	// Getting first opened port and run on it
	addr, err := net.ResolveTCPAddr("tcp", s.options.Address)
	if err != nil {
		return err
	}

	// Staring default listener
	listener, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}

	defer func() {
		s.server.GracefulStop()
		listener.Close()
		s.options.Logger.Info("Closing server listener")
	}()

	// Default channel for global stop
	stop := make(chan struct{}, 1)

	// Listening
	go func() {
		s.options.Logger.Infof("Starting listen on %s", listener.Addr().String())

		if serveErr := s.server.Serve(listener); err != nil {
			s.options.Logger.Errorf("Error while trying to Serve: %v", serveErr)
			stop <- struct{}{}
		}
	}()

	// <nil> - is not valid address
	var ip = addr.IP.String()
	if len(addr.IP) == 0 {
		ip = ""
	}

	if ip, err = internal.Extract(ip); err != nil {
		return err
	}

	// Registering service in registry
	sid := s.options.Name + "-" + uuid.NewV4().String()
	s.service = registry.Service{
		ID:      sid,
		Name:    s.options.Name,
		Version: s.options.Version,
		Addr:    ip,
		Port:    addr.Port,
	}
	if err := s.register(); err != nil {
		return err
	}
	s.options.Logger.Infof("Registered in registry with name %s", s.service.ID)

	// Running heartbeat
	go s.run(stop)

	// Catching sigterm and process them
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)

		sig := <-ch
		s.options.Logger.Infof("Received signal %s", sig)

		stop <- struct{}{}
	}()

	<-stop
	if err := s.options.Registry.Deregister(s.service); err != nil {
		return err
	}

	return nil
}

func (s *DefaultService) register() error {
	if err := s.options.Registry.Register(s.service); err != nil {
		return err
	}

	s.options.Logger.Debugf("Updated registry record")
	return nil
}

func (s *DefaultService) run(exit chan struct{}) {
	t := time.NewTicker(time.Second * 15)

	for {
		select {
		case <-t.C:
			s.register()
		case <-exit:
			t.Stop()
			return
		}
	}
}

func (s *DefaultService) Server() *grpc.Server {
	return s.server
}

func (s *DefaultService) Client(name string) (*grpc.ClientConn, error) {
	services, err := s.options.Registry.GetService(name)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, errors.New("cannot find any service")
	}

	rr := RoundRobin(services)
	service, err := rr()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", service.Addr, service.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

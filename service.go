package gemstone

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptopay-dev/gemstone/logger"
	"github.com/cryptopay-dev/gemstone/registry"
	"github.com/cryptopay-dev/gemstone/registry/consul"
	uuid "github.com/satori/go.uuid"
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
		registry, err := consul.New()
		if err != nil {
			return nil, err
		}

		options.Registry = registry
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
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return err
	}

	// Staring default listener
	listener, err := net.Listen("tcp", addr.String())
	defer listener.Close()
	if err != nil {
		return err
	}

	// Listening
	go s.server.Serve(listener)

	// Registering service in registry
	sid := s.options.Name + "-" + uuid.NewV4().String()
	s.service = registry.Service{
		ID:      sid,
		Name:    s.options.Name,
		Version: s.options.Version,
		Addr:    listener.Addr().String(),
	}
	if err := s.register(); err != nil {
		return err
	}
	s.options.Logger.Infof("Registered in registry with name %s", s.service.ID)

	// Running hearthbeat
	ex := make(chan struct{}, 1)
	go s.run(ex)

	// Catching sigterm and process them
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	s.options.Logger.Infof("Received signal %s", <-ch)
	ex <- struct{}{}

	// Deregistering on SIGTERM
	if err := s.options.Registry.Deregister(s.service); err != nil {
		return err
	}

	s.options.Logger.Info("Closing server listener")
	return nil
}

func (s *DefaultService) register() error {
	if err := s.options.Registry.Register(s.service); err != nil {
		return err
	}

	s.options.Logger.Debugf("Updated registery record")
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
		return nil, errors.New("Cannot find any service")
	}

	rr := RoundRobin(services)
	service, err := rr()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(service.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

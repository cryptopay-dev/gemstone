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
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
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

		if err := s.server.Serve(listener); err != nil {
			s.options.Logger.Errorf("Error while trying to Serve: %v", err)
			stop <- struct{}{}
		}
	}()

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

package gemstone

import (
	"github.com/cryptopay-dev/gemstone/logger"
	"github.com/cryptopay-dev/gemstone/logger/zap"
	"github.com/cryptopay-dev/gemstone/registry"
)

type Options struct {
	Logger   logger.Logger
	Registry registry.Registry
	Address  string
	Name     string
	Version  string
}

type Option func(*Options)

var (
	DefaultServiceName = "microservice"
	DefaultVersion     = "0.0.0"
	DefaultAddress     = "localhost:0"
	DefaultLogger      = zap.New()
)

func newOptions(opts ...Option) Options {
	opt := Options{
		Address: DefaultAddress,
		Logger:  DefaultLogger,
		Name:    DefaultServiceName,
		Version: DefaultVersion,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func Logger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func Registry(registry registry.Registry) Option {
	return func(o *Options) {
		o.Registry = registry
	}
}

func Address(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

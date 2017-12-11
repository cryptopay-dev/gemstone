package consul

import (
	"fmt"
	"sync"
	"time"

	"github.com/cryptopay-dev/gemstone/registry"
	"github.com/hashicorp/consul/api"
	"github.com/mitchellh/hashstructure"
)

type Registry struct {
	client *api.Client

	registry map[string]uint64
	sync.Mutex
}

func New() (registry.Registry, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{
		client:   client,
		registry: make(map[string]uint64),
	}, nil
}

func (r *Registry) Register(s registry.Service) error {
	h, err := hashstructure.Hash(s, nil)
	if err != nil {
		return err
	}

	r.Lock()
	v, ok := r.registry[s.Name]
	r.Unlock()

	if ok && v == h {
		if err := r.client.Agent().PassTTL("service:"+s.ID, ""); err == nil {
			return nil
		}
	}

	// TTL of service
	splay := time.Second * 5
	deregTTL := time.Minute + splay

	check := &api.AgentServiceCheck{
		TTL: fmt.Sprintf("%v", time.Minute),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%v", deregTTL),
	}

	if err := r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      s.ID,
		Name:    s.Name,
		Port:    s.Port,
		Address: s.Addr,
		Tags:    []string{s.Name},
		Check:   check,
	}); err != nil {
		return err
	}

	return r.client.Agent().PassTTL("service:"+s.ID, "")
}

func (r *Registry) Deregister(s registry.Service) error {
	r.Lock()
	delete(r.registry, s.Name)
	r.Unlock()

	return r.client.Agent().ServiceDeregister(s.ID)
}

func (r *Registry) GetService(name string) ([]*registry.Service, error) {
	res, _, err := r.client.Health().Service(name, "", false, nil)
	if err != nil {
		return nil, err
	}

	var services []*registry.Service
	for _, s := range res {
		if s.Service.Service != name {
			continue
		}

		if s.Checks.AggregatedStatus() == api.HealthCritical {
			continue
		}

		services = append(services, &registry.Service{
			ID:   s.Service.ID,
			Name: s.Service.Service,
			Port: s.Service.Port,
			Addr: s.Service.Address,
		})
	}

	return services, nil
}

func (r *Registry) List() ([]*registry.Service, error) {
	res, _, err := r.client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	var services []*registry.Service
	for service := range res {
		services = append(services, &registry.Service{
			Name: service,
		})
	}

	return services, nil
}

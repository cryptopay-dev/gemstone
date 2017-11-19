package gemstone

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/cryptopay-dev/gemstone/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Next func() (*registry.Service, error)

func RoundRobin(services []*registry.Service) Next {
	var i = rand.Int()
	var mtx sync.Mutex

	return func() (*registry.Service, error) {
		if len(services) == 0 {
			return nil, errors.New("no services available")
		}

		mtx.Lock()
		service := services[i%len(services)]
		i++
		mtx.Unlock()

		return service, nil
	}
}

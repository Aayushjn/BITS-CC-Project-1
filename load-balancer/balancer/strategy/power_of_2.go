package strategy

import (
	"math/rand"
	"net/http"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type PowerOfTwoStrategy struct {
	backends []*backend.Backend
}

func (po2 *PowerOfTwoStrategy) nextIndex(numBackends int) int {
	return rand.Intn(numBackends)
}

func (po2 *PowerOfTwoStrategy) Backends() []*backend.Backend {
	return po2.backends
}

func (po2 *PowerOfTwoStrategy) Next(req *http.Request) *backend.Backend {
	numBackends := len(po2.backends)
	if numBackends == 0 {
		return nil
	}

	for i := 1; i <= numBackends; i++ {
		first := po2.nextIndex(numBackends)
		second := po2.nextIndex(numBackends)
		for second == first {
			second = po2.nextIndex(numBackends)
		}

		backend1 := po2.backends[first]
		backend2 := po2.backends[second]

		if !backend1.IsAlive() && !backend2.IsAlive() {
			continue
		}

		if backend1.IsAlive() {
			if backend2.IsAlive() {
				if backend1.GetInFlightRequests() <= backend2.GetInFlightRequests() {
					return backend1
				}
				return backend2
			}
			return backend1
		}
		return backend2
	}
	return nil
}

func (po2 *PowerOfTwoStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(po2.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	po2.backends = append(po2.backends, b)
	return nil
}

func (po2 *PowerOfTwoStrategy) Unregister(backendUrl string) error {
	if len(po2.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range po2.backends {
		if b.URL.String() == backendUrl {
			po2.backends = append(po2.backends[:i], po2.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewPowerOfTwoStrategy() *PowerOfTwoStrategy {
	return &PowerOfTwoStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
	}
}

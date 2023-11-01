package strategy

import (
	"math/rand"
	"net/http"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type RandomStrategy struct {
	backends []*backend.Backend
}

func (r *RandomStrategy) nextIndex(numBackends int) int {
	return rand.Intn(numBackends)
}

func (r *RandomStrategy) Backends() []*backend.Backend {
	return r.backends
}

func (r *RandomStrategy) Next(req *http.Request) *backend.Backend {
	numBackends := len(r.backends)
	if numBackends == 0 {
		return nil
	}

	for i := 1; i <= numBackends; i++ {
		next := r.nextIndex(numBackends)
		if r.backends[next].IsAlive() {
			return r.backends[next]
		}
	}
	return nil
}

func (r *RandomStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(r.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	r.backends = append(r.backends, b)
	return nil
}

func (r *RandomStrategy) Unregister(backendUrl string) error {
	if len(r.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range r.backends {
		if b.URL.String() == backendUrl {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
	}
}

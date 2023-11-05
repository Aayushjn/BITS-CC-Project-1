package strategy

import (
	"net/http"
	"sync/atomic"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type RoundRobinStrategy struct {
	backends []*backend.Backend
	current  atomic.Uint64
}

func (rr *RoundRobinStrategy) nextIndex(numBackends int) int {
	return int(rr.current.Add(uint64(1)) % uint64(numBackends))
}

func (rr *RoundRobinStrategy) Backends() []*backend.Backend {
	return rr.backends
}

func (rr *RoundRobinStrategy) Next(req *http.Request) *backend.Backend {
	numBackends := len(rr.backends)
	if numBackends == 0 {
		return nil
	}

	next := rr.nextIndex(numBackends)
	for i := next; i < numBackends+next; i++ {
		idx := i % numBackends
		if rr.backends[idx].IsAlive() {
			if i != next {
				rr.current.Store(uint64(idx))
			}
			return rr.backends[idx]
		}
	}
	return nil
}

func (rr *RoundRobinStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(rr.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	rr.backends = append(rr.backends, b)
	return nil
}

func (rr *RoundRobinStrategy) Unregister(backendUrl string) error {
	if len(rr.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range rr.backends {
		if b.URL.String() == backendUrl {
			rr.backends = append(rr.backends[:i], rr.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
	}
}

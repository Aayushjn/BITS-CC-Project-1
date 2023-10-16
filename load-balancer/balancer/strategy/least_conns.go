package strategy

import (
	"sync/atomic"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type LeastConnsStrategy struct {
	backends []*backend.Backend
	current  atomic.Uint64
}

func (lc *LeastConnsStrategy) Backends() []*backend.Backend {
	return lc.backends
}

func (lc *LeastConnsStrategy) Next() *backend.Backend {
	if len(lc.backends) == 0 {
		return nil
	}

	leastConnIdx := 0
	for i, backend := range lc.backends {
		if i == int(lc.current.Load()) {
			continue
		}
		if backend.IsAlive() && backend.GetInFlightRequests() <= lc.backends[leastConnIdx].GetInFlightRequests() {
			leastConnIdx = i
		}
	}
	lc.current.Store(uint64(leastConnIdx))
	return lc.backends[leastConnIdx]
}

func (lc *LeastConnsStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(lc.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	lc.backends = append(lc.backends, b)
	return nil
}

func (lc *LeastConnsStrategy) Unregister(backendUrl string) error {
	if len(lc.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range lc.backends {
		if b.URL.String() == backendUrl {
			lc.backends = append(lc.backends[:i], lc.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewLeastConnsStrategy() *LeastConnsStrategy {
	return &LeastConnsStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
	}
}

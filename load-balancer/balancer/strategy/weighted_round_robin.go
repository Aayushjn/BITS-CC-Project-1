package strategy

import (
	"net/http"
	"sync/atomic"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type WeightedRoundRobinStrategy struct {
	backends []*backend.Backend
	weights  cmap.ConcurrentMap[string, int64]
	current  atomic.Uint64
	counter  atomic.Uint64
}

func (wrr *WeightedRoundRobinStrategy) nextIndex(numBackends int) int {
	currentBackend := wrr.backends[wrr.current.Load()]
	if wt, ok := wrr.weights.Get(currentBackend.URL.String()); ok && wrr.counter.Load() < uint64(wt) {
		return int(wrr.current.Load())
	}
	wrr.counter.Store(0)
	return int(wrr.current.Add(uint64(1)) % uint64(numBackends))
}

func (wrr *WeightedRoundRobinStrategy) Backends() []*backend.Backend {
	return wrr.backends
}

func (wrr *WeightedRoundRobinStrategy) Next(req *http.Request) *backend.Backend {
	numBackends := len(wrr.backends)
	if numBackends == 0 {
		return nil
	}

	next := wrr.nextIndex(numBackends)
	for i := next; i < numBackends+next; i++ {
		idx := i % numBackends
		if wrr.backends[idx].IsAlive() {
			if i != next {
				wrr.current.Store(uint64(idx))
			}
			wrr.counter.Add(1)
			return wrr.backends[idx]
		}
	}
	return nil
}

func (wrr *WeightedRoundRobinStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(wrr.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	wrr.backends = append(wrr.backends, b)
	wrr.weights.Set(b.URL.String(), int64(params["weight"].(float64)))

	return nil
}

func (lc *WeightedRoundRobinStrategy) Unregister(backendUrl string) error {
	if len(lc.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range lc.backends {
		if b.URL.String() == backendUrl {
			lc.weights.Remove(backendUrl)
			lc.backends = append(lc.backends[:i], lc.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewWeightedRoundRobinStrategy() *WeightedRoundRobinStrategy {
	return &WeightedRoundRobinStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
		weights:  cmap.New[int64](),
	}
}

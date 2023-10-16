package strategy

import (
	"sync/atomic"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type WeightedRoundRobinStrategy struct {
	backends      []*backend.Backend
	weights       cmap.ConcurrentMap[string, int64]
	current       atomic.Uint64
	maxWeight     atomic.Int64
	currentGcd    atomic.Int64
	currentWeight atomic.Int64
}

func (wrr *WeightedRoundRobinStrategy) nextIndex(numBackends int) int {
	return int(wrr.current.Add(uint64(1)) % uint64(numBackends))
}

func (wrr *WeightedRoundRobinStrategy) Backends() []*backend.Backend {
	return wrr.backends
}

func (wrr *WeightedRoundRobinStrategy) Next() *backend.Backend {
	numBackends := len(wrr.backends)
	if numBackends == 0 {
		return nil
	}

	for i := 1; i <= numBackends; i++ {
		next := wrr.nextIndex(numBackends)
		if next == 0 {
			wrr.currentWeight.Add(^wrr.currentGcd.Load())
			if wrr.currentWeight.Load() <= 0 {
				wrr.currentWeight.Store(wrr.maxWeight.Load())
			}
		}

		nextBackend := wrr.backends[next]
		if !nextBackend.IsAlive() {
			continue
		}
		nextBackendWeight, ok := wrr.weights.Get(nextBackend.URL.String())
		if !ok {
			nextBackendWeight = 1
		}
		if nextBackendWeight >= wrr.currentWeight.Load() {
			return wrr.backends[next]
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
	backendWeight, ok := wrr.weights.Get(b.URL.String())
	if !ok {
		backendWeight = 1
	}
	if wrr.currentGcd.Load() == 0 {
		wrr.currentGcd.Store(backendWeight)
		wrr.maxWeight.Store(backendWeight)
	} else {
		wrr.currentGcd.Store(util.Gcd(wrr.currentGcd.Load(), backendWeight))
		if wrr.maxWeight.Load() < backendWeight {
			wrr.maxWeight.Store(backendWeight)
		}
	}
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

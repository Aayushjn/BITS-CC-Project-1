package strategy

import (
	"hash"
	"hash/fnv"
	"net/http"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type IpHashStrategy struct {
	backends []*backend.Backend
	hasher   hash.Hash32
}

func (ih *IpHashStrategy) Backends() []*backend.Backend {
	return ih.backends
}

func (ih *IpHashStrategy) Next(req *http.Request) *backend.Backend {
	numBackends := len(ih.backends)
	if numBackends == 0 {
		return nil
	}

	ih.hasher.Reset()
	ih.hasher.Write([]byte(req.RemoteAddr))
	idx := ih.hasher.Sum32() % uint32(numBackends)
	return ih.backends[idx]
}

func (ih *IpHashStrategy) Register(b *backend.Backend, params map[string]any) error {
	if len(ih.backends) == util.MaxBackends {
		return &errors.ErrBackendLimitOverflow{}
	}
	ih.backends = append(ih.backends, b)
	return nil
}

func (ih *IpHashStrategy) Unregister(backendUrl string) error {
	if len(ih.backends) == util.MinBackends {
		return &errors.ErrBackendLimitUnderflow{}
	}
	for i, b := range ih.backends {
		if b.URL.String() == backendUrl {
			ih.backends = append(ih.backends[:i], ih.backends[i+1:]...)
			return nil
		}
	}
	return &errors.ErrBackendNotFound{BackendUrl: backendUrl}
}

func NewIpHashStrategy() *IpHashStrategy {
	return &IpHashStrategy{
		backends: make([]*backend.Backend, 0, util.MaxBackends),
		hasher:   fnv.New32a(),
	}
}

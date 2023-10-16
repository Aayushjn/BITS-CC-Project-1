package strategy

import (
	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/errors"
	"github.com/aayushjn/load-balancer/util"
)

type IpHashStrategy struct {
	backends []*backend.Backend
}

func (ih *IpHashStrategy) Backends() []*backend.Backend {
	return ih.backends
}

func (ih *IpHashStrategy) Next() *backend.Backend {
	// TODO: Implement this
	return nil
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
	}
}

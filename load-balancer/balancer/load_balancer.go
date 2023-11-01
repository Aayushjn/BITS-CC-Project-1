package balancer

import (
	"fmt"
	"net/http"

	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/aayushjn/load-balancer/balancer/strategy"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type LoadBalancer struct {
	balancingStrategy strategy.LoadBalancingStrategy
	sessionTable      cmap.ConcurrentMap[string, *backend.Backend]
}

func (lb *LoadBalancer) Register(b *backend.Backend, params map[string]any) error {
	b.SetProxyHandler(lb.RequestHandler())
	return lb.balancingStrategy.Register(b, params)
}

func (lb *LoadBalancer) Unregister(url string) error {
	return lb.balancingStrategy.Unregister(url)
}

func (lb *LoadBalancer) HealthCheck() {
	var status string
	var err error
	for _, b := range lb.balancingStrategy.Backends() {
		status = "up"
		err = b.Test()
		b.SetAlive(err == nil)
		if err != nil {
			status = "down"
		}
		fmt.Printf("%s [%s]\n", b.URL, status)
	}
}

func (lb *LoadBalancer) RequestHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend := lb.balancingStrategy.Next(r)
		if backend != nil {
			backend.IncrementInFlightRequests()
			backend.ReverseProxy.ServeHTTP(w, r)
			backend.DecrementInFlightRequests()
			return
		}
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	})
}

func NewLoadBalancer(balancingStrategy strategy.LoadBalancingStrategy) *LoadBalancer {
	return &LoadBalancer{
		balancingStrategy: balancingStrategy,
		sessionTable:      cmap.New[*backend.Backend](),
	}
}

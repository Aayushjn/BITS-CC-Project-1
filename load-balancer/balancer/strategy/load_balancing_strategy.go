package strategy

import "github.com/aayushjn/load-balancer/balancer/backend"

type LoadBalancingStrategy interface {
	Backends() []*backend.Backend
	Next() *backend.Backend
	Register(*backend.Backend, map[string]any) error
	Unregister(string) error
}

var AllowedStrategies = []string{"least_conns", "power_of_2", "random", "round_robin", "weighted_round_robin"}

func NewLoadBalancingStrategy(strategy string) LoadBalancingStrategy {
	switch strategy {
	case "least_conns":
		return NewLeastConnsStrategy()
	case "power_of_2":
		return NewPowerOfTwoStrategy()
	case "random":
		return NewRandomStrategy()
	case "round_robin":
		return NewRoundRobinStrategy()
	case "weighted_round_robin":
		return NewWeightedRoundRobinStrategy()
	default:
		return nil
	}
}

package config

import "github.com/BurntSushi/toml"

type Config struct {
	Strategy        string             `toml:"strategy"`
	CpuThreshold    float64            `toml:"cpu_threshold"`
	MemoryThreshold float64            `toml:"memory_threshold"`
	Network         string             `toml:"network"`
	Service         ServiceConfig      `toml:"service"`
	LoadBalancer    LoadBalancerConfig `toml:"load_balancer"`
}

type ServiceConfig struct {
	Name  string `toml:"name"`
	Image string `toml:"image"`
}

type LoadBalancerConfig struct {
	Image    string `toml:"image"`
	Strategy string `toml:"strategy"`
}

func LoadConfig(configFile string) (Config, error) {
	var conf Config
	_, err := toml.DecodeFile(configFile, &conf)
	return conf, err
}

package config

import (
	"net"
	"strconv"
)

type Config struct {
	Host      string // IP address of host machine
	Port      int
	StaticDir string // Path to static assets
	LogFile   string
}

func (config *Config) Addr() string {
	return net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
}

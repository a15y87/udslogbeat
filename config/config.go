// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	SocketPath string `config:"socketpath"`
	MaxMessageSize int `config:"maxmessagesize"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	MaxMessageSize: 8192,
	SocketPath: "/run/udslogbeat.sock",
}

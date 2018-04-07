// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period          time.Duration `config:"period"`
	NatsHost        string        `config:"nats_host"`
	NatsMonitorPort int           `config:"nats_monitor_port"`
}

var DefaultConfig = Config{
	Period:          1 * time.Second,
	NatsHost:        "127.0.0.1",
	NatsMonitorPort: 8200,
}

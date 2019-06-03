package udpout

import (
	"github.com/elastic/beats/libbeat/outputs/codec"
)

type udpoutConfig struct {
	Host string `config:"host"`
	Port int    `config:"port"`
	Codec codec.Config `config:"codec"`
}

var (
	defaultConfig = udpoutConfig{
		Port: 5556,
	}
)

func (c *udpoutConfig) Validate() error {
	return nil
}

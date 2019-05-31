package udpout

import (
	"github.com/elastic/beats/libbeat/outputs/codec"
)

type config struct {
	Host string `config:"host"`
	Port int    `config:"port"`
	Codec codec.Config `config:"codec"`
}

var (
	defaultConfig = config{
		Port: 5556,
	}
)

func (c *config) Validate() error {
	return nil
}

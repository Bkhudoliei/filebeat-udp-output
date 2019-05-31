package udpout

import (
	"fmt"
	"github.com/elastic/beats/libbeat/common/file"
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
	if c.NumberOfFiles < 2 || c.NumberOfFiles > file.MaxBackupsLimit {
		return fmt.Errorf("The number_of_files to keep should be between 2 and %v",
			file.MaxBackupsLimit)
	}
	return nil
}

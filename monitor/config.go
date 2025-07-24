package monitor

import (
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Hosts    []string `toml:"hosts"`
	Port     int      `toml:"port"`
	Interval int      `toml:"interval"`
}

func ReadConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

    return config
}

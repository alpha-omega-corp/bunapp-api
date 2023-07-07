package app

import (
	"embed"
	"io/fs"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed config
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

func FS() fs.FS {
	unwrapFSOnce.Do(func() {
		fileSys, err := fs.Sub(embedFS, "config")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fileSys
	})
	return unwrappedFS
}

type Config struct {
	Service string
	Env     string

	Debug     bool   `yaml:"debug"`
	SecretKey string `yaml:"secret_key"`

	DB struct {
		ADDR string `yaml:"addr"`
		NAME string `yaml:"name"`
		USER string `yaml:"user"`
		PASS string `yaml:"pass"`
	} `yaml:"db"`

	GPT struct {
		HOST   string `yaml:"host"`
		BEARER string `yaml:"bearer"`
	}
}

func ReadConfig(fileSys fs.FS, service, env string) (*Config, error) {
	b, err := fs.ReadFile(fileSys, env+".yaml")
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	cfg.Service = service
	cfg.Env = env

	return cfg, nil
}

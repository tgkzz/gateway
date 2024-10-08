package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Env        string        `json:"env"`
	Host       string        `json:"host"`
	Port       string        `json:"port"`
	AuthServer AuthServerCfg `json:"auth_server"`
}

type AuthServerCfg struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func MustRead(configPath string) *Config {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Config
	if err = json.Unmarshal(b, &config); err != nil {
		panic(err)
	}

	return &config
}

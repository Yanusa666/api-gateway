package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LogLevel   string           `json:"log_level"`
	HttpServer HttpServerConfig `json:"http_server"`
	Censor     CensorConfig     `json:"censor"`
	News       NewsConfig       `json:"news"`
	Comments   CommentsConfig   `json:"comments"`
}

func NewConfig() *Config {
	conf := new(Config)

	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

type HttpServerConfig struct {
	ListenAddress string `json:"listen_address"`
}

type NewsConfig struct {
	URI string `json:"URI"`
}

type CommentsConfig struct {
	URI string `json:"URI"`
}

type CensorConfig struct {
	URI string `json:"URI"`
}

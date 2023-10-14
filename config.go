package main

import (
	"encoding/json"
	"os"
)

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Device struct {
	Address uint8 `json:"address"`
	Length  uint8 `json:"length"`
}

type RelayConfig struct {
	Auth         AuthConfig   `json:"auth"`
	ReadDevices  []Device     `json:"readDevices"`
	ReadInterval int          `json:"readInterval"`
	Server       ServerConfig `json:"server"`
	TopicPrefix  string       `json:"topicPrefix"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func getConfig(config *RelayConfig) {
	var buffer = make([]byte, 1024)
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	length, err := file.Read(buffer)
	if err != nil {
		panic(err)
	}

	buffer = buffer[0:length]
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		panic(err)
	}
}

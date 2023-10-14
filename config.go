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
	Server       ServerConfig `json:"server"`
	ReadDevices  []Device     `json:"readDevices"`
	ReadInterval uint         `json:"readInterval"`
	TopicPrefix  string       `json:"topicPrefix"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func getConfig() RelayConfig {
	var buffer = make([]byte, 1024)
	var data RelayConfig
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(buffer, &data)
	return data
}

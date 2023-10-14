package main

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	gocron "github.com/go-co-op/gocron"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload = msg.Payload()
	var state []uint8
	_, addressStr := path.Split(msg.Topic())
	address, err := strconv.ParseUint(addressStr, 16, 8)
	if err != nil {
		fmt.Printf("Error parsing address: %v", err)
		return
	}

	err = json.Unmarshal(payload, &state)
	if err != nil {
		fmt.Printf("Error parsing payload: %v", err)
		return
	}

	fmt.Printf("Writing %d bytes to device: %x", len(state), address)
	write(uint8(address), state)
}

var readDevice = func(client mqtt.Client, device Device, outputTopicPrefix string) {
	fmt.Printf("Reading %d bytes from device: %x", device.Length, device.Address)
	outputTopic := path.Join(outputTopicPrefix, strconv.FormatUint(uint64(device.Address), 16))
	var buffer = make([]uint8, device.Length)
	received := read(device.Address, buffer)
	if received != int(device.Length) {
		fmt.Printf("Only read %d bytes from device: %x", received, device.Address)
	}

	payload, err := json.Marshal(buffer)
	if err != nil {
		fmt.Printf("Error building payload: %v", err)
		return
	}

	client.Publish(outputTopic, 0, false, payload)
}

func main() {
	var config = getConfig()

	opts := mqtt.NewClientOptions()
	server := fmt.Sprintf("tcp://%s:%d", config.Server.Host, config.Server.Port)
	fmt.Printf("Connecting to: %s", server)
	opts.AddBroker(server)
	opts.SetClientID(config.Auth.Username)
	opts.SetUsername(config.Auth.Username)
	opts.SetPassword(config.Auth.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	inputTopic := path.Join(config.TopicPrefix, "input", "#")
	if token := client.Subscribe(inputTopic, 1, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Subscribed to topic %s", inputTopic)

	outputTopicPrefix := path.Join(config.TopicPrefix, "output")
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(config.ReadInterval).Seconds().Do(func() {
		for _, device := range config.ReadDevices {
			readDevice(client, device, outputTopicPrefix)
		}
	})

	if err != nil {
		panic(err)
	}

	for client.IsConnected() {

	}

	s.Stop()
}

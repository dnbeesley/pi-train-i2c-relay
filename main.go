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
	fmt.Printf("Connect lost: %v\n", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload = msg.Payload()
	var state []uint8
	_, addressStr := path.Split(msg.Topic())
	address, err := strconv.ParseUint(addressStr, 16, 8)
	if err != nil {
		fmt.Printf("Error parsing address: %v\n", err)
		return
	}

	err = json.Unmarshal(payload, &state)
	if err != nil {
		fmt.Printf("Error parsing payload: %v\n", err)
		return
	}

	fmt.Printf("Writing %d bytes to device: %x\n", len(state), address)
	write(uint8(address), state)
}

var readAllDevices = func(client mqtt.Client, devices []Device, outputTopicPrefix string) {
	fmt.Println("Reading i2c device states")
	for _, device := range devices {
		readDevice(client, device, outputTopicPrefix)
	}
}

var readDevice = func(client mqtt.Client, device Device, outputTopicPrefix string) {
	fmt.Printf("Reading %d bytes from device: %x\n", device.Length, device.Address)
	outputTopic := path.Join(outputTopicPrefix, strconv.FormatUint(uint64(device.Address), 16))
	var buffer = make([]uint8, device.Length)
	received := read(device.Address, buffer)
	if received != int(device.Length) {
		fmt.Printf("Only read %d bytes from device: %x\n", received, device.Address)
	}

	var array = make([]int, received)
	for i, v := range buffer {
		array[i] = int(v)
	}

	payload, err := json.Marshal(array)
	if err != nil {
		fmt.Printf("Error building payload: %v\n", err)
		return
	}

	client.Publish(outputTopic, 0, false, payload)
}

func main() {
	var config RelayConfig
	getConfig(&config)

	opts := mqtt.NewClientOptions()
	server := fmt.Sprintf("tcp://%s:%d", config.Server.Host, config.Server.Port)
	fmt.Println("Connecting to:", server)
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

	fmt.Println("Subscribed to topic:", inputTopic)

	outputTopicPrefix := path.Join(config.TopicPrefix, "output")
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(config.ReadInterval).Seconds().Do(readAllDevices, client, config.ReadDevices, outputTopicPrefix)

	if err != nil {
		panic(err)
	}

	s.StartBlocking()
}

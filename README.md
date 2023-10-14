# PI train i2c relay

This program is designed to run on a raspberry pi to act a relay between MQTT and a i2c bus.

## Configuration

Configuration file, config.json, needs to be in the applications context path. The file needs to define the following:

```JSON
{
  "auth": {
    "username": "username",
    "password": "password"
  },
  "server": {
    "host": "192.168.1.2",
    "port": 1883
  },
  "readDevices": [
    {
      "address": 80, // Corresponding to device: 0x50
      "length": 2
    },
    {
      "address": 81, // Corresponding to device: 0x51
      "length": 3
    }
  ],
  "readInterval": 1, // seconds
  "topicPrefix": "i2c-relay"
}
```

This example file will subscribe to: "i2c-relay/input/#"" and publish to: "i2c-relay/output/50" and "i2c-relay/output/51". In the MQTT topic hexidecimal strings are used from the numbers. The message payloads should be a JSON array of bytes, for both inout and output. Each input will be read roughly every 1 second.

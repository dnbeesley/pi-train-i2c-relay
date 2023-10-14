package main

import (
	"sync"

	i2c "github.com/d2r2/go-i2c"
)

var mutex sync.Mutex

func read(address uint8, buffer []byte) int {
	mutex.Lock()
	defer mutex.Unlock()

	return readImpl(address, buffer)
}

func readImpl(address uint8, buffer []byte) int {
	var connection *i2c.I2C
	connection, err := i2c.NewI2C(address, 1)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
	}()

	length, err := connection.ReadBytes(buffer)
	if err != nil {
		panic(err)
	}

	return length
}

func write(address uint8, payload []byte) int {
	mutex.Lock()
	defer mutex.Unlock()

	return writeImpl(address, payload)
}

func writeImpl(address uint8, payload []byte) int {
	var connection *i2c.I2C
	connection, err := i2c.NewI2C(address, 1)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
	}()

	length, err := connection.WriteBytes(payload)
	if err != nil {
		panic(err)
	}

	return length
}

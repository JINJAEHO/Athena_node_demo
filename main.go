package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"os"
	"time"
)

func main() {
	// Fisrt argument is for node's group
	arg_first := os.Args[1]
	// Second argument is for node's port
	arg_second := os.Args[2]

	// Check if config file is modified with GoRoutine
	go func() {
		check := true
		for {
			ConfigData = LoadConfig()
			check = checkConfig(ConfigData)
			if !check {
				log.Println("Config file was modified!!")
				Hash = MakeHashOfConfig(ConfigData)
			}
			time.Sleep(time.Millisecond * 3000)
		}
	}()

	NewServer(arg_second, arg_first)
	ServerStart(arg_second)

}

var Hash [32]byte

// Check if config file is modified by Hash
func checkConfig(config Config) bool {
	tmpHash := MakeHashOfConfig(config)

	if tmpHash == Hash {
		return true
	}
	return false
}

// Create data of config file to make hash
func prepareData(config Config) []byte {
	data := bytes.Join([][]byte{
		[]byte(config.URL.GotWalletReq),
		[]byte(config.URL.GotTxsReq),
		[]byte(config.URL.RegTx),
		[]byte(config.URL.GotDeatilinfo),
		[]byte(config.MspPort),
		[]byte(config.Public),
		[]byte(config.TestIp),
	}, []byte{})
	return data
}

// Create hash of config file
func MakeHashOfConfig(config Config) [32]byte {
	data := prepareData(config)
	hash := sha256.Sum256(data)

	return hash
}

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

	NewServer(arg_second, arg_first)

	check := true
	for {
		check = checkConfig()
		if !check {
			log.Println("Config file was modified!!")
			break
		}
		time.Sleep(time.Millisecond * 3000)
	}
}

var Hash [32]byte

// Check if config file is modified by Hash
func checkConfig() bool {
	// Load config.json
	config := LoadConfig()

	tmpHash := MakeHashOfConfig(config)

	if tmpHash == Hash {
		return true
	}
	return false
}

// create data of config file to make hash
func prepareData(config Config) []byte {
	data := bytes.Join([][]byte{
		[]byte(config.URL.GetWalletReq),
		[]byte(config.URL.GotTxsReq),
		[]byte(config.URL.RegTx),
		[]byte(config.URL.GotDeatilinfo),
		[]byte(config.MspPort),
		[]byte(config.Public),
		[]byte(config.TestIp),
	}, []byte{})
	return data
}

// create hash of config file
func MakeHashOfConfig(config Config) [32]byte {
	data := prepareData(config)
	hash := sha256.Sum256(data)

	return hash
}

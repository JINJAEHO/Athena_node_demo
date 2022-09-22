package main

import (
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

func main() {
	// Fisrt argument is for node's group
	arg_first := os.Args[1]
	// Second argument is for node's port
	arg_second := os.Args[2]
	// Third argument is for node's name
	arg_third := os.Args[3]

	// Check if config file is modified with GoRoutine
	go func() {
		check := true
		for {
			conf := LoadConfig()
			check = checkConfig(conf)
			if !check {
				log.Println("Config file was modified!!")
				ConfigData = conf
				Hash = MakeHashOfConfig(conf)
			}
			time.Sleep(time.Millisecond * 3000)
		}
	}()

	// Clean node's buffer when node is zombie
	go func() {
		for {
			if InitValue.Group == "zombie" {
				usageStr := GetMemoryUsage()
				usage, _ := strconv.ParseFloat(usageStr, 32)
				if usage >= 20 {
					log.Println("Free Memory")
					debug.FreeOSMemory()
				} else {
					time.Sleep(time.Millisecond * 1000)
				}
			} else {
				time.Sleep(time.Millisecond * 3000)
			}
		}
	}()

	NewServer(arg_second, arg_first, arg_third)
	ServerStart(arg_second)
	go TcpStart(arg_second)
}

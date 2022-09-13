package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
)

//var PUBLIC string = "http://175.45.201.176"
var TEST_IP string
var MSP_PORT string

var strategy string

// structure for config.json
type Config struct {
	URL []struct {
		GetWalletReq  string `json:"/GetWalletReq"`
		RegTx         string `json:"/RegTx"`
		GotTxsReq     string `json:"/GotTxsReq"`
		GotDeatilinfo string `json:"/GotDeatilinfo"`
	} `json:"url"`
	Public  string `json:"public"`
	TestIp  string `json:"testIp"`
	MspPort string `json:"mspPort"`
}

// new node notify and try to connect to MSP
func NewNode(myPort string) {
	// Load config.json
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	var config Config
	json.NewDecoder(file).Decode(&config)
	log.Println(config)
	//PUBLIC = config.Public
	TEST_IP = config.TestIp
	MSP_PORT = config.MspPort

	// Get my ip address as string type
	var myIP string

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			myIP = ipv4.String()
		}
	}

	myIpMap := make(map[string]string)
	myIpMap[myPort] = myIP + ":" + myPort

	ipMarshal, _ := json.Marshal(myIpMap)

	// Notify that new node want to join to MSP
	res, err := http.Post("http://"+TEST_IP+":"+MSP_PORT+"/newNodeAlarm", "application/json", bytes.NewBuffer(ipMarshal))
	log.Println(res)

	if err != nil {
		log.Println(err)
	}

	// variable for Blacklist ip table
	blacklist := make(map[string]string)
	json.NewDecoder(res.Body).Decode(&blacklist)
	for _, v := range blacklist {
		log.Println(v)
	}

	// init status variable
	strategy = "normal"
}

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

var PUBLIC string
var TEST_IP string
var MSP_PORT string

var Strategy string
var Blacklist map[string]string

var ConfigData Config

// structure for config.json
type Config struct {
	URL struct {
		GotWalletReq  string `json:"/GotWalletReq"`
		RegTx         string `json:"/RegTx"`
		GotTxsReq     string `json:"/GotTxsReq"`
		GotDeatilinfo string `json:"/GotDeatilinfo"`
	} `json:"url"`
	Public  string `json:"public"`
	TestIp  string `json:"testIp"`
	MspPort string `json:"mspPort"`
}

// structure for Node info
type Addr struct {
	NewNode string `json:"node"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

// new node notify and try to connect to MSP
func NewNode(myPort string, status string) {
	// Load config.json
	ConfigData = LoadConfig()
	log.Println(ConfigData)

	PUBLIC = ConfigData.Public
	TEST_IP = ConfigData.TestIp
	MSP_PORT = ConfigData.MspPort

	// Save initial hash
	Hash = MakeHashOfConfig(ConfigData)

	// Get my ip address as string type
	var myIP string

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			myIP = ipv4.String()
		}
	}

	var myIpStruct Addr
	myIpStruct.Address = myIP
	myIpStruct.NewNode = myPort
	myIpStruct.Type = status

	ipMarshal, _ := json.Marshal(myIpStruct)

	// Notify that new node want to join to MSP
	res, err := http.Post("http://"+TEST_IP+":"+MSP_PORT+"/RegNewNode", "application/json", bytes.NewBuffer(ipMarshal))
	log.Println(res)

	if err != nil {
		log.Println(err)
	}

	// Get blacklist ip table from MSP
	UpdateBlacklist(res.Body)

	// Init status variable
	if status == "1" {
		Strategy = "normal"
	} else {
		Strategy = "abnormal"
	}
}

// Update blacklist ip table from MSP
func UpdateBlacklist(body io.Reader) {
	// Empty blacklist table
	Blacklist = make(map[string]string, 0)

	json.NewDecoder(body).Decode(&Blacklist)
	for _, v := range Blacklist {
		log.Println(v)
	}
}

// Load config.json file
func LoadConfig() Config {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	var config Config
	json.NewDecoder(file).Decode(&config)

	return config
}

package main

import (
	"log"
	"net/http"
)

func NewServer(myPort string, status string) {
	NewNode(myPort, status)
	setRoute()
}

func setRoute() {
	http.HandleFunc("/PingReq", GetStatus)
	http.HandleFunc("/ChangeStrategy", nil)
	http.HandleFunc("/TableUpdateAlarm", TableUpdate)
}

// Check client IP if it is blacklist
func CheckBlacklist(clientIP string) bool {
	for _, v := range Blacklist {
		if clientIP == v {
			return false
		}
	}
	return false
}

// Select Back-end URL
func SelectURL(reqURL string) string {
	var targetURL string

	switch reqURL {
	case "GotWalletReq":
		targetURL = ConfigData.URL.GotWalletReq
	case "RegTx":
		targetURL = ConfigData.URL.RegTx
	case "GotTxsReq":
		targetURL = ConfigData.URL.GotTxsReq
	case "GotDeatilinfo":
		targetURL = ConfigData.URL.GotDeatilinfo
	}

	return targetURL
}

func GetStatus(w http.ResponseWriter, req *http.Request) {

}

// Manage blacklist ip table from MSP
func TableUpdate(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request(/TableUpdateAlarm) from MSP")
	UpdateBlacklist(req.Body)
}

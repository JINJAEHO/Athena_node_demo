package main

import "net/http"

func NewServer(myPort string, status string) {
	NewNode(myPort, status)
	setRoute()
}

func setRoute() {
	http.HandleFunc("/PingReq", GetStatus)
	http.HandleFunc("/ChangeStrategy", nil)
	http.HandleFunc("/TableUpdateAlarm", nil)
}

func GetStatus(w http.ResponseWriter, req *http.Request) {

}

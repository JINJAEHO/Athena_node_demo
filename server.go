package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func NewServer(myPort string, status string) {
	NewNode(myPort, status)
	setRoute()
}

func ServerStart(myPort string) {
	log.Println("노드 실행 : ", ConfigData.TestIp+":"+myPort)
	if err := http.ListenAndServe(ConfigData.TestIp+":"+myPort, nil); err != nil {
		log.Println(err)
		return
	}
}

func setRoute() {
	http.HandleFunc("/PingReq", GetStatus)
	http.HandleFunc("/ChangeStrategy", GetStrategy)
	http.HandleFunc("/TableUpdateAlarm", TableUpdate)

	http.HandleFunc("/SendWalletReq", SendWalletReq)
	http.HandleFunc("/SendTxsReq", SendTxsReq)
	http.HandleFunc("/SendDetailinfo", SendDetailinfo)
	http.HandleFunc("/SendRegTx", SendRegTx)
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

// Get ping and check my memory status and then send response with memory status to MSP
func GetStatus(w http.ResponseWriter, req *http.Request) {
	vm, _ := mem.VirtualMemory()
	used := vm.Used

	pid := os.Getpid()
	ps, _ := process.NewProcess(int32(pid))
	mem, _ := ps.MemoryInfo()
	percent, _ := ps.MemoryPercent()
	vms := mem.VMS

	usage := fmt.Sprint(((float32(used) * percent) / float32(vms)) * 100)

	w.Header().Set("Content-Type", "text/plain")
	json.NewEncoder(w).Encode(usage)
}

// Get and change strategy
func GetStrategy(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request for changing strategy")

	var stgy string
	json.NewDecoder(req.Body).Decode(&stgy)
	Strategy = stgy
}

// Manage blacklist ip table from MSP
func TableUpdate(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request(/TableUpdateAlarm) from MSP")
	UpdateBlacklist(req.Body)
}

// Sending semi-blackIP to MSP
func SendIP(semiBlack string) {
	ipMarshal, _ := json.Marshal(semiBlack)

	res, err := http.Post("http://"+ConfigData.TestIp+":"+ConfigData.MspPort+"/SendBlackIP", "application/json", bytes.NewBuffer(ipMarshal))
	if res != nil {
		res.Body.Close()
	}

	if err != nil {
		log.Println(err)
	}

	_, errs := io.Copy(ioutil.Discard, res.Body)
	if errs != nil {
		log.Println(errs)
	}
}

// Return client ip from http request
func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func SendWalletReq(w http.ResponseWriter, req *http.Request) {
	if Strategy == "normal" {
		ip, err := getIP(req)
		if err != nil {
			log.Println(err)
		}
		log.Println(ip)
	}
}

func SendTxsReq(w http.ResponseWriter, req *http.Request) {
	logFile := OpenLogFile()
	WriteLog(logFile, "TEST")
}

func SendDetailinfo(w http.ResponseWriter, req *http.Request) {

}

func SendRegTx(w http.ResponseWriter, req *http.Request) {

}

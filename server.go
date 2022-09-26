package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func NewServer(myPort string, group string, name string) {
	NewNode(myPort, group, name)
	setRoute()
}

func ServerStart(myPort string) {
	log.Println("노드 실행 : ", ConfigData.Public+":"+myPort)
	if err := http.ListenAndServe(":"+myPort, nil); err != nil {
		log.Println(err)
		return
	}
}

func TcpStart(myPort string) {
	port, _ := strconv.Atoi(myPort)
	port += 100
	log.Println(port)
	ln, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		log.Println(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept() // 클라이언트가 연결되면 TCP 연결을 리턴
		log.Println("conn", conn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer conn.Close() // main 함수가 끝나기 직전에 TCP 연결을 닫음
		go GetStatus(conn) // 패킷을 처리할 함수를 고루틴으로 실행
	}
}

func setRoute() {
	// http.HandleFunc("/PingReq", GetStatus)
	http.HandleFunc("/ChangeStrategy", GetStrategy)
	http.HandleFunc("/TableUpdateAlarm", TableUpdate)
	//http.HandleFunc("/DelayResult", delayResult)

	for k := range ConfigData.URL {
		http.HandleFunc(k, ServiceReq)
	}
}

// Check client IP if it is blacklist
func CheckBlacklist(clientIP string) bool {
	for _, v := range InitValue.Blacklist {
		if clientIP == v {
			return true
		}
	}
	return false
}

// Select Back-end URL
func SelectURL(reqURL string) string {
	var targetURL string

	for k, v := range ConfigData.URL {
		if reqURL == k {
			targetURL = v
			break
		}
	}
	return targetURL
}

// Get ping and check my memory status and then send response with memory status to MSP
func GetStatus(conn net.Conn) {
	for {
		var groupName string
		json.NewDecoder(conn).Decode(&groupName)
		if len(groupName) > 0 {
			usage := GetMemoryUsage()

			InitValue.Group = groupName

			logData := "Nodename," + InitValue.NodeName + ",clientIP,null,url,null,address," + ConfigData.Public + ":" + InitValue.MyPort + ",cpuUsed," + fmt.Sprint(usage) + ",group," + InitValue.Group
			// logFile := OpenLogFile(InitValue.NodeName + "-Status")
			// defer logFile.Close()
			// WriteLog(logFile, logData)
			statusQue <- logData
			json.NewEncoder(conn).Encode(usage)
		}
	}
}

func GetMemoryUsage() string {
	vm, _ := mem.VirtualMemory()
	//used := vm.Used
	total := vm.Total

	pid := os.Getpid()
	ps, _ := process.NewProcess(int32(pid))
	Mem, _ := ps.MemoryInfo()
	percent, _ := ps.MemoryPercent()
	vms := Mem.VMS
	usage := fmt.Sprint(((float32(total) * (percent / 100.0)) / float32(vms)) * 100.0)
	//usage := fmt.Sprint(100 * float32(Mem.RSS) / float32(vms))
	cpuPercent, _ := ps.CPUPercent()
	log.Println("=====================================================")
	log.Println("cpu percent:", cpuPercent)
	log.Println("memory percent:", percent)
	log.Println("memory usage:", usage)
	return fmt.Sprint(cpuPercent)
}

// Get and change strategy
func GetStrategy(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request for changing strategy")

	var stgy string
	json.NewDecoder(req.Body).Decode(&stgy)
	InitValue.Strategy = stgy
}

// Manage blacklist ip table from MSP
func TableUpdate(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request(/TableUpdateAlarm) from MSP")
	UpdateBlacklist(req.Body)
}

// Sending semi-blackIP to MSP
func SendIP(ip string, code string) {
	// logFile := OpenLogFile(InitValue.NodeName + "-Warning")
	// defer logFile.Close()
	if code == "warning" {
		data := "Nodename," + InitValue.NodeName + ",warning," + ip + ",danger,null"
		warningQue <- data
		//WriteLog(logFile, "Nodename,"+InitValue.NodeName+",warning,"+ip+",danger,null")
	} else if code == "danger" {
		data := "Nodename," + InitValue.NodeName + ",warning,null,danger," + ip
		warningQue <- data
		//WriteLog(logFile, "Nodename,"+InitValue.NodeName+",warning,null,danger,"+ip)
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

// HandleFunc about all of service requests
func ServiceReq(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	log.Println("Start time:", startTime)

	ip, err := getIP(req)
	if err != nil {
		log.Println(err)
	}

	// Check blacklist
	IsBlack := CheckBlacklist(ip)
	if IsBlack {
		log.Println("Access of blacklist is found. Block the request!")
		return
	}
	url_path := req.URL.Path

	log.Println(url_path, "접속, ClientIP:", ip)

	// Write log

	// logFile := OpenLogFile(InitValue.NodeName + "-Status")
	// defer logFile.Close()
	// WriteLog(logFile, "Nodename,"+InitValue.NodeName+",clientIP,"+ip+",url,"+url_path+",address,null,memUsed,null,group,null")
	logData := "Nodename," + InitValue.NodeName + ",clientIP," + ip + ",url," + url_path + ",address,null,memUsed,null,group,null"
	statusQue <- logData

	if InitValue.Strategy == "ABNORMAL" {
		SendIP(ip, "danger")
	} else {
		usage, _ := strconv.ParseFloat(GetMemoryUsage(), 32)
		if usage >= 1 && usage < 5 {
			SendIP(ip, "warning")
		}
	}
	// targetURL := SelectURL(url_path)
	// res, err := http.Post("http://"+ConfigData.Public+":"+ConfigData.GatePort+targetURL, "application/json", req.Body)
	// closeResponse(res, err)
	totalTime := time.Since(startTime)
	vps := float64(totalTime) / float64(time.Millisecond)
	// logFile = OpenLogFile(InitValue.NodeName + "-Performance")
	// defer logFile.Close()
	// WriteLog(logFile, "Nodename,"+InitValue.NodeName+",vps,"+fmt.Sprint(vps))
	performanceQue <- "Nodename," + InitValue.NodeName + ",vps," + fmt.Sprint(vps)
}

// start pBFT for delay
func SendReqPBFT(body io.Reader, ip string) {
	res, err := http.Post("http://"+ip+"/StartDelay", "application/json", body)
	closeResponse(res, err)
}

// Close response for preventing memory leak
func closeResponse(res *http.Response, err error) {
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.Println("log close error:", err)
	}

	_, errs := io.Copy(ioutil.Discard, res.Body)
	if errs != nil {
		log.Println("io discard error: ", errs)
	}
}

// Get result of delay and forwarding request to BackEnd
func delayResult(w http.ResponseWriter, req *http.Request) {
	targetURL := SelectURL(req.URL.Path)
	res, err := http.Post("http://"+ConfigData.Public+":"+ConfigData.MspPort+targetURL, "application/json", req.Body)
	closeResponse(res, err)
}

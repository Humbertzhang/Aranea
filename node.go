package main

import (
	"github.com/humbertzhang/Aranea/core"
	"github.com/humbertzhang/Aranea/core/status"

	"github.com/gorilla/mux"

	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var node = new(core.Node)

func main() {
	//处理命令行参数
	runPort := "9699"
	name := "node"
	myip := "localhost"
	masterHost := ""
	masterPort := ""

	argv := os.Args
	argc := len(os.Args)
	if argc >= 2 {
		for i := 0; i < argc; i += 1 {
			if argv[i] == "-mh" {
				//master host
				masterHost = argv[i+1]
			}
			if argv[i] == "-mp" {
				//master port
				masterPort = argv[i+1]
			}
			//自己监听的端口
			if argv[i] == "-p" || argv[i] == "--port" {
				runPort = argv[i+1]
			}
			if argv[i] == "-n" {
				name = argv[i+1]
			}
		}
	}

	//未指定master host & port 退出
	if masterHost == "" || masterPort == "" {
		fmt.Println("need master's port and master's host")
		os.Exit(1)
	}

	/*获取本地ip*/
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				myip = ipnet.IP.String()
				fmt.Println(myip)
				// 疑惑：这里可以打印出4个ip地址，不清楚哪个才是正确的，取的en0那个
				break;
			}
		}
	}
	node.Port = runPort
	node.Name = name
	node.IP = myip
	node.Status = status.STATUSNORMAL

	/*注册到master*/
	err = node.RegisterToMaster(masterHost, masterPort)
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	/*http端口等待任务*/
	router := mux.NewRouter()
	router.HandleFunc("/jobs", WaitJobs).Methods("POST")
	http.Handle("/", router)

	// run
	port := ":"+runPort
	fmt.Println("Node listening localhost",port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func WaitJobs(writer http.ResponseWriter, request *http.Request) {

}


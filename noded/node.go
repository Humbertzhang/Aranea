package main

import (
	"encoding/json"
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
	// TODO:这里ip选择有错误，现在只有在我开着ssr时，master才能ping到node
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
	router.HandleFunc("/node/health", Health).Methods(http.MethodGet)
	router.HandleFunc("/node/jobs", WaitJobs).Methods(http.MethodPost)
	router.HandleFunc("/node/pong", Pong).Methods(http.MethodPost)
	http.Handle("/", router)

	// run
	port := ":"+runPort
	fmt.Println("Node listening localhost",port)
	log.Fatal(http.ListenAndServe(port, nil))
}


func Health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	return
}


// TODO:实现node接受Master ping过来的结构体之后进行一次爬取.200ok则返回statusOK，否则返回失败状态码
// 心跳函数
func Pong(writer http.ResponseWriter, request *http.Request) {
	job := &core.Job{}
	err := json.NewDecoder(request.Body).Decode(job)
	fmt.Printf("JOB:%+v:\n", job)
	if err != nil {
		println("pong error")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	statusCode, _, err := job.Crawler.Request()
	if statusCode != http.StatusOK {
		writer.WriteHeader(statusCode)
		writer.Write([]byte("status NOT ok"))
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("status ok"))
	}
	fmt.Println("NODE: MADE response.")
	return
}



func WaitJobs(writer http.ResponseWriter, request *http.Request) {

}


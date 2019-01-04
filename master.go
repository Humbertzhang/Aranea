package main

import (
	"Aranea/core"
	"Aranea/core/status"

	"github.com/gorilla/mux"

	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"os"
)


var master = new(core.Master)

func main() {
	/*处理命令行参数*/
	runPort := "9999"

	argv := os.Args
	argc := len(os.Args)
	if argc >= 2 {
		if argv[1] == "-p" || argv[1] == "--port" {
			runPort = argv[2]
		}
		if argv[1] == "-h" || argv[1] == "--help" {
			help()
		}
		if argv[1] == "-v" || argv[1] == "--version" {
			version()
		}
	}

	/* init master */
	master.IP = "0.0.0.0"
	master.Port = runPort

	/* init server */
	router := mux.NewRouter()
	router.HandleFunc("/master/register", registerNode).Methods("POST")
	router.HandleFunc("/master/unregister", unregisterNode).Methods("DELETE")
	http.Handle("/", router)
	port := ":"+runPort
	fmt.Println("listening localhost",port)
	log.Fatal(http.ListenAndServe(port, nil))
}

/*
node json:
{
	"name": String,
	"ip": string,
	"port": string,
}
*/
type NodeJson struct {
	Name 		string  `json:"name"`
	Ip 			string 	`json:"ip"`
	Port 		string  `json:"port"`
}

// TODO:对IP和PORT格式做检查
func registerNode(writer http.ResponseWriter, request *http.Request) {
	node := &core.Node{}
	err := json.NewDecoder(request.Body).Decode(node)
	//println(node.Name)
	if err != nil {
		println("...")
		/*返回错误信息*/
		writer.WriteHeader(500)
		return
	}
	node.Status = status.STATUSNORMAL

	/*检查是否已经被注册*/
	for _, v := range(master.Nodes) {
		if v.IP == node.IP && v.Port == node.Port {
			println("error:重复注册\n")
			writer.WriteHeader(400)
			return
		}
	}

	/*注册并打印已经有了的node*/
	master.RegisterNode(node)
	master.TraverseNodes()
}

func unregisterNode(writer http.ResponseWriter, request *http.Request) {

}


func version() {
	version := "0.0.1"
	println("Aranea Master Version: ", version)
	os.Exit(0)
}

func help() {
	println("usage:")
	println("run: ./master")
	println("run in specified port: ./master -p 8888")
	os.Exit(0)
}


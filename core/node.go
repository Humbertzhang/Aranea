package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Node struct {
	Name   string  			`json:"name"`
	IP     string			`json:"ip"`
	Port   string			`json:"port"`
	Status int				`json:"status"`

	// outtime ping失败次数
	OutTimes int
	// unix时间戳秒数, 代表下次ping的时间
	NextPing time.Time

	JQ 		JobQueue
}

type NodeCreateJson struct {
	Name   string  			`json:"name"`
	IP     string			`json:"ip"`
	Port   string			`json:"port"`
}



// 注册到master
func (node *Node) RegisterToMaster(masterIP string, masterPort string) (err error) {
	// 准备json
	n := NodeCreateJson{
		Name: node.Name,
		IP: node.IP,
		Port: node.Port,
	}
	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(n)
	if err != nil {
		return err
	}
	// post json to master.
	postUrl := "http://" + masterIP + ":" + masterPort + "/master/register"
	res, err := http.Post(postUrl, "application/json; charset=utf-8", data)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("created failed")
	}
	println("Node Registered to", masterIP, ":", masterPort)
	return nil
}


/*
TODO:需要一个消息队列或者etcd？，去同步任务执行结果

*/
// 执行任务
/*
func (node *Node) LoopDoJobs() {
	for j := range node.JQ.Queue {
		resp, err := DoJob(j)
		if err != nil {

		}
	}
}

func DoJob(job *Job) (*http.Response, error) {


	return &http.Response{}, nil
}
*/







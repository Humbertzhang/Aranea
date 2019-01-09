package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/humbertzhang/Aranea/core/status"
	"github.com/humbertzhang/Aranea/core/utils"
	"net/http"
	"time"
)

type Master struct {
	JQ 		JobQueue
	IP 		string
	Port 	string
	Nodes 	[]*Node
}

// 在register时对node做检查
func (master *Master) RegisterNode(node *Node) {
	master.Nodes = append(master.Nodes, node)
}

// 移除node
func (master *Master) RemoveNode(name string) error {
	for k, v := range master.Nodes {
		if v.Name == name {
			master.Nodes = append(master.Nodes[:k], master.Nodes[k+1:]...)
			return nil
		}
	}
	return errors.New("Node Not Found.")
}

// 遍历打印node
func (master *Master) TraverseNodes() {
	fmt.Println("All Nodes Master have:")
	for k, v := range master.Nodes {
		fmt.Printf("%d:%+v\n", k, v)
	}
}


func (master *Master) DistributeJob() {
	nodec := 0
	for j := range master.JQ.Queue {

		// 任务超出最大失败次数，便打印并放弃
		if j.FailedTimes >= len(master.Nodes) {
			fmt.Println("JobId:", j.Id, "JobURL:", j.URL, "超出最大失败次数.")
			continue
		}

		nodeTryTimes := 0
		for {
			// tryTime记录遍历node的个数，当大于node总个数时打印error，放弃循环
			if nodeTryTimes > len(master.Nodes) {
				fmt.Println("Error: 无node可用!")
				break
			}

			// 检查node是否正常，不正常继续
			if master.Nodes[nodec].Status != status.STATUSNORMAL {
				nodec = (nodec + 1) % len(master.Nodes)
				nodeTryTimes += 1
				continue
			}
		}


		resp, err := master.DistributeJobToNode(j, master.Nodes[nodec])
		if err != nil {
			fmt.Println(err)
			j.FailedTimes += 1

			// 失败后判断ping是否成功，若失败，则标记node为unknown， 若成功，则将任务重新加入到队列中
			nodePingURL := "http://" + master.Nodes[nodec].IP + ":" + master.Nodes[nodec].Port + "/node/pong"
			if err = master.pingOnce(nodePingURL); err != nil {
				fmt.Println(err)
				master.Nodes[nodec].Status = status.STATUSUNKNOWN
			} else {
				master.JQ.PushJob(j)
				continue
			}
		}
		// 第一次成功了，打印返回信息
		fmt.Println(resp)
	}
}

func (master *Master) DistributeJobToNode(job *Job, node *Node) (*http.Response, error) {
	nodeURL := "http://" + node.IP + ":" + node.Port + "/node/job"

	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(job); err != nil{
		fmt.Println(err)
		return &http.Response{}, err
	}

	res, err := http.Post(nodeURL, "application/json; charset=utf-8", data)
	if res == nil {
		fmt.Println("Master:No response")
		return &http.Response{}, err
	}
	if err != nil {
		return res, err
	}
	return res, nil
}




// master 发送http心跳任务, 设置超时时间
// node接受心跳之后执行, 若超时node仍未返回则设为失败。
// node返回的为一个core/status包中的一个status，master根据返回的status设置node的status
// ping机制: 正常(NORMAL) 1min 1次，且发送任务前需要额外一次ping
// 1次失败后进入 (unknown) 状态, 不可以被调度任务， 15s一次 ping
// unknown 状态下 10次失败后unconnected, 15min ping一次
// unconnected 10次之后下线. status -> offline
// offline的会在下次循环中被取出来

func (master *Master) Ping() {
	// 初始化时间，将所有node的nextPing时间设为now()+某s
	// 初始化失败次数
	for _, node := range master.Nodes {
		node.NextPing = time.Now()
		node.OutTimes = 0
		node.Status = status.STATUSNORMAL
	}

	for {
		// 无限循环.

		for _, node := range master.Nodes {
			// 时间是否到了可以ping的时候
			// 小于代表不可以，直接continue
			if time.Now().Before(node.NextPing) {
				continue
			}

			// status setter 状态控制器，控制unknown向unconnected状态转换和unconnected向offline转换、
			switch node.Status {
			case status.STATUSNORMAL:
				break
			case status.STATUSUNKNOWN:
				fmt.Println("STATUSUNUNKNOWN")
				if node.OutTimes > 10 {
					node.Status = status.STATUSUNCONNECTED
					node.OutTimes = 0
				}
				break
			case status.STATUSUNCONNECTED:
				fmt.Println("STATUSUNCONNECTED")
				if node.OutTimes > 10 {
					node.Status = status.STATUSOFFLINE
				}
				break
			case status.STATUSOFFLINE:
				//去掉node
				fmt.Println("Rmoving Node:", node.Name)
				err := master.RemoveNode(node.Name)
				if err != nil {
					panic(err)
				}
			default:
				fmt.Println("status setter error:", node.Status)
				panic(nil)
				break
			}

			// ping 操作
			// 失败会报错,导致状态转变为unknown
			NodeURL := "http://" + node.IP + ":" + node.Port + "/node/pong"
			err := master.pingOnce(NodeURL)
			// 状态控制
			// 成功会将非normal状态转换回来
			// 失败，如果原来为Normal,则会将Nromal状态转换回正常状态
			if err != nil {
				if node.Status == status.STATUSNORMAL {
					node.Status = status.STATUSUNKNOWN
				}
				node.OutTimes += 1
			} else {
				if node.Status != status.STATUSNORMAL {
					node.Status = status.STATUSNORMAL
					node.OutTimes = 0
				}
			}


			// ping time setter 时间设置器

			now := time.Now()

			// 15s
			exprNormal := cronexpr.MustParse("*/15 * * * * * *")
			// 5s
			exprUnKnown := cronexpr.MustParse("*/5 * * * * * *")
			// 5min
			exprUnConnected := cronexpr.MustParse("* */5 * * * * *")
			exprOffline := cronexpr.MustParse("* * * * * * *")

			switch node.Status {
			case status.STATUSNORMAL:
				node.NextPing = exprNormal.Next(now)
				break
			case status.STATUSUNKNOWN:
				node.NextPing = exprUnKnown.Next(now)
				break
			case status.STATUSUNCONNECTED:
				node.NextPing = exprUnConnected.Next(now)
				break
			case status.STATUSOFFLINE:
				node.NextPing = exprOffline.Next(now)
				break
			default:
				fmt.Println("ping time setter error:", node.Status)
				panic(nil)
				break
			}

		}

		// 每100ms扫描一次
		// 通过select的阻塞机制实现
		// 其中timer的原理是 timer中有个channel 叫 C, 当到期时timer会向这个channel投递一个对象
		select {
		case <- time.NewTimer(100*time.Millisecond).C:
		}
	}
}


var URLs = [...]string{"https://www.baidu.com/", "https://www.tencent.com", "https://m.taobao.com/", "https://www.bytedance.com/"}
var size = 4
// ping如果失败会选取下一个URL再试一次,如果还是失败才是真正失败
func (master *Master) pingOnce(NodeURL string) (err error) {
	// 根据毫秒数随机选取一个URL
	millSecondNow := int(time.Now().UnixNano()/1000000)
	CrawlerURL := URLs[millSecondNow%size]
	// 创建Job Json
	data, err := PingJobJsonCreater(CrawlerURL)
	if err != nil {
		return err
	}

	// http post 发送给node url
	// 并检查运行结果,若失败有第二次机会尝试.
	statusCode, err := master.postPingJobToNode(data, NodeURL)
	if err != nil  || statusCode != http.StatusOK{
		// 检查第二次
		CrawlerURL2 := URLs[(millSecondNow+1)%size]
		data, err = PingJobJsonCreater(CrawlerURL2)
		if err != nil {
			return err
		}
		statusCode, err = master.postPingJobToNode(data, NodeURL)
		// 第二次还是失败
		if err != nil || statusCode != http.StatusOK {
			return errors.New("ping error")
		}
	}
	// 成功返回nil
	return nil
}

func PingJobJsonCreater(CrawlerURL string) (buffer *bytes.Buffer, err error){
	crawler := &Crawler{
		URL:         CrawlerURL,
		Method:      http.MethodGet,
	}
	job := &Job{
		Crawler: crawler,
		Id:      utils.StringTimeStampNanoSecond(),
		Name:    "ping",
		Delay:   0,
	}
	//bytes, err := json.Marshal(job) ?
	data := new(bytes.Buffer)
	err = json.NewEncoder(data).Encode(job)
	if err != nil {
		return data, err
	}
	return data, nil
}


func (master *Master) postPingJobToNode(data *bytes.Buffer, nodeURL string) (statusCode int, err error){
	fmt.Println("Master: NodeURL:", nodeURL)
	res, err := http.Post(nodeURL, "application/json; charset=utf-8", data)
	if res == nil {
		fmt.Println("Master:No response")
		return http.StatusInternalServerError, err
	}
	if err != nil {
		return res.StatusCode, err
	}
	fmt.Println("status code:", res.StatusCode)
	return res.StatusCode, nil
}

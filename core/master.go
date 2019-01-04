package core

import (
	"errors"
	"fmt"
)

type Master struct {
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

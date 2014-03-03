package bbbikeng

import (
	"sync"
)

type nodeData []*Node

type NodeSet struct {
	data nodeData
	closedData nodeData
	mu  sync.Mutex
}

func (this *NodeSet) Add(value *Node) {

	contains := this.Contains(value)
	this.mu.Lock()
	if !contains {
		this.data = append(this.data, value)
	}
	this.mu.Unlock()
}

func (this *NodeSet) GetByKey(key int) (value *Node) {
	this.mu.Lock()
	for _, node := range this.data {
		if node.NodeID == key {
			this.mu.Unlock()
			return node
		}
	}
	this.mu.Unlock()
	return nil
}

func (this *NodeSet) Remove(value *Node) {
	this.RemoveByKey(value.NodeID)
}

func (this *NodeSet) RemoveByKey(key int) () {
	this.mu.Lock()
	var newData nodeData
	for _, node := range this.data {
		if node.NodeID != key {
			newData = append(newData, node)
		}
	}
	this.data = newData
	this.mu.Unlock()
}

func (this *NodeSet) ContainsByKey(key int) (exists bool) {
	this.mu.Lock()
	for _, node := range this.data{
		if node.NodeID == key {
			this.mu.Unlock()
			return true
		}
	}
	this.mu.Unlock()
	return false
}

func (this *NodeSet) Contains(value *Node) (exists bool) {
	return this.ContainsByKey(value.NodeID)
}

func (this *NodeSet) Length() (int) {
	return len(this.data)
}

func NewNodeSet() (*NodeSet) {
	return &NodeSet{}
}

func (d nodeData) Len() int {
	return len(d)
}

func (d nodeData) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d nodeData) Less(i, j int) bool {
	return d[i].F < d[j].F
}

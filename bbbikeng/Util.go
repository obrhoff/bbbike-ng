package bbbikeng

import (
	"sort"
)

type nodeData []*Node

type NodeSet struct {
	data nodeData
	closedData nodeData
}

func (this *NodeSet) Add(value *Node) {
	contains := this.Contains(value)
	if !contains {
		this.data = append(this.data, value)
		sort.Sort(this.data)
	}
}

func (this *NodeSet) GetByKey(key int) (value *Node) {
	for _, node := range this.data {
		if node.NodeID == key {
			return node
		}
	}
	return nil
}

func (this *NodeSet) Remove(value *Node) {
	this.RemoveByKey(value.NodeID)
}

func (this *NodeSet) RemoveByKey(key int) () {
	var newData nodeData
	for _, node := range this.data {
		if node.NodeID != key {
			newData = append(newData, node)
		}
	}
	sort.Sort(newData)
	this.data = newData
}

func (this *NodeSet) ContainsByKey(key int) (exists bool) {
	for _, node := range this.data{
		if node.NodeID == key {
			return true
		}
	}
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

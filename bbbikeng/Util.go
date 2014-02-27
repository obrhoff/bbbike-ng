package bbbikeng

import (
	"math"

)


func (f *Point) PointIsValid() bool {

	return f.Lat != 0.0 && f.Lng != 0.0

}

func (f *Point) Compare(comparePoint Point) (equal bool) {

	thresholdLat := math.Abs(math.Abs(f.Lat) - math.Abs(comparePoint.Lat))
	thresholdLng := math.Abs(math.Abs(f.Lng) - math.Abs(comparePoint.Lng))
	return (thresholdLat <= 0.0000001 && thresholdLng <= 0.0000001)

}

type NodeSet struct {
	data map[int]Node
}

func (this *NodeSet) Add(value Node) {
	contains := this.Contains(value)
	if !contains {
		this.data[value.NodeID] = value
	}
}

func (this *NodeSet) GetByKey(key int) (value Node) {
	return this.data[key]
}

func (this *NodeSet) Remove(value Node) {
	this.RemoveByKey(value.NodeID)
}

func (this *NodeSet) RemoveByKey(key int) () {
	delete(this.data, key)
}

func (this *NodeSet) ContainsByKey(key int) (exists bool) {
	_, exists = this.data[key]
	return exists
}

func (this *NodeSet) Contains(value Node) (exists bool) {
	return this.ContainsByKey(value.NodeID)
}

func (this *NodeSet) Length() (int) {
	return len(this.data)
}

func NewNodeSet() (*NodeSet) {
	return &NodeSet{make(map[int]Node)}
}

type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].F < h[j].F }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NodeHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Node))
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}


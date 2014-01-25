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

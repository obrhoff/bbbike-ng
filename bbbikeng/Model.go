package bbbikeng

import (
)

type Node struct {

	NodeID int
	NodeGeometry Point
	Streets []Street
	Neigbors []Node

	DistanceFromParentNode int
	PathFromParentNode []Node

	Walkable bool
	Heuristic int

	G int
	F int

	Value interface{}

}

type base struct {
	ID int
	Type string
	Path []Point
}

type City struct {
	ID int
	base
	Country string
	Geometry []Point
}

type Bla struct {
	base
	Node
}


type Street struct {
	Name string
	Nodes []Node
	base
}

type Cyclepath struct {
	Name string
	base
}

type Quality struct {
	Name string
	base
}

type Greenway struct {
	Name string
	base
}

func (f *base) SetPathFromGeoJSON(jsonInput string) {
	f.Path = ConvertGeoJSONtoPath(jsonInput)
}

func (f base) GetGeoJSONPath() (jsonOutput string) {
	return ConvertPathToGeoJSON(f.Path)
}

/*
func (this *Node) Neighbors(potentialNeighbors []Node) (neighbors []Node) {

	for _, street := range this.Streets {

		var baseIndex int
		for i, point := range street.Path {
			if (this.NodeGeometry.Compare(point)) {
				baseIndex = i
			}
		}

		leftPath := []Point{this.NodeGeometry}
		rightPath := []Point{this.NodeGeometry}

		leftDistance := 0
		rightDistance := 0

		leftLoop: for i := baseIndex + 1;  i < len(street.Path); i++ {
			point := street.Path[i]
			leftPath := append(leftPath, point)
			leftDistance += DistanceFromPointToPoint(point, street.Path[i-1])
			for _, node := range potentialNeighbors {
				if node.NodeGeometry.Compare(point) {
					node.DistanceFromParentNode = leftDistance
					node.ParentNode = this
					node.PathFromParentNode = leftPath
					neighbors = append(neighbors, node)
					break leftLoop
				}
			}
		}

		rightLoop: for i := baseIndex - 1;  i >= 0; i-- {
			point := street.Path[i]
			rightPath := append(rightPath, point)
			rightDistance += DistanceFromPointToPoint(point, street.Path[i+1])
			for _, node := range potentialNeighbors {
				if node.NodeGeometry.Compare(point) {
					node.DistanceFromParentNode = rightDistance
					node.ParentNode = this
					node.PathFromParentNode = rightPath
					neighbors = append(neighbors, node)
					break rightLoop
				}
			}
		}
	}

	return neighbors

}

*/

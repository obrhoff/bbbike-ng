package bbbikeng

type Node struct {

	NodeID int
	NodeGeometry Point
	Neigbors []Node

	DistanceFromParentNode int
	StreetFromParentNode Path

	ParentNodes *Node
	Walkable bool
	Heuristic float64

	G float64
	F float64

	Value interface{}

}


type Path struct{

	ID int
	WayID int
	Name string
	Type string
	PathIndex int
	Path []Point
	Attributes attributes

}

type attributes struct {

	Quality string
	Cyclepath string
	Greenway string
	TrafficLight bool

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

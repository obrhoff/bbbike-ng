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

	TrafficLight bool
	Value interface{}
}

type Path struct{
	ID int
	WayID int
	Name string
	Type string
	PathIndex int
	Path []Point
	Attributes []AttributeInterface
}

type AttributeInterface interface {

	Type() string
	SetType(Type string)

	Geometry() []Point
	SetGeometry(geometry []Point)

	Relevance() bool
	SetRelevance(relevance bool)

	CalculateScore (preference *Preferences) float64
}

type Attribute struct {
	attributeType string
	geometry []Point
	isRelevant bool
	AttributeInterface
}

type CyclepathAttribute struct {
	Attribute
}

type GreenwayAttribute struct {
	Attribute
}

type QualityAttribute struct {
	Attribute
}

type UnlitAttribute struct {
	Attribute
}

type TrafficLightAttribute struct {
	Attribute
}

func (this *Attribute) Type() (Type string) {
	return this.attributeType
}

func (this *Attribute) Geometry() (geometry []Point) {
	return this.geometry
}

func (this *Attribute) Relevance() (relevance bool) {
	return this.isRelevant
}

func (this *Attribute) SetType(newType string)  () {
	this.attributeType = newType
}

func (this *Attribute) SetGeometry(newGeometry []Point)() {
	this.geometry = newGeometry
}

func (this *Attribute) SetRelevance(relevance bool)() {
	this.isRelevant = relevance
}

type base struct {
	ID int
	Type string
	Path []Point
}

type City struct {
	ID int
	Name string
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

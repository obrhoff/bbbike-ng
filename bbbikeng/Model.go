package bbbikeng
import (
	"strings"
)

type Node struct {

	NodeID int
	NodeGeometry Point
	Neigbors []Node

	DistanceFromParentNode int
	StreetFromParentNode Path

	ParentNodes *Node
	Walkable bool
	Heuristic int

	G int
	F int

	flippedDirection bool
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

	GlobalAttribute []AttributeInterface
	NormalAttribute []AttributeInterface
	FlippedAttribute []AttributeInterface

	Attributes []AttributeInterface
}

type AttributeInterface interface {

	Id() int
	SetId (id int)

	Name() string
	SetName(Name string)

	Type() string
	SetType(Type string)

	Geometry() []Point
	SetGeometry(geometry []Point)

	CalculateScore (preference *Preferences) float64
	SetPathFromGeoJSON (jsonInput interface {})

}


type Attribute struct {
	id int
	attributeType string
	geometry []Point
	AttributeInterface
}

type WayAttribute struct {
	name string
	Attribute
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

type HandicapAttribute struct {
	name string
	Attribute
}

func (this *WayAttribute) SetName(Name string) () {
	this.name = strings.Replace(Name, "'", "''", -1)
}

func (this *WayAttribute) Name() (Name string) {
	return this.name
}

func (this *HandicapAttribute) SetName(Name string) () {
	this.name = strings.Replace(Name, "'", "''", -1)
}

func (this *HandicapAttribute) Name() (Name string) {
	return this.name
}

func (this *Attribute) SetId(id int) {
	this.id = id
}

func (this *Attribute) Id() (id int) {
	return this.id
}

func (this *Attribute) SetName(Name string) () {
}

func (this *Attribute) Name() (Name string) {
	return ""
}

func (this *Attribute) Type() (Type string) {
	return this.attributeType
}

func (this *Attribute) Geometry() (geometry []Point) {
	return this.geometry
}

func (this *Attribute) SetType(newType string)  () {
	this.attributeType = newType
}

func (this *Attribute) SetGeometry(newGeometry []Point)() {
	this.geometry = newGeometry
}

func (this *Attribute) GetJsonGeometry()(geojson GeoJSON) {
	return geojson
}

func (this *Attribute) SetFromGeojson(geojson string) () {

}



package bbbikeng

type Street struct {
	PathID     int
	Name       string
	StreetType string
	Path       []Point
	Intersections []Intersection
}

type Intersection struct {
	Coordinate Point
	Street Street
}

type Attributes struct {
	Greenways     string
	Quality       string
	TrafficLights string
}

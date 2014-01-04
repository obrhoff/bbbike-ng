package bbbikeng

type Street struct {
	PathID     int
	Name       string
	StreetType string
	Path       []Point
}

type Attributes struct {
	Greenways     string
	Quality       string
	TrafficLights string
}

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

func (f *Street) SetIntersections(){
	f.Intersections = GetStreetIntersections(f)
}

func (f *Street) SetPathFromGeoJSON(jsonInput string) {
	f.Path = ConvertGeoJSONtoPath(jsonInput)
}

func (f Street) GetGeoJSONPath() (jsonOutput string) {
	return ConvertPathToGeoJSON(f.Path)
}

func (f *Intersection) SetCoordinationFromGeoJSON(jsonInput string) {
	f.Coordinate = ConvertGeoJSONtoPoint(jsonInput)
}


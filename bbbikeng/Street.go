package bbbikeng


type basePath struct {
	Type string
	Path []Point
}

type advancedPath struct {
	Name string
	ID int
}

type City struct {
	CityID int
	Name string
	Country string
	Geometry []Point
}

type Street struct {
	basePath
	advancedPath
	Intersections []Intersection
}

type Cyclepath struct {
	basePath
}

type Quality struct {
	basePath
}

type Greenway struct {
	basePath
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

func (f *basePath) SetPathFromGeoJSON(jsonInput string) {
	f.Path = ConvertGeoJSONtoPath(jsonInput)
}

func (f basePath) GetGeoJSONPath() (jsonOutput string) {
	return ConvertPathToGeoJSON(f.Path)
}

func (f *Intersection) SetCoordinationFromGeoJSON(jsonInput string) {
	point := ConvertGeoJSONtoPath(jsonInput)
	if len(point) > 0 {
		f.Coordinate = point[0]
	}
}


package bbbikeng

type BBJSON struct {

	Response bool
	Preferences Preferences
	Distance int
	Time int
	Lights int
	Instruction []BBJSONInstruction
	Path [][2]float64

}

type BBJSONInstruction struct {
	Index int
	Name string
	Type string
	Instruction string
}

type BBJSONAttribute struct {
	Type string
	Path []Point
}

func (this *Route) GetGeojson() (geojson GeoJSON) {

	geojson.Type = "LineString"
	for _, segment := range this.way {
		var newPoint [2]float64
		newPoint[0] = segment.Lng
		newPoint[1] = segment.Lat
		geojson.Coordinates = append(geojson.Coordinates, newPoint)
	}
	return geojson

}

func (this *Route) GetBBJson() (json BBJSON) {

	for _, segment := range this.way {
		var newPoint [2]float64
		newPoint[0] = segment.Lng
		newPoint[1] = segment.Lat
		json.Path = append(json.Path, newPoint)
	}

	for _, path := range this.detailed {
		var newInstruction BBJSONInstruction
		newInstruction.Name = path.Name
		newInstruction.Index = path.PathIndex
		newInstruction.Type = path.Type
		json.Instruction = append(json.Instruction, newInstruction)
	}

	json.Preferences = this.Preferences
	json.Lights = this.TrafficLights
	json.Distance = int(this.distance * 1000.0)
	json.Response = true

	return json
}

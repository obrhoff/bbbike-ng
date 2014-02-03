package bbbikeng

import (
"log"
)

type Route struct {

	way []Point
	detailed []*Path
	distance float64
	time int
	nodes []*Node

}

func constructRoute (finalNode Node) (route Route) {

	var parentNode *Node
	parentNode = &finalNode
	// gather all nodes
	for parentNode != nil {
		route.nodes = append(route.nodes, parentNode)
		parentNode = parentNode.ParentNodes
	}
	// reverse list
	var tmpNodeList []*Node
	for i := len(route.nodes)-1; i >= 0; i-- {
		tmpNodeList = append(tmpNodeList, route.nodes[i])
	}
	route.nodes = tmpNodeList

	for _, node := range route.nodes {
		// build up the actually path
		if len(node.StreetFromParentNode.Path) > 0 {
			firstPoint := node.StreetFromParentNode.Path[0]
			if !firstPoint.Compare(node.NodeGeometry) {
				for i := 1; i < len(node.StreetFromParentNode.Path); i++ {
					route.way = append(route.way, node.StreetFromParentNode.Path[i])
				}
			} else {
				for i := len(node.StreetFromParentNode.Path)-2; i >= 0; i-- {
					route.way = append(route.way, node.StreetFromParentNode.Path[i])
				}
			}
			// needs to be fixed for first node
			node.StreetFromParentNode.PathIndex = len(route.way)-1
			if (len(route.detailed) > 0  && (route.detailed[len(route.detailed)-1].WayID != node.StreetFromParentNode.WayID) && route.detailed[len(route.detailed)-1].Name != node.StreetFromParentNode.Name ) || len(route.detailed) < 1{
				route.detailed = append(route.detailed, &node.StreetFromParentNode)
			}
		}

	}

	return route

}

func CalculateHeuristic(parentNode Node, neighborNode Node) (heuristic float64) {

	heuristic = DistanceFromLinePoint(neighborNode.StreetFromParentNode.Path)
	return heuristic

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
		newInstruction.Roadname = path.Name
		newInstruction.PathIndex = path.PathIndex
		newInstruction.Type = path.Type
		newInstruction.Quality = path.Attributes.Quality

		json.Instruction = append(json.Instruction, newInstruction)

	}

	log.Println("Output:", json)

	json.Response = true

	return json
}

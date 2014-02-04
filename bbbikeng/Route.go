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
	for i := len(route.nodes)-2; i >= 0; i-- {
		tmpNodeList = append(tmpNodeList, route.nodes[i])
	}

	route.nodes = tmpNodeList

	startNode := route.nodes[0]
	endNode := route.nodes[len(route.nodes)-1]

	for _, node := range route.nodes {

		streetPath := flippPath(node)
		if node.NodeID == startNode.NodeID || node.NodeID == endNode.NodeID {

		// needs to be fixed for first node

			var index int
			if node.NodeID == startNode.NodeID {
				index = 0
			} else {
				index = len(route.way)-1
			}
			node.StreetFromParentNode.PathIndex = index
			route.detailed = append(route.detailed, &node.StreetFromParentNode)

			for i := 0; i < len(node.StreetFromParentNode.Path); i++ {
				route.way = append(route.way, node.StreetFromParentNode.Path[i])
			}


		} else {

			if route.detailed[len(route.detailed)-1].WayID != node.StreetFromParentNode.WayID && route.detailed[len(route.detailed)-1].Name != node.StreetFromParentNode.Name {
				node.StreetFromParentNode.PathIndex = len(route.way)-1
				route.detailed = append(route.detailed, &node.StreetFromParentNode)
			}

			for i := 1; i < len(streetPath); i++ {
				route.way = append(route.way, streetPath[i])
			}

		}
	}

	return route

}

func CalculateHeuristic(parentNode Node, neighborNode Node) (heuristic float64) {

	heuristic = DistanceFromLinePoint(neighborNode.StreetFromParentNode.Path)
	return heuristic

}

func flippPath(node *Node) (path []Point) {

	firstPoint := node.NodeGeometry
	if firstPoint.Compare(node.StreetFromParentNode.Path[0]) {
		var flippedPath []Point
		for i := len(node.StreetFromParentNode.Path)-1; i >= 0; i-- {

			point := node.StreetFromParentNode.Path[i]
			flippedPath = append(flippedPath, point)
		}
		return flippedPath
	}

	return node.StreetFromParentNode.Path
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
		newInstruction.PathIndex = path.PathIndex
		newInstruction.Type = path.Type
		newInstruction.Quality = path.Attributes.Quality

		json.Instruction = append(json.Instruction, newInstruction)

	}

	log.Println("Output:", json)

	json.Response = true

	return json
}

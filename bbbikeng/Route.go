package bbbikeng

import (

)

type Route struct {

	way []Point
	distance float64
	time int
	nodes []*Node

}

func constructRoute (finalNode Node) (route Route) {

	// todo: reverse node array since it starts from the end

	var parentNode *Node
	parentNode = &finalNode

	route.way = append(route.way, parentNode.NodeGeometry)
	for parentNode != nil {

		// add node. could be useful for building up later more detailed
		route.nodes = append(route.nodes, parentNode)

			// build up the actually path
			if len(parentNode.StreetFromParentNode.Path) > 0 {
				firstPoint := parentNode.StreetFromParentNode.Path[0]
				if firstPoint.Compare(parentNode.NodeGeometry) {
					for i := 1; i < len(parentNode.StreetFromParentNode.Path); i++ {
						route.way = append(route.way, parentNode.StreetFromParentNode.Path[i])
					}
				} else {
					for i := len(parentNode.StreetFromParentNode.Path)-2; i >= 0; i-- {
						route.way = append(route.way, parentNode.StreetFromParentNode.Path[i])
					}
				}
			}
		parentNode = parentNode.ParentNodes
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
		json.path = append(json.path, newPoint)
	}

	return json
}

package bbbikeng

import (
	"log"
)

type Route struct {

	distance float64
	nodes []*Node

}

func constructRoute (finalNode Node) (route Route) {

	var parentNode *Node
	parentNode = finalNode.ParentNodes
	for parentNode != nil {
		route.nodes = append(route.nodes, parentNode)
		if parentNode.ParentNodes != nil {
			route.distance =+ DistanceFromPointToPoint(parentNode.NodeGeometry, parentNode.ParentNodes.NodeGeometry)
			log.Println("Node:", parentNode.NodeID, "Geometry:", parentNode.NodeGeometry.Lat,",",parentNode.NodeGeometry.Lng)
		}
		parentNode = parentNode.ParentNodes
	}

	return route

}

func (this *Route) GetGeojson() (geojson GeoJSON) {

	geojson.Type = "LineString"
	for _, node := range this.nodes {
		var newPoint [2]float64
		newPoint[0] = node.NodeGeometry.Lng
		newPoint[1] = node.NodeGeometry.Lat
		geojson.Coordinates = append(geojson.Coordinates, newPoint)
	}

	return geojson
}

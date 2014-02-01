package bbbikeng

import (
	"log"
)

type Route struct {

	way []Point
	distance float64
	nodes []*Node

}

func constructRoute (finalNode Node) (route Route) {

	var parentNode *Node
	parentNode = finalNode.ParentNodes
	for parentNode != nil {
		route.nodes = append(route.nodes, parentNode)
		if parentNode.ParentNodes != nil {
			firstPoint := parentNode.StreetFromParentNode.Path[0]
			if firstPoint.Compare(parentNode.NodeGeometry) {
				for i := 0; i < len(parentNode.StreetFromParentNode.Path); i++ {
					route.way = append(route.way, parentNode.StreetFromParentNode.Path[i])
				}
			} else {
				for i := len(parentNode.StreetFromParentNode.Path)-1; i >= 0; i-- {
					route.way = append(route.way, parentNode.StreetFromParentNode.Path[0])
				}
			}
		}
		parentNode = parentNode.ParentNodes
	}

	return route

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

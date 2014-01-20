package bbbikeng

import (
	"log"
	"fmt"
)

type Route struct {




}

func GetNearestNode(point Point) (nearestNode Node) {

	return FindNearestNode(point)

}

func GetRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	var openList = NewNodeSet();
	var closedList = NewNodeSet();

	openList.Add(startNode)

	for openList.Length() > 0 {

		lowLnd := -1
		var currentNode Node

		for _, value := range openList.data {
			score := value.DistanceFromParentNode + DistanceFromPointToPoint(value.NodeGeometry, endNode.NodeGeometry)
			if score < lowLnd || lowLnd < 0 {
				currentNode = value
			}
		}

		log.Println("CurrentNode", currentNode)

		if currentNode.NodeID == endNode.NodeID {
			fmt.Println("Done!", currentNode)
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		for _, neighbor := range currentNode.Neigbors {

			neighbor.Neigbors = GetNeighborNodesFromNode(neighbor)

			if closedList.Contains(neighbor) || len(neighbor.Neigbors) < 1 {
				continue
			}

			gScore := currentNode.DistanceFromParentNode + neighbor.DistanceFromParentNode
			gScoreIsBest := false;

			if !openList.Contains(neighbor) {
				gScoreIsBest = true;
				neighbor.Heuristic = DistanceFromPointToPoint(neighbor.NodeGeometry, endNode.NodeGeometry)
				openList.Add(neighbor)
			} else if(gScore < neighbor.Heuristic) {
				gScoreIsBest = true;
			}

			if (gScoreIsBest) {


			}


		}


	}

	return route
}

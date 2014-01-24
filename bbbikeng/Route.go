package bbbikeng

import (
	"log"
	"fmt"
)

type Route struct {


}


func GetRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	log.Println("StartNode:", startNode.NodeID)
	log.Println("EndNode:", endNode.NodeID)

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


			return route
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(currentNode);

		for _, neighbor := range neighbors {

			if closedList.Contains(neighbor) || !neighbor.Walkable  {
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

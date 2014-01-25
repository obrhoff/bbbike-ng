package bbbikeng

import (
	"fmt"
)

type Route struct {


}


func GetRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	fmt.Println("StartNode:", startNode.NodeID)
	fmt.Println("EndNode:", endNode.NodeID)

	var openList = NewNodeSet()
	var closedList = NewNodeSet()

	openList.Add(startNode)

	for openList.Length() > 0 {

		fmt.Println("OpenList Next:", openList.data)

		var bestNode Node
		bestNode.NodeID = -1
		for _, node := range openList.data {
			if bestNode.NodeID < 0 {
				bestNode = node
			} else {
				if bestNode.F >= node.F {
					bestNode = node
				}
			}
		}

		currentNode := bestNode
		fmt.Println("Test Node: ", currentNode.NodeID, " H: ", currentNode.F)

		if currentNode.NodeID == endNode.NodeID {

			var points []Point
			for _, point := range currentNode.PathFromParentNode{
				points = append(points, point.NodeGeometry)
			}
			fmt.Println("Done:", ConvertPathToGeoJSON(points))
			return route
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(currentNode);

		for _, neighbor := range neighbors {

			if closedList.Contains(neighbor) || !neighbor.Walkable  {
				continue
			}

			gScore := DistanceFromPointToPoint(currentNode.NodeGeometry, neighbor.NodeGeometry)
			gScoreIsBest := false;

			if !openList.ContainsByKey(neighbor.NodeID) {
				gScoreIsBest = true;
				neighbor.Heuristic = DistanceFromPointToPoint(neighbor.NodeGeometry, endNode.NodeGeometry)
				openList.Add(neighbor)
			}
			if(gScore < neighbor.Heuristic) {
				gScoreIsBest = true;
			}

			if (gScoreIsBest) {
					fmt.Println("Next Node: ", neighbor.NodeID, " H: ", neighbor.F)
					openList.Remove(neighbor)
					neighbor.G = gScore
					neighbor.F = neighbor.G + neighbor.Heuristic
					neighbor.PathFromParentNode = append(neighbor.PathFromParentNode, currentNode)
					openList.Add(neighbor)
					currentNode = neighbor
			}
		}


	}

	return route
}

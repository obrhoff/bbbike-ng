package bbbikeng

import (
	"log"
)

func GetAStarRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	log.Println("StartNode:", startNode.NodeID)
	log.Println("EndNode:", endNode.NodeID)

	var openList = NewNodeSet()
	var closedList = NewNodeSet()

	openList.Add(startNode)

	for openList.Length() > 0 {


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
		log.Println("ParentNode:", currentNode.NodeID , " Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng)

		if currentNode.NodeID == endNode.NodeID {
			return constructRoute(currentNode)
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(currentNode);

		for _, neighbor := range neighbors {

			if closedList.Contains(neighbor) || !neighbor.Walkable  {
				continue
			}

			gScore := DistanceFromLinePoint(neighbor.StreetFromParentNode.Path)
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
				openList.Remove(neighbor)
				neighbor.G = gScore
				neighbor.F = neighbor.G + neighbor.Heuristic
				neighbor.ParentNodes = &currentNode
				log.Println("Next Node: ", neighbor.NodeID, " H: ", neighbor.F, " ParentNode:", neighbor.ParentNodes.NodeID)
				openList.Add(neighbor)
				//currentNode = neighbor
			}
		}
	}

	return route
}

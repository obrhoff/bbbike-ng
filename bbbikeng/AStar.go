package bbbikeng

import (
	"log"
)

func GetAStarRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	log.Println("StartNode:", startNode.NodeID , "(",startNode.StreetFromParentNode.Name,") Geometry:", startNode.NodeGeometry.Lat, "," ,startNode.NodeGeometry.Lng)
	log.Println("EndNode:", endNode.NodeID , "(",endNode.StreetFromParentNode.Name,") Geometry:", endNode.NodeGeometry.Lat, "," ,endNode.NodeGeometry.Lng)

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
		log.Println("ParentNode:", currentNode.NodeID , "(",currentNode.StreetFromParentNode.ID, currentNode.StreetFromParentNode.Name,") Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng)
		if currentNode.NodeID == endNode.NodeID {
			return constructRoute(currentNode)
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(currentNode);

		for _, neighbor := range neighbors {

			if closedList.Contains(neighbor) || (!neighbor.Walkable && neighbor.NodeID != endNode.NodeID)  {
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
				log.Println("Possible Node:", neighbor.NodeID , "(",neighbor.StreetFromParentNode.Name,") Geometry:", neighbor.NodeGeometry.Lat, "," ,neighbor.NodeGeometry.Lng)
				openList.Add(neighbor)
			}
		}
	}

	return route
}

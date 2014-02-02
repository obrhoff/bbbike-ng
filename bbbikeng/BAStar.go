package bbbikeng

import (
	"log"
)

func GetBAStarRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	beginNodeChannel := make(chan Node)
	endNodeChannel := make(chan Node)

	log.Println("StartNode:", startNode.NodeID , " Geometry:", startNode.NodeGeometry.Lat, "," ,startNode.NodeGeometry.Lng)
	log.Println("EndNode:", endNode.NodeID , " Geometry:", endNode.NodeGeometry.Lat, "," ,endNode.NodeGeometry.Lng)



	go startAstarRoute(startNode, endNode, beginNodeChannel, endNodeChannel);
	go startAstarRoute(endNode, startNode, endNodeChannel, beginNodeChannel);

	return route

}

func startAstarRoute(startNode Node, endNode Node, inNodeChannel chan Node, outNodeChannel chan Node){

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

		outNodeChannel <- currentNode
		currentConcurrentNode := <-inNodeChannel

		log.Println("ParentNode:", currentNode.NodeID , " Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng)

		if currentNode.NodeID == endNode.NodeID || currentConcurrentNode.NodeID == currentNode.NodeID || closedList.Contains(currentConcurrentNode)  {
			// return constructRoute(currentNode)
			log.Println("jeogjegegheiugheu")
			break

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
				log.Println("Next Node: ", neighbor.NodeID, " H: ", neighbor.F, " Geometry:", neighbor.NodeGeometry.Lat, "," ,neighbor.NodeGeometry.Lng)
				openList.Add(neighbor)
				//currentNode = neighbor
			}
		}
	}
}

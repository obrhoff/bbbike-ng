package bbbikeng

import (
	"log"
)

func GetBAStarRoute(from Point, to Point) (route Route){

	startNode := FindNearestNode(from)
	endNode := FindNearestNode(to)

	beginNodeChannel := make(chan Node, 3)
	endNodeChannel := make(chan Node, 3)

	finalBeginChannel := make(chan Node, 1)
	finalEndChannel := make(chan Node, 1)

	log.Println("StartNode:", startNode.NodeID , " Geometry:", startNode.NodeGeometry.Lat, "," ,startNode.NodeGeometry.Lng)
	log.Println("EndNode:", endNode.NodeID , " Geometry:", endNode.NodeGeometry.Lat, "," ,endNode.NodeGeometry.Lng)

	go startAstarRoute(startNode, endNode, beginNodeChannel, endNodeChannel, finalBeginChannel, finalEndChannel);
	go startAstarRoute(endNode, startNode, endNodeChannel, beginNodeChannel, finalEndChannel, finalBeginChannel);

	finalEndNode := <-finalEndChannel;
	finalBeginNode := <-finalBeginChannel;

	log.Println("FinalBeginNode:", finalBeginNode)
	log.Println("FinalEndNode:", finalEndNode)

//	route = constructRoute(finalEndNode)

	return route

}

func startAstarRoute(startNode Node, endNode Node, inNodeChannel chan Node, outNodeChannel chan Node, finalBeginChannel chan Node, finalEndChannel chan Node){

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

		if closedList.Contains(currentConcurrentNode) {
			log.Println("Already Visited:", closedList.GetByKey(currentConcurrentNode.NodeID))
			finalBeginChannel <- closedList.GetByKey(currentConcurrentNode.NodeID)
			finalEndChannel <- currentConcurrentNode
			return
		} else if currentConcurrentNode.NodeID == currentNode.NodeID  {
			finalBeginChannel <- currentNode
			finalEndChannel <- currentConcurrentNode
			return
 		} else if currentNode.NodeID == endNode.NodeID {
			log.Println("Error")
			return
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
				openList.Add(neighbor)
			}
		}
	}
}

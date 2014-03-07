package bbbikeng

import (
	"log"
)

func (this *Route) StartBiRouting(startPoint Point, endPoint Point) {

	this.startNode = FindNearestNode(startPoint)
	this.endNode = FindNearestNode(endPoint)

	forwardChannel := make(chan *Node, 10)
	backwardChannel := make(chan *Node, 10)
	doneChannel := make(chan Node, 2)
	go this.GetBiAStarRoute(this.startNode, this.endNode, forwardChannel, backwardChannel, doneChannel, false)
	go this.GetBiAStarRoute(this.endNode, this.startNode, backwardChannel, forwardChannel, doneChannel, true)

	var forwardNode Node
	var backwardNode Node

	for node := range doneChannel {
		if node.flippedDirection {
			backwardNode = node
		} else {
			forwardNode = node
		}
		if backwardNode.NodeID == forwardNode.NodeID {
			break
		}
	}


	log.Println("ForwardNode:", forwardNode)
	log.Println("BackwardNode:", backwardNode)
	this.constructRoute(forwardNode)
	this.constructRoute(backwardNode)

}


func (this *Route) GetBiAStarRoute(startNode Node, endNode Node, forwardChannel chan *Node, backwardChannel chan *Node, doneChannel chan <- Node, reverse bool) (){

	log.Println("StartNode:", startNode.NodeID , "(",startNode.StreetFromParentNode.Name,") Geometry:", startNode.NodeGeometry.Lat, "," ,startNode.NodeGeometry.Lng)
	log.Println("EndNode:", endNode.NodeID , "(",endNode.StreetFromParentNode.Name,") Geometry:", endNode.NodeGeometry.Lat, "," ,endNode.NodeGeometry.Lng)

	var openList = NewNodeSet()
	var closedList = NewNodeSet()
	var concurrentNode *Node
	openList.Add(&startNode)

	for openList.Length() > 0 {

		currentNode := openList.data[0]
		currentNode.flippedDirection = reverse
		forwardChannel <- currentNode
		concurrentNode = <- backwardChannel
		log.Println("Concurrent Node:", concurrentNode)
		log.Println("ParentNode:", currentNode.NodeID , "(",currentNode.StreetFromParentNode.ID, currentNode.StreetFromParentNode.Name, currentNode.StreetFromParentNode.Path, ", Attributes:", currentNode.StreetFromParentNode.Attributes,") Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng, "")
		log.Println("Score:", currentNode.G)

		if currentNode.NodeID == endNode.NodeID {
			doneChannel <- *currentNode
			return
		}

		if concurrentNode != nil && closedList.Contains(concurrentNode) {
			finishNode := closedList.GetByKey(concurrentNode.NodeID)
			log.Println("Found a way:", finishNode.NodeID)
			log.Println("Found a way (Concurrent):", concurrentNode.NodeID)
			doneChannel <- *finishNode
			doneChannel <- *concurrentNode
			/*close(doneChannel)
			close(forwardChannel)
			close(backwardChannel) */
			return
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(*currentNode);

		for _, neighbor := range neighbors {

			neighbor.StreetFromParentNode.CorrectPath(currentNode)
			if closedList.Contains(neighbor) || !neighbor.Walkable  {
				continue
			}

			/*
			if openList.ContainsByKey(neighbor.NodeID) {
				neighbor = openList.GetByKey(neighbor.NodeID)
			} */

			neighbor.G = currentNode.G + DistanceFromLinePoint(neighbor.StreetFromParentNode.Path)
			neighbor.Heuristic = this.CalculateHeuristic(currentNode, neighbor, &endNode)
			//neighbor.Heuristic += int(float64(neighbor.Heuristic) * 0.15)

			neighbor.F = neighbor.G + neighbor.Heuristic
			neighbor.ParentNodes = currentNode
			log.Println("Possible Node:", neighbor.NodeID , "(",neighbor.StreetFromParentNode.Name,") Geometry:", neighbor.NodeGeometry.Lat, "," ,neighbor.NodeGeometry.Lng)
			log.Println("Street:", neighbor.StreetFromParentNode.Name, "Path:", neighbor.StreetFromParentNode.Path, " Attributes:", neighbor.StreetFromParentNode.Attributes)
			openList.Add(neighbor)
		}

		openList.Sort()
	}
}

package bbbikeng

import (
	"log"
)

type Route struct {

	way []Point
	detailed []*Path
	distance int
	time int
	nodes []*Node

	startNode Node
	endNode Node

	Preferences Preferences
	TrafficLights int
	Attributes []AttributeInterface

}

type Preferences struct {

	Speed int64
	Quality string
	Types string
	Greenways string
	AvoidUnlit bool
	AvoidLight bool
	IncludeFerries bool

}

func (this *Route) constructRoute(finalNode Node) {

	var parentNode *Node
	parentNode = &finalNode
	// gather all nodes

	var routeNodes []*Node
	for parentNode != nil {
		routeNodes = append(routeNodes, parentNode)
		parentNode = parentNode.ParentNodes
	}
	// reverse list for forward node. Keep it for backward node
	var tmpNodeList []*Node
	if !finalNode.flippedDirection {
		for i := len(routeNodes)-2; i >= 0; i-- {
			tmpNodeList = append(tmpNodeList, routeNodes[i])
		}
		routeNodes = tmpNodeList
	}

	for _, node := range routeNodes {
		this.nodes = append(this.nodes, node)
	}

	startNode := routeNodes[0]
	endNode := routeNodes[len(routeNodes)-1]

	for i := 0; i < len(routeNodes)-1; i++ {

		node := routeNodes[i]
		for _, attribute := range node.StreetFromParentNode.Attributes{
			this.Attributes = append(this.Attributes, attribute)
		}

		if node.NodeID == startNode.NodeID || node.NodeID == endNode.NodeID {

			// needs to be fixed for first node

			var index int
			if node.NodeID == startNode.NodeID {
				index = 0
			} else {
				index = len(this.way)-1
			}
			node.StreetFromParentNode.PathIndex = index
			this.detailed = append(this.detailed, &node.StreetFromParentNode)

			for i := 0; i < len(node.StreetFromParentNode.Path); i++ {
				this.way = append(this.way, node.StreetFromParentNode.Path[i])
			}


		} else {

			if this.detailed[len(this.detailed)-1].WayID != node.StreetFromParentNode.WayID && this.detailed[len(this.detailed)-1].Name != node.StreetFromParentNode.Name {
				node.StreetFromParentNode.PathIndex = len(this.way)-1
				this.detailed = append(this.detailed, &node.StreetFromParentNode)
			}
			for i := 1; i < len(node.StreetFromParentNode.Path); i++ {
				this.way = append(this.way, node.StreetFromParentNode.Path[i])
			}

		}

	}

	this.distance = DistanceFromLinePoint(this.way)

}

func (this *Route) StartRouting(startPoint Point, endPoint Point) {

	this.startNode = FindNearestNode(startPoint)
	this.endNode = FindNearestNode(endPoint)

	this.GetAStarRoute()

}


func (this *Route) GetAStarRoute() (){

	log.Println("StartNode:", this.startNode.NodeID , "(",this.startNode.StreetFromParentNode.Name,") Geometry:", this.startNode.NodeGeometry.Lat, "," ,this.startNode.NodeGeometry.Lng)
	log.Println("EndNode:", this.endNode.NodeID , "(",this.endNode.StreetFromParentNode.Name,") Geometry:", this.endNode.NodeGeometry.Lat, "," ,this.endNode.NodeGeometry.Lng)

	var openList = NewNodeSet()
	var closedList = NewNodeSet()
	openList.Add(&this.startNode)

	for openList.Length() > 0 {

		currentNode := openList.data[0]
		currentNode.flippedDirection = false
		log.Println("ParentNode:", currentNode.NodeID , "(",currentNode.StreetFromParentNode.ID, currentNode.StreetFromParentNode.Name, currentNode.StreetFromParentNode.Path, ", Attributes:", currentNode.StreetFromParentNode.Attributes,") Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng, "")
		log.Println("Score:", currentNode.G)
		if currentNode.NodeID == this.endNode.NodeID {
			this.constructRoute(*currentNode)
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


			neighbor.G = this.CalculateCosts(currentNode, neighbor)
			neighbor.Heuristic = this.CalculateHeuristic(neighbor, &this.endNode)
			neighbor.F = neighbor.G + neighbor.Heuristic
			neighbor.ParentNodes = currentNode
			log.Println("Possible Node:", neighbor.NodeID , "(",neighbor.StreetFromParentNode.Name,") Geometry:", neighbor.NodeGeometry.Lat, "," ,neighbor.NodeGeometry.Lng)
			log.Println("Street:", neighbor.StreetFromParentNode.Name, "Path:", neighbor.StreetFromParentNode.Path, " Attributes:", neighbor.StreetFromParentNode.Attributes)
			openList.Add(neighbor)
		}
		openList.Sort()
	}
}


func (this *Preferences) SetPreferedQuality(preferedQuality string){

	if preferedQuality != "Q0" && preferedQuality != "Q2" {
		log.Printf("Unknown Streetypes Setting")
		return
	}

	this.Quality = preferedQuality

}

func (this *Preferences) SetPreferedTypes(preferedTypes string){

	if preferedTypes != "N1" && preferedTypes != "N2" && preferedTypes != "H1" && preferedTypes != "H2" && preferedTypes != "N_RW" && preferedTypes != "N_RW1" {
		log.Printf("Unknown Prefered Types Setting")
		return
	}
	this.Types = preferedTypes

}

func (this *Preferences) SetPreferedGreen(preferedGreen string){

	if preferedGreen != "GR1"  && preferedGreen != "GR2" {
		log.Printf("Unknown Greenway Setting")
		return
	}
	this.Greenways = preferedGreen

}

func (this *Preferences) SetAvoidUnlit(avoidUnlit bool){
	this.AvoidUnlit = avoidUnlit
}

func (this *Preferences) SetAvoidTrafficLight(avoidTrafficLight bool){
	this.AvoidLight = avoidTrafficLight
}

func (this *Preferences) SetIncludeFerries(includeFerries bool){
	this.IncludeFerries = includeFerries
}

func (this *Preferences) SetPreferedSpeed(speed int64){
	this.Speed = speed
}

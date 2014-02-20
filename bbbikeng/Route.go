package bbbikeng

import (
	"log"
)

type Route struct {

	way []Point
	detailed []*Path
	distance float64
	time int
	nodes []*Node

	startNode Node
	endNode Node
	Preferences Preferences

	TrafficLights int

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
	for parentNode != nil {
		this.nodes = append(this.nodes, parentNode)
		parentNode = parentNode.ParentNodes
	}
	// reverse list
	var tmpNodeList []*Node
	for i := len(this.nodes)-2; i >= 0; i-- {
		tmpNodeList = append(tmpNodeList, this.nodes[i])
	}

	this.nodes = tmpNodeList

	startNode := this.nodes[0]
	endNode := this.nodes[len(this.nodes)-1]

	for _, node := range this.nodes {

		streetPath := flippPath(node)
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

			for i := 1; i < len(streetPath); i++ {
				this.way = append(this.way, streetPath[i])
			}

		}
	}

	this.distance = DistanceFromLinePoint(this.way)

}

func (this *Route) CalculateHeuristic(parentNode Node, neighborNode Node) (heuristic float64) {

	distanceToEnd := DistanceFromPointToPoint(neighborNode.NodeGeometry, this.endNode.NodeGeometry)
	//pathDistance := DistanceFromLinePoint(neighborNode.StreetFromParentNode.Path)
	return distanceToEnd

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

	openList.Add(this.startNode)

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
		log.Println("ParentNode:", currentNode.NodeID , "(",currentNode.StreetFromParentNode.ID, currentNode.StreetFromParentNode.Name,", Attributes:", currentNode.StreetFromParentNode.Attributes,") Geometry:", currentNode.NodeGeometry.Lat, "," ,currentNode.NodeGeometry.Lng, "")
		if currentNode.NodeID == this.endNode.NodeID {
			this.constructRoute(currentNode)
			return
		}

		openList.Remove(currentNode)
		closedList.Add(currentNode)

		neighbors := GetNeighborNodesFromNode(currentNode);

		for _, neighbor := range neighbors {

			if closedList.Contains(neighbor) || (!neighbor.Walkable && neighbor.NodeID != this.endNode.NodeID)  {
				continue
			}

			gScore := DistanceFromLinePoint(neighbor.StreetFromParentNode.Path)
			gScoreIsBest := false;

			if !openList.ContainsByKey(neighbor.NodeID) {
				gScoreIsBest = true;
				neighbor.Heuristic = this.CalculateHeuristic(currentNode, neighbor)
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
				log.Println("Street:", neighbor.StreetFromParentNode.Name, "Path:", neighbor.StreetFromParentNode.Path, " Attributes:", neighbor.StreetFromParentNode.Attributes)
				log.Println("Score:", neighbor.Heuristic)
				openList.Add(neighbor)
			}
		}
	}

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

func (this *Route) GetBBJson() (json BBJSON) {

	for _, segment := range this.way {
		var newPoint [2]float64
		newPoint[0] = segment.Lng
		newPoint[1] = segment.Lat
		json.Path = append(json.Path, newPoint)
	}

	for _, path := range this.detailed {
		var newInstruction BBJSONInstruction
		newInstruction.Name = path.Name
		newInstruction.Index = path.PathIndex
		newInstruction.Type = path.Type
//		newInstruction.Quality = path.Attributes.Quality
		json.Instruction = append(json.Instruction, newInstruction)
	}

	json.Preferences = this.Preferences
	json.Lights = this.TrafficLights
	json.Distance = int(this.distance * 1000.0)
	json.Response = true

	return json
}

func (this *Preferences) SetPreferedQuality(preferedQuality string){

	if preferedQuality != "Q0"  && preferedQuality != "Q2" {
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

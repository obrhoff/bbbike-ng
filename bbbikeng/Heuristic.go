package bbbikeng

import (
	"log"

)



/*func (this *Route) CalculateHeuristic(parentNode *Node, neighborNode *Node, endNode *Node) (heuristic int) {

	distanceToDestiny := DistanceFromPointToPoint(neighborNode.NodeGeometry, endNode.NodeGeometry)
	pathDistance := DistanceFromLinePoint(neighborNode.StreetFromParentNode.Path)
	score := 1.0

	sortedAttributes, attributesPerIndex, distancePerIndex := GetRelevantAttributes(parentNode, neighborNode)
	neighborNode.StreetFromParentNode.Attributes = sortedAttributes

	for key, attributeSegments := range attributesPerIndex {
		distanceFromSegment := distancePerIndex[key]
		weightOfTotal := (distanceFromSegment / pathDistance)
		segmentScore := 1.0
		for _, attribute := range attributeSegments{
			attr := *attribute
			segmentScore += attr.CalculateScore(&this.Preferences)
		}
		score *= (segmentScore * float64(weightOfTotal))
	}

	return  distanceToDestiny

} */

func (this *Route) CalculateCosts(parentNode *Node, neighborNode *Node) (heuristic int) {

	streetPathDistance := DistanceFromLinePoint(neighborNode.StreetFromParentNode.Path)
	streetPathDistance += parentNode.G

	attributesMap := GetRelevantAttributes(parentNode, neighborNode)
	log.Println("ATTRIBUTES:", attributesMap)

	return streetPathDistance

}

func (this *Route) CalculateHeuristic(neighborNode *Node, endNode *Node) (heuristic int) {
	return DistanceFromPointToPoint(neighborNode.NodeGeometry, endNode.NodeGeometry)
}

func GetRelevantAttributes (parentNode *Node, neighborNode *Node) (attributesPerIndex []map[string]*AttributeInterface){

	streetPath := neighborNode.StreetFromParentNode
	var relevantAttributes *[]AttributeInterface

	if parentNode.flippedDirection {
		relevantAttributes = &streetPath.FlippedAttribute
	} else {
		relevantAttributes = &streetPath.NormalAttribute
	}
	for _, globalAttribute := range streetPath.GlobalAttribute {
		*relevantAttributes = append(*relevantAttributes, globalAttribute)
	}

	log.Println("RELEVANT:", relevantAttributes)

	for _, attr := range *relevantAttributes {
		log.Println("Relevant Attribute", attr)
	}

	attributesPerIndex = make([]map[string]*AttributeInterface,0,len(streetPath.Path))

	/*	distancePerIndex = make(map([int]int)

	for i := 0; i < len(neighborNode.StreetFromParentNode.Path)-1; i++ {
		if i+1 <= len(neighborNode.StreetFromParentNode.Path)-1 {
			firstPoint := neighborNode.StreetFromParentNode.Path[i]
			secondPoint := neighborNode.StreetFromParentNode.Path[i+1]
			distancePerIndex[i] = DistanceFromPointToPoint(firstPoint, secondPoint)
		}
	} */

	for i :=0; i < len(streetPath.Path); i++ {
		attributesPerIndex = append(attributesPerIndex, make(map[string]*AttributeInterface));
	}

	for _, attribute := range *relevantAttributes {

		attributeGeometry := attribute.Geometry()
		firstGeometry := attributeGeometry[0]
		lastGeometry := attributeGeometry[len(attributeGeometry)-1]
		firstIndex := -1
		lastIndex := -1
		for pathIndex, pathSegment := range neighborNode.StreetFromParentNode.Path {
			if pathSegment.Compare(firstGeometry) {
				firstIndex = pathIndex
			} else if (pathSegment.Compare(lastGeometry)){
				lastIndex = pathIndex
			}
		}
		if lastIndex >= 0 && firstIndex >= 0 {
			for pathIndex := firstIndex; pathIndex < lastIndex; pathIndex++ {
				attributeMap := attributesPerIndex[pathIndex]
				var newKey string
				switch attribute.(type)  {
					case *CyclepathAttribute: {
						newKey = "CA"
					}
					case *GreenwayAttribute: {
						newKey = "GA"
					}
					case *QualityAttribute: {
						newKey = "QA"
					}
					case *UnlitAttribute: {
						newKey = "UA"
					}
					case *TrafficLightAttribute: {
						newKey = "TA"
					}
					case *HandicapAttribute: {
						newKey = "HA"
					}
				}
				attributeMap[newKey] = &attribute
			}
		}
	}

	return attributesPerIndex

}

func (this *CyclepathAttribute) CalculateScore (preference *Preferences) (score float64) {

	return score
}

func (this *GreenwayAttribute) CalculateScore (preference *Preferences) (score float64) {



	return score
}
func (this *QualityAttribute) CalculateScore (preference *Preferences) (score float64) {

	qualityPreference := preference.Quality
	attribute := this.Type()

	if qualityPreference == "Q0" {
		switch attribute {
			case "Q2":
				score = 0.5
			case "Q1":
				score = 0.75
			case "Q0":
				score = 1.0
		}
	}
	if qualityPreference == "Q2" {
		switch attribute {
		case "Q2":
			score = 0.75
		case "Q1":
			score = 1.0
		case "Q0":
			score = 1.0
		}
	}


	return score
}
func (this *UnlitAttribute) CalculateScore (preference *Preferences) (score float64) {


	return score
}
func (this *TrafficLightAttribute) CalculateScore (preference *Preferences) (score float64) {
	return score
}

func (this *HandicapAttribute) CalculateScore (preference *Preferences) (score float64) {
	return score
}

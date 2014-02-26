package bbbikeng

import (

)

func (this *Route) CalculateHeuristic(parentNode *Node, neighborNode *Node) (heuristic float64) {

	distanceToDestiny := DistanceFromPointToPoint(neighborNode.NodeGeometry, this.endNode.NodeGeometry)
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
		score *= (segmentScore * weightOfTotal)
	}

	return  distanceToDestiny

}

func GetRelevantAttributes (parentNode *Node, neighborNode *Node) (relevantAttributes []AttributeInterface, attributesPerIndex map[int][]*AttributeInterface, distancePerIndex map[int]float64){

	flipped := !parentNode.NodeGeometry.Compare(neighborNode.StreetFromParentNode.Path[0])
	if flipped {
		relevantAttributes = neighborNode.StreetFromParentNode.FlippedAttribute
	} else {
		relevantAttributes = neighborNode.StreetFromParentNode.NormalAttribute
	}
	for _, globalAttribute := range neighborNode.StreetFromParentNode.GlobalAttribute {
		relevantAttributes = append(relevantAttributes, globalAttribute)
	}

	attributesPerIndex = make(map[int][]*AttributeInterface)
	distancePerIndex = make(map[int]float64)

	for i := 0; i < len(neighborNode.StreetFromParentNode.Path)-1; i++ {
		if i+1 <= len(neighborNode.StreetFromParentNode.Path)-1 {
			firstPoint := neighborNode.StreetFromParentNode.Path[i]
			secondPoint := neighborNode.StreetFromParentNode.Path[i+1]
			distancePerIndex[i] = DistanceFromPointToPoint(firstPoint, secondPoint)
		}
	}

	for _, attribute := range relevantAttributes {

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
				attributesPerIndex[pathIndex] = append(attributesPerIndex[pathIndex], &attribute)
			}
		}
	}

	return relevantAttributes, attributesPerIndex, distancePerIndex

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

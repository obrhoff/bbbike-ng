package bbbikeng

import (
	"log"
)

func (this *Route) CalculateHeuristic(parentNode *Node, neighborNode *Node) (heuristic float64) {

	heuristic = DistanceFromPointToPoint(neighborNode.NodeGeometry, this.endNode.NodeGeometry)
	neighborNode.StreetFromParentNode.Attributes = GetRelevantAttributes(parentNode, neighborNode)
	log.Println("Calculating score for:", neighborNode.StreetFromParentNode.Attributes)
	for _, attribute := range neighborNode.StreetFromParentNode.Attributes {

		log.Println("Calculating Attributescore for:", attribute)
		heuristic += attribute.CalculateScore(&this.Preferences)


	}

	return


}

func GetRelevantAttributes (parentNode *Node, neighborNode *Node) (relevantAttributes []AttributeInterface){

	flipped := !parentNode.NodeGeometry.Compare(neighborNode.StreetFromParentNode.Path[0])
	if flipped {
		relevantAttributes = neighborNode.StreetFromParentNode.NormalAttribute
	} else {
		relevantAttributes = neighborNode.StreetFromParentNode.FlippedAttribute
	}
	for _, globalAttribute := range neighborNode.StreetFromParentNode.GlobalAttribute {
		relevantAttributes = append(relevantAttributes, globalAttribute)
	}
	return relevantAttributes

}

func (this *CyclepathAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("Cyclepath Score:")

	return score
}

func (this *GreenwayAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("Greenway Score:")


	return score
}
func (this *QualityAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("Quality Score:")


	return score
}
func (this *UnlitAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("Unlit Score:")


	return score
}
func (this *TrafficLightAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("TrafficLight Score:")


	return score
}

func (this *HandicapAttribute) CalculateScore (preference *Preferences) (score float64) {

	log.Println("Handicap Score:")


	return score
}

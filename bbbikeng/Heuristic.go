package bbbikeng

import (
	"log"
)

func (this *Route) CalculateHeuristic(parentNode Node, neighborNode Node) (heuristic float64) {

	heuristic = DistanceFromPointToPoint(neighborNode.NodeGeometry, this.endNode.NodeGeometry)

	log.Println("Calculating score for:", neighborNode)

	for _, attribute := range neighborNode.StreetFromParentNode.Attributes {

		if !attribute.Relevance() {
			continue
		}

		log.Println("Calculating Attributescore for:", attribute)
		heuristic += attribute.CalculateScore(&this.Preferences)


	}

	return heuristic

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

	log.Println("Unlit Score:")


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

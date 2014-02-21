package bbbikeng

func (this *Route) CalculateHeuristic(parentNode Node, neighborNode Node) (heuristic float64) {

	distanceToEnd := DistanceFromPointToPoint(neighborNode.NodeGeometry, this.endNode.NodeGeometry)
	for _, attribute := range neighborNode.StreetFromParentNode.Attributes {

		if !attribute.isValid {
			continue
		}

		switch attribute.Category {

			case "greenway": {

			}

			case "cyclepath": {

			}

			case "quality": {

			}

			case "unlit": {

			}

			case "trafficlight": {

			}

		}

	}

	return distanceToEnd

}

/*
func (this *Attribute) scoreForCyclepath (preferences *Preferences) (score float64) {

}

func scoreForQuality (qualityPreference string, qualityAttribute Attribute) (score float64) {

}


func scoreForGreenway (greenwayPreference string, greenwayAttribute Attribute) (score float64) {

}

func scoreForTrafficLight (avoidTrafficLight bool, trafficLights int) (score float64) {

}*/


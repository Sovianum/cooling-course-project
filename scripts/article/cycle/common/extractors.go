package common

import "github.com/Sovianum/turbocycle/core/graph"

type PowerNode interface {
	MassRateInput() graph.Port
	PowerOutput() graph.Port
}

func GetPower(node PowerNode) float64 {
	return node.PowerOutput().GetState().Value().(float64) * node.MassRateInput().GetState().Value().(float64)
}

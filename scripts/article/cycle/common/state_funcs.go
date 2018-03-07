package common

import "github.com/Sovianum/turbocycle/core/graph"

func CopyAll(p1s, p2s []graph.Port) {
	for i, p1 := range p1s {
		CopyState(p1, p2s[i])
	}
}

func CopyState(p1, p2 graph.Port) {
	p2.SetState(p1.GetState())
}

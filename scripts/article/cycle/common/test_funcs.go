package common

import "github.com/Sovianum/turbocycle/core/graph"

type NamedFloat struct {
	Val float64
	Name string
}

func TraceWithTags(ports []graph.Port, tags []string) []NamedFloat {
	vals := Trace(ports...)
	result := make([]NamedFloat, len(vals))
	for i, val := range vals {
		result[i] = NamedFloat{Val:val, Name:tags[i]}
	}
	return result
}

func Trace(ports ...graph.Port) []float64 {
	result := make([]float64, len(ports))
	for i, port := range ports {
		result[i] = port.GetState().Value().(float64)
	}
	return result
}

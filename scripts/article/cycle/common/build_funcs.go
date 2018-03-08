package common

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/library/schemes"
)

func GetMassRate(power float64, scheme schemes.Scheme, mrs nodes.MassRateSink) float64 {
	return schemes.GetMassRate(power, scheme) * mrs.MassRateInput().GetState().Value().(float64)
}

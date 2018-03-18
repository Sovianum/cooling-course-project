package common

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/library/schemes"
)

func GetMassRate(power float64, scheme schemes.Scheme, mrs nodes.MassRateSink) float64 {
	mr := schemes.GetMassRate(power, scheme)
	factor := mrs.MassRateInput().GetState().Value().(float64)
	return mr * factor
}

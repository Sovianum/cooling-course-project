package two_shafts

import (
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	tAtm           = 288
	pAtm           = 1e5
	sigmaInlet     = 0.98
	etaComp        = 0.86
	piComp         = 11
	tGas           = 1223
	tFuel          = 300
	sigmaBurn      = 0.99
	etaBurn        = 0.98
	initAlpha      = 3
	t0             = 300
	etaCompTurbine = 0.9
	lambdaOut      = 0.3
	dgCompTurbine  = -0.01
	etaM           = 0.99

	sigmaCompTurbinePipe = 0.98

	etaFreeTurbine               = 0.92
	dgFreeTurbine                = -0.01
	freeTurbinePressureLossSigma = 0.93

	precision = 0.05
)

func GetInitedTwoShaftsScheme() schemes.TwoShaftsScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var gasGenerator = compose.NewGasGeneratorNode(
		etaComp, piComp, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaCompTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaM, precision,
	)
	var compressorTurbinePipe = constructive.NewPressureLossNode(sigmaCompTurbinePipe)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
		pAtm,
		etaFreeTurbine, lambdaOut, precision,
		func(node constructive.TurbineNode) float64 {
			return dgFreeTurbine
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		freeTurbinePressureLossSigma,
	)

	return schemes.NewTwoShaftsScheme(gasSource, inletPressureDrop, gasGenerator, compressorTurbinePipe, freeTurbineBlock)
}

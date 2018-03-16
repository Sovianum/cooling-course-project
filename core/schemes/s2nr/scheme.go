package s2nr

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

const (
	tAtm           = 288
	pAtm           = 1e5
	sigmaInlet     = 0.98
	etaComp        = 0.86
	piComp         = 11
	tGas           = 1450
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

	regeneratorSigma = 0.8

	precision = 0.05
)

func GetInitedTwoShaftsRegeneratorScheme() schemes.TwoShaftsRegeneratorScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var turboCascade = compose.NewTurboCascadeNode(
		etaComp, piComp, etaCompTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return dgCompTurbine
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaM, precision,
	)
	var burner = constructive.NewBurnerNode(fuel.GetCH4(), tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0, precision, 1, nodes.DefaultN)
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
	var regenerator = constructive.NewRegeneratorNode(regeneratorSigma, precision)

	return schemes.NewTwoShaftsRegeneratorScheme(
		gasSource, inletPressureDrop, turboCascade, burner, compressorTurbinePipe, freeTurbineBlock, regenerator,
	)
}

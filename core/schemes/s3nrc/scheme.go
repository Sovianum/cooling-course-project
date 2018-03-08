package s3nrc

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
	tAtm = 288
	pAtm = 1e5

	sigmaInlet = 0.98

	etaMiddlePressureComp    = 0.86
	piCompTotal              = 30
	piCompFactor             = 0.18
	etaMiddlePressureTurbine = 0.9
	dgMiddlePressureTurbine  = 0.01
	etaMMiddleCascade        = 0.99

	etaHighPressureComp = 0.83

	tGas                   = 1223
	tFuel                  = 300
	sigmaBurn              = 0.99
	etaBurn                = 0.98
	initAlpha              = 3
	t0                     = 300
	etaHighPressureTurbine = 0.9
	lambdaOut              = 0.3
	dgHighPressureTurbine  = -0.01
	etaM                   = 0.99

	middlePressureCompressorPipeSigma = 0.98
	highPressureTurbinePipeSigma      = 0.98
	middlePressureTurbinePipeSigma    = 0.98

	etaFreeTurbine               = 0.92
	dgFreeTurbine                = -0.01
	freeTurbinePressureLossSigma = 0.93

	regeneratorSigma     = 0.8
	regeneratorPipeSigma = 0.98

	coolerTOut  = 320
	coolerSigma = 0.98

	precision = 0.05
)

func GetInitedThreeShaftsCoolRegeneratorScheme() schemes.ThreeShaftsCoolingRegeneratorScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var middlePressureCascade = compose.NewTurboCascadeNode(
		etaMiddlePressureComp, piCompTotal*piCompFactor,
		etaMiddlePressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return dgMiddlePressureTurbine
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaMMiddleCascade, precision,
	)
	var regenerativeGasGenerator = compose.NewRegenerativeGasGeneratorNode(
		etaHighPressureComp, 1/piCompFactor, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaHighPressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return dgHighPressureTurbine
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		regeneratorSigma, regeneratorPipeSigma, etaM, precision, 1, nodes.DefaultN,
	)
	var middlePressureCompressorPipe = constructive.NewPressureLossNode(middlePressureCompressorPipeSigma)
	var cooler = constructive.NewCoolerNode(coolerTOut, coolerSigma)
	var highPressureTurbinePipe = constructive.NewPressureLossNode(highPressureTurbinePipeSigma)
	var middlePressureTurbinePipe = constructive.NewPressureLossNode(middlePressureTurbinePipeSigma)
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

	return schemes.NewThreeShaftsCoolingRegeneratorScheme(
		gasSource, inletPressureDrop, middlePressureCascade, cooler, regenerativeGasGenerator,
		middlePressureCompressorPipe, highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
	)
}

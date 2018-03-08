package s3nb

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/compose"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/source"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
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

	tGas                  = 1450
	tFuel                 = 300
	sigmaBurn             = 0.99
	etaBurn               = 0.98
	initAlpha             = 3
	t0                    = 300
	etaCompTurbine        = 0.9
	lambdaOut             = 0.3
	dgHighPressureTurbine = -0.01
	etaM                  = 0.99

	middlePressureCompressorPipeSigma = 0.98
	highPressureTurbinePipeSigma      = 0.98
	middlePressureTurbinePipeSigma    = 0.98

	etaFreeTurbine               = 0.92
	dgFreeTurbine                = -0.01
	freeTurbinePressureLossSigma = 0.93

	midBurnTGas  = 1450
	midBurnEta = 0.98
	midBurnSigma = 0.93

	precision = 0.05
	relaxCoef = 0.1
	iterLimit = 100
)

func GetInitedThreeShaftsBurnScheme() schemes.ThreeShaftsBurnScheme {
	gasSource := source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	inletPressureDrop := constructive.NewPressureLossNode(sigmaInlet)
	middlePressureCascade := compose.NewTurboCascadeNode(
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
	gasGenerator := compose.NewGasGeneratorNode(
		etaHighPressureComp, 1/piCompFactor, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaCompTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return dgHighPressureTurbine
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaM, precision, relaxCoef, iterLimit,
	)
	middlePressureCompressorPipe := constructive.NewPressureLossNode(middlePressureCompressorPipeSigma)
	highPressureTurbinePipe := constructive.NewPressureLossNode(highPressureTurbinePipeSigma)
	middlePressureTurbinePipe := constructive.NewPressureLossNode(middlePressureTurbinePipeSigma)
	freeTurbineBlock := compose.NewFreeTurbineBlock(
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
	midBurner := constructive.NewBurnerNode(fuel.GetCH4(), midBurnTGas, tFuel, midBurnSigma, midBurnEta, initAlpha, t0, precision, relaxCoef, iterLimit)
	midBurner.SetName("MidBurner")

	return schemes.NewThreeShaftsBurnScheme(
		gasSource, inletPressureDrop, middlePressureCascade, gasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock, midBurner,
	)
}

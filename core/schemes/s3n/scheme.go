package s3n

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
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

	etaMiddlePressureComp = 0.84
	piCompTotal           = 30
	piCompFactor          = 0.18

	PiDiplomaTotal = 19.
	PiDiplomaHigh  = 4
	PiDiplomaLow   = PiDiplomaTotal / PiDiplomaHigh // proto value: 19/4

	etaMiddlePressureTurbine = 0.90
	etaMMiddleCascade        = 0.99

	etaHighPressureComp = 0.82

	tGas                   = 1450
	tFuel                  = 300
	sigmaBurn              = 0.99
	etaBurn                = 0.98
	initAlpha              = 3
	t0                     = 300
	etaHighPressureTurbine = 0.88
	lambdaOut              = 0.3
	etaMHighCascade        = 0.99

	middlePressureCompressorPipeSigma = 0.98
	highPressureTurbinePipeSigma      = 0.98
	middlePressureTurbinePipeSigma    = 0.98

	etaFreeTurbine               = 0.92
	dgFreeTurbine                = -0.01
	freeTurbinePressureLossSigma = 0.93

	hptCoolMassRate = -0.1
	lptCoolMassRate = 0

	hptLeakMassRate = -0.01
	lptLeakMassRate = -0.01

	precision = 0.05
)

func GetDiplomaInitedThreeShaftsScheme() schemes.ThreeShaftsScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var middlePressureCascade = compose.NewTurboCascadeNode(
		etaMiddlePressureComp, PiDiplomaLow,
		etaMiddlePressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return -lptLeakMassRate
		},
		func(node constructive.TurbineNode) float64 {
			return -lptCoolMassRate
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaMMiddleCascade, precision,
	)
	var gasGenerator = compose.NewGasGeneratorNode(
		etaHighPressureComp, PiDiplomaHigh, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaHighPressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return -hptLeakMassRate
		},
		func(node constructive.TurbineNode) float64 {
			return -hptCoolMassRate
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaMHighCascade, precision, 0.1, nodes.DefaultN,
	)
	var middlePressureCompressorPipe = constructive.NewPressureLossNode(middlePressureCompressorPipeSigma)
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

	return schemes.NewThreeShaftsScheme(
		gasSource, inletPressureDrop, middlePressureCascade, gasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
	)
}

func GetInitedThreeShaftsScheme() schemes.ThreeShaftsScheme {
	var gasSource = source.NewComplexGasSourceNode(gases.GetAir(), tAtm, pAtm, 1)
	var inletPressureDrop = constructive.NewPressureLossNode(sigmaInlet)
	var middlePressureCascade = compose.NewTurboCascadeNode(
		etaMiddlePressureComp, piCompTotal*piCompFactor,
		etaMiddlePressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaMMiddleCascade, precision,
	)
	var gasGenerator = compose.NewGasGeneratorNode(
		etaHighPressureComp, 1/piCompFactor, fuel.GetCH4(),
		tGas, tFuel, sigmaBurn, etaBurn, initAlpha, t0,
		etaHighPressureTurbine, lambdaOut,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		etaMHighCascade, precision, 1, nodes.DefaultN,
	)
	var middlePressureCompressorPipe = constructive.NewPressureLossNode(middlePressureCompressorPipeSigma)
	var highPressureTurbinePipe = constructive.NewPressureLossNode(highPressureTurbinePipeSigma)
	var middlePressureTurbinePipe = constructive.NewPressureLossNode(middlePressureTurbinePipeSigma)
	var freeTurbineBlock = compose.NewFreeTurbineBlock(
		pAtm,
		etaFreeTurbine, lambdaOut, precision,
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		func(node constructive.TurbineNode) float64 {
			return 0
		},
		freeTurbinePressureLossSigma,
	)

	return schemes.NewThreeShaftsScheme(
		gasSource, inletPressureDrop, middlePressureCascade, gasGenerator, middlePressureCompressorPipe,
		highPressureTurbinePipe, middlePressureTurbinePipe, freeTurbineBlock,
	)
}

func GetInitDF() InitDF {
	return InitDF{
		PAtm:      pAtm,
		TAtm:      tAtm,
		SigmaIn:   sigmaInlet,
		EtaLPC:    etaMiddlePressureComp,
		EtaHPC:    etaHighPressureComp,
		EtaLPT:    etaMiddlePressureTurbine,
		EtaHPT:    etaHighPressureTurbine,
		EtaFT:     etaFreeTurbine,
		EtaMLow:   etaMMiddleCascade,
		EtaMHigh:  etaMHighCascade,
		TGas:      tGas,
		TFuel:     tFuel,
		T0:        t0,
		SigmaBurn: sigmaBurn,
		EtaBurn:   etaBurn,
		SigmaLPC:  middlePressureCompressorPipeSigma,
		SigmaHPT:  highPressureTurbinePipeSigma,
		SigmaLPT:  middlePressureTurbinePipeSigma,
		SigmaFT:   freeTurbinePressureLossSigma,
		LambdaOut: lambdaOut,
	}
}

type InitDF struct {
	Ne        float64
	EtaR      float64
	PAtm      float64
	TAtm      float64
	SigmaIn   float64
	EtaLPC    float64
	EtaHPC    float64
	EtaHPT    float64
	EtaLPT    float64
	EtaFT     float64
	EtaMLow   float64
	EtaMHigh  float64
	TGas      float64
	TFuel     float64
	T0        float64
	SigmaBurn float64
	EtaBurn   float64
	SigmaLPC  float64
	SigmaHPT  float64
	SigmaLPT  float64
	SigmaFT   float64
	LambdaOut float64
}

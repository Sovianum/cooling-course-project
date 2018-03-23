package p2nr

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/schemes/s2nr"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	etaM       = 0.99
	cRpm0      = 10000
	cLambdaIn0 = 0.3

	ctID       = 0.5
	ctLambdaU0 = 0.3
	ctStageNum = 1

	ftID       = 0.7
	ftLambdaU0 = 0.3
	ftStageNum = 1

	payloadRpm0 = 3000

	velocityHotIn0        = 20
	velocityColdIn0       = 20
	hydraulicDiameterHot  = 1e-3
	hydraulicDiameterCold = 1e-3

	power = 16e6

	relaxCoef       = 1
	schemeRelaxCoef = 0.5
	iterNum         = 10000
	precision       = 1e-7
	schemePrecision = 1e-5

	startPi   = 5
	piStep    = 0.5
	piStepNum = 12
)

func SolveParametric(pScheme free2n.DoubleShaftRegFreeScheme) (Data2nr, error) {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return Data2nr{}, pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, schemePrecision,
	)

	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, common.DetailedLog2Shaft),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), precision, relaxCoef, 10000)
	if sErr != nil {
		return Data2nr{}, sErr
	}

	data := NewData2nr()
	for i := 0; i != 20; i++ {
		data.Load(pScheme)
		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 20)

		r := 1.
		_, sErr = vSolver.Solve(vSolver.GetInit(), precision, r, 1000)
		if sErr != nil {
			return Data2nr{}, sErr
		}
		fmt.Println(i)
	}
	return data, nil
}

func GetParametric(scheme schemes.TwoShaftsRegeneratorScheme) (free2n.DoubleShaftRegFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	solveErr := network.Solve(relaxCoef, 2, iterNum, schemePrecision)
	if solveErr != nil {
		return nil, solveErr
	}

	return getParametricScheme(scheme), nil
}

func getParametricScheme(scheme schemes.TwoShaftsRegeneratorScheme) free2n.DoubleShaftRegFreeScheme {
	builder := NewBuilder(
		scheme,
		power,
		cRpm0, cLambdaIn0,
		ctID, ctLambdaU0, ctStageNum,
		ftID, ftLambdaU0, ftStageNum,
		payloadRpm0, etaM,
		velocityHotIn0, velocityColdIn0,
		hydraulicDiameterHot, hydraulicDiameterCold,
		constructive.DefaultNuFunc, constructive.DefaultNuFunc,
		constructive.CounterTDrop,
		schemePrecision, schemeRelaxCoef, iterNum,
	)
	return builder.Build()
}

func OptimizeScheme(scheme schemes.TwoShaftsRegeneratorScheme, data core.SingleCompressorData) {
	optPi := 0.
	maxEta := -1.
	for i := range data.Efficiency {
		if data.Efficiency[i] > maxEta {
			optPi = data.Pi[i]
			maxEta = data.Efficiency[i]
		}
	}
	scheme.Compressor().SetPiStag(optPi)
}

func GetSchemeData(scheme schemes.TwoShaftsRegeneratorScheme) (core.SingleCompressorData, error) {
	points, err := io.GetTwoShaftSchemeData(scheme, power, startPi, piStep, piStepNum)
	if err != nil {
		return core.SingleCompressorData{}, err
	}
	return core.ConvertSingleCompressorDataPoint(points), nil
}

func GetScheme(piStag float64) schemes.TwoShaftsRegeneratorScheme {
	scheme := s2nr.GetInitedTwoShaftsRegeneratorScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

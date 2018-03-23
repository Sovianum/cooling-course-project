package p2n

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/schemes/s2n"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
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

	power     = 16e6
	relaxCoef = 0.5
	iterNum   = 10000
	precision = 1e-4

	piStag = 10

	startPi   = 8
	piStep    = 0.5
	piStepNum = 20
)

func SolveParametric(pScheme free2n.DoubleShaftFreeScheme) (Data2n, error) {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return Data2n{}, pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, precision,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, common.DetailedLog2Shaft),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 1, 10000)
	if sErr != nil {
		return Data2n{}, sErr
	}

	data := NewData2n()
	step := 10.
	r := 1.
	for i := 0; i != 30; i++ {
		//if i == 15 {
		//	step = 1
		//}
		//if i == 19 {
		//	step = 0.5
		//}

		data.Load(pScheme)
		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - step)

		fmt.Println(i)
		_, sErr = vSolver.Solve(vSolver.GetInit(), 1e-5, r, 10000)
		if sErr != nil {
			return Data2n{}, sErr
		}
	}
	return data, nil
}

func GetParametric(scheme schemes.TwoShaftsScheme) (free2n.DoubleShaftFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	solveErr := network.Solve(relaxCoef, 2, iterNum, precision/10)
	if solveErr != nil {
		return nil, solveErr
	}

	return getParametricScheme(scheme), nil
}

func getParametricScheme(scheme schemes.TwoShaftsScheme) free2n.DoubleShaftFreeScheme {
	builder := NewBuilder(
		scheme, power, cRpm0, cLambdaIn0,
		ctID, ctLambdaU0, ctStageNum,
		ftID, ftLambdaU0, ftStageNum,
		payloadRpm0, etaM, precision, relaxCoef, iterNum,
	)
	return builder.Build()
}

func OptimizeScheme(scheme schemes.TwoShaftsScheme, data core.SingleCompressorData) {
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

func GetSchemeData(scheme schemes.TwoShaftsScheme) (core.SingleCompressorData, error) {
	points, err := io.GetTwoShaftSchemeData(scheme, power, startPi, piStep, piStepNum)
	if err != nil {
		return core.SingleCompressorData{}, err
	}
	return core.ConvertSingleCompressorDataPoint(points), nil
}

func GetScheme(piStag float64) schemes.TwoShaftsScheme {
	scheme := s2n.GetInitedTwoShaftsScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

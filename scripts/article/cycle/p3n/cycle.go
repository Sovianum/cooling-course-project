package p3n

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3n"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/core"
)

const (
	power = 16e6

	lpcRpm0 = 6000
	hpcRpm0 = 10000

	lambdaIn0 = 0.3

	hptInletDiameter = 0.5
	hptLambdaU0      = 0.3
	hptStageNum      = 1

	lptInletDiameter = 1
	lptLambdaU0      = 0.3
	lptStageNum      = 3

	ftInletDiameter = 1
	ftLambdaU0      = 0.3
	ftStageNum      = 3

	payloadRpm0 = 3000

	lpEtaM = 0.99
	hpEtaM = 0.99

	relaxCoef = 0.5
	iterNum   = 10000
	precision = 1e-5

	schemePrecision = 1e-5

	lpcPiStag = 4
	hpcPiStag = 2.5

	startPi   = 8
	piStep    = 0.5
	piStepNum = 30
)

func SolveParametric(pScheme free3n.ThreeShaftFreeScheme) (Data3n, error) {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return Data3n{}, pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, schemePrecision,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, common.DetailedLog3Shaft),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 1, 10000)
	if sErr != nil {
		return Data3n{}, sErr
	}

	data := NewData3n()
	for i := 0; i != 15; i++ {
		data.Load(pScheme)

		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 1.
		_, sErr = vSolver.Solve(vSolver.GetInit(), 1e-5, r, 1000)
		if sErr != nil {
			return Data3n{}, sErr
		}
		fmt.Println(i)
	}
	return data, nil
}

func GetParametric(scheme schemes.ThreeShaftsScheme) (free3n.ThreeShaftFreeScheme, error) {
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

func getParametricScheme(scheme schemes.ThreeShaftsScheme) free3n.ThreeShaftFreeScheme {
	builder := NewBuilder(
		scheme, power,
		lpcRpm0, hpcRpm0,
		lambdaIn0,
		lptInletDiameter, lptLambdaU0, lptStageNum,
		hptInletDiameter, hptLambdaU0, hptStageNum,
		ftInletDiameter, ftLambdaU0, ftStageNum,
		payloadRpm0,
		lpEtaM, hpEtaM,
		precision, relaxCoef, iterNum,
	)
	return builder.Build()
}

func OptimizeScheme(scheme schemes.ThreeShaftsScheme, data core.DoubleCompressorData) {
	optPiLow := 0.
	optPiHigh := 0.
	maxEta := -1.
	for i := range data.Efficiency {
		if data.Efficiency[i] > maxEta {
			optPiLow = data.PiLow[i]
			optPiHigh = data.PiHigh[i]
			maxEta = data.Efficiency[i]
		}
	}
	scheme.HPC().SetPiStag(optPiHigh)
	scheme.LPC().SetPiStag(optPiLow)
}

func GetSchemeData(scheme schemes.ThreeShaftsScheme) (core.DoubleCompressorData, error) {
	points, err := io.GetThreeShaftsSchemeData(scheme, power, startPi, piStep, piStepNum, 0.1, 0.1, 8)
	if err != nil {
		return core.DoubleCompressorData{}, err
	}
	return core.ConvertDoubleCompressorDataPoints(points), nil
}

func GetScheme(piStagLow, piStagHigh float64) schemes.ThreeShaftsScheme {
	scheme := s3n.GetInitedThreeShaftsScheme()
	scheme.LPC().SetPiStag(piStagLow)
	scheme.HPC().SetPiStag(piStagHigh)
	return scheme
}

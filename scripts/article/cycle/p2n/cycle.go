package p2n

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/s2n"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"os"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
)

const (
	etaM       = 0.99
	cRpm0      = 10000
	cLambdaIn0 = 0.3

	ctID       = 0.3
	ctLambdaU0 = 0.3
	ctStageNum = 1

	ftID       = 0.5
	ftLambdaU0 = 0.3
	ftStageNum = 1

	payloadRpm0 = 3000

	power     = 16e6
	relaxCoef = 1
	iterNum   = 10000
	precision = 1e-4

	piStag = 10
)

func SolveParametric(pScheme free2n.DoubleShaftFreeScheme) error {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, 0.1,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, common.DetailedLog),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 1, 10000)
	if sErr != nil {
		return sErr
	}

	data := NewData2n()
	for i := 0; i != 17; i++ {
		data.Load(pScheme)
		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 1.
		_, sErr = vSolver.Solve(vSolver.GetInit(), 1e-5, r, 1000)
		if sErr != nil {
			break
		}
		fmt.Println(i)
	}

	b, _ := json.Marshal(data)
	f, e := os.Create("/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/notebooks/data/2n.json")
	if e != nil {
		return e
	}

	_, e = f.WriteString(string(b))
	return e
}

func GetParametric(scheme schemes.TwoShaftsScheme) (free2n.DoubleShaftFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	solveErr := network.Solve(relaxCoef, 2, iterNum, precision / 10)
	if solveErr != nil {
		return nil, solveErr
	}

	return get2nParametricScheme(scheme), nil
}

func get2nParametricScheme(scheme schemes.TwoShaftsScheme) free2n.DoubleShaftFreeScheme {
	builder := NewBuilder(
		scheme, power, cRpm0, cLambdaIn0,
		ctID, ctLambdaU0, ctStageNum,
		ftID, ftLambdaU0, ftStageNum,
		payloadRpm0, etaM, precision, relaxCoef, iterNum,
	)
	return builder.Build()
}

func GetScheme(piStag float64) schemes.TwoShaftsScheme {
	scheme := s2n.GetInitedTwoShaftsScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

package p3nb

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3nb"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"os"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
)

const (
	power = 16e6
	t0    = 300
	p0    = 1e5

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

	relaxCoef = 1
	iterNum   = 1000
	precision = 1e-4

	schemePrecision = 1e-4

	lpcPiStag = 4
	hpcPiStag = 2.5
)

func SolveParametric(pScheme free3n.ThreeShaftFreeScheme) error {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, precision,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, common.DetailedLog3ShaftMidBurn),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 0.5, 10000)
	if sErr != nil {
		return sErr
	}

	data := NewData3n()
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
	f, _ := os.Create("/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/notebooks/data/3nb.json")
	f.WriteString(string(b))
	return nil
}

func GetParametric(scheme schemes.ThreeShaftsBurnScheme) (free3n.ThreeShaftBurnFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	solveErr := network.Solve(relaxCoef, 2, iterNum, schemePrecision)
	if solveErr != nil {
		return nil, solveErr
	}

	return get3nbBurnParametricScheme(scheme), nil
}

func get3nbBurnParametricScheme(scheme schemes.ThreeShaftsBurnScheme) free3n.ThreeShaftBurnFreeScheme {
	builder := NewBuilder(
		scheme, power, t0, p0,
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

func GetScheme(piStagLow, piStagHigh float64) schemes.ThreeShaftsBurnScheme {
	scheme := s3nb.GetInitedThreeShaftsBurnScheme()
	scheme.LPC().SetPiStag(piStagLow)
	scheme.HPC().SetPiStag(piStagHigh)
	return scheme
}

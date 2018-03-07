package p2nr

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/two_shafts_regenerator"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"gonum.org/v1/gonum/mat"
	"os"
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

	t0 = 300
	p0 = 1e5

	velocityHotIn0        = 20
	velocityColdIn0       = 20
	hydraulicDiameterHot  = 1e-3
	hydraulicDiameterCold = 1e-3

	power = 20e6

	relaxCoef       = 1
	schemeRelaxCoef = 0.5
	iterNum         = 10000
	precision       = 0.01
	schemePrecision = 0.1
)

func SolveParametric(pScheme free2n.DoubleShaftRegFreeScheme) error {
	network, pErr := pScheme.GetNetwork()
	if pErr != nil {
		return pErr
	}

	sysCall := variator.SysCallFromNetwork(
		network, pScheme.Assembler().GetVectorPort(),
		relaxCoef, 2, iterNum, schemePrecision,
	)
	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, func(iterNum int, precision float64, residual *mat.VecDense) {
			result := fmt.Sprintf("i: %d\t", iterNum)
			result += fmt.Sprintf("precision: %f\t", precision)
			result += fmt.Sprintf(
				"ggmr: %.5f\t ggPower: %.5f\t ftMR: %.5f\t ftPower: %.5f\t ftPressure: %.5f\t ggTemp: %.5f\t",
				residual.At(0, 0),
				residual.At(1, 0),
				residual.At(2, 0),
				residual.At(3, 0),
				residual.At(4, 0),
				residual.At(5, 0),
			)
			result += fmt.Sprintf("residual: %f", mat.Norm(residual, 2))
			fmt.Println(result)
		}),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), precision, relaxCoef, 10000)
	if sErr != nil {
		return sErr
	}

	data := NewData2nr()
	for i := 0; i != 17; i++ {
		data.Load(pScheme)
		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 1.
		_, sErr = vSolver.Solve(vSolver.GetInit(), precision, r, 1000)
		if sErr != nil {
			break
		}
		fmt.Println(i)
	}

	b, _ := json.Marshal(data)
	f, _ := os.Create("/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/notebooks/data/2nr.json")
	f.WriteString(string(b))
	return nil
}

func GetParametric(scheme schemes.TwoShaftsRegeneratorScheme) (free2n.DoubleShaftRegFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	converged, solveErr := network.Solve(relaxCoef, 2, iterNum, schemePrecision)
	if solveErr != nil {
		return nil, solveErr
	}
	if !converged {
		return nil, fmt.Errorf("failed to converge")
	}

	return getParametricScheme(scheme), nil
}

func getParametricScheme(scheme schemes.TwoShaftsRegeneratorScheme) free2n.DoubleShaftRegFreeScheme {
	builder := NewBuilder(
		scheme,
		power, t0, p0,
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

func GetScheme(piStag float64) schemes.TwoShaftsRegeneratorScheme {
	scheme := two_shafts_regenerator.GetInitedTwoShaftsRegeneratorScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

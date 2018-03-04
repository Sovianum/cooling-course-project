package p2nr

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"os"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/cooling-course-project/core/schemes/two_shafts_regenerator"
	"gonum.org/v1/gonum/mat"
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

	velocityHotIn0 = 20
	velocityColdIn0 = 20
	hydraulicDiameterHot = 1e-3
	hydraulicDiameterCold = 1e-3

	power     = 20e6

	relaxCoef = 0.1
	schemeRelaxCoef = 0.5
	iterNum   = 10000
	precision = 0.01
	schemePrecision = 0.01
)

func SolveParametric(pScheme free2n.DoubleShaftFreeScheme) error {
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

	var data struct {
		T       []float64 `json:"t"`
		P       []float64 `json:"p"`
		G       []float64 `json:"g"`
		PiC     []float64 `json:"pi_c"`
		PiTC    []float64 `json:"pi_tc"`
		PiF     []float64 `json:"pi_f"`
		GNormTC []float64 `json:"g_norm"`
		GNormTF []float64 `json:"g_norm_tf"`
		GNormC  []float64 `json:"g_norm_c"`
		RpmTC   []float64 `json:"rpm_tc"`
		RpmFT   []float64 `json:"rpm_ft"`
	}

	for i := 0; i != 17; i++ {
		t := pScheme.TemperatureSource().GetTemperature()
		labour := pScheme.FreeTurbine().PowerOutput().GetState().Value().(float64)
		massRate := pScheme.Compressor().MassRate()
		normMassRateTC := pScheme.CompressorTurbine().NormMassRate()
		normMassRateFT := pScheme.FreeTurbine().NormMassRate()

		data.T = append(data.T, t)
		data.P = append(data.P, labour*massRate/1e6)
		data.G = append(data.G, massRate)
		data.PiC = append(data.PiC, pScheme.Compressor().PiStag())
		data.PiTC = append(data.PiTC, pScheme.CompressorTurbine().PiTStag())
		data.PiF = append(data.PiF, pScheme.FreeTurbine().PiTStag())
		data.GNormC = append(data.GNormC, pScheme.Compressor().NormMassRate())
		data.GNormTC = append(data.GNormTC, normMassRateTC)
		data.GNormTF = append(data.GNormTF, normMassRateFT)
		data.RpmTC = append(data.RpmTC, pScheme.CompressorTurbine().RPMInput().GetState().Value().(float64))
		data.RpmFT = append(data.RpmFT, pScheme.FreeTurbine().RPMInput().GetState().Value().(float64))

		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 0.1
		//if i >= 8 {
		//	r = 0.1
		//}
		_, sErr = vSolver.Solve(vSolver.GetInit(), precision, r, 1000)
		if sErr != nil {
			break
		}
		fmt.Println(i)
	}

	b, _ := json.Marshal(data)
	f, _ := os.Create("/tmp/dat")
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

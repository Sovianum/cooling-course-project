package p2n

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/two_shafts"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
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

	power     = 20e6
	relaxCoef = 0.1
	iterNum   = 10000
	precision = 0.01

	piStag = 10
)

func SolveParametric(pScheme free2n.DoubleShaftFreeScheme) error {
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
		newton.NewUniformNewtonSolverGen(1e-5, newton.DefaultLog),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 0.5, 10000)
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

		r := 1.
		//if i >= 8 {
		//	r = 0.1
		//}
		_, sErr = vSolver.Solve(vSolver.GetInit(), 1e-5, r, 1000)
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

func GetParametric(scheme schemes.TwoShaftsScheme) (free2n.DoubleShaftFreeScheme, error) {
	network, err := scheme.GetNetwork()
	if err != nil {
		return nil, err
	}
	converged, solveErr := network.Solve(relaxCoef, 2, iterNum, precision)
	if solveErr != nil {
		return nil, solveErr
	}
	if !converged {
		return nil, fmt.Errorf("failed to converge")
	}

	return get2nParametricScheme(scheme), nil
}

func get2nParametricScheme(scheme schemes.TwoShaftsScheme) free2n.DoubleShaftFreeScheme {
	builder := NewBuilder(
		scheme, power, t0, p0, cRpm0, cLambdaIn0,
		ctID, ctLambdaU0, ctStageNum,
		ftID, ftLambdaU0, ftStageNum,
		payloadRpm0, etaM, precision, relaxCoef, iterNum,
	)
	return builder.Build()
}

func GetScheme(piStag float64) schemes.TwoShaftsScheme {
	scheme := two_shafts.GetInitedTwoShaftsScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

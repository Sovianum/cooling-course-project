package p3n

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/turbocycle/core/math/solvers/newton"
	"github.com/Sovianum/turbocycle/core/math/variator"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"os"
)

const (
	power = 20e6
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

	relaxCoef = 0.1
	iterNum   = 10000
	precision = 0.01

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
		newton.NewUniformNewtonSolverGen(1e-5, newton.DefaultLog),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), 1e-6, 0.5, 10000)
	if sErr != nil {
		return sErr
	}

	var data struct {
		T []float64 `json:"t"`
		P []float64 `json:"p"`
		G []float64 `json:"g"`

		PiLPC []float64 `json:"pi_lpc"`
		PiHPC []float64 `json:"pi_hpc"`

		PiHPT []float64 `json:"pi_hpt"`
		PiLPT []float64 `json:"pi_hpt"`
		PiFT  []float64 `json:"pi_ft"`

		GNormHPT []float64 `json:"g_norm_hpt"`
		GNormLPT []float64 `json:"g_norm_lpt"`
		GNormFT  []float64 `json:"g_norm_ft"`

		RpmHPT []float64 `json:"rpm_hpt"`
		RpmLPT []float64 `json:"rpm_lpt"`
		RpmFT  []float64 `json:"rpm_ft"`
	}

	for i := 0; i != 17; i++ {
		t := pScheme.TemperatureSource().GetTemperature()
		labour := pScheme.FT().PowerOutput().GetState().Value().(float64)
		massRate := pScheme.FT().MassRateInput().GetState().Value().(float64)
		normMassRateHPT := pScheme.HPT().NormMassRate()
		normMassRateLPT := pScheme.LPT().NormMassRate()
		normMassRateFT := pScheme.FT().NormMassRate()

		data.T = append(data.T, t)
		data.P = append(data.P, labour*massRate/1e6)
		data.G = append(data.G, massRate)

		data.PiLPC = append(data.PiLPC, pScheme.LPC().PiStag())
		data.PiHPC = append(data.PiHPC, pScheme.HPC().PiStag())

		data.PiHPT = append(data.PiHPT, pScheme.HPT().PiTStag())
		data.PiLPT = append(data.PiLPT, pScheme.LPT().PiTStag())
		data.PiFT = append(data.PiFT, pScheme.FT().PiTStag())

		data.GNormHPT = append(data.GNormHPT, normMassRateHPT)
		data.GNormLPT = append(data.GNormLPT, normMassRateLPT)
		data.GNormFT = append(data.GNormFT, normMassRateFT)

		data.RpmHPT = append(data.RpmHPT, pScheme.HPT().RPMInput().GetState().Value().(float64))
		data.RpmLPT = append(data.RpmLPT, pScheme.LPT().RPMInput().GetState().Value().(float64))
		data.RpmFT = append(data.RpmFT, pScheme.FT().RPMInput().GetState().Value().(float64))

		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 0.5
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
	f, _ := os.Create("/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/notebooks/data/3n.csv")
	f.WriteString(string(b))
	return nil
}

func GetParametric(scheme schemes.ThreeShaftsScheme) (free3n.ThreeShaftFreeScheme, error) {
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

	return get3nParametricScheme(scheme), nil
}

func get3nParametricScheme(scheme schemes.ThreeShaftsScheme) free3n.ThreeShaftFreeScheme {
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

func GetScheme(piStagLow, piStagHigh float64) schemes.ThreeShaftsScheme {
	scheme := three_shafts.GetInitedThreeShaftsScheme()
	scheme.LPC().SetPiStag(piStagLow)
	scheme.HPC().SetPiStag(piStagHigh)
	return scheme
}

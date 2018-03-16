package p2nr

import (
	"encoding/json"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/s2nr"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
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

	type monitorData struct {
		CTIn  common.FloatArr `json:"ct_in"`
		CTOut common.FloatArr `json:"ct_out"`
		CL    common.FloatArr `json:"cl"`
		CMR   common.FloatArr `json:"cmr"`
		CEta common.FloatArr `json:"c_eta"`

		TTIn  common.FloatArr `json:"tt_in"`
		TTOut common.FloatArr   `json:"tt_out"`
		TL    common.FloatArr `json:"tl"`
		TMR   common.FloatArr `json:"tmr"`
		TEta common.FloatArr `json:"t_eta"`
	}

	md0 := monitorData{}
	md1 := monitorData{}

	ggPowerMonitor := func(i int) {
		if i == -1 {
			md0.CTIn.Append(pScheme.Compressor().TStagIn())
			md0.CTOut.Append(pScheme.Compressor().TStagOut())
			md0.CL.Append(pScheme.Compressor().PowerOutput().GetState().Value().(float64))
			md0.CMR.Append(pScheme.Compressor().MassRateInput().GetState().Value().(float64))
			md0.CEta.Append(pScheme.Compressor().Eta())

			md0.TTIn.Append(pScheme.CompressorTurbine().TStagIn())
			md0.TTOut.Append(pScheme.CompressorTurbine().TStagOut())
			md0.TL.Append(pScheme.CompressorTurbine().PowerOutput().GetState().Value().(float64))
			md0.TMR.Append(pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64))
			md0.TEta.Append(pScheme.CompressorTurbine().Eta())
		}
		if i == 0 {
			md1.CTIn.Append(pScheme.Compressor().TStagIn())
			md1.CTOut.Append(pScheme.Compressor().TStagOut())
			md1.CL.Append(pScheme.Compressor().PowerOutput().GetState().Value().(float64))
			md1.CMR.Append(pScheme.Compressor().MassRateInput().GetState().Value().(float64))
			md1.CEta.Append(pScheme.Compressor().Eta())

			md1.TTIn.Append(pScheme.CompressorTurbine().TStagIn())
			md1.TTOut.Append(pScheme.CompressorTurbine().TStagOut())
			md1.TL.Append(pScheme.CompressorTurbine().PowerOutput().GetState().Value().(float64))
			md1.TMR.Append(pScheme.CompressorTurbine().MassRateInput().GetState().Value().(float64))
			md1.TEta.Append(pScheme.CompressorTurbine().Eta())
		}
	}

	vSolver := variator.NewVariatorSolver(
		sysCall, pScheme.Variators(),
		newton.NewUniformNewtonSolverGen(1e-5, func(iterNum int, precision float64, residual *mat.VecDense) {
			common.DetailedLog2Shaft(iterNum, precision, residual)
			ggPowerMonitor(iterNum)
		}),
	)

	_, sErr := vSolver.Solve(vSolver.GetInit(), precision, relaxCoef, 10000)
	if sErr != nil {
		return sErr
	}
	//pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() + 50)
	//vSolver.Solve(vSolver.GetInit(), precision, relaxCoef, 10000)

	data := NewData2nr()
	for i := 0; i != 30; i++ {
		data.Load(pScheme)
		pScheme.TemperatureSource().SetTemperature(pScheme.TemperatureSource().GetTemperature() - 10)

		r := 1.
		_, sErr = vSolver.Solve(vSolver.GetInit(), precision, r, 1000)
		if sErr != nil {
			break
		}
		fmt.Println(i)
	}

	monitorB0, _ := json.Marshal(md0)
	monitorF0, _ := os.Create("/tmp/monitor0.json")
	monitorF0.WriteString(string(monitorB0))

	monitorB1, _ := json.Marshal(md1)
	monitorF1, _ := os.Create("/tmp/monitor1.json")
	monitorF1.WriteString(string(monitorB1))

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

func GetScheme(piStag float64) schemes.TwoShaftsRegeneratorScheme {
	scheme := s2nr.GetInitedTwoShaftsRegeneratorScheme()
	scheme.Compressor().SetPiStag(piStag)
	return scheme
}

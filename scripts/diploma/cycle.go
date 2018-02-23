package diploma

import (
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/io"
)

func saveCycleTemplate(scheme schemes.ThreeShaftsScheme) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cycleTemplate,
		buildDir+"/"+cycleOut,
	)
	var df = dataframes.NewThreeShaftsDF(power, etaR, scheme)
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func solveParticularScheme(scheme schemes.ThreeShaftsScheme, lowPiStag, highPiStag float64) {
	scheme.LowPressureCompressor().SetPiStag(lowPiStag)
	scheme.HighPressureCompressor().SetPiStag(highPiStag)
	network, netErr := scheme.GetNetwork()
	if netErr != nil {
		panic(netErr)
	}

	if converged, err := network.Solve(relaxCoef, 2, iterNum, precision); !converged || err != nil {
		if err != nil {
			panic(err)
		}
		if !converged {
			panic(fmt.Errorf("not converged"))
		}
	}
}

func saveVariantTemplate(schemeData []core.DoubleCompressorDataPoint) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+variantTemplate,
		buildDir+"/"+variantOut,
	)
	var df = dataframes.VariantDF{
		MaxEta:    core.EtaOptimalPoint(schemeData),
		MaxLabour: core.LabourOptimalPoint(schemeData),
		PiLow:     lowPiStag,
		PiHigh:    highPiStag,
		PiTotal:   totalPiStag,
	}
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func saveInputTemplates() {
	var cycleInputInserter = templ.NewDataInserter(
		templatesDir+"/"+cycleInputTemplate,
		buildDir+"/"+cycleInputOut,
	)
	var projectInputInserter = templ.NewDataInserter(
		templatesDir+"/"+projectInputTemplate,
		buildDir+"/"+projectInputOut,
	)

	var df = three_shafts.GetInitDF()
	df.Ne = power
	df.EtaR = etaR
	if err := cycleInputInserter.Insert(df); err != nil {
		panic(err)
	}
	if err := projectInputInserter.Insert(df); err != nil {
		panic(err)
	}
}

func saveSchemeData(data []core.DoubleCompressorDataPoint) {
	var matrix = make([][]float64, len(data))
	for i, point := range data {
		matrix[i] = point.ToArray()
	}

	if err := profiling.SaveMatrix(dataDir+"3n.csv", matrix); err != nil {
		panic(err)
	}
}

func getSchemeData(scheme schemes.ThreeShaftsScheme) []core.DoubleCompressorDataPoint {
	if data, err := io.GetThreeShaftsSchemeData(
		scheme,
		power/etaR,
		startPi, piStep, piStepNum,
		startPiFactor, piFactorStep, piFactorStepNum,
	); err != nil {
		panic(err)
	} else {
		return data
	}
}

func getScheme(lowPiStag, highPiStag float64) schemes.ThreeShaftsScheme {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	scheme.LowPressureCompressor().SetPiStag(lowPiStag)
	scheme.HighPressureCompressor().SetPiStag(highPiStag)
	return scheme
}

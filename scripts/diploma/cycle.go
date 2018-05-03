package diploma

import (
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3n"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/library/schemes"
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
	scheme.LPC().SetPiStag(lowPiStag)
	scheme.HPC().SetPiStag(highPiStag)
	network, netErr := scheme.GetNetwork()
	if netErr != nil {
		panic(netErr)
	}

	if err := network.Solve(relaxCoef, 2, iterNum, precision); err != nil {
		if err != nil {
			panic(err)
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
		PiLow:     s3n.PiDiplomaLow,
		PiHigh:    s3n.PiDiplomaHigh,
		PiTotal:   s3n.PiDiplomaTotal,
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

	var df = s3n.GetInitDF()
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
		0.3, 0.1, 4,
	); err != nil {
		panic(err)
	} else {
		return data
	}
}

func getScheme(lowPiStag, highPiStag float64) schemes.ThreeShaftsScheme {
	var scheme = s3n.GetDiplomaInitedThreeShaftsScheme()
	scheme.LPC().SetPiStag(lowPiStag)
	scheme.HPC().SetPiStag(highPiStag)
	return scheme
}

package main

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/core/midline"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
)

const (
	power     = 16e6
	relaxCoef = 0.1
	iterNum   = 100
	precision = 0.05

	startPi   = 8
	piStep    = 0.1
	piStepNum = 150

	startPiFactor   = 0.15
	piFactorStep    = 0.1
	piFactorStepNum = 8

	totalPiStag = 12
	piFactor    = 0.5

	dataDir = "postprocessing/notebooks/cycle/data/"

	templatesDir = "postprocessing/templates"
	buildDir     = "build"

	cycleTemplate = "cycle_calc_template.tex"
	cycleOut      = "cycle_calc.tex"

	rootTemplate = "root.tex"
	rootOut      = "root.tex"

	stageTemplate = "mean_line_calc_template.tex"
	stageOut      = "mean_line_calc.tex"
)

func main() {
	var lowPiStag = totalPiStag * piFactor
	var highPiStag = 1 / piFactor
	var scheme = getScheme(lowPiStag, highPiStag)

	//saveSchemeData(scheme)
	solveParticularScheme(scheme, lowPiStag, highPiStag)
	saveCycleTemplate(scheme)

	var stage = midline.GetInitedStageNode(scheme)
	solveParticularStage(stage)
	saveStageTemplate(stage)

	saveRootTemplate()
}

func saveRootTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir + "/" + rootTemplate,
		buildDir + "/" + rootOut,
	)
	if err := inserter.Insert(nil); err != nil {
		panic(err)
	}
}

func saveStageTemplate(stage nodes.TurbineStageNode) {
	var inserter = templ.NewDataInserter(
		templatesDir + "/" + stageTemplate,
		buildDir + "/" + stageOut,
	)
	var df, err = dataframes.NewStageDF(stage)
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func solveParticularStage(stage nodes.TurbineStageNode) {
	if err := stage.Process(); err != nil {
		panic(err)
	}
}

func saveCycleTemplate(scheme schemes.ThreeShaftsScheme) {
	var inserter = templ.NewDataInserter(
		templatesDir + "/" + cycleTemplate,
		buildDir + "/" + cycleOut,
	)
	var df = dataframes.NewThreeShaftsDF(power, scheme)
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func solveParticularScheme(scheme schemes.ThreeShaftsScheme, lowPiStag, highPiStag float64) {
	scheme.LowPressureCompressor().SetPiStag(lowPiStag)
	scheme.HighPressureCompressor().SetPiStag(highPiStag)
	if converged, err := scheme.GetNetwork().Solve(relaxCoef, iterNum, precision); !converged || err != nil {
		if err != nil {
			panic(err)
		}
		if !converged {
			panic(fmt.Errorf("not converged"))
		}
	}
}

func saveSchemeData(scheme schemes.ThreeShaftsScheme)  {
	if err := io.SaveThreeShaftsSchemeData(
		scheme,
		power,
		startPi, piStep, piStepNum,
		startPiFactor, piFactorStep, piFactorStepNum,
		dataDir+"3n.csv",
	); err != nil {
		panic(err)
	}
}

func getScheme(lowPiStag, highPiStag float64) schemes.ThreeShaftsScheme {
	var scheme = three_shafts.GetInitedThreeShaftsScheme()
	scheme.LowPressureCompressor().SetPiStag(lowPiStag)
	scheme.HighPressureCompressor().SetPiStag(highPiStag)
	return scheme
}

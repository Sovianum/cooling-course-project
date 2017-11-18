package main

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/midline"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/library/schemes"
	"os/exec"
	"os"
)

const (
	power     = 16e6
	relaxCoef = 0.1
	iterNum   = 100
	precision = 0.05

	startPi   = 7
	piStep    = 0.1
	piStepNum = 200

	startPiFactor   = 0.15
	piFactorStep    = 0.1
	piFactorStepNum = 8

	totalPiStag = 12
	piFactor    = 0.5

	templatesDir = "postprocessing/templates"

	buildDir     = "build"
	dataDir      = "build/data/"
	imgDir       = "build/img"

	cycleTemplate = "cycle_calc_template.tex"
	cycleOut      = "cycle_calc.tex"

	rootTemplate = "root.tex"
	rootOut      = "root.tex"

	titleTemplate = "title.tex"
	titleOut      = "title.tex"

	stageTemplate = "mean_line_calc_template.tex"
	stageOut      = "mean_line_calc.tex"

	cooling1Template = "cooling_calc1_template.tex"
	cooling1Out      = "cooling_calc1.tex"

	cooling2Template = "cooling_calc2_template.tex"
	cooling2Out      = "cooling_calc2.tex"

	plotterPath = "plot_all.py"

	cycleFileName = "3n.csv"
)

func main() {
	io.PrepareDirectories(
		buildDir, dataDir, imgDir,
	)

	var lowPiStag = totalPiStag * piFactor
	var highPiStag = 1 / piFactor
	var scheme = getScheme(lowPiStag, highPiStag)

	saveSchemeData(scheme)
	solveParticularScheme(scheme, lowPiStag, highPiStag)
	saveCycleTemplate(scheme)

	var stage = midline.GetInitedStageNode(scheme)
	solveParticularStage(stage)
	saveStageTemplate(stage)

	saveCooling1Template()
	saveCooling2Template()

	saveRootTemplate()
	saveTitleTemplate()

	buildPlots()
}

func buildPlots() {
	var arguments = []string{
		imgDir,
		dataDir + "/" + cycleFileName,
	}

	var cmd = exec.Command("./plot_all.py", arguments...)
	cmd.Stdout = os.Stdout
	var err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func saveRootTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+rootTemplate,
		buildDir+"/"+rootOut,
	)
	if err := inserter.Insert(nil); err != nil {
		panic(err)
	}
}

func saveTitleTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+titleTemplate,
		buildDir+"/"+titleOut,
	)
	if err := inserter.Insert(nil); err != nil {
		panic(err)
	}
}

func saveCooling2Template() {
	var geomDF = dataframes.TProfileGeomDF{}
	var gasDF = dataframes.TProfileGasDF{
		LengthPSArr:   []float64{1, 1, 1},
		AlphaAirPSArr: []float64{1, 1, 1},
		AlphaGasPSArr: []float64{1, 1, 1},
		TAirPSArr:     []float64{1, 1, 1},

		LengthSSArr:   []float64{1, 1, 1},
		AlphaAirSSArr: []float64{1, 1, 1},
		AlphaGasSSArr: []float64{1, 1, 1},
		TAirSSArr:     []float64{1, 1, 1},
	}
	var calcDF = dataframes.TProfileCalcDF{
		Geom: geomDF,
		Gas:  gasDF,
	}

	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cooling2Template,
		buildDir+"/"+cooling2Out,
	)

	//var df, err = gapCalcDF, nil
	var df = calcDF
	var err error = nil
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func saveCooling1Template() {
	var geomDF = dataframes.GapGeometryDF{}
	var metalDF = dataframes.GapMetalDF{}
	var gasDF = dataframes.GapGasDF{
		AirMassRate: []float64{1, 1, 1},
		DCoef:       []float64{1, 1, 1},
		EpsCoef:     []float64{1, 1, 1},
		AirGap:      []float64{1, 1, 1},
	}
	var gapCalcDF = dataframes.GapCalcDF{
		geomDF, metalDF, gasDF,
	}

	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cooling1Template,
		buildDir+"/"+cooling1Out,
	)

	//var df, err = gapCalcDF, nil
	var df = gapCalcDF
	var err error = nil
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func saveStageTemplate(stage nodes.TurbineStageNode) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+stageTemplate,
		buildDir+"/"+stageOut,
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
		templatesDir+"/"+cycleTemplate,
		buildDir+"/"+cycleOut,
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

func saveSchemeData(scheme schemes.ThreeShaftsScheme) {
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

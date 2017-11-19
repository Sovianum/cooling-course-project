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
	"github.com/Sovianum/cooling-course-project/postprocessing/builder"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
	"math"
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

	buildDir     = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/build"
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

	cycleData     = "3n.csv"
	statorMidData = "stator_mid.csv"

	inletAngleData  = "inlet_angle.csv"
	outletAngleData = "outlet_angle.csv"

	hPointNum = 50
)

func main() {
	io.PrepareDirectories(
		buildDir, dataDir, imgDir,
	)

	var lowPiStag = totalPiStag * piFactor
	var highPiStag = 1 / piFactor
	var scheme = getScheme(lowPiStag, highPiStag)

	//saveSchemeData(scheme)
	solveParticularScheme(scheme, lowPiStag, highPiStag)
	saveCycleTemplate(scheme)

	var stage = midline.GetInitedStageNode(scheme)
	solveParticularStage(stage)
	saveStageTemplate(stage)

	var statorProfiler = getStatorProfiler(stage)
	saveProfiles(
		statorProfiler,
		[]float64{0, 0.5, 1.0},
		[]string{"stator_root.csv", "stator_mid.csv", "stator_top.csv"},
		false,
	)

	var rotorProfiler = getRotorProfiler(stage)
	saveAngleData(rotorProfiler, func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle {
		var triangle = profiler.InletTriangle(hRel)
		return triangle
	}, inletAngleData)
	saveAngleData(rotorProfiler, func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle {
		var triangle = profiler.OutletTriangle(hRel)
		return triangle
	}, outletAngleData)
	saveProfiles(
		rotorProfiler,
		[]float64{0, 0.5, 1},
		[]string{"rotor_root.csv", "rotor_mid.csv", "rotor_top.csv"},
		true,
	)

	//saveCooling1Template()
	//saveCooling2Template()
	//
	//saveRootTemplate()
	//saveTitleTemplate()
	//
	buildPlots()
	//buildReport()
	////cleanup()
}

func cleanup() {
	var cmd = exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("cd %s && rm !(*.pdf)", buildDir),
	)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func buildReport() {
	if err := builder.BuildLatex(buildDir, rootOut); err != nil {
		panic(err)
	}
}

func buildPlots() {
	var arguments = []string{
		imgDir,
		dataDir,
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

func saveProfiles(profiler profilers.Profiler, hRelArr []float64, dataNames []string, isRotor bool) {
	var profileArr = make([]profiles.BladeProfile, len(hRelArr))
	for i, hRel := range hRelArr {
		profileArr[i] = profiles.NewBladeProfileFromProfiler(
			hRel,
			0.01, 0.01,
			0.2, 0.2,
			profiler,
		)
	}

	var installationAngleArr = make([]float64, len(hRelArr))
	for i, hRel := range hRelArr {
		installationAngleArr[i] = profiler.InstallationAngle(hRel)
	}

	var segments = make([]geom.Segment, len(hRelArr))
	for i, profile := range profileArr {
		if isRotor {
			profile.Transform(geom.Reflection(0))
		}
		profile.Transform(geom.Translation(mat.NewVecDense(2, []float64{-1, 0})))
		if !isRotor {
			profile.Transform(geom.Rotation(installationAngleArr[i] - math.Pi))
		} else {
			profile.Transform(geom.Rotation(-installationAngleArr[i]))
		}

		segments[i] = profiles.CircularSegment(profile)
	}

	for i := range hRelArr {
		var coordinates = geom.GetCoordinates(common.LinSpace(0, 1, 200), segments[i])
		if err := profiling.SaveMatrix(dataDir + "/" + dataNames[i], coordinates); err != nil {
			panic(err)
		}
	}
}

func saveAngleData(profiler profilers.Profiler, triangleExtractor func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle, filename string)  {
	var hRelArr = common.LinSpace(0, 1, hPointNum)

	var angleArr = make([][]float64, hPointNum)
	for i, hRel := range hRelArr {
		var triangle = triangleExtractor(hRel, profiler)
		angleArr[i] = make([]float64, 3)

		angleArr[i][0] = hRel
		angleArr[i][1] = triangle.Alpha()
		angleArr[i][2] = triangle.Beta()
	}

	if err := profiling.SaveMatrix(dataDir + "/" + filename, angleArr); err != nil {
		panic(err)
	}
}

func getRotorProfiler(stage nodes.TurbineStageNode) profilers.Profiler {
	var pack = stage.GetDataPack()
	var profiler = profiling.GetInitedRotorProfiler(
		stage.StageGeomGen().RotorGenerator(),
		pack.RotorInletTriangle,
		pack.RotorOutletTriangle,
	)
	return profiler
}

func getStatorProfiler(stage nodes.TurbineStageNode) profilers.Profiler {
	var pack = stage.GetDataPack()
	var profiler = profiling.GetInitedStatorProfiler(
		stage.StageGeomGen().StatorGenerator(),
		stage.VelocityInput().GetState().(states.VelocityPortState).Triangle,
		pack.RotorInletTriangle,
	)
	return profiler
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

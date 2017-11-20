package main

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/cooling-course-project/core/midline"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/builder"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"gonum.org/v1/gonum/mat"
	"math"
	"os"
	"os/exec"
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

	totalPiStag = 16
	lowPiStag   = 16 / 2
	highPiStag  = 2

	templatesDir = "postprocessing/templates"

	buildDir = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/build"
	dataDir  = "build/data/"
	imgDir   = "build/img"

	cycleInputTemplate = "cycle_input_data_template.tex"
	cycleInputOut      = "cycle_input_data.tex"

	variantTemplate = "variant_template.tex"
	variantOut      = "variant.tex"

	cycleTemplate = "cycle_calc_template.tex"
	cycleOut      = "cycle_calc.tex"

	rootTemplate = "root.tex"
	rootOut      = "root.tex"

	titleTemplate = "title.tex"
	titleOut      = "title.tex"

	stageTemplate = "mean_line_calc_template.tex"
	stageOut      = "mean_line_calc.tex"

	profilingTemplate = "profiling_template.tex"
	profilingOut      = "profiling.tex"

	cooling1Template = "cooling_calc1_template.tex"
	cooling1Out      = "cooling_calc1.tex"

	cooling2Template = "cooling_calc2_template.tex"
	cooling2Out      = "cooling_calc2.tex"

	inletAngleData  = "inlet_angle.csv"
	outletAngleData = "outlet_angle.csv"

	hPointNum = 50
)

func main() {
	io.PrepareDirectories(
		buildDir, dataDir, imgDir,
	)

	var scheme = getScheme(lowPiStag, highPiStag)

	saveCycleInputTemplate()

	//var schemeData = getSchemeData(scheme)
	//saveSchemeData(schemeData)
	//saveVariantTemplate(schemeData)

	solveParticularScheme(scheme, lowPiStag, highPiStag)
	saveCycleTemplate(scheme)

	var stage = midline.GetInitedStageNode(scheme)
	solveParticularStage(stage)
	saveStageTemplate(stage)

	var statorProfiler = getStatorProfiler(stage)
	saveProfiles(
		statorProfiler,
		stage.StageGeomGen().StatorGenerator(),
		[]float64{0, 0.5, 1.0},
		[][]string{
			{"rotor_root_1.csv", "rotor_root_2.csv"},
			{"rotor_mid_1.csv", "rotor_mid_2.csv"},
			{"rotor_top_1.csv", "rotor_top_2.csv"},
		},
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
		stage.StageGeomGen().RotorGenerator(),
		[]float64{0, 0.5, 1},
		[][]string{
			{"stator_root_1.csv", "stator_root_2.csv"},
			{"stator_mid_1.csv", "stator_mid_2.csv"},
			{"stator_top_1.csv", "stator_top_2.csv"},
		},
		true,
	)

	saveProfilingTemplate()

	saveCooling1Template()
	saveCooling2Template()

	saveRootTemplate()
	saveTitleTemplate()

	buildPlots()
	buildReport()
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

func saveProfilingTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+profilingTemplate,
		buildDir+"/"+profilingOut,
	)
	if err := inserter.Insert(nil); err != nil {
		panic(err)
	}
}

func saveProfiles(
	profiler profilers.Profiler,
	geomGen geometry.BladingGeometryGenerator,
	hRelArr []float64,
	dataNames [][]string,
	isRotor bool,
) {
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
	var tRelArr = make([]float64, len(hRelArr))
	var tArr = common.LinSpace(0, 1, 200)

	for i, hRel := range hRelArr {
		installationAngleArr[i] = profiler.InstallationAngle(hRel)
		tRelArr[i] = geometry.TRel(hRel, geomGen)
	}

	var coordinatesArr = make([][][][]float64, len(hRelArr))
	for i, profile := range profileArr {
		coordinatesArr[i] = make([][][]float64, 2)

		if isRotor {
			profile.Transform(geom.Reflection(0))
		}
		profile.Transform(geom.Translation(mat.NewVecDense(2, []float64{-1, 0})))
		if !isRotor {
			profile.Transform(geom.Rotation(installationAngleArr[i] - math.Pi))
		} else {
			profile.Transform(geom.Rotation(-installationAngleArr[i]))
		}

		coordinatesArr[i][0] = geom.GetCoordinates(tArr, profiles.CircularSegment(profile))

		profile.Transform(geom.Translation(mat.NewVecDense(2, []float64{
			tRelArr[i] * profiler.InstallationAngle(0.5), 0,
		})))
		coordinatesArr[i][1] = geom.GetCoordinates(tArr, profiles.CircularSegment(profile))
	}

	for i := range hRelArr {
		for j := 0; j != 2; j++ {
			if err := profiling.SaveMatrix(dataDir+"/"+dataNames[i][j], coordinatesArr[i][j]); err != nil {
				panic(err)
			}
		}
	}
}

func saveAngleData(
	profiler profilers.Profiler,
	triangleExtractor func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle,
	filename string,
) {
	var hRelArr = common.LinSpace(0, 1, hPointNum)

	var angleArr = make([][]float64, hPointNum)
	for i, hRel := range hRelArr {
		var triangle = triangleExtractor(hRel, profiler)
		angleArr[i] = make([]float64, 3)

		angleArr[i][0] = hRel
		angleArr[i][1] = triangle.Alpha()
		angleArr[i][2] = triangle.Beta()
	}

	if err := profiling.SaveMatrix(dataDir+"/"+filename, angleArr); err != nil {
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

func saveCycleInputTemplate() {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cycleInputTemplate,
		buildDir+"/"+cycleInputOut,
	)
	var df = three_shafts.GetInitDF()
	df.Ne = power
	if err := inserter.Insert(df); err != nil {
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
		power,
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

package diploma

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core/midall/inited"
	"github.com/Sovianum/cooling-course-project/core/schemes/s3n"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/builder"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"os"
	"os/exec"
)

const (
	power     = 16e6
	etaR      = 0.98
	relaxCoef = 0.1
	iterNum   = 100
	precision = 0.05

	startPi   = 10
	piStep    = 0.5
	piStepNum = 100

	templatesDir = "postprocessing/templates"

	buildDir = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/build"
	dataDir  = "build/data/"
	imgDir   = "build/img"

	projectInputTemplate = "project_input_data_template.tex"
	projectInputOut      = "project_input_data.tex"

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

	turbineStageTemplate = "mean_line_calc_template.tex"
	turbineStageOut      = "mean_line_calc.tex"

	compressorStageTemplate = "compressor_calc_template.tex"
	compressorStageOut      = "compressor_calc.tex"

	lpcTotalTableTemplate = "lpc_total_table_template.tex"
	lpcTotalTableOut      = "lpc_total_table.tex"

	hpcTotalTableTemplate = "hpc_total_table_template.tex"
	hpcTotalTableOut      = "hpc_total_table.tex"

	turbineTotalTableTemplate = "turbine_total_table_template.tex"
	turbineTotalTableOut      = "turbine_total_table.tex"

	profilingTemplate = "profiling_template.tex"
	profilingOut      = "profiling.tex"

	cooling1Template = "cooling_calc1_template.tex"
	cooling1Out      = "cooling_calc1.tex"

	cooling2Template = "cooling_calc2_template.tex"
	cooling2Out      = "cooling_calc2.tex"

	cooling2PSData = "cooling_2_ps.json"
	cooling2SSData = "cooling_2_ss.json"

	inletAngleData  = "inlet_angle.csv"
	outletAngleData = "outlet_angle.csv"

	hPointNum       = 50
	coolAirMassRate = 0.04
	theta0          = 500
	gapWidth        = 1e-3

	velocityCoef = 0.98
	massRateCoef = 0.98

	coolingBladeLength = 40e-3
	coolingHoleNum     = 20

	dInlet = 2.2e-3
)

func Entry() {
	io.PrepareDirectories(
		buildDir, dataDir, imgDir,
	)

	saveInputTemplates()

	scheme := getScheme(s3n.PiDiplomaLow, s3n.PiDiplomaHigh)
	schemeData := getSchemeData(scheme)
	saveSchemeData(schemeData)
	saveVariantTemplate(schemeData)

	solveParticularScheme(scheme, s3n.PiDiplomaLow, s3n.PiDiplomaHigh)
	saveCycleTemplate(scheme)

	saveCompressorStageTemplate()
	saveCompressorTotalTableTemplates()

	initedMachines, err := inited.GetInitedStagedNodes()
	if err != nil {
		panic(err)
	}
	stage := initedMachines.HPT.Stages()[0]
	saveTurbineStageTemplate(stage)
	saveTurbineTotalTableTemplates()

	statorProfiler := getStatorProfiler(stage)
	saveProfiles(
		statorProfiler,
		stage.StageGeomGen().StatorGenerator(),
		[]float64{0, 0.5, 1.0},
		[][]string{
			{"stator_root_1.csv", "stator_root_2.csv"},
			{"stator_mid_1.csv", "stator_mid_2.csv"},
			{"stator_top_1.csv", "stator_top_2.csv"},
		},
		false,
	)

	rotorProfiler := getRotorProfiler(stage)
	fmt.Println(fmt.Sprintf("profile %.1f", 0.5), getProfileMsg(rotorProfiler, 0.5))
	fmt.Println(fmt.Sprintf("profile %.1f", 0.0), getProfileMsg(rotorProfiler, 0.0))
	fmt.Println(fmt.Sprintf("profile %.1f", 1.0), getProfileMsg(rotorProfiler, 1.0))

	saveAngleData(rotorProfiler, func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle {
		triangle := profiler.InletTriangle(hRel)
		return triangle
	}, inletAngleData)
	saveAngleData(rotorProfiler, func(hRel float64, profiler profilers.Profiler) states.VelocityTriangle {
		triangle := profiler.OutletTriangle(hRel)
		return triangle
	}, outletAngleData)
	saveProfiles(
		rotorProfiler,
		stage.StageGeomGen().RotorGenerator(),
		[]float64{0, 0.5, 1},
		[][]string{
			{"rotor_root_1.csv", "rotor_root_2.csv"},
			{"rotor_mid_1.csv", "rotor_mid_2.csv"},
			{"rotor_top_1.csv", "rotor_top_2.csv"},
		},
		true,
	)

	inletGasProfiler, outletGasProfiler := getGasProfilers(stage, rotorProfiler)
	fmt.Println(profilers.Reactivity(0, 0.5, inletGasProfiler, outletGasProfiler))
	fmt.Println(profilers.Reactivity(0.5, 0.5, inletGasProfiler, outletGasProfiler))
	fmt.Println(profilers.Reactivity(1, 0.5, inletGasProfiler, outletGasProfiler))

	saveProfilingTemplate()

	statorMidProfile := profiles.NewBladeProfileFromProfiler(
		0.5,
		0.01, 0.01,
		0.2, 0.2,
		statorProfiler,
	)
	stagePack := stage.GetDataPack()
	statorMidProfile.Transform(geom.Scale(geometry.ChordProjection(stagePack.StageGeometry.StatorGeometry())))

	gapCalculator := getGapCalculator(stage, statorMidProfile)
	gapPack := gapCalculator.GetPack(coolAirMassRate)

	gapCalcDF := getGapDF(common.LinSpace(0.01, 0.10, 10), gapCalculator)
	saveCooling1Template(gapCalcDF)

	psTemperatureSystem := getPSConvFilmTemperatureSystem(
		gapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{3e-3, 0.4e-3},
			{27e-3, 0.5e-3},
			{37e-3, 0.5e-3},
		},
	)
	psSolution := psTemperatureSystem.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(psSolution, cooling2PSData)

	ssTemperatureSystem := getSSConvFilmTemperatureSystem(
		gapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{3e-3, 0.4e-3},
			{16e-3, 0.25e-3},
			{22e-3, 0.25e-3},
			{27e-3, 0.3e-3},
			{31e-3, 0.35e-3},
			{36.5e-3, 0.40e-3},
		},
	)
	ssSolution := ssTemperatureSystem.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(ssSolution, cooling2SSData)

	tempProfileDF := getTempProfileDF(gapCalcDF, stage, statorMidProfile, psSolution, ssSolution)
	saveCooling2Template(tempProfileDF)

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
	var cmd = exec.Command("./plot_all.py", imgDir, dataDir)
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

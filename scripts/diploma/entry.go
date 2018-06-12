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

	cooling2NoFrontPSData = "cooling_2_no_front_ps.json"
	cooling2NoFrontSSData = "cooling_2_no_front_ss.json"

	cooling2FrontPSData = "cooling_2_front_ps.json"
	cooling2FrontSSData = "cooling_2_front_ss.json"

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
	if err := copyPassiveFiles(); err != nil {
		panic(err)
	}

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
	fmt.Println("stator")
	fmt.Println(fmt.Sprintf("profile %.1f", 0.0), getProfileMsg(statorProfiler, 0.0))
	fmt.Println(fmt.Sprintf("profile %.1f", 0.5), getProfileMsg(statorProfiler, 0.5))
	fmt.Println(fmt.Sprintf("profile %.1f", 1.0), getProfileMsg(statorProfiler, 1.0))

	fmt.Println("rotor")
	fmt.Println(fmt.Sprintf("profile %.1f", 0.0), getProfileMsg(rotorProfiler, 0.0))
	fmt.Println(fmt.Sprintf("profile %.1f", 0.5), getProfileMsg(rotorProfiler, 0.5))
	fmt.Println(fmt.Sprintf("profile %.1f", 1.0), getProfileMsg(rotorProfiler, 1.0))

	for _, hRel := range []float64{0, 0.5, 1} {
		fmt.Println(fmt.Sprintf("triangle %.1f\n", hRel), getTrianglesMsg(
			rotorProfiler.InletTriangle(hRel), rotorProfiler.OutletTriangle(hRel),
		))
	}
	rotorGeom := stage.GetDataPack().StageGeometry.RotorGeometry()
	fmt.Println(
		rotorGeom.InnerProfile().Diameter(0)*1000,
		rotorGeom.MeanProfile().Diameter(0)*1000,
		rotorGeom.OuterProfile().Diameter(0)*1000,
	)

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

	noFrontGapPack := gapCalculator.GetPack(coolAirMassRate)
	gapCalcDF := getGapDF(common.LinSpace(0.01, 0.10, 10), gapCalculator)
	saveCooling1Template(gapCalcDF)

	psTemperatureSystemNoFront := getPSConvFilmTemperatureSystem(
		coolAirMassRate,
		noFrontGapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{4e-3, 0.45e-3},
			{18e-3, 0.4e-3},
			{30e-3, 0.5e-3},
			{37e-3, 0.40e-3},
		},
	)
	psSolutionNoFront := psTemperatureSystemNoFront.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(psSolutionNoFront, cooling2NoFrontPSData)

	ssTemperatureSystemNoFront := getSSConvFilmTemperatureSystem(
		coolAirMassRate,
		noFrontGapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{7e-3, 0.45e-3},
			{22e-3, 0.25e-3},
			{27e-3, 0.25e-3},
			{32e-3, 0.32e-3},
			{38e-3, 0.35e-3},
			{43e-3, 0.45e-3},
		},
	)
	ssSolutionNoFront := ssTemperatureSystemNoFront.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(ssSolutionNoFront, cooling2NoFrontSSData)

	minCoolAirMassRate := coolAirMassRate * 0.91
	frontGapPack := gapCalculator.GetPack(minCoolAirMassRate)
	tempProfileDF := getTempProfileDF(gapCalcDF, stage, statorMidProfile, psSolutionNoFront, ssSolutionNoFront)
	saveCooling2Template(tempProfileDF)

	psTemperatureSystemFront := getPSConvFilmTemperatureSystem(
		minCoolAirMassRate,
		frontGapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{0, 0.15e-3},
			{10e-3, 0.30e-3},
			{18e-3, 0.30e-3},
			{25e-3, 0.55e-3},
			{36.5e-3, 0.53e-3},
		},
	)
	psSolutionFront := psTemperatureSystemFront.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(psSolutionFront, cooling2FrontPSData)

	ssTemperatureSystemFront := getSSConvFilmTemperatureSystem(
		minCoolAirMassRate,
		frontGapPack.AlphaGas,
		stage,
		statorMidProfile,
		[]SlitGeom{
			{0, 0.115e-3},
			{18e-3, 0.25e-3},
			{24e-3, 0.25e-3},
			{30e-3, 0.3e-3},
			{35e-3, 0.55e-3},
			{43.5e-3, 0.55e-3},
		},
	)
	ssSolutionFront := ssTemperatureSystemFront.Solve(0, theta0, 1, 0.001)
	saveCoolingSolution(ssSolutionFront, cooling2FrontSSData)

	saveRootTemplate()
	saveTitleTemplate()

	buildPlots()
	buildReport()
}

func copyPassiveFiles() error {
	imgNames := []string{
		"cost.png",
		"cycle_2n_opt.png", "cycle_2n_part.png", "cycle_2n_scheme.png",
		"cycle_3n_opt.png", "cycle_3n_part.png", "cycle_3n_scheme.png",
		"cycle_2nr_opt.png", "cycle_2nr_part.png", "cycle_2nr_scheme.png",
		"cycle_eta_comparison.png",
		"ecology_bc.png", "ecology_cloud.png",
		"ecology_plan.png", "ecology_result.png",
		"lock_control.png", "profile_shape_control.png",
		"profile_thk_control.png",
	}
	templateNames := []string{
		"ecology.tex", "economics.tex", "ending.tex",
		"referat.tex", "technology.tex", "intro.tex",
		"literature.tex",
	}

	imgSrc := "postprocessing/media/img/"
	templateSrc := "postprocessing/templates/"

	for _, name := range imgNames {
		if err := io.CopyFile(imgSrc+name, imgDir+"/"+name); err != nil {
			return err
		}
	}

	for _, name := range templateNames {
		if err := io.CopyFile(templateSrc+name, "build/"+name); err != nil {
			return err
		}
	}
	return nil
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

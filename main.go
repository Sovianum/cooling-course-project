package main

import (
	"fmt"
	"github.com/Sovianum/cooling-course-project/core"
	cooling2 "github.com/Sovianum/cooling-course-project/core/cooling"
	"github.com/Sovianum/cooling-course-project/core/midline"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/core/schemes/three_shafts"
	"github.com/Sovianum/cooling-course-project/io"
	"github.com/Sovianum/cooling-course-project/postprocessing/builder"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/common"
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"gonum.org/v1/gonum/mat"
	"math"
	"os"
	"os/exec"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/gap"
	"github.com/Sovianum/turbocycle/common/gdf"
	"github.com/Sovianum/turbocycle/material/gases"
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

	startPiFactor   = 0.05
	piFactorStep    = 0.05
	piFactorStepNum = 19

	totalPiStag = 20
	lowPiStag   = 5.7
	highPiStag  = 3.5

	templatesDir = "postprocessing/templates"

	buildDir = "/home/artem/gowork/src/github.com/Sovianum/cooling-course-project/build"
	dataDir  = "build/data/"
	imgDir   = "build/img"

	projectInputTemplate = "project_input_data_template.tex"
	projectInputOut = "project_input_data.tex"

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

	hPointNum       = 50
	coolAirMassRate = 0.05
	theta0          = 500
	gapWidth        = 2.4e-3

	velocityCoef = 0.98
	massRateCoef = 0.98
)

func main() {
	io.PrepareDirectories(
		buildDir, dataDir, imgDir,
	)

	var scheme = getScheme(lowPiStag, highPiStag)

	saveInputTemplates()

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
			{"stator_root_1.csv", "stator_root_2.csv"},
			{"stator_mid_1.csv", "stator_mid_2.csv"},
			{"stator_top_1.csv", "stator_top_2.csv"},
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
			{"rotor_root_1.csv", "rotor_root_2.csv"},
			{"rotor_mid_1.csv", "rotor_mid_2.csv"},
			{"rotor_top_1.csv", "rotor_top_2.csv"},
		},
		true,
	)

	saveProfilingTemplate()

	var statorMidProfile = profiles.NewBladeProfileFromProfiler(
		0.5,
		0.01, 0.01,
		0.2, 0.2,
		statorProfiler,
	)
	var stagePack = stage.GetDataPack()
	statorMidProfile.Transform(geom.Scale(geometry.ChordProjection(stagePack.StageGeometry.StatorGeometry())))

	var gapCalculator = getGapCalculator(stage, statorMidProfile)
	var gapPack = gapCalculator.GetPack(coolAirMassRate)

	var gapCalcDF = getGapDF(common.LinSpace(0.01, 0.10, 10), gapCalculator)
	saveCooling1Template(gapCalcDF)

	//var psTemperatureSystem = getPSConvTemperatureSystem(gapPack.AlphaGas, stage, statorMidProfile)
	var psTemperatureSystem = getPSConvFilmTemperatureSystem(gapPack.AlphaGas, stage, statorMidProfile)
	var psSolution = psTemperatureSystem.Solve(0, theta0, 1, 0.001)

	//var ssTemperatureSystem = getSSConvTemperatureSystem(gapPack.AlphaGas, stage, statorMidProfile)
	var ssTemperatureSystem = getSSConvFilmTemperatureSystem(gapPack.AlphaGas, stage, statorMidProfile)
	var ssSolution = ssTemperatureSystem.Solve(0, theta0, 1, 0.001)

	var tempProfileDF = getTempProfileDF(gapCalcDF, stage, statorMidProfile, psSolution, ssSolution)
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

func saveCooling2Template(df dataframes.TProfileCalcDF) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cooling2Template,
		buildDir+"/"+cooling2Out,
	)

	var err error = nil
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func getTempProfileDF(
	gapDF dataframes.GapCalcDF,
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
	psSolution profile.TemperatureSolution,
	ssSolution profile.TemperatureSolution,
) dataframes.TProfileCalcDF {
	var inletTriangle = stage.VelocityInput().GetState().(states.VelocityPortState).Triangle

	var dInlet = 2 * geom.CurvRadius2(profile.InletEdge(), 0.5, 1e-3)

	var gas = stage.GasInput().GetState().(states2.GasPortState).Gas
	var tStagIn = stage.TemperatureInput().GetState().(states2.TemperaturePortState).TStag
	var pStagIn = stage.PressureInput().GetState().(states2.PressurePortState).PStag
	var density0 = pStagIn / (gas.R() * tStagIn)
	var ca = inletTriangle.CA()
	var massRateIntensity = density0 * ca

	var geomDF = dataframes.TProfileGeomDF{
		DInlet: dInlet,
	}

	var alphaMean = gapDF.Gas.AlphaGas

	var alphaGasSS = cooling.InletSSAlphaLaw(alphaMean)(0, tStagIn)
	var alphaGasOutlet = cooling.OutletSSAlphaLaw(alphaMean)(0, tStagIn)
	var alphaGasPS = cooling.PSAlphaLaw(alphaMean)(0, tStagIn)
	var alphaInlet = cooling.CylinderAlphaLaw(gas, massRateIntensity, dInlet)(0, tStagIn)

	var gasDF = dataframes.TProfileGasDF{
		Ca:             ca,
		RhoGas:         density0,
		MuGas:          gapDF.Gas.MuGas,
		LambdaGas:      gapDF.Gas.LambdaGas,
		AlphaMean:      alphaMean,
		AlphaGasSS:     alphaGasSS,
		AlphaGasPS:     alphaGasPS,
		AlphaGasOutlet: alphaGasOutlet,
		AlphaGasInlet:  alphaInlet,
		SkipSteps:      50,
	}
	gasDF.SetPSSolutionInfo(psSolution)
	gasDF.SetSSSolutionInfo(ssSolution)

	var calcDF = dataframes.TProfileCalcDF{
		Geom: geomDF,
		Gas:  gasDF,
	}
	return calcDF
}

func getSSConvTemperatureSystem(
	meanAlphaGas float64,
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
) profile.TemperatureSystem {
	var segment = profiles.SSSegment(profile, 0.5, 0.5)
	var alphaGasFunc, alphaAirFunc = getAlphaLaws(meanAlphaGas, stage, profile, cooling2.SSProfileGasAlphaLaw)
	if system, err := cooling2.GetInitedStatorConvTemperatureSystem(
		coolAirMassRate, stage, segment, alphaAirFunc, alphaGasFunc,
	); err != nil {
		panic(err)
	} else {
		return system
	}
}

func getPSConvTemperatureSystem(
	meanAlphaGas float64,
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
) profile.TemperatureSystem {
	var segment = profiles.PSSegment(profile, 0.5, 0.5)
	var alphaGasFunc, alphaAirFunc = getAlphaLaws(meanAlphaGas, stage, profile, cooling2.PSProfileGasAlphaLaw)
	if system, err := cooling2.GetInitedStatorConvTemperatureSystem(
		coolAirMassRate, stage, segment, alphaAirFunc, alphaGasFunc,
	); err != nil {
		panic(err)
	} else {
		return system
	}
}

func getSSConvFilmTemperatureSystem(
	meanAlphaGas float64,
	stage nodes.TurbineStageNode,
	bladeProfile profiles.BladeProfile,
) profile.TemperatureSystem {
	var segment = profiles.SSSegment(bladeProfile, 0.5, 0.5)
	var alphaGasFunc, alphaAirFunc = getAlphaLaws(meanAlphaGas, stage, bladeProfile, cooling2.SSProfileGasAlphaLaw)
	var lambdaLaw = getLambdaLaw(stage, cooling.SSLambdaLaw)

	var slitInfoArr = []profile.SlitInfo{
		{
			Coord:0,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
		{
			Coord:1e-2,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
		{
			Coord:2e-2,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
	}

	if system, err := cooling2.GetInitedStatorConvFilmTemperatureSystem(
		coolAirMassRate, stage, segment, alphaAirFunc, alphaGasFunc, lambdaLaw, slitInfoArr,
	); err != nil {
		panic(err)
	} else {
		return system
	}
}

func getPSConvFilmTemperatureSystem(
	meanAlphaGas float64,
	stage nodes.TurbineStageNode,
	bladeProfile profiles.BladeProfile,
) profile.TemperatureSystem {
	var segment = profiles.PSSegment(bladeProfile, 0.5, 0.5)
	var alphaGasFunc, alphaAirFunc = getAlphaLaws(meanAlphaGas, stage, bladeProfile, cooling2.PSProfileGasAlphaLaw)
	var lambdaLaw = getLambdaLaw(stage, cooling.PSLambdaLaw)

	var slitInfoArr = []profile.SlitInfo{
		{
			Coord:0,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
		{
			Coord:1e-2,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
		{
			Coord:2e-2,
			Thickness:3e-4,
			Area:5e-8,
			VelocityCoef:velocityCoef,
			MassRateCoef:massRateCoef,
		},
	}

	if system, err := cooling2.GetInitedStatorConvFilmTemperatureSystem(
		coolAirMassRate, stage, segment, alphaAirFunc, alphaGasFunc, lambdaLaw, slitInfoArr,
	); err != nil {
		panic(err)
	} else {
		return system
	}
}

func getLambdaLaw(stage nodes.TurbineStageNode, lambdaGenerator func(float64, float64) cooling.LambdaLaw) cooling.LambdaLaw {
	var gas = stage.GasOutput().GetState().(states2.GasPortState).Gas
	var tStagOut = stage.TemperatureOutput().GetState().(states2.TemperaturePortState).TStag
	var velocityOut = stage.VelocityOutput().GetState().(states.VelocityPortState).Triangle.C()

	var lambdaIn = 0.3
	var lambdaOut = velocityOut / gdf.ACrit(gases.K(gas, tStagOut), gas.R(), tStagOut)

	return lambdaGenerator(lambdaIn, lambdaOut)
}

func getAlphaLaws(
	meanAlphaGas float64,
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
	gasAlphaGenerator func(profiles.BladeProfile, float64, float64) cooling.AlphaLaw,
) (alphaGas cooling.AlphaLaw, alphaAir cooling.AlphaLaw) {
	var pack = stage.GetDataPack()
	var inletTriangle = stage.VelocityInput().GetState().(states.VelocityPortState).Triangle

	var gas = stage.GasInput().GetState().(states2.GasPortState).Gas
	var tStagIn = stage.TemperatureInput().GetState().(states2.TemperaturePortState).TStag
	var pStagIn = stage.PressureInput().GetState().(states2.PressurePortState).PStag
	var density0 = pStagIn / (gas.R() * tStagIn)

	var massRateIntensity = density0 * inletTriangle.CA()

	var dInlet = 2 * geom.CurvRadius2(profile.InletEdge(), 0.5, 1e-3)
	var alphaInlet = cooling.CylinderAlphaLaw(gas, massRateIntensity, dInlet)(0, tStagIn)

	alphaGas = gasAlphaGenerator(
		profile, alphaInlet, meanAlphaGas,
	)
	alphaAir = cooling.DefaultAirAlphaLaw(
		stage.GasInput().GetState().(states2.GasPortState).Gas,
		geometry.Height(0, pack.StageGeometry.StatorGeometry()),
		gapWidth, coolAirMassRate,
	)
	return
}

func saveCooling1Template(
	df dataframes.GapCalcDF,
) {
	var inserter = templ.NewDataInserter(
		templatesDir+"/"+cooling1Template,
		buildDir+"/"+cooling1Out,
	)

	var err error = nil
	if err != nil {
		panic(err)
	}
	if err := inserter.Insert(df); err != nil {
		panic(err)
	}
}

func getGapDF(
	massRateArr []float64,
	calculator gap.GapCalculator,
) dataframes.GapCalcDF {
	var dataPackArr = make([]gap.DataPack, len(massRateArr))

	for i, massRate := range massRateArr {
		var pack = calculator.GetPack(massRate)
		if pack.Err != nil {
			panic(pack.Err)
		}
		dataPackArr[i] = pack
	}
	var gapCalcDF = dataframes.GapCalcFromDataPacks(dataPackArr)
	gapCalcDF.Gas.NuCoef = 0.079 // todo remove hardcode

	return gapCalcDF
}

func getGapCalculator(
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
) gap.GapCalculator {
	if result, err := cooling2.GetInitedStatorGapCalculator(stage, profile); err != nil {
		panic(err)
	} else {
		return result
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
			tRelArr[i], 0,
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
	var df = dataframes.NewThreeShaftsDF(power, etaR, scheme)
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
		power / etaR,
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

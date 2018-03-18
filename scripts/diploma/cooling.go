package diploma

import (
	"fmt"
	cooling2 "github.com/Sovianum/cooling-course-project/core/cooling"
	"github.com/Sovianum/cooling-course-project/core/profiling"
	"github.com/Sovianum/cooling-course-project/postprocessing/dataframes"
	"github.com/Sovianum/cooling-course-project/postprocessing/templ"
	"github.com/Sovianum/turbocycle/common/gdf"
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/gap"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"math"
)

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

func saveCoolingSolution(solution profile.TemperatureSolution, fileName string) {
	if err := profiling.SaveMatrix(dataDir+"/"+fileName, solution.ToMatrix()); err != nil {
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
		Geom:       geomDF,
		Gas:        gasDF,
		PSSolution: psSolution,
		SSSolution: ssSolution,
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

	var slitGeomData = []struct {
		coord float64
		d     float64
	}{
		{3e-3, 0.4e-3},
		{27e-3, 0.5e-3},
		{37e-3, 0.5e-3},
		//{41e-3, 0.3e-3},
		//{30e-3, 0.5e-3},
		//{35e-3, 0.3e-3},
	}

	var slitInfoArr = make([]profile.SlitInfo, len(slitGeomData))
	for i, item := range slitGeomData {
		slitInfoArr[i] = profile.SlitInfo{
			Coord:        item.coord,
			Thickness:    getSlitThk(item.d),
			Area:         getSlitArea(item.d),
			VelocityCoef: velocityCoef,
			MassRateCoef: massRateCoef,
		}
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

	var slitGeomData = []struct {
		coord float64
		d     float64
	}{
		{3e-3, 0.4e-3},
		{16e-3, 0.25e-3},
		{22e-3, 0.25e-3},
		{27e-3, 0.3e-3},
		{31e-3, 0.35e-3},
		{36.5e-3, 0.40e-3},
	}

	var slitInfoArr = make([]profile.SlitInfo, len(slitGeomData))
	for i, item := range slitGeomData {
		slitInfoArr[i] = profile.SlitInfo{
			Coord:        item.coord,
			Thickness:    getSlitThk(item.d),
			Area:         getSlitArea(item.d),
			VelocityCoef: velocityCoef,
			MassRateCoef: massRateCoef,
		}
	}

	if system, err := cooling2.GetInitedStatorConvFilmTemperatureSystem(
		coolAirMassRate, stage, segment, alphaAirFunc, alphaGasFunc, lambdaLaw, slitInfoArr,
	); err != nil {
		panic(err)
	} else {
		return system
	}
}

func getSlitThk(diameter float64) float64 {
	return math.Pi * diameter * diameter / 4 * coolingHoleNum / coolingBladeLength
}

func getSlitArea(diameter float64) float64 {
	return math.Pi * diameter * diameter / 4 * coolingHoleNum
}

func getLambdaLaw(stage nodes.TurbineStageNode, lambdaGenerator func(float64, float64) cooling.LambdaLaw) cooling.LambdaLaw {
	var gas = stage.GasInput().GetState().(states2.GasPortState).Gas
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

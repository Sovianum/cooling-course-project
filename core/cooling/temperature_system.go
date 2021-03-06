package cooling

import (
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode/forward"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
)

func GetInitedStatorConvTemperatureSystem(
	airMassRate float64,
	stage turbine.StageNode,
	segment geom.Segment,
	alphaAirFunc cooling.AlphaLaw,
	alphaGasFunc cooling.AlphaLaw,
) (profile.TemperatureSystem, error) {
	var dataPack = stage.GetDataPack()
	if dataPack.Err != nil {
		return nil, dataPack.Err
	}
	var tGas = stage.TemperatureInput().GetState().(states.TemperaturePortState).TStag

	return profile.NewConvectiveTemperatureSystem(
		forward.NewEulerSolver(),
		airMassRate,
		gases.GetAir().Cp,
		func(x float64) float64 {
			return tGas
		},
		alphaAirFunc,
		alphaGasFunc,
		func(x float64) float64 {
			return wallThk
		},
		func(t float64) float64 {
			return lambdaM
		},
		segment,
	), nil
}

func GetInitedStatorConvFilmTemperatureSystem(
	coolerMassRate0 float64,
	stage turbine.StageNode,
	segment geom.Segment,
	alphaAirFunc cooling.AlphaLaw,
	alphaGasFunc cooling.AlphaLaw,
	law cooling.LambdaLaw,
	slitInfoArray []profile.SlitInfo,
) (profile.TemperatureSystem, error) {
	var dataPack = stage.GetDataPack()
	if dataPack.Err != nil {
		return nil, dataPack.Err
	}
	var gas = stage.GasInput().GetState().(states.GasPortState).Gas
	var tGas = stage.TemperatureInput().GetState().(states.TemperaturePortState).TStag
	var pGas = stage.PressureInput().GetState().(states.PressurePortState).PStag

	return profile.NewConvFilmTemperatureSystem(
		forward.NewEulerSolver(),
		coolerMassRate0,
		gases.GetAir(), gas,
		func(x float64) float64 {
			return tGas
		},
		func(x float64) float64 {
			return pGas * 0.95
		},

		func(x float64) float64 {
			return pGas
		},
		law,
		alphaAirFunc, alphaGasFunc,
		slitInfoArray,
		func(x float64) float64 {
			return wallThk
		},
		func(t float64) float64 {
			return lambdaM
		},
		segment,
	), nil
}

func PSProfileGasAlphaLaw(
	profile profiles.BladeProfile,
	inletAlpha float64,
	meanAlpha float64,
) cooling.AlphaLaw {
	var inletEdgeLength = geom.ApproxLength(profile.InletEdge(), 0.5, 1, 100)
	var psLength = geom.ApproxLength(profile.PSLine(), 0, 1, 100)

	var boundaryValue = inletEdgeLength
	var totalLength = inletEdgeLength + psLength

	return cooling.JoinedAlphaLaw(
		[]cooling.AlphaLaw{
			cooling.ConstantAlphaLaw(inletAlpha),
			cooling.PSAlphaLaw(meanAlpha),
		}, []float64{0, boundaryValue, totalLength},
	)
}

func SSProfileGasAlphaLaw(
	profile profiles.BladeProfile,
	inletAlpha float64,
	meanAlpha float64,
) cooling.AlphaLaw {
	var inletEdgeLength = geom.ApproxLength(profile.InletEdge(), 0.5, 1, 100)
	var ssLength = geom.ApproxLength(profile.SSLine(), 0, 1, 100)
	var totalLength = inletEdgeLength + ssLength

	var boundary1 = inletEdgeLength
	var boundary2 = inletEdgeLength + 2*ssLength/3

	return cooling.JoinedAlphaLaw(
		[]cooling.AlphaLaw{
			cooling.ConstantAlphaLaw(inletAlpha),
			cooling.InletSSAlphaLaw(meanAlpha),
			cooling.OutletSSAlphaLaw(meanAlpha),
		}, []float64{0, boundary1, boundary2, totalLength},
	)
}

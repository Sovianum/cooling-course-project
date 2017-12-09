package cooling

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
)

func GetInitedStatorTemperatureSystem(
	airMassRate float64,
	stage nodes.TurbineStageNode,
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
		ode.NewEulerSolver(),
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
	var boundary2 = inletEdgeLength + 2 * ssLength / 3

	return cooling.JoinedAlphaLaw(
		[]cooling.AlphaLaw{
			cooling.ConstantAlphaLaw(inletAlpha),
			cooling.InletSSAlphaLaw(meanAlpha),
			cooling.OutletSSAlphaLaw(meanAlpha),
		}, []float64{0, boundary1, boundary2, totalLength},
	)
}

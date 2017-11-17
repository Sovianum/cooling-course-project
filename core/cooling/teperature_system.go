package cooling

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/ode"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
)

const (
	airMassRate = 0.05
)

func GetInitedStatorTemperatureSystem(
	stage nodes.TurbineStageNode,
	segment geom.Segment,
	alphaAirFunc cooling.AlphaLaw,
	alphaGasFunc cooling.AlphaLaw,
) (cooling.TemperatureSystem, error) {
	var dataPack = stage.GetDataPack()
	if dataPack.Err != nil {
		return nil, dataPack.Err
	}
	var tGas = stage.TemperatureInput().GetState().(states.TemperaturePortState).TStag

	return cooling.NewTemperatureSystem(
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

func PSProfileAlphaLaw(
	gas gases.Gas,
	profile profiles.BladeProfile,
	massRateIntensity float64,
	meanAlpha float64,
) cooling.AlphaLaw {
	var inletEdgeDiameter = 2 * geom.CurvRadius2(profile.InletEdge(), 0.5, 1e-4)
	var inletEdgeLength = geom.ApproxLength(profile.InletEdge(), 0.5, 1, 100)
	var psLength = geom.ApproxLength(profile.PSLine(), 0, 1, 100)

	var boundaryValue = inletEdgeLength / (inletEdgeLength + psLength)

	return cooling.JoinedAlphaLaw(
		[]cooling.AlphaLaw{
			cooling.CylinderAlphaLaw(gas, massRateIntensity, inletEdgeDiameter),
			cooling.PSAlphaLaw(meanAlpha),
		}, []float64{boundaryValue},
	)
}

func SSProfileAlphaLaw(
	gas gases.Gas,
	profile profiles.BladeProfile,
	massRateIntensity float64,
	meanAlpha float64,
) cooling.AlphaLaw {
	var inletEdgeDiameter = 2 * geom.CurvRadius2(profile.InletEdge(), 0.5, 1e-4)
	var inletEdgeLength = geom.ApproxLength(profile.InletEdge(), 0.5, 1, 100)
	var psLength = geom.ApproxLength(profile.PSLine(), 0, 1, 100)

	var boundaryValue = inletEdgeLength / (inletEdgeLength + psLength)

	return cooling.JoinedAlphaLaw(
		[]cooling.AlphaLaw{
			cooling.CylinderAlphaLaw(gas, massRateIntensity, inletEdgeDiameter),
			cooling.PSAlphaLaw(meanAlpha),
		}, []float64{boundaryValue},
	)
}

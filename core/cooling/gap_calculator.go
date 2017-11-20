package cooling

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	states2 "github.com/Sovianum/turbocycle/impl/turbine/states"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"math"
)

const (
	wallThk = 1e-3
	lambdaM = 20
	tWallOuter = 1000
	tCoolerInlet = 500
)

func GetInitedStatorGapCalculator(
	stage nodes.TurbineStageNode,
	profile profiles.BladeProfile,
) (cooling.GapCalculator, error) {
	var dataPack = stage.GetDataPack()
	if dataPack.Err != nil {
		return nil, dataPack.Err
	}

	var gas = stage.GasInput().GetState().(states.GasPortState).Gas
	var ca = stage.VelocityInput().GetState().(states2.VelocityPortState).Triangle.CA()	// todo fix axial velocity
	var pGas = stage.PressureInput().GetState().(states.PressurePortState).PStag
	var tGas = stage.TemperatureInput().GetState().(states.TemperaturePortState).TStag
	return cooling.NewGapCalculator(
		gases.GetAir(), gas,
		ca, pGas,
		dataPack.StageGeometry.StatorGeometry(),
		profile,
		wallThk,
		lambdaM,
		func(re float64) float64 {
			return 0.079 * math.Pow(re, 0.68)
		},
		tGas,
		tWallOuter,
		tCoolerInlet,
	), nil
}

package cooling

import (
	"github.com/Sovianum/turbocycle/impl/engine/states"
	states2 "github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/gap"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profiles"
	"math"
)

const (
	wallThk      = 1e-3
	lambdaM      = 20
	tWallOuter   = 1000
	tCoolerInlet = 500
)

func GetInitedStatorGapCalculator(
	stage turbine.StageNode,
	profile profiles.BladeProfile,
) (gap.GapCalculator, error) {
	var dataPack = stage.GetDataPack()
	if dataPack.Err != nil {
		return nil, dataPack.Err
	}

	var gas = stage.GasInput().GetState().(states.GasPortState).Gas
	var ca = stage.VelocityInput().GetState().(states2.VelocityPortState).Triangle.CA()
	var pGas = stage.PressureInput().GetState().(states.PressurePortState).PStag
	var tGas = stage.TemperatureInput().GetState().(states.TemperaturePortState).TStag
	return gap.NewGapCalculator(
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

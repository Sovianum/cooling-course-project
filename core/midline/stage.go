package midline

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	n          = 1.15e4
	reactivity = 0.30
	phi        = 0.98
	psi        = 0.98
	airGapRel  = 0.001
	precision  = 0.05

	lRelOut = 0.08

	gammaIn  = -2
	gammaOut = 2

	baRelStator = 1.3
	baRelRotor  = 1.75

	deltaRelStator = 0.1
	deltaRelRotor  = 0.1

	approxTRelStator = 0.7
	approxTRelRotor  = 0.7

	alpha1 = 14
	power  = 16e6 / 0.98
)

func GetInitedStageNode(scheme schemes.ThreeShaftsScheme) turbine.StageNode {
	var geomGen = turbine.NewStageGeometryGenerator(
		lRelOut,
		turbine.NewIncompleteGenerator(
			baRelStator,
			deltaRelStator,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
			approxTRelStator,
		),
		turbine.NewIncompleteGenerator(
			baRelRotor,
			deltaRelRotor,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
			approxTRelRotor,
		),
	)
	var stage = turbine.NewTurbineSingleStageNode(
		n,
		constructive.Ht(scheme.GasGenerator().TurboCascade().Turbine()),
		reactivity, phi, psi, airGapRel, precision, geomGen,
	)
	var massRate = schemes.GetMassRate(power, scheme)
	turbine.InitFromTurbineNode(
		stage, scheme.GasGenerator().TurboCascade().Turbine(),
		massRate,
		common.ToRadians(alpha1),
	)
	return stage
}

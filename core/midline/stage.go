package midline

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	n          = 1.3e4
	reactivity = 0.4
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
	power  = 16e6/0.98
)

func GetInitedStageNode(scheme schemes.ThreeShaftsScheme) nodes.TurbineStageNode {
	var geomGen = geometry.NewStageGeometryGenerator(
		lRelOut,
		geometry.NewIncompleteGeneratorFromProfileAngles(
			baRelStator,
			deltaRelStator,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
			approxTRelStator,
		),
		geometry.NewIncompleteGeneratorFromProfileAngles(
			baRelRotor,
			deltaRelRotor,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
			approxTRelRotor,
		),
	)
	var stage = nodes.NewTurbineStageNode(
		n,
		constructive.Ht(scheme.GasGenerator().TurboCascade().Turbine()),
		reactivity, phi, psi, airGapRel, precision, geomGen,
	)
	var massRate = schemes.GetMassRate(power, scheme)
	nodes.InitFromTurbineNode(
		stage, scheme.GasGenerator().TurboCascade().Turbine(),
		massRate,
		common.ToRadians(alpha1),
	)
	return stage
}

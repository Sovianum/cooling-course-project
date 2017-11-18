package midline

import (
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
)

const (
	n             = 1e4
	reactivity    = 0.5
	phi           = 0.98
	psi           = 0.98
	airGapRel     = 0.001
	precision     = 0.05

	lRelOut = 0.15

	gammaIn = -10
	gammaOut = 10

	baRelStator = 4
	baRelRotor = 4

	deltaRelStator = 0.1
	deltaRelRotor = 0.1

	alpha1 = 14
	power = 16e6
)

func GetInitedStageNode(scheme schemes.ThreeShaftsScheme) nodes.TurbineStageNode {
	var geomGen = geometry.NewStageGeometryGenerator(
		lRelOut,
		geometry.NewIncompleteGeneratorFromProfileAngles(
			baRelStator,
			deltaRelStator,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
		),
		geometry.NewIncompleteGeneratorFromProfileAngles(
			baRelRotor,
			deltaRelRotor,
			common.ToRadians(gammaIn),
			common.ToRadians(gammaOut),
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

package dataframes

import (
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/turbine/geometry"
	"github.com/Sovianum/turbocycle/impl/turbine/nodes"
	"github.com/Sovianum/turbocycle/impl/turbine/states"
)

func NewBladingGeometryDF(relGen geometry.BladingGeometryGenerator, geom geometry.BladingGeometry) BladingGeometryDF {
	return BladingGeometryDF{
		LRelOut:    relGen.LRelOut(),
		Elongation: relGen.Elongation(),
		DeltaRel:   relGen.DeltaRel(),
		GammaIn:    relGen.GammaIn(),
		GammaOut:   relGen.GammaOut(),
		AreaOut:    geometry.Area(geom.XBladeOut(), geom),
		DMeanOut:   geom.MeanProfile().Diameter(geom.XBladeOut()),
		DMeanIn:    geom.MeanProfile().Diameter(geom.XBladeOut()),
		LOut:       geometry.Height(geom.XBladeOut(), geom),
	}
}

type BladingGeometryDF struct {
	LRelOut    float64
	Elongation float64
	DeltaRel   float64
	GammaIn    float64
	GammaOut   float64
	AreaOut    float64
	DMeanOut   float64
	DMeanIn    float64
	LOut       float64
}

func NewStageDF(node nodes.TurbineStageNode) (StageDF, error) {
	var pack = node.GetDataPack()
	if pack.Err != nil {
		return StageDF{}, nil
	}

	var xBladeStatorOut = pack.StageGeometry.StatorGeometry().XBladeOut()
	var xStatorOut = pack.StageGeometry.StatorGeometry().XGapOut()
	var xBladeRotorOut = pack.StageGeometry.RotorGeometry().XBladeOut()
	var x = xStatorOut - xBladeStatorOut + xBladeRotorOut

	return StageDF{
		Reactivity: pack.Reactivity,

		Ht:     node.Ht(),
		Hs:     pack.StatorHeatDrop,
		Hr:     pack.RotorHeatDrop,
		HtStag: pack.StageHeatDropStag,

		C1Ad: pack.C1Ad,
		Phi:  pack.Phi,

		Tg: node.TemperatureInput().GetState().(states2.TemperaturePortState).TStag,

		P1:      pack.P1,
		T1:      pack.T1,
		T1Prime: pack.T1Prime,
		Rho1:    pack.Density1,

		MassRate: node.MassRateInput().GetState().(states.MassRatePortState).MassRate,
		RPM:      pack.RPM,

		Pw1: pack.Pw1,
		Tw1: pack.Tw1,

		W2Ad: pack.WAd2,
		Psi:  pack.Psi,

		P2:      pack.P2,
		T2:      pack.T2,
		T2Prime: pack.T2Prime,
		T2Stag:  pack.T2Stag,
		Rho2:    pack.Density2,

		PStagIn:  node.PressureInput().GetState().(states2.PressurePortState).PStag,
		PStagOut: node.PressureOutput().GetState().(states2.PressurePortState).PStag,
		Pi:       pack.Pi,

		Lu:   pack.MeanRadiusLabour,
		EtaU: pack.EtaU,

		LossStator:  pack.StatorSpecificLoss,
		LossRotor:   pack.RotorSpecificLoss,
		LossOutflow: pack.OutletVelocitySpecificLoss,
		LossRadial:  pack.AirGapSpecificLoss,
		LossVent:    pack.VentilationSpecificLoss,

		EtaPower: pack.EtaT,
		Lt:       pack.StageLabour,
		EtaTStag: pack.EtaTStag,

		StatorGeom: NewBladingGeometryDF(
			node.StageGeomGen().StatorGenerator(),
			pack.StageGeometry.StatorGeometry(),
		),
		RotorGeom: NewBladingGeometryDF(
			node.StageGeomGen().RotorGenerator(),
			pack.StageGeometry.RotorGeometry(),
		),
		X:      x,
		DeltaR: pack.AirGapRel,

		StatorGas: NewGasMeanDF(
			pack.P0,
			node.TemperatureInput().GetState().(states2.TemperaturePortState).TStag,
			pack.T1,
			node.GasInput().GetState().(states2.GasPortState).Gas,
		),
		RotorGas: NewGasMeanDF(
			pack.P1,
			pack.T1,
			pack.T2,
			node.GasInput().GetState().(states2.GasPortState).Gas,
		),

		InletTriangle:  pack.RotorInletTriangle,
		OutletTriangle: pack.RotorOutletTriangle,
	}, nil
}

type StageDF struct {
	Reactivity float64

	Ht     float64
	Hs     float64
	Hr     float64
	HtStag float64

	C1Ad float64
	Phi  float64

	Tg float64

	P1      float64
	T1      float64
	T1Prime float64
	Rho1    float64

	MassRate float64
	RPM      float64

	Pw1 float64
	Tw1 float64

	W2Ad float64
	Psi  float64

	P2      float64
	T2      float64
	T2Prime float64
	T2Stag  float64
	Rho2    float64

	PStagIn  float64
	PStagOut float64

	Pi float64

	Lu   float64
	EtaU float64

	LossStator  float64
	LossRotor   float64
	LossOutflow float64
	LossRadial  float64
	LossVent    float64

	EtaPower float64
	Lt       float64
	EtaTStag float64

	StatorGeom BladingGeometryDF
	RotorGeom  BladingGeometryDF
	X          float64
	DeltaR     float64

	StatorGas GasMeanDF
	RotorGas  GasMeanDF

	InletTriangle  states.VelocityTriangle
	OutletTriangle states.VelocityTriangle
}

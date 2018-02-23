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
	LRelOut    float64 `json:"l_rel_out"`
	Elongation float64 `json:"elongation"`
	DeltaRel   float64 `json:"delta_rel"`
	GammaIn    float64 `json:"gamma_in"`
	GammaOut   float64 `json:"gamma_out"`
	AreaOut    float64 `json:"area_out"`
	DMeanOut   float64 `json:"d_mean_out"`
	DMeanIn    float64 `json:"d_mean_in"`
	LOut       float64 `json:"l_out"`
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

		MassRate: node.MassRateInput().GetState().Value().(float64),
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
		Gas:NewGasMeanDF(
			pack.P0,
			node.TemperatureInput().GetState().(states2.TemperaturePortState).TStag,
			node.TemperatureOutput().GetState().(states2.TemperaturePortState).TStag,
			node.GasInput().GetState().(states2.GasPortState).Gas,
		),

		InletTriangle:  pack.RotorInletTriangle,
		OutletTriangle: pack.RotorOutletTriangle,
	}, nil
}

type StageDF struct {
	Reactivity float64 `json:"reactivity"`

	Ht     float64 `json:"ht"`
	Hs     float64 `json:"hs"`
	Hr     float64 `json:"hr"`
	HtStag float64 `json:"ht_stag"`

	C1Ad float64 `json:"c_1_ad"`
	Phi  float64 `json:"phi"`

	Tg float64 `json:"tg"`

	P1      float64 `json:"p_1"`
	T1      float64 `json:"t_1"`
	T1Prime float64 `json:"t_1_prime"`
	Rho1    float64 `json:"rho_1"`

	MassRate float64 `json:"mass_rate"`
	RPM      float64 `json:"rpm"`

	Pw1 float64 `json:"pw_1"`
	Tw1 float64 `json:"tw_1"`

	W2Ad float64 `json:"w_2_ad"`
	Psi  float64 `json:"psi"`

	P2      float64 `json:"p_2"`
	T2      float64 `json:"t_2"`
	T2Prime float64 `json:"t_2_prime"`
	T2Stag  float64 `json:"t_2_stag"`
	Rho2    float64 `json:"rho_2"`

	PStagIn  float64 `json:"p_stag_in"`
	PStagOut float64 `json:"p_stag_out"`

	Pi float64 `json:"pi"`

	Lu   float64 `json:"lu"`
	EtaU float64 `json:"eta_u"`

	LossStator  float64 `json:"loss_stator"`
	LossRotor   float64 `json:"loss_rotor"`
	LossOutflow float64 `json:"loss_outflow"`
	LossRadial  float64 `json:"loss_radial"`
	LossVent    float64 `json:"loss_vent"`

	EtaPower float64 `json:"eta_power"`
	Lt       float64 `json:"lt"`
	EtaTStag float64 `json:"eta_t_stag"`

	StatorGeom BladingGeometryDF `json:"stator_geom"`
	RotorGeom  BladingGeometryDF `json:"rotor_geom"`
	X          float64           `json:"x"`
	DeltaR     float64           `json:"delta_r"`

	StatorGas GasMeanDF `json:"stator_gas"`
	RotorGas  GasMeanDF `json:"rotor_gas"`
	Gas       GasMeanDF `json:"gas"`

	InletTriangle  states.VelocityTriangle `json:"inlet_triangle"`
	OutletTriangle states.VelocityTriangle `json:"outlet_triangle"`
}

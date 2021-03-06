package dataframes

import (
	"github.com/Sovianum/turbocycle/common"
	states2 "github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/impl/stage/turbine"
)

func NewStagedTurbineDF(node turbine.StagedTurbineNode) StagedTurbineDF {
	result := StagedTurbineDF{StageDFS: make([]TurbineStageDF, len(node.Stages()))}
	for i, stage := range node.Stages() {
		df, err := NewTurbineStageDF(stage)
		if err != nil {
			panic(err)
		}
		result.StageDFS[i] = df
	}
	return result
}

type StagedTurbineDF struct {
	StageDFS []TurbineStageDF
	rows     []StageRow
}

func (df StagedTurbineDF) Join(another StagedTurbineDF) StagedTurbineDF {
	return StagedTurbineDF{StageDFS: append(df.StageDFS, another.StageDFS...)}
}

func (df StagedTurbineDF) Rows() []StageRow {
	if df.rows == nil {
		getRow := func(name, dim string, f func(df TurbineStageDF) float64) StageRow {
			return NewTurbineStageRow(name, dim, df.StageDFS, f)
		}
		df.rows = []StageRow{
			getRow("$H_с$", "$10^6 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.Hs }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$c_{1ад}$", "$м/с$", func(df TurbineStageDF) float64 { return df.C1Ad }).FormatString(Round1),
			getRow("$c_{1}$", "$м/с$", func(df TurbineStageDF) float64 { return df.InletTriangle.C() }).FormatString(Round1),
			getRow("$T_1$", "$К$", func(df TurbineStageDF) float64 { return df.T1 }).FormatString(Round1),
			getRow("$T_1^\\prime$", "$К$", func(df TurbineStageDF) float64 { return df.T1Prime }).FormatString(Round1),
			getRow("$p_1$", "$МПа$", func(df TurbineStageDF) float64 { return df.P1 }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$\\rho_1$", "$кг/м^3$", func(df TurbineStageDF) float64 { return df.Rho1 }).FormatString(Round2),
			getRow("$\\alpha_1$", "$\\degree$", func(df TurbineStageDF) float64 { return df.InletTriangle.Alpha() }).FormatFloat(common.ToDegrees).FormatString(Round1),
			getRow("$c_{1a}$", "$м/с$", func(df TurbineStageDF) float64 { return df.InletTriangle.CA() }).FormatString(Round1),
			getRow("$A_1$", "$м^2$", func(df TurbineStageDF) float64 { return df.StatorGeom.AreaOut }).FormatString(Round2),
			getRow("$D_1$", "$м$", func(df TurbineStageDF) float64 { return df.StatorGeom.DMeanOut }).FormatString(Round3),
			getRow("$u_1$", "$м/с$", func(df TurbineStageDF) float64 { return df.InletTriangle.U() }).FormatString(Round1),
			getRow("$w_1$", "$м/с$", func(df TurbineStageDF) float64 { return df.InletTriangle.W() }).FormatString(Round1),
			getRow("$T_{w1}$", "$К$", func(df TurbineStageDF) float64 { return df.Tw1 }).FormatString(Round1),
			getRow("$p_{w1}$", "$МПа$", func(df TurbineStageDF) float64 { return df.Pw1 }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$H_л$", "$10^6 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.Hr }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$x$", "$м$", func(df TurbineStageDF) float64 { return df.X }).FormatString(Round3),
			getRow("$D_2$", "$м$", func(df TurbineStageDF) float64 { return df.RotorGeom.DMeanOut }).FormatString(Round3),
			getRow("$l_2$", "$м$", func(df TurbineStageDF) float64 { return df.RotorGeom.LOut }).FormatString(Round3),
			getRow("$\\left( \\frac{l}{D} \\right)_2$", "$-$", func(df TurbineStageDF) float64 { return df.RotorGeom.LRelOut }).FormatString(Round3),
			getRow("$u_2$", "$м/с$", func(df TurbineStageDF) float64 { return df.OutletTriangle.U() }).FormatString(Round1),
			getRow("$w_{2ад}$", "$м/с$", func(df TurbineStageDF) float64 { return df.W2Ad }).FormatString(Round1),
			getRow("$w_2$", "$м/с$", func(df TurbineStageDF) float64 { return df.OutletTriangle.W() }).FormatString(Round1),
			getRow("$T_2$", "$К$", func(df TurbineStageDF) float64 { return df.T2 }).FormatString(Round1),
			getRow("$T_2^\\prime$", "$К$", func(df TurbineStageDF) float64 { return df.T2Prime }).FormatString(Round1),
			getRow("$p_2$", "$МПа$", func(df TurbineStageDF) float64 { return df.P2 }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$\\beta_2$", "$\\degree$", func(df TurbineStageDF) float64 { return df.OutletTriangle.Beta() }).FormatFloat(common.ToDegrees).FormatString(Round1),
			getRow("$\\alpha_2$", "$\\degree$", func(df TurbineStageDF) float64 { return df.OutletTriangle.Alpha() }).FormatFloat(common.ToDegrees).FormatString(Round1),
			getRow("$c_2$", "$м/с$", func(df TurbineStageDF) float64 { return df.OutletTriangle.CU() }).FormatString(Round1),
			getRow("$\\pi_т$", "$-$", func(df TurbineStageDF) float64 { return df.Pi }).FormatString(Round2),
			getRow("$c_{2a}$", "$м/с$", func(df TurbineStageDF) float64 { return df.OutletTriangle.U() }).FormatString(Round2),
			getRow("$\\rho_2$", "$кг/м^3$", func(df TurbineStageDF) float64 { return df.Rho2 }).FormatString(Round2),
			getRow("$L_u$", "$10^6 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.Lu }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$\\eta_u$", "$-$", func(df TurbineStageDF) float64 { return df.EtaU }).FormatString(Round2),
			getRow("$h_с$", "$10^3 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.LossStator }).FormatFloat(DivideE3).FormatString(Round2),
			getRow("$h_р$", "$10^3 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.LossRotor }).FormatFloat(DivideE3).FormatString(Round2),
			getRow("$h_{вых}$", "$10^3 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.LossOutflow }).FormatFloat(DivideE3).FormatString(Round2),
			getRow("$h_з$", "$10^3 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.LossRadial }).FormatFloat(DivideE3).FormatString(Round2),
			getRow("$h_{вент}$", "$10^3 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.LossVent }).FormatFloat(DivideE3).FormatString(Round2),
			getRow("$T_2^*$", "$К$", func(df TurbineStageDF) float64 { return df.T2Stag }).FormatString(Round1),
			getRow("$p_2^*$", "$МПа$", func(df TurbineStageDF) float64 { return df.PStagOut }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$\\eta_{т \\/\\ мощн}$", "$-$", func(df TurbineStageDF) float64 { return df.EtaPower }).FormatString(Round2),
			getRow("$L_т$", "$10^6 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.Lt }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$H_т^*$", "$10^6 \\cdot Дж/кг$", func(df TurbineStageDF) float64 { return df.HtStag }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$\\eta_т^*$", "$-$", func(df TurbineStageDF) float64 { return df.EtaTStag }).FormatString(Round2),
		}
	}
	for i := range df.rows {
		df.rows[i].ID = i + 1
	}
	return df.rows
}

func NewTurbineStageDF(node turbine.StageNode) (TurbineStageDF, error) {
	var pack = node.GetDataPack()
	if pack.Err != nil {
		return TurbineStageDF{}, nil
	}

	var xBladeStatorOut = pack.StageGeometry.StatorGeometry().XBladeOut()
	var xStatorOut = pack.StageGeometry.StatorGeometry().XGapOut()
	var xBladeRotorOut = pack.StageGeometry.RotorGeometry().XBladeOut()
	var x = xStatorOut - xBladeStatorOut + xBladeRotorOut

	return TurbineStageDF{
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

		StatorGeom: NewTurbineBladingGeometryDF(
			node.StageGeomGen().StatorGenerator(),
			pack.StageGeometry.StatorGeometry(),
		),
		RotorGeom: NewTurbineBladingGeometryDF(
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
		Gas: NewGasMeanDF(
			pack.P0,
			node.TemperatureInput().GetState().(states2.TemperaturePortState).TStag,
			node.TemperatureOutput().GetState().(states2.TemperaturePortState).TStag,
			node.GasInput().GetState().(states2.GasPortState).Gas,
		),

		InletTriangle:  pack.RotorInletTriangle,
		OutletTriangle: pack.RotorOutletTriangle,
	}, nil
}

type TurbineStageDF struct {
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

	StatorGeom TurbineBladingGeometryDF `json:"stator_geom"`
	RotorGeom  TurbineBladingGeometryDF `json:"rotor_geom"`
	X          float64                  `json:"x"`
	DeltaR     float64                  `json:"delta_r"`

	StatorGas GasMeanDF `json:"stator_gas"`
	RotorGas  GasMeanDF `json:"rotor_gas"`
	Gas       GasMeanDF `json:"gas"`

	InletTriangle  states.VelocityTriangle `json:"inlet_triangle"`
	OutletTriangle states.VelocityTriangle `json:"outlet_triangle"`
}

func NewTurbineBladingGeometryDF(relGen turbine.BladingGeometryGenerator, geom geometry.BladingGeometry) TurbineBladingGeometryDF {
	return TurbineBladingGeometryDF{
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

type TurbineBladingGeometryDF struct {
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

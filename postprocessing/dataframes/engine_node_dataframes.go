package dataframes

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/material/fuel"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/core/graph"
)

func NewGasDF(p, t float64, gas gases.Gas) GasDF {
	return GasDF{
		T:       t,
		P:       p,
		Density: gases.Density(gas, t, p),
		K:       gases.K(gas, t),
		Cp:      gas.Cp(t),
		R:       gas.R(),
	}
}

type GasDF struct {
	T       float64 `json:"t"`
	P       float64 `json:"p"`
	Density float64 `json:"density"`
	K       float64 `json:"k"`
	Cp      float64 `json:"cp"`
	R       float64 `json:"r"`
}

func NewGasMeanDF(p, t1, t2 float64, gas gases.Gas) GasMeanDF {
	return GasMeanDF{
		T1:     t1,
		T2:     t2,
		P:      p,
		KMean:  gases.KMean(gas, t1, t2, nodes.DefaultN),
		CpMean: gases.CpMean(gas, t1, t2, nodes.DefaultN),
		R:      gas.R(),
	}
}

type GasMeanDF struct {
	T1     float64 `json:"t_1"`
	T2     float64 `json:"t_2"`
	P      float64 `json:"p"`
	KMean  float64 `json:"k_mean"`
	CpMean float64 `json:"cp_mean"`
	R      float64 `json:"r"`
}

func NewCompressorDF(node constructive.CompressorNode) CompressorDF {
	return CompressorDF{
		PIn:  node.PStagIn(),
		POut: node.PStagOut(),

		TIn:  node.TStagIn(),
		TOut: node.TStagOut(),

		Pi:     node.PiStag(),
		Labour: node.LSpecific(),
		Eta:    node.Eta(),
		EtaPol: node.EtaPol(),

		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			node.GasInput().GetState().Value().(gases.Gas),
		),
	}
}

type CompressorDF struct {
	PIn  float64 `json:"p_in"`
	POut float64 `json:"p_out"`

	TIn  float64 `json:"t_in"`
	TOut float64 `json:"t_out"`

	Pi     float64 `json:"pi"`
	Eta    float64 `json:"eta"`
	EtaPol float64 `json:"eta_pol"`

	Labour float64 `json:"labour"`

	GasData GasMeanDF `json:"gas_data"`
}

func NewPressureDropDF(node constructive.PressureLossNode) PressureDropDF {
	return PressureDropDF{
		PIn:   node.PStagIn(),
		POut:  node.PStagOut(),
		TIn:   node.TStagIn(),
		TOut:  node.TStagOut(),
		Sigma: node.Sigma(),
	}
}

type PressureDropDF struct {
	PIn  float64 `json:"p_in"`
	POut float64 `json:"p_out"`

	TIn  float64 `json:"t_in"`
	TOut float64 `json:"t_out"`

	Sigma float64 `json:"sigma"`
}

func NewRegeneratorNode(node constructive.RegeneratorNode) RegeneratorDF {
	extractor := func(port graph.Port) float64 {return port.GetState().Value().(float64)}

	return RegeneratorDF{
		PColdIn:  extractor(node.ColdInput().PressureInput()),
		PColdOut: extractor(node.ColdOutput().PressureOutput()),
		PHotIn:   extractor(node.HotInput().PressureInput()),
		PHotOut:  extractor(node.HotOutput().PressureOutput()),

		TColdIn:  extractor(node.ColdInput().TemperatureInput()),
		TColdOut: extractor(node.ColdOutput().TemperatureOutput()),
		THotIn:   extractor(node.HotInput().TemperatureInput()),
		THotOut:  extractor(node.HotOutput().TemperatureOutput()),

		Sigma: node.Sigma(),
	}
}

type RegeneratorDF struct {
	PColdIn  float64 `json:"p_cold_in"`
	PColdOut float64 `json:"p_cold_out"`
	PHotIn   float64 `json:"p_hot_in"`
	PHotOut  float64 `json:"p_hot_out"`

	TColdIn  float64 `json:"t_cold_in"`
	TColdOut float64 `json:"t_cold_out"`
	THotIn   float64 `json:"t_hot_in"`
	THotOut  float64 `json:"t_hot_out"`

	Sigma float64
}

func NewFuelDF(TInit, T0 float64, fuel fuel.GasFuel) FuelDF {
	return FuelDF{
		C:      fuel.Cp(T0),
		TInit:  TInit,
		T0:     T0,
		QLower: fuel.QLower(),
		L0:     fuel.GasMassTheory(gases.GetAir()),
	}
}

type FuelDF struct {
	C      float64 `json:"c"`
	TInit  float64 `json:"t_init"`
	T0     float64 `json:"t_0"`
	QLower float64 `json:"q_lower"`
	L0     float64 `json:"l_0"`
}

func NewBurnerDF(node constructive.BurnerNode) BurnerDF {
	var t0 = node.T0()
	var inletGas = node.GasInput().GetState().Value().(gases.Gas)
	var outletGas = node.GasOutput().GetState().Value().(gases.Gas)

	var df = BurnerDF{
		Tg:              node.TStagOut(),
		Eta:             node.Eta(),
		Alpha:           node.Alpha(),
		FuelMassRateRel: node.FuelRateRel(),
		Sigma:           node.Sigma(),

		Fuel: NewFuelDF(
			node.TFuel(),
			t0,
			node.Fuel(),
		),

		AirDataInlet:  NewGasDF(
			node.PressureInput().GetState().Value().(float64),
			node.TemperatureInput().GetState().Value().(float64),
			inletGas,
		),
		AirData0:      NewGasDF(
			node.PressureInput().GetState().Value().(float64),
			t0, inletGas,
		),
		GasData0:      NewGasDF(
			node.PressureInput().GetState().Value().(float64),
			t0, outletGas,
		),
		GasDataOutlet: NewGasDF(
			node.PressureOutput().GetState().Value().(float64),
			node.TemperatureOutput().GetState().Value().(float64),
			node.GasOutput().GetState().Value().(gases.Gas),
		),
	}

	df.A = df.GasDataOutlet.Cp*df.Tg - df.AirDataInlet.Cp*node.TemperatureInput().GetState().Value().(float64)
	df.B = (df.GasData0.Cp - df.AirData0.Cp) * t0
	df.C = df.GasDataOutlet.Cp*df.Tg - df.GasData0.Cp*t0
	df.D = df.Fuel.C * (node.TFuel() - t0)

	return df
}

type BurnerDF struct {
	Tg              float64 `json:"tg"`
	Eta             float64 `json:"eta"`
	Alpha           float64 `json:"alpha"`
	FuelMassRateRel float64 `json:"fuel_mass_rate_rel"`
	Sigma           float64 `json:"sigma"`

	A float64
	B float64
	C float64
	D float64

	Fuel FuelDF `json:"fuel"`

	AirDataInlet  GasDF `json:"air_data_inlet"`
	AirData0      GasDF `json:"air_data_0"`
	GasData0      GasDF `json:"gas_data_0"`
	GasDataOutlet GasDF `json:"gas_data_outlet"`
}

func NewTurbineDFFromBlockedTurbine(node constructive.BlockedTurbineNode) TurbineDF {
	inletGas := node.GasInput().GetState().Value().(gases.Gas)
	outletGas := node.GasOutput().GetState().Value().(gases.Gas)

	return TurbineDF{
		PIn:  node.PStagIn(),
		POut: node.PStagOut(),

		TIn:  node.TStagIn(),
		TOut: node.TStagOut(),

		InletGasData:  NewGasDF(node.PStagIn(), node.TStagIn(), inletGas),
		OutletGasData: NewGasDF(node.PStagOut(), node.TStagOut(), outletGas),
		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			outletGas,
		),

		MassRateRel:     node.MassRateRel(),
		LeakMassRateRel: node.LeakMassRateRel(),
		CoolMassRateRel: node.CoolMassRateRel(),

		LambdaOut: node.LambdaOut(),
		POutStat:  constructive.POut(node),
		TOutStat:  constructive.TOut(node),

		Labour: node.LSpecific(),
		Eta:    node.Eta(),
	}
}

func NewTurbineDFFromFreeTurbine(node constructive.FreeTurbineNode) TurbineDF {
	inletGas := node.GasInput().GetState().Value().(gases.Gas)
	outletGas := node.GasOutput().GetState().Value().(gases.Gas)

	return TurbineDF{
		PIn:  node.PStagIn(),
		POut: node.PStagOut(),

		TIn:  node.TStagIn(),
		TOut: node.TStagOut(),

		InletGasData:  NewGasDF(node.PStagIn(), node.TStagIn(), inletGas),
		OutletGasData: NewGasDF(node.PStagOut(), node.TStagOut(), outletGas),
		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			outletGas,
		),

		MassRateRel:     node.MassRateRel(),
		LeakMassRateRel: node.LeakMassRateRel(),
		CoolMassRateRel: node.CoolMassRateRel(),

		LambdaOut: node.LambdaOut(),
		POutStat:  constructive.POut(node),
		TOutStat:  constructive.TOut(node),

		Labour: node.LSpecific(),
		Eta:    node.Eta(),
	}
}

type TurbineDF struct {
	PIn  float64 `json:"p_in"`
	POut float64 `json:"p_out"`

	TIn  float64 `json:"t_in"`
	TOut float64 `json:"t_out"`

	InletGasData  GasDF     `json:"inlet_gas_data"`
	OutletGasData GasDF     `json:"outlet_gas_data"`
	GasData       GasMeanDF `json:"gas_data"`

	MassRateRel     float64 `json:"mass_rate_rel"`
	LeakMassRateRel float64 `json:"leak_mass_rate_rel"`
	CoolMassRateRel float64 `json:"cool_mass_rate_rel"`

	LambdaOut float64 `json:"lambda_out"`
	POutStat  float64 `json:"p_out_stat"`
	TOutStat  float64 `json:"t_out_stat"`

	Labour float64 `json:"labour"`
	Eta    float64 `json:"eta"`
}

func NewShaftDF(node constructive.TransmissionNode) ShaftDF {
	return ShaftDF{
		Eta: node.Eta(),
	}
}

type ShaftDF struct {
	Eta float64 `json:"eta"`
}

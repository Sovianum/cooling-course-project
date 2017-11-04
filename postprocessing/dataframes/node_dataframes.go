package dataframes

import (
	"github.com/Sovianum/turbocycle/helpers/fuel"
	"github.com/Sovianum/turbocycle/helpers/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
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
	T       float64
	P       float64
	Density float64
	K       float64
	Cp      float64
	R       float64
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
	T1     float64
	T2     float64
	P      float64
	KMean  float64
	CpMean float64
	R      float64
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

		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			node.ComplexGasInput().GetState().(states.ComplexGasPortState).Gas,
		),
	}
}

type CompressorDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	Pi  float64
	Eta float64

	Labour float64

	GasData GasMeanDF
}

func NewPressureDropDF(node constructive.PressureLossNode) PressureDropDF {
	return PressureDropDF{
		PIn:node.PStagIn(),
		POut:node.PStagOut(),
		TIn:node.TStagIn(),
		TOut:node.TStagOut(),
		Sigma:node.Sigma(),
	}
}

type PressureDropDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	Sigma float64
}

func NewRegeneratorNode(node constructive.RegeneratorNode) RegeneratorDF {
	var coldInputState = node.ColdInput().GetState().(states.ComplexGasPortState)
	var hotInputState = node.HotInput().GetState().(states.ComplexGasPortState)
	var coldOutputState = node.ColdOutput().GetState().(states.ComplexGasPortState)
	var hotOutputSTate = node.HotOutput().GetState().(states.ComplexGasPortState)

	return RegeneratorDF{
		PColdIn:  coldInputState.PStag,
		PColdOut: coldOutputState.PStag,
		PHotIn:   hotInputState.PStag,
		PHotOut:  hotOutputSTate.PStag,

		TColdIn:  coldInputState.TStag,
		TColdOut: coldOutputState.TStag,
		THotIn:   hotInputState.TStag,
		THotOut:  hotOutputSTate.TStag,

		Sigma: node.Sigma(),
	}
}

type RegeneratorDF struct {
	PColdIn  float64
	PColdOut float64
	PHotIn   float64
	PHotOut  float64

	TColdIn  float64
	TColdOut float64
	THotIn   float64
	THotOut  float64

	Sigma float64
}

func NewFuelDF(TInit, T0 float64, fuel fuel.GasFuel) FuelDF {
	return FuelDF{
		C:      fuel.Cp(T0),
		TInit:  TInit,
		T0:     T0,
		QLower: fuel.QLower(),
		L0:     fuel.AirMassTheory(),
	}
}

type FuelDF struct {
	C      float64
	TInit  float64
	T0     float64
	QLower float64
	L0     float64
}

func NewBurnerDF(node constructive.BurnerNode) BurnerDF {
	var inletGasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)
	var outletGasState = node.ComplexGasOutput().GetState().(states.ComplexGasPortState)
	var t0 = node.T0()
	var inletGas = inletGasState.Gas
	var outletGas = outletGasState.Gas

	return BurnerDF{
		Tg:              node.TStagOut(),
		Eta:             node.Eta(),
		Alpha:           node.Alpha(),
		FuelMassRateRel: node.GetFuelRateRel(),
		Sigma:           node.Sigma(),

		Fuel: NewFuelDF(
			node.TFuel(),
			t0,
			node.Fuel(),
		),

		AirDataInlet:  NewGasDF(inletGasState.PStag, inletGasState.TStag, inletGas),
		AirData0:      NewGasDF(inletGasState.PStag, t0, inletGas),
		GasData0:      NewGasDF(inletGasState.PStag, t0, outletGas),
		GasDataOutlet: NewGasDF(outletGasState.PStag, outletGasState.TStag, outletGasState.Gas),
	}
}

type BurnerDF struct {
	Tg              float64
	Eta             float64
	Alpha           float64
	FuelMassRateRel float64
	Sigma           float64

	Fuel FuelDF

	AirDataInlet  GasDF
	AirData0      GasDF
	GasData0      GasDF
	GasDataOutlet GasDF
}

func NewTurbineDFFromBlockedTurbine(node constructive.BlockedTurbineNode) TurbineDF {
	var inletGasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)
	var outletGasState = node.ComplexGasOutput().GetState().(states.ComplexGasPortState)

	return TurbineDF{
		PIn:  node.PStagIn(),
		POut: node.PStagOut(),

		TIn:  node.TStagIn(),
		TOut: node.TStagOut(),

		InletGasData:  NewGasDF(node.PStagIn(), node.TStagIn(), inletGasState.Gas),
		OutletGasData: NewGasDF(node.PStagOut(), node.TStagOut(), outletGasState.Gas),
		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			outletGasState.Gas,
		),

		MassRateRel:     node.MassRateRel(),
		LeakMassRateRel: node.LeakMassRateRel(),
		CoolMassRateRel: node.CoolMassRateRel(),

		LambdaOut: node.LambdaOut(),
		POutStat:  node.PStatOut(),
		TOutStat:  node.TStatOut(),

		Labour: node.LSpecific(),
		Eta:    node.Eta(),
	}
}

func NewTurbineDFFromFreeTurbine(node constructive.FreeTurbineNode) TurbineDF {
	var inletGasState = node.ComplexGasInput().GetState().(states.ComplexGasPortState)
	var outletGasState = node.GasOutput().GetState().(states.GasPortState)

	return TurbineDF{
		PIn:  node.PStagIn(),
		POut: node.PStagOut(),

		TIn:  node.TStagIn(),
		TOut: node.TStagOut(),

		InletGasData:  NewGasDF(node.PStagIn(), node.TStagIn(), inletGasState.Gas),
		OutletGasData: NewGasDF(node.PStagOut(), node.TStagOut(), outletGasState.Gas),
		GasData: NewGasMeanDF(
			node.PStagIn(),
			node.TStagIn(),
			node.TStagOut(),
			outletGasState.Gas,
		),

		MassRateRel:     node.MassRateRel(),
		LeakMassRateRel: node.LeakMassRateRel(),
		CoolMassRateRel: node.CoolMassRateRel(),

		LambdaOut: node.LambdaOut(),
		POutStat:  node.PStatOut(),
		TOutStat:  node.TStatOut(),

		Labour: node.LSpecific(),
		Eta:    node.Eta(),
	}
}

type TurbineDF struct {
	PIn  float64
	POut float64

	TIn  float64
	TOut float64

	InletGasData  GasDF
	OutletGasData GasDF
	GasData       GasMeanDF

	MassRateRel     float64
	LeakMassRateRel float64
	CoolMassRateRel float64

	LambdaOut float64
	POutStat  float64
	TOutStat  float64

	Labour float64
	Eta    float64
}

func NewShaftDF(node constructive.TransmissionNode) ShaftDF {
	return ShaftDF{
		Eta: node.Eta(),
	}
}

type ShaftDF struct {
	Eta float64
}

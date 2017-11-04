package dataframes

import (
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
)


func NewThreeShaftsDF(power float64, scheme schemes.ThreeShaftsScheme) ThreeShaftsDF {
	var gasSourceState = scheme.GasSource().ComplexGasOutput().GetState().(states.ComplexGasPortState)

	return ThreeShaftsDF{
		GasSource:NewGasDF(gasSourceState.PStag, gasSourceState.TStag, gasSourceState.Gas),
		InletPipe:NewPressureDropDF(scheme.InletPressureDrop()),

		LPCompressor:NewCompressorDF(scheme.LowPressureCompressor()),
		LPCompressorPipe:NewPressureDropDF(scheme.MiddlePressureCompressorPipe()),
		LPTurbine:NewTurbineDFFromBlockedTurbine(
			scheme.MiddlePressureCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		LPTurbinePipe:NewPressureDropDF(scheme.MiddlePressureTurbinePipe()),
		LPShaft:NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		HPCompressor:NewCompressorDF(scheme.HighPressureCompressor()),
		HPTurbine:NewTurbineDFFromBlockedTurbine(
			scheme.GasGenerator().TurboCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		HPTurbinePipe:NewPressureDropDF(scheme.HighPressureTurbinePipe()),
		HPShaft:NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		Burner:NewBurnerDF(scheme.GasGenerator().Burner()),
		FreeTurbine:NewTurbineDFFromFreeTurbine(scheme.FreeTurbineBlock().FreeTurbine()),

		EngineLabour:scheme.FreeTurbineBlock().FreeTurbine().LSpecific(),
		Ce:schemes.GetSpecificFuelRate(scheme),
		MassRate:schemes.GetMassRate(power, scheme),
		Eta:schemes.GetEfficiency(scheme),
		Ne:power,
	}
}

type ThreeShaftsDF struct {
	GasSource   GasDF
	InletPipe   PressureDropDF

	LPCompressor     CompressorDF
	LPCompressorPipe PressureDropDF
	LPTurbine        TurbineDF
	LPTurbinePipe    PressureDropDF
	LPShaft          ShaftDF

	HPCompressor  CompressorDF
	HPTurbine     TurbineDF
	HPTurbinePipe PressureDropDF
	HPShaft       ShaftDF

	Burner BurnerDF

	FreeTurbine TurbineDF
	OutletPipe  PressureDropDF

	EngineLabour float64
	Ce           float64
	MassRate     float64
	Eta          float64
	Ne           float64
}

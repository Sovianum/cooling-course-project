package dataframes

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/impl/engine/states"
	"github.com/Sovianum/turbocycle/library/schemes"
)

func NewThreeShaftsDF(nE float64, etaR float64, scheme schemes.ThreeShaftsScheme) ThreeShaftsDF {
	var gasSourceState = scheme.GasSource().ComplexGasOutput().GetState().(states.ComplexGasPortState)

	return ThreeShaftsDF{
		GasSource: NewGasDF(gasSourceState.PStag, gasSourceState.TStag, gasSourceState.Gas),
		InletPipe: NewPressureDropDF(scheme.InletPressureDrop()),

		LPCompressor:     NewCompressorDF(scheme.LowPressureCompressor()),
		LPCompressorPipe: NewPressureDropDF(scheme.MiddlePressureCompressorPipe()),
		LPTurbine: NewTurbineDFFromBlockedTurbine(
			scheme.MiddlePressureCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		LPTurbinePipe: NewPressureDropDF(scheme.MiddlePressureTurbinePipe()),
		LPShaft:       NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		HPCompressor: NewCompressorDF(scheme.HighPressureCompressor()),
		HPTurbine: NewTurbineDFFromBlockedTurbine(
			scheme.GasGenerator().TurboCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		HPTurbinePipe: NewPressureDropDF(scheme.HighPressureTurbinePipe()),
		HPShaft:       NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		Burner:      NewBurnerDF(scheme.GasGenerator().Burner()),
		FreeTurbine: NewTurbineDFFromFreeTurbine(scheme.FreeTurbineBlock().FreeTurbine()),
		OutletPipe:  NewPressureDropDF(scheme.FreeTurbineBlock().OutletPressureLoss()),

		EngineLabour: scheme.FreeTurbineBlock().FreeTurbine().LSpecific(),
		Ce:           schemes.GetSpecificFuelRate(scheme),
		MassRate:     schemes.GetMassRate(nE, scheme),
		Eta:          schemes.GetEfficiency(scheme),
		Ne:           nE,

		EtaR:   etaR,
		NeMech: nE / etaR,
	}
}

type ThreeShaftsDF struct {
	GasSource GasDF          `json:"gas_source"`
	InletPipe PressureDropDF `json:"inlet_pipe"`

	LPCompressor     CompressorDF   `json:"lp_compressor"`
	LPCompressorPipe PressureDropDF `json:"lp_compressor_pipe"`
	LPTurbine        TurbineDF      `json:"lp_turbine"`
	LPTurbinePipe    PressureDropDF `json:"lp_turbine_pipe"`
	LPShaft          ShaftDF        `json:"lp_shaft"`

	HPCompressor  CompressorDF   `json:"hp_compressor"`
	HPTurbine     TurbineDF      `json:"hp_turbine"`
	HPTurbinePipe PressureDropDF `json:"hp_turbine_pipe"`
	HPShaft       ShaftDF        `json:"hp_shaft"`

	Burner BurnerDF `json:"burner"`

	FreeTurbine TurbineDF      `json:"free_turbine"`
	OutletPipe  PressureDropDF `json:"outlet_pipe"`

	EngineLabour float64 `json:"engine_labour"`
	Ce           float64 `json:"ce"`
	MassRate     float64 `json:"mass_rate"`
	Eta          float64 `json:"eta"`
	Ne           float64 `json:"ne"`

	EtaR   float64 `json:"eta_r"`
	NeMech float64 `json:"ne_mech"`
}

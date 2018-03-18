package dataframes

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewThreeShaftsDF(nE float64, etaR float64, scheme schemes.ThreeShaftsScheme) ThreeShaftsDF {
	gs := scheme.GasSource()
	pStag := gs.PressureOutput().GetState().Value().(float64)
	tStag := gs.PressureOutput().GetState().Value().(float64)
	gas := gs.GasOutput().GetState().Value().(gases.Gas)

	return ThreeShaftsDF{
		GasSource: NewGasDF(pStag, tStag, gas),
		InletPipe: NewPressureDropDF(scheme.InletPressureDrop()),

		LPCompressor:     NewCompressorDF(scheme.LPC()),
		LPCompressorPipe: NewPressureDropDF(scheme.LPCPipe()),
		LPTurbine: NewTurbineDFFromBlockedTurbine(
			scheme.MiddlePressureCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		LPTurbinePipe: NewPressureDropDF(scheme.LPTPipe()),
		LPShaft:       NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		HPCompressor: NewCompressorDF(scheme.HPC()),
		HPTurbine: NewTurbineDFFromBlockedTurbine(
			scheme.GasGenerator().TurboCascade().Turbine().(constructive.BlockedTurbineNode),
		),
		HPTurbinePipe: NewPressureDropDF(scheme.HPTPipe()),
		HPShaft:       NewShaftDF(scheme.MiddlePressureCascade().Transmission()),

		Burner:      NewBurnerDF(scheme.GasGenerator().Burner()),
		FreeTurbine: NewTurbineDFFromFreeTurbine(scheme.FTBlock().FreeTurbine()),
		OutletPipe:  NewPressureDropDF(scheme.FTBlock().OutletPressureLoss()),

		EngineLabour: scheme.FTBlock().FreeTurbine().LSpecific(),
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

package p3n

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
)

func NewData3n() Data3n {
	return Data3n{}
}

type Data3n struct {
	T common.FloatArr `json:"t"`
	Power common.FloatArr `json:"power"`
	MassRate common.FloatArr `json:"mass_rate"`
	Eta common.FloatArr `json:"eta"`

	PiLPC common.FloatArr `json:"pi_lpc"`
	PiHPC common.FloatArr `json:"pi_hpc"`

	PiHPT common.FloatArr `json:"pi_hpt"`
	PiLPT common.FloatArr `json:"pi_hpt"`
	PiFT  common.FloatArr `json:"pi_ft"`

	GNormHPT common.FloatArr `json:"g_norm_hpt"`
	GNormLPT common.FloatArr `json:"g_norm_lpt"`
	GNormFT  common.FloatArr `json:"g_norm_ft"`

	RpmHPT common.FloatArr `json:"rpm_hpt"`
	RpmLPT common.FloatArr `json:"rpm_lpt"`
	RpmFT  common.FloatArr `json:"rpm_ft"`
}

func (data *Data3n) Load(scheme free3n.ThreeShaftFreeScheme) {
	t := scheme.TemperatureSource().GetTemperature()
	labour := scheme.FT().PowerOutput().GetState().Value().(float64)
	massRate := scheme.LPC().MassRate()

	data.T.Append(t)

	data.Power.Append(labour*massRate/1e6)
	data.MassRate.Append(massRate)
	data.Eta.Append(scheme.Efficiency())

	data.PiLPC.Append(scheme.LPC().PiStag())
	data.PiHPC.Append(scheme.HPC().PiStag())

	data.PiHPT.Append(scheme.HPT().PiTStag())
	data.PiLPT.Append(scheme.LPT().PiTStag())
	data.PiFT.Append(scheme.FT().PiTStag())

	data.GNormHPT.Append(scheme.HPT().NormMassRate())
	data.GNormLPT.Append(scheme.LPT().NormMassRate())
	data.GNormFT.Append(scheme.FT().NormMassRate())

	data.RpmHPT.Append(scheme.HPT().RPMInput().GetState().Value().(float64))
	data.RpmLPT.Append(scheme.LPT().RPMInput().GetState().Value().(float64))
	data.RpmFT.Append(scheme.FT().RPMInput().GetState().Value().(float64))
}

package p2n

import (
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
)

func NewData2n() Data2n {
	return Data2n{}
}

type Data2n struct {
	Power    common.FloatArr `json:"power"`
	MassRate common.FloatArr `json:"mass_rate"`
	Eta      common.FloatArr `json:"eta"`

	T        common.FloatArr `json:"t"`
	PiC      common.FloatArr `json:"pi_c"`
	PiTC     common.FloatArr `json:"pi_tc"`
	PiF      common.FloatArr `json:"pi_f"`
	GNormTC  common.FloatArr `json:"g_norm_tc"`
	GNormTF  common.FloatArr `json:"g_norm_tf"`
	RpmTC    common.FloatArr `json:"rpm_tc"`
	RpmFT    common.FloatArr `json:"rpm_ft"`
}

func (data *Data2n) Load(scheme free2n.DoubleShaftFreeScheme) {
	t := scheme.TemperatureSource().GetTemperature()
	labour := scheme.FreeTurbine().PowerOutput().GetState().Value().(float64)
	freeTurbineMassRate := scheme.FreeTurbine().MassRateInput().GetState().Value().(float64)

	normMassRateTC := scheme.CompressorTurbine().NormMassRate()
	normMassRateFT := scheme.FreeTurbine().NormMassRate()

	data.Power.Append(labour * freeTurbineMassRate / 1e6)
	data.MassRate.Append(scheme.Compressor().MassRate())
	data.Eta.Append(scheme.Efficiency())

	data.T.Append(t)
	data.PiC.Append(scheme.Compressor().PiStag())
	data.PiTC.Append(scheme.CompressorTurbine().PiTStag())
	data.PiF.Append(scheme.FreeTurbine().PiTStag())
	data.GNormTC.Append(normMassRateTC)
	data.GNormTF.Append(normMassRateFT)
	data.RpmTC.Append(scheme.CompressorTurbine().RPMInput().GetState().Value().(float64))
	data.RpmFT.Append(scheme.FreeTurbine().RPMInput().GetState().Value().(float64))
}

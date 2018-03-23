package subcompress

import (
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	power = 16e6
)

type SchemeData struct {
	core.DoubleCompressorData
	SubCompressorPi []float64 `json:"sub_compressor_pi"`
	SplitFactor     []float64 `json:"split_factor"`
	SubCoolerT      []float64 `json:"sub_cooler_t"`
}

func updateSchemeData(scheme schemes.ThreeShaftsSubCompressScheme, data SchemeData) error {
	n, e := scheme.GetNetwork()
	if e != nil {
		return e
	}
	for i := range data.Pi {
		scheme.LPC().SetPiStag(data.PiLow[i])
		scheme.HPC().SetPiStag(data.PiHigh[i])
		scheme.SubCompressor().SetPiStag(data.SubCompressorPi[i])
		scheme.GasSplitter().SetExtraWeight(data.SplitFactor[i])

		if e := n.Solve(1, 2, 100, 1e-2); e != nil {
			return e
		}

		data.PiTLow[i] = scheme.LPT().PiTStag()
		data.PiTHigh[i] = scheme.HPT().PiTStag()
		data.Efficiency[i] = schemes.GetEfficiency(scheme)
		data.MassRate[i] = schemes.GetMassRate(power, scheme)
		data.SpecificPower[i] = scheme.GetSpecificPower()
		data.LabourLPC[i] = constructive.CompressorLabour(scheme.LPC())
		data.LabourHPC[i] = constructive.CompressorLabour(scheme.HPC())
		data.LabourLPT[i] = constructive.TurbineLabour(scheme.LPT())
		data.LabourHPT[i] = constructive.TurbineLabour(scheme.HPT())
		data.LabourFT[i] = constructive.TurbineLabour(scheme.FT())
		data.Heat[i] = constructive.FuelMassRate(scheme.MainBurner()) * scheme.GetQLower()
		data.SubCoolerT[i] = scheme.SubCooler().TemperatureOutput().GetState().Value().(float64)
	}
	return nil
}

func getSchemeDataTemplate(totalPiArr, piFactorArr, subCompressPiArr, splitFactorArr []float64) SchemeData {
	totalLen := len(totalPiArr) * len(piFactorArr) * len(subCompressPiArr) * len(splitFactorArr)
	result := SchemeData{
		DoubleCompressorData: core.NewDoubleCompressorSchemeData(totalLen),
		SubCompressorPi:      make([]float64, totalLen),
		SplitFactor:          make([]float64, totalLen),
		SubCoolerT:           make([]float64, totalLen),
	}

	cnt := 0
	for _, pi := range totalPiArr {
		for _, piFactor := range piFactorArr {
			piLow, piHigh := core.GetCompressorPiPair(pi, piFactor)
			for _, subCompressPi := range subCompressPiArr {
				for _, splitFactor := range splitFactorArr {
					result.Pi[cnt] = pi
					result.PiFactor[cnt] = piFactor
					result.PiLow[cnt] = piLow
					result.PiHigh[cnt] = piHigh
					result.SubCompressorPi[cnt] = subCompressPi
					result.SplitFactor[cnt] = splitFactor
					cnt++
				}
			}
		}
	}
	return result
}

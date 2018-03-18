package core

import (
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/schemes"
	"math"
)

type DoubleCompressorScheme interface {
	schemes.Scheme
	schemes.DoubleCompressor
	MainBurner() constructive.BurnerNode
	HPT() constructive.StaticTurbineNode
	LPT() constructive.StaticTurbineNode
	FT() constructive.StaticTurbineNode
}

func EtaOptimalPoint(points []DoubleCompressorDataPoint) DoubleCompressorDataPoint {
	var eta = -1.
	var ind = 0

	for i, point := range points {
		if point.Efficiency > eta {
			eta = point.Efficiency
			ind = i
		}
	}

	return points[ind]
}

func LabourOptimalPoint(points []DoubleCompressorDataPoint) DoubleCompressorDataPoint {
	var labour = -1.
	var ind = 0

	for i, point := range points {
		if point.SpecificPower > labour {
			labour = point.SpecificPower
			ind = i
		}
	}

	return points[ind]
}

type DoubleCompressorData struct {
	Pi            []float64 `json:"pi"`
	PiFactor      []float64 `json:"pi_factor"`
	MassRate      []float64 `json:"mass_rate"`
	SpecificPower []float64 `json:"specific_power"`
	Efficiency    []float64 `json:"efficiency"`
	PiLow         []float64 `json:"pi_low"`
	PiHigh        []float64 `json:"pi_high"`
	PiTLow        []float64 `json:"pi_t_low"`
	PiTHigh       []float64 `json:"pi_t_high"`
	LabourHPC     []float64 `json:"labour_hpc"`
	LabourLPC     []float64 `json:"labour_lpc"`
	LabourLPT     []float64 `json:"labour_lpt"`
	LabourHPT     []float64 `json:"labour_hpt"`
	LabourFT      []float64 `json:"labour_ft"`
	Heat          []float64 `json:"heat"`
}

func ConvertDoubleCompressorDataPoints(points []DoubleCompressorDataPoint) DoubleCompressorData {
	result := DoubleCompressorData{
		Pi:            make([]float64, len(points)),
		PiFactor:      make([]float64, len(points)),
		MassRate:      make([]float64, len(points)),
		SpecificPower: make([]float64, len(points)),
		Efficiency:    make([]float64, len(points)),
		PiLow:         make([]float64, len(points)),
		PiHigh:        make([]float64, len(points)),
		PiTLow:        make([]float64, len(points)),
		PiTHigh:       make([]float64, len(points)),
		LabourHPC:     make([]float64, len(points)),
		LabourLPC:     make([]float64, len(points)),
		LabourLPT:     make([]float64, len(points)),
		LabourHPT:     make([]float64, len(points)),
		LabourFT:      make([]float64, len(points)),
		Heat:          make([]float64, len(points)),
	}
	for i, point := range points {
		result.Pi[i] = point.Pi
		result.PiFactor[i] = point.PiFactor
		result.MassRate[i] = point.MassRate
		result.SpecificPower[i] = point.SpecificPower
		result.Efficiency[i] = point.Efficiency
		result.PiLow[i] = point.PiLow
		result.PiHigh[i] = point.PiHigh
		result.PiTLow[i] = point.PiTLow
		result.PiTHigh[i] = point.PiTHigh
		result.LabourHPC[i] = point.LabourHPC
		result.LabourLPC[i] = point.LabourLPC
		result.LabourHPT[i] = point.LabourHPT
		result.LabourLPT[i] = point.LabourLPT
		result.LabourFT[i] = point.LabourFT
		result.Heat[i] = point.Heat
	}
	return result
}

// todo remove (deprecated). use DoubleCompressorData instead
type DoubleCompressorDataPoint struct {
	Pi            float64
	PiFactor      float64
	MassRate      float64
	SpecificPower float64
	Efficiency    float64
	PiLow         float64
	PiHigh        float64
	PiTLow        float64
	PiTHigh       float64
	LabourHPC     float64
	LabourLPC     float64
	LabourLPT     float64
	LabourHPT     float64
	LabourFT      float64
	Heat          float64
}

func (point DoubleCompressorDataPoint) ToArray() []float64 {
	return []float64{
		point.Pi, point.PiFactor,
		point.MassRate, point.SpecificPower, point.Efficiency,
		point.PiLow, point.PiHigh,
		point.PiTLow, point.PiTHigh,
		point.LabourHPC, point.LabourLPC,
		point.LabourHPT, point.LabourLPT, point.LabourFT,
		point.Heat,
	}
}

func GetDoubleCompressorDataGenerator(
	scheme DoubleCompressorScheme, power float64, relaxCoef float64, iterNum int,
) func(pi, piFactor float64) (DoubleCompressorDataPoint, error) {
	return func(pi, piFactor float64) (DoubleCompressorDataPoint, error) {
		var piLow, piHigh = getCompressorPiPair(pi, piFactor)

		scheme.LPC().SetPiStag(piLow)
		scheme.HPC().SetPiStag(piHigh)
		network, netErr := scheme.GetNetwork()
		if netErr != nil {
			panic(netErr)
		}

		var err = network.Solve(relaxCoef, 2, iterNum, 0.001)
		if err != nil {
			return DoubleCompressorDataPoint{}, err
		}

		return DoubleCompressorDataPoint{
			Pi:            pi,
			PiFactor:      piFactor,
			PiLow:         piLow,
			PiHigh:        piHigh,
			PiTLow:        scheme.LPT().PiTStag(),
			PiTHigh:       scheme.HPT().PiTStag(),
			Efficiency:    schemes.GetEfficiency(scheme),
			MassRate:      schemes.GetMassRate(power, scheme),
			SpecificPower: scheme.GetSpecificPower(),
			LabourLPC:     constructive.CompressorLabour(scheme.LPC()),
			LabourHPC:     constructive.CompressorLabour(scheme.HPC()),
			LabourLPT:     constructive.TurbineLabour(scheme.LPT()),
			LabourHPT:     constructive.TurbineLabour(scheme.HPT()),
			LabourFT:      constructive.TurbineLabour(scheme.FT()),
			Heat:          constructive.FuelMassRate(scheme.MainBurner()) * scheme.GetQLower(),
		}, nil
	}
}

func getCompressorPiPair(piTotal, piFactor float64) (float64, float64) {
	piLow := math.Pow(piTotal, piFactor)
	piHigh := math.Pow(piTotal, 1 - piFactor)
	return piLow, piHigh
}

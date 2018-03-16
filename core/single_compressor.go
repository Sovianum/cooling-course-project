package core

import (
	"github.com/Sovianum/turbocycle/library/schemes"
	"strconv"
)

type SingleCompressorScheme interface {
	schemes.Scheme
	schemes.SingleCompressor
}

type SingleCompressorData struct {
	Pi            []float64
	MassRate      []float64
	SpecificPower []float64
	Efficiency    []float64
}

func ConvertSingleCompressorDataPoint(points []SingleCompressorDataPoint) SingleCompressorData {
	result := SingleCompressorData{
		Pi:            make([]float64, len(points)),
		MassRate:      make([]float64, len(points)),
		SpecificPower: make([]float64, len(points)),
		Efficiency:    make([]float64, len(points)),
	}
	for i, point := range points {
		result.Pi[i] = point.Pi
		result.MassRate[i] = point.MassRate
		result.SpecificPower[i] = point.SpecificPower
		result.Efficiency[i] = point.Efficiency
	}
	return result
}

// todo remove this structure and use SingleCompressorData instead
type SingleCompressorDataPoint struct {
	Pi            float64
	MassRate      float64
	SpecificPower float64
	Efficiency    float64
}

func (point SingleCompressorDataPoint) ToRecord() []string {
	return []string{
		strconv.FormatFloat(point.Pi, 'f', -1, 64),
		strconv.FormatFloat(point.MassRate, 'f', -1, 64),
		strconv.FormatFloat(point.SpecificPower, 'f', -1, 64),
		strconv.FormatFloat(point.Efficiency, 'f', -1, 64),
	}
}

func GetSingleCompressorDataGenerator(
	scheme SingleCompressorScheme, power float64, relaxCoef float64, iterNum int,
) func(pi float64) (SingleCompressorDataPoint, error) {
	return func(pi float64) (SingleCompressorDataPoint, error) {
		scheme.Compressor().SetPiStag(pi)
		network, netErr := scheme.GetNetwork()
		if netErr != nil {
			panic(netErr)
		}

		var err = network.Solve(relaxCoef, 2, iterNum, 0.001)
		if err != nil {
			return SingleCompressorDataPoint{}, err
		}

		return SingleCompressorDataPoint{
			Pi:            pi,
			Efficiency:    schemes.GetEfficiency(scheme),
			MassRate:      schemes.GetMassRate(power, scheme),
			SpecificPower: scheme.GetSpecificPower(),
		}, nil
	}
}

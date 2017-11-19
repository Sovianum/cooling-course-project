package core

import (
	"github.com/Sovianum/turbocycle/library/schemes"
	"errors"
)

type DoubleCompressorScheme interface {
	schemes.Scheme
	schemes.DoubleCompressor
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

type DoubleCompressorDataPoint struct {
	Pi            float64
	PiLow         float64
	PiHigh        float64
	PiFactor      float64
	MassRate      float64
	SpecificPower float64
	Efficiency    float64
}

func (point DoubleCompressorDataPoint) ToArray() []float64 {
	return []float64{
		point.Pi, point.PiFactor, point.MassRate, point.SpecificPower, point.Efficiency,
	}
}

func GetDoubleCompressorDataGenerator(
	scheme DoubleCompressorScheme, power float64, relaxCoef float64, iterNum int,
) func(pi, piFactor float64) (DoubleCompressorDataPoint, error) {
	return func(pi, piFactor float64) (DoubleCompressorDataPoint, error) {
		scheme.LowPressureCompressor().SetPiStag(pi * piFactor)
		scheme.HighPressureCompressor().SetPiStag(1 / piFactor)
		var converged, err = scheme.GetNetwork().Solve(relaxCoef, iterNum, 0.001)
		if err != nil {
			return DoubleCompressorDataPoint{}, err
		}
		if !converged {
			return DoubleCompressorDataPoint{}, errors.New("not converged")
		}

		return DoubleCompressorDataPoint{
			Pi:            pi,
			PiFactor:      piFactor,
			PiLow:         (pi - 1) * piFactor + 1,
			PiHigh:        pi / ((pi - 1) * piFactor + 1),
			Efficiency:    schemes.GetEfficiency(scheme),
			MassRate:      schemes.GetMassRate(power, scheme),
			SpecificPower: scheme.GetSpecificPower(),
		}, nil
	}
}


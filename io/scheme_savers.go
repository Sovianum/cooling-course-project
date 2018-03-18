package io

import (
	"github.com/Sovianum/cooling-course-project/core"
	"github.com/Sovianum/turbocycle/library/schemes"
)

const (
	relaxCoef = 0.1
	iterNum   = 100
)

func GetThreeShaftsSchemeData(
	scheme schemes.ThreeShaftsScheme,
	power float64,
	startPi, piStep float64, piStepNum int,
	startPiFactor, piFactorStep float64, piFactorStepNum int,
) ([]core.DoubleCompressorDataPoint, error) {
	var piArr []float64
	for i := 0; i != piStepNum; i++ {
		piArr = append(piArr, startPi+float64(i)*piStep)
	}

	var piFactorArr []float64
	for i := 0; i != piFactorStepNum; i++ {
		piFactorArr = append(piFactorArr, startPiFactor+float64(i)*piFactorStep)
	}

	var points []core.DoubleCompressorDataPoint
	var generator = core.GetDoubleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	for _, piFactor := range piFactorArr {
		for _, pi := range piArr {
			var point, err = generator(pi, piFactor)
			if err != nil {
				return nil, err
			}
			points = append(points, point)
		}
	}
	return points, nil
}

func GetTwoShaftSchemeData(
	scheme core.SingleCompressorScheme,
	power float64,
	startPi, piStep float64,
	stepNum int,
) ([]core.SingleCompressorDataPoint, error) {
	piArr := make([]float64, stepNum)
	for i := range piArr{
		piArr[i] = startPi+float64(i)*piStep
	}

	points := make([]core.SingleCompressorDataPoint, stepNum)
	generator := core.GetSingleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)

	for i, pi := range piArr {
		var point, err = generator(pi)
		if err != nil {
			return nil, err
		}
		points[i] = point
	}

	return points, nil
}

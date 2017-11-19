package io

import (
	"os"
	"encoding/csv"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/cooling-course-project/core"
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

func SaveTwoShaftSchemeData(
	scheme core.SingleCompressorScheme,
	power float64,
	startPi, piStep float64,
	stepNum int, filename string,
) error {
	var piArr []float64

	for i := 0; i != stepNum; i++ {
		piArr = append(piArr, startPi+float64(i)*piStep)
	}

	var records [][]string
	var generator = core.GetSingleCompressorDataGenerator(scheme, power, relaxCoef, iterNum)
	for _, pi := range piArr {
		var point, err = generator(pi)
		if err != nil {
			return err
		}
		records = append(records, point.ToRecord())
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.WriteAll(records)

	return nil
}
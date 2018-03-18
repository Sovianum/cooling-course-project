package dataframes

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/profile"
)

type TProfileCalcDF struct {
	Geom       TProfileGeomDF
	Gas        TProfileGasDF
	PSSolution profile.TemperatureSolution
	SSSolution profile.TemperatureSolution
}

type TProfileGeomDF struct {
	DInlet float64
}

type TProfileGasDF struct {
	RhoGas float64
	Ca     float64

	MuGas float64

	LambdaGas float64

	AlphaMean float64

	AlphaGasInlet  float64
	AlphaGasSS     float64
	AlphaGasPS     float64
	AlphaGasOutlet float64

	LengthPSArr   []float64
	AlphaAirPSArr []float64
	AlphaGasPSArr []float64
	TAirPSArr     []float64
	TWallPSArr    []float64

	LengthSSArr   []float64
	AlphaAirSSArr []float64
	AlphaGasSSArr []float64
	TAirSSArr     []float64
	TWallSSArr    []float64

	SkipSteps int
}

func (df *TProfileGasDF) SetSSSolutionInfo(solution profile.TemperatureSolution) {
	df.LengthSSArr = solution.LengthCoord
	df.AlphaAirSSArr = solution.AlphaAir
	df.AlphaGasSSArr = solution.AlphaGas
	df.TAirSSArr = solution.AirTemperature
	df.TWallSSArr = solution.WallTemperature
}

func (df *TProfileGasDF) SetPSSolutionInfo(solution profile.TemperatureSolution) {
	df.LengthPSArr = solution.LengthCoord
	df.AlphaAirPSArr = solution.AlphaAir
	df.AlphaGasPSArr = solution.AlphaGas
	df.TAirPSArr = solution.AirTemperature
	df.TWallPSArr = solution.WallTemperature
}

type TProfileRow struct {
	Id       int
	X        float64
	AlphaAir float64
	AlphaGas float64
	TAir     float64
	TWall    float64
}

func (df TProfileGasDF) PSRows() chan TProfileRow {
	var rowFunc = func(ch chan TProfileRow) {
		for i, j := 0, 1; i < len(df.LengthPSArr); i, j = i+df.SkipSteps, j+1 {
			ch <- TProfileRow{
				Id:       j,
				X:        df.LengthPSArr[i],
				AlphaAir: df.AlphaAirPSArr[i],
				AlphaGas: df.AlphaGasPSArr[i],
				TAir:     df.TAirPSArr[i],
				TWall:    df.TWallPSArr[i],
			}
		}
		close(ch)
	}
	var result = make(chan TProfileRow)
	go rowFunc(result)
	return result
}

func (df TProfileGasDF) SSRows() chan TProfileRow {
	var rowFunc = func(ch chan TProfileRow) {
		for i, j := 0, 1; i < len(df.LengthPSArr); i, j = i+df.SkipSteps, j+1 {
			ch <- TProfileRow{
				Id:       j,
				X:        df.LengthSSArr[i],
				AlphaAir: df.AlphaAirSSArr[i],
				AlphaGas: df.AlphaGasSSArr[i],
				TAir:     df.TAirSSArr[i],
				TWall:    df.TWallSSArr[i],
			}
		}
		close(ch)
	}
	var result = make(chan TProfileRow)
	go rowFunc(result)
	return result
}

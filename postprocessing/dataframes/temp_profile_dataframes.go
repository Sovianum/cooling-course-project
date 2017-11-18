package dataframes

type TProfileCalcDF struct {
	Geom TProfileGeomDF
	Gas  TProfileGasDF
}

type TProfileGeomDF struct {
	DInlet float64
	DMean  float64
}

type TProfileGasDF struct {
	MassRateGas float64
	MassRateAir float64

	MuGas float64
	MuAir float64

	LambdaGas float64
	LambdaAir float64

	Alpha0 float64

	AlphaMean     float64
	AlphaGasInlet float64
	AlphaGasSS    float64
	AlphaGasPS    float64
	AlphaGasOutlet   float64

	AlphaAir float64

	LengthPSArr   []float64
	AlphaAirPSArr []float64
	AlphaGasPSArr []float64
	TAirPSArr     []float64

	LengthSSArr   []float64
	AlphaAirSSArr []float64
	AlphaGasSSArr []float64
	TAirSSArr     []float64
}

type TProfileRow struct {
	Id       int
	X        float64
	AlphaAir float64
	AlphaGas float64
	TAir     float64
}

func (df TProfileGasDF) PSRows() chan TProfileRow {
	var rowFunc = func(ch chan TProfileRow) {
		for i := range df.LengthPSArr {
			ch <- TProfileRow{
				Id:       i,
				X:        df.LengthPSArr[i],
				AlphaAir: df.AlphaAirPSArr[i],
				AlphaGas: df.AlphaGasPSArr[i],
				TAir:     df.TAirPSArr[i],
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
		for i := range df.LengthPSArr {
			ch <- TProfileRow{
				Id:       i,
				X:        df.LengthSSArr[i],
				AlphaAir: df.AlphaAirSSArr[i],
				AlphaGas: df.AlphaGasSSArr[i],
				TAir:     df.TAirSSArr[i],
			}
		}
		close(ch)
	}
	var result = make(chan TProfileRow)
	go rowFunc(result)
	return result
}
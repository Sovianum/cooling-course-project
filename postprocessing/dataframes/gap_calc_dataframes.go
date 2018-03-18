package dataframes

import (
	"github.com/Sovianum/turbocycle/utils/turbine/cooling/gap"
)

func GapCalcFromDataPacks(packArr []gap.DataPack) GapCalcDF {
	return GapCalcDF{
		Geom:  GapGeomFromDataPack(packArr[0]),
		Metal: GapMetalFromDataPack(packArr[0]),
		Gas:   GapGasFromDataPacks(packArr),
	}
}

type GapCalcDF struct {
	Geom  GapGeometryDF
	Metal GapMetalDF
	Gas   GapGasDF
}

func GapGeomFromDataPack(pack gap.DataPack) GapGeometryDF {
	return GapGeometryDF{
		BladeLength:     pack.BladeLength,
		ChordProjection: pack.ChordProjection,
		BladeArea:       pack.BladeArea,
		Perimeter:       pack.Perimeter,
		WallThk:         pack.WallThk,
	}
}

type GapGeometryDF struct {
	BladeLength     float64
	DMean           float64
	ChordProjection float64
	BladeArea       float64
	Perimeter       float64
	WallThk         float64
	DInlet          float64
}

func GapMetalFromDataPack(pack gap.DataPack) GapMetalDF {
	return GapMetalDF{
		TWallOuter: pack.TWallOuter,
		TWallInner: pack.TWallInner,
		TWallMean:  pack.TMean,
		DTWall:     pack.TDrop,
		LambdaM:    pack.LambdaM,
	}
}

type GapMetalDF struct {
	TWallOuter float64
	TWallInner float64
	TWallMean  float64
	DTWall     float64
	LambdaM    float64
}

func GapGasFromDataPacks(packArr []gap.DataPack) GapGasDF {
	var airMassRate = make([]float64, len(packArr))
	var dCoef = make([]float64, len(packArr))
	var epsCoef = make([]float64, len(packArr))
	var airGap = make([]float64, len(packArr))

	for i, pack := range packArr {
		airMassRate[i] = pack.MassRateCooler
		dCoef[i] = pack.DComplex
		epsCoef[i] = pack.EpsComplex
		airGap[i] = pack.AirGap
	}

	return GapGasDF{
		Tg:         packArr[0].TGas,
		CaGas:      packArr[0].CaGas,
		DensityGas: packArr[0].DensityGas,
		MuGas:      packArr[0].MuGas,
		LambdaGas:  packArr[0].LambdaGas,

		ReGas: packArr[0].ReGas,
		NuGas: packArr[0].NuGas,

		Theta0:   packArr[0].TAir0,
		AlphaGas: packArr[0].AlphaGas,
		Heat:     packArr[0].BladeHeat,

		AirMassRate: airMassRate,
		DCoef:       dCoef,
		EpsCoef:     epsCoef,
		AirGap:      airGap,
	}
}

type GapGasDF struct {
	Tg         float64
	CaGas      float64
	DensityGas float64
	MuGas      float64
	LambdaGas  float64

	GasMassRate float64
	ReGas       float64
	NuGas       float64
	NuCoef      float64

	Theta0 float64

	AlphaGas float64
	Heat     float64

	AirMassRate []float64
	DCoef       []float64
	EpsCoef     []float64
	AirGap      []float64
}

type GapTableRow struct {
	Id          int
	AirMassRate float64
	DCoef       float64
	EpsCoef     float64
	AirGap      float64
}

func (df GapGasDF) TableRows() chan GapTableRow {
	var iterFunc = func(ch chan GapTableRow) {
		for i := range df.DCoef {
			ch <- GapTableRow{
				Id:          i + 1,
				AirMassRate: df.AirMassRate[i],
				DCoef:       df.DCoef[i],
				EpsCoef:     df.EpsCoef[i],
				AirGap:      df.AirGap[i],
			}
		}
		close(ch)
	}

	var result = make(chan GapTableRow)
	go iterFunc(result)

	return result
}

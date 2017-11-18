package dataframes

type GapCalcDF struct {
	Geom  GapGeometryDF
	Metal GapMetalDF
	Gas   GapGasDF
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

type GapMetalDF struct {
	TWallOuter float64
	TWallInner float64
	TWallMean  float64
	DTWall     float64
	LambdaM    float64
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

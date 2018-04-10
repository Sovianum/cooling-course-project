package dataframes

import (
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
)

func NewCompressorBladingGeomDF(
	geom geometry.BladingGeometry,
	geomGen compressor.IncompleteBladingGeometryGenerator,
) CompressorBladingGeomDF {
	return CompressorBladingGeomDF{
		DRelIn:     geometry.DRel(0, geom),
		RRelIn:     geometry.RRel(geometry.DRel(0, geom)),
		DRelOut:    geometry.DRel(geom.XGapOut(), geom),
		RRelOut:    geometry.RRel(geometry.DRel(geom.XGapOut(), geom)),
		Elongation: geomGen.Elongation(),
		DeltaRel:   geomGen.DeltaRel(),
		GammaIn:    geomGen.GammaIn(),
		GammaOut:   geomGen.GammaOut(),
		AreaIn:     geometry.Area(0, geom),
		AreaOut:    geometry.Area(geom.XGapOut(), geom),
		DOutIn:     geom.OuterProfile().Diameter(0),
		DInIn:      geom.InnerProfile().Diameter(0),
		DOutOut:    geom.OuterProfile().Diameter(geom.XGapOut()),
		DInOut:     geom.InnerProfile().Diameter(geom.XGapOut()),
		BladeWidth: geom.XBladeOut(),
	}
}

type CompressorBladingGeomDF struct {
	DRelIn     float64
	RRelIn     float64
	DRelOut    float64
	RRelOut    float64
	Elongation float64
	DeltaRel   float64
	GammaIn    float64
	GammaOut   float64
	AreaIn     float64
	AreaOut    float64
	DOutIn     float64
	DInIn      float64
	DOutOut    float64
	DInOut     float64
	BladeWidth float64
}

type CompressorStageDF struct {
	StatorDF   CompressorBladingGeomDF
	RotorDF    CompressorBladingGeomDF
	BladeWidth float64

	Ht             float64
	HtCoefCurr     float64
	HtCoefNext     float64
	ReactivityCurr float64
	ReactivityNext float64
	UOut           float64
	Lz             float64
	Kh             float64
	HAd            float64
	Eta            float64
	DT             float64
	Pi             float64
	P1             float64
	P3             float64
	T1             float64
	T3             float64
	ACrit1         float64
	ACrit3         float64

	CpAir float64
	KAir  float64
	RAir  float64

	Q1 float64

	MassRate float64

	CURel1 float64
	CURel2 float64
	CURel3 float64
	CARel1 float64
	CARel2 float64
	CARel3 float64

	Lambda1 float64
	Lambda3 float64

	Triangle1 states.VelocityTriangle
	Triangle2 states.VelocityTriangle
	Triangle3 states.VelocityTriangle
}

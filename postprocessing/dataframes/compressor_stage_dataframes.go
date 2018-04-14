package dataframes

import (
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/stage/compressor"
	"github.com/Sovianum/turbocycle/impl/stage/geometry"
	"github.com/Sovianum/turbocycle/impl/stage/states"
	"github.com/Sovianum/turbocycle/material/gases"
)

func NewStagedCompressorDF(c compressor.StagedCompressorNode) StagedCompressorDF {
	result := StagedCompressorDF{StageDFS: make([]CompressorStageDF, len(c.Stages()))}
	for i, stage := range c.Stages() {
		result.StageDFS[i] = NewCompressorStageDF(stage)
	}
	return result
}

type StagedCompressorDF struct {
	StageDFS []CompressorStageDF

	rows []StageRow
}

func (df StagedCompressorDF) Rows() []StageRow {
	if df.rows == nil {
		getRow := func(name, dim string, f func(df CompressorStageDF) float64) StageRow {
			return NewCompressorStageRow(name, dim, df.StageDFS, f)
		}
		df.rows = []StageRow{
			getRow("$H_т$", "$10^5 \\cdot Дж/кг$", func(df CompressorStageDF) float64 { return df.Ht }).FormatFloat(DivideE5).FormatString(Round3),
			getRow("$L_z$", "$10^5 \\cdot Дж/кг$", func(df CompressorStageDF) float64 { return df.Lz }).FormatFloat(DivideE5).FormatString(Round3),
			getRow("$H_{ад}$", "$10^5 \\cdot Дж/кг$", func(df CompressorStageDF) float64 { return df.HAd }).FormatFloat(DivideE5).FormatString(Round3),
			getRow("$\\Delta T$", "$К$", func(df CompressorStageDF) float64 { return df.DT }).FormatString(Round1),
			getRow("$\\pi^*$", "-", func(df CompressorStageDF) float64 { return df.Pi }).FormatString(Round3),
			getRow("$p_1^*$", "$10^6 \\cdot Па$", func(df CompressorStageDF) float64 { return df.P1 }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$p_3^*$", "$10^6 \\cdot Па$", func(df CompressorStageDF) float64 { return df.P3 }).FormatFloat(DivideE6).FormatString(Round3),
			getRow("$a_{кр1}$", "$м/с$", func(df CompressorStageDF) float64 { return df.ACrit1 }).FormatString(Round2),
			getRow("$a_{кр3}$", "$м/с$", func(df CompressorStageDF) float64 { return df.ACrit3 }).FormatString(Round2),
			getRow("$\\overline{r_{ср1}}$", "-", func(df CompressorStageDF) float64 { return df.RotorDF.RRelIn }).FormatString(Round3),
			getRow("$\\overline{c_{u1}}$", "-", func(df CompressorStageDF) float64 { return df.CURel1 }).FormatString(Round3),
			getRow("$\\alpha_1$", "$\\degree$", func(df CompressorStageDF) float64 { return common.ToDegrees(df.Triangle1.Alpha()) }).FormatString(Round1),
			getRow("$\\lambda_1$", "-", func(df CompressorStageDF) float64 { return df.Lambda1 }).FormatString(Round2),
			getRow("$F_1$", "$м^2$", func(df CompressorStageDF) float64 { return df.RotorDF.AreaIn }).FormatString(Round3),
			getRow("$D_1$", "м", func(df CompressorStageDF) float64 { return df.RotorDF.DInOut }).FormatString(Round3),
			getRow("$d_1$", "м", func(df CompressorStageDF) float64 { return df.RotorDF.DInIn }).FormatString(Round3),
			getRow("$x_{ступ}$", "м", func(df CompressorStageDF) float64 { return df.StageWidth }).FormatString(Round3),
			getRow("$D_3$", "м", func(df CompressorStageDF) float64 { return df.StatorDF.DOutOut }).FormatString(Round3),
			getRow("$d_3$", "м", func(df CompressorStageDF) float64 { return df.StatorDF.DOutIn }).FormatString(Round3),
			getRow("$F_3$", "$м^2$", func(df CompressorStageDF) float64 { return df.StatorDF.AreaOut }).FormatString(Round3),
			getRow("$\\overline{d_3}$", "-", func(df CompressorStageDF) float64 { return df.StatorDF.DRelOut }).FormatString(Round3),
			getRow("$\\overline{r_{ср3}}$", "-", func(df CompressorStageDF) float64 { return df.StatorDF.RRelOut }).FormatString(Round3),
			getRow("$\\overline{c_{u3}}$", "-", func(df CompressorStageDF) float64 { return df.CURel3 }).FormatString(Round3),
			getRow("$\\lambda_3$", "-", func(df CompressorStageDF) float64 { return df.Lambda3 }).FormatString(Round3),
			getRow("$\\alpha_3$", "$\\degree$", func(df CompressorStageDF) float64 { return common.ToDegrees(df.Triangle3.Alpha()) }).FormatString(Round1),
			getRow("$\\overline{c_{u2}}$", "-", func(df CompressorStageDF) float64 { return df.CURel2 }).FormatString(Round2),
			getRow("$\\beta_1$", "$\\degree$", func(df CompressorStageDF) float64 { return common.ToDegrees(df.Triangle1.Beta()) }).FormatString(Round1),
			getRow("$\\beta_2$", "$\\degree$", func(df CompressorStageDF) float64 { return common.ToDegrees(df.Triangle2.Beta()) }).FormatString(Round1),
			getRow("$\\alpha_2$", "$\\degree$", func(df CompressorStageDF) float64 { return common.ToDegrees(df.Triangle2.Alpha()) }).FormatString(Round1),
			getRow("$w_2$", "$м/с$", func(df CompressorStageDF) float64 { return df.Triangle2.W() }).FormatString(Round2),
			getRow("$c_2$", "$м/с$", func(df CompressorStageDF) float64 { return df.Triangle2.C() }).FormatString(Round2),
		}

		for i := range df.rows {
			df.rows[i].ID = i + 1
		}
	}
	return df.rows
}

func NewCompressorStageDF(stage compressor.StageNode) CompressorStageDF {
	dp := stage.GetDataPack()
	return CompressorStageDF{
		StatorDF:   NewCompressorBladingGeomDF(dp.StageGeometry.StatorGeometry(), stage.GeomGen().StatorGenerator()),
		RotorDF:    NewCompressorBladingGeomDF(dp.StageGeometry.RotorGeometry(), stage.GeomGen().RotorGenerator()),
		StageWidth: dp.StageGeometry.RotorGeometry().XGapOut() + dp.StageGeometry.StatorGeometry().XGapOut(),
		RPM:        stage.RPM(),

		Ht:             dp.HT,
		HtCoefCurr:     stage.HtCoef(),
		HtCoefNext:     stage.HtCoefNext(),
		ReactivityCurr: stage.Reactivity(),
		ReactivityNext: stage.ReactivityNext(),
		UOut:           dp.UOut,
		Lz:             dp.Labour,
		Kh:             dp.LabourCoef,
		HAd:            dp.AdiabaticLabour,
		Eta:            dp.EtaAd,
		DT:             dp.TemperatureDrop,
		Pi:             dp.PiStag,
		P1:             dp.P1Stag,
		P3:             dp.P3Stag,
		T1:             dp.T1Stag,
		T3:             dp.T3Stag,
		ACrit1:         dp.ACrit1,
		ACrit3:         dp.ACrit3,

		CpAir: gases.CpMean(stage.Gas(), dp.T1Stag, dp.T3Stag, nodes.DefaultN),
		KAir:  gases.KMean(stage.Gas(), dp.T1Stag, dp.T3Stag, nodes.DefaultN),
		RAir:  stage.Gas().R(),

		Q1:       dp.Q1,
		MassRate: stage.MassRate(),

		CURel1: dp.InletTriangle.CU() / dp.UOut,
		CURel2: dp.MidTriangle.CU() / dp.UOut,
		CURel3: dp.OutletTriangle.CU() / dp.UOut,

		CARel1: dp.InletTriangle.CA() / dp.UOut,
		CARel2: dp.MidTriangle.CA() / dp.UOut,
		CARel3: dp.OutletTriangle.CA() / dp.UOut,

		Lambda1: dp.Lambda1,
		Lambda3: dp.Lambda3,

		Triangle1: dp.InletTriangle,
		Triangle2: dp.MidTriangle,
		Triangle3: dp.OutletTriangle,
	}
}

type CompressorStageDF struct {
	StatorDF   CompressorBladingGeomDF
	RotorDF    CompressorBladingGeomDF
	StageWidth float64
	RPM        float64

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

func NewCompressorBladingGeomDF(
	geom geometry.BladingGeometry,
	geomGen compressor.BladingGeometryGenerator,
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

package p2n

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free2n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
)

func NewBuilder(
	source schemes.TwoShaftsScheme,
	power, t0, p0,
	cRpm0, lambdaIn0,
	ctInletMeanDiameter, ctLambdaU0, ctStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0, etaM,
	precision, relaxCoef float64, iterLimit int,
) *Builder {
	return &Builder{
		Source:              source,
		Power:               power,
		T0:                  t0,
		P0:                  p0,
		CRpm0:               cRpm0,
		LambdaIn0:           lambdaIn0,
		CtInletMeanDiameter: ctInletMeanDiameter,
		CtLambdaU0:          ctLambdaU0,
		CtStageNum:          ctStageNum,
		FtInletMeanDiameter: ftInletMeanDiameter,
		FtLambdaU0:          ftLambdaU0,
		FtStageNum:          ftStageNum,
		PayloadRpm0:         payloadRpm0,
		EtaM:                etaM,
		Precision:           precision,
		RelaxCoef:           relaxCoef,
		IterLimit:           iterLimit,
	}
}

type Builder struct {
	Source schemes.TwoShaftsScheme
	Power  float64
	T0     float64
	P0     float64

	CRpm0 float64

	LambdaIn0 float64

	CtInletMeanDiameter float64
	CtLambdaU0          float64
	CtStageNum          float64

	FtInletMeanDiameter float64
	FtLambdaU0          float64
	FtStageNum          float64

	PayloadRpm0 float64

	EtaM float64

	Precision float64
	RelaxCoef float64
	IterLimit int
}

func (b *Builder) Build() free2n.DoubleShaftFreeScheme {
	return free2n.NewDoubleShaftFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.T0, b.P0, b.Source.Burner().TStagOut(),
		b.EtaM, b.BuildCompressor(), b.BuildCompressorPipe(),
		b.BuildBurner(), b.BuildCompressorTurbine(), b.BuildFreeTurbinePipe(),
		b.BuildFreeTurbine(), b.BuildFreeTurbinePipe(), b.BuildPayload(),
	)
}

func (b *Builder) BuildCompressor() constructive.ParametricCompressorNode {
	c := b.Source.Compressor()
	massRate0 := common.GetMassRate(b.Power, b.Source, c)
	return common.BuildCompressor(
		c,
		methodics.NewCompressorCharGen(
			c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, b.IterLimit,
		),
		b.CRpm0,
		common.GetMassRate(b.Power, b.Source, b.Source.Compressor()),
		b.Precision,
	)
}

func (b *Builder) BuildCompressorPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(1) // todo extract data about pipe from somewhere
}

func (b *Builder) BuildBurner() constructive.ParametricBurnerNode {
	burn := b.Source.Burner()
	return common.BuildBurner(
		burn, b.LambdaIn0,
		common.GetMassRate(b.Power, b.Source, burn),
		b.Precision,
	)
}

func (b *Builder) BuildCompressorTurbine() constructive.ParametricTurbineNode {
	ct := b.Source.TurboCascade().Turbine()
	return common.BuildTurbine(
		ct,
		methodics.NewKazandjanTurbineCharacteristic(),
		common.GetMassRate(b.Power, b.Source, ct),
		b.CtInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) BuildCTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.CompressorTurbinePipe().Sigma())
}

func (b *Builder) BuildFreeTurbine() constructive.ParametricTurbineNode {
	ft := b.Source.FreeTurbineBlock().FreeTurbine()
	return common.BuildTurbine(
		ft,
		methodics.NewKazandjanTurbineCharacteristic(),
		common.GetMassRate(b.Power, b.Source, ft),
		b.FtInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) BuildFreeTurbinePipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.FreeTurbineBlock().OutletPressureLoss().Sigma())
}

func (b *Builder) BuildPayload() constructive.Payload {
	return constructive.NewPayload(
		b.PayloadRpm0, b.Power, func(normRpm float64) float64 {
			//delta := normRpm - 1
			//return normRpm - delta * delta
			return normRpm * normRpm * normRpm // todo add smth more precise
		},
	)
}

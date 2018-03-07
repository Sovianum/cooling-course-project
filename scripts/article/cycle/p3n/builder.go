package p3n

import (
	"github.com/Sovianum/cooling-course-project/scripts/article/cycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/library/parametric/free3n"
	"github.com/Sovianum/turbocycle/library/schemes"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/methodics"
)

func NewBuilder(
	source schemes.ThreeShaftsScheme,
	power, t0, p0,
	lpcRpm0, hpcRpm0,
	lambdaIn0 float64,
	lptInletMeanDiameter, lptLambdaU0, lptStageNum,
	hptInletMeanDiameter, hptLambdaU0, hptStageNum,
	ftInletMeanDiameter, ftLambdaU0, ftStageNum,
	payloadRpm0,
	lpEtaM, hpEtaM,
	precision, relaxCoef float64, iterLimit int,
) *Builder {
	return &Builder{
		Source: source,
		Power:  power,
		T0:     t0,
		P0:     p0,

		LPCRpm0: lpcRpm0,
		HPCRpm0: hpcRpm0,

		LambdaIn0: lambdaIn0,

		LPTInletMeanDiameter: lptInletMeanDiameter,
		LPTLambdaU0:          lptLambdaU0,
		LPTStageNum:          lptStageNum,

		HPTInletMeanDiameter: hptInletMeanDiameter,
		HPTLambdaU0:          hptLambdaU0,
		HPTStageNum:          hptStageNum,

		FTInletMeanDiameter: ftInletMeanDiameter,
		FTLambdaU0:          ftLambdaU0,
		FTStageNum:          ftStageNum,

		PayloadRpm0: payloadRpm0,

		LPEtaM: lpEtaM,
		HPEtaM: hpEtaM,

		Precision: precision,
		RelaxCoef: relaxCoef,
		IterLimit: iterLimit,
	}
}

type Builder struct {
	Source schemes.ThreeShaftsScheme
	Power  float64
	T0     float64
	P0     float64

	LPCRpm0 float64
	HPCRpm0 float64

	LambdaIn0 float64

	LPTInletMeanDiameter float64
	LPTLambdaU0          float64
	LPTStageNum          float64

	HPTInletMeanDiameter float64
	HPTLambdaU0          float64
	HPTStageNum          float64

	FTInletMeanDiameter float64
	FTLambdaU0          float64
	FTStageNum          float64

	PayloadRpm0 float64

	LPEtaM float64
	HPEtaM float64

	Precision float64
	RelaxCoef float64
	IterLimit int
}

func (b *Builder) Build() free3n.ThreeShaftFreeScheme {
	return free3n.NewThreeShaftFreeScheme(
		b.Source.GasSource().GasOutput().GetState().Value().(gases.Gas),
		b.T0, b.P0, b.Source.MainBurner().TStagOut(),

		b.buildLPC(), b.buildLPCPipe(),
		b.buildLPT(), b.buildLPTPipe(),
		b.LPEtaM,

		b.buildHPC(), b.buildHPCPipe(),
		b.buildHPT(), b.buildHPTPipe(),
		b.HPEtaM,

		b.buildFT(), b.buildFTPipe(),
		b.BuildBurner(), b.buildPayload(),
	)
}

func (b *Builder) buildLPC() constructive.ParametricCompressorNode {
	c := b.Source.LPC()
	massRate0 := common.GetMassRate(b.Power, b.Source, c)
	charGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, b.IterLimit,
	)

	return common.BuildCompressor(
		c,
		charGen,
		b.LPCRpm0, common.GetMassRate(b.Power, b.Source, b.Source.LPC()),
		b.Precision,
	)
}

func (b *Builder) buildLPCPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.LPCPipe().Sigma())
}

func (b *Builder) buildHPC() constructive.ParametricCompressorNode {
	c := b.Source.LPC()
	massRate0 := common.GetMassRate(b.Power, b.Source, c)
	charGen := methodics.NewCompressorCharGen(
		c.PiStag(), c.Eta(), massRate0, precision, relaxCoef, b.IterLimit,
	)

	return common.BuildCompressor(
		c,
		charGen,
		b.HPCRpm0, common.GetMassRate(b.Power, b.Source, b.Source.HPC()),
		b.Precision,
	)
}

func (b *Builder) buildHPCPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.HPCPipe().Sigma())
}

func (b *Builder) BuildBurner() constructive.ParametricBurnerNode {
	burn := b.Source.MainBurner()
	return common.BuildBurner(
		burn, b.LambdaIn0,
		common.GetMassRate(b.Power, b.Source, burn),
		b.Precision,
	)
}

func (b *Builder) buildHPT() constructive.ParametricTurbineNode {
	return common.BuildTurbine(
		b.Source.HPT(),
		methodics.NewKazandjanTurbineCharacteristic(),
		common.GetMassRate(b.Power, b.Source, b.Source.HPT()),
		b.HPTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) buildHPTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.HPTPipe().Sigma())
}

func (b *Builder) buildLPT() constructive.ParametricTurbineNode {
	return common.BuildTurbine(
		b.Source.LPT(),
		methodics.NewKazandjanTurbineCharacteristic(),
		common.GetMassRate(b.Power, b.Source, b.Source.LPT()),
		b.LPTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) buildLPTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.LPTPipe().Sigma())
}

func (b *Builder) buildFT() constructive.ParametricTurbineNode {
	return common.BuildTurbine(
		b.Source.FT(),
		methodics.NewKazandjanTurbineCharacteristic(),
		common.GetMassRate(b.Power, b.Source, b.Source.FT()),
		b.FTInletMeanDiameter, b.Precision,
	)
}

func (b *Builder) buildFTPipe() constructive.PressureLossNode {
	return constructive.NewPressureLossNode(b.Source.FTBlock().OutletPressureLoss().Sigma())
}

func (b *Builder) buildPayload() constructive.Payload {
	return constructive.NewPayload(
		b.PayloadRpm0, b.Power, func(normRpm float64) float64 {
			//delta := normRpm - 1
			//return normRpm - delta * delta
			return normRpm * normRpm * normRpm // todo add smth more precise
		},
	)
}
